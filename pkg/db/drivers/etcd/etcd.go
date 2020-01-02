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

/*
This module implements the etcd database operation of data structure
defined in api module.

*/

package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/opensds/opensds/pkg/utils/config"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
	"github.com/opensds/opensds/pkg/utils/urls"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultSortKey          = "ID"
	defaultBlockProfileName = "default_block"
	defaultFileProfileName  = "default_file"
	typeBlock               = "block"
	typeFile                = "file"
)

var validKey = []string{"limit", "offset", "sortDir", "sortKey"}

const (
	typeFileShares         string = "FileShares"
	typeFileShareSnapshots string = "FileShareSnapshots"
	typeDocks              string = "Docks"
	typePools              string = "Pools"
	typeProfiles           string = "Profiles"
	typeVolumes            string = "Volumes"
	typeAttachments        string = "Attachments"
	typeVolumeSnapshots    string = "VolumeSnapshots"
)

var sortableKeysMap = map[string][]string{
	typeFileShares:         {"ID", "NAME", "STATUS", "AVAILABILITYZONE", "PROFILEID", "TENANTID", "SIZE", "POOLID", "DESCRIPTION"},
	typeFileShareSnapshots: {"ID", "NAME", "VOLUMEID", "STATUS", "USERID", "TENANTID", "SIZE"},
	typeDocks:              {"ID", "NAME", "STATUS", "ENDPOINT", "DRIVERNAME", "DESCRIPTION"},
	typePools:              {"ID", "NAME", "STATUS", "AVAILABILITYZONE", "DOCKID"},
	typeProfiles:           {"ID", "NAME", "DESCRIPTION"},
	typeVolumes:            {"ID", "NAME", "STATUS", "AVAILABILITYZONE", "PROFILEID", "TENANTID", "SIZE", "POOLID", "DESCRIPTION", "GROUPID"},
	typeAttachments:        {"ID", "VOLUMEID", "STATUS", "USERID", "TENANTID", "SIZE"},
	typeVolumeSnapshots:    {"ID", "NAME", "VOLUMEID", "STATUS", "USERID", "TENANTID", "SIZE"},
}

func IsAdminContext(ctx *c.Context) bool {
	return ctx.IsAdmin
}

func AuthorizeProjectContext(ctx *c.Context, tenantId string) bool {
	return ctx.TenantId == tenantId
}

// NewClient
func NewClient(etcd *config.Database) *Client {
	return &Client{
		clientInterface: Init(etcd),
	}
}

// Client
type Client struct {
	clientInterface
}

//Parameter
type Parameter struct {
	beginIdx, endIdx int
	sortDir, sortKey string
}

//IsInArray
func (c *Client) IsInArray(e string, s []string) bool {
	for _, v := range s {
		if strings.EqualFold(e, v) {
			return true
		}
	}
	return false
}

func (c *Client) SelectOrNot(m map[string][]string) bool {
	for key := range m {
		if !utils.Contained(key, validKey) {
			return true
		}
	}
	return false
}

//Get parameter limit
func (c *Client) GetLimit(m map[string][]string) int {
	var limit int
	var err error
	v, ok := m["limit"]
	if ok {
		limit, err = strconv.Atoi(v[0])
		if err != nil || limit < 0 {
			log.Warning("Invalid input limit:", limit, ",use default value instead:50")
			return constants.DefaultLimit
		}
	} else {
		log.Warning("The parameter limit is not present,use default value instead:50")
		return constants.DefaultLimit
	}
	return limit
}

//Get parameter offset
func (c *Client) GetOffset(m map[string][]string, size int) int {

	var offset int
	var err error
	v, ok := m["offset"]
	if ok {
		offset, err = strconv.Atoi(v[0])

		if err != nil || offset < 0 || offset > size {
			log.Warning("Invalid input offset or input offset is out of bounds:", offset, ",use default value instead:0")

			return constants.DefaultOffset
		}

	} else {
		log.Warning("The parameter offset is not present,use default value instead:0")
		return constants.DefaultOffset
	}
	return offset
}

//Get parameter sortDir
func (c *Client) GetSortDir(m map[string][]string) string {
	var sortDir string
	v, ok := m["sortDir"]
	if ok {
		sortDir = v[0]
		if !strings.EqualFold(sortDir, "desc") && !strings.EqualFold(sortDir, "asc") {
			log.Warning("Invalid input sortDir:", sortDir, ",use default value instead:desc")
			return constants.DefaultSortDir
		}
	} else {
		log.Warning("The parameter sortDir is not present,use default value instead:desc")
		return constants.DefaultSortDir
	}
	return sortDir
}

//Get parameter sortKey
func (c *Client) GetSortKey(m map[string][]string, sortKeys []string) string {
	var sortKey string
	v, ok := m["sortKey"]
	if ok {
		sortKey = strings.ToUpper(v[0])
		if !c.IsInArray(sortKey, sortKeys) {
			log.Warning("Invalid input sortKey:", sortKey, ",use default value instead:ID")
			return defaultSortKey
		}

	} else {
		log.Warning("The parameter sortKey is not present,use default value instead:ID")
		return defaultSortKey
	}
	return sortKey
}

func (c *Client) FilterAndSort(src interface{}, params map[string][]string, sortableKeys []string) interface{} {
	var ret interface{}
	ret = utils.Filter(src, params)
	if len(params["sortKey"]) > 0 && utils.ContainsIgnoreCase(sortableKeys, params["sortKey"][0]) {
		ret = utils.Sort(ret, params["sortKey"][0], c.GetSortDir(params))
	}
	ret = utils.Slice(ret, c.GetOffset(params, reflect.ValueOf(src).Len()), c.GetLimit(params))
	return ret
}

//ParameterFilter
func (c *Client) ParameterFilter(m map[string][]string, size int, sortKeys []string) *Parameter {

	limit := c.GetLimit(m)
	offset := c.GetOffset(m, size)

	beginIdx := offset
	endIdx := limit + offset

	// If use not specified the limit return all the items.
	if limit == constants.DefaultLimit || endIdx > size {
		endIdx = size
	}

	sortDir := c.GetSortDir(m)
	sortKey := c.GetSortKey(m, sortKeys)

	return &Parameter{beginIdx, endIdx, sortDir, sortKey}
}

// *************   FileShare code block  *************

var fileshare_sortKey string

type FileShareSlice []*model.FileShareSpec

func (fileshare FileShareSlice) Len() int { return len(fileshare) }

func (fileshare FileShareSlice) Swap(i, j int) {
	fileshare[i], fileshare[j] = fileshare[j], fileshare[i]
}

func (fileshare FileShareSlice) Less(i, j int) bool {
	switch fileshare_sortKey {

	case "ID":
		return fileshare[i].Id < fileshare[j].Id
	case "NAME":
		return fileshare[i].Name < fileshare[j].Name
	case "STATUS":
		return fileshare[i].Status < fileshare[j].Status
	case "AVAILABILITYZONE":
		return fileshare[i].AvailabilityZone < fileshare[j].AvailabilityZone
	case "PROFILEID":
		return fileshare[i].ProfileId < fileshare[j].ProfileId
	case "TENANTID":
		return fileshare[i].TenantId < fileshare[j].TenantId
	case "SIZE":
		return fileshare[i].Size < fileshare[j].Size
	case "POOLID":
		return fileshare[i].PoolId < fileshare[j].PoolId
	case "DESCRIPTION":
		return fileshare[i].Description < fileshare[j].Description
	}
	return false
}

func (c *Client) FindFileShareValue(k string, p *model.FileShareSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAt":
		return p.UpdatedAt
	case "TenantId":
		return p.TenantId
	case "UserId":
		return p.UserId
	case "Name":
		return p.Name
	case "Description":
		return p.Description
	case "AvailabilityZone":
		return p.AvailabilityZone
	case "Size":
		return strconv.FormatInt(p.Size, 10)
	case "Status":
		return p.Status
	case "PoolId":
		return p.PoolId
	case "ProfileId":
		return p.ProfileId
	}
	return ""
}

