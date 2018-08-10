// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package selector

import (
	"fmt"
	"reflect"
	"testing"

	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	dbtest "github.com/opensds/opensds/testutils/db/testing"
)

func TestSelectSupportedPool(t *testing.T) {
	mockClient := new(dbtest.Client)
	mockClient.On("ListPools", c.NewAdminContext()).Return(fakePools, nil)
	db.C = mockClient

	testCases := []struct {
		request  map[string]interface{}
		expected *model.StoragePoolSpec
	}{
		{
			request: map[string]interface{}{
				"freeCapacity":         ">= 5001",
				"availabilityZone":     "az1",
				"extras.advanced.thin": true,
			},
			expected: fakePools[1],
		},
		{
			request: map[string]interface{}{
				"freeCapacity":             ">= 400",
				"availabilityZone":         "s== default",
				"extras.advanced.diskType": "s== SSD",
			},
			expected: nil,
		},
	}

	s := NewSelector()
	for _, testCase := range testCases {
		result, _ := s.SelectSupportedPool(testCase.request)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

var (
	fakePools = []*model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
				Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
				CreatedAt: "2017-10-24T15:04:05",
			},
			Name:             "fakePool",
			Description:      "fake pool for testing",
			Status:           "available",
			AvailabilityZone: "az1",
			TotalCapacity:    99999,
			FreeCapacity:     5000,
			DockId:           "ccac4f33-e603-425a-8813-371bbe10566e",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					RecoveryTimeObjective: 1,
					ProvisioningPolicy:    "Thin",
					IsSpaceEfficient:      true,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "rbd",
					MaxIOPS:        1,
					MaxBWS:         10,
				},
				DataProtection: model.DataProtectionLos{},
				Advanced: model.ExtraSpec{
					"thin":     true,
					"dedupe":   false,
					"compress": false,
					"diskType": "SSD",
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
				CreatedAt: "2017-10-24T15:04:05",
			},
			Name:             "fakePool",
			Description:      "fake pool for testing",
			Status:           "available",
			AvailabilityZone: "az1",
			TotalCapacity:    99999,
			FreeCapacity:     6999,
			DockId:           "ccac4f33-e603-425a-8813-371bbe10566e",
			Extras: model.StoragePoolExtraSpec{
				Advanced: model.ExtraSpec{
					"thin":     true,
					"dedupe":   true,
					"compress": true,
					"diskType": "SATA",
				},
			},
		},
	}
)

var (
	PoolA = model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e11",
			CreatedAt: "2017-10-24T15:04:05",
		},
		FreeCapacity:     int64(50),
		AvailabilityZone: "az1",
		Extras: model.StoragePoolExtraSpec{
			Advanced: model.ExtraSpec{
				"thin":     true,
				"dedupe":   true,
				"compress": true,
				"diskType": "SSD",
			},
		},
	}
	PoolB = model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e12",
			CreatedAt: "2017-10-24T15:04:06",
		},
		FreeCapacity:     int64(60),
		AvailabilityZone: "az2",
		Extras: model.StoragePoolExtraSpec{
			Advanced: model.ExtraSpec{
				"thin":     false,
				"dedupe":   false,
				"compress": false,
				"diskType": "SATA",
			},
		},
	}
	PoolC = model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id:        "f4486139-78d5-462d-a7b9-fdaf6c797e13",
			CreatedAt: "2017-10-24T15:04:07",
		},
		FreeCapacity:     int64(70),
		AvailabilityZone: "az3",
		Extras: model.StoragePoolExtraSpec{
			Advanced: model.ExtraSpec{
				"thin":     true,
				"dedupe":   false,
				"compress": true,
				"diskType": "SSD",
			},
		},
	}
)

func TestSelectSupportedPools_00(t *testing.T) {
	request := make(map[string]interface{})
	request["extras.advanced.thin"] = true

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(1, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	delete(request, "extras.advanced.thin")
	request["freeCapacity"] = float64(70)
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolC, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolC, supportedPools[0])
	}
}

func TestSelectSupportedPools_01(t *testing.T) {
	request := make(map[string]interface{})
	request["extras.advanced.thin"] = 1

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "the type of extras.advanced.thin must be bool"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	fmt.Println(err.Error())
	delete(request, "extras.advanced.thin")
	delete(request, "availabilityZone")
	request["freeCapacity"] = float64(80)
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "no available pool to meet user's requirement"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_03(t *testing.T) {
	request := make(map[string]interface{})
	// bool:1、0、t、f、T、F、true、false、True、False、TRUE、FALSE
	request["extras.advanced.thin"] = "1"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	delete(request, "extras.advanced.thin")
	request["freeCapacity"] = "70"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolC, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolC, supportedPools[0])
	}

	delete(request, "freeCapacity")
	request["extras.advanced.diskType"] = "SATA"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolB, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolB, supportedPools[0])
	}
}

