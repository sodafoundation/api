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
	"log"
	"strings"

	api "github.com/opensds/opensds/pkg/model"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	"github.com/satori/go.uuid"
	"os/exec"
	"strconv"
)

const (
	OPENSDS_PREFIX string = "OPENSDS"
	SPLIT_CHAR            = ":"
	SIZE_SHIFT_BIT        = 20
)

const (
	GLOBAL_SIZE = iota
	GLOBAL_AVAIL
	GLOBAL_RAW_USED
	GLOBAL_RAW_USED_PERCENTAGE
)

const (
	POOL_NAME = iota
	POOL_ID
	POOL_USED
	POOL_USED_PER
	POOL_MAX_AVAIL
	POOL_OBJECTS
)

const (
	POOL_TYPE = iota
	POOL_TYPE_SIZE
	POOL_CRUSH_RULESET
)

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
	if !strings.HasPrefix(fullName, OPENSDS_PREFIX) {
		return nil
	}

	nameInfo := strings.Split(fullName, SPLIT_CHAR)

	return &Name{
		Name: nameInfo[1],
		ID:   nameInfo[2],
	}
}

func (name *Name) GetFullName() string {
	return OPENSDS_PREFIX + ":" + name.Name + ":" + name.ID
}

func (name *Name) GetName() string {
	return name.Name
}

func (name *Name) GetUUID() string {
	return name.ID
}

// Response is a structure for all properties of
// a volume for a non detailed query
type Response struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Size              int64  `json:"size"`
	Availability_zone string `json:"availability_zone"`
}

// SnapshotResponse is a structure for all properties of
// a volume snapshot for a non detailed query
type SnapshotResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Volume_id string `json:"volume_id"`
	Size      int64  `json:"size"`
}

