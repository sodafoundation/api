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

/*
This module implements manila driver for OpenSDS. Manila driver will pass
these operation requests about fileshare to gophercloud which is an OpenStack
Go SDK.
*/

package manila

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/golang/glog"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	sharesv2 "github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	snapshotsv2 "github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/snapshots"
	driverConfig "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	"github.com/opensds/opensds/pkg/utils/pwd"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultConfPath = "/etc/opensds/driver/manila.yaml"
	// KManilaShareID is the UUID of the share in mannila.
	KManilaShareID = "manilaShareID"
	// KManilaSnapID is the UUID of the share snapshot in mannila.
	KManilaSnapID = "manilaSnapId"
	// KManilaShareACLID is the UUID of the share acl in mannila.
	KManilaShareACLID = "manilaAclId"
)

// Driver is a struct of manila backend.
type Driver struct {
	sharedFileSystemV2 *gophercloud.ServiceClient
	conf               *Config
}

// AuthOptions corresponds to the authentication configuration in manila.yaml
type AuthOptions struct {
	IdentityEndpoint string `yaml:"endpoint,omitempty"`
	DomainID         string `yaml:"domainId,omitempty"`
	DomainName       string `yaml:"domainName,omitempty"`
	Username         string `yaml:"username,omitempty"`
	Password         string `yaml:"password,omitempty"`
	PwdEncrypter     string `yaml:"pwdEncrypter,omitempty"`
	EnableEncrypted  bool   `yaml:"enableEncrypted,omitempty"`
	TenantID         string `yaml:"tenantId,omitempty"`
	TenantName       string `yaml:"tenantName,omitempty"`
}

// Config is a struct for parsing manila.yaml
type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]driverConfig.PoolProperties `yaml:"pool,flow"`
}

// Setup implementation
func (d *Driver) Setup() error {
	// Read manila config file
	d.conf = &Config{}
	p := config.CONF.OsdsDock.Backends.Manila.ConfigPath
	if "" == p {
		p = defaultConfPath
	}

	driverConfig.Parse(d.conf, p)
	var pwdCiphertext = d.conf.Password
	if d.conf.EnableEncrypted {
		// Decrypte the password
		pwdTool := pwd.NewPwdEncrypter(d.conf.PwdEncrypter)
		password, err := pwdTool.Decrypter(pwdCiphertext)
		if err != nil {
			return err
		}
		pwdCiphertext = password
	}

	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: d.conf.IdentityEndpoint,
		DomainID:         d.conf.DomainID,
		DomainName:       d.conf.DomainName,
		Username:         d.conf.Username,
		Password:         pwdCiphertext,
		TenantID:         d.conf.TenantID,
		TenantName:       d.conf.TenantName,
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		log.Error("openstack.AuthenticatedClient failed:", err)
		return err
	}

	d.sharedFileSystemV2, err = openstack.NewSharedFileSystemV2(provider,
		gophercloud.EndpointOpts{})
	if err != nil {
		log.Error("openstack.NewSharedFileSystemV2 failed:", err)
		return err
	}

	log.V(5).Info("setup succeeded\n")
	return nil
}

// Unset implementation
func (d *Driver) Unset() error { return nil }

// ListPools implementation
func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	// This feature is currently not implemented in gophercloud.
	// See issue: https://github.com/gophercloud/gophercloud/issues/1546
	// "Support Shared File Systems Storage Pools resource #1546"
	var pols []*model.StoragePoolSpec
	poolName := "pool1"

	pol := &model.StoragePoolSpec{
		BaseModel: &model.BaseModel{
			Id: uuid.NewV5(uuid.NamespaceOID, poolName).String(),
		},
		Name:             poolName,
		TotalCapacity:    100,
		FreeCapacity:     100,
		StorageType:      d.conf.Pool[poolName].StorageType,
		AvailabilityZone: d.conf.Pool[poolName].AvailabilityZone,
		Extras:           d.conf.Pool[poolName].Extras,
	}

	pols = append(pols, pol)
	log.V(5).Infof("function ListPools succeeded, pols:%+v\n", pols)
	return pols, nil
}

