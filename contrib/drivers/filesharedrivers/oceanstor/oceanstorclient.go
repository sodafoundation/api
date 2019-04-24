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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	//pb "github.com/opensds/opensds/pkg/model/proto"
	//"github.com/opensds/opensds/pkg/utils/exec"
	"github.com/opensds/opensds/pkg/utils/pwd"
)

type Cli struct {
	username   string
	password   string
	urlPrefix  string
	deviceId   string
	iBaseToken string
	cookie     string
	header     map[string]string
}

func newRestCommon(conf *Config) (*Cli, error) {

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

	client := &Cli{
		urlPrefix: conf.Uri,
		username:  conf.Username,
		password:  pwdCiphertext,
		header:    map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	err := client.login()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Cli) login() error {
	data := map[string]string{
		"username": c.username,
		"password": c.password,
		"scope":    "0",
	}

	url := "/xxxxx/sessions"
	resp, err := c.request(url, "POST", data)
	if err != nil {
		log.Errorf("login failed: %v", err)
		return err
	}

	var auth *AuthResp

	if err := json.Unmarshal(resp, &auth); err != nil {
		return err
	}

	c.iBaseToken = auth.Data.IBaseToken
	if auth.Data.AccountState == PwdReset || auth.Data.AccountState == PwdExpired {
		msg := "Password has expired or must be reset,please change the password."
		log.Error(msg)
		c.logout()
		return errors.New(msg)
	}
	if auth.Data.DeviceId == "" {
		msg := fmt.Sprintf("failed to login with rest URLs %s", c.urlPrefix)
		log.Error(msg)
		return errors.New(msg)
	}

	c.deviceId = auth.Data.DeviceId
	c.urlPrefix += "/" + auth.Data.DeviceId

	c.header["Connection"] = "keep-alive"
	c.header["iBaseToken"] = c.iBaseToken
	c.header["Cookie"] = c.cookie

	return nil
}

func (c *Cli) logout() ([]byte, error) {
	return c.request("DELETE", "/sessions", nil)
}

func (c *Cli) createNFSShare(shareName, fsID string) (*NFSShareData, error) {
	sharePath := getSharePath(shareName)
	data := map[string]string{
		"DESCRIPTION": "",
		"FSID":        fsID,
		"SHAREPATH":   sharePath,
	}
	fmt.Println(data)
	url := "/NFSHARE"

	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var nfsShare *NFSShare
	if err := json.Unmarshal(resp, &nfsShare); err != nil {
		return nil, err
	}

	if nfsShare.Error.Code != 0 {
		return nil, errors.New(nfsShare.Error.Description)
	}

	return &nfsShare.NFSShareData, nil
}

func (c *Cli) createCIFSShare(shareName, fsId string) (*CIFSShareData, error) {
	sharePath := getSharePath(shareName)
	data := map[string]string{
		"SHAREPATH":    sharePath,
		"DESCRIPTION":  "",
		"ABEENABLE":    "false",
		"ENABLENOTIFY": "true",
		"ENABLEOPLOCK": "true",
		"NAME":         strings.Replace(shareName, "-", "_", -1),
		"FSID":         fsId,
		"TENANCYID":    "0",
	}

	url := "/CIFSHARE"

	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var cifsShare *CIFSShare
	if err := json.Unmarshal(resp, &cifsShare); err != nil {
		return nil, err
	}

	if cifsShare.Error.Code != 0 {
		return nil, errors.New(cifsShare.Error.Description)
	}

	return &cifsShare.CIFSShareData, nil
}

func getSharePath(shareName string) string {
	sharePath := "/" + strings.Replace(shareName, "-", "_", -1) + "/"
	return sharePath
}

func getShareUrlType(shareProto string) (string, error) {
	if shareProto == NFS {
		return "NFSHARE", nil
	}

	if shareProto == CIFS {
		return "CIFSHARE", nil
	}

	return "", errors.New(shareProto + " protocol is not supported")
}

func (c *Cli) getNFSShare(shareName string) (*NFSShareData, error) {
	url := fmt.Sprintf("/NFSHARE?filter=SHAREPATH::%s&range=[0-100]", getSharePath(shareName))
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var nfsShareList *NFSShareList
	if err := json.Unmarshal(resp, &nfsShareList); err != nil {
		return nil, err
	}

	if nfsShareList.Error.Code != 0 {
		return nil, errors.New(nfsShareList.Error.Description)
	}

	if len(nfsShareList.Data) > 0 {
		return &nfsShareList.Data[0], nil
	}

	//	for _, v := range nfsShare.Data {
	//		if v.Name == "/"+shareName {
	//			return &v, nil
	//		}
	//	}
	//	if nfsShareList.ID == "" {
	//		return nil, nil
	//	}

	return nil, nil
}

func (c *Cli) getCIFSShare(shareName string) (*CIFSShareData, error) {
	url := fmt.Sprintf("/CIFSHARE?filter=NAME:%s&range=[0-100]", strings.Replace(shareName, "-", "_", -1))
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var cifsShareList *CIFSShareList
	if err := json.Unmarshal(resp, &cifsShareList); err != nil {
		return nil, err
	}

	if cifsShareList.Error.Code != 0 {
		return nil, errors.New(cifsShareList.Error.Description)
	}

	//	for _, v := range cifsShareList.Data {
	//		if v.Name == shareName {
	//			return &v, nil
	//		}
	//	}

	if len(cifsShareList.Data) > 0 {
		return &cifsShareList.Data[0], nil
	}

	return nil, nil
}

func (c *Cli) listCIFSShares() (*CIFSShareData, error) {
	url := "/CIFSHARE"
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var cifsShareList *CIFSShareList
	if err := json.Unmarshal(resp, &cifsShareList); err != nil {
		return nil, err
	}

	if cifsShareList.Error.Code != 0 {
		return nil, errors.New(cifsShareList.Error.Description)
	}

	return cifsShareList, nil
}

func (c *Cli) createFileSystem(name, poolID, tenantID string) (*FileSystem, error) {
	data := map[string]interface{}{
		"PARENTID":   poolID,
		"NAME":       name,
		"PARENTTYPE": 216,
		"ALLOCTYPE":  1,
	}

	url := "/filesystem"
	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var fileSystem *FileSystem
	if err := json.Unmarshal(resp, &fileSystem); err != nil {
		return nil, err
	}

	if fileSystem.Error.Code != 0 {
		return nil, errors.New(fileSystem.Error.Description)
	}

	return fileSystem, nil
}

func (c *Cli) deleteNFSShare(shareID string) error {
	url := "/nfshare/" + shareID
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete *DeleteError
	if err := json.Unmarshal(resp, &errDelete); err != nil {
		return err
	}

	if errDelete.Code != 0 {
		msg := errDelete.Description
		return errors.New(msg)
	}

	return nil
}

func (c *Cli) deleteFS(fsID string) error {
	url := "/filesystem/" + fsID
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete *DeleteError
	if err := json.Unmarshal(resp, &errDelete); err != nil {
		return err
	}

	if errDelete.Code != 0 {
		msg := errDelete.Description
		return errors.New(msg)
	}

	return nil
}

func (c *Cli) getNFSShareByID(shareID, shareProto string) (*NFSShare, error) {
	url := "/NFSHARE/" + shareID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var nfsShare *NFSShare
	if err := json.Unmarshal(resp, &nfsShare); err != nil {
		return nil, err
	}

	if nfsShare.Error.Code != 0 {
		return nil, errors.New(nfsShare.Error.Description)
	}

	if nfsShare.ID == "" {
		return nil, nil
	}

	return nfsShare, nil
}

func (c *Cli) getCIFSShareByID(shareID, shareProto string) (*CIFSShare, error) {
	url := "/CIFSHARE/" + shareID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var cifsShare *CIFSShare
	if err := json.Unmarshal(resp, &cifsShare); err != nil {
		return nil, err
	}

	if cifsShare.Error.Code != 0 {
		return nil, errors.New(cifsShare.Error.Description)
	}

	if cifsShare.ID == "" {
		return nil, nil
	}

	return cifsShare, nil
}

func (c *Cli) deleteCIFSShare(shareID string) error {
	url := "/cifshare/" + shareID
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete *DeleteError
	if err := json.Unmarshal(resp, &errDelete); err != nil {
		return err
	}

	if errDelete.Code != 0 {
		msg := errDelete.Description
		return errors.New(msg)
	}

	return nil
}

func (c *Cli) getFileSystem(fsid string) (*FileSystem, error) {
	url := "/filesystem/" + fsid
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fileSystem *FileSystem
	if err := json.Unmarshal(resp, &fileSystem); err != nil {
		return nil, err
	}

	if fileSystem.Error.Code != 0 {
		msg := "get file system failed, " + fileSystem.Error.Description
		return nil, errors.New(msg)
	}

	return fileSystem, nil
}

func (c *Cli) getFileSystemByName(name string) (*FileSystemList, error) {
	url := "/filesystem?filter=NAME::" + name
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fsList *FileSystemList
	if err := json.Unmarshal(resp, &fsList); err != nil {
		return nil, err
	}

	if fsList.Error.Code != 0 {
		return nil, errors.New(fsList.Error.Description)
	}

	return fsList, nil
}

func (c *Cli) deleteFileSystem(name, poolID, tenantID string) (*FileSystem, error) {
	data := map[string]interface{}{
		"PARENTID":   poolID,
		"NAME":       name,
		"PARENTTYPE": 216,
		"ALLOCTYPE":  1,
	}

	url := "/filesystem"
	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var fileSystem *FileSystem
	if err := json.Unmarshal(resp, &fileSystem); err != nil {
		return nil, err
	}

	if fileSystem.Error.Code != 0 {
		msg := "create file system failed, " + fileSystem.Error.Description
		return nil, errors.New(msg)
	}

	return fileSystem, nil
}

func (c *Cli) ListStoragePools() ([]StoragePool, error) {
	var pools *StoragePoolsResp
	resp, err := c.request("/storagepool", "GET", nil)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resp, &pools); err != nil {
		return nil, err
	}

	return pools.Data, nil
}

