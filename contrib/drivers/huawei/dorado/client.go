// Copyright 2017 The OpenSDS Authors.
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

package dorado

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/astaxie/beego/httplib"
	log "github.com/golang/glog"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/pwd"
)

type ArrayInnerError struct {
	Data interface{} `json:"data"`
	Err  Error       `json:"error"`
}

func (e *ArrayInnerError) Error() string {
	return fmt.Sprintf("Array internal error, error code:%d, description:%s",
		e.Err.Code, e.Err.Description)
}

type NotFoundError struct {
	name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("can not find %s", e.name)
}

func IsNotFoundError(err error) bool {
	if err != nil {
		_, ok := err.(*NotFoundError)
		return ok
	}
	return false
}

type DoradoClient struct {
	user       string
	passwd     string
	endpoints  []string
	urlPrefix  string
	deviceId   string
	iBaseToken string
	cookie     string
	insecure   bool
}

func NewClient(opt *AuthOptions) (*DoradoClient, error) {
	endpoints := strings.Split(opt.Endpoints, ",")
	var pwdCiphertext = opt.Password

	if opt.EnableEncrypted {
		// Decrypte the password
		pwdTool := pwd.NewPwdEncrypter(opt.PwdEncrypter)
		password, err := pwdTool.Decrypter(pwdCiphertext)
		if err != nil {
			return nil, err
		}
		pwdCiphertext = password
	}

	c := &DoradoClient{
		user:      opt.Username,
		passwd:    pwdCiphertext,
		endpoints: endpoints,
		insecure:  opt.Insecure,
	}
	err := c.login()
	return c, err
}

func (c *DoradoClient) Destroy() error {
	return c.logout()
}

func (c *DoradoClient) doRequest(method, url string, in interface{}) ([]byte, http.Header, error) {
	req := httplib.NewBeegoRequest(url, method)
	req.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: c.insecure,
	})
	req.Header("Connection", "keep-alive")
	req.Header("Content-Type", "application/json;charset=utf-8")
	req.Header("iBaseToken", c.iBaseToken)
	req.Header("Cookie", c.cookie)

	if in != nil {
		body, _ := json.Marshal(in)
		req.Body(body)
	}

	resp, err := req.Response()
	if err != nil {
		log.Errorf("Do http request failed, method: %s\n url: %s\n error: %v", method, url, err)
		return nil, nil, err
	}

	b, err := req.Bytes()
	if err != nil {
		log.Errorf("Get byte[] from response failed, method: %s\n url: %s\n error: %v", method, url, err)
		return nil, nil, err
	}

	inErr := &ArrayInnerError{}
	json.Unmarshal(b, inErr)
	if inErr.Err.Code != 0 {
		log.Errorf("Get error info from response failed, method: %s\n url: %s\n error: %v", method, url, inErr)
		return nil, nil, inErr
	}
	return b, resp.Header, nil
}

