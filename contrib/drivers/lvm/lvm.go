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

package lvm

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/lvm/targets"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

const (
	defaultTgtConfDir = "/etc/tgt/conf.d"
	defaultTgtBindIp  = "127.0.0.1"
	defaultConfPath   = "/etc/opensds/driver/lvm.yaml"
	volumePrefix      = "volume-"
	snapshotPrefix    = "_snapshot-"
	blocksize         = 4096
	sizeShiftBit      = 30
)

type LvInfo struct {
	Name string
	Vg   string
	Size int64
}

type LVMConfig struct {
	TgtBindIp      string                    `yaml:"tgtBindIp"`
	TgtConfDir     string                    `yaml:"tgtConfDir"`
	EnableChapAuth bool                      `yaml:"enableChapAuth"`
	Pool           map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf *LVMConfig

	handler func(script string, cmd []string) (string, error)
}

func (d *Driver) Setup() error {
	// Read lvm config file
	d.conf = &LVMConfig{TgtBindIp: defaultTgtBindIp, TgtConfDir: defaultTgtConfDir}
	p := config.CONF.OsdsDock.Backends.LVM.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	if _, err := Parse(d.conf, p); err != nil {
		return err
	}
	d.handler = execCmd

	return nil
}

func (*Driver) Unset() error { return nil }

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	var size = fmt.Sprint(opt.GetSize()) + "G"
	var polName = opt.GetPoolName()
	var id = opt.GetId()
	var name = volumePrefix + id

	if _, err := d.handler("lvcreate", []string{
		"-Z", "n",
		"-n", name, // use uuid instead of name.
		"-L", size,
		polName,
	}); err != nil {
		log.Error("Failed to create logic volume:", err)
		return nil, err
	}

	var lvPath, lvStatus string
	// Display and parse some metadata in logic volume returned.
	lvPath = path.Join("/dev", polName, name)
	lv, err := d.handler("lvdisplay", []string{lvPath})
	if err != nil {
		log.Error("Failed to display logic volume:", err)
		return nil, err
	}

	for _, line := range strings.Split(lv, "\n") {
		if strings.Contains(line, "LV Path") {
			lvPath = strings.Fields(line)[2]
		}
		if strings.Contains(line, "LV Status") {
			lvStatus = strings.Fields(line)[2]
		}
	}

	// Copy snapshot to volume
	var snap = opt.GetSnapshotId()
	if snap != "" {
		var snapSize = uint64(opt.GetSnapshotSize())
		var count = (snapSize<<sizeShiftBit)/blocksize
		var snapName = snapshotPrefix + snap
		var snapPath = path.Join("/dev", polName, snapName)
		if _, err := d.handler("dd", []string{
			"if=" + snapPath,
			"of=" + lvPath,
			"count=" + fmt.Sprint(count),
			"bs=" + fmt.Sprint(blocksize),
		}); err != nil {
			log.Error("Failed to create logic volume:", err)
			return nil, err
		}
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Status:      lvStatus,
		Metadata: map[string]string{
			"lvPath": lvPath,
		},
	}, nil
}

func (d *Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
	// Display and parse some metadata in logic volume returned.
	lv, err := d.handler("lvdisplay", []string{volIdentifier})
	if err != nil {
		log.Error("Failed to display logic volume:", err)
		return nil, err
	}
	var lvStatus string
	for _, line := range strings.Split(lv, "\n") {
		if strings.Contains(line, "LV Status") {
			lvStatus = strings.Fields(line)[2]
		}
	}

	return &model.VolumeSpec{
		Status: lvStatus,
	}, nil
}

func (d *Driver) geLvInfos() ([]*LvInfo, error) {
	var lvList []*LvInfo
	args := []string{"--noheadings", "--unit=g", "-o", "vg_name,name,size", "--nosuffix"}
	info, err := d.handler("lvs", args)
	if err != nil {
		log.Error("Get volume failed", err)
		return lvList, err
	}
	for _, line := range strings.Split(info, "\n") {
		if len(line) == 0 {
			continue
		}
		words := strings.Fields(line)
		size, _ := strconv.ParseInt(words[2], 10, 64)
		lv := &LvInfo{
			Vg:   words[0],
			Name: words[1],
			Size: size,
		}
		lvList = append(lvList, lv)
	}
	return lvList, nil
}

