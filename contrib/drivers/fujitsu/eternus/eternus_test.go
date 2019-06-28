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

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	uuid "github.com/satori/go.uuid"

	mock "github.com/stretchr/testify/mock"
)

func checkArg(actual string, expected string) bool {
	actualArr := deleteEmptyStr(strings.Split(actual, " "))
	expectedArr := deleteEmptyStr(strings.Split(expected, " "))
	// same number of element
	if len(actualArr) != len(expectedArr) {
		return false
	}

	// same command and args
	for i := 0; i < len(actualArr); i = i + 2 {
		match := false
		for k := 0; k < len(expectedArr); k = k + 2 {
			if actualArr[i] == expectedArr[k] &&
				actualArr[i+1] == expectedArr[k+1] {
				match = true
				break
			}
		}
		if match {
			continue
		}
		return false
	}

	return true
}
func deleteEmptyStr(arr []string) []string {
	ret := []string{}
	for _, v := range arr {
		if v != "" && v != "\n" {
			ret = append(ret, v)
		}
	}
	return ret
}

func TestListPools(t *testing.T) {

	execCmd := "show thin-pro-pools\n"
	outStr := "\r\nCLI> show thin-pro-pools\r\n"
	outStr += "00\r\n"
	outStr += "00000000\r\n"
	outStr += "0004\r\n"
	outStr += "0011\tosdstest\t01\t00\t04\t0010\t00000000666FC000\t0000000000000000\t01\t5A\t4B\t00\t00\t00\t0000000000200000\t00000001\tFF\t01\r\n"
	outStr += "0012\tpoolname\t01\t00\t04\t0010\t00000000666FC000\t0000000000000000\t01\t5A\t4B\t00\t00\t00\t0000000000200000\t00000001\tFF\t01\r\n"
	outStr += "0013\tpoolname2\t01\t00\t04\t0010\t0000000019000000\t0000000018000000\t01\t5A\t4B\t00\t00\t00\t0000000000200000\t00000001\tFF\t01\r\n"
	outStr += "0014\tpoolname3\t01\t00\t04\t0010\t0000000000000000\t0000000000000000\t01\t5A\t4B\t00\t00\t00\t0000000000200000\t00000001\tFF\t01\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)

	d := &Driver{
		conf: &EternusConfig{
			AuthOptions: AuthOptions{
				Endpoint: "1.2.3.4",
			},
			Pool: map[string]PoolProperties{
				"poolname": PoolProperties{
					Extras: model.StoragePoolExtraSpec{
						DataStorage: model.DataStorageLoS{
							ProvisioningPolicy: "Thin",
						},
					},
					StorageType:      "block",
					AvailabilityZone: "az-test",
				},
				"poolname2": PoolProperties{
					Extras: model.StoragePoolExtraSpec{
						DataStorage: model.DataStorageLoS{
							ProvisioningPolicy: "Thin",
						},
					},
					StorageType: "block",
				},
			},
		},
		client: client,
	}

	ret, err := d.ListPools()
	if err != nil {
		t.Error("Test ListPools failed")
	}
	if len(ret) != 2 {
		t.Error("Test ListPools failed")
	}

	host, _ := os.Hostname()
	name := fmt.Sprintf("%s:%s:%s", host, d.conf.Endpoint, "18")
	id := uuid.NewV5(uuid.NamespaceOID, name).String()
	if ret[0].BaseModel.Id != id || ret[0].Name != "poolname" ||
		ret[0].TotalCapacity != 819 || ret[0].FreeCapacity != 819 ||
		ret[0].StorageType != "block" || ret[0].AvailabilityZone != "az-test" ||
		ret[0].Extras.DataStorage.ProvisioningPolicy != "Thin" {
		t.Error("Test ListPools failed")
	}
	name = fmt.Sprintf("%s:%s:%s", host, d.conf.Endpoint, "19")
	id = uuid.NewV5(uuid.NamespaceOID, name).String()
	if ret[1].BaseModel.Id != id || ret[1].Name != "poolname2" ||
		ret[1].TotalCapacity != 200 || ret[1].FreeCapacity != 8 ||
		ret[1].StorageType != "block" || ret[1].AvailabilityZone != "default" ||
		ret[1].Extras.DataStorage.ProvisioningPolicy != "Thin" {
		t.Error("Test ListPools failed")
	}
}

