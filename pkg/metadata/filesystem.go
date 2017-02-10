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
This module implements some operations to file system resource.

*/

package metadata

import (
	"encoding/json"
	"strings"

	"github.com/opensds/opensds/pkg/metadata/plugins/config"
)

func CreateFileSystem(name string, size int) (string, error) {
	result := "Create file system success!"
	return result, nil
}

func GetFileSystem(id int, name string) (string, error) {
	fsInfo := config.FsInfo{}

	switch name {
	case "ext4":
		fsInfo = *fsInfo.GetExt4Info()
	case "btrfs":
		fsInfo = *fsInfo.GetBtrfsInfo()
	case "xfs":
		fsInfo = *fsInfo.GetXfsInfo()
	default:
		return "Null", nil
	}

	a, _ := json.Marshal(fsInfo)
	result := string(a)
	return result, nil
}

func GetAllFileSystems() (string, error) {
	fsInfo := config.FsInfo{}
	resInfo := make([]config.FsInfo, 0, 3)
	resInfo = append(resInfo, *fsInfo.GetExt4Info())
	resInfo = append(resInfo, *fsInfo.GetBtrfsInfo())
	resInfo = append(resInfo, *fsInfo.GetXfsInfo())

	fsSlice := make([]string, 3, 6)
	for i, _ := range resInfo {
		a, _ := json.Marshal(resInfo[i])
		fsSlice[i] = string(a)
	}
	result := strings.Join(fsSlice[:], ",")
	return result, nil
}

func UpdateFileSystem(id int, size int, name string) (string, error) {
	result := "Update file system success!"
	return result, nil
}

func DeleteFileSystem(id int, name string, cascade bool) (string, error) {
	result := "Delete file system success!"
	return result, nil
}
