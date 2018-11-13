// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package fusionstorage

type EncryptOpts struct {
	cmkId     string
	authToken string
}

type SnapshotResp struct {
	Name           string `fsc:"snap_name"`
	FatherName     string `fsc:"father_name"`
	Status         int    `fsc:"status"`
	Size           int64  `fsc:"snap_size"`
	RealSize       int64  `fsc:"real_size"`
	PoolId         string `fsc:"pool_id"`
	DeletePriority int    `fsc:"delete_priority"`
	CreateTime     int64  `fsc:"create_time"`
	EncryptFlag    bool   `fsc:"encrypt_flag"`
	SmartCacheFlag bool   `fsc:"smartCacheFlag"`
	TreeId         string `fsc:"tree_id"`
	BranchId       string `fsc:"branch_id"`
	SnapId         string `fsc:"snap_id"`
}

type VolumeResp struct {
	Name        string `fsc:"vol_name"`
	FatherName  string `fsc:"father_name"`
	Status      int    `fsc:"status"`
	Size        int64  `fsc:"vol_size"`
	RealSize    int64  `fsc:"real_size"`
	PoolId      string `fsc:"pool_id"`
	CreateTime  int64  `fsc:"create_time"`
	EncryptFlag bool   `fsc:"encrypt_flag"`
	LunId       string `fsc:"lun_id"`
	LLDProgress int    `fsc:"lld_progress"`
	RWRight     int    `fsc:"rw_right"`
	WWN         int    `fsc:"wwn"`
}

type PoolResp struct {
	PoolId        string `fsc:"pool_id"`
	TotalCapacity int64  `fsc:"total_capacity"`
	AllocCapacity int64  `fsc:"alloc_capacity"`
	UsedCapacity  int64  `fsc:"used_capacity"`
	PoolModel     int64  `fsc:"pool_model"`
}
