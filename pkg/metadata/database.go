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
This module implements some operations to database resource.

*/

package metadata

import (
	"encoding/json"
	"strings"

	"github.com/opensds/opensds/pkg/metadata/plugins/config"
)

func CreateDatabase(name string, size int) (string, error) {
	result := "Create database success!"
	return result, nil
}

func GetDatabase(id int, name string) (string, error) {
	dbInfo := config.DbInfo{}

	switch name {
	case "sqlBtree":
		dbInfo = *dbInfo.GetSqlBtreeInfo()
	case "sqlHash":
		dbInfo = *dbInfo.GetSqlHashInfo()
	case "mongodb":
		dbInfo = *dbInfo.GetMongodbInfo()
	case "rocksdb":
		dbInfo = *dbInfo.GetRocksdbInfo()
	default:
		return "Null", nil
	}

	a, _ := json.Marshal(dbInfo)
	result := string(a)
	return result, nil
}

func GetAllDatabases() (string, error) {
	dbInfo := config.DbInfo{}
	resInfo := make([]config.DbInfo, 0, 4)
	resInfo = append(resInfo, *dbInfo.GetSqlBtreeInfo())
	resInfo = append(resInfo, *dbInfo.GetSqlHashInfo())
	resInfo = append(resInfo, *dbInfo.GetMongodbInfo())
	resInfo = append(resInfo, *dbInfo.GetRocksdbInfo())

	dbSlice := make([]string, 4, 8)
	for i, _ := range resInfo {
		a, _ := json.Marshal(resInfo[i])
		dbSlice[i] = string(a)
	}
	result := strings.Join(dbSlice[:], ",")
	return result, nil
}

func UpdateDatabase(id int, size int, name string) (string, error) {
	result := "Update database success!"
	return result, nil
}

func DeleteDatabase(id int, name string, cascade bool) (string, error) {
	result := "Delete database success!"
	return result, nil
}