func (c *DoradoClient) request(method, url string, in, out interface{}) error {
	var b []byte
	var err error
	for i := 0; i < 2; i++ {
		// For debugging
		log.V(5).Infof("URL:%s %s\n BODY:%v", method, c.urlPrefix+url, in)
		b, _, err = c.doRequest(method, c.urlPrefix+url, in)
		if err == nil {
			break
		}
		log.Errorf("URL:%s %s\n BODY:%v", method, c.urlPrefix+url, in)
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
		// This will print tones of info, so set to level 8
		log.V(8).Infof("Response Body: %s", string(b))
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
		url := ep + "/xxxxx/sessions"
		auth := &AuthResp{}
		b, header, err := c.doRequest("POST", url, data)
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
		c.urlPrefix = ep + "/" + auth.Data.DeviceId
		// Get the first controller that can be connected, then break

		c.cookie = header.Get("set-cookie")
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
	if c.urlPrefix == "" {
		return nil
	}
	return c.request("DELETE", "/sessions", nil, nil)
}

func (c *DoradoClient) CreateVolume(name string, size int64, desc string, poolId string) (*Lun, error) {
	data := map[string]interface{}{
		"NAME":        name,
		"CAPACITY":    Gb2Sector(size),
		"DESCRIPTION": desc,
		"ALLOCTYPE":   ThinLunType,
		"PARENTID":    poolId,
		"WRITEPOLICY": 1,
	}
	lun := &LunResp{}
	err := c.request("POST", "/lun", data, lun)
	return &lun.Data, err
}

func (c *DoradoClient) CreateLunCopy(name, srcid, tgtid, copyspeed string) (string, error) {
	url := "/luncopy"
	if !utils.Contains(LunCopySpeedTypes, copyspeed) {
		log.Warningf("The copy speed %s is invalid, using Medium Speed instead", copyspeed)
		copyspeed = LunCopySpeedMedium
	}
	data := map[string]interface{}{
		"TYPE":        ObjectTypeLunCopy,
		"NAME":        name,
		"DESCRIPTION": name,
		"COPYSPEED":   copyspeed,
		"LUNCOPYTYPE": "1",
		"SOURCELUN":   fmt.Sprintf("INVALID;%s;INVALID;INVALID;INVALID", srcid),
		"TARGETLUN":   fmt.Sprintf("INVALID;%s;INVALID;INVALID;INVALID", tgtid),
	}
	lun := &LunResp{}
	err := c.request("POST", url, data, lun)
	if err != nil {
		log.Errorf("Create LunCopy failed :%v", err)
		return "", err
	}
	return lun.Data.Id, err
}
func (c *DoradoClient) GetLunInfo(id string) (*Lun, error) {
	url := "/LUNCOPY/" + id
	lun := &LunResp{}
	err := c.request("GET", url, nil, lun)
	if err != nil {
		return nil, err
	}
	return &lun.Data, nil
}
func (c *DoradoClient) StartLunCopy(luncopyid string) error {
	url := "/LUNCOPY/start"
	data := map[string]interface{}{
		"TYPE": ObjectTypeLunCopy,
		"ID":   luncopyid,
	}
	err := c.request("PUT", url, data, nil)
	return err
}

func (c *DoradoClient) DeleteLunCopy(luncopyid string) error {
	url := "/LUNCOPY/" + luncopyid
	err := c.request("DELETE", url, nil, nil)
	return err
}

func (c *DoradoClient) GetVolume(id string) (*Lun, error) {
	lun := &LunResp{}
	err := c.request("GET", "/lun/"+id, nil, lun)
	if err != nil {
		return nil, err
	}
	return &lun.Data, err
}

func (c *DoradoClient) GetVolumeByName(name string) (*Lun, error) {
	lun := &LunResp{}
	err := c.request("GET", "/lun?filter=NAME::"+name, nil, lun)
	if err != nil {
		return nil, err
	}
	return &lun.Data, err
}
func (c *DoradoClient) DeleteVolume(id string) error {
	err := c.request("DELETE", "/lun/"+id, nil, nil)
	return err
}

// ExtendVolume ...
func (c *DoradoClient) ExtendVolume(size int64, id string) error {
	data := map[string]interface{}{
		"CAPACITY": Gb2Sector(size),
		"ID":       id,
	}

	err := c.request("PUT", "/lun/expand", data, nil)
	return err
}

func (c *DoradoClient) CheckLunExist(id, wwn string) bool {
	lun := &LunResp{}
	err := c.request("GET", "/lun/"+id, nil, lun)
	if err != nil {
		return false
	}
	if wwn != "" && lun.Data.Wwn != wwn {
		log.Infof("LUN Id %s with WWN %s does not on the array.", id, wwn)
		return false
	}
	return true
}

func (c *DoradoClient) CreateSnapshot(lunId, name, desc string) (*Snapshot, error) {
	data := map[string]interface{}{
		"PARENTTYPE":  ObjectTypeLun,
		"PARENTID":    lunId,
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

func (c *DoradoClient) GetSnapshotByName(name string) (*Snapshot, error) {
	snap := &SnapshotsResp{}
	err := c.request("GET", "/snapshot?filter=NAME::"+name, nil, snap)
	return &snap.Data[0], err
}

func (c *DoradoClient) DeleteSnapshot(id string) error {
	return c.request("DELETE", "/snapshot/"+id, nil, nil)
}

func (c *DoradoClient) ListStoragePools() ([]StoragePool, error) {
	pools := &StoragePoolsResp{}
	err := c.request("GET", "/storagepool?range=[0-100]", nil, pools)
	return pools.Data, err
}

func (c *DoradoClient) ListAllStoragePools() ([]StoragePool, error) {
	pools := &StoragePoolsResp{}
	err := c.request("GET", "/storagepool", nil, pools)
	return pools.Data, err
}

func (c *DoradoClient) GetPoolIdByName(poolName string) (string, error) {
	pools, err := c.ListAllStoragePools()
	if err != nil {
		return "", err
	}
	for _, p := range pools {
		if p.Name == poolName {
			return p.Id, nil
		}
	}
	return "", fmt.Errorf("not found specified pool '%s'", poolName)
}

func (c *DoradoClient) AddHostWithCheck(hostInfo *pb.HostInfo) (string, error) {
	hostName := EncodeHostName(hostInfo.Host)

	hostId, _ := c.GetHostIdByName(hostInfo.Host)
	if hostId != "" {
		return hostId, nil
	}

	reqBody := map[string]interface{}{
		"NAME":            hostName,
		"OPERATIONSYSTEM": 0, /*linux*/
		"IP":              hostInfo.Ip,
	}
	hostResp := &HostResp{}
	if err := c.request("POST", "/host", reqBody, hostResp); err != nil {
		log.Errorf("Create host failed, host name: %s, error: %v", hostName, err)
		return "", err
	}

	if hostResp.Data.Id != "" {
		return hostResp.Data.Id, nil
	}

	log.Errorf("Create host failed by host name: %s, error code:%d, description:%s",
		hostInfo.Host, hostResp.Error.Code, hostResp.Error.Description)
	return "", fmt.Errorf("Create host failed by host name: %s, error code:%d, description:%s",
		hostInfo.Host, hostResp.Error.Code, hostResp.Error.Description)
}

func (c *DoradoClient) GetHostIdByName(hostName string) (string, error) {
	hostName = EncodeHostName(hostName)
	hostsResp := &HostsResp{}

	if err := c.request("GET", "/host?filter=NAME::"+hostName, nil, hostsResp); err != nil {
		log.Errorf("Get host failed by host name: %s, error: %v", hostName, err)
		return "", err
	}

	if len(hostsResp.Data) > 0 {
		return hostsResp.Data[0].Id, nil
	}

	log.Errorf("Get host failed by host name: %s, error code:%d, description:%s",
		hostName, hostsResp.Error.Code, hostsResp.Error.Description)
	return "", &NotFoundError{name: hostName}
}

func (c *DoradoClient) AddInitiatorToHostWithCheck(hostId, initiatorName string) error {

	if !c.IsArrayContainInitiator(initiatorName) {
		if err := c.AddInitiatorToArray(initiatorName); err != nil {
			return err
		}
	}
	if !c.IsHostContainInitiator(hostId, initiatorName) {
		if err := c.AddInitiatorToHost(hostId, initiatorName); err != nil {
			return err
		}
	}
	return nil
}

func (c *DoradoClient) IsArrayContainInitiator(initiatorName string) bool {
	initiatorResp := &InitiatorResp{}

	if err := c.request("GET", "/iscsi_initiator/"+initiatorName, nil, initiatorResp); err != nil {
		log.Errorf("Get iscsi initiator failed by initiator name: %s, error: %v", initiatorName, err)
		return false
	}

	if initiatorResp.Data.Id == "" {
		log.Infof("Array does not contains the initiator: %s", initiatorName)
		return false
	}

	log.Infof("Array contains the initiator: %s", initiatorName)
	return true
}

func (c *DoradoClient) IsHostContainInitiator(hostId, initiatorName string) bool {
	initiatorsResp := &InitiatorsResp{}

	if err := c.request("GET", "/iscsi_initiator?ISFREE=false&PARENTID="+hostId, nil, initiatorsResp); err != nil {
		log.Errorf("Get iscsi initiator failed by host id: %s, initiator name: %s, error: %v", hostId, initiatorName, err)
		return false
	}

	for _, initiator := range initiatorsResp.Data {
		if initiator.Id == initiatorName {
			log.Infof("Host:%s contains the initiator: %s", hostId, initiatorName)
			return true
		}
	}

	log.Infof("Host:%s does not contains the initiator: %s", hostId, initiatorName)
	return false
}

func (c *DoradoClient) AddInitiatorToArray(initiatorName string) error {

	reqBody := map[string]interface{}{
		"ID": initiatorName,
	}
	initiatorResp := &InitiatorResp{}

	if err := c.request("POST", "/iscsi_initiator", reqBody, initiatorResp); err != nil {
		log.Errorf("Create iscsi initiator failed, initiator name: %s, error: %v", initiatorName, err)
		return err
	}

	if initiatorResp.Error.Code != 0 {
		log.Errorf("Add iscsi initiator to array failed, error code:%d, description:%s",
			initiatorResp.Error.Code, initiatorResp.Error.Description)
		return fmt.Errorf("code: %d, description: %s",
			initiatorResp.Error.Code, initiatorResp.Error.Description)
	}

	log.Infof("Create the initiator: %s successfully.", initiatorName)
	return nil
}

func (c *DoradoClient) AddInitiatorToHost(hostId, initiatorName string) error {

	reqBody := map[string]interface{}{
		"ID":       initiatorName,
		"PARENTID": hostId,
	}
	initiatorResp := &InitiatorResp{}

	if err := c.request("PUT", "/iscsi_initiator/"+initiatorName, reqBody, initiatorResp); err != nil {
		log.Errorf("Modify iscsi initiator failed, initiator name: %s, error: %v", initiatorName, err)
		return err
	}

	if initiatorResp.Error.Code != 0 {
		log.Errorf("Add iscsi initiator to host failed, error code:%d, description:%s",
			initiatorResp.Error.Code, initiatorResp.Error.Description)
		return fmt.Errorf("code: %d, description: %s",
			initiatorResp.Error.Code, initiatorResp.Error.Description)
	}

	log.Infof("Add the initiator: %s to host: %s successfully.", initiatorName, hostId)
	return nil
}

func (c *DoradoClient) AddHostToHostGroup(hostId string) (string, error) {

	hostGrpName := PrefixHostGroup + hostId
	hostGrpId, err := c.CreateHostGroupWithCheck(hostGrpName)
	if err != nil {
		log.Errorf("Create host group witch check failed, host group id: %s, error: %v", hostGrpId, err)
		return "", err
	}

	contained := c.IsHostGroupContainHost(hostGrpId, hostId)
	if contained {
		return hostGrpId, nil
	}

	err = c.AssociateHostToHostGroup(hostGrpId, hostId)
	if err != nil {
		log.Errorf("Associate host to host group failed, host group id: %s, host id: %s, error: %v", hostGrpId, hostId, err)
		return "", err
	}

	return hostGrpId, nil
}

func (c *DoradoClient) CreateHostGroupWithCheck(hostGrpName string) (string, error) {

	hostGrpId, _ := c.FindHostGroup(hostGrpName)
	if hostGrpId != "" {
		return hostGrpId, nil
	}

	hostGrpId, err := c.CreateHostGroup(hostGrpName)
	if err != nil {
		log.Errorf("Create host group with name: %s failed, error: %v", hostGrpName, err)
		return "", err
	}
	return hostGrpId, nil
}

func (c *DoradoClient) FindHostGroup(groupName string) (string, error) {

	hostGrpsResp := &HostGroupsResp{}

	if err := c.request("GET", "/hostgroup?filter=NAME::"+groupName, nil, hostGrpsResp); err != nil {
		log.Errorf("Get host groups failed by filter name: %s, error: %v", groupName, err)
		return "", err
	}

	if hostGrpsResp.Error.Code != 0 {
		log.Errorf("Get host groups failed by filter name: %s, error code:%d, description:%s",
			groupName, hostGrpsResp.Error.Code, hostGrpsResp.Error.Description)
		return "", fmt.Errorf("code: %d, description: %s",
			hostGrpsResp.Error.Code, hostGrpsResp.Error.Description)
	}

	if len(hostGrpsResp.Data) == 0 {
		log.Infof("No host group with name %s was found.", groupName)
		return "", &NotFoundError{name: groupName}
	}

	return hostGrpsResp.Data[0].Id, nil
}

func (c *DoradoClient) CreateHostGroup(groupName string) (string, error) {

	reqBody := map[string]interface{}{
		"NAME": groupName,
	}
	hostGrpResp := &HostGroupResp{}

	if err := c.request("POST", "/hostgroup", reqBody, hostGrpResp); err != nil {
		log.Errorf("Create host group failed, group name: %s, error: %v", groupName, err)
		return "", err
	}

	if hostGrpResp.Error.Code != 0 {
		log.Errorf("Create host group failed, group name: %s, error code:%d, description:%s",
			groupName, hostGrpResp.Error.Code, hostGrpResp.Error.Description)
		return "", fmt.Errorf("code: %d, description: %s",
			hostGrpResp.Error.Code, hostGrpResp.Error.Description)
	}

	return hostGrpResp.Data.Id, nil
}

func (c *DoradoClient) IsHostGroupContainHost(hostGrpId, hostId string) bool {
	hostsResp := &HostsResp{}

	if err := c.request("GET", "/host/associate?ASSOCIATEOBJTYPE=14&ASSOCIATEOBJID="+hostGrpId, nil, hostsResp); err != nil {
		log.Errorf("List hosts failed by parent id: %s, error: %v", hostGrpId, err)
		return false
	}

	for _, host := range hostsResp.Data {
		if host.Id == hostId {
			log.Infof("HostGroup: %s contains the host: %s", hostGrpId, hostId)
			return true
		}
	}

	log.Infof("HostGroup: %s does not contain the host: %s", hostGrpId, hostId)
	return false
}

func (c *DoradoClient) AssociateHostToHostGroup(hostGrpId, hostId string) error {

	reqBody := map[string]interface{}{
		"ID":               hostGrpId,
		"ASSOCIATEOBJTYPE": ObjectTypeHost,
		"ASSOCIATEOBJID":   hostId,
	}
	resp := &GenericResult{}

	if err := c.request("POST", "/hostgroup/associate", reqBody, resp); err != nil {
		log.Errorf("Associate host:%s to host group:%s failed, error: %v", hostId, hostGrpId, err)
		return err
	}

	if resp.Error.Code != 0 {
		log.Errorf("Associate host:%s to host group:%s failed, error code:%d, description:%s",
			hostId, hostGrpId, resp.Error.Code, resp.Error.Description)
		return fmt.Errorf("code: %d, description: %s",
			resp.Error.Code, resp.Error.Description)
	}

	return nil
}

func (c *DoradoClient) DoMapping(lunId, hostGrpId, hostId string) error {

	var err error
	// Find or create lun group and add lun into lun group.
	lunGrpName := PrefixLunGroup + hostId
	lunGrpId, _ := c.FindLunGroup(lunGrpName)
	if lunGrpId == "" {
		lunGrpId, err = c.CreateLunGroup(lunGrpName)
		if err != nil {
			log.Errorf("Create lun group failed, group name:%s, error: %v", lunGrpName, err)
			return err
		}
	}
	if !c.IsLunGroupContainLun(lunGrpId, lunId) {
		if err := c.AssociateLunToLunGroup(lunGrpId, lunId); err != nil {
			log.Errorf("Associate lun to lun group failed, group id:%s, lun id:%s, error: %v", lunGrpId, lunId, err)
			return err
		}
	}

	// Find or create mapping view
	mappingViewName := PrefixMappingView + hostId
	mappingViewId, _ := c.FindMappingView(mappingViewName)
	if mappingViewId == "" {
		mappingViewId, err = c.CreateMappingView(mappingViewName)
		if err != nil {
			log.Errorf("Create mapping view failed, view name:%s, error: %v", mappingViewName, err)
			return err
		}
	}

	// Associate host group and lun group to mapping view.
	if !c.IsMappingViewContainHostGroup(mappingViewId, hostGrpId) {
		if err := c.AssocateHostGroupToMappingView(mappingViewId, hostGrpId); err != nil {
			log.Errorf("Assocate host group to mapping view failed, view id:%s, host group id:%s, error: %v",
				mappingViewId, hostGrpId, err)
			return err
		}
	}
	if !c.IsMappingViewContainLunGroup(mappingViewId, lunGrpId) {
		if err := c.AssocateLunGroupToMappingView(mappingViewId, lunGrpId); err != nil {
			log.Errorf("Assocate lun group to mapping view failed, view id:%s, lun group id:%s, error: %v",
				mappingViewId, lunGrpId, err)
			return err
		}
	}

	log.Infof("DoMapping successfully, with params lunId:%s, hostGrpId:%s, hostId:%s",
		lunId, lunGrpId, hostId)
	return nil
}

func (c *DoradoClient) FindLunGroup(groupName string) (string, error) {

	lunGrpsResp := &LunGroupsResp{}

	if err := c.request("GET", "/lungroup?filter=NAME::"+groupName, nil, lunGrpsResp); err != nil {
		log.Errorf("Get lun groups failed by filter name: %s, error: %v", groupName, err)
		return "", err
	}

	if lunGrpsResp.Error.Code != 0 {
		log.Errorf("Get lun groups failed by filter name: %s, error code:%d, description:%s",
			groupName, lunGrpsResp.Error.Code, lunGrpsResp.Error.Description)
		return "", fmt.Errorf("code: %d, description: %s",
			lunGrpsResp.Error.Code, lunGrpsResp.Error.Description)
	}

	if len(lunGrpsResp.Data) == 0 {
		log.Infof("No lun group with name %s was found.", groupName)
		return "", &NotFoundError{name: groupName}
	}

	return lunGrpsResp.Data[0].Id, nil
}

func (c *DoradoClient) FindMappingView(name string) (string, error) {

	mvsResp := &MappingViewsResp{}

	if err := c.request("GET", "/mappingview?filter=NAME::"+name, nil, mvsResp); err != nil {
		log.Errorf("Get mapping views failed by filter name: %s, error: %v", name, err)
		return "", err
	}

	if mvsResp.Error.Code != 0 {
		log.Errorf("Get mapping views failed by filter name: %s, error code:%d, description:%s",
			name, mvsResp.Error.Code, mvsResp.Error.Description)
		return "", fmt.Errorf("code: %d, description: %s",
			mvsResp.Error.Code, mvsResp.Error.Description)
	}

	if len(mvsResp.Data) == 0 {
		log.Infof("No mapping view with name %s was found.", name)
		return "", &NotFoundError{name: name}
	}

	return mvsResp.Data[0].Id, nil
}

func (c *DoradoClient) CreateLunGroup(groupName string) (string, error) {

	reqBody := map[string]interface{}{
		"NAME":      groupName,
		"APPTYPE":   0,
		"GROUPTYPE": 0,
	}
	lunGrpResp := &LunGroupResp{}

	if err := c.request("POST", "/lungroup", reqBody, lunGrpResp); err != nil {
		log.Errorf("Create lun group failed, group name: %s, error: %v", groupName, err)
		return "", err
	}

	if lunGrpResp.Error.Code != 0 {
		log.Errorf("Create lun group failed, group name: %s, error code:%d, description:%s",
			groupName, lunGrpResp.Error.Code, lunGrpResp.Error.Description)
		return "", fmt.Errorf("code: %d, description: %s",
			lunGrpResp.Error.Code, lunGrpResp.Error.Description)
	}

	return lunGrpResp.Data.Id, nil
}

func (c *DoradoClient) CreateMappingView(name string) (string, error) {

	reqBody := map[string]interface{}{
		"NAME": name,
	}
	mvResp := &MappingViewResp{}

	if err := c.request("POST", "/mappingview", reqBody, mvResp); err != nil {
		log.Errorf("Create mapping view failed, view name: %s, error: %v", name, err)
		return "", err
	}

	if mvResp.Error.Code != 0 {
		log.Errorf("Create mapping view failed, view name: %s, error code:%d, description:%s",
			name, mvResp.Error.Code, mvResp.Error.Description)
		return "", fmt.Errorf("code: %d, description: %s",
			mvResp.Error.Code, mvResp.Error.Description)
	}

	return mvResp.Data.Id, nil
}

func (c *DoradoClient) IsLunGroupContainLun(lunGrpId, lunId string) bool {
	lunsResp := &LunsResp{}

	if err := c.request("GET", "/lun/associate?ASSOCIATEOBJTYPE=256&ASSOCIATEOBJID="+lunGrpId, nil, lunsResp); err != nil {
		log.Errorf("List luns failed by lun group id: %s, error: %v", lunGrpId, err)
		return false
	}

	for _, lun := range lunsResp.Data {
		if lun.Id == lunId {
			log.Infof("LunGroup: %s contains the lun: %s", lunGrpId, lunId)
			return true
		}
	}

	log.Infof("LunGroup: %s does not contain the lun: %s", lunGrpId, lunId)
	return false
}

func (c *DoradoClient) AssociateLunToLunGroup(lunGrpId, lunId string) error {

	reqBody := map[string]interface{}{
		"ID":               lunGrpId,
		"ASSOCIATEOBJTYPE": ObjectTypeLun,
		"ASSOCIATEOBJID":   lunId,
	}
	resp := &GenericResult{}

	if err := c.request("POST", "/lungroup/associate", reqBody, resp); err != nil {
		log.Errorf("Associate lun:%s to lun group:%s failed, error: %v", lunId, lunGrpId, err)
		return err
	}

	if resp.Error.Code != 0 {
		log.Errorf("Associate lun:%s to lun group:%s failed, error code:%d, description:%s",
			lunId, lunGrpId, resp.Error.Code, resp.Error.Description)
		return fmt.Errorf("code: %d, description: %s",
			resp.Error.Code, resp.Error.Description)
	}

	return nil
}

func (c *DoradoClient) IsMappingViewContainHostGroup(viewId, groupId string) bool {
	mvsResp := &MappingViewsResp{}
	if err := c.request("GET", "/mappingview/associate?ASSOCIATEOBJTYPE=14&ASSOCIATEOBJID="+groupId, nil, mvsResp); err != nil {
		log.Errorf("List mapping views failed by host group id: %s, error: %v", groupId, err)
		return false
	}

	for _, view := range mvsResp.Data {
		if view.Id == viewId {
			log.Infof("Mapping view: %s contains the host group: %s", viewId, groupId)
			return true
		}
	}

	log.Infof("Mapping view: %s does not contain the host group: %s", viewId, groupId)
	return false
}

func (c *DoradoClient) AssocateHostGroupToMappingView(viewId, groupId string) error {

	reqBody := map[string]interface{}{
		"ID":               viewId,
		"ASSOCIATEOBJTYPE": ObjectTypeHostGroup,
		"ASSOCIATEOBJID":   groupId,
	}
	resp := &GenericResult{}
	if err := c.request("PUT", "/mappingview/create_associate", reqBody, resp); err != nil {
		log.Errorf("Associate host group:%s to mapping view:%s failed, error: %v", groupId, viewId, err)
		return err
	}

	if resp.Error.Code != 0 {
		log.Errorf("Associate host group:%s to mapping view:%s failed, error code:%d, description:%s",
			groupId, viewId, resp.Error.Code, resp.Error.Description)
		return fmt.Errorf("code: %d, description: %s",
			resp.Error.Code, resp.Error.Description)
	}

	return nil
}

func (c *DoradoClient) IsMappingViewContainLunGroup(viewId, groupId string) bool {
	mvsResp := &MappingViewsResp{}
	if err := c.request("GET", "/mappingview/associate?ASSOCIATEOBJTYPE=256&ASSOCIATEOBJID="+groupId, nil, mvsResp); err != nil {
		log.Errorf("List mapping views failed by lun group id: %s, error: %v", groupId, err)
		return false
	}

	for _, view := range mvsResp.Data {
		if view.Id == viewId {
			log.Infof("Mapping view: %s contains the lun group: %s", viewId, groupId)
			return true
		}
	}

	log.Infof("Mapping view: %s does not contain the lun group: %s", viewId, groupId)
	return false
}

func (c *DoradoClient) AssocateLunGroupToMappingView(viewId, groupId string) error {

	reqBody := map[string]interface{}{
		"ID":               viewId,
		"ASSOCIATEOBJTYPE": ObjectTypeLunGroup,
		"ASSOCIATEOBJID":   groupId,
	}
	resp := &GenericResult{}
	if err := c.request("PUT", "/mappingview/create_associate", reqBody, resp); err != nil {
		log.Errorf("Associate lun group:%s to mapping view:%s failed, error: %v", groupId, viewId, err)
		return err
	}

	if resp.Error.Code != 0 {
		log.Errorf("Associate lun group:%s to mapping view:%s failed, error code:%d, description:%s",
			groupId, viewId, resp.Error.Code, resp.Error.Description)
		return fmt.Errorf("code: %d, description: %s",
			resp.Error.Code, resp.Error.Description)
	}

	return nil
}

func (c *DoradoClient) ListTgtPort() (*IscsiTgtPortsResp, error) {
	resp := &IscsiTgtPortsResp{}
	if err := c.request("GET", "/iscsi_tgt_port", nil, resp); err != nil {
		log.Errorf("Get tgt port failed, error: %v", err)
		return nil, err
	}

	if resp.Error.Code != 0 {
		log.Errorf("Get tgt port failed, error code:%d, description:%s",
			resp.Error.Code, resp.Error.Description)
		return nil, fmt.Errorf("code: %d, description: %s", resp.Error.Code, resp.Error.Description)
	}
	return resp, nil
}

func (c *DoradoClient) ListHostAssociateLuns(hostId string) (*HostAssociateLunsResp, error) {
	resp := &HostAssociateLunsResp{}
	url := fmt.Sprintf("/lun/associate?TYPE=11&ASSOCIATEOBJTYPE=21&ASSOCIATEOBJID=%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *DoradoClient) GetHostLunId(hostId, lunId string) (int, error) {
	resp, err := c.ListHostAssociateLuns(hostId)
	if err != nil {
		return -1, err
	}
	type Metadata struct {
		HostLunId int `json:"HostLUNID"`
	}
	for _, lun := range resp.Data {
		if lun.Id != lunId {
			continue
		}
		md := &Metadata{}
		if err := json.Unmarshal([]byte(lun.AssociateMetadata), md); err != nil {
			log.Error("Decoding json error,", err)
			return -1, err
		}
		return md.HostLunId, nil
	}
	log.Info("Do not find the host lun id, using default id 1")
	return 1, nil
}

func (c *DoradoClient) RemoveLunFromLunGroup(lunGrpId, lunId string) error {
	url := fmt.Sprintf("/lungroup/associate?ID=%s&ASSOCIATEOBJTYPE=11&ASSOCIATEOBJID=%s", lunGrpId, lunId)
	if err := c.request("DELETE", url, nil, nil); err != nil {
		log.Errorf("Remove lun %s from lun group %s failed, %v", lunId, lunGrpId, err)
		return err
	}
	log.Infof("Remove lun %s from lun group %s success", lunId, lunGrpId)
	return nil
}

func (c *DoradoClient) RemoveLunGroupFromMappingView(viewId, lunGrpId string) error {
	if !c.IsMappingViewContainLunGroup(viewId, lunGrpId) {
		log.Infof("Lun group %s has already been removed from mapping view %s", lunGrpId, viewId)
		return nil
	}
	url := "/mappingview/REMOVE_ASSOCIATE"
	data := map[string]interface{}{
		"ASSOCIATEOBJTYPE": ObjectTypeLunGroup,
		"ASSOCIATEOBJID":   lunGrpId,
		"TYPE":             ObjectTypeMappingView,
		"ID":               viewId}
	if err := c.request("PUT", url, data, nil); err != nil {
		log.Errorf("Remove lun group %s from mapping view %s failed", lunGrpId, viewId)
		return err
	}
	log.Infof("Remove lun group %s from mapping view %s success", lunGrpId, viewId)
	return nil
}

func (c *DoradoClient) RemoveHostGroupFromMappingView(viewId, hostGrpId string) error {
	if !c.IsMappingViewContainHostGroup(viewId, hostGrpId) {
		log.Infof("Host group %s has already been removed from mapping view %s", hostGrpId, viewId)
		return nil
	}
	url := "/mappingview/REMOVE_ASSOCIATE"
	data := map[string]interface{}{
		"ASSOCIATEOBJTYPE": ObjectTypeHostGroup,
		"ASSOCIATEOBJID":   hostGrpId,
		"TYPE":             ObjectTypeMappingView,
		"ID":               viewId}
	if err := c.request("PUT", url, data, nil); err != nil {
		log.Errorf("Remove host group %s from mapping view %s failed", hostGrpId, viewId)
		return err
	}
	log.Infof("Remove host group %s from mapping view %s success", hostGrpId, viewId)
	return nil
}

func (c *DoradoClient) RemoveHostFromHostGroup(hostGrpId, hostId string) error {

	url := fmt.Sprintf("/host/associate?TYPE=14&ID=%s&ASSOCIATEOBJTYPE=21&ASSOCIATEOBJID=%s",
		hostGrpId, hostId)
	if err := c.request("DELETE", url, nil, nil); err != nil {
		log.Errorf("Remove host %s from host group %s failed", hostId, hostGrpId)
		return err
	}
	log.Infof("Remove host %s from host group %s success", hostId, hostGrpId)
	return nil
}

func (c *DoradoClient) RemoveIscsiFromHost(initiator string) error {

	url := "/iscsi_initiator/remove_iscsi_from_host"
	data := map[string]interface{}{"TYPE": ObjectTypeIscsiInitiator, "ID": initiator}
	if err := c.request("PUT", url, data, nil); err != nil {
		log.Errorf("Remove initiator %s failed", initiator)
		return err
	}
	log.Infof("Remove initiator %s success", initiator)
	return nil
}

func (c *DoradoClient) DeleteHostGroup(id string) error {
	return c.request("DELETE", "/hostgroup/"+id, nil, nil)
}

func (c *DoradoClient) DeleteLunGroup(id string) error {
	return c.request("DELETE", "/LUNGroup/"+id, nil, nil)
}

func (c *DoradoClient) DeleteHost(id string) error {
	return c.request("DELETE", "/host/"+id, nil, nil)
}

func (c *DoradoClient) DeleteMappingView(id string) error {
	return c.request("DELETE", "/mappingview/"+id, nil, nil)
}

func (c *DoradoClient) GetArrayInfo() (*System, error) {
	sys := &SystemResp{}
	err := c.request("GET", "/system/", nil, sys)
	if err != nil {
		log.Error("Get system info failed,", err)
		return nil, err
	}
	return &sys.Data, nil
}

func (c *DoradoClient) ListRemoteDevices() (*[]RemoteDevice, error) {
	dev := &RemoteDevicesResp{}
	err := c.request("GET", "/remote_device", nil, dev)
	if err != nil {
		log.Error("List remote devices failed,", err)
		return nil, err
	}
	return &dev.Data, nil
}

func (c *DoradoClient) CreatePair(params map[string]interface{}) (*ReplicationPair, error) {
	pair := &ReplicationPairResp{}
	err := c.request("POST", "/REPLICATIONPAIR", params, pair)
	return &pair.Data, err
}

func (c *DoradoClient) GetPair(id string) (*ReplicationPair, error) {
	pair := &ReplicationPairResp{}
	err := c.request("GET", "/REPLICATIONPAIR/"+id, nil, pair)
	if err != nil {
		log.Errorf("Get pair failed, %v", err)
		return nil, err
	}
	return &pair.Data, err
}

func (c *DoradoClient) SwitchPair(id string) error {
	data := map[string]interface{}{"ID": id, "TYPE": ObjectTypeReplicationPair}
	err := c.request("PUT", "/REPLICATIONPAIR/switch", data, nil)
	if err != nil {
		log.Errorf("Switch pair failed, %v", err)
	}
	return err
}

func (c *DoradoClient) SplitPair(id string) error {
	data := map[string]interface{}{"ID": id, "TYPE": ObjectTypeReplicationPair}
	err := c.request("PUT", "/REPLICATIONPAIR/split", data, nil)
	if err != nil {
		log.Errorf("Split pair failed, %v", err)
	}
	return err
}

func (c *DoradoClient) SyncPair(id string) error {
	data := map[string]interface{}{"ID": id, "TYPE": ObjectTypeReplicationPair}
	err := c.request("PUT", "/REPLICATIONPAIR/sync", data, nil)
	if err != nil {
		log.Errorf("Sync pair failed, %v", err)
	}
	return err
}

func (c *DoradoClient) SetPairSecondAccess(id string, access string) error {
	data := map[string]interface{}{"ID": id, "SECRESACCESS": access}
	err := c.request("PUT", "/REPLICATIONPAIR/"+id, data, nil)
	if err != nil {
		log.Errorf("Set pair secondary access failed, %v", err)
	}
	return err
}

func (c *DoradoClient) DeletePair(id string) error {
	return c.request("DELETE", "/REPLICATIONPAIR/"+id, nil, nil)
}

func (c *DoradoClient) CheckPairExist(id string) bool {
	resp := &SimpleResp{}
	err := c.request("GET", "/REPLICATIONPAIR/"+id, nil, resp)
	return err == nil
}

func (c *DoradoClient) GetHostOnlineFCInitiators(hostId string) ([]string, error) {
	resp := &FCInitiatorsResp{}
	url := fmt.Sprintf("/fc_initiator?PARENTTYPE=21&PARENTID=%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get host online fc initiators from host %s failed.", hostId)
		return nil, err
	}

	var initiators []string
	if len(resp.Data) > 0 {
		for _, item := range resp.Data {
			if item.ParentId != "" && item.ParentId == hostId && item.RunningStatus == RunningStatusOnline {
				initiators = append(initiators, item.Id)
			}
		}
	}
	log.Infof("Get host online fc initiators from host %s success.", hostId)
	return initiators, nil
}

func (c *DoradoClient) GetOnlineFreeWWNs() ([]string, error) {
	resp := &FCInitiatorsResp{}
	url := "/fc_initiator?ISFREE=true&range=[0-65535]"
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get online free wwns failed.")
		return nil, err
	}

	var wwns []string
	if len(resp.Data) > 0 {
		for _, item := range resp.Data {
			if item.RunningStatus == RunningStatusOnline {
				wwns = append(wwns, item.Id)
			}
		}
	}

	log.Infof("Get online free wwns success.")
	return wwns, nil
}

func (c *DoradoClient) GetOnlineFCInitiatorOnArray() ([]string, error) {
	resp := &FCInitiatorsResp{}
	url := "/fc_initiator?range=[0-65535]"
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get online FC initiator on array failed.")
		return nil, err
	}

	var fcInitiators []string
	for _, item := range resp.Data {
		if item.RunningStatus == RunningStatusOnline {
			fcInitiators = append(fcInitiators, item.Id)
		}
	}

	log.Infof("Get online fc initiators success.")
	return fcInitiators, nil
}

func (c *DoradoClient) GetHostFCInitiators(hostId string) ([]string, error) {
	resp := &FCInitiatorsResp{}
	url := fmt.Sprintf("/fc_initiator?PARENTTYPE=21&PARENTID=%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get host fc initiators failed.")
		return nil, err
	}

	var initiators []string
	if len(resp.Data) > 0 {
		for _, item := range resp.Data {
			if item.ParentId != "" && item.ParentId == hostId {
				initiators = append(initiators, item.Id)
			}
		}
	}
	return initiators, nil
}