func TestCreateVolume(t *testing.T) {
	id := "volumeid"
	size := "1"
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	hashname := GetFnvHash(id)
	poolname := "poolname"

	opt := &pb.CreateVolumeOpts{
		Id:               id,
		Name:             "volumename",
		Size:             sizeInt,
		Description:      "test description",
		AvailabilityZone: "default",
		ProfileId:        "profileid",
		PoolId:           "poolid",
		PoolName:         poolname,
		Metadata:         map[string]string{},
		DriverName:       "drivername",
		Context:          "",
	}
	execCmd := "create volume -name " + hashname
	execCmd += " -size " + size + "gb"
	execCmd += " -pool-name " + poolname
	execCmd += " -type tpv -allocation thin \n"

	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "11\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmd)
			}),
	).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser := new(MockReadCloser)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			Pool: map[string]PoolProperties{
				"poolname": PoolProperties{
					Extras: model.StoragePoolExtraSpec{
						DataStorage: model.DataStorageLoS{
							ProvisioningPolicy: "Thin",
						},
					},
				},
			},
		},
		client: client,
	}

	ret, err := d.CreateVolume(opt)
	if err != nil {
		t.Error("Test CreateVolume failed")
	}
	if ret.BaseModel.Id != id || ret.Name != opt.Name ||
		ret.Size != sizeInt || ret.Description != opt.Description ||
		ret.AvailabilityZone != opt.AvailabilityZone ||
		ret.Metadata[KLunId] != "17" {
		t.Error("Test CreateVolume failed")
	}
}

func TestDeleteVolume(t *testing.T) {
	opt := &pb.DeleteVolumeOpts{
		Id: "id",
		Metadata: map[string]string{
			KLunId: "21",
		},
		DriverName: "drivername",
		Context:    "",
	}

	execCmd := "delete volume -volume-number 21 \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)

	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	err := d.DeleteVolume(opt)
	if err != nil {
		t.Error("Test DeleteVolume failed")
	}
}

func TestExtendVolume(t *testing.T) {
	id := "volumeid"
	lunid := "21"
	size := "2"
	sizeInt, _ := strconv.ParseInt(size, 10, 64)
	poolname := "poolname"

	opt := &pb.ExtendVolumeOpts{
		Id:               id,
		Name:             "volumename",
		Size:             sizeInt,
		Description:      "test description",
		AvailabilityZone: "default",
		ProfileId:        "profileid",
		PoolId:           "poolid",
		PoolName:         poolname,
		Metadata: map[string]string{
			KLunId: lunid,
		},
		DriverName: "drivername",
		Context:    "",
	}
	execCmd := "expand volume -volume-number " + lunid
	execCmd += " -size " + size + "gb \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmd)
			}),
	).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser := new(MockReadCloser)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	ret, err := d.ExtendVolume(opt)
	if err != nil {
		t.Error("Test ExtendVolume failed")
	}
	if ret.BaseModel.Id != id || ret.Name != opt.Name ||
		ret.Size != opt.Size || ret.Description != opt.Description ||
		ret.AvailabilityZone != opt.AvailabilityZone {
		t.Error("Test ExtendVolume failed")
	}
}

