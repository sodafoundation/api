// Copyright (c) 2017 OpenSDS Authors.
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
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

const (
	defaultConfPath = "/etc/opensds/driver/lvm.yaml"
)

type LVMConfig struct {
	TgtBindIp string                    `yaml:"tgtBindIp"`
	Pool      map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf *LVMConfig

	handler func(script string, cmd []string) (string, error)
}

func (d *Driver) Setup() error {
	// Read lvm config file
	d.conf = &LVMConfig{TgtBindIp: "127.0.0.1"}
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

	if _, err := d.handler("lvcreate", []string{
		"-n", opt.GetName(),
		"-L", size,
		polName,
	}); err != nil {
		log.Error("Failed to create logic volume:", err)
		return nil, err
	}

	var lvPath, lvStatus string
	// Display and parse some metadata in logic volume returned.
	lvPath = path.Join("/dev", polName, opt.GetName())
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

	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: uuid.NewV4().String(),
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

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume metadata!")
		log.Error(err)
		return err
	}
	if _, err := d.handler("lvremove", []string{
		"-f", lvPath,
	}); err != nil {
		log.Error("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
	var initiator string
	if initiator = opt.HostInfo.GetInitiator(); initiator == "" {
		initiator = "ALL"
	}
	// TODO	Add lvm path in Metadata field.
	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume attachment metadata!")
		log.Error(err)
		return nil, err
	}

	t := targets.NewTarget(d.conf.TgtBindIp)
	expt, err := t.CreateExport(lvPath, initiator)
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
	var initiator string
	if initiator = opt.HostInfo.GetInitiator(); initiator == "" {
		initiator = "ALL"
	}
	// TODO	Add lvm path in Metadata field.
	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume attachment metadata!")
		log.Error(err)
		return err
	}

	t := targets.NewTarget(d.conf.TgtBindIp)
	if err := t.RemoveExport(lvPath, initiator); err != nil {
		log.Error("Failed to initialize connection of logic volume:", err)
		return err
	}

	return nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	var size = fmt.Sprint(opt.GetSize()) + "G"
	lvPath, ok := opt.GetMetadata()["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume snapshot metadata!")
		log.Error(err)
		return nil, err
	}

	if _, err := d.handler("lvcreate", []string{
		"-n", opt.GetName(),
		"-L", size,
		"-p", "r",
		"-s", lvPath,
	}); err != nil {
		log.Error("Failed to create logic volume snapshot:", err)
		return nil, err
	}

	var lvsDir, lvsPath string
	lvsDir, _ = path.Split(lvPath)
	lvsPath = path.Join(lvsDir, opt.GetName())
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
			Id: uuid.NewV4().String(),
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
}

func (d *Driver) getVGList() (*[]VolumeGroup, error) {
	const vgInfoLineCount = 10
	info, err := d.handler("vgdisplay", []string{})
	if err != nil {
		return nil, err
	}

	log.Info("Got vgs info:", info)
	lines := strings.Split(info, "\n")
	vgs := make([]VolumeGroup, len(lines)/vgInfoLineCount)

	var vgIdx = -1
	for _, line := range lines {
		if strings.Contains(line, "--- Volume group ---") {
			vgIdx++
			continue
		}
		if strings.Contains(line, "VG Name") {
			slice := strings.Fields(line)
			vgs[vgIdx].Name = slice[2]
		}
		if strings.Contains(line, "VG Size") {
			slice := strings.Fields(line)
			capa, _ := strconv.ParseFloat(slice[2], 64)
			vgs[vgIdx].TotalCapacity = int64(capa)
		}
		if strings.Contains(line, "Free  PE / Size") {
			slice := strings.Fields(line)
			capa, _ := strconv.ParseFloat(slice[len(slice)-2], 64)
			vgs[vgIdx].FreeCapacity = int64(capa)
		}
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
		param := d.buildPoolParam(d.conf.Pool[vg.Name])
		pol := &model.StoragePoolSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, vg.Name).String(),
			},
			Name:             vg.Name,
			TotalCapacity:    vg.TotalCapacity,
			FreeCapacity:     vg.FreeCapacity,
			Extras:           *param,
			AvailabilityZone: d.conf.Pool[vg.Name].AZ,
		}
		if pol.AvailabilityZone == "" {
			pol.AvailabilityZone = "default"
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

func (*Driver) buildPoolParam(proper PoolProperties) *map[string]interface{} {
	var param = make(map[string]interface{})
	param["diskType"] = proper.DiskType
	return &param
}

func execCmd(script string, cmd []string) (string, error) {
	ret, err := exec.Command(script, cmd...).Output()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return string(ret), nil
}
