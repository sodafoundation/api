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

package fusionstorage

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/exec"
	"github.com/opensds/opensds/pkg/utils/pwd"
)

var CliErrorMap = map[string]string{
	"50000001": "DSware error",
	"50150001": "Receive a duplicate request",
	"50150002": "Command type is not supported",
	"50150003": "Command format is error",
	"50150004": "Lost contact with major VBS",
	"50150005": "Volume does not exist",
	"50150006": "Snapshot does not exist",
	"50150007": "Volume already exists or name exists or name duplicates with a snapshot name",
	"50150008": "The snapshot has already existed",
	"50150009": "VBS space is not enough",
	"50150010": "The node type is error",
	"50150011": "Volume and snapshot number is beyond max",
	"50150012": "VBS is not ready",
	"50150013": "The ref num of node is not 0",
	"50150014": "The volume is not in the pre-deletion state.",
	"50150015": "The storage resource pool is faulty",
	"50150016": "VBS handle queue busy",
	"50150017": "VBS handle request timeout",
	"50150020": "VBS metablock is locked",
	"50150021": "VBS pool dose not exist",
	"50150022": "VBS is not ok",
	"50150023": "VBS pool is not ok",
	"50150024": "VBS dose not exist",
	"50150064": "VBS load SCSI-3 lock pr meta failed",
	"50150100": "The disaster recovery relationship exists",
	"50150101": "The DR relationship does not exist",
	"50150102": "Volume has existed mirror",
	"50150103": "The volume does not have a mirror",
	"50150104": "Incorrect volume status",
	"50150105": "The mirror volume already exists",
}

func NewCliError(code string) error {
	if msg, ok := CliErrorMap[code]; ok {
		return NewCliErrorBase(msg, code)
	}
	return NewCliErrorBase("CLI execute error", code)
}

type CliError struct {
	Msg  string
	Code string
}

func (c *CliError) Error() string {
	return fmt.Sprintf("msg: %s, code:%s", c.Msg, c.Code)
}

func NewCliErrorBase(msg, code string) *CliError {
	return &CliError{Msg: msg, Code: code}
}

type FsClient struct {
	username string
	password string
	version  string
	addess   string
	headers  map[string]string
	fmIp     string
	fsaIp    []string
}

// Command Root exectuer
var rootExecuter = exec.NewRootExecuter()