func (c *Cli) createSnapshot(fsID, snapName string) (*FSSnapshot, error) {
	data := map[string]string{
		"PARENTTYPE":  "40",
		"TYPE":        "48",
		"PARENTID":    fsID,
		"NAME":        strings.Replace(snapName, "-", "_", -1),
		"DESCRIPTION": "",
	}
	url := "/FSSNAPSHOT"
	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var fsSnapshot *FSSnapshot
	if err := json.Unmarshal(resp, &fsSnapshot); err != nil {
		return nil, err
	}

	if fsSnapshot.Error.Code != 0 {
		msg := fsSnapshot.Error.Description
		return nil, errors.New(msg)
	}

	return fsSnapshot, nil
}

func (c *Cli) listSnapshots(fsID string) (*FSSnapshotList, error) {
	url := "/FSSNAPSHOT?sortby=TIMESTAMP,d&range=[0-100]&PARENTID=" + fsID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fsSnapshotList *FSSnapshotList
	if err := json.Unmarshal(resp, &fsSnapshotList); err != nil {
		return nil, err
	}

	if fsSnapshotList.Error.Code != 0 {
		msg := fsSnapshotList.Error.Description
		return nil, errors.New(msg)
	}

	return fsSnapshotList, nil
}

func (c *Cli) deleteFSSnapshot(snapID string) error {
	url := "/FSSNAPSHOT/" + snapID
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete *DeleteError
	if err := json.Unmarshal(resp, &errDelete); err != nil {
		return err
	}

	if errDelete.Code != 0 {
		msg := errDelete.Description
		return errors.New(msg)
	}

	return nil
}

