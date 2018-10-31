// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements ceph driver for OpenSDS. Ceph driver will pass these
operation requests about volume to go-ceph module.
*/

package ceph

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

const (
	opensdsPrefix   = "opensds-"
	sizeShiftBit    = 30
	defaultConfPath = "/etc/opensds/driver/ceph.yaml"
	defaultAZ       = "default"
)

const (
	globalSize = iota
	globalAvail
	globalRawUsed
	globalRawUsedPercentage
)

const (
	poolName = iota
	poolId
	poolUsed
	poolUsedPer
	poolMaxAvail
	poolObjects
)

const (
	poolType = iota
	poolTypeSize
	poolCrushRuleset
)

type CephConfig struct {
	ConfigFile string                    `yaml:"configFile,omitempty"`
	Pool       map[string]PoolProperties `yaml:"pool,flow"`
}

func execCmd(cmd string) (string, error) {
	ret, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return string(ret[:len(ret)-1]), nil
}

type Driver struct {
	conn  *rados.Conn
	ioctx *rados.IOContext
	conf  *CephConfig
}

func (d *Driver) Setup() error {
	d.conf = &CephConfig{ConfigFile: "/etc/ceph/ceph.conf"}
	p := config.CONF.OsdsDock.Backends.Ceph.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	_, err := Parse(d.conf, p)
	return err
}

func (d *Driver) Unset() error { return nil }

func (d *Driver) initConn() error {
	conn, err := rados.NewConn()
	if err != nil {
		log.Error("New connect failed:", err)
		return err
	}

	if err = conn.ReadConfigFile(d.conf.ConfigFile); err != nil {
		log.Error("Read config file failed:", err)
		return err
	}
	if err = conn.Connect(); err != nil {
		log.Error("Connect failed:", err)
		return err
	}

	d.conn = conn
	return nil
}

func (d *Driver) init(poolName string) error {
	err := d.initConn()
	if err != nil {
		return err
	}
	d.ioctx, err = d.conn.OpenIOContext(poolName)
	if err != nil {
		log.Error("Open IO context failed:", err)
		return err
	}
	return nil
}

func (d *Driver) destroy() {
	defer d.conn.Shutdown()
	defer d.ioctx.Destroy()
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	size := opt.GetSize()
	id := opt.GetId()
	name := opensdsPrefix + id
	if err := d.init(opt.GetPoolName()); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroy()

	_, err := rbd.Create(d.ioctx, name, uint64(size)<<sizeShiftBit, 20)
	if err != nil {
		log.Errorf("Create rbd image (%s) failed, (%v)", name, err)
		return nil, err
	}

	log.Infof("Create volume %s (%s) success.", opt.GetName(), id)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: id,
		},
		Name:             opt.GetName(),
		Size:             size,
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			"poolName": opt.GetPoolName(),
		},
	}, nil
}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	if err := d.init(opt.GetPoolName()); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroy()

	img, _, err := d.getImage(opt.GetId())
	if err != nil {
		log.Error("When get image:", err)
		return nil, err
	}

	if err = img.Open(); err != nil {
		log.Error("When open image:", err)
		return nil, err
	}
	defer img.Close()

	size := opt.GetSize()
	if err = img.Resize(uint64(size) << sizeShiftBit); err != nil {
		log.Error("When resize image:", err)
		return nil, err
	}
	log.Info("Resize image success, volume id =", opt.GetId())

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             size,
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
	}, nil
}

func (d *Driver) getImage(volID string) (*rbd.Image, string, error) {
	imgName := opensdsPrefix + volID
	imgNames, err := rbd.GetImageNames(d.ioctx)
	if err != nil {
		log.Error("When getImageNames:", err)
		return nil, "", err
	}
	for _, name := range imgNames {
		if name == imgName {
			return rbd.GetImage(d.ioctx, imgName), name, nil
		}
	}
	return nil, "", rbd.RbdErrorNotFound
}

func (d *Driver) getSize(img *rbd.Image) int64 {
	if img.Open() != nil {
		log.Error("When open image!")
		return 0
	}
	defer img.Close()

	size, err := img.GetSize()
	if err != nil {
		log.Error("When get image size:", err)
		return 0
	}
	return int64(size >> sizeShiftBit)
}

