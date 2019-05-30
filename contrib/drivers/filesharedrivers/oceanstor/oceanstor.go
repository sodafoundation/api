// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oceanstor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	model "github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

var once sync.Once

func (d *Driver) Setup() error {
	var err error

	once.Do(func() {
		d.IniConf()
		cli, err := newRestCommon(d.Config)
		if err == nil {
			d.Client = cli
			log.Info("get oceanstor client successfully")
		}
	})

	if err != nil {
		msg := fmt.Sprintf("get new client failed: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	return nil
}

func (d *Driver) IniConf() {
	path := config.CONF.OsdsDock.Backends.HuaweiOceanstor.ConfigPath
	if path == "" {
		path = DefaultConfPath
	}

	conf := &Config{}
	d.Config = conf
	Parse(conf, path)
}

func (d *Driver) Unset() error {

	if err := d.logout(); err != nil {
		msg := fmt.Sprintf("logout failed: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	return nil
}
func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	fsName := opt.GetName()
	size := opt.GetSize()
	prf := opt.GetProfile()
	poolID := opt.GetPoolName()
	shareProto := ""

	err := d.parameterCheck(poolID, prf, size, &fsName, &shareProto)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// create file system
	fs, err := d.createFileSystemIfNotExist(fsName, poolID, size)
	if err != nil {
		msg := fmt.Sprintf("create file system failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	shareDriver := NewProtocol(shareProto, d.Client)
	// create file share if not exist
	shareID, err := d.createShareIfNotExist(fsName, fs.ID, shareDriver)
	if err != nil {
		return nil, err
	}

	// get location
	location, err := d.getShareLocation(fsName, shareDriver)
	if err != nil {
		msg := fmt.Sprintf("get share location failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	share := &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Protocols:        []string{shareProto},
		Description:      opt.GetDescription(),
		Size:             size,
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           poolID,
		ExportLocations:  location,
		Metadata:         map[string]string{FileShareName: fsName, FileShareID: shareID},
	}
	return share, nil
}

func (d *Driver) parameterCheck(poolID, prf string, size int64, fsName, shareProto *string) error {
	// Parameter check
	if poolID == "" {
		msg := "pool id cannot be empty"
		log.Error(msg)
		return errors.New(msg)
	}

	if *fsName == "" {
		log.Infof("use default file system name %s", defaultFileSystem)
		*fsName = defaultFileSystem
	}

	proto, err := d.GetProtoFromProfile(prf)
	if err != nil {
		return err
	}

	if !checkProtocol(proto) {
		return fmt.Errorf("%s protocol is not supported, support is NFS and CIFS", proto)
	}

	*shareProto = proto

	if size == 0 {
		return errors.New("size must be greater than 0")
	}

	return nil
}

func (d *Driver) createShareIfNotExist(fsName, fsID string, shareDriver Protocol) (string, error) {
	sharePath := getSharePath(fsName)
	share, err := shareDriver.getShare(fsName)
	if err != nil {
		return "", fmt.Errorf("get share %s failed: %v", sharePath, err)
	}

	if share != nil {
		log.Infof("share %s already exist", sharePath)
		return "", nil
	}

	share, err = shareDriver.createShare(fsName, fsID)
	if err != nil {
		shareDriver.deleteShare(fsName)
		return "", fmt.Errorf("create share %s failed: %v", sharePath, err)
	}

	log.Infof("create share %s successfully", sharePath)
	return shareDriver.getShareID(share), nil
}

func (d *Driver) getShareLocation(fsName string, shareDriver Protocol) ([]string, error) {
	logicalPortList, err := d.getAllLogicalPort()
	if err != nil {
		return nil, err
	}

	location, err := d.getLocationPath(fsName, logicalPortList, shareDriver)
	if err != nil {
		return nil, err
	}

	return location, nil
}

// createFileSystemIfNotExist
func (d *Driver) createFileSystemIfNotExist(fsName, poolID string, size int64) (*FileSystemData, error) {
	fsList, err := d.getFileSystemByName(fsName)
	if err != nil {
		return nil, fmt.Errorf("get filesystem %s by name failed: %v", fsName, err)
	}

	if len(fsList) == 0 {
		fs, err := d.createFileSystem(fsName, poolID, size)
		if err != nil {
			errDelete := d.deleteFileSystem(fsName, poolID)
			fmt.Println("deleteFileSystem", errDelete)
			return nil, fmt.Errorf("create file system %s failed, %v", fsName, err)
		}

		err = d.checkFsStatus(fs, poolID)
		if err != nil {
			return nil, err
		}

		log.Infof("create filesystem %s successfully", fsName)
		return fs, nil
	}

	log.Infof("filesystem %s already exist", fsName)
	return &fsList[0], nil
}

func (d *Driver) getLocationPath(sharePath string, logicalPortList []LogicalPortData, shareDriver Protocol) ([]string, error) {
	if len(logicalPortList) == 0 {
		return nil, errors.New("cannot find file share server end logical ip")
	}

	var location []string

	for _, port := range logicalPortList {
		location = append(location, shareDriver.getLocation(sharePath, port.IpAddr))
	}

	return location, nil
}

func (d *Driver) checkFsStatus(fs *FileSystemData, poolID string) error {
	ticker := time.NewTicker(3 * time.Second)
	timeout := time.After(1 * time.Minute)
	var fsStable *FileSystemData
	var err error
	for {
		select {
		case <-ticker.C:
			fsStable, err = d.getFileSystem(fs.ID)
			if err != nil {
				d.deleteFileSystem(fs.Name, poolID)
				return fmt.Errorf("check file system status failed: %v", err)
			}

			if fsStable.HealthStatus == StatusFSHealth && fsStable.RunningStatus == StatusFSRunning {
				return nil
			}

		case <-timeout:
			d.deleteFileSystem(fs.Name, poolID)
			return fmt.Errorf("timeout occured waiting for checking file system status %s or invalid status health:%s, running:%s", fsStable.ID, fsStable.HealthStatus, fsStable.RunningStatus)
		}
	}
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec
	sp, err := d.ListStoragePools()
	if err != nil {
		msg := fmt.Sprintf("list pools from storage failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	c := d.Config
	for _, p := range sp {
		if _, ok := c.Pool[p.Name]; !ok {
			continue
		}
		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, d.Uri, p.Id)

		userTotalCapacity, _ := strconv.ParseInt(p.UserTotalCapacity, 10, 64)
		userFreeCapacity, _ := strconv.ParseInt(p.UserFreeCapacity, 10, 64)

		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             p.Id,
			TotalCapacity:    Sector2Gb(userTotalCapacity),
			FreeCapacity:     Sector2Gb(userFreeCapacity),
			StorageType:      c.Pool[p.Name].StorageType,
			Extras:           c.Pool[p.Name].Extras,
			AvailabilityZone: c.Pool[p.Name].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}

	if len(pols) == 0 {
		msg := "pools in configuration file not found"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	log.Info("list pools successfully")

	return pols, nil
}

func (d *Driver) getFSInfo(fsName string) (*FileSystemData, error) {
	fsList, err := d.getFileSystemByName(fsName)
	if err != nil {
		return nil, fmt.Errorf("get filesystem %s by name failed: %v", fsName, err)
	}

	if len(fsList) == 0 {
		return nil, fmt.Errorf("filesystem %s does not exist", fsName)
	}

	return &fsList[0], nil
}

func (d *Driver) GetProtoFromProfile(prf string) (string, error) {
	if prf == "" {
		msg := "profile cannot be empty"
		return "", errors.New(msg)
	}

	log.V(5).Infof("file share profile is %s", prf)
	profile := &model.ProfileSpec{}
	err := json.Unmarshal([]byte(prf), profile)
	if err != nil {
		msg := fmt.Sprintf("unmarshal profile failed: %v", err)
		return "", errors.New(msg)
	}

	shareProto := profile.ProvisioningProperties.IOConnectivity.AccessProtocol
	if shareProto == "" {
		msg := "file share protocol cannot be empty"
		return "", errors.New(msg)
	}

	return shareProto, nil
}

func (d *Driver) DeleteFileShare(opt *pb.DeleteFileShareOpts) (*model.FileShareSpec, error) {
	shareProto, err := d.GetProtoFromProfile(opt.GetProfile())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	meta := opt.GetMetadata()
	if meta == nil || (meta != nil && meta[FileShareName] == "" && meta[FileShareID] == "") {
		msg := "cannot get file share name and id"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	fsName := meta[FileShareName]
	shareID := meta[FileShareID]

	shareDriver := NewProtocol(shareProto, d.Client)

	sharePath := getSharePath(fsName)
	if err := shareDriver.deleteShare(shareID); err != nil {
		msg := fmt.Sprintf("delete file share %s failed: %v", sharePath, err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	log.Infof("delete share %s successfully", sharePath)

	err = d.DeleteFileSystem(fsName)
	if err != nil {
		msg := fmt.Sprintf("delete filesystem %s failed: %v", fsName, err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	log.Infof("delete file system %s successfully", fsName)

	return nil, nil
}

func (d *Driver) DeleteFileSystem(fsName string) error {
	fs, err := d.getFSInfo(fsName)
	if err != nil {
		return err
	}

	err = d.deleteFS(fs.ID)
	if err != nil {
		return err
	}

	log.Infof("delete filesystem %s successfully", fs.ID)
	return nil
}

func (d *Driver) CreateFileShareSnapshot(opt *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	snapID := opt.GetId()
	if snapID == "" {
		msg := "snapshot id cannot be empty"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	meta := opt.GetMetadata()

	if meta == nil || (meta != nil && meta[FileShareName] == "" && meta[FileShareID] == "") {
		msg := "cannot get file share name and id"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	fsName := meta[FileShareName]

	fs, err := d.getFSInfo(fsName)
	if err != nil {
		msg := fmt.Sprintf("get file system %s failed: %v", fsName, err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	snapName := EncodeName(snapID)

	fsSnapshot, err := d.createSnapshot(fs.ID, snapName)
	if err != nil {
		msg := fmt.Sprintf("create filesystem snapshot failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	return &model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:     opt.GetName(),
		Metadata: map[string]string{FileShareSnapshotID: fsSnapshot.ID},
	}, nil
}

func (d *Driver) DeleteFileShareSnapshot(opt *pb.DeleteFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	//opt.GetContext()
	//TODO real snapshot id is needed
	snapID, err := d.getSnapshotID(opt.GetId())
	if err != nil {
		return nil, err
	}

	err = d.deleteFSSnapshot(snapID)
	if err != nil {
		msg := fmt.Sprintf("delete filesystem snapshot %s failed, %v", snapID, err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	log.Infof("delete file share snapshot %s successfully", snapID)

	return nil, nil
}

func (d *Driver) getSnapshotID(snapID string) (string, error) {
	snapName := EncodeName(snapID)

	//TODO change listsnapshots
	snaps, err := d.listSnapshots()
	if err != nil {
		msg := fmt.Sprintf("list share snapshots failed: %v", err)
		log.Error(msg)
		return "", errors.New(msg)
	}

	for _, snap := range snaps {
		if snap.Name == snapName {
			return snap.ID, nil
		}
	}

	return "", fmt.Errorf("cannot find snapshot %s", snapID)
}

// AllowAccess allow access to the share
func (d *Driver) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error) {
	accessLevels := opt.GetAccessCapability()
	accessToShares := opt.GetAccessTo()
	fsName := opt.GetName()
	profile := opt.Profile
	if len(accessLevels) == 0 || len(accessToShares) == 0 {
		msg := "access level and access to cannot be empty"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	// Take only the first value
	accessLevel := accessLevels[0]
	accessTo := accessToShares[0]

	//shareProto, shareName is needed
	if !checkAccessLevel(accessLevel) {
		return fmt.Errorf("access level %s is unsupported", accessLevel)
	}

	pattern := "\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}"
	matchIp, _ := regexp.MatchString(pattern, accessTo)
	if !matchIp {
		return nil, fmt.Errorf("ip %s is invalid", accessTo)
	}

	if shareProto == NFSProto {
		if accessLevel == AccessLevelRW {
			accessLevel = AccessNFSRw
		} else {
			accessLevel = AccessNFSRo
		}
	}

	if shareProto == CIFSProto {
		if accessLevel == AccessLevelRW {
			accessLevel = AccessCIFSFullControl
		} else {
			accessLevel = AccessCIFSRo
		}
	}

	share, err := d.getFileShare(shareProto, shareName)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if share == nil {
		return nil, fmt.Errorf("share %s is not exist", getSharePath(shareName))
	}

	shareID := d.getShareID(shareProto, share)

	err = d.createAccessIfNotExist(shareID, accessTo, shareProto, accessLevel)
	if err != nil {
		msg := fmt.Sprintf("allow access %s to %s failed %v", accessTo, shareName, err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	return nil
}

func (d *Driver) createAccessIfNotExist(shareID, accessTo, shareProto, accessLevel string) error {
	// Check if access already exists
	accessID, err := d.getAccessFromShare(shareID, accessTo, shareProto)
	if err != nil {
		return err
	}

	if accessID == "" {
		return d.allowAccessToShare(shareID, accessTo, shareProto, accessLevel)
	}

	return nil
}

func (d *Driver) allowAccessToShare(shareID, accessTo, shareProto, accessLevel string) error {
	switch shareProto {
	case NFS:
		if _, err := d.allowNFSAccess(shareID, accessTo, accessLevel); err != nil {
			return err
		}
	case CIFS:
		if _, err := d.allowCIFSAccess(shareID, accessTo, accessLevel); err != nil {
			return err
		}
	}

	return nil
}