func (d *Driver) volumeExists(id string) bool {
	lvList, _ := d.geLvInfos()
	name := volumePrefix + id
	for _, lv := range lvList {
		if lv.Name == name {
			return true
		}
	}
	return false
}

func (d *Driver) lvHasSnapshot(lvPath string) bool {
	args := []string{"--noheading", "-C", "-o", "Attr", lvPath}
	info, err := d.handler("lvdisplay", args)
	if err != nil {
		log.Error("Failed to remove logic volume:", err)
		return false
	}
	info = strings.Trim(info, " ")
	return info[0] == 'o' || info[0] == 'O'
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {

	id := opt.GetId()
	if !d.volumeExists(id) {
		log.Warningf("Volume(%s) does not exist, nothing to remove", id)
		return nil
	}

	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("failed to find logic volume path in volume metadata")
		log.Error(err)
		return err
	}

	if d.lvHasSnapshot(lvPath) {
		err := fmt.Errorf("unable to delete due to existing snapshot for volume: %s", id)
		log.Error(err)
		return err
	}

	if _, err := d.handler("lvremove", []string{"-f", lvPath}); err != nil {
		log.Error("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

// ExtendVolume ...
func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("failed to find logic volume path in volume metadata")
		log.Error(err)
		return nil, err
	}

	var size = fmt.Sprint(opt.GetSize()) + "G"

	if _, err := d.handler("lvresize", []string{
		"-L", size,
		lvPath,
	}); err != nil {
		log.Error("Failed to extend logic volume:", err)
		return nil, err
	}

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Metadata:    opt.GetMetadata(),
	}, nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	initiator := opt.HostInfo.GetInitiator()
	if initiator == "" {
		initiator = "ALL"
	}

	hostIP := opt.HostInfo.GetIp()
	if hostIP == "" {
		hostIP = "ALL"
	}

	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume attachment metadata!")
		log.Error(err)
		return nil, err
	}
	var chapAuth []string
	if d.conf.EnableChapAuth {
		chapAuth = []string{utils.RandSeqWithAlnum(20), utils.RandSeqWithAlnum(16)}
	}
	t := targets.NewTarget(d.conf.TgtBindIp, d.conf.TgtConfDir)
	expt, err := t.CreateExport(opt.GetVolumeId(), lvPath, hostIP, initiator, chapAuth)
	if err != nil {
		log.Error("Failed to initialize connection of logic volume:", err)
		return nil, err
	}

	return &model.ConnectionInfo{
		DriverVolumeType: "iscsi",
		ConnectionData:   expt,
	}, nil
}