func TestInitializeConnection_IscsiNoPort(t *testing.T) {
	opt := &pb.CreateVolumeAttachmentOpts{
		Id:            "id",
		VolumeId:      "volumeid",
		DoLocalAttach: false,
		MultiPath:     false,
		HostInfo: &pb.HostInfo{
			Platform:  "linux",
			OsType:    "ubuntu",
			Host:      "hostname",
			Ip:        "1.1.1.1",
			Initiator: "iqn.testtest",
		},
		Metadata: map[string]string{
			KLunId: "1",
		},
		DriverName:     "drivername",
		Context:        "",
		AccessProtocol: "iscsi",
	}

	execCmd := "show iscsi-parameters\n"
	outStr := "CLI> show iscsi-parameters\r\n"
	outStr += "00\r\n"
	outStr += "04\r\n"
	outStr += "50\t00\t01\t00\tiqn.eternus-dx1\t\tFF\tDefault\t01\t00\t192.168.1.1\t255.255.255.0\t0.0.0.0\t000000000000\t0CBC\t02\t00000000\t0000\t0.0.0.0\t0C85\t\t00\t00\t00\t0001\t00\tFFFF\t0514\t0000\t\t::\t::\t::\tFF\tFF\t80000000\tFF\tFF\r\n"
	outStr += "50\t01\t00\t01\tiqn.eternus-dx2\t\tFF\tDefault\t01\t00\t192.166.1.2\t255.255.255.0\t0.0.0.0\t000000000000\t0CBC\t02\t00000000\t0000\t0.0.0.0\t0C85\t\t00\t00\t00\t0000\t00\tFFFF\t0514\t0000\t\t::\t::\t::\t00\t00\t80000000\t00\tFF\r\n"
	outStr += "51\t00\t04\t01\tiqn.eternus-dx3\t\t00\tDefault\t01\t00\t192.168.1.2\t255.255.255.0\t0.0.0.0\t000000000000\t0CBC\t02\t00000000\t0000\t0.0.0.0\t0C85\t\t00\t00\t00\t0000\t00\tFFFF\t0514\t0000\t\t::\t::\t::\tFF\tFF\t80000000\tFF\tFF\r\n"
	outStr += "51\t01\t01\t01\tiqn.eternus-dx4\t\tFF\tDefault\t01\t00\t192.166.1.4\t255.255.255.0\t0.0.0.0\t000000000000\t0CBC\t02\t00000000\t0000\t0.0.0.0\t0C85\t\t00\t00\t00\t0000\t00\tFFFF\t0514\t0000\t\t::\t::\t::\t00\t00\t80000000\t00\tFF\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)

	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	_, err := d.InitializeConnection(opt)
	if err == nil {
		t.Error("Test NewClient failed")
	}
}

func TestInitializeConnection_Iscsi(t *testing.T) {
	initiator := "iqn.testtest"
	hostname := "hostname"
	ipAddr := "1.1.1.1"
	hashhostname := GetFnvHash(initiator + ipAddr)
	opt := &pb.CreateVolumeAttachmentOpts{
		Id:            "id",
		VolumeId:      "volumeid",
		DoLocalAttach: false,
		MultiPath:     false,
		HostInfo: &pb.HostInfo{
			Platform:  "linux",
			OsType:    "ubuntu",
			Host:      hostname,
			Ip:        ipAddr,
			Initiator: initiator,
		},
		Metadata: map[string]string{
			KLunId: "21",
		},
		DriverName:     "drivername",
		Context:        "",
		AccessProtocol: ISCSIProtocol,
	}

	// Get iscsi port
	execCmd := "show iscsi-parameters\n"
	outStr := "CLI> show iscsi-parameters\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "50\t00\t00\t00\tiqn.eternus-dx1\t\tFF\tDefault\t01\t00\t192.168.1.1\t255.255.255.0\t0.0.0.0\t000000000000\t0CBC\t02\t00000000\t0000\t0.0.0.0\t0C85\t\t00\t00\t00\t0001\t00\tFFFF\t0514\t0000\t\t::\t::\t::\tFF\tFF\t80000000\tFF\tFF\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get iscsi host
	execCmd = "show host-iscsi-names\n"
	outStr = "\r\nCLI> show host-iscsi-names\r\n"
	outStr += "00\r\n"
	outStr += "0003\r\n"
	outStr += "0000\tHOST_NAME#0\t00\tDefault\t7F000001\tiqn.testtesttest\t00\r\n"
	outStr += "0001\tHOST_NAME#1\t00\tDefault\t7F000001\tiqn.testtesttesttest\t00\r\n"
	outStr += "0002\ttest_0\t00\tDefault\t02020202\tiqn.test\t00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(3, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(4, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Create iscsi host
	execCmdCreateHost := "create host-iscsi-name"
	execCmdCreateHost += " -name " + hashhostname
	execCmdCreateHost += " -ip " + ipAddr + " -ip-version ipv4"
	execCmdCreateHost += " -iscsi-name " + initiator + " \n"
	outStr = "\r\nCLI> create host-iscsi-name\r\n"
	outStr += "00\r\n"
	outStr += "11\r\n"
	outStr += "01\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdCreateHost)
			}),
	).Return(5, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(6, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get Lun group
	execCmd = "show lun-groups\n"
	outStr = "\r\nCLI> show lun-groups\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\ttest\tFFFF\tFFFF\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(7, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(8, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Create Lun group
	execCmdCreateLunGrp := "create lun-group"
	execCmdCreateLunGrp += " -name " + hashhostname
	execCmdCreateLunGrp += " -volume-number 21 -lun 0 \n"
	outStr = "\r\nCLI> create lun-group\r\n"
	outStr += "00\r\n"
	outStr += "12\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdCreateLunGrp)
			}),
	).Return(9, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(10, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Add host affinity
	execCmdSetHostAffinity := "set host-affinity \n"
	execCmdSetHostAffinity += " -port 010"
	execCmdSetHostAffinity += " -lg-number 18"
	execCmdSetHostAffinity += " -host-number 17 \n"
	outStr = "\r\nCLI> set host-affinity\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdSetHostAffinity)
			}),
	).Return(11, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get host lun
	execCmd = "show lun-group -lg-number 18 \n"
	outStr = "\r\nCLI> show lun-group -lg-number 18\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\ttest\tFFFF\tFFFF\r\n"
	outStr += "0003\r\n"
	outStr += "0000\t0014\tvolname1\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "0015\t0015\tvolname2\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "0002\t0016\tvolname3\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(13, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(14, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	ret, err := d.InitializeConnection(opt)
	if err != nil {
		t.Error("Test InitializeConnection failed")
	}
	connData := ret.ConnectionData
	if connData["targetIQN"] != "iqn.eternus-dx1" ||
		connData["targetPortal"] != "192.168.1.1:3260" ||
		connData["targetLun"] != "21" {
		t.Error("Test InitializeConnection failed")
	}
}

