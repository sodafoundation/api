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
	"bytes"
	"io"
	"testing"

	mock "github.com/stretchr/testify/mock"
)

// MockWriteCloser
type MockWriteCloser struct {
	mock.Mock
}

// MockWriteCloser function
func (_m *MockWriteCloser) Write(in []byte) (int, error) {
	ret := _m.Called(in)
	return ret.Get(0).(int), ret.Error(1)
}
func (_m *MockWriteCloser) Close() error {
	ret := _m.Called()
	return ret.Error(0)
}

// MockWriteCloser
type MockReadCloser struct {
	mock.Mock
}

// MockWriteCloser function
func (_m *MockReadCloser) Read(in []byte) (int, error) {
	ret := _m.Called(in)
	out := ret.Get(2).([]byte)
	for i, v := range out {
		in[i] = v
	}
	return ret.Get(0).(int), ret.Error(1)
}
func (_m *MockReadCloser) Close() error {
	ret := _m.Called()
	return ret.Error(0)
}

func TestNewClient(t *testing.T) {
	authOptions := &AuthOptions{
		Username:        "testuser",
		Password:        "testpassword",
		AdminUsername:   "testadminuser",
		AdminPassword:   "testadminpassword",
		PwdEncrypter:    "",
		EnableEncrypted: false,
		Endpoint:        "endpoint",
		Insecure:        false,
	}
	c, err := NewClient(authOptions)
	if err != nil {
		t.Error("Test NewClient failed")
	}
	if c.user != authOptions.Username {
		t.Error("Test NewClient failed")
	}
	if c.password != authOptions.Password {
		t.Error("Test NewClient failed")
	}
	if c.endpoint != authOptions.Endpoint {
		t.Error("Test NewClient failed")
	}
}

func TestNewClient_Encryption(t *testing.T) {
	authOptions := &AuthOptions{
		Username:        "testuser",
		Password:        "d82cf4dd2523237a240b4d400e3bef67c694f1afce5d96a09e10cd7f",
		AdminUsername:   "testadminuser",
		AdminPassword:   "testadminpassword",
		PwdEncrypter:    "aes",
		EnableEncrypted: true,
		Endpoint:        "endpoint",
		Insecure:        false,
	}
	c, err := NewClient(authOptions)
	if err != nil {
		t.Error("Test NewClient failed")
	}
	if c.user != authOptions.Username {
		t.Error("Test NewClient failed")
	}
	if c.password != "testpassword" {
		t.Error("Test NewClient failed")
	}
	if c.endpoint != authOptions.Endpoint {
		t.Error("Test NewClient failed")
	}
}

func TestNewClient_Error(t *testing.T) {
	authOptions := &AuthOptions{
		Username:        "testuser",
		Password:        "d",
		AdminUsername:   "testadminuser",
		AdminPassword:   "testadminpassword",
		PwdEncrypter:    "aes",
		EnableEncrypted: true,
		Endpoint:        "endpoint",
		Insecure:        false,
	}
	_, err := NewClient(authOptions)
	if err == nil {
		t.Error("Test NewClient failed")
	}
}

func TestNewClientForAdmin(t *testing.T) {
	authOptions := &AuthOptions{
		Username:        "testuser",
		Password:        "testpassword",
		AdminUsername:   "testadminuser",
		AdminPassword:   "testadminpassword",
		PwdEncrypter:    "",
		EnableEncrypted: false,
		Endpoint:        "endpoint",
		Insecure:        false,
	}
	c, err := NewClientForAdmin(authOptions)
	if err != nil {
		t.Error("Test NewClientForAdmin failed")
	}
	if c.user != authOptions.AdminUsername {
		t.Error("Test NewClientForAdmin failed")
	}
	if c.password != authOptions.AdminPassword {
		t.Error("Test NewClientForAdmin failed")
	}
	if c.endpoint != authOptions.Endpoint {
		t.Error("Test NewClientForAdmin failed")
	}
}

func TestNewClientForAdmin_Encryption(t *testing.T) {
	authOptions := &AuthOptions{
		Username:        "testuser",
		Password:        "d82cf4dd2523237a240b4d400e3bef67c694f1afce5d96a09e10cd7f",
		AdminUsername:   "testadminuser",
		AdminPassword:   "ac867dd8ca873e2285c69c2f9678c13945cad4471992919cfb38345052797b7f83",
		PwdEncrypter:    "aes",
		EnableEncrypted: true,
		Endpoint:        "endpoint",
		Insecure:        false,
	}
	c, err := NewClientForAdmin(authOptions)
	if err != nil {
		t.Error("Test NewClientForAdmin failed")
	}
	if c.user != authOptions.AdminUsername {
		t.Error("Test NewClientForAdmin failed")
	}
	if c.password != "testadminpassword" {
		t.Error("Test NewClientForAdmin failed")
	}
	if c.endpoint != authOptions.Endpoint {
		t.Error("Test NewClientForAdmin failed")
	}
}