// CreateFileShare implementation
func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	prf := opt.GetProfile()
	shareProto, err := d.GetProtoFromProfile(prf)
	if err != nil {
		return nil, err
	}

	// Configure create request body.
	opts := &sharesv2.CreateOpts{
		ShareProto:  shareProto,
		Size:        int(opt.GetSize()),
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		ShareType:   "dhss_false",
		Metadata:    opt.GetMetadata(),
		// Manila's default AvailabilityZone is not "default", but ""
		//AvailabilityZone: opt.GetAvailabilityZone(),
	}

	share, err := sharesv2.Create(d.sharedFileSystemV2, opts).Extract()
	if err != nil {
		log.Errorf("cannot create share, err:%v, CreateOpts:%+v\n", err, opts)
		return nil, err
	}

	log.V(5).Infof("sharesv2.Create succeeded\n")
	// Currently dock framework doesn't support sync data from storage system,
	// therefore, it's necessary to wait for the result of resource's creation.
	// Timout after 10s.
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				tmpShare, err := d.PullFileShare(share.ID)
				if err != nil {
					continue
				}
				if tmpShare.Status != "creating" {
					share.Status = tmpShare.Status
					close(done)
					return
				}
			case <-timeout:
				close(done)
				return
			}

		}
	}()
	<-done

	d.sharedFileSystemV2.Microversion = "2.14"
	manilaExportLocations, err := sharesv2.GetExportLocations(d.sharedFileSystemV2, share.ID).Extract()
	if err != nil {
		log.Errorf("function GetExportLocations failed, err:%v", err)
		return nil, err
	}
	log.V(5).Infof("sharesv2.GetExportLocations succeeded\n")

	var exportLocations []string
	for _, v := range manilaExportLocations {
		exportLocations = append(exportLocations, v.Path)
	}

	respShare := model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Protocols:        []string{shareProto},
		Description:      opt.GetDescription(),
		Size:             opt.GetSize(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		PoolId:           opt.GetPoolId(),
		Status:           share.Status,
		Metadata:         map[string]string{KManilaShareID: share.ID},
		ExportLocations:  exportLocations,
	}

	log.V(5).Infof("function CreateFileShare succeeded, share:%+v\n", respShare)
	return &respShare, nil
}

// DeleteFileShare implementation
func (d *Driver) DeleteFileShare(opt *pb.DeleteFileShareOpts) error {
	manilaShareID := opt.Metadata[KManilaShareID]
	if err := sharesv2.Delete(d.sharedFileSystemV2, manilaShareID).ExtractErr(); err != nil {
		log.Error("cannot delete share:", err)
		return err
	}

	log.V(5).Info("function DeleteFileShare succeeded\n")
	return nil
}

// PullFileShare implementation
func (d *Driver) PullFileShare(ID string) (*model.FileShareSpec, error) {
	share, err := sharesv2.Get(d.sharedFileSystemV2, ID).Extract()
	if err != nil {
		log.Error("cannot get share:", err)
		return nil, err
	}

	respShare := model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: ID,
		},
		Name:        share.Name,
		Description: share.Description,
		Size:        int64(share.Size),
		Status:      share.Status,
	}

	log.V(5).Infof("function PullFileShare succeeded, share:%+v\n", respShare)
	return &respShare, nil
}

// CreateFileShareAcl implementation
func (d *Driver) CreateFileShareAcl(opt *pb.CreateFileShareAclOpts) (fshare *model.FileShareAclSpec, err error) {
	var accessLevel string
	accessCapability := opt.GetAccessCapability()
	canRead, canWrite, canExecute := false, false, false
	for _, v := range accessCapability {
		switch strings.ToLower(v) {
		case "read":
			canRead = true
		case "write":
			canWrite = true
		case "execute":
			canExecute = true
		default:
			return nil, errors.New("accessCapability can only be read, write or execute")
		}
	}

	switch {
	case canRead && !canWrite && !canExecute:
		accessLevel = "ro"
	case canRead && canWrite && !canExecute:
		accessLevel = "rw"
	default:
		return nil, errors.New("only read only and read write access level are supported")
	}

	// Configure request body.
	opts := &sharesv2.GrantAccessOpts{
		AccessType:  opt.Type,
		AccessTo:    opt.GetAccessTo(),
		AccessLevel: accessLevel,
	}

	d.sharedFileSystemV2.Microversion = "2.7"
	shareACL, err := sharesv2.GrantAccess(d.sharedFileSystemV2, opt.Metadata[KManilaShareID], opts).Extract()
	if err != nil {
		log.Errorf("cannot grant access, err:%v, mailaShareID:%v, opts:%+v\n", err, opt.Metadata[KManilaShareID], opts)
		return nil, err
	}

	log.V(5).Infof("sharesv2.GrantAccess succeeded\n")
	respShareACL := model.FileShareAclSpec{
		BaseModel: &model.BaseModel{
			Id: opt.Id,
		},
		FileShareId:      opt.FileshareId,
		Type:             opt.Type,
		AccessCapability: opt.GetAccessCapability(),
		AccessTo:         opt.GetAccessTo(),
		Description:      opt.Description,
		Metadata:         map[string]string{KManilaShareACLID: shareACL.ID},
	}

	log.V(5).Infof("function CreateFileShareAcl succeeded, respShareAcl:%+v\n", respShareACL)
	return &respShareACL, nil
}