func TestInitializeConnection_FC(t *testing.T) {
	initiator := "AAAAAAAAAAAAAAAA"
	hostname := "hostname"
	ipAddr := "1.1.1.1"
	hashhostname := GetFnvHash(initiator)
	opt := &pb.CreateVolumeAttachmentOpts{
		Id:            "id",
		VolumeId:      "volumeid",
		DoLocalAttach: false,
		MultiPath:     false,
		HostInfo: &pb.HostInfo{
			Platform:  "linux",
			OsType:    "ubuntu",
			Host:      hostname,
			Ip:        ipAddr,
			Initiator: initiator,
		},
		Metadata: map[string]string{
			KLunId: "21",
		},
		DriverName:     "drivername",
		Context:        "",
		AccessProtocol: FCProtocol,
	}

	// Get iscsi port
	execCmd := "show fc-parameters\n"
	outStr := "CLI> show fc-parameters\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "40\t00\t04\t01\tFF\tFF\t00\t0800\t00\tFF\t\t01\t00\t00\t00\t00\tFF\t0000000000000001\t0000000000000000\tFF\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get iscsi host
	execCmd = "show host-wwn-names\n"
	outStr = "\r\nCLI> show host-wwn-names\r\n"
	outStr += "00\r\n"
	outStr += "0003\r\n"
	outStr += "0000\tHOST_NAME#0\t1234567890123456\t0000\tDefault\r\n"
	outStr += "0001\tHOST_NAME#1\t1234567890123457\t0000\tDefault\r\n"
	outStr += "0002\tHOST_NAME#2\t1234567890123458\t0000\tDefault\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(3, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(4, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Create iscsi host

	execCmdCreateHost := "create host-wwn-name"
	execCmdCreateHost += " -name " + hashhostname
	execCmdCreateHost += " -wwn " + initiator + " \n"
	outStr = "\r\nCLI> create host-wwn-name\r\n"
	outStr += "00\r\n"
	outStr += "11\r\n"
	outStr += "01\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdCreateHost)
			}),
	).Return(5, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(6, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get Lun group
	execCmd = "show lun-groups\n"
	outStr = "\r\nCLI> show lun-groups\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\ttest\tFFFF\tFFFF\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(7, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(8, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Create Lun group
	execCmdCreateLunGrp := "create lun-group"
	execCmdCreateLunGrp += " -name " + hashhostname
	execCmdCreateLunGrp += " -volume-number 21 -lun 0 \n"
	outStr = "\r\nCLI> create lun-group\r\n"
	outStr += "00\r\n"
	outStr += "12\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdCreateLunGrp)
			}),
	).Return(9, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(10, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Add host affinity
	execCmdSetHostAffinity := "set host-affinity \n"
	execCmdSetHostAffinity += " -port 000"
	execCmdSetHostAffinity += " -lg-number 18"
	execCmdSetHostAffinity += " -host-number 17 \n"
	outStr = "\r\nCLI> set host-affinity\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdSetHostAffinity)
			}),
	).Return(11, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get host lun
	execCmd = "show lun-group -lg-number 18 \n"
	outStr = "\r\nCLI> show lun-group -lg-number 18\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\ttest\tFFFF\tFFFF\r\n"
	outStr += "0003\r\n"
	outStr += "0000\t0014\tvolname1\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "0015\t0015\tvolname2\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "0002\t0016\tvolname3\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(13, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(14, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	ret, err := d.InitializeConnection(opt)
	if err != nil {
		t.Error("Test InitializeConnection failed")
	}
	connData := ret.ConnectionData
	if connData["targetWwn"] != "0000000000000001" ||
		connData["hostname"] != hostname ||
		connData["targetLun"] != "21" {

		t.Error("Test InitializeConnection failed")
	}
}

