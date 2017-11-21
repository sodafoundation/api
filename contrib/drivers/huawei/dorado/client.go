// Copyright (c) 2017 OpenSDS Authors.
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

package dorado

import (
	"encoding/json"
	"errors"

	"fmt"

	"github.com/astaxie/beego/httplib"
	log "github.com/golang/glog"
)

const (
	ThickLuntype      = 0
	ThinLuntype       = 1
	MaxNameLength     = 31
	MaxVolDescription = 170
	PortNumPerContr   = 2
	PwdExpired        = 3
	PwdReset          = 4
)
const (
	ErrorConnectToServer      = -403
	ErrorUnauthorizedToServer = -401
)
const urlPathPrefix = "/deviceManager/rest"

type ArrayInnerError struct {
	Data interface{} `json:"data"`
	Err  Error       `json:"error"`
}

func (e *ArrayInnerError) Error() string {
	return fmt.Sprintf("Array internal error, error code:%d, description:%s",
		e.Err.Code, e.Err.Description)
}

type HttpError struct {
	code int
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Http error, error code:%v", e.code)
}

type DoradoClient struct {
	user       string
	passwd     string
	endpoints  []string
	urlPrefix  string
	deviceId   string
	iBaseToken string
}

func NewClient(user, passwd string, endpoints []string) (*DoradoClient, error) {
	c := &DoradoClient{
		user:      user,
		passwd:    passwd,
		endpoints: endpoints,
	}
	err := c.login()
	return c, err
}

func (c *DoradoClient) Destroy() error {
	return c.logout()
}

func (c *DoradoClient) doRequest(method, url string, in interface{}) ([]byte, error) {
	req := httplib.NewBeegoRequest(url, method)
	req.Header("Connection", "keep-alive")
	req.Header("Content-Type", "application/json")
	if c.iBaseToken != "" {
		req.Header("iBaseToken", c.iBaseToken)
	}

	if in != nil {
		body, _ := json.Marshal(in)
		req.Body(body)
	}

	resp, err := req.Response()
	if err != nil {
		log.Errorf("Get Http response error:", err)
		return nil, err
	}
	//Handle the http status code
	switch resp.StatusCode {
	case 200, 201, 202, 204, 206:
		break
	default:
		err := &HttpError{resp.StatusCode}
		log.Error(err)
		return nil, err
	}

	b, err := req.Bytes()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	inErr := &ArrayInnerError{}
	json.Unmarshal(b, inErr)
	if inErr.Err.Code != 0 {
		log.Error(inErr)
		return nil, inErr
	}
	return b, nil
}

func (c *DoradoClient) request(method, url string, in, out interface{}) error {
	var b []byte
	var err error
	for i := 0; i < 2; i++ {
		b, err = c.doRequest(method, c.urlPrefix+url, in)
		if err == nil {
			break
		}
		if inErr, ok := err.(*ArrayInnerError); ok {
			errCode := inErr.Err.Code
			if errCode == ErrorConnectToServer || errCode == ErrorUnauthorizedToServer {
				log.Error("Can't open the recent url, relogin...")
				if err = c.login(); err == nil {
					log.Info("Relogin success")
					continue
				}
			}
			err = inErr
		}
		log.Errorf("Request %d times error:%v", i+1, err)
		return err
	}

	if out != nil {
		json.Unmarshal(b, out)
	}
	return nil
}

func (c *DoradoClient) login() error {
	data := map[string]string{
		"username": c.user,
		"password": c.passwd,
		"scope":    "0",
	}
	c.deviceId = ""
	for _, ep := range c.endpoints {
		url := ep + urlPathPrefix + "/xxxxx/sessions"
		auth := &AuthResp{}
		b, err := c.doRequest("POST", url, data)
		if err != nil {
			log.Error("Login failed,", err)
			continue
		}
		json.Unmarshal(b, auth)
		c.iBaseToken = auth.Data.IBaseToken
		if auth.Data.AccountState == PwdReset || auth.Data.AccountState == PwdExpired {
			msg := "Password has expired or must be reset,please change the password."
			log.Error(msg)
			c.logout()
			return errors.New(msg)
		}
		if auth.Data.DeviceId == "" {
			continue
		}
		c.deviceId = auth.Data.DeviceId
		c.urlPrefix = ep + urlPathPrefix + "/" + auth.Data.DeviceId
		// Get the first controller that can be connected, then break
		break
	}

	if c.deviceId == "" {
		msg := "Failed to login with all rest URLs"
		log.Error(msg)
		return errors.New(msg)
	}
	return nil
}

func (c *DoradoClient) logout() error {
	url := c.urlPrefix + "/sessions"
	if c.urlPrefix == "" {
		return nil
	}
	return c.request("DELETE", url, nil, nil)
}

func (c *DoradoClient) CreateVolume(name string, sectors int64, desc string) (*Lun, error) {
	data := map[string]interface{}{
		"NAME":        name,
		"CAPACITY":    sectors,
		"DESCRIPTION": desc,
		"ALLOCTYPE":   1,
		"PARENTID":    "0",
		"WRITEPOLICY": 1,
	}
	lun := &LunResp{}
	err := c.request("POST", "/lun", data, lun)
	return &lun.Data, err
}

func (c *DoradoClient) GetVolume(id string) (*Lun, error) {
	lun := &LunResp{}
	err := c.request("GET", "/lun/"+id, nil, lun)
	if err != nil {
		return nil, err
	}
	return &lun.Data, err
}

func (c *DoradoClient) DeleteVolume(id string) error {
	lun := &LunResp{}
	err := c.request("DELETE", "/lun/"+id, nil, lun)
	return err
}

func (c *DoradoClient) CreateSnapshot(volId, name, desc string) (*Snapshot, error) {
	data := map[string]interface{}{
		"PARENTTYPE":  11,
		"PARENTID":    volId,
		"NAME":        name,
		"DESCRIPTION": desc,
	}
	snap := &SnapshotResp{}
	err := c.request("POST", "/snapshot", data, snap)
	return &snap.Data, err
}

func (c *DoradoClient) GetSnapshot(id string) (*Snapshot, error) {
	snap := &SnapshotResp{}
	err := c.request("GET", "/snapshot/"+id, nil, snap)
	return &snap.Data, err
}

func (c *DoradoClient) DeleteSnapshot(id string) error {
	return c.request("GET", "/snapshot/"+id, nil, nil)
}

func (c *DoradoClient) ListStoragePools() ([]StoragePool, error) {
	pools := &StoragePoolsResp{}
	err := c.request("GET", "storagepool?range=[0-100]", nil, pools)
	return pools.Data, err
}

