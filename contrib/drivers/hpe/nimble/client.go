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

package nimble

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/astaxie/beego/httplib"
	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

const (
	requestRetryTimes            = 1
	FcInitiatorDefaultGrpName    = "OsdsFcGrp"
	IscsiInitiatorDefaultGrpName = "OsdsIscsiGrp"

	tokenUrlPath         = "/tokens"
	poolUrlPath          = "/pools"
	poolDetailUrlPath    = "/pools/detail"
	volumeUrlPath        = "/volumes"
	snapshotUrlPath      = "/snapshots"
	initiatorUrlPath     = "/initiators"
	initiatorGrpUrlPath  = "/initiator_groups"
	accessControlUrlPath = "/access_control_records"
)

// *****************Original Errors***********************
func (e *ClinetErrors) Error() string {
	var errStrings []string
	for _, err := range e.Errs {
		errStrings = append(errStrings, fmt.Sprint(err))
	}
	return fmt.Sprint(strings.Join(errStrings, "\n"))
}

func (e *ArrayInnerErrorBody) Error() string {
	var errStrings []string
	for _, err := range e.Errs {
		errStrings = append(errStrings, fmt.Sprint(err.Error()))
	}
	return fmt.Sprint(strings.Join(errStrings, "\n"))
}

func (e *ArrayInnerErrorResp) Error() string {
	return fmt.Sprintf("code:%v severity:%v text:%v", e.Code, e.Severity, e.Text)
}

// *******************************************************

func unset(strings []string, search string) []string {
	result := []string{}
	deleteFlag := false
	for _, v := range strings {
		if v == search && !deleteFlag {
			deleteFlag = true
			continue
		}
		result = append(result, v)
	}
	return result
}

func NewClient(opt *AuthOptions) (*NimbleClient, error) {
	edp := strings.Split(opt.Endpoints, ",")
	c := &NimbleClient{
		user:      opt.Username,
		passwd:    opt.Password,
		endpoints: edp,
		insecure:  opt.Insecure,
	}
	err := c.login()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *NimbleClient) login() error {
	reqBody := &LoginReqBody{
		Data: LoginReqData{
			Username: c.user,
			Password: c.passwd,
		},
	}
	var tokens []string
	var errs []error
	edp := c.endpoints

	for _, ep := range edp {
		url := ep + tokenUrlPath
		auth := &AuthRespBody{}
		token := ""
		log.Infof("%v: trying login to %v", DriverName, ep)
		b, _, err := c.doRequest("POST", url, reqBody, token)

		// Basic HTTP Request Error
		if err != nil {
			log.Errorf("%v: login failed.", DriverName)
			c.endpoints = unset(c.endpoints, ep) // Delete invalid endpoint from client
			errs = append(errs, err)
			continue
		}
		json.Unmarshal(b, auth)
		tokens = append(tokens, auth.Data.SessionToken)
		log.Infof("%v: got token from %v", DriverName, url)
	}

	c.tokens = tokens // Insert valid tokes into client

	if len(errs) != 0 {
		err := &ClinetErrors{errs}
		return err
	}
	return nil
}

func (c *NimbleClient) doRequest(method, url string, in interface{}, token string) ([]byte, http.Header, error) {
	req := httplib.NewBeegoRequest(url, method)
	req.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: c.insecure,
	})
	req.Header("Connection", "keep-alive")
	req.Header("Content-Type", "application/json;charset=utf-8")
	req.Header("X-Auth-Token", token) // Set token

	if in != nil {
		body, _ := json.Marshal(in)
		req.Body(body)

	}

	resp, err := req.Response()
	if err != nil {
		log.Errorf("%v: http request failed, method: %s url: %s error: %v", DriverName, method, url, err)
		return nil, nil, err
	}

	b, err := req.Bytes()
	if err != nil {
		log.Errorf("%v: get byte[] from response failed, method: %s url: %s error: %v", DriverName, method, url, err)
		return nil, nil, err
	}

	inErr := &ArrayInnerErrorBody{}
	json.Unmarshal(b, inErr)
	if len(inErr.Errs) != 0 {
		log.Errorf("%v: get error Infof from response failed, method: %s url: %s error: %v", DriverName, method, url, inErr)
		return nil, nil, inErr
	}

	return b, resp.Header, nil
}

