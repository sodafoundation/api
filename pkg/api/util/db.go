// Copyright 2018 The OpenSDS Authors.
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
This module implements API-related database operations.
*/

package util

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/golang/glog"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
	uuid "github.com/satori/go.uuid"
)

//function to store filesahreAcl metadata into database
func CreateFileShareAclDBEntry(ctx *c.Context, in *model.FileShareAclSpec) (*model.FileShareAclSpec, error) {
	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}

	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	if in.UpdatedAt == "" {
		in.UpdatedAt = time.Now().Format(constants.TimeFormat)
	}
	// validate type
	if in.Type != "ip" {
		errMsg := fmt.Sprintf("invalid fileshare type: %v. Supported type is: ip", in.Type)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	// validate accessTo
	accessto := in.AccessTo
	if in.AccessTo == "" {
		errMsg := fmt.Sprintf("accessTo is empty. Please give valid ip segment")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	} else if strings.Contains(accessto, "/") {
		first, cidr, bool := net.ParseCIDR(accessto)
		log.Info(first, cidr)
		if bool != nil {
			errMsg := fmt.Sprintf("invalid IP segment %v", accessto)
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
	} else {
		server := net.ParseIP(in.AccessTo)
		if server == nil {
			errMsg := fmt.Sprintf("%v is not a valid ip. Please give the proper ip", in.AccessTo)
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	// validate accesscapability
	accessCapability := in.AccessCapability
	if len(accessCapability) == 0 {
		errMsg := fmt.Sprintf("empty fileshare accesscapability. Supported accesscapability are: {read, write}")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	permissions := []string{"write", "read"}
	for _, value := range accessCapability {
		value = strings.ToLower(value)
		if !(utils.Contains(permissions, value)) {
			errMsg := fmt.Sprintf("invalid fileshare accesscapability: %v. Supported accesscapability are: {read, write}", value)
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	// get fileshare details
	fileshare, err := db.C.GetFileShare(ctx, in.FileShareId)
	if err != nil {
		log.Error("file shareid is not valid: ", err)
		return nil, err
	}

	if fileshare.Status != model.FileShareAvailable {
		var errMsg = "only the status of file share is available, the acl can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	// Store the fileshare meadata into database.
	return db.C.CreateFileShareAcl(ctx, in)
}

func DeleteFileShareAclDBEntry(ctx *c.Context, in *model.FileShareAclSpec) error {
	// If fileshare id is invalid, it would mean that fileshare snapshot
	// creation failed before the create method in storage driver was
	// called, and delete its db entry directly.
	if _, err := db.C.GetFileShare(ctx, in.FileShareId); err != nil {
		if err := db.C.DeleteFileShareAcl(ctx, in.Id); err != nil {
			log.Error("when delete fileshare acl in db:", err)
			return err
		}
	}
	return nil
}

// Function to store metadeta of fileshare into database
func CreateFileShareDBEntry(ctx *c.Context, in *model.FileShareSpec) (*model.FileShareSpec, error) {
	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	// validate the size
	if in.Size <= 0 {
		errMsg := fmt.Sprintf("invalid fileshare size: %d", in.Size)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if in.AvailabilityZone == "" {
		log.Warning("Use default availability zone when user doesn't specify availabilityZone.")
		in.AvailabilityZone = "default"
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	if in.UpdatedAt == "" {
		in.UpdatedAt = time.Now().Format(constants.TimeFormat)
	}
	//validate the name
	if in.Name == "" {
		errMsg := fmt.Sprintf("empty fileshare name is not allowed. Please give valid name.")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if len(in.Name) > 128 {
		errMsg := fmt.Sprintf("fileshare name length should not more than 128 characters. current length is : %d", len(in.Name))
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	in.UserId = ctx.UserId
	in.Status = model.FileShareCreating
	// Store the fileshare meadata into database.
	return db.C.CreateFileShare(ctx, in)
}

// DeleteFileShareDBEntry just modifies the state of the fileshare to be deleting in
// the DB, the real deletion operation would be executed in another new thread.
func DeleteFileShareDBEntry(ctx *c.Context, in *model.FileShareSpec) error {
	validStatus := []string{model.FileShareAvailable, model.FileShareError,
		model.FileShareErrorDeleting}
	if !utils.Contained(in.Status, validStatus) {
		errMsg := fmt.Sprintf("only the fileshare with the status available, error, errorDeleting, can be deleted, the fileshare status is %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	snaps, err := db.C.ListSnapshotsByShareId(ctx, in.Id)
	if err != nil {
		return err
	}
	if len(snaps) > 0 {
		errMsg := fmt.Sprintf("file share %s can not be deleted, because it still has snapshots", in.Id)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	acls, err := db.C.ListFileShareAclsByShareId(ctx, in.Id)
	if err != nil {
		return err
	}
	if len(acls) > 0 {
		errMsg := fmt.Sprintf("file share %s can not be deleted, because it still has acls", in.Id)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.Status = model.FileShareDeleting
	_, err = db.C.UpdateFileShare(ctx, in)
	if err != nil {
		return err
	}
	return nil
}

// To create entry in database
func CreateFileShareSnapshotDBEntry(ctx *c.Context, in *model.FileShareSnapshotSpec) (*model.FileShareSnapshotSpec, error) {
	fshare, err := db.C.GetFileShare(ctx, in.FileShareId)
	if err != nil {
		log.Error("get fileshare failed in create fileshare snapshot method: ", err)
		return nil, err
	}
	if fshare.Status != model.FileShareAvailable && fshare.Status != model.FileShareInUse {
		var errMsg = "only the status of fileshare is available or in-use, the snapshot can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	//validate the snapshot name
	if in.Name == "" {
		errMsg := fmt.Sprintf("snapshot name can not be empty. Please give valid snapshot name")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if strings.HasPrefix(in.Name, "snapshot") {
		errMsg := fmt.Sprintf("names starting 'snapshot' are reserved. Please choose a different snapshot name.")
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	in.Status = model.FileShareSnapCreating
	in.Metadata = fshare.Metadata
	return db.C.CreateFileShareSnapshot(ctx, in)
}

func DeleteFileShareSnapshotDBEntry(ctx *c.Context, in *model.FileShareSnapshotSpec) error {
	validStatus := []string{model.FileShareSnapAvailable, model.FileShareSnapError,
		model.FileShareSnapErrorDeleting}
	if !utils.Contained(in.Status, validStatus) {
		errMsg := fmt.Sprintf("only the fileshare snapshot with the status available, error, error_deleting can be deleted, the fileshare status is %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	// If fileshare id is invalid, it would mean that fileshare snapshot creation failed before the create method
	// in storage driver was called, and delete its db entry directly.
	_, err := db.C.GetFileShare(ctx, in.FileShareId)
	if err != nil {
		if err := db.C.DeleteFileShareSnapshot(ctx, in.Id); err != nil {
			log.Error("when delete fileshare snapshot in db:", err)
			return err
		}
		return nil
	}

	in.Status = model.FileShareSnapDeleting
	_, err = db.C.UpdateFileShareSnapshot(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

func CreateVolumeDBEntry(ctx *c.Context, in *model.VolumeSpec) (*model.VolumeSpec, error) {
	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.Size <= 0 {
		errMsg := fmt.Sprintf("invalid volume size: %d", in.Size)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if in.SnapshotId != "" {
		snap, err := db.C.GetVolumeSnapshot(ctx, in.SnapshotId)
		if err != nil {
			log.Error("get snapshot failed in create volume method: ", err)
			return nil, err
		}
		if snap.Status != model.VolumeSnapAvailable {
			var errMsg = "only if the snapshot is available, the volume can be created"
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		if snap.Size > in.Size {
			var errMsg = "size of volume must be equal to or bigger than size of the snapshot"
			log.Error(errMsg)
			return nil, errors.New(errMsg)
		}
	}
	if in.AvailabilityZone == "" {
		log.Warning("Use default availability zone when user doesn't specify availabilityZone.")
		in.AvailabilityZone = "default"
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	in.UserId = ctx.UserId
	in.Status = model.VolumeCreating
	// Store the volume data into database.
	return db.C.CreateVolume(ctx, in)
}

// DeleteVolumeDBEntry just modifies the state of the volume to be deleting in
// the DB, the real deletion operation would be executed in another new thread.
func DeleteVolumeDBEntry(ctx *c.Context, in *model.VolumeSpec) error {
	validStatus := []string{model.VolumeAvailable, model.VolumeError,
		model.VolumeErrorDeleting, model.VolumeErrorExtending}
	if !utils.Contained(in.Status, validStatus) {
		errMsg := fmt.Sprintf("only the volume with the status available, error, error_deleting, error_extending can be deleted, the volume status is %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	snaps, err := db.C.ListSnapshotsByVolumeId(ctx, in.Id)
	if err != nil {
		return err
	}
	if len(snaps) > 0 {
		return fmt.Errorf("volume %s can not be deleted, because it still has snapshots", in.Id)
	}

	volAttachments, err := db.C.ListAttachmentsByVolumeId(ctx, in.Id)
	if err != nil {
		return err
	}
	if len(volAttachments) > 0 {
		return fmt.Errorf("volume %s can not be deleted, because it's in use", in.Id)
	}

	in.Status = model.VolumeDeleting
	_, err = db.C.UpdateVolume(ctx, in)
	if err != nil {
		return err
	}
	return nil
}

// ExtendVolumeDBEntry just modifies the state of the volume to be extending in
// the DB, the real operation would be executed in another new thread, and the
// new size would be updated in controller module.
func ExtendVolumeDBEntry(ctx *c.Context, volID string, in *model.ExtendVolumeSpec) (*model.VolumeSpec, error) {
	volume, err := db.C.GetVolume(ctx, volID)
	if err != nil {
		log.Error("get volume failed in extend volume method: ", err)
		return nil, err
	}

	if volume.Status != model.VolumeAvailable {
		errMsg := "the status of the volume to be extended must be available!"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	if in.NewSize <= volume.Size {
		errMsg := fmt.Sprintf("new size for extend must be greater than current size."+
			"(current: %d GB, extended: %d GB).", volume.Size, in.NewSize)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	volume.Status = model.VolumeExtending
	// Store the volume data into database.
	return db.C.ExtendVolume(ctx, volume)
}

// CreateVolumeAttachmentDBEntry just modifies the state of the volume attachment
// to be creating in the DB, the real operation would be executed in another new
// thread.
func CreateVolumeAttachmentDBEntry(ctx *c.Context, volAttachment *model.VolumeAttachmentSpec) (
	*model.VolumeAttachmentSpec, error) {
	vol, err := db.C.GetVolume(ctx, volAttachment.VolumeId)
	if err != nil {
		msg := fmt.Sprintf("get volume failed in create volume attachment method: %v", err)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	if vol.Status == model.VolumeAvailable {
		db.UpdateVolumeStatus(ctx, db.C, vol.Id, model.VolumeAttaching)
	} else if vol.Status == model.VolumeInUse {
		if vol.MultiAttach {
			db.UpdateVolumeStatus(ctx, db.C, vol.Id, model.VolumeAttaching)
		} else {
			msg := "volume is already attached or volume multiattach must be true if attach more than once"
			log.Error(msg)
			return nil, errors.New(msg)
		}
	} else {
		errMsg := "only the status of volume is available, attachment can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if volAttachment.Id == "" {
		volAttachment.Id = uuid.NewV4().String()
	}

	if volAttachment.CreatedAt == "" {
		volAttachment.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	if volAttachment.AttachMode != "ro" && volAttachment.AttachMode != "rw" {
		volAttachment.AttachMode = "rw"
	}

	volAttachment.Status = model.VolumeAttachCreating
	volAttachment.Metadata = utils.MergeStringMaps(volAttachment.Metadata, vol.Metadata)
	return db.C.CreateVolumeAttachment(ctx, volAttachment)
}

// DeleteVolumeAttachmentDBEntry just modifies the state of the volume attachment to
// be deleting in the DB, the real deletion operation would be executed in
// another new thread.
func DeleteVolumeAttachmentDBEntry(ctx *c.Context, in *model.VolumeAttachmentSpec) error {
	validStatus := []string{model.VolumeAttachAvailable, model.VolumeAttachError,
		model.VolumeAttachErrorDeleting}
	if !utils.Contained(in.Status, validStatus) {
		errMsg := fmt.Sprintf("only the volume attachment with the status available, error, error_deleting can be deleted, the volume status is %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	// If volume id is invalid, it would mean that volume attachment creation failed before the create method
	// in storage driver was called, and delete its db entry directly.
	_, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		if err := db.C.DeleteVolumeAttachment(ctx, in.Id); err != nil {
			log.Error("when delete volume attachment in db:", err)
			return err
		}
		return nil
	}

	in.Status = model.VolumeAttachDeleting
	_, err = db.C.UpdateVolumeAttachment(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

// CreateVolumeSnapshotDBEntry just modifies the state of the volume snapshot
// to be creating in the DB, the real operation would be executed in another new
// thread.
func CreateVolumeSnapshotDBEntry(ctx *c.Context, in *model.VolumeSnapshotSpec) (*model.VolumeSnapshotSpec, error) {
	vol, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		log.Error("get volume failed in create volume snapshot method: ", err)
		return nil, err
	}
	if vol.Status != model.VolumeAvailable && vol.Status != model.VolumeInUse {
		var errMsg = "only the status of volume is available or in-use, the snapshot can be created"
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	in.Status = model.VolumeSnapCreating
	return db.C.CreateVolumeSnapshot(ctx, in)
}

// DeleteVolumeSnapshotDBEntry just modifies the state of the volume snapshot to
// be deleting in the DB, the real deletion operation would be executed in
// another new thread.
func DeleteVolumeSnapshotDBEntry(ctx *c.Context, in *model.VolumeSnapshotSpec) error {
	validStatus := []string{model.VolumeSnapAvailable, model.VolumeSnapError,
		model.VolumeSnapErrorDeleting}
	if !utils.Contained(in.Status, validStatus) {
		errMsg := fmt.Sprintf("only the volume snapshot with the status available, error, error_deleting can be deleted, the volume status is %s", in.Status)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	// If volume id is invalid, it would mean that volume snapshot creation failed before the create method
	// in storage driver was called, and delete its db entry directly.
	_, err := db.C.GetVolume(ctx, in.VolumeId)
	if err != nil {
		if err := db.C.DeleteVolumeSnapshot(ctx, in.Id); err != nil {
			log.Error("when delete volume snapshot in db:", err)
			return err
		}
		return nil
	}

	in.Status = model.VolumeSnapDeleting
	_, err = db.C.UpdateVolumeSnapshot(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

// CreateReplicationDBEntry just modifies the state of the volume replication
// to be creating in the DB, the real deletion operation would be executed
// in another new thread.
func CreateReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) (*model.ReplicationSpec, error) {
	pVol, err := db.C.GetVolume(ctx, in.PrimaryVolumeId)
	if err != nil {
		log.Error("get primary volume failed in create volume replication method: ", err)
		return nil, err
	}
	if pVol.Status != model.VolumeAvailable && pVol.Status != model.VolumeInUse {
		var errMsg = fmt.Errorf("only the status of primary volume is available or in-use, the replicaiton can be created")
		log.Error(errMsg)
		return nil, errMsg
	}
	sVol, err := db.C.GetVolume(ctx, in.SecondaryVolumeId)
	if err != nil {
		log.Error("get secondary volume failed in create volume replication method: ", err)
		return nil, err
	}
	if sVol.Status != model.VolumeAvailable && sVol.Status != model.VolumeInUse {
		var errMsg = fmt.Errorf("only the status of secondary volume is available or in-use, the replicaiton can be created")
		log.Error(errMsg)
		return nil, errMsg
	}

	// Check if specified volume has already been used in other replication.
	v, err := db.C.GetReplicationByVolumeId(ctx, in.PrimaryVolumeId)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); !ok {
			var errMsg = fmt.Errorf("get replication by primary volume id %s failed: %v",
				in.PrimaryVolumeId, err)
			log.Error(errMsg)
			return nil, errMsg
		}
	}

	if v != nil {
		var errMsg = fmt.Errorf("specified primary volume(%s) has already been used in replication(%s)",
			in.PrimaryVolumeId, v.Id)
		log.Error(errMsg)
		return nil, errMsg
	}

	// check if specified volume has already been used in other replication.
	v, err = db.C.GetReplicationByVolumeId(ctx, in.SecondaryVolumeId)
	if err != nil {
		if _, ok := err.(*model.NotFoundError); !ok {
			var errMsg = fmt.Errorf("get replication by secondary volume id %s failed: %v",
				in.SecondaryVolumeId, err)
			log.Error(errMsg)
			return nil, errMsg
		}
	}
	if v != nil {
		var errMsg = fmt.Errorf("specified secondary volume(%s) has already been used in replication(%s)",
			in.SecondaryVolumeId, v.Id)
		log.Error(errMsg)
		return nil, errMsg
	}

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}

	in.ReplicationStatus = model.ReplicationCreating
	return db.C.CreateReplication(ctx, in)
}

// DeleteReplicationDBEntry just modifies the state of the volume replication to
// be deleting in the DB, the real deletion operation would be executed in
// another new thread.
func DeleteReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}

	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("can't delete the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.ReplicationStatus = model.ReplicationDeleting
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

// EnableReplicationDBEntry just modifies the state of the volume replication to
// be enabling in the DB, the real deletion operation would be executed in
// another new thread.
func EnableReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}
	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("can't enable the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.ReplicationStatus = model.ReplicationEnabling
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

// DisableReplicationDBEntry just modifies the state of the volume replication to
// be disabling in the DB, the real deletion operation would be executed in
// another new thread.
func DisableReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}
	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("can't disable the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	in.ReplicationStatus = model.ReplicationDisabling
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

// FailoverReplicationDBEntry just modifies the state of the volume replication
// to be failing_over or failing_back in the DB, the real deletion operation
// would be executed in another new thread.
func FailoverReplicationDBEntry(ctx *c.Context, in *model.ReplicationSpec, secondaryBackendId string) error {
	invalidStatus := []string{model.ReplicationCreating, model.ReplicationDeleting, model.ReplicationEnabling,
		model.ReplicationDisabling, model.ReplicationFailingOver, model.ReplicationFailingBack}
	if utils.Contained(in.ReplicationStatus, invalidStatus) {
		errMsg := fmt.Sprintf("can't fail over/back the replication in %s", in.ReplicationStatus)
		log.Error(errMsg)
		return errors.New(errMsg)
	}

	if secondaryBackendId == model.ReplicationDefaultBackendId {
		in.ReplicationStatus = model.ReplicationFailingOver
	} else {
		in.ReplicationStatus = model.ReplicationFailingBack
	}
	_, err := db.C.UpdateReplication(ctx, in.Id, in)
	if err != nil {
		return err
	}
	return nil
}

// CreateVolumeGroupDBEntry just modifies the state of the volume group
// to be creating in the DB, the real deletion operation would be
// executed in another new thread.
func CreateVolumeGroupDBEntry(ctx *c.Context, in *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	if len(in.Profiles) == 0 {
		msg := fmt.Sprintf("profiles must be provided to create volume group.")
		log.Error(msg)
		return nil, errors.New(msg)
	}

	if in.Id == "" {
		in.Id = uuid.NewV4().String()
	}
	if in.CreatedAt == "" {
		in.CreatedAt = time.Now().Format(constants.TimeFormat)
	}
	if in.AvailabilityZone == "" {
		log.Warning("Use default availability zone when user doesn't specify availabilityZone.")
		in.AvailabilityZone = "default"
	}

	in.Status = model.VolumeGroupCreating
	return db.C.CreateVolumeGroup(ctx, in)
}

// UpdateVolumeGroupDBEntry just modifies the state of the volume group
// to be updating in the DB, the real deletion operation would be
// executed in another new thread.
func UpdateVolumeGroupDBEntry(ctx *c.Context, vgUpdate *model.VolumeGroupSpec) (*model.VolumeGroupSpec, error) {
	vg, err := db.C.GetVolumeGroup(ctx, vgUpdate.Id)
	if err != nil {
		return nil, err
	}

	var name string
	if vg.Name == vgUpdate.Name {
		name = ""
	} else {
		name = vgUpdate.Name
	}
	var description string
	if vg.Description == vgUpdate.Description {
		description = ""
	} else {
		description = vgUpdate.Description
	}
	vgUpdate.Profiles = vg.Profiles
	vgUpdate.PoolId = vg.PoolId

	var invalidUuids []string
	for _, uuidAdd := range vgUpdate.AddVolumes {
		for _, uuidRemove := range vgUpdate.RemoveVolumes {
			if uuidAdd == uuidRemove {
				invalidUuids = append(invalidUuids, uuidAdd)
			}
		}
	}
	if len(invalidUuids) > 0 {
		msg := fmt.Sprintf("uuid %s is in both add and remove volume list", strings.Join(invalidUuids, ","))
		log.Error(msg)
		return nil, errors.New(msg)
	}

	volumes, err := db.C.ListVolumesByGroupId(ctx, vgUpdate.Id)
	if err != nil {
		return nil, err
	}

	var addVolumesNew, removeVolumeNew []string
	// Validate volumes in AddVolumes and RemoveVolumes.
	if len(vgUpdate.AddVolumes) > 0 {
		if addVolumesNew, err = ValidateAddVolumes(ctx, volumes, vgUpdate.AddVolumes, vgUpdate); err != nil {
			return nil, err
		}
	}
	if len(vgUpdate.RemoveVolumes) > 0 {
		if removeVolumeNew, err = ValidateRemoveVolumes(ctx, volumes, vgUpdate.RemoveVolumes, vgUpdate); err != nil {
			return nil, err
		}
	}

	if name == "" && description == "" && len(addVolumesNew) == 0 && len(removeVolumeNew) == 0 {
		msg := fmt.Sprintf("update group %s faild, because no valid name, description, addvolumes or removevolumes were provided", vgUpdate.Id)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	vgNew := &model.VolumeGroupSpec{
		BaseModel: &model.BaseModel{
			Id: vg.Id,
		},
	}

	vgNew.UpdatedAt = time.Now().Format(constants.TimeFormat)
	// Only update name or description. No need to send them over through an RPC call and set status to available.
	if name != "" {
		vgNew.Name = name
	}
	if description != "" {
		vgNew.Description = description
	}
	if len(addVolumesNew) == 0 && len(removeVolumeNew) == 0 {
		vgNew.Status = model.VolumeGroupAvailable
	} else {
		vgNew.Status = model.VolumeGroupUpdating
	}

	return db.C.UpdateVolumeGroup(ctx, vgNew)
}

func ValidateAddVolumes(ctx *c.Context, volumes []*model.VolumeSpec, addVolumes []string, vg *model.VolumeGroupSpec) ([]string, error) {
	var addVolumeRef []string
	var flag bool
	for _, volumeId := range addVolumes {
		flag = true
		for _, volume := range volumes {
			if volumeId == volume.Id {
				// Volume already in group. Remove it from addVolumes.
				flag = false
				break
			}
		}
		if flag {
			addVolumeRef = append(addVolumeRef, volumeId)
		}
	}

	var addVolumesNew []string
	for _, addVol := range addVolumeRef {
		addVolRef, err := db.C.GetVolume(ctx, addVol)
		if err != nil {
			log.Error(fmt.Sprintf("cannot add volume %s to group %s, volume cannot be found.", addVol, vg.Id))
			return nil, err
		}
		if addVolRef.GroupId != "" {
			return nil, fmt.Errorf("cannot add volume %s to group %s because it is already in group %s", addVolRef.Id, vg.Id, addVolRef.GroupId)
		}
		if addVolRef.ProfileId == "" {
			return nil, fmt.Errorf("cannot add volume %s to group %s , volume has no profile.", addVolRef.Id, vg.Id)
		}
		if !utils.Contained(addVolRef.ProfileId, vg.Profiles) {
			return nil, fmt.Errorf("cannot add volume %s to group %s , volume profile is not supported by the group.", addVolRef.Id, vg.Id)
		}
		if addVolRef.Status != model.VolumeAvailable && addVolRef.Status != model.VolumeInUse {
			return nil, fmt.Errorf("cannot add volume %s to group %s because volume is in invalid status %s", addVolRef.Id, vg.Id, addVolRef.Status)
		}
		if addVolRef.PoolId != vg.PoolId {
			return nil, fmt.Errorf("cannot add volume %s to group %s , volume is not local to the pool of group.", addVolRef.Id, vg.Id)
		}

		addVolumesNew = append(addVolumesNew, addVolRef.Id)
	}

	return addVolumesNew, nil
}

func ValidateRemoveVolumes(ctx *c.Context, volumes []*model.VolumeSpec, removeVolumes []string, vg *model.VolumeGroupSpec) ([]string, error) {

	for _, v := range removeVolumes {
		for _, volume := range volumes {
			if v == volume.Id {
				if volume.Status != model.VolumeAvailable && volume.Status != model.VolumeInUse && volume.Status != model.VolumeError && volume.Status != model.VolumeErrorDeleting {
					return nil, fmt.Errorf("cannot remove volume %s from group %s, volume is in invalid status %s", volume.Id, vg.Id, volume.Status)
				}
				break
			}

		}
	}
	for _, v := range removeVolumes {
		var available = false
		for _, volume := range volumes {
			if v == volume.Id {
				available = true
				break
			}
		}
		if available == false {
			return nil, fmt.Errorf("cannot remove volume %s from group %s, volume is not in group ", v, vg.Id)
		}
	}

	return removeVolumes, nil
}

// DeleteVolumeGroupDBEntry just modifies the state of the volume group
// to be deleting in the DB, the real deletion operation would be
// executed in another new thread.
func DeleteVolumeGroupDBEntry(ctx *c.Context, volumeGroupId string) error {
	vg, err := db.C.GetVolumeGroup(ctx, volumeGroupId)
	if err != nil {
		return err
	}

	// If pool id is invalid, it would mean that volume group creation failed before the create method
	// in storage driver was called, and delete its db entry directly.
	_, err = db.C.GetDockByPoolId(ctx, vg.PoolId)
	if err != nil {
		if err := db.C.DeleteVolumeGroup(ctx, vg.Id); err != nil {
			log.Error("when delete volume group in db:", err)
			return err
		}
		return nil
	}

	//TODO DeleteVolumes tag is set by policy.
	deleteVolumes := true

	if deleteVolumes == false && vg.Status != model.VolumeGroupAvailable && vg.Status != model.VolumeGroupError {
		msg := fmt.Sprintf("the status of the Group must be available or error , group can be deleted. But current status is %s", vg.Status)
		log.Error(msg)
		return errors.New(msg)
	}

	if vg.GroupSnapshots != nil {
		msg := fmt.Sprintf("group can not be deleted, because group has existing snapshots")
		log.Error(msg)
		return errors.New(msg)
	}

	volumes, err := db.C.ListVolumesByGroupId(ctx, vg.Id)
	if err != nil {
		return err
	}

	if len(volumes) > 0 && deleteVolumes == false {
		msg := fmt.Sprintf("group %s still contains volumes. The deleteVolumes flag is required to delete it.", vg.Id)
		log.Error(msg)
		return errors.New(msg)
	}

	var volumesUpdate []*model.VolumeSpec
	for _, value := range volumes {
		if value.AttachStatus == model.VolumeAttached {
			msg := fmt.Sprintf("volume %s in group %s is attached. Need to deach first.", value.Id, vg.Id)
			log.Error(msg)
			return errors.New(msg)
		}

		snapshots, err := db.C.ListSnapshotsByVolumeId(ctx, value.Id)
		if err != nil {
			return err
		}
		if len(snapshots) > 0 {
			msg := fmt.Sprintf("volume %s in group still has snapshots", value.Id)
			log.Error(msg)
			return errors.New(msg)
		}

		volumesUpdate = append(volumesUpdate, &model.VolumeSpec{
			BaseModel: &model.BaseModel{
				Id: value.Id,
			},
			Status:  model.VolumeDeleting,
			GroupId: volumeGroupId,
		})
	}

	db.C.UpdateStatus(ctx, volumesUpdate, "")
	db.C.UpdateStatus(ctx, vg, model.VolumeGroupDeleting)

	return nil
}