func (c *DoradoClient) GetHostIscsiInitiators(hostId string) ([]string, error) {
	resp := &InitiatorsResp{}
	url := fmt.Sprintf("/iscsi_initiator?PARENTTYPE=21&PARENTID=%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get host iscsi initiators failed.")
		return nil, err
	}

	var initiators []string
	if len(resp.Data) > 0 {
		for _, item := range resp.Data {
			if item.ParentId != "" && item.ParentId == hostId {
				initiators = append(initiators, item.Id)
			}
		}
	}

	log.Infof("Get host iscsi initiators success.")
	return initiators, nil
}

func (c *DoradoClient) IsHostAssociatedToHostgroup(hostId string) (bool, error) {
	resp := &HostResp{}
	url := fmt.Sprintf("/host/%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get host iscsi initiators failed.")
		return false, err
	}

	if resp.Data.IsAddToHostGroup {
		return true, nil
	}

	return false, nil
}

func (c *DoradoClient) RemoveHost(hostId string) error {
	return c.request("DELETE", fmt.Sprintf("/host/%s", hostId), nil, nil)
}

func (c *DoradoClient) AddFCPortTohost(hostId string, wwn string) error {
	url := fmt.Sprintf("/fc_initiator/%s", wwn)
	data := map[string]interface{}{
		"TYPE":       ObjectTypeFcInitiator,
		"ID":         wwn,
		"PARENTTYPE": ObjectTypeHost,
		"PARENTID":   hostId,
	}

	if err := c.request("PUT", url, data, nil); err != nil {
		log.Errorf("Add fc port to host failed.")
		return err
	}

	return nil
}

