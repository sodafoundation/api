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
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/lvm/targets"
	pb "github.com/opensds/opensds/pkg/dock/proto"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

const (
	vgName = "vg001"
)

var conf = LVMConfig{}

type Driver struct {
	config LVMConfig
}

type LVMConfig struct {
	Pool map[string]PoolProperties `yaml:"pool,flow"`
}

type PoolProperties struct {
	DiskType  string `yaml:"diskType"`
	IOPS      int64  `yaml:"iops"`
	BandWidth int64  `yaml:"bandwidth"`
}

func (d *Driver) Setup() error {
	// Read lvm config file
	confYaml, err := ioutil.ReadFile(config.CONF.LVMConfig)
	if err != nil {
		log.Fatalf("Read lvm config yaml file (%s) failed, reason:(%v)", config.CONF.LVMConfig, err)
		return err
	}
	if err = yaml.Unmarshal(confYaml, &conf); err != nil {
		log.Fatal("Parse error: %v", err)
		return err
	}
	d.config = conf

	return nil
}

func (*Driver) Unset() error { return nil }

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	var size = fmt.Sprint(opt.GetSize()) + "G"

	cmd := strings.Join([]string{"lvcreate", "-n", opt.GetName(), "-L", size, vgName}, " ")
	if _, err := d.execCmd(cmd); err != nil {
		log.Error("Failed to create logic volume:", err)
		return nil, err
	}

	var lvPath, lvStatus string
	// Display and parse some metadata in logic volume returned.
	lvPath = strings.Join([]string{"/dev", vgName, opt.GetName()}, "/")
	lv, err := d.execCmd("lvdisplay " + lvPath)
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
	lv, err := d.execCmd("lvmdisplay " + volIdentifier)
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
	cmd := strings.Join([]string{"lvremove", "-f", lvPath}, " ")
	if _, err := d.execCmd(cmd); err != nil {
		log.Error("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

func (*Driver) InitializeConnection(opt *pb.CreateAttachmentOpts) (*model.ConnectionInfo, error) {
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

	t := targets.NewTarget()
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

func (*Driver) TerminateConnection(opt *pb.DeleteAttachmentOpts) error {
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

	t := targets.NewTarget()
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

	cmd := strings.Join([]string{"lvcreate", "-n", opt.GetName(), "-L", size, "-p r", "-s", lvPath}, " ")
	if _, err := d.execCmd(cmd); err != nil {
		log.Error("Failed to create logic volume snapshot:", err)
		return nil, err
	}

	var lvsPath, lvStatus string
	lvsPath = strings.Join([]string{"/dev", vgName, opt.GetName()}, "/")
	// Display and parse some metadata in logic volume snapshot returned.
	lvs, err := d.execCmd("lvdisplay " + lvsPath)
	if err != nil {
		log.Error("Failed to display logic volume snapshot:", err)
		return nil, err
	}
	for _, line := range strings.Split(lvs, "\n") {
		if strings.Contains(line, "LV Path") {
			lvsPath = strings.Fields(line)[2]
		}
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
	lv, err := d.execCmd("lvmdisplay " + snapIdentifier)
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
	cmd := strings.Join([]string{"lvremove", "-f", lvsPath}, " ")
	if _, err := d.execCmd(cmd); err != nil {
		log.Error("Failed to remove logic volume:", err)
		return err
	}

	return nil
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	vgs, err := d.execCmd("vgdisplay")
	if err != nil {
		return nil, err
	}
	log.Info("Got vgs info:", vgs)

	var tCapacity, fCapacity int64
	for _, line := range strings.Split(vgs, "\n") {
		if strings.Contains(line, "VG Size") {
			slice := strings.Fields(line)
			cap, _ := strconv.ParseFloat(slice[2], 64)
			tCapacity = int64(cap)
		}
		if strings.Contains(line, "Free  PE / Size") {
			slice := strings.Fields(line)
			cap, _ := strconv.ParseFloat(slice[len(slice)-2], 64)
			fCapacity = int64(cap)
		}
	}

	var pols []*model.StoragePoolSpec
	if _, ok := d.config.Pool[vgName]; !ok {
		return pols, nil
	}
	param := d.buildPoolParam(d.config.Pool[vgName])
	pol := &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: uuid.NewV5(uuid.NamespaceOID, vgName).String(),
		},
		Name:          vgName,
		TotalCapacity: tCapacity,
		FreeCapacity:  fCapacity,
		Parameters:    *param,
	}
	pols = append(pols, pol)

	return pols, nil
}

func (*Driver) buildPoolParam(proper PoolProperties) *map[string]interface{} {
	var param = make(map[string]interface{})
	param["diskType"] = proper.DiskType
	param["iops"] = proper.IOPS
	param["bandwidth"] = proper.BandWidth

	return &param
}

func (*Driver) execCmd(cmd string) (string, error) {
	ret, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return string(ret), nil
}