func (c *NimbleClient) ListStoragePools() ([]StoragePoolRespData, error) {
	resp := &StoragePoolsRespBody{}
	pools := []StoragePoolRespData{}
	var err error
	var errs []error

	if len(c.endpoints) == 0 {
		log.Errorf("%v: there are no valid endpoints.", DriverName)
		return nil, fmt.Errorf("%v: cannot get storage pools\n", DriverName)
	}
	for i, ep := range c.endpoints {
		err = c.request("GET", ep+poolDetailUrlPath, nil, resp, c.tokens[i])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		// Set endpoint which belonging to the storage pool
		for j, _ := range resp.Data {
			resp.Data[j].Endpoint = ep
			resp.Data[j].Token = c.tokens[i]
		}
		pools = append(pools, resp.Data...)
	}

	if len(errs) != 0 {
		err = &ClinetErrors{errs}
	}
	return pools, err
}

func (c *NimbleClient) request(method, url string, in, out interface{}, token string) error {
	var b []byte
	var errReq error
	var errs []error
	for i := 0; i < requestRetryTimes; i++ {
		b, _, errReq = c.doRequest(method, url, in, token)
		if errReq == nil {
			json.Unmarshal(b, out)
			log.Infof("%v: got response from %v.", DriverName, url)
			break
		} else {

			log.Errorf("%v: url:%s %s body:%+v", DriverName, method, url, in)

			// Token expired handling
			if inErr, ok := errReq.(*ArrayInnerErrorBody); ok {
				for j := range inErr.Errs {
					if inErr.Errs[j].Code == ErrorUnauthorizedToServer {
						log.Errorf("%v: auth failure, trying re-login....", DriverName)
						if errLogin := c.login(); errLogin == nil {
							log.Infof("%v: relogin success!!", DriverName)
							break
						} else {
							log.Errorf("%v: relogin failed.", DriverName)
						}
					}
				}
			}
		}

		if i == requestRetryTimes-1 {
			log.Errorf("%v: finally, could not get response from %v.", DriverName, url)
			errs = append(errs, errReq)
		}
	}

	if len(errs) != 0 {
		err := &ClinetErrors{errs}
		return err
	}
	return nil
}

func (c *NimbleClient) GetPoolIdByName(poolName string) (string, error) {
	pools, err := c.ListStoragePools()
	if err != nil {
		return "", err
	}
	for _, p := range pools {
		if p.Name == poolName {
			return p.Id, nil
		}
	}
	return "", fmt.Errorf("%v: not found specified pool '%s'\n", DriverName, poolName)
}

func (c *NimbleClient) GetTokenByPoolId(poolId string) (string, string, error) {
	pools, err := c.ListStoragePools()
	if err != nil {
		return "", "", err
	}
	for _, p := range pools {
		if p.Id == poolId {
			return p.Endpoint, p.Token, nil
		}
	}
	return "", "", fmt.Errorf("%v: not found specified pool '%s'\n", DriverName, poolId)
}

