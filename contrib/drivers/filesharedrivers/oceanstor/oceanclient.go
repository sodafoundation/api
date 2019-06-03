// Copyright 2019 The OpenSDS Authors.
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
	"github.com/opensds/opensds/pkg/utils/pwd"
)

func newRestCommon(conf *Config) (*Client, error) {
	pwdCiphertext, err := decrypte(conf)
	if err != nil {
		msg := fmt.Sprintf("decryption failed: %v", err)
		log.Error(msg)
		return nil, err
	}

	client := &Client{
		urlPrefix: conf.Uri,
		username:  conf.Username,
		password:  pwdCiphertext,
		header:    map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	if err := tryTimes(client.login); err != nil {
		msg := fmt.Sprintf("login failed: %v", err)
		log.Error(msg)
		return nil, err
	}

	return client, nil
}

func decrypte(conf *Config) (string, error) {
	var pwdCiphertext = conf.Password

	if conf.EnableEncrypted {
		// Decrypte the password
		pwdTool := pwd.NewPwdEncrypter(conf.PwdEncrypter)
		password, err := pwdTool.Decrypter(pwdCiphertext)
		if err != nil {
			return "", err
		}
		pwdCiphertext = password
	}

	return pwdCiphertext, nil
}

func (c *Client) login() error {
	auth, err := c.getAuth()
	if err != nil {
		return err
	}

	if auth.AccountState == PwdReset || auth.AccountState == PwdExpired {
		c.logout()
		return errors.New("password has expired or must be reset, please change the password")
	}

	if auth.DeviceId == "" {
		c.logout()
		return fmt.Errorf("failed to login with rest URLs %s", c.urlPrefix)
	}

	c.urlPrefix += "/" + auth.DeviceId

	c.header["Connection"] = "keep-alive"
	c.header["iBaseToken"] = auth.IBaseToken
	c.header["Cookie"] = c.cookie

	return nil
}

func (c *Client) getAuth() (*AuthData, error) {
	data := map[string]string{
		"username": c.username,
		"password": c.password,
		"scope":    "0",
	}

	url := "/xxxxx/sessions"
	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var auth Auth

	if err := handleReponse(resp, &auth); err != nil {
		return nil, err
	}

	return &auth.Data, nil
}

type Protocol interface {
	createShare(fsName, fsID string) (interface{}, error)
	getShare(fsName string) (interface{}, error)
	getShareID(share interface{}) string
	deleteShare(shareID string) error
	getShareByID(shareID string) (interface{}, error)
	getLocation(sharePath, ipAddr string) string
	allowAccess(shareID, accessTo, accessLevel string) (interface{}, error)
	getAccessLevel(accessLevel string) string
}

func NewProtocol(proto string, c *Client) Protocol {
	switch proto {
	case NFSProto:
		return &NFS{Client: c}
	case CIFSProto:
		return &CIFS{Client: c}
	}

	return nil
}

func (c *Client) createFileSystem(name, poolID string, size int64) (*FileSystemData, error) {
	data := map[string]interface{}{
		"PARENTID":   poolID,
		"NAME":       name,
		"PARENTTYPE": 216,
		"ALLOCTYPE":  1,
		"CAPACITY":   Gb2Sector(size),
	}

	url := "/filesystem"
	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var fileSystem FileSystem

	if err := handleReponse(resp, &fileSystem); err != nil {
		return nil, err
	}

	return &fileSystem.Data, nil
}

func (c *Client) deleteFS(fsID string) error {
	url := "/filesystem/" + fsID
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete DeleteError
	if err := handleReponse(resp, &errDelete); err != nil {
		return err
	}

	return nil
}

func (c *Client) getAllLogicalPort() ([]LogicalPortData, error) {
	url := "/LIF"
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var logicalPortList LogicalPortList

	if err := handleReponse(resp, &logicalPortList); err != nil {
		return nil, err
	}

	return logicalPortList.Data, nil
}

func (c *Client) getFileSystem(fsid string) (*FileSystemData, error) {
	url := "/filesystem/" + fsid
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fileSystem FileSystem
	if err := handleReponse(resp, &fileSystem); err != nil {
		return nil, err
	}

	return &fileSystem.Data, nil
}

func (c *Client) getFileSystemByName(name string) ([]FileSystemData, error) {
	url := "/filesystem?filter=NAME::" + name
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fsList FileSystemList

	if err := handleReponse(resp, &fsList); err != nil {
		return nil, err
	}

	return fsList.Data, nil
}

func (c *Client) ListStoragePools() ([]StoragePool, error) {
	resp, err := c.request("/storagepool", "GET", nil)
	if err != nil {
		return nil, err
	}

	var pools StoragePoolList
	if err := handleReponse(resp, &pools); err != nil {
		return nil, err
	}

	return pools.Data, nil
}

func (c *Client) createSnapshot(fsID, snapName string) (*FSSnapshotData, error) {
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

	var fsSnapshot FSSnapshot
	if err := handleReponse(resp, &fsSnapshot); err != nil {
		return nil, err
	}

	return &fsSnapshot.Data, nil
}

func (c *Client) listSnapshots(fsID string) ([]FSSnapshotData, error) {
	url := "/FSSNAPSHOT?sortby=TIMESTAMP,d&range=[0-100]&PARENTID=" + fsID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fsSnapshotList FSSnapshotList
	if err := handleReponse(resp, &fsSnapshotList); err != nil {
		return nil, err
	}

	return fsSnapshotList.Data, nil
}

func (c *Client) deleteFSSnapshot(snapID string) error {
	url := "/FSSNAPSHOT/" + snapID
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete DeleteError
	if err := handleReponse(resp, &errDelete); err != nil {
		return err
	}

	return nil
}

func (c *Client) showFSSnapshot(snapID string) (*FSSnapshotData, error) {
	url := "/FSSNAPSHOT/" + snapID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fsSnapshot FSSnapshot
	if err := json.Unmarshal(resp, &fsSnapshot); err != nil {
		return nil, err
	}

	return &fsSnapshot.Data, nil
}

func (c *Client) getAllFilesystem() ([]FileSystemData, error) {
	url := "/filesystem"
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var fsList FileSystemList
	if err := handleReponse(resp, &fsList); err != nil {
		return nil, err
	}

	return fsList.Data, nil
}

func (c *Client) getAccessCount(shareID, shareClientType string) (string, error) {
	url := fmt.Sprintf("/%s/count?filter=PARENTID::%s", shareClientType, shareID)
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return "", err
	}

	var count shareAuthClientCount
	if err := handleReponse(resp, &count); err != nil {
		return "", err
	}

	return count.Data.Counter, nil
}