func TestInitializeConnection_FCNoInitiator(t *testing.T) {
	initiator := ""
	hostname := "hostname"
	ipAddr := "1.1.1.1"
	opt := &pb.CreateVolumeAttachmentOpts{
		Id:            "id",
		VolumeId:      "volumeid",
		DoLocalAttach: false,
		MultiPath:     false,
		HostInfo: &pb.HostInfo{
			Platform:  "linux",
			OsType:    "ubuntu",
			Host:      hostname,
			Ip:        ipAddr,
			Initiator: initiator,
		},
		Metadata: map[string]string{
			KLunId: "21",
		},
		DriverName:     "drivername",
		Context:        "",
		AccessProtocol: FCProtocol,
	}

	// Get iscsi port
	execCmd := "show fc-parameters\n"
	outStr := "CLI> show fc-parameters\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "40\t00\t04\t01\tFF\tFF\t00\t0800\t01\tFF\t\t01\t00\t00\t00\t00\tFF\t0000000000000001\t0000000000000000\tFF\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Show mapping
	execCmd = "show mapping -port 000 \n"
	outStr = "\r\nCLI> show mapping -port 000\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "11005100\t00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\t0001\tosds-643e8232-1b\tA000\t20\t0000000000200000\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Add mapping
	execCmdSetMapping := "set mapping"
	execCmdSetMapping += " -port 000"
	execCmdSetMapping += " -volume-number 21"
	execCmdSetMapping += " -lun 1 \n"
	outStr = "\r\nCLI> set mapping\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdSetMapping)
			}),
	).Return(11, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	ret, err := d.InitializeConnection(opt)
	if err != nil {
		t.Error("Test InitializeConnection failed")
	}
	connData := ret.ConnectionData
	if connData["targetWwn"] != "0000000000000001" ||
		connData["hostname"] != hostname ||
		connData["targetLun"] != "1" {
		t.Error("Test InitializeConnection failed")
	}
}

