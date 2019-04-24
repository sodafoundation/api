// Copyright (c) 2019 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package oceanstor

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	model "github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	//	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

type Driver struct {
	conf   *Config
	client *Cli
}

type AuthOptions struct {
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Uri             string `yaml:"uri"`
	PwdEncrypter    string `yaml:"PwdEncrypter,omitempty"`
	EnableEncrypted bool   `yaml:"EnableEncrypted,omitempty"`
}

type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

const (
	KFileshareName = "OceanFileshareName"
	KFileshareID   = "OceanFileshareID"
)

func (d *Driver) Setup() error {
	if d.client != nil {
		return nil
	}
	conf := &Config{}

	d.conf = conf

	path := ""
	if path == "" {
		path = "./testdata/oceanstor.yaml"
	}

	Parse(conf, path)

	cli, err := newRestCommon(conf)
	if err != nil {
		log.Errorf("Get new client failed, %v", err)
		return err
	}

	d.client = cli

	log.Info("Get new client success")
	return nil
}

func (d *Driver) Unset() error {
	return nil
}

func EncodeName(id string) string {
	return NamePrefix + "-" + id
}

func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	// need tenantID, shareProto parameters
	var fs FileSystemData

	fsName := defaultFileSystem

	// Parameter check
	poolID := opt.GetPoolName()
	tenantID := ""
	shareProto := ""

	if poolID == "" {
		msg := "pool id cannot be empty"
		log.Error(msg)
		return nil, errors.New(msg)
	}

	if shareProto != CIFS && shareProto != NFS {
		return nil, errors.New(shareProto + " protocol is not supported, support is NFS and CIFS")
	}

	// create filesystem if not exist
	fsList, err := d.client.getFileSystemByName(fsName)
	if err != nil {
		msg := fmt.Sprintf("get filesystem %s by name failed, %v", fsName, err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	if len(fsList.Data) == 0 {
		fsObj, err := d.client.createFileSystem(fsName, poolID, tenantID)
		if err != nil {
			//NOTE
			//d.client.deleteFileSystem(name)
			msg := fmt.Sprintf("create file system %s failed, %v", fsName, err)
			log.Error(msg)
			return nil, errors.New(msg)
		}

		err = d.checkFsStatus(fsObj)
		if err != nil {
			return nil, err
		}

		fs = fsObj.FileSystemData
		log.Infof("filesystem %s creation success", fsName)
	} else {
		fs = fsList.Data[0]
		log.Infof("filesystem %s already exist", fsName)
	}

	var NFSShare *NFSShareData
	var CIFSShare *CIFSShareData
	var shareID, shareName string
	// create file share based on protocol
	if shareProto == NFS {
		// create nfs share if not exist
		NFSShare, err = d.client.getNFSShare(fsName)
		if err != nil {
			msg := fmt.Sprintf("get nfs share /%s/ failed, %v", fsName, err)
			log.Error(msg)
			return nil, errors.New(msg)
		}
		fmt.Println(NFSShare, err)
		if NFSShare == nil {
			NFSShare, err = d.client.createNFSShare(fsName, fs.ID)
			if err != nil {
				//d.client.deleteShare()
				msg := fmt.Sprintf("create nfs share /%s/ failed, %v", fsName, err)
				log.Error(msg)
				return nil, errors.New(msg)
			}
			log.Infof("nfs share /%s/ creation success", fsName)
		} else {
			log.Infof("nfs share /%s/ already exist", fsName)
		}
		shareID = NFSShare.ID
		shareName = NFSShare.Name
	}

	if shareProto == CIFS {
		// create cifs share if not exist
		CIFSShare, err = d.client.getCIFSShare(fsName)
		if err != nil {
			msg := fmt.Sprintf("get cifs share /%s/ failed, %v", fsName, err)
			log.Error(msg)
			return nil, errors.New(msg)
		}

		if CIFSShare == nil {
			CIFSShare, err = d.client.createCIFSShare(fsName, fs.ID)
			if err != nil {
				//d.client.deleteShare()
				msg := fmt.Sprintf("create cifs share /%s/ failed, %v", fsName, err)
				log.Error(msg)
				return nil, errors.New(msg)
			}

			log.Infof("cifs share /%s/ creation success", fsName)
		} else {
			log.Infof("cifs share /%s/ already exist", fsName)
		}
		shareID = CIFSShare.ID
		shareName = CIFSShare.Name
	}

	u, _ := url.Parse(d.conf.Uri)
	ip := strings.Split(u.Host, ":")[0]
	location := d.getLocationPath(fsName, shareProto, ip)

	return &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		TenantId:         "123",
		PoolId:           poolID,
		ExportLocations:  []string{location},
		Metadata: map[string]string{
			KFileshareName: shareName,
			KFileshareID:   shareID,
		},
	}, nil
}

func (d *Driver) getLocationPath(sharePath, shareProto, ip string) string {
	var location string
	if shareProto == NFS {
		if isIPv6(ip) {
			location = fmt.Sprintf("[%s]:/%s", ip, strings.Replace(sharePath, "-", "_", -1))
		} else {
			location = fmt.Sprintf("%s:/%s", ip, strings.Replace(sharePath, "-", "_", -1))
		}
	}
	if shareProto == CIFS {
		location = fmt.Sprintf("\\\\%s\\%s", ip, strings.Replace(sharePath, "-", "_", -1))
	}

	return location
}

func isIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

