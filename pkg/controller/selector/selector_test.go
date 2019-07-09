// Copyright 2017 The OpenSDS Authors.
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

func TestSelectSupportedPoolForVolume(t *testing.T) {
	mockClient := new(dbtest.Client)
	mockClient.On("GetDefaultProfile", c.NewAdminContext()).Return(fakeProfiles[0], nil)
	mockClient.On("GetProfile", c.NewAdminContext(), "2f9c0a04-66ef-11e7-ade2-43158893e017").Return(fakeProfiles[1], nil)
	mockClient.On("GetProfile", c.NewAdminContext(), "c611ab76-b4a8-11e8-b76f-97665ba92921").Return(fakeProfiles[2], nil)
	mockClient.On("ListPools", c.NewAdminContext()).Return(fakePools, nil)
	db.C = mockClient

	testCases := []struct {
		request  *model.VolumeSpec
		expected *model.StoragePoolSpec
	}{
		{
			request: &model.VolumeSpec{
				Size:             40,
				AvailabilityZone: "az1",
				PoolId:           "f4486139-78d5-462d-a7b9-fdaf6c797e1b",
			},
			expected: fakePools[0],
		},
		{
			request: &model.VolumeSpec{
				Size:             5001,
				ProfileId:        "2f9c0a04-66ef-11e7-ade2-43158893e017",
				AvailabilityZone: "az1",
			},
			expected: fakePools[1],
		},
		{
			request: &model.VolumeSpec{
				Size:             400,
				ProfileId:        "c611ab76-b4a8-11e8-b76f-97665ba92921",
				AvailabilityZone: "default",
			},
			expected: nil,
		},
	}

	s := NewSelector()
	for _, testCase := range testCases {
		result, _ := s.SelectSupportedPoolForVolume(testCase.request)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Expected %v, get %v", testCase.expected, result)
		}
	}
}

var (
	fakeProfiles = []*model.ProfileSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:        "default",
			Description: "default policy",
			StorageType: "block",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "profile-01",
			Description: "silver policy",
			StorageType: "block",
			ProvisioningProperties: model.ProvisioningPropertiesSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					IsSpaceEfficient:   true,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "iscsi",
					MaxIOPS:        50000,
					MaxBWS:         500,
				},
			},
			ReplicationProperties: model.ReplicationPropertiesSpec{
				DataProtection: model.DataProtectionLoS{
					IsIsolated:  true,
					ReplicaType: "Clone",
				},
			},
			DataProtectionProperties: model.DataProtectionPropertiesSpec{
				DataProtection: model.DataProtectionLoS{
					IsIsolated:  true,
					ReplicaType: "Clone",
				},
				ConsistencyEnabled: true,
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "c611ab76-b4a8-11e8-b76f-97665ba92921",
			},
			Name:        "profile-02",
			Description: "silver policy",
			StorageType: "block",
			ProvisioningProperties: model.ProvisioningPropertiesSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					IsSpaceEfficient:   true,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "rbd",
					MaxIOPS:        500,
					MaxBWS:         500,
				},
			},
			ReplicationProperties: model.ReplicationPropertiesSpec{
				DataProtection: model.DataProtectionLoS{
					IsIsolated:  true,
					ReplicaType: "Clone",
				},
			},
			DataProtectionProperties: model.DataProtectionPropertiesSpec{
				DataProtection: model.DataProtectionLoS{
					IsIsolated:  true,
					ReplicaType: "Clone",
				},
				ConsistencyEnabled: true,
			},
			CustomProperties: model.CustomPropertiesSpec{
				"diskType": "SSD",
			},
		},
	}

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
			StorageType:      "block",
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
					AccessProtocol: "iscsi",
					MaxIOPS:        100000,
					MaxBWS:         1000,
				},
				DataProtection: model.DataProtectionLoS{},
				Advanced: map[string]interface{}{
					"compress": false,
					"diskType": "SSD",
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id:        "42a4c0e0-b497-11e8-a14c-3bb1bb1e8caf",
				CreatedAt: "2017-10-24T15:04:05",
			},
			Name:             "fakePool",
			Description:      "fake pool for testing",
			Status:           "available",
			StorageType:      "block",
			AvailabilityZone: "az1",
			TotalCapacity:    99999,
			FreeCapacity:     6999,
			DockId:           "ccac4f33-e603-425a-8813-371bbe10566e",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					RecoveryTimeObjective: 1,
					ProvisioningPolicy:    "Thin",
					IsSpaceEfficient:      true,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "iscsi",
					MaxIOPS:        60000,
					MaxBWS:         600,
				},
				DataProtection: model.DataProtectionLoS{
					IsIsolated:  true,
					ReplicaType: "Clone",
				},
				Advanced: map[string]interface{}{
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
			Advanced: map[string]interface{}{
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
			Advanced: map[string]interface{}{
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
			Advanced: map[string]interface{}{
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
