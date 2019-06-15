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
	"strconv"
	"time"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	model "github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

func (d *Driver) Setup() error {
	if d.Client != nil {
		return nil
	}

	var err error

	d.InitConf()
	cli, err := newRestCommon(d.Config)
	if err != nil {
		msg := fmt.Sprintf("get new client failed: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	d.Client = cli
	log.Info("get oceanstor client successfully")

	return nil
}

func (d *Driver) InitConf() {
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
		return fmt.Errorf("%s protocol is not supported, support is %s and %s", proto, NFSProto, CIFSProto)
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
			d.DeleteFileSystem(fsName)
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
				d.DeleteFileSystem(fs.Name)
				return fmt.Errorf("check file system status failed: %v", err)
			}

			if fsStable.HealthStatus == StatusFSHealth && fsStable.RunningStatus == StatusFSRunning {
				return nil
			}

		case <-timeout:
			d.DeleteFileSystem(fs.Name)
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

func (d *Driver) DeleteFileShare(opt *pb.DeleteFileShareOpts) error {
	shareProto, err := d.GetProtoFromProfile(opt.GetProfile())
	if err != nil {
		log.Error(err.Error())
		return err
	}

	meta := opt.GetMetadata()
	if meta == nil || (meta != nil && meta[FileShareName] == "" && meta[FileShareID] == "") {
		msg := "cannot get file share name and id"
		log.Error(msg)
		return errors.New(msg)
	}

	fsName := meta[FileShareName]
	shareID := meta[FileShareID]

	shareDriver := NewProtocol(shareProto, d.Client)

	sharePath := getSharePath(fsName)
	if err := shareDriver.deleteShare(shareID); err != nil {
		msg := fmt.Sprintf("delete file share %s failed: %v", sharePath, err)
		log.Error(msg)
		return errors.New(msg)
	}

	log.Infof("delete share %s successfully", sharePath)

	if err := d.DeleteFileSystem(fsName); err != nil {
		msg := fmt.Sprintf("delete filesystem %s failed: %v", fsName, err)
		log.Error(msg)
		return errors.New(msg)
	}

	log.Infof("delete file system %s successfully", fsName)

	return nil
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

	if fs == nil {
		msg := fmt.Sprintf("%s does not exist", fsName)
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

	snapSize, _ := strconv.ParseInt(fsSnapshot.Capacity, 10, 64)

	return &model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:         snapName,
		Description:  opt.GetDescription(),
		SnapshotSize: snapSize,
		Metadata:     map[string]string{FileShareSnapshotID: fsSnapshot.ID},
	}, nil
}

func (d *Driver) DeleteFileShareSnapshot(opt *pb.DeleteFileShareSnapshotOpts) error {
	meta := opt.GetMetadata()
	if meta == nil || (meta != nil && meta[FileShareSnapshotID] == "") {
		msg := "cannot get file share snapshot id"
		log.Error(msg)
		return errors.New(msg)
	}

	snapID := meta[FileShareSnapshotID]

	err := d.deleteFSSnapshot(snapID)
	if err != nil {
		msg := fmt.Sprintf("delete filesystem snapshot %s failed, %v", snapID, err)
		log.Error(msg)
		return errors.New(msg)
	}

	log.Infof("delete file share snapshot %s successfully", snapID)

	return nil
}

func (d *Driver) getAccessLevel(accessLevels []string, shareProto string) (string, error) {
	var accessLevel string

	if accessLevels == nil || (accessLevels != nil && len(accessLevels) == 0) {
		return "", errors.New("access level cannot be empty")
	}

	supportAccessLevels := []string{AccessLevelRead, AccessLevelWrite}

	if len(accessLevels) > len(supportAccessLevels) {
		return "", errors.New("invalid access level")
	}

	accessLevel = "ro"
	for _, v := range accessLevels {
		if !utils.Contained(v, supportAccessLevels) {
			return "", errors.New("only read only or read write access level are supported")
		}
		if v == AccessLevelWrite {
			accessLevel = "rw"
		}
	}

	shareDriver := NewProtocol(shareProto, d.Client)
	return shareDriver.getAccessLevel(accessLevel), nil
}

func (d *Driver) CreateFileShareAclParamCheck(opt *pb.CreateFileShareAclOpts) (string, string, string, string, error) {
	log.V(5).Infof("create file share access client parameters %#v", opt)
	meta := opt.GetMetadata()

	if meta == nil || (meta != nil && meta[FileShareName] == "" && meta[FileShareID] == "") {
		msg := "cannot get file share name and id"
		log.Error(msg)
		return "", "", "", "", errors.New(msg)
	}

	fsName := meta[FileShareName]
	if fsName == "" {
		return "", "", "", "", errors.New("fileshare name cannot be empty")
	}

	shareProto, err := d.GetProtoFromProfile(opt.Profile)
	if err != nil {
		return "", "", "", "", err
	}

	if !checkProtocol(shareProto) {
		return "", "", "", "", fmt.Errorf("%s protocol is not supported, support is NFS and CIFS", shareProto)
	}

	accessLevels := opt.GetAccessCapability()

	accessLevel, err := d.getAccessLevel(accessLevels, shareProto)
	if err != nil {
		return "", "", "", "", err
	}

	accessType := opt.Type
	if !checkAccessType(accessType) {
		return "", "", "", "", fmt.Errorf("only access type %s and %s are supported", AccessTypeUser, AccessTypeIp)
	}
	if shareProto == CIFSProto && accessType != AccessTypeUser {
		return "", "", "", "", errors.New("only USER access type is allowed for CIFS shares")
	}

	accessTo := opt.GetAccessTo()
	if accessTo == "" {
		return "", "", "", "", errors.New("access client cannot be empty")
	}

	if shareProto == NFSProto {
		if accessType == AccessTypeUser {
			accessTo += "@"
		} else {
			accessTo = "*"
		}
	}

	return fsName, shareProto, accessLevel, accessTo, nil
}

// AllowAccess allow access to the share
func (d *Driver) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (*model.FileShareAclSpec, error) {
	shareName, shareProto, accessLevel, accessTo, err := d.CreateFileShareAclParamCheck(opt)
	if err != nil {
		msg := fmt.Sprintf("create fileshare access client failed: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	shareDriver := NewProtocol(shareProto, d.Client)

	share, err := shareDriver.getShare(shareName)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if share == nil {
		return nil, fmt.Errorf("share %s does not exist", shareName)
	}

	shareID := shareDriver.getShareID(share)

	err = d.createAccessIfNotExist(shareID, accessTo, shareProto, accessLevel, shareDriver)
	if err != nil {
		msg := fmt.Sprintf("allow access %s to %s failed %v", accessTo, shareName, err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	shareAccess := &model.FileShareAclSpec{
		BaseModel: &model.BaseModel{
			Id: opt.Id,
		},
		AccessTo: accessTo,
		Metadata: map[string]string{FileShareName: shareName},
	}

	return shareAccess, nil
}

func (d *Driver) createAccessIfNotExist(shareID, accessTo, shareProto, accessLevel string, shareDriver Protocol) error {
	// Check if access already exists
	accessID, err := d.getAccessFromShare(shareID, accessTo, shareProto)
	if err != nil {
		return err
	}

	if accessID != "" {
		log.Infof("fileshare access %s already exists", accessID)
		return nil
	}

	if _, err := shareDriver.allowAccess(shareID, accessTo, accessLevel); err != nil {
		return err
	}

	log.Infof("create fileshare access successfully")

	return nil
}

func (d *Driver) DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error {
	shareName, shareProto, accessTo, err := d.DeleteFileShareAclParamCheck(opt)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	accessTo = "*"

	shareDriver := NewProtocol(shareProto, d.Client)

	share, err := shareDriver.getShare(shareName)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if share == nil {
		msg := fmt.Sprintf("share %s does not exist", shareName)
		log.Error(msg)
		return errors.New(msg)
	}

	shareID := shareDriver.getShareID(share)

	accessID, err := d.getAccessFromShare(shareID, accessTo, shareProto)
	if err != nil {
		msg := fmt.Sprintf("get access from share failed: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	if accessID == "" {
		msg := fmt.Sprintf("can not get access id from share %s", shareName)
		log.Error(msg)
		return errors.New(msg)
	}

	if err := d.removeAccessFromShare(accessID, shareProto); err != nil {
		msg := fmt.Sprintf("remove access from share failed: %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	return nil
}

func (d *Driver) DeleteFileShareAclParamCheck(opt *pb.DeleteFileShareAclOpts) (string, string, string, error) {
	meta := opt.GetMetadata()
	if meta == nil || (meta != nil && meta[FileShareName] == "") {
		return "", "", "", errors.New("fileshare name cannot be empty when deleting file share access client")
	}

	fsName := meta[FileShareName]

	shareProto, err := d.GetProtoFromProfile(opt.Profile)
	if err != nil {
		return "", "", "", err
	}

	if !checkProtocol(shareProto) {
		return "", "", "", fmt.Errorf("%s protocol is not supported, support is NFS and CIFS", shareProto)
	}

	accessType := opt.Type
	if !checkAccessType(accessType) {
		return "", "", "", fmt.Errorf("only access type %s and %s are supported", AccessTypeUser, AccessTypeIp)
	}
	if shareProto == CIFSProto && accessType != AccessTypeUser {
		return "", "", "", fmt.Errorf("only USER access type is allowed for CIFS shares")
	}

	accessTo := opt.GetAccessTo()
	if accessTo == "" {
		return "", "", "", errors.New("cannot find access client")
	}

	return fsName, shareProto, accessTo, nil
}