func (d *Driver) checkFsStatus(fs *FileSystem) error {
	ticker := time.NewTicker(3 * time.Second)
	timeout := time.After(1 * time.Minute)
	var fsStable *FileSystem
	var err error
	for {
		select {
		case <-ticker.C:
			fsStable, err = d.client.getFileSystem(fs.ID)
			if err != nil {
				//NOTE
				//d.client.deleteFileSystem(name)
				msg := fmt.Sprintf("check file system status failed, %v", err)
				log.Error(msg)
				return errors.New(msg)
			}

			if fsStable.HealthStatus == StatusFSHealth && fsStable.RunningStatus == StatusFSRunning {
				return nil
			}

		case <-timeout:
			//NOTE
			//d.client.deleteFileSystem(name)
			msg := fmt.Sprintf("timeout occured waiting for checking file system status %s or invalid status health:%s, running:%s", fsStable.ID, fsStable.HealthStatus, fsStable.RunningStatus)
			log.Errorf(msg)
			return errors.New(msg)
		}
	}
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	var pols []*model.StoragePoolSpec
	sp, err := d.client.ListStoragePools()
	if err != nil {
		msg := fmt.Sprintf("list pools from storage failed, %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	c := d.conf
	for _, p := range sp {
		if _, ok := c.Pool[p.Name]; !ok {
			continue
		}
		host, _ := os.Hostname()
		name := fmt.Sprintf("%s:%s:%s", host, d.conf.Uri, p.Id)
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name:             p.Id,
			TotalCapacity:    Sector2Gb(p.UserTotalCapacity),
			FreeCapacity:     Sector2Gb(p.UserFreeCapacity),
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

func Sector2Gb(sec string) int64 {
	size, err := strconv.ParseInt(sec, 10, 64)
	if err != nil {
		log.Error("Convert capacity from string to number failed, error:", err)
		return 0
	}
	return size * 512 / UnitGi
}

func (d *Driver) DeleteShare(shareID, shareProto, fsID string) error {

	if shareProto == NFS {
		err := d.client.deleteNFSShare(shareID)
		if err != nil {
			msg := fmt.Sprintf("delete nfs share %s failed, %v", shareID, err)
			log.Error(msg)
			return errors.New(msg)
		}

		log.Infof("delete nfs share %s successfully", shareID)
	}

	if shareProto == CIFS {
		err := d.client.deleteCIFSShare(shareID)
		if err != nil {
			msg := fmt.Sprintf("delete cifs share %s failed, %v", shareID, err)
			log.Error(msg)
			return errors.New(msg)
		}

		log.Infof("delete cifs share %s successfully", shareID)
	}

	err := d.client.deleteFS(fsID)
	if err != nil {
		msg := fmt.Sprintf("delete filesystem %s failed, %v", fsID, err)
		log.Error(msg)
		return errors.New(msg)
	}

	log.Infof("delete filesystem %s successfully", fsID)
	return nil
}

func (d *Driver) ListAllShares(shareID, shareProto, fsID string) error {
	return nil
}

func (d *Driver) CreateSnapshotFromShare(snapID, shareProto, shareID string) error {
	var fsID string

	if shareProto == NFS {
		nfsShare, _ := d.client.getNFSShareByID(shareID, shareProto)
		if nfsShare == nil || nfsShare.FSID == "" {
			msg := "can not create snapshot due to FS not exist"
			log.Error(msg)
			return errors.New(msg)
		}
		fsID = nfsShare.FSID
	}

	if shareProto == CIFS {
		cifsShare, _ := d.client.getCIFSShareByID(shareID, shareProto)
		if cifsShare == nil || cifsShare.FSID == "" {
			msg := "can not create snapshot due to FS not exist"
			log.Error(msg)
			return errors.New(msg)
		}
		fsID = cifsShare.FSID
	}

	snapName := "share_snapshot_" + snapID

	fsSnapshot, err := d.client.createSnapshot(fsID, snapName)
	if err != nil {
		msg := fmt.Sprintf("create filesystem snapshot failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	fmt.Printf("ss %+v\n", fsSnapshot)
	return nil
}

func (d *Driver) ListAllSnapshots() error {
	fsList, err := d.client.getAllFilesystem()
	if err != nil {
		msg := fmt.Sprintf("list filesystem failed, %v", err)
		log.Error(msg)
		return errors.New(msg)
	}

	var fsSnapshotList []*FSSnapshotData

	for _, v := range fsList.Data {
		fsSnapshots, err := d.client.listSnapshots(v.ID)
		if err != nil {
			msg := fmt.Sprintf("list filesystem snapshots failed, %v", err)
			log.Error(msg)
			return errors.New(msg)
		}

		for _, snap := range fsSnapshots.Data {
			fsSnapshotList = append(fsSnapshotList, &snap)
		}
	}

	for _, v := range fsSnapshotList {
		fmt.Printf("dddd %+v\n", v)
	}

	return nil
}

func (d *Driver) DeleteFSSnapshot(snapID string) error {
	err := d.client.deleteFSSnapshot(snapID)
	if err != nil {
		msg := fmt.Sprintf("delete filesystem snapshot %s failed, %v", snapID, err)
		log.Error(msg)
		return errors.New(msg)
	}

	return nil
}

func (d *Driver) ShowFSSnapshot(snapID string) error {
	snap, err := d.client.showFSSnapshot(snapID)
	if err != nil {
		msg := fmt.Sprintf("show filesystem snapshot %s failed, %v", snapID, err)
		log.Error(msg)
		return errors.New(msg)
	}

	fmt.Printf("ShowFSSnapshot %+v\n", snap)
	return nil
}