func (c *DoradoClient) GetIniTargMap(wwns []string) ([]string, map[string][]string, error) {
	initTargMap := make(map[string][]string)
	var tgtPortWWNs []string
	for _, wwn := range wwns {
		tgtwwpns, err := c.getFCTargetWWPNs(wwn)
		if err != nil {
			return nil, nil, err
		}
		if tgtwwpns == nil {
			continue
		}

		initTargMap[wwn] = tgtwwpns
		for _, tgtwwpn := range tgtwwpns {
			if !c.isInStringArray(tgtwwpn, tgtPortWWNs) {
				tgtPortWWNs = append(tgtPortWWNs, tgtwwpn)
			}
		}
	}

	if tgtPortWWNs == nil {
		msg := fmt.Sprintf("Get fc target wwpns error, tgt_port_wwns:%s, init_targ_map:%s", tgtPortWWNs, initTargMap)
		log.Errorf(msg)
		return nil, nil, errors.New(msg)
	}

	return tgtPortWWNs, initTargMap, nil
}

func (c *DoradoClient) isInStringArray(s string, source []string) bool {
	for _, i := range source {
		if s == i {
			return true
		}
	}
	return false
}

func (c *DoradoClient) getFCTargetWWPNs(wwn string) ([]string, error) {
	resp := &FCTargWWPNResp{}
	url := fmt.Sprintf("/host_link?INITIATOR_TYPE=223&INITIATOR_PORT_WWN=%s", wwn)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get fc target wwpn failed.")
		return nil, err
	}

	var fcWWPNs []string
	if len(resp.Data) > 0 {
		for _, item := range resp.Data {
			if wwn == item.IniPortWWN {
				fcWWPNs = append(fcWWPNs, item.TargPortWWN)
			}
		}
	}

	return fcWWPNs, nil
}