func TestNewClientForAdmin_Error(t *testing.T) {
	authOptions := &AuthOptions{
		Username:        "testuser",
		Password:        "d",
		AdminUsername:   "testadminuser",
		AdminPassword:   "testadminpassword",
		PwdEncrypter:    "aes",
		EnableEncrypted: true,
		Endpoint:        "endpoint",
		Insecure:        false,
	}
	_, err := NewClientForAdmin(authOptions)
	if err == nil {
		t.Error("Test NewClientForAdmin failed")
	}
}

func TestDestroy(t *testing.T) {
	// create mock
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte{0x65, 0x78, 0x69, 0x74, 0xa}).Return(1, nil)

	client := &EternusClient{
		user:        "testuser",
		password:    "testpassword",
		endpoint:    "testendpoint",
		stdin:       mockWriteCloser,
		cliConfPath: "./config/cli_response.yml",
	}

	err := client.Destroy()
	if err != nil {
		t.Error(err)
	}
}

func TestSetConfig(t *testing.T) {
	client := &EternusClient{
		user:        "testuser",
		password:    "testpassword",
		endpoint:    "testendpoint",
		cliConfPath: "./config/cli_response.yml",
	}

	config := client.setConfig()
	if config.User != client.user {
		t.Error("Test setConfig failed")
	}
}

func TestDoRequest(t *testing.T) {
	cmd := "show test"
	param := map[string]string{
		"a": "arg1",
	}
	execCmd := "show test -a arg1 \n"
	// create stdin mock
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil)
	mockWriteCloser.On("Write", []byte("\n")).Return(1, nil)
	// create stdout mock
	mockReadCloser := new(MockReadCloser)
	expectStr := "\r\nCLI> show test -a arg1\r\n00\r\nCLI> "
	buff := make([]byte, 65535)
	outStr := "\r\nCLI> "
	out := []byte(outStr)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out)
	buff2 := make([]byte, 65535)
	for i, v := range out {
		buff2[i] = v
	}
	outStr = "show test -a arg1\r\n00\r\nCLI> "
	out = []byte(outStr)
	mockReadCloser.On("Read", buff2).Return(len(out), nil, out)

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	ret, err := client.doRequest(cmd, param)
	if err != nil {
		t.Error("Test doRequest failed")
	}
	if ret.String() != expectStr {
		t.Error("Test doRequest failed")
	}
}

func TestDoRequest_Error(t *testing.T) {
	cmd := "show test"
	param := map[string]string{
		"a": "arg1",
	}
	execCmd := "show test -a arg1 \n"
	// create mock
	mockWriteCloser := new(MockWriteCloser)
	mockWriteCloser.On("Write", []byte(execCmd)).Return(1, nil)
	mockWriteCloser.On("Write", []byte("\n")).Return(1, nil)
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	mockReadCloser.On("Read", buff).Return(0, io.ErrClosedPipe, []byte{})

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	_, err := client.doRequest(cmd, param)
	if err == nil {
		t.Error("Test doRequest failed")
	}
}

func TestRequest(t *testing.T) {
	colName := []string{"lun_group_no", "lun_group_name",
		"previous_affinity_group_no", "next_affinity_group_no",
	}
	data := [][]string{
		[]string{"0000", "test", "FFFF", "FFFF"},
		[]string{"0001", "test2", "FFFF", "FFFF"},
		[]string{"0003", "test3", "FFFF", "FFFF"},
		[]string{"0004", "test4", "FFFF", "FFFF"},
	}
	cmd := "show lun-groups"
	resultArray := [][]string{
		[]string{"00"},
		[]string{"0004"},
	}
	resultArray = append(resultArray, data...)
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	ret, err := client.parseResult(cmd, resultArray)
	if err != nil {
		t.Error("Test doRequest failed")
	}
	for k, e := range ret {
		for i, v := range colName {
			if data[k][i] != e[v] {
				t.Error("Test doRequest failed")
			}
		}
	}
}