func (c *NimbleClient) CreateVolume(poolId string, opt *pb.CreateVolumeOpts) (*VolumeRespData, error) {

	/* Get endpoint and token for spwcific pool */
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return nil, err
	}

	// Parse options from Profile----------------------------------------------
	profileOpt := &model.ProfileSpec{}
	reqOptions := CreateVolumeReqData{}
	if err := json.Unmarshal([]byte(opt.GetProfile()), profileOpt); err != nil {
		return nil, err
	}

	options, err := json.Marshal(profileOpt.CustomProperties)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(options, &reqOptions); err != nil {
		return nil, err
	}
	// -------------------------------------------------------------------------

	reqOptions.Name = opt.GetId()
	reqOptions.PoolId = poolId
	reqOptions.Size = Gib2Mebi(opt.GetSize())
	reqOptions.Description = TruncateDescription(opt.GetDescription())

	// Create volume from snapshot
	if opt.GetSnapshotId() != "" {
		log.Infof("%v: try to create volume from snapshot...", DriverName)
		log.Infof("%v: snap id: %v", DriverName, opt.GetSnapshotId())
		storageSnapshotId, err := c.GetStorageSnapshotId(poolId, opt.GetSnapshotId())
		if err != nil {
			return nil, err
		}
		if storageSnapshotId == "" {
			err = fmt.Errorf("%v: there is no such snapshot name on storage => %v\n", DriverName, opt.GetSnapshotId())
			return nil, err
		}

		reqOptions.BaseSnapId = storageSnapshotId
		reqOptions.Clone = true
	}

	lunResp := &VolumeRespBody{}
	reqBody := &CreateVolumeReqBody{Data: reqOptions}
	err = c.request("POST", ep+volumeUrlPath, reqBody, lunResp, token)
	return &lunResp.Data, err
}

func (c *NimbleClient) ListVolume(poolId string) (*AllVolumeRespBody, error) {
	/* Get endpoint and token for spwcific pool */
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return nil, err
	}
	respBody := &AllVolumeRespBody{}
	err = c.request("GET", ep+volumeUrlPath+"/detail", nil, respBody, token)
	return respBody, err
}

func (c *NimbleClient) GetStorageVolumeId(poolId string, volName string) (string, error) {
	storageVolumeId := ""

	respBody, err := c.ListVolume(poolId)
	if err != nil {
		return storageVolumeId, err
	}

	for _, data := range respBody.Data {
		if data.Name == volName {
			storageVolumeId = data.Id
			return storageVolumeId, nil
		}
	}
	return storageVolumeId, fmt.Errorf("%v: could not get storage volume ID of %v\n", DriverName, volName)
}

func (c *NimbleClient) GetStorageSnapshotId(poolId string, baseSnapName string) (string, error) {
	storageSnapshotId := ""

	// Get all volume names
	volResp, err := c.ListVolume(poolId)
	if err != nil {
		return storageSnapshotId, err
	}

	for _, volData := range volResp.Data {
		ep, token, err := c.GetTokenByPoolId(poolId)
		if err != nil {
			return storageSnapshotId, err
		}

		respBody := &AllSnapshotRespBody{}
		err = c.request("GET", ep+snapshotUrlPath+"?vol_name="+volData.Name, nil, respBody, token)
		if err != nil {
			return storageSnapshotId, err
		}

		for _, snapData := range respBody.Data {
			if baseSnapName == snapData.Name {
				storageSnapshotId = snapData.Id
				return storageSnapshotId, nil
			}
		}

	}

	return storageSnapshotId, err
}

func (c *NimbleClient) DeleteVolume(poolId string, opt *pb.DeleteVolumeOpts) error {
	lunId := opt.GetMetadata()["LunId"]
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return err
	}

	err = c.OfflineVolume(lunId, poolId)
	if err != nil {
		return err
	}

	// Nimble storage daes not support delete options. No need to parse values in profile.

	err = c.request("DELETE", ep+volumeUrlPath+"/"+lunId, nil, nil, token)
	return err
}

func (c *NimbleClient) OfflineVolume(id string, poolId string) error {
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return err
	}

	reqOptions := OfflineVolumeReqData{
		Online: false,
	}

	reqBody := &OfflineVolumeReqBody{Data: reqOptions}
	err = c.request("PUT", ep+volumeUrlPath+"/"+id, reqBody, nil, token)
	return err
}