func TestTerminateConnection_Iscsi(t *testing.T) {
	initiator := "iqn.testtest"
	hostname := "hostname"
	ipAddr := "1.1.1.1"
	hashhostname := GetFnvHash(initiator + ipAddr)
	opt := &pb.DeleteVolumeAttachmentOpts{
		Id:       "id",
		VolumeId: "volumeid",
		HostInfo: &pb.HostInfo{
			Platform:  "linux",
			OsType:    "ubuntu",
			Host:      hostname,
			Ip:        ipAddr,
			Initiator: initiator,
		},
		Metadata: map[string]string{
			KLunId: "21",
		},
		DriverName:     "drivername",
		Context:        "",
		AccessProtocol: ISCSIProtocol,
	}

	// Get iscsi port
	execCmd := "show iscsi-parameters\n"
	outStr := "CLI> show iscsi-parameters\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "50\t00\t00\t00\tiqn.eternus-dx1\t\tFF\tDefault\t01\t00\t192.168.1.1\t255.255.255.0\t0.0.0.0\t000000000000\t0CBC\t02\t00000000\t0000\t0.0.0.0\t0C85\t\t00\t00\t00\t0001\t00\tFFFF\t0514\t0000\t\t::\t::\t::\tFF\tFF\t80000000\tFF\tFF\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get host lun
	execCmd = "show lun-group -lg-name " + hashhostname + " \n"
	outStr = "\r\nCLI> show lun-group -lg-name " + hashhostname + "\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\ttest\tFFFF\tFFFF\r\n"
	outStr += "0003\r\n"
	outStr += "0000\t0014\tvolname1\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "0015\t0015\tvolname2\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "0002\t0016\tvolname3\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Remove volume from lun group
	execCmd = "delete lun-group -lg-name " + hashhostname + " -lun 21 \n"
	outStr = "\r\nCLI> " + execCmd + "\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmd)
			}),
	).Return(3, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(4, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	err := d.TerminateConnection(opt)
	if err != nil {
		t.Error("Test TerminateConnection failed")
	}
}

func TestTerminateConnection_IscsiDLunGroup(t *testing.T) {
	initiator := "iqn.testtest"
	hostname := "hostname"
	ipAddr := "1.1.1.1"
	hashhostname := GetFnvHash(initiator + ipAddr)
	opt := &pb.DeleteVolumeAttachmentOpts{
		Id:       "id",
		VolumeId: "volumeid",
		HostInfo: &pb.HostInfo{
			Platform:  "linux",
			OsType:    "ubuntu",
			Host:      hostname,
			Ip:        ipAddr,
			Initiator: initiator,
		},
		Metadata: map[string]string{
			KLunId: "21",
		},
		DriverName:     "drivername",
		Context:        "",
		AccessProtocol: ISCSIProtocol,
	}

	// Get iscsi port
	execCmd := "show iscsi-parameters\n"
	outStr := "CLI> show iscsi-parameters\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "50\t00\t00\t00\tiqn.eternus-dx1\t\tFF\tDefault\t01\t00\t192.168.1.1\t255.255.255.0\t0.0.0.0\t000000000000\t0CBC\t02\t00000000\t0000\t0.0.0.0\t0C85\t\t00\t00\t00\t0001\t00\tFFFF\t0514\t0000\t\t::\t::\t::\tFF\tFF\t80000000\tFF\tFF\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get host lun
	execCmd = "show lun-group -lg-name " + hashhostname + " \n"
	outStr = "\r\nCLI> show lun-group -lg-name " + hashhostname + "\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\ttest\tFFFF\tFFFF\r\n"
	outStr += "0001\r\n"
	outStr += "0015\t0015\tvolname2\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Remove volume from lun group
	execCmdReleaseHostAffinity := "release host-affinity -port 010"
	execCmdReleaseHostAffinity += " -host-name " + hashhostname + " -mode all \n"
	outStr = "\r\nCLI> " + execCmd + "\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdReleaseHostAffinity)
			}),
	).Return(3, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(4, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Delete lun group
	execCmd = "delete lun-group -lg-name " + hashhostname + " \n"
	outStr = "\r\nCLI> " + execCmd + "\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(6, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Delete host
	execCmd = "delete host-iscsi-name -host-name " + hashhostname + " \n"
	outStr = "\r\nCLI> " + execCmd + "\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(6, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	err := d.TerminateConnection(opt)
	if err != nil {
		t.Error("Test TerminateConnection failed")
	}
}

func TestTerminateConnection_FcDLunGroup(t *testing.T) {
	initiator := "AAAAAAAAAAAAAAAA"
	hostname := "hostname"
	ipAddr := "1.1.1.1"
	hashhostname := GetFnvHash(initiator)
	opt := &pb.DeleteVolumeAttachmentOpts{
		Id:       "id",
		VolumeId: "volumeid",
		HostInfo: &pb.HostInfo{
			Platform:  "linux",
			OsType:    "ubuntu",
			Host:      hostname,
			Ip:        ipAddr,
			Initiator: initiator,
		},
		Metadata: map[string]string{
			KLunId: "21",
		},
		DriverName:     "drivername",
		Context:        "",
		AccessProtocol: FCProtocol,
	}

	// Get iscsi port
	execCmd := "show fc-parameters\n"
	outStr := "CLI> show fc-parameters\r\n"
	outStr += "00\r\n"
	outStr += "01\r\n"
	outStr += "40\t00\t04\t01\tFF\tFF\t00\t0800\t00\tFF\t\t01\t00\t00\t00\t00\tFF\t0000000000000001\t0000000000000000\tFF\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Get host lun
	execCmd = "show lun-group -lg-name " + hashhostname + " \n"
	outStr = "\r\nCLI> show lun-group -lg-name " + hashhostname + "\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\ttest\tFFFF\tFFFF\r\n"
	outStr += "0001\r\n"
	outStr += "0015\t0015\tvolname2\tA000\t20\t0000000000000000\t00000000000000000000000000000000\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Remove volume from lun group
	execCmdReleaseHostAffinity := "release host-affinity -port 000"
	execCmdReleaseHostAffinity += " -host-name " + hashhostname + " -mode all \n"
	outStr = "\r\nCLI> " + execCmd + "\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdReleaseHostAffinity)
			}),
	).Return(3, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(4, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Delete lun group
	execCmd = "delete lun-group -lg-name " + hashhostname + " \n"
	outStr = "\r\nCLI> " + execCmd + "\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(6, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Delete host
	execCmd = "delete host-wwn-name -host-name " + hashhostname + " \n"
	outStr = "\r\nCLI> " + execCmd + "\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(6, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	err := d.TerminateConnection(opt)
	if err != nil {
		t.Error("Test TerminateConnection failed")
	}
}

func TestAddMapping(t *testing.T) {
	// Show mapping
	execCmd := "show mapping -port 000 \n"
	outStr := "\r\nCLI> show mapping -port 000\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "11005100\t00\r\n"
	outStr += "0003\r\n"
	outStr += "0000\t0001\tvol1\tA000\t20\t0000000000200000\r\n"
	outStr += "0001\t0002\tvol2\tA000\t20\t0000000000200000\r\n"
	outStr += "0003\t0002\tvol2\tA000\t20\t0000000000200000\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Add mapping
	execCmdSetMapping := "set mapping"
	execCmdSetMapping += " -port 000"
	execCmdSetMapping += " -volume-number 21"
	execCmdSetMapping += " -lun 2 \n"
	outStr = "\r\nCLI> set mapping\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdSetMapping)
			}),
	).Return(11, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	ret, err := d.addMapping("000", "21")
	if err != nil {
		t.Error("Test addMapping failed")
	}
	if ret != "2" {
		t.Error("Test addMapping failed")
	}
}

