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

import (
	"fmt"
	//	"reflect"
	"testing"
	//	"github.com/opensds/opensds/contrib/drivers/utils/config"
	//	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	//	"github.com/opensds/opensds/pkg/model"
)

func TestSetup(t *testing.T) {
	driver := &Driver{}
	err := driver.Setup()
	if err != nil {
		fmt.Println(err)
	}
	pools, err := driver.ListPools()
	fmt.Printf("pool %+v %+v\n", pools[0].Name, err)
	//	err = driver.CreateShare("", pools[0].Name, "", "NFS")
	//	if err != nil {
	//		fmt.Println(err)
	//	}

	//	err = driver.DeleteShare("9", "CIFS", "51")
	//	if err != nil {
	//		fmt.Println(err)
	//	}

	//driver.CreateSnapshotFromShare("opensds", "NFS", "7")
	driver.ListAllSnapshots()
	driver.ShowFSSnapshot("60@share_snapshot_opensds")
	driver.DeleteFSSnapshot("60@share_snapshot_opensds")
}