func (c *NimbleClient) ExtendVolume(poolId string, opt *pb.ExtendVolumeOpts) (*VolumeRespData, error) {
	lunId := opt.GetMetadata()["LunId"]
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return nil, err
	}

	// Parse options from Profile----------------------------------------------
	reqOptions := ExtendVolumeReqData{}

	profileOpt := &model.ProfileSpec{}
	if err := json.Unmarshal([]byte(opt.GetProfile()), profileOpt); err != nil {
		return nil, err
	}

	options, err := json.Marshal(profileOpt.CustomProperties)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(options, &reqOptions); err != nil {
		return nil, err
	}

	// -------------------------------------------------------------------------
	reqOptions.Size = Gib2Mebi(opt.GetSize())
	reqBody := &ExtendVolumeReqBody{Data: reqOptions}
	respBody := &ExtendVolumeRespBody{}
	err = c.request("PUT", ep+volumeUrlPath+"/"+lunId, reqBody, respBody, token)
	return &respBody.Data, err
}

func (c *NimbleClient) CreateSnapshot(poolId string, opt *pb.CreateVolumeSnapshotOpts) (*SnapshotRespData, error) {
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return nil, err
	}

	// Parse options from Profile----------------------------------------------
	reqOptions := CreateSnapshotReqData{}

	profileOpt := &model.ProfileSpec{}
	if err := json.Unmarshal([]byte(opt.GetProfile()), profileOpt); err != nil {
		return nil, err
	}

	options, err := json.Marshal(profileOpt.CustomProperties)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(options, &reqOptions); err != nil {
		return nil, err
	}
	// -------------------------------------------------------------------------

	reqOptions.Name = opt.GetId()
	reqOptions.VolId = opt.GetMetadata()["LunId"]
	reqOptions.Description = TruncateDescription(opt.GetDescription())

	reqBody := &CreateSnapshotReqBody{Data: reqOptions}
	respBody := &SnapshotRespBody{}
	err = c.request("POST", ep+snapshotUrlPath, reqBody, respBody, token)
	return &respBody.Data, err
}

func (c *NimbleClient) DeleteSnapshot(poolId string, opt *pb.DeleteVolumeSnapshotOpts) error {
	snapId := opt.GetMetadata()["SnapId"]
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return err
	}
	// Nimble storage daes not support delete options. No need to parse values in profile.
	return c.request("DELETE", ep+snapshotUrlPath+"/"+snapId, nil, nil, token)
}

func (c *NimbleClient) ListInitiator(poolId string, initiatorResp *AllInitiatorRespBody) error {
	// Get endpoint and token for spwcific pool
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return err
	}

	err = c.request("GET", ep+initiatorUrlPath+"/detail", nil, initiatorResp, token)
	return err
}

func (c *NimbleClient) GetStorageInitiatorGrpId(poolId string, initiatorIqn string) (string, error) {
	// List all registered initiators
	respBody := &AllInitiatorRespBody{}
	err := c.ListInitiator(poolId, respBody)
	if err != nil {
		return "", err
	}
	for _, data := range respBody.Data {
		if data.Iqn == initiatorIqn {
			return data.InitiatorGroupId, nil
		}
	}

	return "", nil
}

func (c *NimbleClient) RegisterInitiatorIntoDefaultGrp(poolId string, opt *pb.CreateVolumeAttachmentOpts, initiatorGrpId string) (*InitiatorRespData, error) {
	// Get endpoint and token for spwcific pool
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return nil, err
	}
	reqOptions := CreateInitiatorReqData{}
	reqOptions.IpAddress = opt.GetHostInfo().Ip
	reqOptions.InitiatorGroupId = initiatorGrpId
	// For iSCSI initiator
	if opt.GetAccessProtocol() == ISCSIProtocol {
		reqOptions.AccessProtocol = "iscsi"
		reqOptions.Iqn = opt.GetHostInfo().Initiator
		reqOptions.Label = opt.GetId()
	}
	// For FC initiator
	if opt.GetAccessProtocol() == FCProtocol {
		reqOptions.AccessProtocol = "fc"
		reqOptions.Wwpn = opt.GetHostInfo().Initiator
		reqOptions.Alias = opt.GetId()
	}

	reqBody := &CreateInitiatorReqBody{Data: reqOptions}
	respBody := &InitiatorRespBody{}
	err = c.request("POST", ep+initiatorUrlPath, reqBody, respBody, token)
	return &respBody.Data, err
}

