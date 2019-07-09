// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package eternus

type System struct {
	Id          string `json:"ID"`
	Name        string `json:"NAME"`
	Location    string `json:"LOCATION"`
	ProductMode string `json:"PRODUCTMODE"`
	Wwn         string `json:"wwn"`
}

// StoragePool is Pool information
type StoragePool struct {
	Id            string
	Name          string
	TotalCapacity int64
	FreeCapacity  int64
}

// Volume is Pool information
type Volume struct {
	Id            string
	Name          string
	Status        string
	Size          int64
	TotalCapacity int64
	FreeCapacity  int64
	PoolName      string
	PoolId        string
}

// IscsiPortInfo is iscsi port info
type IscsiPortInfo struct {
	PortNumber     string
	IscsiName      string
	Ip             string
	TcpPort        int
	IsnsServerIp   string
	IsnsServerPort int
}

// FcPortInfo is iscsi port info
type FcPortInfo struct {
	PortNumber string
	Wwpn       string
}

// LunGroup is lun group info
type LunGroup struct {
	Volumes []LunGroupVolume
}

// LunGroupVolume is lun group info
type LunGroupVolume struct {
	Id          string
	Name        string
	RawStatus   string
	RoundStatus string
	Size        int64
	Uid         string
	Lun         string
}

// Mapping
type Mapping struct {
	Lun               string
	VolumeNumber      string
	VolumeName        string
	VolumeRawStatus   string
	VolumeRoundStatus string
	VolumeSize        int64
}

// SnapShot
type SnapShot struct {
	Sid        string
	Gen        string
	GenTotal   string
	Type       string
	VolumeType string
	SrcNo      string
	SrcName    string
	DestNo     string
	DestName   string
	Status     string
	Phase      string
	ErrorCode  string
	Requestor  string
}