func (d *Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error {
	t := targets.NewTarget(d.conf.TgtBindIp, d.conf.TgtConfDir)
	if err := t.RemoveExport(opt.GetVolumeId()); err != nil {
		log.Error("Failed to initialize connection of logic volume:", err)
		return err
	}
	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	var size = fmt.Sprint(opt.GetSize()) + "G"
	var id = opt.GetId()
	var snapName = snapshotPrefix + id
	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume snapshot metadata!")
		log.Error(err)
		return nil, err
	}

	if _, err := d.handler("lvcreate", []string{
		"-n", snapName,
		"-L", size,
		"-p", "r",
		"-s", lvPath,
	}); err != nil {
		log.Error("Failed to create logic volume snapshot:", err)
		return nil, err
	}

	var lvsDir, lvsPath string
	lvsDir, _ = path.Split(lvPath)
	lvsPath = path.Join(lvsDir, snapName)
	// Display and parse some metadata in logic volume snapshot returned.
	lvs, err := d.handler("lvdisplay", []string{lvsPath})
	if err != nil {
		log.Error("Failed to display logic volume snapshot:", err)
		return nil, err
	}
	var lvStatus string
	for _, line := range strings.Split(lvs, "\n") {
		if strings.Contains(line, "LV Status") {
			lvStatus = strings.Fields(line)[2]
		}
	}

	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: id,
		},
		Name:        opt.GetName(),
		Size:        opt.GetSize(),
		Description: opt.GetDescription(),
		Status:      lvStatus,
		VolumeId:    opt.GetVolumeId(),
		Metadata: map[string]string{
			"lvsPath": lvsPath,
		},
	}, nil
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	// Display and parse some metadata in logic volume snapshot returned.
	lv, err := d.handler("lvdisplay", []string{snapIdentifier})
	if err != nil {
		log.Error("Failed to display logic volume snapshot:", err)
		return nil, err
	}
	var lvStatus string
	for _, line := range strings.Split(lv, "\n") {
		if strings.Contains(line, "LV Status") {
			lvStatus = strings.Fields(line)[2]
		}
	}

	return &model.VolumeSnapshotSpec{
		Status: lvStatus,
	}, nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	lvsPath, ok := opt.GetMetadata()["lvsPath"]
	if !ok {
		err := errors.New("Failed to find logic volume snapshot path in volume snapshot metadata!")
		log.Error(err)
		return err
	}
	if _, err := d.handler("lvremove", []string{
		"-f", lvsPath,
	}); err != nil {
		log.Error("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

type VolumeGroup struct {
	Name          string
	TotalCapacity int64
	FreeCapacity  int64
	UUID          string
}

func (d *Driver) getVGList() (*[]VolumeGroup, error) {
	info, err := d.handler("vgs", []string{
		"--noheadings", "--nosuffix",
		"--unit=g",
		"-o", "name,size,free,uuid",
	})
	if err != nil {
		return nil, err
	}

	lines := strings.Split(info, "\n")
	var vgs []VolumeGroup
	for _, line := range lines {
		val := strings.Fields(line)
		if len(val) != 4 {
			continue
		}

		capa, _ := strconv.ParseFloat(val[1], 64)
		total := int64(capa)
		capa, _ = strconv.ParseFloat(val[2], 64)
		free := int64(capa)

		vg := VolumeGroup{
			Name:          val[0],
			TotalCapacity: total,
			FreeCapacity:  free,
			UUID:          val[3],
		}
		vgs = append(vgs, vg)
	}
	return &vgs, nil
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {

	vgs, err := d.getVGList()
	if err != nil {
		return nil, err
	}

	var pols []*model.StoragePoolSpec
	for _, vg := range *vgs {
		if _, ok := d.conf.Pool[vg.Name]; !ok {
			continue
		}

		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, vg.UUID).String(),
			},
			Name:             vg.Name,
			TotalCapacity:    vg.TotalCapacity,
			FreeCapacity:     vg.FreeCapacity,
			StorageType:      d.conf.Pool[vg.Name].StorageType,
			Extras:           d.conf.Pool[vg.Name].Extras,
			AvailabilityZone: d.conf.Pool[vg.Name].AvailabilityZone,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = "default"
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"Method CreateVolumeGroup did not implement."}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts, vg *model.VolumeGroupSpec, addVolumesRef []*model.VolumeSpec, removeVolumesRef []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, []*model.VolumeSpec, error) {
	return nil, nil, nil, &model.NotImplementError{"Method UpdateVolumeGroup did not implement."}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts, vg *model.VolumeGroupSpec, volumes []*model.VolumeSpec) (*model.VolumeGroupSpec, []*model.VolumeSpec, error) {
	return nil, nil, &model.NotImplementError{"Method UpdateVolumeGroup did not implement."}
}

func execCmd(script string, cmd []string) (string, error) {
	log.Infof("Command: %s %s", script, strings.Join(cmd, " "))
	info, err := exec.Command(script, cmd...).Output()
	if err != nil {
		log.Error(info, err.Error())
		return "", err
	}
	log.V(8).Infof("Command Result:\n%s", string(info))
	return string(info), nil
}