func TestRequest_Dropdata(t *testing.T) {
	colName := []string{"tpp_number", "tpp_name", "disk_type",
		"nearline", "raid_level", "tpp_status", "total_capacity",
		"used_capacity", "alarm_status", "warning_level_range",
		"attention_level_range", "encryption_status", "dedup",
		"data_reduction_rate", "provisioned_capacity",
		"provisioned_rate", "dedup_status", "chunk_size"}
	cmd := "show thin-pro-pools"
	data := [][]string{
		[]string{"0002", "osdstest", "01", "00", "04", "0010",
			"00000000666FC000", "0000000000000000", "01", "5A", "4B",
			"00", "00", "00", "0000000000200000", "00000001", "FF", "01"},
	}
	resultArray := [][]string{
		[]string{"00"},
		[]string{"00000000"},
		[]string{"0001"},
	}
	resultArray = append(resultArray, data...)
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	ret, err := client.parseResult(cmd, resultArray)
	if err != nil {
		t.Error("Test doRequest failed")
	}

	for k, e := range ret {
		for i, v := range colName {
			if data[k][i] != e[v] {
				t.Error("Test doRequest failed")
			}
		}
	}
}

func TestRequest_Dropmultidata(t *testing.T) {
	colName := []string{"lun", "volume_no", "volume_name",
		"volume_raw_status", "volume_round_status", "volume_size", "uid",
	}
	cmd := "show lun-group"
	data := [][]string{
		[]string{"0000", "0019", "170-10-vol0", "A000", "20", "0000000000000000", "60000000000000000000000000000000"},
		[]string{"0001", "001A", "170-10-vol1", "A000", "20", "0000000000000000", "60000000000000000000000000000000"},
	}
	resultArray := [][]string{
		[]string{"00"},
		[]string{"0001"},
		[]string{"0001", "170-10", "FFFF", "FFFF"},
		[]string{"0002"},
	}
	resultArray = append(resultArray, data...)
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	ret, err := client.parseResult(cmd, resultArray)
	if err != nil {
		t.Error("Test doRequest failed")
	}

	for k, e := range ret {
		for i, v := range colName {
			if data[k][i] != e[v] {
				t.Error("Test doRequest failed")
			}
		}
	}
}

func TestRequest_ErrorStatus(t *testing.T) {
	cmd := "show lun-group"
	data := [][]string{
		[]string{"0000", "0019", "170-10-vol0", "A000", "20", "0000000000000000", "60000000000000000000000000000000"},
	}
	resultArray := [][]string{
		[]string{"01"},
		[]string{"0001"},
		[]string{"0001", "170-10", "FFFF", "FFFF"},
		[]string{"0001"},
	}
	resultArray = append(resultArray, data...)
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	_, err := client.parseResult(cmd, resultArray)
	if err == nil {
		t.Error("Test doRequest failed")
	}
}

func TestRequest_ErrorNoStatus(t *testing.T) {
	cmd := "show lun-group"
	data := [][]string{
		[]string{"0000", "0019", "170-10-vol0", "A000", "20", "0000000000000000", "60000000000000000000000000000000"},
	}
	resultArray := [][]string{
		[]string{},
		[]string{"0001"},
		[]string{"0001", "170-10", "FFFF", "FFFF"},
		[]string{"0001"},
	}
	resultArray = append(resultArray, data...)
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	_, err := client.parseResult(cmd, resultArray)
	if err == nil {
		t.Error("Test doRequest failed")
	}
}

func TestRequest_ErrorInconsistent(t *testing.T) {
	cmd := "show lun-group"
	data := [][]string{
		[]string{"0000", "0019", "170-10-vol0", "A000", "20", "0000000000000000", "60000000000000000000000000000000"},
	}
	resultArray := [][]string{
		[]string{},
		[]string{"0001"},
		[]string{"0001", "170-10", "FFFF", "FFFF"},
		[]string{"0002"},
	}
	resultArray = append(resultArray, data...)
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	_, err := client.parseResult(cmd, resultArray)
	if err == nil {
		t.Error("Test doRequest failed")
	}
}

func TestGetData(t *testing.T) {
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	var adatas []interface{}
	e := map[interface{}]interface{}{
		"name": "col1",
	}
	adatas = append(adatas, e)
	e = map[interface{}]interface{}{
		"name": "col2",
	}
	adatas = append(adatas, e)
	e = map[interface{}]interface{}{
		"name": "col3",
	}
	adatas = append(adatas, e)

	result := []string{"a", "b", "c"}
	ret, err := client.getData(adatas, result)
	if err != nil {
		t.Error("Test getData failed")
	}
	if ret["col1"] != "a" {
		t.Error("Test getData failed")
	}
	if ret["col2"] != "b" {
		t.Error("Test getData failed")
	}
	if ret["col3"] != "c" {
		t.Error("Test getData failed")
	}
}