func (c *Cli) showFSSnapshot(snapID string) (*FSSnapshot, error) {
	url := "/FSSNAPSHOT/" + snapID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fsSnapshot *FSSnapshot
	if err := json.Unmarshal(resp, &fsSnapshot); err != nil {
		return nil, err
	}

	if fsSnapshot.Error.Code != 0 {
		msg := fsSnapshot.Error.Description
		return nil, errors.New(msg)
	}

	return fsSnapshot, nil
}

func (c *Cli) getAllFilesystem() (*FileSystemList, error) {
	url := "/filesystem"
	resp, err := c.request(url, "Get", nil)
	if err != nil {
		return nil, err
	}

	var fsList *FileSystemList
	if err := json.Unmarshal(resp, &fsList); err != nil {
		return nil, err
	}

	if fsList.Error.Code != 0 {
		msg := fsList.Error.Description
		return nil, errors.New(msg)
	}

	return fsList, nil
}

func (c *Cli) request(url, method string, reqParams interface{}) ([]byte, error) {

	callUrl := c.urlPrefix + url
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
			return nil, fmt.Errorf("Failed to marshal the request parameters, url is %s, error is %v", callUrl, err)
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), callUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Failed to initiate the request, url is %s, error is %v", callUrl, err)
	}

	// initiate the header
	for k, v := range c.header {
		req.Header.Set(k, v)
	}

	// do the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Process request failed: %v, url is %s", err, callUrl)
	}
	defer resp.Body.Close()

	respContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read from response body failed: %v, url is %s", err, callUrl)
	}

	if 400 <= resp.StatusCode && resp.StatusCode <= 599 {
		pc, _, line, _ := runtime.Caller(1)
		return nil, fmt.Errorf("return status code is: %s, return content is: %s, error function is: %s, error line is: %s, url is %s",
			strconv.Itoa(resp.StatusCode), string(respContent), runtime.FuncForPC(pc).Name(), strconv.Itoa(line), callUrl)
	}
	//	// Check the error code in the returned content
	//	var respResult *responseResult
	//	if err := json.Unmarshal(respContent, &respResult); err != nil {
	//		return nil, err
	//	}

	//	if respResult.RespCode != 0 {
	//		return nil, fmt.Errorf(string(respContent))
	//	}
	if c.cookie == "" && resp.Header != nil {
		cookie := resp.Header.Get("set-cookie")
		if cookie != "" {
			c.cookie = cookie
		}
	}

	return respContent, nil
}
