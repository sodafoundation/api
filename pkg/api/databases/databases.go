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
This module implements the entry into CRUD operation of databases.

*/

package databases

import (
	"github.com/opensds/opensds/pkg/api/grpcapi"
)

func Create(name string, size int) (string, error) {
	result, err := grpcapi.CreateDatabase(name, size)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func Show(id int, name string) (string, error) {
	result, err := grpcapi.GetDatabase(id, name)

	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

func List() (string, error) {
	result, err := grpcapi.GetAllDatabases()

	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}

func Update(id int, size int, name string) (string, error) {
	result, err := grpcapi.UpdateDatabase(id, size, name)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}

func Delete(id int, name string, cascade bool) (string, error) {
	result, err := grpcapi.DeleteDatabase(id, name, cascade)

	if err != nil {
		return "Error", err
	} else {
		return result, nil
	}
}