func (d *Driver) PullVolume(volID string) (*model.VolumeSpec, error) {

	err := d.initConn()
	if err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.conn.Shutdown()

	var img *rbd.Image
	var name string
	for poolName, _ := range d.conf.Pool {
		d.ioctx, err = d.conn.OpenIOContext(poolName)
		if err != nil {
			log.Error("Open IO context failed:", err)
			return nil, err
		}
		img, name, err = d.getImage(volID)
		d.ioctx.Destroy()
		if err != nil {
			if err.Error() == rbd.RbdErrorNotFound.Error() {
				continue
			}
			log.Error("When get image:", err)
			return nil, err
		}
		break
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: name,
		},
		Size: d.getSize(img),
	}, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	poolName, ok := opt.GetMetadata()["poolName"]
	if !ok {
		err := errors.New("Failed to find poolName in volume metadata!")
		log.Error(err)
		return err
	}
	if err := d.init(poolName); err != nil {
		log.Error("Connect ceph failed.")
		return err
	}
	defer d.destroy()

	img, _, err := d.getImage(opt.GetId())
	if err != nil {
		return err
	}
	if err = img.Remove(); err != nil {
		log.Error("When remove image:", err)
		return err
	}
	log.Info("Remove image success, volume id =", opt.GetId())
	return nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	poolName, ok := opt.GetMetadata()["poolName"]
	if !ok {
		err := errors.New("Failed to find poolName in volume metadata!")
		log.Error(err)
		return nil, err
	}
	return &model.ConnectionInfo{
		DriverVolumeType: RBDProtocol,
		ConnectionData: map[string]interface{}{
			"secret_type":  "ceph",
			"name":         poolName + "/" + opensdsPrefix + opt.GetVolumeId(),
			"cluster_name": "ceph",
			"hosts":        []string{opt.GetHostInfo().Host},
			"volume_id":    opt.GetVolumeId(),
			"access_mode":  "rw",
			"ports":        []string{"6789"},
		},
	}, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error { return nil }

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	poolName, ok := opt.GetMetadata()["poolName"]
	if !ok {
		err := errors.New("Failed to find poolName in volume metadata!")
		log.Error(err)
		return nil, err
	}
	if err := d.init(poolName); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroy()

	img, _, err := d.getImage(opt.GetVolumeId())
	if err != nil {
		log.Error("When get image:", err)
		return nil, err
	}
	if err = img.Open(); err != nil {
		log.Error("When open image:", err)
		return nil, err
	}
	id := opt.GetId()
	name := opensdsPrefix + id
	if _, err = img.CreateSnapshot(name); err != nil {
		log.Error("When create snapshot:", err)
		return nil, err
	}

	img.Close()

	log.Infof("Create snapshot (name:%s, id:%s, volID:%s) success",
		opt.GetName(), id, opt.GetVolumeId())

	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: id,
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        d.getSize(img),
		Metadata: map[string]string{
			"poolName": poolName,
		},
	}, nil

}

func (d *Driver) visitSnapshot(snapID string, fn func(imgName string, img *rbd.Image, snap *rbd.SnapInfo) error) error {
	imageNames, err := rbd.GetImageNames(d.ioctx)
	if err != nil {
		log.Error("When getImageNames:", err)
		return err
	}
	for _, imgName := range imageNames {
		//Filter the snapshots that not belong OpenSDS
		if !strings.HasPrefix(imgName, opensdsPrefix) {
			continue
		}
		img := rbd.GetImage(d.ioctx, imgName)
		if err = img.Open(); err != nil {
			log.Error("When open image:", err)
			return err
		}
		snapInfos, err := img.GetSnapshotNames()
		img.Close()
		if err != nil {
			log.Error("When GetSnapshotNames:", err)
			continue
		}

		snapName := opensdsPrefix + snapID
		for _, snapInfo := range snapInfos {
			if snapName == snapInfo.Name {
				return fn(imgName, img, &snapInfo)
			}
		}
	}
	return fmt.Errorf("Not found the snapshot(%s)", snapID)
}

func (d *Driver) PullSnapshot(snapID string) (*model.VolumeSnapshotSpec, error) {
	return nil, fmt.Errorf("Ceph PullSnapshot has not implemented yet.")
	/*
		if err := d.init(); err != nil {
			log.Error("Connect ceph failed.")
			return nil, err
		}
		defer d.destroy()
		var snapshot *model.VolumeSnapshotSpec
		err := d.visitSnapshot(snapID, func(volName *Name, img *rbd.Image, snap *rbd.SnapInfo) error {
			snapName := ParseName(snap.Name)
			snapshot = &model.VolumeSnapshotSpec{
				BaseModel: &model.BaseModel{
					Id: snapName.GetUUID(),
				},
				Name:     snapName.GetName(),
				Size:     int64(snap.Size >> sizeShiftBit),
				VolumeId: volName.ID,
			}
			return nil
		})
		return snapshot, err
	*/
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	poolName, ok := opt.GetMetadata()["poolName"]
	if !ok {
		err := errors.New("Failed to find poolName in volume metadata!")
		log.Error(err)
		return err
	}
	if err := d.init(poolName); err != nil {
		log.Error("Connect ceph failed.")
		return err
	}
	defer d.destroy()
	err := d.visitSnapshot(opt.GetId(), func(volName string, img *rbd.Image, snap *rbd.SnapInfo) error {
		if err := img.Open(snap.Name); err != nil {
			log.Error("When open image:", err)
		}
		snapshot := img.GetSnapshot(snap.Name)
		if err := snapshot.Remove(); err != nil {
			log.Error("When remove snapshot:", err)
			return err
		}
		img.Close()
		log.Infof("Delete snapshot (%s) success", opt.GetVolumeId())
		return nil
	})
	return err
}