func (c *DoradoClient) getObjCountFromLungroupByType(lunGroupId, lunType string) (int, error) {
	// Get obj count associated to the lungroup.
	var cmdType string
	if lunType == ObjectTypeLun {
		cmdType = "lun"
	} else {
		cmdType = "snapshot"
	}

	resp := &ObjCountResp{}
	url := fmt.Sprintf("/%s/count?TYPE=%s&ASSOCIATEOBJTYPE=256&ASSOCIATEOBJID=%s", cmdType, lunType, lunGroupId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get Obj count from lungroup by type failed.")
		return 0, err
	}
	if resp.Error.Code == ErrorObjectUnavailable {
		log.Errorf("LUN group %s not exist.", lunGroupId)
		return 0, nil
	}
	count, _ := strconv.Atoi(resp.Data.Count)
	return count, nil
}

func (c *DoradoClient) getObjectCountFromLungroup(lunGrpId string) (int, error) {
	lunCount, err := c.getObjCountFromLungroupByType(lunGrpId, ObjectTypeLun)
	if err != nil {
		return 0, nil
	}
	snapshotCount, err := c.getObjCountFromLungroupByType(lunGrpId, ObjectTypeSnapshot)
	if err != nil {
		return 0, nil
	}
	return lunCount + snapshotCount, nil
}

