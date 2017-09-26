// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"fmt"
	"sync"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	"github.com/go-yaml/yaml"
	log "github.com/golang/glog"
	api "github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

const (
	opensdsPrefix string = "OPENSDS"
	splitChar            = ":"
	sizeShiftBit         = 20
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

type PoolProperties struct {
	DiskType  string `yaml:"diskType"`
	IOPS      int    `yaml:"iops"`
	BandWitdh string `yaml:"bandWitdh"`
}

type CephConfig struct {
	ConfigFile string                    `yaml:"configFile,omitempty"`
	Pool       map[string]PoolProperties `yaml:"pool,flow"`
}

var cephConfig *CephConfig
var once sync.Once

func (c *CephConfig) Load(file string) error {
	// Set /etc/ceph/ceph.conf as default value
	confYaml, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Read ceph config yaml file (%s) failed, reason:(%v)", file, err)
		return err
	}
	err = yaml.Unmarshal([]byte(confYaml), c)
	if err != nil {
		log.Fatal("Parse error: %v", err)
		return err
	}
	return nil
}

func getConfig() *CephConfig {
	once.Do(func() {
		cephConfig = &CephConfig{ConfigFile: "/etc/ceph/ceph.conf"}
		cephConfig.Load(config.CONF.OsdsDock.CephConfig)
	})
	return cephConfig
}

type Name struct {
	Name string
	ID   string
}

func NewName(name string) *Name {
	return &Name{
		Name: name,
		ID:   uuid.NewV4().String(),
	}
}

func ParseName(fullName string) *Name {
	if !strings.HasPrefix(fullName, opensdsPrefix) {
		return nil
	}

	nameInfo := strings.Split(fullName, splitChar)

	return &Name{
		Name: nameInfo[1],
		ID:   nameInfo[2],
	}
}

func (name *Name) GetFullName() string {
	return opensdsPrefix + ":" + name.Name + ":" + name.ID
}

func (name *Name) GetName() string {
	return name.Name
}

func (name *Name) GetUUID() string {
	return name.ID
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
}

func (d *Driver) Setup() {}

func (d *Driver) Unset() {}

func (d *Driver) initConn() error {
	conn, err := rados.NewConn()
	if err != nil {
		log.Error("New connect failed:", err)
		return err
	}

	if err = conn.ReadConfigFile(getConfig().ConfigFile); err != nil {
		log.Error("Read config file failed:", err)
		return err
	}
	if err = conn.Connect(); err != nil {
		log.Error("Connect failed:", err)
		return err
	}
	d.ioctx, err = conn.OpenIOContext("rbd")
	if err != nil {
		log.Error("Open IO context failed:", err)
		return err
	}
	d.conn = conn
	return nil
}

func (d *Driver) destroyConn() {
	defer d.conn.Shutdown()
	defer d.ioctx.Destroy()
}

func (d *Driver) CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	if err := d.initConn(); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroyConn()

	imgName := NewName(name)
	_, err := rbd.Create(d.ioctx, imgName.GetFullName(), uint64(size)<<sizeShiftBit, 20)
	if err != nil {
		log.Errorf("Create rbd image (%s) failed, (%v)", name, err)
		return nil, err
	}

	log.Infof("Create volume %s (%s) success.", name, imgName.GetUUID())
	return &api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id: imgName.GetUUID(),
		},
		Name:             imgName.GetName(),
		Size:             size,
		Description:      "",
		AvailabilityZone: "ceph",
	}, nil
}

