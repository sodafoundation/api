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
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"openstack/golang-client/volume"

	"github.com/ceph/go-ceph/rados"
	"github.com/ceph/go-ceph/rbd"
	"github.com/opensds/opensds/pkg/dock/plugins/connector"
	"github.com/satori/go.uuid"
)

const (
	OPENSDS_PREFIX string = "OPENSDS"
	SPLIT_CHAR            = ":"
	SIZE_SHIFT_BIT        = 20
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
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Status      string              `json:"status"`
	Size        int                 `json:"size"`
	Volume_type string              `json:"volume_type"`
	Attachments []map[string]string `json:"attachments"`
}

// SnapshotResponse is a structure for all properties of
// a volume snapshot for a non detailed query
type SnapshotResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Volume_id string `json:"volume_id"`
	Size      int    `json:"size"`
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

func (imgMgr *ImageMgr) CreateImage(name string, size int32) (*Response, error) {
	imageName := NewName(name)

	_, err := rbd.Create(imgMgr.Ioctx, imageName.GetFullName(), uint64(size)<<SIZE_SHIFT_BIT, 20)
	if err != nil {
		log.Println("[Error] When create rbd image:", err)
		return &Response{}, err
	}

	return &Response{
		Name:   imageName.GetName(),
		ID:     imageName.GetUUID(),
		Status: "available",
		Size:   int(size),
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

func (imgMgr *ImageMgr) getSize(img *rbd.Image) int {
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
	return int(size >> SIZE_SHIFT_BIT)
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
		Name:   name.GetName(),
		ID:     name.GetUUID(),
		Status: "available",
		Size:   imgMgr.getSize(img),
	}, nil
}

func (imgMgr *ImageMgr) GetImages(allowDetails bool) (*[]Response, error) {
	imageNames, err := rbd.GetImageNames(imgMgr.Ioctx)
	if err != nil {
		log.Println("[Error] When getImageNames:", err)
		return nil, err
	}

	var images []Response

	for _, name := range imageNames {
		in := ParseName(name)
		if in == nil {
			continue
		}
		img := rbd.GetImage(imgMgr.Ioctx, name)
		image := Response{
			Name:   in.GetName(),
			ID:     in.GetUUID(),
			Status: "available",
			Size:   imgMgr.getSize(img),
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
		Status:    "available",
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
	for _, snapshot := range snapshots {
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
		return nil, err
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
			return nil, err
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
				Status:    "available",
				Size:      int(snapInfo.Size >> SIZE_SHIFT_BIT),
				Volume_id: in.ID,
			}
			snapshots = append(snapshots, snapshot)
		}
	}
	return &snapshots, nil
}

type CephPlugin struct{}

func (plugin *CephPlugin) Setup() {}

func (plugin *CephPlugin) Unset() {}

func (plugin *CephPlugin) CreateVolume(name, volType string, size int32) (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	vol, err := imgMgr.CreateImage(name, size)
	if err != nil {
		log.Println("[Error] When create volume:", err)
		return "", err
	}

	a, _ := json.Marshal(vol)
	result := string(a)
	log.Println("[Info] Create volume success, dls =", result)
	return result, nil
}

func (plugin *CephPlugin) GetVolume(volID string) (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	vol, err := imgMgr.GetImage(volID)
	if err != nil {
		log.Println("[Error] When get volume:", err)
		return "", err
	}

	a, _ := json.Marshal(vol)
	result := string(a)
	log.Println("[Info] Get volume success, dls =", result)
	return result, nil
}

func (plugin *CephPlugin) GetAllVolumes(allowDetails bool) (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	vols, err := imgMgr.GetImages(allowDetails)
	if err != nil {
		log.Println("[Error] When get all volumes:", err)
		return "", err
	}

	a, _ := json.Marshal(vols)
	result := string(a)
	log.Println("[Info] Get all volumes success, dls =", result)
	return result, nil
}

func (plugin *CephPlugin) DeleteVolume(volID string) (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	if err := imgMgr.RemoveImage(volID); err != nil {
		log.Println("[Error] When delete volume:", err)
		return "", err
	}

	result := "Delete volume success!"
	return result, nil
}

func (plugin *CephPlugin) AttachVolume(volID string) (string, error) {
	conn, err := getConnectionInfo(volID)
	if err != nil {
		log.Println("[Error] When get connection info:", err)
		return "", err
	}

	log.Printf("[Info] Receive connection info: %+v\n", conn)

	devPath, err := conn.ConnectVolume()
	if err != nil {
		log.Println("[Error] When connect volume:", err)
		return "", err
	}
	return devPath, nil
}

func (plugin *CephPlugin) DetachVolume(device string) (string, error) {
	ind := strings.Index(device, "by-id/")
	if ind < 0 {
		return "", fmt.Errorf("[Error] Detach disk: no volume id in %s", device)
	}

	var volID = device[ind+6:]
	conn, err := getConnectionInfo(volID)
	if err != nil {
		log.Println("[Error] When get connection info:", err)
		return "", err
	}

	log.Printf("[Info] Receive connection info: %+v\n", conn)

	_, err = conn.DisconnectVolume()
	if err != nil {
		log.Println("[Error] When disconnect volume:", err)
		return "", err
	}
	return "Detach volume success!", nil
}

func (plugin *CephPlugin) CreateSnapshot(name, volID, description string, forced bool) (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	snapshot, err := imgMgr.CreateSnapshot(volID, name)
	if err != nil {
		log.Println("[Error] When create snapshot:", err)
		return "", err
	}

	a, _ := json.Marshal(snapshot)
	result := string(a)
	log.Println("[Info] Create snapshot success, dls =", result)
	return result, nil
}

func (plugin *CephPlugin) GetSnapshot(snapID string) (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	snapshot, err := imgMgr.GetSnapshot(snapID)
	if err != nil {
		log.Println("[Error] When get snapshot:", err)
		return "", err
	}

	a, _ := json.Marshal(snapshot)
	result := string(a)
	log.Println("[Info] Get volume snapshot success, dls =", result)
	return result, nil
}

func (plugin *CephPlugin) GetAllSnapshots() (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	snapshots, err := imgMgr.GetSnapshots()
	if err != nil {
		log.Println("[Error] When get all snapshots:", err)
		return "", err
	}

	a, _ := json.Marshal(snapshots)
	result := string(a)
	log.Println("[Info] Get all snapshots success, dls =", result)
	return result, nil
}

func (plugin *CephPlugin) DeleteSnapshot(snapID string) (string, error) {
	var imgMgr = &ImageMgr{}
	if imgMgr.Init() != nil {
		log.Println("[Error] When ceph connection")
	}

	defer imgMgr.Destory()

	if err := imgMgr.RemoveSnapshot(snapID); err != nil {
		log.Println("[Error] When delete snapshot:", err)
		return "", err
	}

	result := "Delete snapshot success!"
	return result, nil
}

type ConnectionCephData struct {
	SecretType   string   `json:"secret_type"`
	Name         string   `json:"name"`
	Encrypted    bool     `json:"encrypted"`
	ClusterName  string   `json:"cluster_name"`
	SecretUUID   string   `json:"secret_uuid"`
	QosSpes      string   `json:"qos_specs"`
	Hosts        []string `json:"hosts"`
	VolumeId     string   `json:"volume_id"`
	AuthEnabled  bool     `json:"auth_enabled"`
	AccessMode   string   `json:"access_mode"`
	AuthUsername string   `json:"auth_username"`
	Ports        []string `json:"ports"`
}

func getConnectionInfo(volID string) (*connector.Connector, error) {
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

	properties, err := connector.GetConnectorProperties(false)
	if err != nil {
		log.Println("[Error] When get connector properties:", err)
		return nil, err
	}

	return &connector.Connector{
		ConnInfo: volume.ConnectionInfo{
			DriverVolumeType: "rbd",
			ConnectionDataContainer: ConnectionCephData{
				SecretType:  "ceph",
				Name:        "rbd/" + OPENSDS_PREFIX + ":" + img.Name + ":" + img.ID,
				ClusterName: "ceph",
				Hosts:       []string{properties.Host},
				VolumeId:    img.ID,
				AccessMode:  "rw",
				Ports:       []string{"6789"},
			},
		},
	}, nil
}