func (c *Client) getAccessFromShare(shareID, accessTo, shareProto string) (string, error) {
	shareClientType, err := c.getShareClientType(shareProto)
	if err != nil {
		return "", err
	}

	count, err := c.getAccessCount(shareID, shareClientType)
	if err != nil {
		return "", err
	}

	counter, _ := strconv.Atoi(count)
	rangeBegin := 0

	for counter > 0 {
		accessRange, err := c.getAccessFromShareRange(shareID, shareClientType, rangeBegin)
		if err != nil {
			return "", nil
		}
		for _, v := range accessRange {
			if v.Name == accessTo {
				return v.ID, nil
			}
		}

		rangeBegin += 100
		counter -= 100
	}

	return "", nil
}

func (c *Client) getAccessFromShareRange(shareID, shareClientType string, rangeBegin int) ([]ShareAuthClientData, error) {
	rangeEnd := rangeBegin + 100
	url := fmt.Sprintf("/%s?filter=PARENTID::%s&range=[%d-%d]", shareClientType, shareID, rangeBegin, rangeEnd)
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}
	var shareClientList ShareAuthClientList

	if err := handleReponse(resp, &shareClientList); err != nil {
		return nil, err
	}

	return shareClientList.Data, nil
}

func (c *Client) getShareClientType(shareProto string) (string, error) {
	switch shareProto {
	case NFSProto:
		return "NFS_SHARE_AUTH_CLIENT", nil
	case CIFSProto:
		return "CIFS_SHARE_AUTH_CLIENT", nil

	}

	return "", fmt.Errorf("invalid NAS protocol supplied: %s", shareProto)
}

func (c *Client) removeAccessFromShare(accessID, shareProto string) error {
	shareClientType, err := c.getShareClientType(shareProto)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/%s/%s", shareClientType, accessID)
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete DeleteError
	if err := handleReponse(resp, &errDelete); err != nil {
		return err
	}

	return nil
}

func (c *Client) logout() error {
	_, err := c.request("/sessions", "DELETE", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) request(url, method string, reqParams interface{}) ([]byte, error) {
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
			return nil, fmt.Errorf("failed to marshal the request parameters, url is %s, error is %v", callUrl, err)
		}
	}

	req, err := http.NewRequest(strings.ToUpper(method), callUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to initiate the request, url is %s, error is %v", callUrl, err)
	}

	// initiate the header
	for k, v := range c.header {
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

	if c.cookie == "" && resp.Header != nil {
		if cookie := resp.Header.Get("set-cookie"); cookie != "" {
			c.cookie = cookie
		}
	}

	return respContent, nil
}