func (d *Driver) parseCapStr(cap string) int64 {
	if cap == "0" {
		return 0
	}
	UnitMapper := map[string]uint64{
		"K": 0, //shift bit
		"M": 10,
		"G": 20,
		"T": 30,
		"P": 40,
	}
	unit := strings.ToUpper(cap[len(cap)-1:])
	num, err := strconv.ParseInt(cap[:len(cap)-1], 10, 64)
	if err != nil {
		log.Error("Cannot convert this number", err)
		return 0
	}
	if val, ok := UnitMapper[unit]; ok {
		return num << val >> 20 //Convert to unit GB
	} else {
		log.Error("strage unit is not found.")
		return 0
	}
}

func (d *Driver) getPoolsCapInfo() ([][]string, error) {
	output, err := execCmd("ceph df -c " + d.conf.ConfigFile)
	if err != nil {
		log.Error("[Error]:", err)
		return nil, err
	}
	lines := strings.Split(output, "\n")
	var poolsInfo [][]string
	var started = false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if started {
			poolsInfo = append(poolsInfo, strings.Fields(line))
		}
		if strings.HasPrefix(line, "POOLS:") {
			started = true
			i++
		}
	}
	return poolsInfo, nil
}

func (d *Driver) getGlobalCapInfo() ([]string, error) {
	output, err := execCmd("ceph df -c " + d.conf.ConfigFile)
	if err != nil {
		log.Error("[Error]:", err)
		return nil, err
	}
	lines := strings.Split(output, "\n")
	var globalCapInfoLine int
	for i, line := range lines {
		if strings.HasPrefix(line, "GLOBAL:") {
			globalCapInfoLine = i + 2
		}
	}
	return strings.Fields(lines[globalCapInfoLine]), nil
}

func (d *Driver) getPoolsAttr() (map[string][]string, error) {
	cmd := "ceph osd pool ls detail -c " + d.conf.ConfigFile + "| grep \"^pool\"| awk '{print $3, $4, $6, $10}'"
	output, err := execCmd(cmd)
	if err != nil {
		log.Error("[Error]:", err)
		return nil, err
	}
	lines := strings.Split(output, "\n")
	var poolDetail = make(map[string][]string)
	for i := range lines {
		if lines[i] == "" {
			continue
		}
		str := strings.Fields(lines[i])
		key := strings.Replace(str[0], "'", "", -1)
		val := str[1:]
		poolDetail[key] = val
	}
	return poolDetail, nil
}

func (d *Driver) buildPoolExtras(line []string, extras model.StoragePoolExtraSpec) model.StoragePoolExtraSpec {
	extras.Advanced = make(map[string]interface{})
	extras.Advanced["redundancyType"] = line[poolType]
	if extras.Advanced["redundancyType"] == "replicated" {
		extras.Advanced["replicateSize"] = line[poolTypeSize]
	} else {
		extras.Advanced["erasureSize"] = line[poolTypeSize]
	}
	extras.Advanced["crushRuleset"] = line[poolCrushRuleset]

	return extras
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	pc, err := d.getPoolsCapInfo()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	gc, err := d.getGlobalCapInfo()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	pa, err := d.getPoolsAttr()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var pols []*model.StoragePoolSpec
	for i := range pc {
		name := pc[i][poolName]
		c := d.conf
		if _, ok := c.Pool[name]; !ok {
			continue
		}

		extras := d.buildPoolExtras(pa[name], c.Pool[name].Extras)
		totalCap := d.parseCapStr(gc[globalSize])
		maxAvailCap := d.parseCapStr(pc[i][poolMaxAvail])
		availCap := d.parseCapStr(gc[globalAvail])
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name: name,
			//if redundancy type is replicate, MAX AVAIL =  AVAIL / replicate number,
			//and if it is erasure, MAX AVAIL =  AVAIL * k / (m + k)
			TotalCapacity:    totalCap * maxAvailCap / availCap,
			FreeCapacity:     maxAvailCap,
			StorageType:      c.Pool[name].StorageType,
			Extras:           extras,
			AvailabilityZone: c.Pool[name].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = defaultAZ
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {
	return nil, &model.NotImplementError{S: "Method InitializeSnapshotConnection has not been implemented yet"}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &model.NotImplementError{S: "Method TerminateSnapshotConnection has not been implemented yet"}
}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"Method CreateVolumeGroup has not been implemented yet"}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts, vg *model.VolumeGroupSpec, addVolumesRef []*model.VolumeSpec, removeVolumesRef []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, []*model.VolumeSpec, error) {
	return nil, nil, nil, &model.NotImplementError{"Method UpdateVolumeGroup has not been implemented yet"}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts, vg *model.VolumeGroupSpec, volumes []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, error) {
	return nil, nil, &model.NotImplementError{"Method UpdateVolumeGroup has not been implemented yet"}
}