func (c *NimbleClient) GetDefaultInitiatorGrpId(poolId string, opt *pb.CreateVolumeAttachmentOpts) (string, error) {

	respBody := &AllInitiatorGrpRespBody{}
	err := c.ListInitiatorGrp(poolId, respBody)
	if err != nil {
		return "", err
	}

	// For iSCSI
	if opt.GetAccessProtocol() == ISCSIProtocol {
		log.Infof("%v: trying to get default iscsi initiator group ID: %v.", DriverName, IscsiInitiatorDefaultGrpName)
		for _, data := range respBody.Data {
			if data.Name == IscsiInitiatorDefaultGrpName {
				return data.Id, nil
			}
		}
	}

	// For FC
	if opt.GetAccessProtocol() == FCProtocol {
		log.Infof("%v: trying to get default fc initiator group ID: %v.", DriverName, FcInitiatorDefaultGrpName)
		for _, data := range respBody.Data {
			if data.Name == FcInitiatorDefaultGrpName {
				return data.Id, nil
			}
		}
	}
	return "", nil
}

func (c *NimbleClient) ListInitiatorGrp(poolId string, initiatorGrpResp *AllInitiatorGrpRespBody) error {
	// Get endpoint and token for spwcific pool
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return err
	}

	err = c.request("GET", ep+initiatorGrpUrlPath+"/detail", nil, initiatorGrpResp, token)
	return err
}

func (c *NimbleClient) CreateInitiatorDefaultGrp(poolId string, opt *pb.CreateVolumeAttachmentOpts) (*InitiatorGrpRespData, error) {
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return nil, err
	}
	reqOptions := CreateInitiatorGrpReqData{}
	// For iSCSI initiator
	if opt.GetAccessProtocol() == ISCSIProtocol {
		reqOptions.Name = IscsiInitiatorDefaultGrpName
		reqOptions.AccessProtocol = "iscsi"
		reqOptions.Description = "OpenSDS default iSCSI group"
		reqOptions.TargetSubnets = append(reqOptions.TargetSubnets, map[string]string{"label": "management"})
	}
	// For FC initiator
	if opt.GetAccessProtocol() == FCProtocol {
		reqOptions.Name = FcInitiatorDefaultGrpName
		reqOptions.AccessProtocol = "fc"
		reqOptions.Description = "OpenSDS default FC group"
	}

	reqBody := &CreateInitiatorGrpReqBody{Data: reqOptions}
	respBody := &InitiatorGrpRespBody{}
	err = c.request("POST", ep+initiatorGrpUrlPath, reqBody, respBody, token)
	return &respBody.Data, err
}

func (c *NimbleClient) AttachVolume(poolId string, volName string, initiatorGrpId string) (*AccessControlRespData, error) {
	storageVolumeId, err := c.GetStorageVolumeId(poolId, volName)
	if err != nil {
		return nil, err
	}

	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return nil, err
	}
	reqOptions := CreateAccessControlReqData{}
	reqOptions.VolId = storageVolumeId
	reqOptions.InitiatorGroupId = initiatorGrpId

	reqBody := &CreateAccessControlReqBody{Data: reqOptions}
	respBody := &AccessControlRespBody{}
	err = c.request("POST", ep+accessControlUrlPath, reqBody, respBody, token)

	return &respBody.Data, err
}

func (c *NimbleClient) DetachVolume(poolId string, storageAceessId string) error {
	// Detach request does not give us response body
	ep, token, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return err
	}

	err = c.request("DELETE", ep+accessControlUrlPath+"/"+storageAceessId, nil, nil, token)
	return err
}

func (c *NimbleClient) GetTargetVolumeInfo(poolId string, volName string) (string, string, error) {
	ep, _, err := c.GetTokenByPoolId(poolId)
	if err != nil {
		return "", "", err
	}
	respBody, err := c.ListVolume(poolId)
	for _, data := range respBody.Data {
		if data.Name == volName {
			return data.TargetName, ep, err
		}
	}

	return "", "", fmt.Errorf("%v: couldnot get volume target.\n", DriverName)
}