func (d *Driver) getImage(volID string) (*rbd.Image, *Name, error) {
	imgNames, err := rbd.GetImageNames(d.ioctx)
	if err != nil {
		log.Error("When getImageNames:", err)
		return nil, nil, err
	}
	for _, fullName := range imgNames {
		name := ParseName(fullName)
		if name != nil && name.ID == volID {
			return rbd.GetImage(d.ioctx, fullName), name, nil
		}
	}
	return nil, nil, rbd.RbdErrorNotFound
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

func (d *Driver) GetVolume(volID string) (*api.VolumeSpec, error) {
	if err := d.initConn(); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroyConn()

	img, name, err := d.getImage(volID)
	if err != nil {
		log.Error("When get image:", err)
		return nil, err
	}

	return &api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id: name.GetUUID(),
		},
		Name:             name.GetName(),
		Size:             d.getSize(img),
		Description:      "",
		AvailabilityZone: "ceph",
	}, nil
}

func (d *Driver) DeleteVolume(volID string) error {
	if err := d.initConn(); err != nil {
		log.Error("Connect ceph failed.")
		return err
	}
	defer d.destroyConn()

	img, _, err := d.getImage(volID)
	if err != nil {
		return err
	}
	if err = img.Remove(); err != nil {
		log.Error("When remove image:", err)
		return err
	}
	log.Info("Remove image success, volume id =", volID)
	return nil
}

func (d *Driver) InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error) {
	if err := d.initConn(); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroyConn()

	vol, err := d.GetVolume(volID)
	if err != nil {
		log.Error("When get image:", err)
		return nil, err
	}

	return &api.ConnectionInfo{
		DriverVolumeType: "rbd",
		ConnectionData: map[string]interface{}{
			"secret_type":  "ceph",
			"name":         "rbd/" + opensdsPrefix + ":" + vol.Name + ":" + vol.Id,
			"cluster_name": "ceph",
			"hosts":        []string{hostInfo.Host},
			"volume_id":    vol.Id,
			"access_mode":  "rw",
			"ports":        []string{"6789"},
		},
	}, nil
}

func (d *Driver) AttachVolume(volID, host, mountpoint string) error {
	return nil
}

func (d *Driver) DetachVolume(volID string) error {
	return nil
}

func (d *Driver) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error) {
	if err := d.initConn(); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroyConn()

	img, _, err := d.getImage(volID)
	if err != nil {
		log.Error("When get image:", err)
		return nil, err
	}
	if err = img.Open(); err != nil {
		log.Error("When open image:", err)
		return nil, err
	}
	defer img.Close()

	fullName := NewName(name)
	if _, err = img.CreateSnapshot(fullName.GetFullName()); err != nil {
		log.Error("When create snapshot:", err)
		return nil, err
	}
	log.Infof("Create snapshot success, name:%s, id:%s, volID:%s", name, volID, fullName.GetUUID())
	return &api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{
			Id: fullName.GetUUID(),
		},
		Name:        fullName.GetName(),
		Description: description,
		VolumeId:    volID,
		Size:        d.getSize(img),
	}, nil
}

func (d *Driver) visitSnapshot(snapID string, fn func(volName *Name, img *rbd.Image, snap *rbd.SnapInfo) error) error {
	imageNames, err := rbd.GetImageNames(d.ioctx)
	if err != nil {
		log.Error("When getImageNames:", err)
		return err
	}
	for _, name := range imageNames {
		in := ParseName(name)
		//Filter the snapshots that not belong OpenSDS
		if in == nil {
			continue
		}
		img := rbd.GetImage(d.ioctx, name)
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
		for _, snapInfo := range snapInfos {
			name := ParseName(snapInfo.Name)
			if snapID == name.GetUUID() {
				return fn(in, img, &snapInfo)
			}
		}
	}
	reason := fmt.Sprintf("Not found the snapshot(%s)", snapID)
	return errors.New(reason)
}

func (d *Driver) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	if err := d.initConn(); err != nil {
		log.Error("Connect ceph failed.")
		return nil, err
	}
	defer d.destroyConn()
	var snapshot *api.VolumeSnapshotSpec
	err := d.visitSnapshot(snapID, func(volName *Name, img *rbd.Image, snap *rbd.SnapInfo) error {
		snapName := ParseName(snap.Name)
		snapshot = &api.VolumeSnapshotSpec{
			BaseModel: &api.BaseModel{
				Id: snapName.GetUUID(),
			},
			Name:     snapName.GetName(),
			Size:     int64(snap.Size >> sizeShiftBit),
			VolumeId: volName.ID,
		}
		return nil
	})
	return snapshot, err
}