func (c *Client) CreateFileShareAcl(ctx *c.Context, fshare *model.FileShareAclSpec) (*model.FileShareAclSpec, error) {
	acls, err := c.ListFileShareAclsByShareId(ctx, fshare.FileShareId)
	if err != nil {
		log.Error("failed to list acls")
		return nil, err
	}

	for _, acl := range acls {
		if acl.AccessTo == fshare.AccessTo {
			errstr := "for fileshareID: " + acl.FileShareId + ", acl is already set with ip: " + acl.AccessTo + ". If you want to set new acl, first delete the existing one"
			log.Error(errstr)
			return nil, fmt.Errorf(errstr)
		}
	}

	fshare.TenantId = ctx.TenantId
	fshareBody, err := json.Marshal(fshare)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateFileShareAclURL(urls.Etcd, ctx.TenantId, fshare.Id),
		Content: string(fshareBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when create fileshare access rules in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return fshare, nil
}

// UpdateFileShareAcl
func (c *Client) UpdateFileShareAcl(ctx *c.Context, acl *model.FileShareAclSpec) (*model.FileShareAclSpec, error) {
	result, err := c.GetFileShareAcl(ctx, acl.Id)
	if err != nil {
		return nil, err
	}

	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)
	result.Metadata = acl.Metadata

	jsonBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateFileShareAclURL(urls.Etcd, result.TenantId, acl.Id),
		NewContent: string(jsonBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when update fileshare acl in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

func (c *Client) CreateFileShare(ctx *c.Context, fshare *model.FileShareSpec) (*model.FileShareSpec, error) {
	fshare.TenantId = ctx.TenantId
	fshareBody, err := json.Marshal(fshare)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateFileShareURL(urls.Etcd, ctx.TenantId, fshare.Id),
		Content: string(fshareBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when create fileshare in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return fshare, nil
}

func (c *Client) SortFileShares(shares []*model.FileShareSpec, p *Parameter) []*model.FileShareSpec {
	volume_sortKey = p.sortKey

	if strings.EqualFold(p.sortDir, "dsc") {
		sort.Sort(FileShareSlice(shares))

	} else {
		sort.Sort(sort.Reverse(FileShareSlice(shares)))
	}
	return shares
}

func (c *Client) ListFileSharesAclWithFilter(ctx *c.Context, m map[string][]string) ([]*model.FileShareAclSpec, error) {
	fileshares, err := c.ListFileSharesAcl(ctx)
	if err != nil {
		log.Error("list fileshare failed: ", err)
		return nil, err
	}
	return fileshares, nil
}

func (c *Client) ListFileSharesAcl(ctx *c.Context) ([]*model.FileShareAclSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateFileShareAclURL(urls.Etcd, ctx.TenantId),
	}

	// Admin user should get all fileshares including the fileshares whose tenant is not admin.
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateFileShareAclURL(urls.Etcd, "")
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when list fileshares in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var fileshares = []*model.FileShareAclSpec{}
	if len(dbRes.Message) == 0 {
		return fileshares, nil
	}
	for _, msg := range dbRes.Message {
		var share = &model.FileShareAclSpec{}
		if err := json.Unmarshal([]byte(msg), share); err != nil {
			log.Error("when parsing fileshare in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		fileshares = append(fileshares, share)
	}
	return fileshares, nil
}

func (c *Client) ListFileShareAclsByShareId(ctx *c.Context, fileshareId string) ([]*model.FileShareAclSpec, error) {
	acls, err := c.ListFileSharesAcl(ctx)
	if err != nil {
		return nil, err
	}

	var aclList []*model.FileShareAclSpec
	for _, acl := range acls {
		if acl.FileShareId == fileshareId {
			aclList = append(aclList, acl)
		}
	}
	return aclList, nil
}

func (c *Client) ListSnapshotsByShareId(ctx *c.Context, fileshareId string) ([]*model.FileShareSnapshotSpec, error) {
	snaps, err := c.ListFileShareSnapshots(ctx)
	if err != nil {
		return nil, err
	}

	var snapList []*model.FileShareSnapshotSpec
	for _, snap := range snaps {
		if snap.FileShareId == fileshareId {
			snapList = append(snapList, snap)
		}
	}
	return snapList, nil
}

func (c *Client) ListFileSharesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.FileShareSpec, error) {
	fileshares, err := c.ListFileShares(ctx)
	if err != nil {
		log.Error("list fileshare failed: ", err)
		return nil, err
	}

	tmpFileshares := c.FilterAndSort(fileshares, m, sortableKeysMap[typeFileShares])
	var res = []*model.FileShareSpec{}
	for _, data := range tmpFileshares.([]interface{}) {
		res = append(res, data.(*model.FileShareSpec))
	}
	return res, nil
}

// ListFileShares
func (c *Client) ListFileShares(ctx *c.Context) ([]*model.FileShareSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateFileShareURL(urls.Etcd, ctx.TenantId),
	}

	// Admin user should get all fileshares including the fileshares whose tenant is not admin.
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateFileShareURL(urls.Etcd, "")
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when list fileshares in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var fileshares = []*model.FileShareSpec{}
	if len(dbRes.Message) == 0 {
		return fileshares, nil
	}
	for _, msg := range dbRes.Message {
		var share = &model.FileShareSpec{}
		if err := json.Unmarshal([]byte(msg), share); err != nil {
			log.Error("when parsing fileshare in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		fileshares = append(fileshares, share)
	}
	return fileshares, nil
}

// ListFileSharesByProfileId
func (c *Client) ListFileSharesByProfileId(ctx *c.Context, prfId string) ([]string, error) {
	fileshares, err := c.ListFileShares(ctx)
	if err != nil {
		return nil, err
	}
	var res_fileshares []string
	for _, shares := range fileshares {
		if shares.ProfileId == prfId {
			res_fileshares = append(res_fileshares, shares.Name)
		}
	}
	return res_fileshares, nil
}

// GetFileShareAcl
func (c *Client) GetFileShareAcl(ctx *c.Context, aclID string) (*model.FileShareAclSpec, error) {
	acl, err := c.getFileShareAcl(ctx, aclID)
	if !IsAdminContext(ctx) || err == nil {
		return acl, err
	}
	acls, err := c.ListFileSharesAcl(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range acls {
		if f.Id == aclID {
			return f, nil
		}
	}
	return nil, fmt.Errorf("specified fileshare acl(%s) can't find", aclID)
}

func (c *Client) getFileShareAcl(ctx *c.Context, aclID string) (*model.FileShareAclSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateFileShareAclURL(urls.Etcd, ctx.TenantId, aclID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when get fileshare acl in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var acl = &model.FileShareAclSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), acl); err != nil {
		log.Error("when parsing fileshare acl in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return acl, nil
}

// GetFileShare
func (c *Client) GetFileShare(ctx *c.Context, fshareID string) (*model.FileShareSpec, error) {
	fshare, err := c.getFileShare(ctx, fshareID)
	if !IsAdminContext(ctx) || err == nil {
		return fshare, err
	}
	fshares, err := c.ListFileShares(ctx)
	if err != nil {
		return nil, err
	}
	for _, f := range fshares {
		if f.Id == fshareID {
			return f, nil
		}
	}
	return nil, fmt.Errorf("specified fileshare(%s) can't find", fshareID)
}

func (c *Client) getFileShare(ctx *c.Context, fshareID string) (*model.FileShareSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateFileShareURL(urls.Etcd, ctx.TenantId, fshareID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when get fileshare in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var fshare = &model.FileShareSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), fshare); err != nil {
		log.Error("when parsing fileshare in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return fshare, nil
}

// UpdateFileShare ...
func (c *Client) UpdateFileShare(ctx *c.Context, fshare *model.FileShareSpec) (*model.FileShareSpec, error) {
	result, err := c.GetFileShare(ctx, fshare.Id)
	if err != nil {
		return nil, err
	}
	if fshare.Name != "" {
		result.Name = fshare.Name
	}
	if fshare.Description != "" {
		result.Description = fshare.Description
	}
	if fshare.ExportLocations != nil {
		result.ExportLocations = fshare.ExportLocations
	}
	if fshare.Protocols != nil {
		result.Protocols = fshare.Protocols
	}
	if fshare.Metadata != nil {
		result.Metadata = fshare.Metadata
	}
	if fshare.Status != "" {
		result.Status = fshare.Status
	}
	if fshare.PoolId != "" {
		result.PoolId = fshare.PoolId
	}

	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	log.V(5).Infof("update file share object %+v into db", result)

	body, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateFileShareURL(urls.Etcd, result.TenantId, fshare.Id),
		NewContent: string(body),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when update fileshare in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteFileShareAcl
func (c *Client) DeleteFileShareAcl(ctx *c.Context, aclID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		fshare, err := c.GetFileShareAcl(ctx, aclID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = fshare.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateFileShareAclURL(urls.Etcd, tenantId, aclID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when delete fileshare in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// DeleteFileShare
func (c *Client) DeleteFileShare(ctx *c.Context, fileshareID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		fshare, err := c.GetFileShare(ctx, fileshareID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = fshare.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateFileShareURL(urls.Etcd, tenantId, fileshareID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when delete fileshare in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreateFileShareSnapshot
func (c *Client) CreateFileShareSnapshot(ctx *c.Context, snp *model.FileShareSnapshotSpec) (*model.FileShareSnapshotSpec, error) {
	snp.TenantId = ctx.TenantId
	snpBody, err := json.Marshal(snp)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateFileShareSnapshotURL(urls.Etcd, ctx.TenantId, snp.Id),
		Content: string(snpBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when create fileshare snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return snp, nil
}

func (c *Client) GetFileShareSnapshot(ctx *c.Context, snpID string) (*model.FileShareSnapshotSpec, error) {
	snap, err := c.getFileShareSnapshot(ctx, snpID)
	if !IsAdminContext(ctx) || err == nil {
		return snap, err
	}
	snaps, err := c.ListFileShareSnapshots(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range snaps {
		if v.Id == snpID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified fileshare snapshot(%s) can't find", snpID)
}

// GetFileShareSnapshot
func (c *Client) getFileShareSnapshot(ctx *c.Context, snpID string) (*model.FileShareSnapshotSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateFileShareSnapshotURL(urls.Etcd, ctx.TenantId, snpID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when get fileshare attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var fs = &model.FileShareSnapshotSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), fs); err != nil {
		log.Error("when parsing fileshare snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return fs, nil
}

// ListFileShareSnapshots
func (c *Client) ListFileShareSnapshots(ctx *c.Context) ([]*model.FileShareSnapshotSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateFileShareSnapshotURL(urls.Etcd, ctx.TenantId),
	}
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateFileShareSnapshotURL(urls.Etcd, "")
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when list fileshare snapshots in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var fss = []*model.FileShareSnapshotSpec{}
	if len(dbRes.Message) == 0 {
		return fss, nil
	}
	for _, msg := range dbRes.Message {
		var fs = &model.FileShareSnapshotSpec{}
		if err := json.Unmarshal([]byte(msg), fs); err != nil {
			log.Error("When parsing fileshare snapshot in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		fss = append(fss, fs)
	}
	return fss, nil
}

func (c *Client) ListFileShareSnapshotsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.FileShareSnapshotSpec, error) {
	fileshareSnapshots, err := c.ListFileShareSnapshots(ctx)
	if err != nil {
		log.Error("list fileshareSnapshots failed: ", err)
		return nil, err
	}

	tmpFileshareSnapshots := c.FilterAndSort(fileshareSnapshots, m, sortableKeysMap[typeFileShareSnapshots])
	var res = []*model.FileShareSnapshotSpec{}
	for _, data := range tmpFileshareSnapshots.([]interface{}) {
		res = append(res, data.(*model.FileShareSnapshotSpec))
	}
	return res, nil
}

// UpdateFileShareSnapshot
func (c *Client) UpdateFileShareSnapshot(ctx *c.Context, snpID string, snp *model.FileShareSnapshotSpec) (*model.FileShareSnapshotSpec, error) {
	result, err := c.GetFileShareSnapshot(ctx, snpID)
	if err != nil {
		return nil, err
	}
	if snp.Name != "" {
		result.Name = snp.Name
	}
	if snp.Description != "" {
		result.Description = snp.Description
	}
	if snp.Status != "" {
		result.Status = snp.Status
	}

	if snp.SnapshotSize > 0 {
		result.SnapshotSize = snp.SnapshotSize
	}

	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	result.Metadata = snp.Metadata

	atcBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateFileShareSnapshotURL(urls.Etcd, result.TenantId, snpID),
		NewContent: string(atcBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when update fileshare snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteFileShareSnapshot
func (c *Client) DeleteFileShareSnapshot(ctx *c.Context, snpID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		snap, err := c.GetFileShareSnapshot(ctx, snpID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = snap.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateFileShareSnapshotURL(urls.Etcd, tenantId, snpID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when delete fileshare snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// ********************** End Of FileShare *********************

// CreateDock
func (c *Client) CreateDock(ctx *c.Context, dck *model.DockSpec) (*model.DockSpec, error) {
	if dck.Id == "" {
		dck.Id = uuid.NewV4().String()
	}

	if dck.CreatedAt == "" {
		dck.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	dckBody, err := json.Marshal(dck)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateDockURL(urls.Etcd, "", dck.Id),
		Content: string(dckBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when create dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return dck, nil
}

// GetDock
func (c *Client) GetDock(ctx *c.Context, dckID string) (*model.DockSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, "", dckID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("when get dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var dck = &model.DockSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), dck); err != nil {
		log.Error("when parsing dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return dck, nil
}

// GetDockByPoolId
func (c *Client) GetDockByPoolId(ctx *c.Context, poolId string) (*model.DockSpec, error) {
	pool, err := c.GetPool(ctx, poolId)
	if err != nil {
		log.Error("Get pool failed in db: ", err)
		return nil, err
	}

	docks, err := c.ListDocks(ctx)
	if err != nil {
		log.Error("List docks failed failed in db: ", err)
		return nil, err
	}
	for _, dock := range docks {
		if pool.DockId == dock.Id {
			return dock, nil
		}
	}
	return nil, errors.New("Get dock failed by pool id: " + poolId)
}

// ListDocks
func (c *Client) ListDocks(ctx *c.Context) ([]*model.DockSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, ""),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list docks in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var dcks = []*model.DockSpec{}
	if len(dbRes.Message) == 0 {
		return dcks, nil
	}
	for _, msg := range dbRes.Message {
		var dck = &model.DockSpec{}
		if err := json.Unmarshal([]byte(msg), dck); err != nil {
			log.Error("When parsing dock in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		dcks = append(dcks, dck)
	}
	return dcks, nil
}

func (c *Client) ListDocksWithFilter(ctx *c.Context, m map[string][]string) ([]*model.DockSpec, error) {
	docks, err := c.ListDocks(ctx)
	if err != nil {
		log.Error("List docks failed: ", err.Error())
		return nil, err
	}

	tmpDocks := c.FilterAndSort(docks, m, sortableKeysMap[typeDocks])
	var res = []*model.DockSpec{}
	for _, data := range tmpDocks.([]interface{}) {
		res = append(res, data.(*model.DockSpec))
	}
	return res, nil
}

// UpdateDock
func (c *Client) UpdateDock(ctx *c.Context, dckID, name, desp string) (*model.DockSpec, error) {
	dck, err := c.GetDock(ctx, dckID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		dck.Name = name
	}
	if desp != "" {
		dck.Description = desp
	}
	dck.UpdatedAt = time.Now().Format(constants.TimeFormat)

	dckBody, err := json.Marshal(dck)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateDockURL(urls.Etcd, "", dckID),
		NewContent: string(dckBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update dock in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return dck, nil
}

// DeleteDock
func (c *Client) DeleteDock(ctx *c.Context, dckID string) error {
	dbReq := &Request{
		Url: urls.GenerateDockURL(urls.Etcd, "", dckID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete dock in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreatePool
func (c *Client) CreatePool(ctx *c.Context, pol *model.StoragePoolSpec) (*model.StoragePoolSpec, error) {
	if pol.Id == "" {
		pol.Id = uuid.NewV4().String()
	}

	if pol.CreatedAt == "" {
		pol.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	polBody, err := json.Marshal(pol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GeneratePoolURL(urls.Etcd, "", pol.Id),
		Content: string(polBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create pol in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return pol, nil
}

func (c *Client) ListPoolsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.StoragePoolSpec, error) {
	pools, err := c.ListPools(ctx)
	if err != nil {
		log.Error("List pools failed: ", err.Error())
		return nil, err
	}

	tmpPools := c.FilterAndSort(pools, m, sortableKeysMap[typePools])
	var res = []*model.StoragePoolSpec{}
	for _, data := range tmpPools.([]interface{}) {
		res = append(res, data.(*model.StoragePoolSpec))
	}
	return res, nil
}

// GetPool
func (c *Client) GetPool(ctx *c.Context, polID string) (*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, "", polID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var pol = &model.StoragePoolSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), pol); err != nil {
		log.Error("When parsing pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return pol, nil
}

//ListAvailabilityZones
func (c *Client) ListAvailabilityZones(ctx *c.Context) ([]string, error) {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, ""),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("Failed to get AZ for pools in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	var azs = []string{}
	if len(dbRes.Message) == 0 {
		return azs, nil
	}
	for _, msg := range dbRes.Message {
		var pol = &model.StoragePoolSpec{}
		if err := json.Unmarshal([]byte(msg), pol); err != nil {
			log.Error("When parsing pool in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		azs = append(azs, pol.AvailabilityZone)
	}
	//remove redundant AZ
	azs = utils.RvRepElement(azs)
	return azs, nil
}

// ListPools
func (c *Client) ListPools(ctx *c.Context) ([]*model.StoragePoolSpec, error) {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, ""),
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list pools in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var pols = []*model.StoragePoolSpec{}
	if len(dbRes.Message) == 0 {
		return pols, nil
	}
	for _, msg := range dbRes.Message {
		var pol = &model.StoragePoolSpec{}
		if err := json.Unmarshal([]byte(msg), pol); err != nil {
			log.Error("When parsing pool in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		pols = append(pols, pol)
	}
	return pols, nil
}

// UpdatePool
func (c *Client) UpdatePool(ctx *c.Context, polID, name, desp string, usedCapacity int64, used bool) (*model.StoragePoolSpec, error) {
	pol, err := c.GetPool(ctx, polID)
	if err != nil {
		return nil, err
	}
	if name != "" {
		pol.Name = name
	}
	if desp != "" {
		pol.Description = desp
	}
	pol.UpdatedAt = time.Now().Format(constants.TimeFormat)

	polBody, err := json.Marshal(pol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GeneratePoolURL(urls.Etcd, "", polID),
		NewContent: string(polBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update pool in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return pol, nil
}

// DeletePool
func (c *Client) DeletePool(ctx *c.Context, polID string) error {
	dbReq := &Request{
		Url: urls.GeneratePoolURL(urls.Etcd, "", polID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete pool in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreateProfile
func (c *Client) CreateProfile(ctx *c.Context, prf *model.ProfileSpec) (*model.ProfileSpec, error) {
	if prf.Id == "" {
		prf.Id = uuid.NewV4().String()
	}
	if prf.CreatedAt == "" {
		prf.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	// profile name must be unique.
	if _, err := c.getProfileByName(ctx, prf.Name); err == nil {
		return nil, fmt.Errorf("the profile name '%s' already exists", prf.Name)
	}

	prfBody, err := json.Marshal(prf)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateProfileURL(urls.Etcd, "", prf.Id),
		Content: string(prfBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return prf, nil
}

// GetProfile
func (c *Client) GetProfile(ctx *c.Context, prfID string) (*model.ProfileSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateProfileURL(urls.Etcd, "", prfID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var prf = &model.ProfileSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), prf); err != nil {
		log.Error("When parsing profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return prf, nil
}

func (c *Client) getProfileByName(ctx *c.Context, name string) (*model.ProfileSpec, error) {
	profiles, err := c.ListProfiles(ctx)
	if err != nil {
		log.Error("List profile failed: ", err)
		return nil, err
	}

	for _, profile := range profiles {
		if profile.Name == name {
			return profile, nil
		}
	}
	var msg = fmt.Sprintf("can't find profile(name: %s)", name)
	return nil, model.NewNotFoundError(msg)
}

func (c *Client) getProfileByNameAndType(ctx *c.Context, name, storageType string) (*model.ProfileSpec, error) {
	profiles, err := c.ListProfiles(ctx)
	if err != nil {
		log.Error("List profile failed: ", err)
		return nil, err
	}

	for _, profile := range profiles {
		if profile.Name == name && profile.StorageType == storageType {
			return profile, nil
		}
	}
	var msg = fmt.Sprintf("can't find profile(name: %s, storageType:%s)", name, storageType)
	return nil, model.NewNotFoundError(msg)
}

// GetDefaultProfile
func (c *Client) GetDefaultProfile(ctx *c.Context) (*model.ProfileSpec, error) {
	return c.getProfileByNameAndType(ctx, defaultBlockProfileName, typeBlock)
}

// GetDefaultProfileFileShare
func (c *Client) GetDefaultProfileFileShare(ctx *c.Context) (*model.ProfileSpec, error) {
	return c.getProfileByNameAndType(ctx, defaultFileProfileName, typeFile)
}

// ListProfiles
func (c *Client) ListProfiles(ctx *c.Context) ([]*model.ProfileSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateProfileURL(urls.Etcd, ""),
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list profiles in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var prfs = []*model.ProfileSpec{}
	if len(dbRes.Message) == 0 {
		return prfs, nil
	}
	for _, msg := range dbRes.Message {
		var prf = &model.ProfileSpec{}
		if err := json.Unmarshal([]byte(msg), prf); err != nil {
			log.Error("When parsing profile in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		prfs = append(prfs, prf)
	}
	return prfs, nil
}

func (c *Client) ListProfilesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.ProfileSpec, error) {
	profiles, err := c.ListProfiles(ctx)
	if err != nil {
		log.Error("List profiles failed: ", err)
		return nil, err
	}

	tmpProfiles := c.FilterAndSort(profiles, m, sortableKeysMap[typeProfiles])
	var res = []*model.ProfileSpec{}
	for _, data := range tmpProfiles.([]interface{}) {
		res = append(res, data.(*model.ProfileSpec))
	}
	return res, nil
}

// UpdateProfile
func (c *Client) UpdateProfile(ctx *c.Context, prfID string, input *model.ProfileSpec) (*model.ProfileSpec, error) {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return nil, err
	}
	if name := input.Name; name != "" {
		prf.Name = name
	}
	if desp := input.Description; desp != "" {
		prf.Description = desp
	}
	prf.UpdatedAt = time.Now().Format(constants.TimeFormat)

	if props := input.CustomProperties; len(props) != 0 {
		if prf.CustomProperties == nil {
			prf.CustomProperties = make(map[string]interface{})
		}
		for k, v := range props {
			prf.CustomProperties[k] = v
		}
	}

	prf.UpdatedAt = time.Now().Format(constants.TimeFormat)

	prfBody, err := json.Marshal(prf)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateProfileURL(urls.Etcd, "", prfID),
		NewContent: string(prfBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update profile in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return prf, nil
}

// DeleteProfile
func (c *Client) DeleteProfile(ctx *c.Context, prfID string) error {
	dbReq := &Request{
		Url: urls.GenerateProfileURL(urls.Etcd, "", prfID),
	}
	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete profile in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// AddCustomProperty
func (c *Client) AddCustomProperty(ctx *c.Context, prfID string, ext model.CustomPropertiesSpec) (*model.CustomPropertiesSpec, error) {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return nil, err
	}

	if prf.CustomProperties == nil {
		prf.CustomProperties = make(map[string]interface{})
	}

	for k, v := range ext {
		prf.CustomProperties[k] = v
	}

	prf.UpdatedAt = time.Now().Format(constants.TimeFormat)

	if _, err = c.CreateProfile(ctx, prf); err != nil {
		return nil, err
	}
	return &prf.CustomProperties, nil
}

// ListCustomProperties
func (c *Client) ListCustomProperties(ctx *c.Context, prfID string) (*model.CustomPropertiesSpec, error) {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return nil, err
	}
	return &prf.CustomProperties, nil
}

// RemoveCustomProperty
func (c *Client) RemoveCustomProperty(ctx *c.Context, prfID, customKey string) error {
	prf, err := c.GetProfile(ctx, prfID)
	if err != nil {
		return err
	}

	delete(prf.CustomProperties, customKey)
	if _, err = c.CreateProfile(ctx, prf); err != nil {
		return err
	}
	return nil
}

// CreateVolume
func (c *Client) CreateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	profiles, err := c.ListProfiles(ctx)
	if err != nil {
		return nil, err
	} else if len(profiles) == 0 {
		return nil, errors.New("No profile in db.")
	}

	vol.TenantId = ctx.TenantId
	volBody, err := json.Marshal(vol)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId, vol.Id),
		Content: string(volBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return vol, nil
}

// GetVolume
func (c *Client) GetVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	vol, err := c.getVolume(ctx, volID)
	if !IsAdminContext(ctx) || err == nil {
		return vol, err
	}
	vols, err := c.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range vols {
		if v.Id == volID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified volume(%s) can't find", volID)
}

func (c *Client) getVolume(ctx *c.Context, volID string) (*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId, volID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vol = &model.VolumeSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vol); err != nil {
		log.Error("When parsing volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vol, nil
}

// ListVolumes
func (c *Client) ListVolumes(ctx *c.Context) ([]*model.VolumeSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId),
	}

	// Admin user should get all volumes including the volumes whose tenant is not admin.
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateVolumeURL(urls.Etcd, "")
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volumes in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vols = []*model.VolumeSpec{}
	if len(dbRes.Message) == 0 {
		return vols, nil
	}
	for _, msg := range dbRes.Message {
		var vol = &model.VolumeSpec{}
		if err := json.Unmarshal([]byte(msg), vol); err != nil {
			log.Error("When parsing volume in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		vols = append(vols, vol)
	}
	return vols, nil
}

// ListVolumesByProfileId
func (c *Client) ListVolumesByProfileId(ctx *c.Context, prfID string) ([]string, error) {
	vols, err := c.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	var resvols []string
	for _, v := range vols {
		if v.ProfileId == prfID {
			resvols = append(resvols, v.Name)
		}
	}

	return resvols, nil

}

var volume_sortKey string

type VolumeSlice []*model.VolumeSpec

func (volume VolumeSlice) Len() int { return len(volume) }

func (volume VolumeSlice) Swap(i, j int) { volume[i], volume[j] = volume[j], volume[i] }

func (volume VolumeSlice) Less(i, j int) bool {
	switch volume_sortKey {

	case "ID":
		return volume[i].Id < volume[j].Id
	case "NAME":
		return volume[i].Name < volume[j].Name
	case "STATUS":
		return volume[i].Status < volume[j].Status
	case "AVAILABILITYZONE":
		return volume[i].AvailabilityZone < volume[j].AvailabilityZone
	case "PROFILEID":
		return volume[i].ProfileId < volume[j].ProfileId
	case "TENANTID":
		return volume[i].TenantId < volume[j].TenantId
	case "SIZE":
		return volume[i].Size < volume[j].Size
	case "POOLID":
		return volume[i].PoolId < volume[j].PoolId
	case "DESCRIPTION":
		return volume[i].Description < volume[j].Description
	case "GROUPID":
		return volume[i].GroupId < volume[j].GroupId
		// TODO:case "lun_id" (admin_only)
	}
	return false
}

func (c *Client) FindVolumeValue(k string, p *model.VolumeSpec) string {
	switch k {
	case "Id":
		return p.Id
	case "CreatedAt":
		return p.CreatedAt
	case "UpdatedAt":
		return p.UpdatedAt
	case "TenantId":
		return p.TenantId
	case "UserId":
		return p.UserId
	case "Name":
		return p.Name
	case "Description":
		return p.Description
	case "AvailabilityZone":
		return p.AvailabilityZone
	case "Size":
		return strconv.FormatInt(p.Size, 10)
	case "Status":
		return p.Status
	case "PoolId":
		return p.PoolId
	case "ProfileId":
		return p.ProfileId
	case "GroupId":
		return p.GroupId
	case "DurableName":
		return p.Identifier.DurableName
	case "DurableNameFormat":
		return p.Identifier.DurableNameFormat
	}
	return ""
}

func (c *Client) SortVolumes(volumes []*model.VolumeSpec, p *Parameter) []*model.VolumeSpec {
	volume_sortKey = p.sortKey

	if strings.EqualFold(p.sortDir, "asc") {
		sort.Sort(VolumeSlice(volumes))

	} else {
		sort.Sort(sort.Reverse(VolumeSlice(volumes)))
	}
	return volumes
}

func (c *Client) ListVolumesWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSpec, error) {
	volumes, err := c.ListVolumes(ctx)
	if err != nil {
		log.Error("List volumes failed: ", err)
		return nil, err
	}
	// If DurableName is there in filter we need to parse the sub structure 'identifier' to filter out matching volume spec
	var vols = []*model.VolumeSpec{}
	if val, ok := m["DurableName"]; ok {
		for _, vol := range volumes {
			v := c.FindVolumeValue("DurableName", vol)
			if v == val[0] {
				vols = append(vols, vol)
				return vols, nil
			}
		}

		return vols, nil
	}

	tmpVolumes := c.FilterAndSort(volumes, m, sortableKeysMap[typeVolumes])
	var res = []*model.VolumeSpec{}
	for _, data := range tmpVolumes.([]interface{}) {
		res = append(res, data.(*model.VolumeSpec))
	}
	return res, nil
}

// UpdateVolume ...
func (c *Client) UpdateVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	result, err := c.GetVolume(ctx, vol.Id)
	if err != nil {
		return nil, err
	}
	if vol.Name != "" {
		result.Name = vol.Name
	}
	if vol.AvailabilityZone != "" {
		result.AvailabilityZone = vol.AvailabilityZone
	}
	if vol.Description != "" {
		result.Description = vol.Description
	}
	if vol.Metadata != nil {
		result.Metadata = utils.MergeStringMaps(result.Metadata, vol.Metadata)
	}
	if vol.Identifier != nil {
		result.Identifier = vol.Identifier
	}
	if vol.PoolId != "" {
		result.PoolId = vol.PoolId
	}
	if vol.ProfileId != "" {
		result.ProfileId = vol.ProfileId
	}
	if vol.Size != 0 {
		result.Size = vol.Size
	}
	if vol.Status != "" {
		result.Status = vol.Status
	}
	if vol.ReplicationDriverData != nil {
		result.ReplicationDriverData = vol.ReplicationDriverData
	}
	if vol.MultiAttach {
		result.MultiAttach = vol.MultiAttach
	}
	if vol.GroupId != "" {
		result.GroupId = vol.GroupId
	}

	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	body, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeURL(urls.Etcd, result.TenantId, vol.Id),
		NewContent: string(body),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteVolume
func (c *Client) DeleteVolume(ctx *c.Context, volID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		vol, err := c.GetVolume(ctx, volID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = vol.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateVolumeURL(urls.Etcd, tenantId, volID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// ExtendVolume ...
func (c *Client) ExtendVolume(ctx *c.Context, vol *model.VolumeSpec) (*model.VolumeSpec, error) {
	result, err := c.GetVolume(ctx, vol.Id)
	if err != nil {
		return nil, err
	}

	if vol.Size > 0 {
		result.Size = vol.Size
	}
	result.Status = vol.Status
	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	body, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeURL(urls.Etcd, ctx.TenantId, vol.Id),
		NewContent: string(body),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When extend volume in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// CreateVolumeAttachment
func (c *Client) CreateVolumeAttachment(ctx *c.Context, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	if attachment.Id == "" {
		attachment.Id = uuid.NewV4().String()
	}
	attachment.CreatedAt = time.Now().Format(constants.TimeFormat)
	attachment.TenantId = ctx.TenantId

	atcBody, err := json.Marshal(attachment)
	if err != nil {
		return nil, err
	}
	dbReq := &Request{
		Url:     urls.GenerateAttachmentURL(urls.Etcd, ctx.TenantId, attachment.Id),
		Content: string(atcBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return attachment, nil
}

func (c *Client) GetVolumeAttachment(ctx *c.Context, attachmentId string) (*model.VolumeAttachmentSpec, error) {
	attach, err := c.getVolumeAttachment(ctx, attachmentId)
	if !IsAdminContext(ctx) || err == nil {
		return attach, err
	}
	attachs, err := c.ListVolumeAttachments(ctx, "")
	if err != nil {
		return nil, err
	}
	for _, v := range attachs {
		if v.Id == attachmentId {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified volume attachment(%s) can't find", attachmentId)
}

// GetVolumeAttachment
func (c *Client) getVolumeAttachment(ctx *c.Context, attachmentId string) (*model.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateAttachmentURL(urls.Etcd, ctx.TenantId, attachmentId),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var atc = &model.VolumeAttachmentSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), atc); err != nil {
		log.Error("When parsing volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return atc, nil
}

// ListVolumeAttachments
func (c *Client) ListVolumeAttachments(ctx *c.Context, volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateAttachmentURL(urls.Etcd, ctx.TenantId),
	}
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateAttachmentURL(urls.Etcd, "")
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volume attachments in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var atcs = []*model.VolumeAttachmentSpec{}
	for _, msg := range dbRes.Message {
		var atc = &model.VolumeAttachmentSpec{}
		if err := json.Unmarshal([]byte(msg), atc); err != nil {
			log.Error("When parsing volume attachment in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}

		if len(volumeId) == 0 || atc.VolumeId == volumeId {
			atcs = append(atcs, atc)
		}
	}
	return atcs, nil

}

func (c *Client) ListVolumeAttachmentsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeAttachmentSpec, error) {
	var volumeId string
	if v, ok := m["VolumeId"]; ok {
		volumeId = v[0]
	}
	attachments, err := c.ListVolumeAttachments(ctx, volumeId)
	if err != nil {
		log.Error("List volumes failed: ", err)
		return nil, err
	}

	tmpAttachments := c.FilterAndSort(attachments, m, sortableKeysMap[typeAttachments])
	var res = []*model.VolumeAttachmentSpec{}
	for _, data := range tmpAttachments.([]interface{}) {
		res = append(res, data.(*model.VolumeAttachmentSpec))
	}
	return res, nil
}

// UpdateVolumeAttachment
func (c *Client) UpdateVolumeAttachment(ctx *c.Context, attachmentId string, attachment *model.VolumeAttachmentSpec) (*model.VolumeAttachmentSpec, error) {
	result, err := c.GetVolumeAttachment(ctx, attachmentId)
	if err != nil {
		return nil, err
	}
	if len(attachment.Mountpoint) > 0 {
		result.Mountpoint = attachment.Mountpoint
	}
	if len(attachment.Status) > 0 {
		result.Status = attachment.Status
	}

	// Update DriverVolumeType
	if len(attachment.DriverVolumeType) > 0 {
		result.DriverVolumeType = attachment.DriverVolumeType
	}
	// Update connectionData
	// Debug
	log.V(8).Infof("etcd: update volume attachment connection data from db: %v", result.ConnectionData)
	log.V(8).Infof("etcd: update volume attachment connection data from target: %v", attachment.ConnectionData)

	if attachment.ConnectionData != nil {
		if result.ConnectionData == nil {
			result.ConnectionData = make(map[string]interface{})
		}

		for k, v := range attachment.ConnectionData {
			result.ConnectionData[k] = v
		}
	}
	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	atcBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateAttachmentURL(urls.Etcd, result.TenantId, attachmentId),
		NewContent: string(atcBody),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteVolumeAttachment
func (c *Client) DeleteVolumeAttachment(ctx *c.Context, attachmentId string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		attach, err := c.GetVolumeAttachment(ctx, attachmentId)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = attach.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateAttachmentURL(urls.Etcd, tenantId, attachmentId),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume attachment in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

// CreateVolumeSnapshot
func (c *Client) CreateVolumeSnapshot(ctx *c.Context, snp *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	snp.TenantId = ctx.TenantId
	snpBody, err := json.Marshal(snp)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateSnapshotURL(urls.Etcd, ctx.TenantId, snp.Id),
		Content: string(snpBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return snp, nil
}

func (c *Client) GetVolumeSnapshot(ctx *c.Context, snpID string) (*model.VolumeSnapshotSpec, error) {
	snap, err := c.getVolumeSnapshot(ctx, snpID)
	if !IsAdminContext(ctx) || err == nil {
		return snap, err
	}
	snaps, err := c.ListVolumeSnapshots(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range snaps {
		if v.Id == snpID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified volume snapshot(%s) can't find", snpID)
}

// GetVolumeSnapshot
func (c *Client) getVolumeSnapshot(ctx *c.Context, snpID string) (*model.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateSnapshotURL(urls.Etcd, ctx.TenantId, snpID),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume attachment in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vs = &model.VolumeSnapshotSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vs); err != nil {
		log.Error("When parsing volume snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vs, nil
}

// ListVolumeSnapshots
func (c *Client) ListVolumeSnapshots(ctx *c.Context) ([]*model.VolumeSnapshotSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateSnapshotURL(urls.Etcd, ctx.TenantId),
	}
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateSnapshotURL(urls.Etcd, "")
	}
	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volume snapshots in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vss = []*model.VolumeSnapshotSpec{}
	if len(dbRes.Message) == 0 {
		return vss, nil
	}
	for _, msg := range dbRes.Message {
		var vs = &model.VolumeSnapshotSpec{}
		if err := json.Unmarshal([]byte(msg), vs); err != nil {
			log.Error("When parsing volume snapshot in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		vss = append(vss, vs)
	}
	return vss, nil
}

func (c *Client) ListVolumeSnapshotsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeSnapshotSpec, error) {
	volumeSnapshots, err := c.ListVolumeSnapshots(ctx)
	if err != nil {
		log.Error("List volumeSnapshots failed: ", err)
		return nil, err
	}

	tmpVolumeSnapshots := c.FilterAndSort(volumeSnapshots, m, sortableKeysMap[typeVolumeSnapshots])
	var res = []*model.VolumeSnapshotSpec{}
	for _, data := range tmpVolumeSnapshots.([]interface{}) {
		res = append(res, data.(*model.VolumeSnapshotSpec))
	}
	return res, nil
}

// UpdateVolumeSnapshot
func (c *Client) UpdateVolumeSnapshot(ctx *c.Context, snpID string, snp *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	result, err := c.GetVolumeSnapshot(ctx, snpID)
	if err != nil {
		return nil, err
	}
	if snp.Name != "" {
		result.Name = snp.Name
	}
	if snp.Metadata != nil {
		result.Metadata = utils.MergeStringMaps(result.Metadata, snp.Metadata)
	}
	if snp.Size > 0 {
		result.Size = snp.Size
	}
	if snp.VolumeId != "" {
		result.VolumeId = snp.VolumeId
	}
	if snp.Description != "" {
		result.Description = snp.Description
	}
	if snp.Status != "" {
		result.Status = snp.Status
	}
	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)

	atcBody, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateSnapshotURL(urls.Etcd, result.TenantId, snpID),
		NewContent: string(atcBody),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume snapshot in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

// DeleteVolumeSnapshot
func (c *Client) DeleteVolumeSnapshot(ctx *c.Context, snpID string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		snap, err := c.GetVolumeSnapshot(ctx, snpID)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = snap.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateSnapshotURL(urls.Etcd, tenantId, snpID),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume snapshot in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) CreateReplication(ctx *c.Context, r *model.ReplicationSpec) (*model.ReplicationSpec, error) {
	if r.Id == "" {
		r.Id = uuid.NewV4().String()
	}

	r.TenantId = ctx.TenantId
	r.CreatedAt = time.Now().Format(constants.TimeFormat)
	rBody, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req := &Request{
		Url:     urls.GenerateReplicationURL(urls.Etcd, ctx.TenantId, r.Id),
		Content: string(rBody),
	}
	resp := c.Create(req)
	if resp.Status != "Success" {
		log.Error("When create replication in db:", resp.Error)
		return nil, errors.New(resp.Error)
	}

	return r, nil
}

func (c *Client) GetReplication(ctx *c.Context, replicationId string) (*model.ReplicationSpec, error) {
	replication, err := c.getReplication(ctx, replicationId)
	if !IsAdminContext(ctx) || err == nil {
		return replication, err
	}
	replications, err := c.ListReplication(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range replications {
		if r.Id == replicationId {
			return r, nil
		}
	}
	return nil, fmt.Errorf("specified replication(%s) can't find", replicationId)
}

func (c *Client) GetReplicationByVolumeId(ctx *c.Context, volumeId string) (*model.ReplicationSpec, error) {
	replications, err := c.ListReplication(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range replications {
		if volumeId == r.PrimaryVolumeId || volumeId == r.SecondaryVolumeId {
			return r, nil
		}
	}
	return nil, model.NewNotFoundError(fmt.Sprintf("can't find specified replication by volume id %s", volumeId))
}

func (c *Client) getReplication(ctx *c.Context, replicationId string) (*model.ReplicationSpec, error) {
	req := &Request{
		Url: urls.GenerateReplicationURL(urls.Etcd, ctx.TenantId, replicationId),
	}
	resp := c.Get(req)
	if resp.Status != "Success" {
		log.Error("When get pool in db:", resp.Error)
		return nil, errors.New(resp.Error)
	}

	var r = &model.ReplicationSpec{}
	if err := json.Unmarshal([]byte(resp.Message[0]), r); err != nil {
		log.Error("When parsing replication in db:", resp.Error)
		return nil, errors.New(resp.Error)
	}
	return r, nil
}

func (c *Client) ListReplication(ctx *c.Context) ([]*model.ReplicationSpec, error) {
	req := &Request{
		Url: urls.GenerateReplicationURL(urls.Etcd, ctx.TenantId),
	}
	if IsAdminContext(ctx) {
		req.Url = urls.GenerateReplicationURL(urls.Etcd, "")
	}
	resp := c.List(req)
	if resp.Status != "Success" {
		log.Error("When list replication in db:", resp.Error)
		return nil, errors.New(resp.Error)
	}

	var replicas = []*model.ReplicationSpec{}
	if len(resp.Message) == 0 {
		return replicas, nil
	}
	for _, msg := range resp.Message {
		var r = &model.ReplicationSpec{}
		if err := json.Unmarshal([]byte(msg), r); err != nil {
			log.Error("When parsing replication in db:", resp.Error)
			return nil, errors.New(resp.Error)
		}
		replicas = append(replicas, r)
	}
	return replicas, nil

}

func (c *Client) filterByName(param map[string][]string, spec interface{}, filterList map[string]interface{}) bool {
	v := reflect.ValueOf(spec)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for key := range param {
		_, ok := filterList[key]
		if !ok {
			continue
		}
		filed := v.FieldByName(key)
		if !filed.IsValid() {
			continue
		}
		paramVal := param[key][0]
		var val string
		switch filed.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = strconv.FormatInt(filed.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = strconv.FormatUint(filed.Uint(), 10)
		case reflect.String:
			val = filed.String()
		default:
			return false
		}
		if !strings.EqualFold(paramVal, val) {
			return false
		}
	}

	return true
}

func (c *Client) SelectReplication(param map[string][]string, replications []*model.ReplicationSpec) []*model.ReplicationSpec {
	if !c.SelectOrNot(param) {
		return replications
	}

	filterList := map[string]interface{}{
		"Id":                nil,
		"CreatedAt":         nil,
		"UpdatedAt":         nil,
		"Name":              nil,
		"Description":       nil,
		"PrimaryVolumeId":   nil,
		"SecondaryVolumeId": nil,
	}

	var rlist = []*model.ReplicationSpec{}
	for _, r := range replications {
		if c.filterByName(param, r, filterList) {
			rlist = append(rlist, r)
		}
	}
	return rlist

}

type ReplicationsCompareFunc func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool

var replicationsCompareFunc ReplicationsCompareFunc

type ReplicationSlice []*model.ReplicationSpec

func (r ReplicationSlice) Len() int           { return len(r) }
func (r ReplicationSlice) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ReplicationSlice) Less(i, j int) bool { return replicationsCompareFunc(r[i], r[j]) }

var replicationSortKey2Func = map[string]ReplicationsCompareFunc{
	"ID":   func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool { return a.Id > b.Id },
	"NAME": func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool { return a.Name > b.Name },
	"REPLICATIONSTATUS": func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool {
		return a.ReplicationStatus > b.ReplicationStatus
	},
	"AVAILABILITYZONE": func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool {
		return a.AvailabilityZone > b.AvailabilityZone
	},
	"PROFILEID": func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool { return a.ProfileId > b.ProfileId },
	"TENANTID":  func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool { return a.TenantId > b.TenantId },
	"POOLID":    func(a *model.ReplicationSpec, b *model.ReplicationSpec) bool { return a.PoolId > b.PoolId },
}

func (c *Client) SortReplications(replications []*model.ReplicationSpec, p *Parameter) []*model.ReplicationSpec {
	replicationsCompareFunc = replicationSortKey2Func[p.sortKey]

	if strings.EqualFold(p.sortDir, "asc") {
		sort.Sort(ReplicationSlice(replications))

	} else {
		sort.Sort(sort.Reverse(ReplicationSlice(replications)))
	}
	return replications
}

func (c *Client) ListReplicationWithFilter(ctx *c.Context, m map[string][]string) ([]*model.ReplicationSpec, error) {
	replicas, err := c.ListReplication(ctx)
	if err != nil {
		log.Error("List replications failed: ", err)
		return nil, err
	}

	rlist := c.SelectReplication(m, replicas)

	var sortKeys []string
	for k := range replicationSortKey2Func {
		sortKeys = append(sortKeys, k)
	}
	p := c.ParameterFilter(m, len(rlist), sortKeys)
	return c.SortReplications(rlist, p)[p.beginIdx:p.endIdx], nil
}

func (c *Client) DeleteReplication(ctx *c.Context, replicationId string) error {
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		r, err := c.GetReplication(ctx, replicationId)
		if err != nil {
			return err
		}
		tenantId = r.TenantId
	}
	req := &Request{
		Url: urls.GenerateReplicationURL(urls.Etcd, tenantId, replicationId),
	}
	reps := c.Delete(req)
	if reps.Status != "Success" {
		log.Error("When delete replication in db:", reps.Error)
		return errors.New(reps.Error)
	}
	return nil
}

func (c *Client) UpdateReplication(ctx *c.Context, replicationId string, input *model.ReplicationSpec) (*model.ReplicationSpec, error) {
	r, err := c.GetReplication(ctx, replicationId)
	if err != nil {
		return nil, err
	}
	if input.ProfileId != "" {
		r.ProfileId = input.ProfileId
	}
	if input.Name != "" {
		r.Name = input.Name
	}
	if input.Description != "" {
		r.Description = input.Description
	}
	if input.PrimaryReplicationDriverData != nil {
		r.PrimaryReplicationDriverData = input.PrimaryReplicationDriverData
	}
	if input.SecondaryReplicationDriverData != nil {
		r.SecondaryReplicationDriverData = input.SecondaryReplicationDriverData
	}
	if input.Metadata != nil {
		r.Metadata = utils.MergeStringMaps(r.Metadata, input.Metadata)
	}
	if input.ReplicationStatus != "" {
		r.ReplicationStatus = input.ReplicationStatus
	}

	r.UpdatedAt = time.Now().Format(constants.TimeFormat)

	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		tenantId = r.TenantId
	}
	req := &Request{
		Url:        urls.GenerateReplicationURL(urls.Etcd, tenantId, replicationId),
		NewContent: string(b),
	}
	resp := c.Update(req)
	if resp.Status != "Success" {
		log.Error("When update replication in db:", resp.Error)
		return nil, errors.New(resp.Error)
	}
	return r, nil
}
func (c *Client) CreateVolumeGroup(ctx *c.Context, vg *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	vg.TenantId = ctx.TenantId
	vgBody, err := json.Marshal(vg)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vg.Id),
		Content: string(vgBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return vg, nil
}

func (c *Client) GetVolumeGroup(ctx *c.Context, vgId string) (*model.VolumeGroupSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vgId),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var vg = &model.VolumeGroupSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), vg); err != nil {
		log.Error("When parsing volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vg, nil
}

func (c *Client) UpdateVolumeGroup(ctx *c.Context, vgUpdate *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	vg, err := c.GetVolumeGroup(ctx, vgUpdate.Id)
	if err != nil {
		return nil, err
	}
	if vgUpdate.Name != "" && vgUpdate.Name != vg.Name {
		vg.Name = vgUpdate.Name
	}
	if vgUpdate.AvailabilityZone != "" && vgUpdate.AvailabilityZone != vg.AvailabilityZone {
		vg.AvailabilityZone = vgUpdate.AvailabilityZone
	}
	if vgUpdate.Description != "" && vgUpdate.Description != vg.Description {
		vg.Description = vgUpdate.Description
	}
	if vgUpdate.PoolId != "" && vgUpdate.PoolId != vg.PoolId {
		vg.PoolId = vgUpdate.PoolId
	}
	if vg.Status != "" && vgUpdate.Status != vg.Status {
		vg.Status = vgUpdate.Status
	}
	if vgUpdate.PoolId != "" && vgUpdate.PoolId != vg.PoolId {
		vg.PoolId = vgUpdate.PoolId
	}
	if vgUpdate.CreatedAt != "" && vgUpdate.CreatedAt != vg.CreatedAt {
		vg.CreatedAt = vgUpdate.CreatedAt
	}
	if vgUpdate.UpdatedAt != "" && vgUpdate.UpdatedAt != vg.UpdatedAt {
		vg.UpdatedAt = vgUpdate.UpdatedAt
	}

	vgBody, err := json.Marshal(vg)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:        urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId, vgUpdate.Id),
		NewContent: string(vgBody),
	}
	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update volume group in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return vg, nil
}

func (c *Client) UpdateStatus(ctx *c.Context, in interface{}, status string) error {
	switch in.(type) {
	case *model.VolumeSnapshotSpec:
		snap := in.(*model.VolumeSnapshotSpec)
		snap.Status = status
		if _, errUpdate := c.UpdateVolumeSnapshot(ctx, snap.Id, snap); errUpdate != nil {
			log.Error("Error occurs when update volume snapshot status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.VolumeAttachmentSpec:
		attm := in.(*model.VolumeAttachmentSpec)
		attm.Status = status
		if _, errUpdate := c.UpdateVolumeAttachment(ctx, attm.Id, attm); errUpdate != nil {
			log.Error("Error occurred in dock module when update volume attachment status in db:", errUpdate)
			return errUpdate
		}

	case *model.VolumeSpec:
		volume := in.(*model.VolumeSpec)
		volume.Status = status
		if _, errUpdate := c.UpdateVolume(ctx, volume); errUpdate != nil {
			log.Error("When update volume status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.FileShareSpec:
		fileshare := in.(*model.FileShareSpec)
		fileshare.Status = status
		if _, errUpdate := c.UpdateFileShare(ctx, fileshare); errUpdate != nil {
			log.Error("when update fileshare status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.FileShareAclSpec:
		fileshareAcl := in.(*model.FileShareAclSpec)
		fileshareAcl.Status = status
		if _, errUpdate := c.UpdateFileShareAcl(ctx, fileshareAcl); errUpdate != nil {
			log.Error("when update fileshare acl status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.FileShareSnapshotSpec:
		fsnap := in.(*model.FileShareSnapshotSpec)
		fsnap.Status = status
		if _, errUpdate := c.UpdateFileShareSnapshot(ctx, fsnap.Id, fsnap); errUpdate != nil {
			log.Error("when update fileshare status in db:", errUpdate.Error())
			return errUpdate
		}

	case *model.VolumeGroupSpec:
		vg := in.(*model.VolumeGroupSpec)
		vg.Status = status
		if _, errUpdate := c.UpdateVolumeGroup(ctx, vg); errUpdate != nil {
			log.Error("When update volume status in db:", errUpdate.Error())
			return errUpdate
		}

	case []*model.VolumeSpec:
		vols := in.([]*model.VolumeSpec)
		if _, errUpdate := c.VolumesToUpdate(ctx, vols); errUpdate != nil {
			return errUpdate
		}

	case *model.ReplicationSpec:
		replica := in.(*model.ReplicationSpec)
		replica.ReplicationStatus = status
		if _, errUpdate := c.UpdateReplication(ctx, replica.Id, replica); errUpdate != nil {
			return errUpdate
		}
	}
	return nil
}

func (c *Client) ListVolumesByGroupId(ctx *c.Context, vgId string) ([]*model.VolumeSpec, error) {
	volumes, err := c.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}

	var volumesInSameGroup []*model.VolumeSpec
	for _, v := range volumes {
		if v.GroupId == vgId {
			volumesInSameGroup = append(volumesInSameGroup, v)
		}
	}

	return volumesInSameGroup, nil
}

func (c *Client) VolumesToUpdate(ctx *c.Context, volumeList []*model.VolumeSpec) ([]*model.VolumeSpec, error) {
	var volumeRefs []*model.VolumeSpec
	for _, values := range volumeList {
		v, err := c.UpdateVolume(ctx, values)
		if err != nil {
			return nil, err
		}
		volumeRefs = append(volumeRefs, v)
	}
	return volumeRefs, nil
}

// ListVolumes
func (c *Client) ListVolumeGroups(ctx *c.Context) ([]*model.VolumeGroupSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateVolumeGroupURL(urls.Etcd, ctx.TenantId),
	}
	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateVolumeGroupURL(urls.Etcd, "")
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list volume groups in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var groups []*model.VolumeGroupSpec
	if len(dbRes.Message) == 0 {
		return groups, nil
	}
	for _, msg := range dbRes.Message {
		var group = &model.VolumeGroupSpec{}
		if err := json.Unmarshal([]byte(msg), group); err != nil {
			log.Error("When parsing volume group in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (c *Client) DeleteVolumeGroup(ctx *c.Context, volumeGroupId string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		group, err := c.GetVolumeGroup(ctx, volumeGroupId)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = group.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateVolumeGroupURL(urls.Etcd, tenantId, volumeGroupId),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete volume group in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}

func (c *Client) ListSnapshotsByVolumeId(ctx *c.Context, volumeId string) ([]*model.VolumeSnapshotSpec, error) {
	snaps, err := c.ListVolumeSnapshots(ctx)
	if err != nil {
		return nil, err
	}

	var snapList []*model.VolumeSnapshotSpec
	for _, snap := range snaps {
		if snap.VolumeId == volumeId {
			snapList = append(snapList, snap)
		}
	}
	return snapList, nil
}

func (c *Client) ListAttachmentsByVolumeId(ctx *c.Context, volumeId string) ([]*model.VolumeAttachmentSpec, error) {
	return c.ListVolumeAttachments(ctx, volumeId)
}

func (c *Client) ListVolumeGroupsWithFilter(ctx *c.Context, m map[string][]string) ([]*model.VolumeGroupSpec, error) {
	vgs, err := c.ListVolumeGroups(ctx)
	if err != nil {
		log.Error("List volume groups failed: ", err)
		return nil, err
	}

	rlist := c.SelectVolumeGroup(m, vgs)

	var sortKeys []string
	for k := range volumeGroupSortKey2Func {
		sortKeys = append(sortKeys, k)
	}
	p := c.ParameterFilter(m, len(rlist), sortKeys)
	return c.SortVolumeGroups(rlist, p)[p.beginIdx:p.endIdx], nil

}

type VolumeGroupCompareFunc func(a *model.VolumeGroupSpec, b *model.VolumeGroupSpec) bool

var volumeGroupCompareFunc VolumeGroupCompareFunc

type VolumeGroupSlice []*model.VolumeGroupSpec

func (v VolumeGroupSlice) Len() int           { return len(v) }
func (v VolumeGroupSlice) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v VolumeGroupSlice) Less(i, j int) bool { return volumeGroupCompareFunc(v[i], v[j]) }

var volumeGroupSortKey2Func = map[string]VolumeGroupCompareFunc{
	"ID":   func(a *model.VolumeGroupSpec, b *model.VolumeGroupSpec) bool { return a.Id > b.Id },
	"NAME": func(a *model.VolumeGroupSpec, b *model.VolumeGroupSpec) bool { return a.Name > b.Name },
	"STATUS": func(a *model.VolumeGroupSpec, b *model.VolumeGroupSpec) bool {
		return a.Status > b.Status
	},
	"AVAILABILITYZONE": func(a *model.VolumeGroupSpec, b *model.VolumeGroupSpec) bool {
		return a.AvailabilityZone > b.AvailabilityZone
	},
	"TENANTID": func(a *model.VolumeGroupSpec, b *model.VolumeGroupSpec) bool { return a.TenantId > b.TenantId },
	"POOLID":   func(a *model.VolumeGroupSpec, b *model.VolumeGroupSpec) bool { return a.PoolId > b.PoolId },
}

func (c *Client) SortVolumeGroups(vgs []*model.VolumeGroupSpec, p *Parameter) []*model.VolumeGroupSpec {
	volumeGroupCompareFunc = volumeGroupSortKey2Func[p.sortKey]

	if strings.EqualFold(p.sortDir, "asc") {
		sort.Sort(VolumeGroupSlice(vgs))

	} else {
		sort.Sort(sort.Reverse(VolumeGroupSlice(vgs)))
	}
	return vgs
}

func (c *Client) SelectVolumeGroup(param map[string][]string, vgs []*model.VolumeGroupSpec) []*model.VolumeGroupSpec {
	if !c.SelectOrNot(param) {
		return vgs
	}

	filterList := map[string]interface{}{
		"Id":               nil,
		"CreatedAt":        nil,
		"UpdatedAt":        nil,
		"Name":             nil,
		"Status":           nil,
		"TenantId":         nil,
		"UserId":           nil,
		"Description":      nil,
		"AvailabilityZone": nil,
		"PoolId":           nil,
	}

	var vglist = []*model.VolumeGroupSpec{}
	for _, vg := range vgs {
		if c.filterByName(param, vg, filterList) {
			vglist = append(vglist, vg)
		}
	}
	return vglist
}

func (c *Client) ListHosts(ctx *c.Context, m map[string][]string) ([]*model.HostSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateHostURL(urls.Etcd, ctx.TenantId),
	}

	if IsAdminContext(ctx) {
		dbReq.Url = urls.GenerateHostURL(urls.Etcd, "")
	}

	dbRes := c.List(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When list hosts in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var hosts = []*model.HostSpec{}
	if len(dbRes.Message) == 0 {
		return hosts, nil
	}
	for _, msg := range dbRes.Message {
		var host = &model.HostSpec{}
		if err := json.Unmarshal([]byte(msg), host); err != nil {
			log.Error("When parsing host in db:", dbRes.Error)
			return nil, errors.New(dbRes.Error)
		}
		hosts = append(hosts, host)
	}

	tmpHosts := utils.Filter(hosts, m)
	if len(m["sortKey"]) > 0 && utils.Contains([]string{"hostName", "createdAt"}, m["sortKey"][0]) {
		tmpHosts = utils.Sort(tmpHosts, m["sortKey"][0], c.GetSortDir(m))
	}

	tmpHosts = utils.Slice(tmpHosts, c.GetOffset(m, len(hosts)), c.GetLimit(m))
	var res = []*model.HostSpec{}
	for _, data := range tmpHosts.([]interface{}) {
		res = append(res, data.(*model.HostSpec))
	}

	return res, nil
}

func (c *Client) ListHostsByName(ctx *c.Context, hostName string) ([]*model.HostSpec, error) {
	hosts, err := c.ListHosts(ctx, map[string][]string{"hostName": []string{hostName}})
	if err != nil {
		log.Error("List hosts failed: ", err)
		return nil, err
	}

	var res []*model.HostSpec
	for _, host := range hosts {
		if hostName == host.HostName {
			res = append(res, host)
		}
	}

	return res, nil
}

func (c *Client) CreateHost(ctx *c.Context, host *model.HostSpec) (*model.HostSpec, error) {
	host.TenantId = ctx.TenantId
	if host.Id == "" {
		host.Id = uuid.NewV4().String()
	}
	host.CreatedAt = time.Now().Format(constants.TimeFormat)
	hostBody, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}

	dbReq := &Request{
		Url:     urls.GenerateHostURL(urls.Etcd, ctx.TenantId, host.Id),
		Content: string(hostBody),
	}
	dbRes := c.Create(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When create host in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	return host, nil
}

func (c *Client) UpdateHost(ctx *c.Context, host *model.HostSpec) (*model.HostSpec, error) {
	result, err := c.GetHost(ctx, host.Id)
	if err != nil {
		return nil, err
	}
	if host.HostName != "" {
		result.HostName = host.HostName
	}
	if host.OsType != "" {
		result.OsType = host.OsType
	}
	if host.IP != "" {
		result.IP = host.IP
	}
	if host.Port > 0 {
		result.Port = host.Port
	}
	if host.AccessMode != "" {
		result.AccessMode = host.AccessMode
	}
	if host.Username != "" {
		result.Username = host.Username
	}
	if host.Password != "" {
		result.Password = host.Password
	}
	if len(host.AvailabilityZones) > 0 {
		result.AvailabilityZones = host.AvailabilityZones
	}
	if len(host.Initiators) > 0 {
		result.Initiators = host.Initiators
	}
	// Set update time
	result.UpdatedAt = time.Now().Format(constants.TimeFormat)
	body, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	// If an admin want to access other tenant's resource just fake other's tenantId.
	if !IsAdminContext(ctx) && !AuthorizeProjectContext(ctx, result.TenantId) {
		return nil, fmt.Errorf("opertaion is not permitted")
	}

	dbReq := &Request{
		Url:        urls.GenerateHostURL(urls.Etcd, result.TenantId, result.Id),
		NewContent: string(body),
	}

	dbRes := c.Update(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When update host in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return result, nil
}

func (c *Client) GetHost(ctx *c.Context, hostId string) (*model.HostSpec, error) {
	host, err := c.getHost(ctx, hostId)
	if !IsAdminContext(ctx) || err == nil {
		return host, err
	}
	hosts, err := c.ListHosts(ctx, map[string][]string{"id": []string{hostId}})
	if err != nil {
		return nil, err
	}
	for _, v := range hosts {
		if v.Id == hostId {
			return v, nil
		}
	}
	return nil, fmt.Errorf("specified host(%s) can't find", hostId)
}

func (c *Client) getHost(ctx *c.Context, hostId string) (*model.HostSpec, error) {
	dbReq := &Request{
		Url: urls.GenerateHostURL(urls.Etcd, ctx.TenantId, hostId),
	}
	dbRes := c.Get(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When get host in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}

	var host = &model.HostSpec{}
	if err := json.Unmarshal([]byte(dbRes.Message[0]), host); err != nil {
		log.Error("When parsing host in db:", dbRes.Error)
		return nil, errors.New(dbRes.Error)
	}
	return host, nil
}

func (c *Client) DeleteHost(ctx *c.Context, hostId string) error {
	// If an admin want to access other tenant's resource just fake other's tenantId.
	tenantId := ctx.TenantId
	if IsAdminContext(ctx) {
		host, err := c.GetHost(ctx, hostId)
		if err != nil {
			log.Error(err)
			return err
		}
		tenantId = host.TenantId
	}
	dbReq := &Request{
		Url: urls.GenerateHostURL(urls.Etcd, tenantId, hostId),
	}

	dbRes := c.Delete(dbReq)
	if dbRes.Status != "Success" {
		log.Error("When delete host in db:", dbRes.Error)
		return errors.New(dbRes.Error)
	}
	return nil
}