func TestAddMapping_Max(t *testing.T) {
	// Show mapping
	execCmd := "show mapping -port 000 \n"
	outStr := "\r\nCLI> show mapping -port 000\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "11005100\t00\r\n"
	outStr += "0400\r\n"
	for i := 0; i < 1024; i++ {
		outStr += strconv.Itoa(i) + "\t0001\tvol1\tA000\t20\t0000000000200000\r\n"
	}
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	_, err := d.addMapping("000", "21")
	if err == nil {
		t.Error("Test addMapping failed")
	}
}

func TestDeleteMapping(t *testing.T) {
	// Show mapping
	execCmd := "show mapping -port 000 \n"
	outStr := "\r\nCLI> show mapping -port 000\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "11005100\t00\r\n"
	outStr += "0003\r\n"
	outStr += "0000\t0001\tvol1\tA000\t20\t0000000000200000\r\n"
	outStr += "0001\t0002\tvol2\tA000\t20\t0000000000200000\r\n"
	outStr += "0003\t0011\tvol2\tA000\t20\t0000000000200000\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	// Delete mapping
	execCmdSetMapping := "release mapping"
	execCmdSetMapping += " -port 000"
	execCmdSetMapping += " -lun 3 \n"
	outStr = "\r\nCLI> release mapping\r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	mockWriteCloser.On("Write",
		mock.MatchedBy(
			func(cmd []byte) bool {
				return checkArg(string(cmd), execCmdSetMapping)
			}),
	).Return(11, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(12, nil).Once()
	buff = make([]byte, 65535)
	out = []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	err := d.deleteMapping("000", "17")
	if err != nil {
		t.Error("Test deleteMapping failed")
	}
}

func TestDeleteMapping_Deleted(t *testing.T) {
	// Show mapping
	execCmd := "show mapping -port 000 \n"
	outStr := "\r\nCLI> show mapping -port 000\r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "11005100\t00\r\n"
	outStr += "0003\r\n"
	outStr += "0000\t0001\tvol1\tA000\t20\t0000000000200000\r\n"
	outStr += "0001\t0002\tvol2\tA000\t20\t0000000000200000\r\n"
	outStr += "0003\t0010\tvol2\tA000\t20\t0000000000200000\r\n"
	outStr += "CLI> "
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil).Once()
	mockWriteCloser.On("Write", []byte("\n")).Return(2, nil).Once()
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out).Once()

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	d := &Driver{
		conf: &EternusConfig{
			CeSupport: false,
		},
		client: client,
	}

	err := d.deleteMapping("000", "17")
	if err != nil {
		t.Error("Test deleteMapping failed")
	}
}