func (d *Driver) DeleteSnapshot(snapID string) error {
	if err := d.initConn(); err != nil {
		log.Error("Connect ceph failed.")
		return err
	}
	defer d.destroyConn()
	err := d.visitSnapshot(snapID, func(volName *Name, img *rbd.Image, snap *rbd.SnapInfo) error {
		if err := img.Open(snap.Name); err != nil {
			log.Error("When open image:", err)
		}
		snapshot := img.GetSnapshot(snap.Name)
		if err := snapshot.Remove(); err != nil {
			log.Error("When remove snapshot:", err)
			return err
		}
		img.Close()
		log.Info("Delete snapshot {%s} success", ParseName(snap.Name).GetUUID())
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
		return num << val >> sizeShiftBit
	} else {
		log.Error("strage unit is not found.")
		return 0
	}
}

func (d *Driver) getPoolsCapInfo() ([][]string, error) {
	const poolStartLine = 5
	output, err := execCmd("ceph df -c " + getConfig().ConfigFile)
	if err != nil {
		log.Error("[Error]:", err)
		return nil, err
	}
	lines := strings.Split(output, "\n")
	var poolsInfo [][]string
	for i := poolStartLine; i < len(lines); i++ {
		poolsInfo = append(poolsInfo, strings.Fields(lines[i]))
	}
	return poolsInfo, nil
}

func (d *Driver) getGlobalCapInfo() ([]string, error) {
	const globalCapInfoLine = 2
	output, err := execCmd("ceph df")
	if err != nil {
		log.Error("[Error]:", err)
		return nil, err
	}
	lines := strings.Split(output, "\n")
	return strings.Fields(lines[globalCapInfoLine]), nil
}

func (d *Driver) getPoolsAttr() (map[string][]string, error) {
	cmd := "ceph osd pool ls detail | grep \"^pool\"| awk '{print $3, $4, $6, $10}'"
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

func (d *Driver) buildPoolParam(line []string, proper PoolProperties) *map[string]interface{} {
	param := make(map[string]interface{})
	param["diskType"] = proper.DiskType
	param["iops"] = proper.IOPS
	param["bandWidth"] = proper.BandWitdh
	param["redundancyType"] = line[poolType]
	if param["redundancyType"] == "replicated" {
		param["replicateSize"] = line[poolTypeSize]
	} else {
		param["erasureSize"] = line[poolTypeSize]
	}
	param["crushRuleset"] = line[poolCrushRuleset]
	return &param
}

func (d *Driver) ListPools() (*[]api.StoragePoolSpec, error) {
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

	var poolList []api.StoragePoolSpec
	for i := range pc {
		name := pc[i][poolName]
		c := getConfig()
		if _, ok := c.Pool[name]; !ok {
			continue
		}
		param := d.buildPoolParam(pa[name], c.Pool[name])
		totalCap := d.parseCapStr(gc[globalSize])
		maxAvailCap := d.parseCapStr(pc[i][poolMaxAvail])
		availCap := d.parseCapStr(gc[globalAvail])
		pool := api.StoragePoolSpec{
			BaseModel: &api.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, name).String(),
			},
			Name: name,
			//if redundancy type is replicate, MAX AVAIL =  AVAIL / replicate number,
			//and it this is erasure, MAX AVAIL =  AVAIL * k / (m + k)
			TotalCapacity: totalCap * maxAvailCap / availCap,
			FreeCapacity:  maxAvailCap,
			Parameters:    *param,
		}
		poolList = append(poolList, pool)
	}
	return &poolList, nil
}