func TestSelectSupportedPools_04(t *testing.T) {
	request := make(map[string]interface{})
	request["extras.advanced.thin"] = "2"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "capability is: extras.advanced.thin, 2 is not bool"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	fmt.Println(err.Error())
	delete(request, "extras.advanced.thin")
	delete(request, "availabilityZone")
	request["freeCapacity"] = "80"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "no available pool to meet user's requirement"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	delete(request, "freeCapacity")
	request["extras.advanced.diskType"] = "SSD1"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "no available pool to meet user's requirement"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_05(t *testing.T) {
	request := make(map[string]interface{})
	request["freeCapacity"] = "<="

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "the format of freeCapacity: <= is incorrect"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	request["freeCapacity"] = ">= z"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "capability is: freeCapacity, z is not float64"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	request["freeCapacity"] = "== 1 2"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "the format of freeCapacity: == 1 2 is incorrect"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	delete(request, "freeCapacity")
	request["availabilityZone"] = "!= 50"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "the value of availabilityZone is not float64"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_06(t *testing.T) {
	request := make(map[string]interface{})

	request["freeCapacity"] = "!= 50"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolB, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolB, supportedPools[0])
	}

	request["freeCapacity"] = "<= 50"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	request["freeCapacity"] = ">= 70"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolC, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolC, supportedPools[0])
	}

	request["freeCapacity"] = "== 50"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}
}

func TestSelectSupportedPools_07(t *testing.T) {
	request := make(map[string]interface{})
	request["availabilityZone"] = "<in>"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "the format of availabilityZone: <in> is incorrect"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	request["availabilityZone"] = "<in> a az"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "the format of availabilityZone: <in> a az is incorrect"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	delete(request, "availabilityZone")
	request["freeCapacity"] = "<in> a"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "freeCapacity is not a string, so <in> can not be used"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_08(t *testing.T) {
	request := make(map[string]interface{})

	request["availabilityZone"] = "<in> az3"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolC, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolC, supportedPools[0])
	}

	request["availabilityZone"] = "<in> z1"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	request["availabilityZone"] = "<in> 2"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolB, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolB, supportedPools[0])
	}

}

func TestSelectSupportedPools_09(t *testing.T) {
	request := make(map[string]interface{})
	request["availabilityZone"] = "<or>"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "when using <or> as an operator, the <or> and value must appear in pairs"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	request["availabilityZone"] = "<or> az1 <in> az"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "the first operator is <or>, the following operators must be <or>"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_10(t *testing.T) {
	request := make(map[string]interface{})

	request["availabilityZone"] = "<or> az3"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolC, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolC, supportedPools[0])
	}

	delete(request, "availabilityZone")
	request["freeCapacity"] = "<or> 50 <or> 60"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	request["freeCapacity"] = "<or> 70 <or> 60"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolB, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolB, supportedPools[0])
	}

}

func TestSelectSupportedPools_11(t *testing.T) {
	request := make(map[string]interface{})
	request["extras.advanced.dedupe"] = "<is>"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "the format of extras.advanced.dedupe: <is> is incorrect"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	request["extras.advanced.dedupe"] = "<is> 2"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "capability is: extras.advanced.dedupe, 2 is not bool"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	delete(request, "extras.advanced.dedupe")
	request["freeCapacity"] = "<is> 1"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "the value of freeCapacity is not bool"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_12(t *testing.T) {
	request := make(map[string]interface{})

	request["extras.advanced.dedupe"] = "<is> t"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	request["extras.advanced.dedupe"] = "<is> f"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolB, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolB, supportedPools[0])
	}
}

func TestSelectSupportedPools_13(t *testing.T) {
	request := make(map[string]interface{})
	request["availabilityZone"] = "s=="

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "the format of availabilityZone: s== is incorrect"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}

	delete(request, "availabilityZone")
	request["extras.advanced.dedupe"] = "s== az"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	ExpectedErr = "extras.advanced.dedupeis not a string"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_14(t *testing.T) {
	request := make(map[string]interface{})

	request["availabilityZone"] = "s== az3"

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolC, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolC, supportedPools[0])
	}

	request["availabilityZone"] = "s!= z1"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	request["availabilityZone"] = "s>= az2"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolB, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolB, supportedPools[0])
	}

	request["availabilityZone"] = "s> az2"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolC, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolC, supportedPools[0])
	}

	request["availabilityZone"] = "s<= az2"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}

	request["availabilityZone"] = "s< az2"
	supportedPools, err = SelectSupportedPools(3, request, pools)
	if nil != err {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(&PoolA, supportedPools[0]) {
		t.Errorf("Expected %v, get %v", PoolA, supportedPools[0])
	}
}

func TestSelectSupportedPools_15(t *testing.T) {
	request := make(map[string]interface{})

	request["size"] = 50

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "no available pool to meet user's requirement"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}

func TestSelectSupportedPools_16(t *testing.T) {
	request := make(map[string]interface{})
	request["freeCapacity"] = true

	pools := []*model.StoragePoolSpec{
		&PoolA,
		&PoolB,
		&PoolC,
	}

	supportedPools, err := SelectSupportedPools(3, request, pools)
	ExpectedErr := "the type of freeCapacity must be float64"

	if ExpectedErr != err.Error() {
		t.Errorf("Expected %v, get %v", ExpectedErr, err)
	}

	if nil != supportedPools {
		t.Errorf("Expected %v, get %v", nil, supportedPools[0])
	}
}