type PoolResponse struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name,omitempty"`
	Description      string                 `json:"description,omitempty"`
	AvailabilityZone string                 `json:"availabilityZone,omitempty"`
	TotalCapacity    int64                  `json:"totalCapacity,omitempty"`
	FreeCapacity     int64                  `json:"freeCapacity,omitempty"`
	StorageType      string                 `json:"-"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
}

type ImageMgr struct {
	Conn  *rados.Conn
	Ioctx *rados.IOContext
}

func (imgMgr *ImageMgr) Init() error {
	conn, err := rados.NewConn()
	if err != nil {
		log.Println("[Error] New connect failed:", err)
		return err
	}
	if err = conn.ReadDefaultConfigFile(); err != nil {
		log.Println("[Error] Read config file failed:", err)
		return err
	}
	if err = conn.Connect(); err != nil {
		log.Println("[Error] Connect failed:", err)
		return err
	}

	log.Println("[Info] Connect ceph cluster ok!")

	imgMgr.Ioctx, err = conn.OpenIOContext("rbd")
	if err != nil {
		log.Println("[Error] Open IO context failed:", err)
		return err
	}

	imgMgr.Conn = conn
	return nil
}

func (imgMgr *ImageMgr) Destory() {
	defer imgMgr.Conn.Shutdown()
	defer imgMgr.Ioctx.Destroy()
}

func (imgMgr *ImageMgr) CreateImage(name string, size int64) (*Response, error) {
	imageName := NewName(name)

	_, err := rbd.Create(imgMgr.Ioctx, imageName.GetFullName(), uint64(size)<<SIZE_SHIFT_BIT, 20)
	if err != nil {
		log.Println("[Error] When create rbd image:", err)
		return &Response{}, err
	}

	return &Response{
		Name:              imageName.GetName(),
		Id:                imageName.GetUUID(),
		Size:              size,
		Availability_zone: "ceph",
	}, nil
}

func (imgMgr *ImageMgr) getImage(volID string) (*rbd.Image, *Name, error) {
	imageNames, err := rbd.GetImageNames(imgMgr.Ioctx)
	if err != nil {
		log.Println("[Error] When getImageNames:", err)
		return nil, nil, err
	}
	for _, fullName := range imageNames {
		name := ParseName(fullName)
		if name != nil && name.ID == volID {
			return rbd.GetImage(imgMgr.Ioctx, fullName), name, nil
		}
	}
	return nil, nil, rbd.RbdErrorNotFound
}

func (imgMgr *ImageMgr) getSize(img *rbd.Image) int64 {
	if img.Open() != nil {
		log.Println("[Error] When open image!")
		return 0
	}
	defer img.Close()

	size, err := img.GetSize()
	if err != nil {
		log.Println("[Error] When get image size:", err)
		return 0
	}
	return int64(size >> SIZE_SHIFT_BIT)
}

func (imgMgr *ImageMgr) RemoveImage(volID string) error {
	img, _, err := imgMgr.getImage(volID)
	if err != nil {
		return err
	}
	if err = img.Remove(); err != nil {
		log.Println("[Error] When remove image:", err)
		return err
	}

	log.Println("[Info] Remove image success, volume id =", volID)
	return nil
}

func (imgMgr *ImageMgr) GetImage(volID string) (*Response, error) {
	img, name, err := imgMgr.getImage(volID)
	if err != nil {
		log.Println("[Error] When get image:", err)
		return &Response{}, err
	}

	return &Response{
		Name: name.GetName(),
		Id:   name.GetUUID(),
		Size: imgMgr.getSize(img),
	}, nil
}

func (imgMgr *ImageMgr) GetImages() (*[]Response, error) {
	imageNames, err := rbd.GetImageNames(imgMgr.Ioctx)
	if err != nil {
		log.Println("[Error] When getImageNames:", err)
		return &[]Response{}, err
	}

	var images []Response

	for _, name := range imageNames {
		in := ParseName(name)
		if in == nil {
			continue
		}
		img := rbd.GetImage(imgMgr.Ioctx, name)
		image := Response{
			Name: in.GetName(),
			Id:   in.GetUUID(),
			Size: imgMgr.getSize(img),
		}
		images = append(images, image)
	}
	return &images, nil
}

func (imgMgr *ImageMgr) CreateSnapshot(volID, snapshotName string) (*SnapshotResponse, error) {
	img, _, err := imgMgr.getImage(volID)
	if err != nil {
		log.Println("[Error] When get image:", err)
		return &SnapshotResponse{}, err
	}

	if err = img.Open(); err != nil {
		log.Println("[Error] When open image:", err)
		return &SnapshotResponse{}, err
	}
	defer img.Close()

	name := NewName(snapshotName)
	if _, err = img.CreateSnapshot(name.GetFullName()); err != nil {
		log.Println("[Error] When create snapshot:", err)
		return &SnapshotResponse{}, err
	}
	return &SnapshotResponse{
		Name:      name.GetName(),
		ID:        name.GetUUID(),
		Size:      imgMgr.getSize(img),
		Volume_id: volID,
	}, nil
}

func (imgMgr *ImageMgr) RemoveSnapshot(id string) error {
	imageNames, err := rbd.GetImageNames(imgMgr.Ioctx)
	if err != nil {
		log.Println("[Error] When getImageNames:", err)
		return err
	}

	var (
		snapInfo rbd.SnapInfo
		img      *rbd.Image
	)
EXIT:
	for _, name := range imageNames {
		in := ParseName(name)
		if in == nil {
			continue
		}
		img = rbd.GetImage(imgMgr.Ioctx, name)
		if err = img.Open(); err != nil {
			log.Println("[Error] When open image:", err)
			return err
		}
		snapInfos, err := img.GetSnapshotNames()
		img.Close()
		if err != nil {
			log.Println("[Error] When GetSnapshotNames:", err)
			continue
		}
		for _, snapInfo = range snapInfos {
			name := ParseName(snapInfo.Name)
			if id == name.GetUUID() {
				break EXIT
			}
		}
	}

	if err = img.Open(snapInfo.Name); err != nil {
		log.Println("[Error] When open image:", err)
	}
	defer img.Close()

	snapshot := img.GetSnapshot(snapInfo.Name)
	if err = snapshot.Remove(); err != nil {
		log.Println("[Error] When remove snapshot:", err)
		return err
	}
	log.Printf("[Info] Delete snapshot {%s} success", ParseName(snapInfo.Name).GetUUID())
	return nil
}

func (imgMgr *ImageMgr) GetSnapshot(id string) (*SnapshotResponse, error) {
	snapshots, err := imgMgr.GetSnapshots()
	if err != nil {
		return &SnapshotResponse{}, err
	}
	for _, snapshot := range *snapshots {
		if snapshot.ID == id {
			return &snapshot, nil
		}
	}
	return &SnapshotResponse{}, rbd.RbdErrorNotFound
}

func (imgMgr *ImageMgr) GetSnapshots() (*[]SnapshotResponse, error) {
	imageNames, err := rbd.GetImageNames(imgMgr.Ioctx)
	if err != nil {
		log.Println("[Error] When getImageNames:", err)
		return &[]SnapshotResponse{}, err
	}

	var snapshots []SnapshotResponse
	for _, name := range imageNames {
		in := ParseName(name)
		if in == nil {
			continue
		}
		img := rbd.GetImage(imgMgr.Ioctx, name)
		if err = img.Open(); err != nil {
			log.Println("[Error] When open image:", err)
			return &[]SnapshotResponse{}, err
		}
		snapInfos, err := img.GetSnapshotNames()
		img.Close()
		if err != nil {
			log.Println("[Error] When GetSnapshotNames:", err)
			continue
		}
		for _, snapInfo := range snapInfos {
			name := ParseName(snapInfo.Name)
			snapshot := SnapshotResponse{
				Name:      name.GetName(),
				ID:        name.GetUUID(),
				Size:      int64(snapInfo.Size >> SIZE_SHIFT_BIT),
				Volume_id: in.ID,
			}
			snapshots = append(snapshots, snapshot)
		}
	}
	return &snapshots, nil
}
func execCmd(cmd string) (string, error) {
	ret, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return string(ret[:len(ret)-1]), nil
}

func (imgMgr *ImageMgr) parseCapStr(cap string) int64 {
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
		log.Println("[Error] Cannot convert this number", err)
		return 0
	}
	if val, ok := UnitMapper[unit]; ok {
		return num << val >> SIZE_SHIFT_BIT
	} else {
		log.Println("[Error] strage unit is not found.")
		return 0
	}
}

func (imgMgr *ImageMgr) getPoolsCapInfo() ([][]string, error) {
	const poolStartLine = 5
	output, err := execCmd("ceph df")
	if err != nil {
		log.Println("[Error]:", err)
		return nil, err
	}
	lines := strings.Split(output, "\n")
	var poolsInfo [][]string
	for i := poolStartLine; i < len(lines); i++ {
		poolsInfo = append(poolsInfo, strings.Fields(lines[i]))
	}
	return poolsInfo, nil
}

func (imgMgr *ImageMgr) getGlobalCapInfo() ([]string, error) {
	const globalCapInfoLine = 2
	output, err := execCmd("ceph df")
	if err != nil {
		log.Println("[Error]:", err)
		return nil, err
	}
	lines := strings.Split(output, "\n")
	return strings.Fields(lines[globalCapInfoLine]), nil
}

func (imgMgr *ImageMgr) getPoolsAttr() (map[string][]string, error) {
	cmd := "ceph osd pool ls detail | grep \"^pool\"| awk '{print $3, $4, $6, $10}'"
	output, err := execCmd(cmd)
	if err != nil {
		log.Println("[Error]:", err)
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

func (imgMgr *ImageMgr) getPoolParam(line []string) map[string]interface{} {
	param := make(map[string]interface{})
	param["redundancyType"] = line[POOL_TYPE]
	if param["redundancyType"] == "replicated" {
		param["replicateSize"] = line[POOL_TYPE_SIZE]
	} else {
		param["erasureSize"] = line[POOL_TYPE_SIZE]
	}
	param["crushRuleset"] = line[POOL_CRUSH_RULESET]
	return param
}

func (imgMgr *ImageMgr) ListPools() (*[]PoolResponse, error) {
	pc, err := imgMgr.getPoolsCapInfo()
	if err != nil {
		log.Println("[Error]:", err)
		return nil, err
	}
	gc, err := imgMgr.getGlobalCapInfo()
	if err != nil {
		log.Println("[Error]:", err)
		return nil, err
	}
	pa, err := imgMgr.getPoolsAttr()
	if err != nil {
		log.Println("[Error]:", err)
		return nil, err
	}

	var pools []PoolResponse
	for i := range pc {
		name := pc[i][POOL_NAME]
		param := imgMgr.getPoolParam(pa[name])
		totalCap := imgMgr.parseCapStr(gc[GLOBAL_SIZE])
		maxAvailCap := imgMgr.parseCapStr(pc[i][POOL_MAX_AVAIL])
		availCap := imgMgr.parseCapStr(gc[GLOBAL_AVAIL])
		pool := PoolResponse{
			Name:       name,
			ID:         uuid.NewV5(uuid.NamespaceOID, name).String(),
			Parameters: param,
			//if redundancy type is replicate, MAX AVAIL =  AVAIL / relicate number,
			//and it this is erasure, MAX AVAIL =  AVAIL * k / (m + k)
			TotalCapacity: totalCap * maxAvailCap / availCap,
			FreeCapacity:  maxAvailCap,
		}
		pools = append(pools, pool)
	}
	return &pools, nil
}

type CephPlugin struct{}

func (plugin *CephPlugin) Setup() {}

func (plugin *CephPlugin) Unset() {}

func (plugin *CephPlugin) CreateVolume(name string, size int64) (*api.VolumeSpec, error) {
	var imgMgr = &ImageMgr{}
	if err := imgMgr.Init(); err != nil {
		log.Println("[Error] Connect ceph error.")
		return &api.VolumeSpec{}, err
	}

	defer imgMgr.Destory()

	vol, err := imgMgr.CreateImage(name, size)
	if err != nil {
		log.Println("[Error] When create volume:", err)
		return &api.VolumeSpec{}, err
	}

	log.Println("[Info] Create volume success, dls =", vol)
	return &api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id: vol.Id,
		},
		Name:             vol.Name,
		Size:             vol.Size,
		Description:      vol.Description,
		AvailabilityZone: vol.Availability_zone,
	}, nil
}

func (plugin *CephPlugin) GetVolume(volID string) (*api.VolumeSpec, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	vol, err := imgMgr.GetImage(volID)
	if err != nil {
		log.Println("[Error] When get volume:", err)
		return &api.VolumeSpec{}, err
	}

	log.Println("[Info] Get volume success, dls =", vol)
	return &api.VolumeSpec{
		BaseModel: &api.BaseModel{
			Id: vol.Id,
		},
		Name:             vol.Name,
		Size:             vol.Size,
		Description:      vol.Description,
		AvailabilityZone: vol.Availability_zone,
	}, nil
}

func (plugin *CephPlugin) DeleteVolume(volID string) error {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	if err := imgMgr.RemoveImage(volID); err != nil {
		log.Println("[Error] When delete volume:", err)
		return err
	}
	return nil
}

func (plugin *CephPlugin) InitializeConnection(volID string, doLocalAttach, multiPath bool, hostInfo *api.HostInfo) (*api.ConnectionInfo, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	img, err := imgMgr.GetImage(volID)
	if err != nil {
		log.Println("[Error] When get image:", err)
		return nil, err
	}

	return &api.ConnectionInfo{
		DriverVolumeType: "rbd",
		ConnectionData: map[string]interface{}{
			"secret_type":  "ceph",
			"name":         "rbd/" + OPENSDS_PREFIX + ":" + img.Name + ":" + img.Id,
			"cluster_name": "ceph",
			"hosts":        []string{hostInfo.Host},
			"volume_id":    img.Id,
			"access_mode":  "rw",
			"ports":        []string{"6789"},
		},
	}, nil
}

func (plugin *CephPlugin) AttachVolume(volID, host, mountpoint string) error {
	return nil
}

func (plugin *CephPlugin) DetachVolume(volID string) error {
	return nil
}

func (plugin *CephPlugin) CreateSnapshot(name, volID, description string) (*api.VolumeSnapshotSpec, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	snapshot, err := imgMgr.CreateSnapshot(volID, name)
	if err != nil {
		log.Println("[Error] When create snapshot:", err)
		return &api.VolumeSnapshotSpec{}, err
	}

	log.Println("[Info] Create snapshot success, dls =", snapshot)
	return &api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{
			Id: snapshot.ID,
		},
		Name:        snapshot.Name,
		Description: description,
		VolumeId:    volID,
		Size:        snapshot.Size,
	}, nil
}

func (plugin *CephPlugin) GetSnapshot(snapID string) (*api.VolumeSnapshotSpec, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	snapshot, err := imgMgr.GetSnapshot(snapID)
	if err != nil {
		log.Println("[Error] When get snapshot:", err)
		return &api.VolumeSnapshotSpec{}, err
	}

	log.Println("[Info] Get volume snapshot success, dls =", snapshot)
	return &api.VolumeSnapshotSpec{
		BaseModel: &api.BaseModel{
			Id: snapshot.ID,
		},
		Name: snapshot.Name,
		Size: snapshot.Size,
	}, nil
}

func (plugin *CephPlugin) DeleteSnapshot(snapID string) error {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	if err := imgMgr.RemoveSnapshot(snapID); err != nil {
		log.Println("[Error] When delete snapshot:", err)
		return err
	}
	return nil
}

func (plugin *CephPlugin) ListPools() ([]api.StoragePoolSpec, error) {
	var imgMgr = &ImageMgr{}

	var poolList []api.StoragePoolSpec
	poolsResp, err := imgMgr.ListPools()
	if err != nil {
		log.Println("[Error] When get snapshot:", err)
		return nil, err
	}
	for _, pl := range *poolsResp {
		pool := api.StoragePoolSpec{
			BaseModel: &api.BaseModel{
				Id: pl.ID,
			},
			Name:             pl.Name,
			Description:      pl.Description,
			AvailabilityZone: pl.AvailabilityZone,
			TotalCapacity:    pl.TotalCapacity,
			FreeCapacity:     pl.FreeCapacity,
			StorageType:      pl.StorageType,
			Parameters:       pl.Parameters,
		}
		poolList = append(poolList, pool)
	}
	return poolList, nil
}