func (c *DoradoClient) getHostGroupNumFromHost(hostId string) (int, error) {
	resp := &ObjCountResp{}
	url := fmt.Sprintf("/hostgroup/count?TYPE=14&ASSOCIATEOBJTYPE=21&ASSOCIATEOBJID=%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get hostgroup num from host failed.")
		return 0, err
	}

	count, _ := strconv.Atoi(resp.Data.Count)
	return count, nil
}

func (c *DoradoClient) removeFCFromHost(wwn string) error {
	data := map[string]interface{}{
		"TYPE": ObjectTypeFcInitiator,
		"ID":   wwn,
	}

	err := c.request("PUT", "/fc_initiator/remove_fc_from_host", data, nil)
	return err
}

func (c *DoradoClient) getHostgroupAssociatedViews(hostGrpId string) ([]MappingView, error) {
	resp := &MappingViewsResp{}
	url := fmt.Sprintf("/mappingview/associate?TYPE=245&ASSOCIATEOBJTYPE=14&ASSOCIATEOBJID=%s", hostGrpId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get hostgroup associated views failed.")
		return nil, err
	}

	return resp.Data, nil
}

func (c *DoradoClient) getHostsInHostgroup(hostGrpId string) ([]Host, error) {
	resp := &HostsResp{}
	url := fmt.Sprintf("/host/associate?ASSOCIATEOBJTYPE=14&ASSOCIATEOBJID=%s", hostGrpId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get hosts in hostgroup failed.")
		return nil, err
	}

	var tempHosts = []Host{}
	if resp.Error.Code == ErrorHostGroupNotExist {
		log.Errorf("Host group %s not exist", hostGrpId)
		return tempHosts, nil
	}

	return resp.Data, nil
}