func newRestCommon(conf *Config) (*FsClient, error) {
	if conf.Version != ClientVersion6_3 && conf.Version != ClientVersion8_0 {
		return nil, fmt.Errorf("version %s does not support", conf.Version)
	}

	if conf.Version == ClientVersion6_3 {
		if len(conf.FmIp) == 0 || len(conf.FsaIp) == 0 {
			return nil, fmt.Errorf("get %s cli failed, FM ip or FSA ip can not be set to empty", ClientVersion6_3)
		}
		err := StartServer()
		if err != nil {
			return nil, fmt.Errorf("get new client failed, %v", err)
		}
	}

	var pwdCiphertext = conf.Password

	if conf.EnableEncrypted {
		// Decrypte the password
		pwdTool := pwd.NewPwdEncrypter(conf.PwdEncrypter)
		password, err := pwdTool.Decrypter(pwdCiphertext)
		if err != nil {
			return nil, err
		}
		pwdCiphertext = password
	}

	client := &FsClient{
		addess:   conf.Url,
		username: conf.Username,
		password: pwdCiphertext,
		fmIp:     conf.FmIp,
		fsaIp:    conf.FsaIp,
		headers:  map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	var err error
	for i := 1; i <= MaxRetry; i++ {
		log.Printf("try to login the client %d time", i)
		err = client.login()
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *FsClient) getVersion() error {
	url := "rest/version"
	c.headers["Referer"] = c.addess + BasicURI
	content, err := c.request(url, "GET", true, nil)
	if err != nil {
		return fmt.Errorf("failed to get version, %v", err)
	}

	var v Version
	err = json.Unmarshal(content, &v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the result, %v", err)
	}

	c.version = v.CurrentVersion

	return nil
}

func (c *FsClient) login() error {
	c.getVersion()
	url := "/sec/login"
	data := map[string]string{"userName": c.username, "password": c.password}
	_, err := c.request(url, "POST", false, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *FsClient) logout() error {
	url := "/iam/logout"
	_, err := c.request(url, "POST", false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) queryPoolInfo() (*PoolResp, error) {
	url := "/storagePool"
	result, err := c.request(url, "GET", false, nil)
	if err != nil {
		return nil, err
	}

	var pools *PoolResp
	if err := json.Unmarshal(result, &pools); err != nil {
		return nil, err
	}
	return pools, nil
}

func (c *FsClient) createVolume(volName, poolId string, volSize int64) error {
	url := "/volume/create"
	polID, _ := strconv.Atoi(poolId)
	params := map[string]interface{}{"volName": volName, "volSize": volSize, "poolId": polID}

	if _, err := c.request(url, "POST", false, params); err != nil {
		return err
	}
	return nil
}

func (c *FsClient) deleteVolume(volName string) error {
	url := "/volume/delete"
	params := map[string]interface{}{"volNames": []string{volName}}
	_, err := c.request(url, "POST", false, params)
	if err != nil {
		return err
	}

	return nil
}

func (c *FsClient) attachVolume(volName, manageIp string) error {
	url := "/volume/attach"
	params := map[string]interface{}{"volName": []string{volName}, "ipList": []string{manageIp}}
	_, err := c.request(url, "POST", false, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) createPort(initiator string) error {
	url := "iscsi/createPort"
	params := map[string]interface{}{"portName": initiator}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) queryPortInfo(initiator string) error {
	url := "iscsi/queryPortInfo"
	params := map[string]interface{}{"portName": initiator}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}

	return nil
}

func (c *FsClient) queryHostInfo(hostName string) (bool, error) {
	url := "iscsi/queryAllHost"
	result, err := c.request(url, "GET", true, nil)
	if err != nil {
		return false, err
	}

	var hostlist *HostList

	if err := json.Unmarshal(result, &hostlist); err != nil {
		return false, err
	}

	for _, v := range hostlist.HostList {
		if v.HostName == hostName {
			return true, nil
		}
	}

	return false, nil
}

func (c *FsClient) createHost(hostInfo *pb.HostInfo) error {
	url := "iscsi/createHost"
	params := map[string]interface{}{"hostName": hostInfo.GetHost(), "ipAddress": hostInfo.GetIp()}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) addPortToHost(hostName, initiator string) error {
	url := "iscsi/addPortToHost"
	params := map[string]interface{}{"hostName": hostName, "portNames": []string{initiator}}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) queryHostByPort(initiator string) (*PortHostMap, error) {
	url := "iscsi/queryHostByPort"
	params := map[string]interface{}{"portName": []string{initiator}}
	result, err := c.request(url, "POST", true, params)
	if err != nil {
		return nil, err
	}

	var portHostmap *PortHostMap

	if err := json.Unmarshal(result, &portHostmap); err != nil {
		return nil, err
	}

	return portHostmap, nil
}

func (c *FsClient) addLunsToHost(hostName, lunId string) error {
	url := "iscsi/addLunsToHost"
	params := map[string]interface{}{"hostName": hostName, "lunNames": []string{lunId}}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) queryHostLunInfo(hostName string) (*HostLunList, error) {
	url := "iscsi/queryHostLunInfo"
	params := map[string]interface{}{"hostName": hostName}
	result, err := c.request(url, "POST", true, params)
	if err != nil {
		return nil, err
	}

	var lunList *HostLunList

	if err := json.Unmarshal(result, &lunList); err != nil {
		return nil, err
	}

	return lunList, nil
}

func (c *FsClient) queryIscsiPortalVersion6(initiator string) ([]string, error) {
	args := []string{
		"--op", "queryIscsiPortalInfo", "--portName", initiator,
	}
	out, err := c.RunCmd(args...)
	if err != nil {
		return nil, fmt.Errorf("query iscsi portal failed: %v", err)
	}

	if len(out) > 0 {
		return out, nil
	}

	return nil, fmt.Errorf("the iscsi target portal is empty.")
}

func (c *FsClient) getDeviceVersion() (*DeviceVersion, error) {
	url := "/version"
	result, err := c.request(url, "Get", false, nil)
	if err != nil {
		return nil, err
	}

	var version *DeviceVersion

	if err := json.Unmarshal(result, &version); err != nil {
		return nil, err
	}

	return version, nil
}

func (c *FsClient) queryIscsiPortalVersion8() (*IscsiPortal, error) {
	url := "cluster/dswareclient/queryIscsiPortal"
	params := map[string]interface{}{}
	result, err := c.request(url, "Post", true, params)
	if err != nil {
		return nil, err
	}

	var iscsiPortals *IscsiPortal

	if err := json.Unmarshal(result, &iscsiPortals); err != nil {
		return nil, err
	}

	return iscsiPortals, nil
}

func (c *FsClient) queryHostFromVolume(lunId string) ([]Host, error) {
	url := "iscsi/queryHostFromVolume"
	params := map[string]interface{}{"lunName": lunId}
	out, err := c.request(url, "POST", true, params)
	if err != nil {
		return nil, err
	}

	var hostlist *HostList

	if err := json.Unmarshal(out, &hostlist); err != nil {
		return nil, err
	}

	return hostlist.HostList, nil
}

func (c *FsClient) deleteLunFromHost(hostName, lunId string) error {
	url := "iscsi/deleteLunFromHost"
	params := map[string]interface{}{"hostName": hostName, "lunNames": []string{lunId}}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) deletePortFromHost(hostName, initiator string) error {
	url := "iscsi/deletePortFromHost"
	params := map[string]interface{}{"hostName": hostName, "portNames": []string{initiator}}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) deleteHost(hostName string) error {
	url := "iscsi/deleteHost"
	params := map[string]interface{}{"hostName": hostName}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) deletePort(initiator string) error {
	url := "iscsi/deletePort"
	params := map[string]interface{}{"portName": initiator}
	_, err := c.request(url, "POST", true, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) request(url, method string, isGetVersion bool, reqParams interface{}) ([]byte, error) {
	var callUrl string
	if !isGetVersion {
		callUrl = c.addess + BasicURI + c.version + url
	} else {
		callUrl = c.addess + BasicURI + url
	}

	// No verify by SSL
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// initialize http client
	client := &http.Client{Transport: tr}

	var body []byte
	var err error
	if reqParams != nil {
		body, err = json.Marshal(reqParams)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal the request parameters, url is %s, error is %v", callUrl, err)
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), callUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to initiate the request, url is %s, error is %v", callUrl, err)
	}

	// initiate the header
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("process request failed: %v, url is %s", err, callUrl)
	}
	defer resp.Body.Close()

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read from response body failed: %v, url is %s", err, callUrl)
	}

	if 400 <= resp.StatusCode && resp.StatusCode <= 599 {
		pc, _, line, _ := runtime.Caller(1)
		return nil, fmt.Errorf("return status code is: %s, return content is: %s, error function is: %s, error line is: %s, url is %s",
			strconv.Itoa(resp.StatusCode), string(respContent), runtime.FuncForPC(pc).Name(), strconv.Itoa(line), callUrl)
	}

	// Check the error code in the returned content
	var respResult *ResponseResult
	if err := json.Unmarshal(respContent, &respResult); err != nil {
		return nil, err
	}

	if respResult.RespCode != 0 {
		return nil, errors.New(string(respContent))
	}

	if c.headers["x-auth-token"] == "" && resp.Header != nil && len(resp.Header["X-Auth-Token"]) > 0 {
		c.headers["x-auth-token"] = resp.Header["X-Auth-Token"][0]
	}

	return respContent, nil
}

func StartServer() error {
	_, err := rootExecuter.Run(CmdBin, "--op", "startServer")
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	return nil
}

func (c *FsClient) RunCmd(args ...string) ([]string, error) {
	var lines []string
	var result string

	args = append(args, "--manage_ip", c.fmIp, "--ip", "")
	for _, ip := range c.fsaIp {
		args[len(args)-1] = ip
		out, _ := rootExecuter.Run(CmdBin, args...)
		lines = strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) > 0 {
			const resultPrefix = "result="
			for _, line := range lines {
				if strings.HasPrefix(line, resultPrefix) {
					result = line[len(resultPrefix):]
				}
			}
			if result == "0" {
				return lines[:len(lines)-1], nil
			}
		}
	}

	return nil, NewCliError(result)
}

func (c *FsClient) extendVolume(name string, newSize int64) error {
	url := "/volume/expand"
	params := map[string]interface{}{"volName": name, "newVolSize": newSize}
	_, err := c.request(url, "POST", false, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) createSnapshot(snapName, volName string) error {
	url := "/snapshot/create"
	params := map[string]interface{}{"volName": volName, "snapshotName": snapName}
	_, err := c.request(url, "POST", false, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *FsClient) deleteSnapshot(snapName string) error {
	url := "/snapshot/delete"
	params := map[string]interface{}{"snapshotName": snapName}
	_, err := c.request(url, "POST", false, params)
	if err != nil {
		return err
	}
	return nil
}