func TestGetData_Error(t *testing.T) {
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	var adatas []interface{}
	e := map[interface{}]interface{}{
		"name": "col1",
	}
	adatas = append(adatas, e)
	e = map[interface{}]interface{}{
		"name": "col2",
	}
	adatas = append(adatas, e)

	result := []string{"a"}
	_, err := client.getData(adatas, result)
	if err == nil {
		t.Error("Test getData failed")
	}
}

func TestConvStringArray(t *testing.T) {
	client := &EternusClient{
		cliConfPath: "./config/cli_response.yml",
	}
	cmd := "show test_cmd"
	resultString := "\r\nCLI> show test_cmd\r\n00\r\na\tb\tc\r\nCLI> "
	var resultByte bytes.Buffer
	resultByte.WriteString(resultString)

	ret := client.convStringArray(cmd, &resultByte)

	if len(ret[0]) != 1 || ret[0][0] != "00" {
		t.Error("Test convStringArray failed")
	}

	if len(ret[1]) != 3 || ret[1][0] != "a" ||
		ret[1][1] != "b" || ret[1][2] != "c" {
		t.Error("Test convStringArray failed")
	}
}

func createIOMock(cmd string, output string) *EternusClient {
	mockWriteCloser := new(MockWriteCloser)
	if cmd != "" {
		mockWriteCloser.On("Write", []byte(cmd)).Return(1, nil)
	} else {
		mockWriteCloser.On("Write", mock.Anything).Return(1, nil)
	}
	mockWriteCloser.On("Write", []byte("\n")).Return(1, nil)
	// create stdout mock
	mockReadCloser := new(MockReadCloser)
	buff := make([]byte, 65535)
	out := []byte(output)
	mockReadCloser.On("Read", buff).Return(len(out), nil, out)

	client := &EternusClient{
		stdin:       mockWriteCloser,
		stdout:      mockReadCloser,
		cliConfPath: "./config/cli_response.yml",
	}
	return client
}

func TestGetVolume(t *testing.T) {
	execCmd := "show volumes -volume-number 1 \n"
	outStr := "\r\nCLI> show volumes -volume-number 1 \r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0012\tosds-643e8232-1b\tA000\t09\t00\t0002\tosdstest\t0000000000000000\t00\t00\t00000000\t0050\tFF\t00\tFF\tFF\t20\t00\tFFFF\t00\t60000000000000000000000000000000\t00\t00\tFF\tFF\tFFFFFFFF\t00\t00\tFF\t00\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)
	ret, err := client.GetVolume("1")
	if err != nil {
		t.Error("Test doRequest failed")
	}
	if ret.Id != "18" && ret.Name != "osds-643e8232-1b" &&
		ret.Size != 1 && ret.Status != "00" &&
		ret.PoolName != "tosdstest" && ret.PoolId != "2" {
		t.Error("Test GetVolume failed")
	}
}

func TestGetVolumeByName(t *testing.T) {
	execCmd := "show volumes -volume-name osds-643e8232-1b \n"
	outStr := "\r\nCLI> show volumes -volume-name osds-643e8232-1b \r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "0012\tosds-643e8232-1b\tA000\t09\t00\t0002\tosdstest\t0000000000000000\t00\t00\t00000000\t0050\tFF\t00\tFF\tFF\t20\t00\tFFFF\t00\t60000000000000000000000000000000\t00\t00\tFF\tFF\tFFFFFFFF\t00\t00\tFF\t00\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)
	ret, err := client.GetVolumeByName("osds-643e8232-1b")
	if err != nil {
		t.Error("Test doRequest failed")
	}
	if ret.Id != "18" && ret.Name != "osds-643e8232-1b" &&
		ret.Size != 1 && ret.Status != "00" &&
		ret.PoolName != "tosdstest" && ret.PoolId != "2" {
		t.Error("Test GetVolume failed")
	}
}

func TestDeleteIscsiHost(t *testing.T) {
	name := "hostname"

	execCmd := "delete host-iscsi-name -host-name " + name + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)
	err := client.DeleteIscsiHostByName(name)
	if err != nil {
		t.Error("Test DeleteIscsiHostByName failed")
	}
}

func TestDeleteLunGroup(t *testing.T) {
	name := "hostname"
	execCmd := "delete lun-group -lg-name " + name + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	client := createIOMock("", outStr)
	err := client.DeleteLunGroupByName(name)
	if err != nil {
		t.Error("Test DeleteLunGroupByName failed")
	}
}

