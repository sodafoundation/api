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

type Error struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type DeleteError struct {
	Error `json:"error"`
}

type GenericResult struct {
	Data  interface{} `json:"data"`
	Error Error       `json:"error"`
}

type Auth struct {
	AccountState  int    `json:"accountstate"`
	DeviceId      string `json:"deviceid"`
	IBaseToken    string `json:"iBaseToken"`
	LastLoginIp   string `json:"lastloginip"`
	LastLoginTime int    `json:"lastlogintime"`
	Level         int    `json:"level"`
	PwdChanGeTime int    `json:"pwdchangetime"`
	UserGroup     string `json:"usergroup"`
	UserId        string `json:"userid"`
	UserName      string `json:"username"`
	UserScope     string `json:"userscope"`
}

type AuthResp struct {
	Data  Auth  `json:"data"`
	Error Error `json:"error"`
}

type StoragePool struct {
	Description       string `json:"DESCRIPTION"`
	Id                string `json:"ID"`
	Name              string `json:"NAME"`
	UserFreeCapacity  string `json:"USERFREECAPACITY"`
	UserTotalCapacity string `json:"USERTOTALCAPACITY"`
}

type StoragePoolsResp struct {
	Data  []StoragePool `json:"data"`
	Error Error         `json:"error"`
}

type FileSystem struct {
	FileSystemData `json:"data"`
	Error          Error `json:"error"`
}

type FileSystemList struct {
	Data  []FileSystemData `json:"data"`
	Error Error            `json:"error"`
}

type FileSystemData struct {
	HealthStatus  string `json:"HEALTHSTATUS"`
	RunningStatus string `json:"RUNNINGSTATUS"`
	ID            string `json:"ID"`
	Capacity      string `json:"CAPACITY"`
	PoolName      string `json:"POOLNAME"`
	AllocType     string `json:"ALLOCTYPE"`
	Name          string `json:"NAME"`
}

type NFSShare struct {
	NFSShareData `json:"data"`
	Error        `json:"error"`
}

type NFSShareData struct {
	Description       string `json:"DESCRIPTION"`
	FSID              string `json:"FSID"`
	ID                string `json:"ID"`
	SharePath         string `json:"SHAREPATH"`
	LockPolicy        string `json:"LOCKPOLICY"`
	Name              string `json:"NAME"`
	CharacterEncoding string `json:"CHARACTERENCODING"`
}

type CIFSShare struct {
	CIFSShareData `json:"data"`
	Error         `json:"error"`
}

type CIFSShareData struct {
	Description               string `json:"DESCRIPTION"`
	FSID                      string `json:"FSID"`
	ID                        string `json:"ID"`
	SharePath                 string `json:"SHAREPATH"`
	Name                      string `json:"NAME"`
	AbeEnable                 string `json:"ABEENABLE"`
	EnableCA                  string `json:"ENABLECA"`
	EnableFileExtensionFilter string `json:"ENABLEFILEEXTENSIONFILTER"`
	EnableNotify              string `json:"ENABLENOTIFY"`
	EnableOpLock              string `json:"ENABLEOPLOCK"`
	EnableIPControl           string `json:"ENABLEIPCONTROL"`
	OfflineFileMode           string `json:"OFFLINEFILEMODE"`
	ApplyDefaultACL           string `json:"APPLYDEFAULTACL"`
}

type NFSShareList struct {
	Data  []NFSShareData `json:"data"`
	Error `json:"error"`
}

type CIFSShareList struct {
	Data  []CIFSShareData `json:"data"`
	Error `json:"error"`
}

type FSSnapshotData struct {
	Type            int    `json:"TYPE"`
	ID              string `json:"ID"`
	Name            string `json:"NAME"`
	ConsumeCapacity string `json:"CONSUMEDCAPACITY"`
	HealthStatus    string `json:"HEALTHSTATUS"`
	ParentID        string `json:"PARENTID"`
	ParentName      string `json:"PARENTNAME"`
	ParentType      int    `json:"PARENTTYPE"`
	Capacity        string `json:"USERCAPACITY"`
}

type FSSnapshot struct {
	FSSnapshotData `json:"data"`
	Error          `json:"error"`
}

type FSSnapshotList struct {
	Data  []FSSnapshotData `json:"data"`
	Error `json:"error"`
}