func (c *DoradoClient) checkFCInitiatorsExistInHost(hostId string) (bool, error) {
	resp := &FCInitiatorsResp{}
	url := fmt.Sprintf("/fc_initiator?range=[0-65535]&PARENTID=%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get FC initiators exist in host failed.")
		return false, err
	}
	if len(resp.Data) > 0 {
		return true, nil
	}

	return false, nil
}

func (c *DoradoClient) checkIscsiInitiatorsExistInHost(hostId string) (bool, error) {
	resp := &FCInitiatorsResp{}
	url := fmt.Sprintf("/iscsi_initiator?range=[0-65535]&PARENTID=%s", hostId)
	if err := c.request("GET", url, nil, resp); err != nil {
		log.Errorf("Get FC initiators exist in host failed.")
		return false, err
	}
	if len(resp.Data) > 0 {
		return true, nil
	}

	return false, nil
}

func (c *DoradoClient) ListControllers() ([]SimpleStruct, error) {
	resp := &SimpleResp{}
	err := c.request("GET", "/controller", nil, resp)
	return resp.Data, err
}

func (c *DoradoClient) GetPerformance(resId string, dataIdList []string) (map[string]string, error) {
	perf := &PerformancesResp{}
	url := fmt.Sprintf("/performace_statistic/cur_statistic_data/"+
		"?CMO_STATISTIC_UUID=%s&CMO_STATISTIC_DATA_ID_LIST=%s", resId, strings.Join(dataIdList, ","))
	if err := c.request("GET", url, nil, perf); err != nil {
		return nil, err
	}
	if len(perf.Data) == 0 {
		return nil, fmt.Errorf("got empty performance data")
	}

	dataList := strings.Split(perf.Data[0].DataList, ",")
	idList := strings.Split(perf.Data[0].DataIdList, ",")
	if len(dataList) != len(idList) {
		return nil, fmt.Errorf("wrong response data, data id list length does not equal to data length")
	}

	perfMap := map[string]string{}
	for i, id := range idList {
		perfMap[id] = dataList[i]
	}

	return perfMap, nil
}