func TestDeleteHostAffinity(t *testing.T) {
	port := "010"
	name := "hostname"
	execCmd := "release host-affinity -port " + port
	execCmd += " -host-name " + name + " -mode all" + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	client := createIOMock("", outStr)
	err := client.DeleteHostAffinity(port, name)
	if err != nil {
		t.Error("Test DeleteHostAffinity failed")
	}
}

func TestDeleteFcHost(t *testing.T) {
	name := "hostname"
	execCmd := "delete host-wwn-name -host-name " + name + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)
	err := client.DeleteFcHostByName(name)
	if err != nil {
		t.Error("Test DeleteFcHost failed")
	}
}

func TestListMapping(t *testing.T) {
	port := "010"

	execCmd := "show mapping -port " + port + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "0001\r\n"
	outStr += "11005100\t00\r\n"
	outStr += "0001\r\n"
	outStr += "0000\t0001\tosds-643e8232-1b\tA000\t20\t0000000000000000\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)
	mapping, err := client.ListMapping(port)
	if err != nil {
		t.Error("Test ListMapping failed")
	}
	for _, v := range mapping {
		if v.Lun != "0" && v.VolumeNumber != "1" &&
			v.VolumeName != "osds-643e8232-1b" && v.VolumeRawStatus != "A000" &&
			v.VolumeRoundStatus != "20" && v.VolumeSize != 1 {
			t.Error("Test ListMapping failed")
		}
	}
}

func TestCreateSnapshot(t *testing.T) {
	srcLunID := "0"
	destLunID := "1"
	execCmd := "start advanced-copy -source-volume-number " + srcLunID + " -destination-volume-number " + destLunID + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	client := createIOMock("", outStr)
	err := client.CreateSnapshot(srcLunID, destLunID)
	if err != nil {
		t.Error("Test DeleteMapping failed")
	}
}

func TestListSnapshot(t *testing.T) {
	execCmd := "show advanced-copy-session -type sopc+  \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "SID   Gene-   Type     Volume    Source Volume                          Destination Volume                     Status        Phase            Error  Requestor\r\n"
	outStr += "      ration           Type      No.   Name                             No.   Name                                                            Code\r\n"
	outStr += "----- ------- -------- --------- ----- -------------------------------- ----- -------------------------------- ------------- ---------------- ------ ---------\r\n"
	outStr += "   50   1/  1 SnapOPC+ Standard      2 testvol0                            10 testvol1                         Active        Copying           0x00  GUI\r\n"
	outStr += "CLI> "
	client := createIOMock("", outStr)
	snapshots, err := client.ListSnapshot()
	if err != nil {
		t.Error("Test ListSnapshot failed")
	}
	for _, v := range snapshots {
		if v.Sid != "50" && v.Gen != "1/" &&
			v.GenTotal != "1" && v.Type != "SnapOPC+" &&
			v.VolumeType != "Standard" && v.SrcNo != "2" &&
			v.SrcName != "testvol0" && v.DestNo != "10" &&
			v.DestName != "testvol1" && v.Status != "Active" &&
			v.Phase != "Copying" && v.ErrorCode != "0x00" &&
			v.Requestor != "GUI" {
			t.Error("Test ListSnapshot failed")
		}
	}
}

func TestDeleteSnapshot(t *testing.T) {
	sid := "0!"
	execCmd := "stop advanced-copy -session-id " + sid + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "00\r\n"
	outStr += "CLI> "
	client := createIOMock("", outStr)
	err := client.DeleteSnapshot(sid)
	if err != nil {
		t.Error("Test DeleteSnapshot failed")
	}
}

func TestDeleteSnapshot_AlreadyDelete(t *testing.T) {
	sid := "0!"
	execCmd := "stop advanced-copy -session-id " + sid + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "Error: E0110 Resource does not exist.\r\n"
	outStr += "CLI> "
	client := createIOMock("", outStr)
	err := client.DeleteSnapshot(sid)
	if err != nil {
		t.Error("Test DeleteSnapshot failed")
	}
}

func TestDeleteSnapshot_Error(t *testing.T) {
	sid := "0"
	execCmd := "stop advanced-copy -session-id " + sid + " \n"
	outStr := "\r\nCLI> " + execCmd + " \r\n"
	outStr += "Error: E0112 XXXXXXXXXXXXX.\r\n"
	outStr += "CLI> "
	client := createIOMock(execCmd, outStr)
	err := client.DeleteSnapshot(sid)
	if err == nil {
		t.Error("Test DeleteSnapshot failed")
	}
}