// DeleteFileShareAcl implementation
func (d *Driver) DeleteFileShareAcl(opt *pb.DeleteFileShareAclOpts) error {
	opts := &sharesv2.RevokeAccessOpts{
		AccessID: opt.Metadata[KManilaShareACLID],
	}

	d.sharedFileSystemV2.Microversion = "2.7"
	if err := sharesv2.RevokeAccess(d.sharedFileSystemV2, opt.Metadata[KManilaShareID], opts).ExtractErr(); err != nil {
		log.Error("cannot revoke access:", err)
		return err
	}

	log.V(5).Info("function DeleteFileShareAcl succeeded\n")
	return nil
}

// CreateFileShareSnapshot implementation
func (d *Driver) CreateFileShareSnapshot(opt *pb.CreateFileShareSnapshotOpts) (*model.FileShareSnapshotSpec, error) {
	mailaShareID := opt.Metadata[KManilaShareID]
	opts := &snapshotsv2.CreateOpts{
		ShareID:            mailaShareID,
		Name:               opt.GetName(),
		Description:        opt.GetDescription(),
		DisplayName:        "",
		DisplayDescription: "",
	}

	snapshot, err := snapshotsv2.Create(d.sharedFileSystemV2, opts).Extract()
	if err != nil {
		log.Errorf("cannot create snapshot, err:%v, CreateOpts:%+v\n", err, opts)
		return nil, err
	}

	// Currently dock framework doesn't support sync data from storage system,
	// therefore, it's necessary to wait for the result of resource's creation.
	// Timout after 10s.
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				tmpSnapshot, err := d.PullFileShareSnapshot(snapshot.ID)
				if err != nil {
					continue
				}
				if tmpSnapshot.Status != "creating" {
					snapshot.Status = tmpSnapshot.Status
					close(done)
					return
				}
			case <-timeout:
				close(done)
				return
			}

		}
	}()
	<-done

	respSnapshot := model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:         opt.GetName(),
		Description:  opt.GetDescription(),
		SnapshotSize: int64(snapshot.Size),
		Status:       snapshot.Status,
		Metadata:     map[string]string{KManilaSnapID: snapshot.ID},
	}

	log.V(5).Infof("function CreateFileShareSnapshot succeeded, snapshot:%+v\n", respSnapshot)
	return &respSnapshot, nil
}

// DeleteFileShareSnapshot implementation
func (d *Driver) DeleteFileShareSnapshot(opt *pb.DeleteFileShareSnapshotOpts) error {
	manilaSnapID := opt.Metadata[KManilaSnapID]
	if err := snapshotsv2.Delete(d.sharedFileSystemV2, manilaSnapID).ExtractErr(); err != nil {
		log.Error("cannot delete share:", err)
		return err
	}

	log.V(5).Info("function DeleteFileShareSnapshot succeeded\n")
	return nil
}

// PullFileShareSnapshot implementation
func (d *Driver) PullFileShareSnapshot(ID string) (*model.FileShareSnapshotSpec, error) {
	snapshot, err := snapshotsv2.Get(d.sharedFileSystemV2, ID).Extract()
	if err != nil {
		log.Error("cannot get snapshot:", err)
		return nil, err
	}

	respShareSnap := model.FileShareSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: ID,
		},
		Name:         snapshot.Name,
		Description:  snapshot.Description,
		SnapshotSize: int64(snapshot.Size),
		Status:       snapshot.Status,
	}

	log.V(5).Infof("function PullFileShareSnapshot succeeded, snapshot:%+v\n", respShareSnap)
	return &respShareSnap, nil
}

// GetProtoFromProfile implementation
func (d *Driver) GetProtoFromProfile(prf string) (string, error) {
	if prf == "" {
		msg := "profile cannot be empty"
		return "", errors.New(msg)
	}

	log.V(5).Infof("file share profile is %s", prf)
	profile := &model.ProfileSpec{}
	err := json.Unmarshal([]byte(prf), profile)
	if err != nil {
		msg := fmt.Sprintf("unmarshal profile failed: %v", err)
		return "", errors.New(msg)
	}

	shareProto := profile.ProvisioningProperties.IOConnectivity.AccessProtocol
	if shareProto == "" {
		msg := "file share protocol cannot be empty"
		return "", errors.New(msg)
	}

	return shareProto, nil
}
