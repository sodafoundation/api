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
	"fmt"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	uuid "github.com/satori/go.uuid"
)

const (
	DefaultConfPath = "/etc/opensds/driver/hpe_nimble.yaml"
	NamePrefix      = "opensds"
)

type Config struct {
	AuthOptions `yaml:"authOptions"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
}

type Driver struct {
	conf   *Config
	client *NimbleClient
}

func (d *Driver) Setup() (err error) {

	conf := &Config{}
	d.conf = conf
	path := config.CONF.OsdsDock.Backends.HpeNimble.ConfigPath
	if "" == path {
		path = DefaultConfPath
	}
	Parse(conf, path)

	d.client, err = NewClient(&d.conf.AuthOptions)
	if err != nil {
		log.Errorf("%v: get new client failed.\n%v", DriverName, err)
		return err
	}

	return nil
}

func (d *Driver) Unset() error {
	return nil
}

func (d *Driver) CreateVolume(opt *pb.CreateVolumeOpts) (*model.VolumeSpec, error) {
	log.Infof("%v: try to create volume...", DriverName)

	poolId, err := d.client.GetPoolIdByName(opt.GetPoolName())
	if err != nil {
		return nil, err
	}
	lun, err := d.client.CreateVolume(poolId, opt)
	if err != nil {
		log.Errorf("%v: create Volume Failed: %v", DriverName, err)
		return nil, err
	}
	log.Infof("%v: create volume success.", DriverName)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             Byte2Gib(lun.Size),
		Description:      lun.Description,
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			"Group":  lun.OwnedByGroup,
			"Iqn":    lun.TargetName,
			"LunId":  lun.Id,
			"PoolId": poolId,
		},
	}, nil
}

func (d *Driver) DeleteVolume(opt *pb.DeleteVolumeOpts) error {
	log.Infof("%v: Trying delete volume ...", DriverName)
	poolId := opt.GetMetadata()["PoolId"]
	err := d.client.DeleteVolume(poolId, opt)
	if err != nil {
		log.Errorf("%v: delete volume failed, volume id =%s , error:%s", DriverName, opt.GetId(), err)
		return err
	}
	log.Infof("%v: remove volume success", DriverName)
	return nil
}

func (d *Driver) ExtendVolume(opt *pb.ExtendVolumeOpts) (*model.VolumeSpec, error) {
	log.Infof("%v: trying Extend volume...", DriverName)
	poolId := opt.GetMetadata()["PoolId"]
	_, err := d.client.ExtendVolume(poolId, opt)
	if err != nil {
		log.Errorf("%v: extend Volume Failed:", DriverName)
		return nil, err
	}

	log.Infof("%v: extend volume success.", DriverName)
	return &model.VolumeSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
	}, nil
}

func (d *Driver) CreateSnapshot(opt *pb.CreateVolumeSnapshotOpts) (*model.VolumeSnapshotSpec, error) {
	log.Infof("%v: trying create snapshot...", DriverName)
	poolId := opt.GetMetadata()["PoolId"]
	snap, err := d.client.CreateSnapshot(poolId, opt)
	if err != nil {
		return nil, err
	}

	log.Infof("%v: create snapshot success.", DriverName)
	return &model.VolumeSnapshotSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:        opt.GetName(),
		Description: opt.GetDescription(),
		VolumeId:    opt.GetVolumeId(),
		Size:        0,
		Metadata: map[string]string{
			"SnapId": snap.Id,
		},
	}, nil
}

func (d *Driver) DeleteSnapshot(opt *pb.DeleteVolumeSnapshotOpts) error {
	log.Infof("%v: trying delete snapshot...", DriverName)
	poolId := opt.GetMetadata()["PoolId"]
	err := d.client.DeleteSnapshot(poolId, opt)
	if err != nil {
		log.Errorf("%v: delete volume snapshot failed, volume snapshot id = %s , error: %v", DriverName, opt.GetId(), err)
		return err
	}
	log.Infof("%v: remove volume snapshot success, volume snapshot id=%v", DriverName, opt.GetId())
	return nil
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	log.Infof("%v: listPools ...", DriverName)
	var pols []*model.StoragePoolSpec
	sp, err := d.client.ListStoragePools()
	if err != nil {
		return nil, err
	}

	c := d.conf

	for _, pool := range sp {
		for grpName, _ := range c.Pool {
			if grpName == pool.Name && c.Pool[grpName].AvailabilityZone == pool.ArrayList[0].ArrayName {

				pol := &model.StoragePoolSpec{
					BaseModel: &model.BaseModel{
						Id: uuid.NewV5(uuid.NamespaceOID, pool.Id).String(),
					},
					Name:             pool.Name,
					Description:      pool.Description,
					TotalCapacity:    Byte2Gib(pool.TotalCapacity),
					FreeCapacity:     Byte2Gib(pool.FreeCapacity),
					StorageType:      "block",
					AvailabilityZone: pool.ArrayList[0].ArrayName + "/" + pool.Name,
					Extras:           c.Pool[grpName].Extras,
				}
				pols = append(pols, pol)
				break
			}
		}
	}

	// Error if there is NO valid storage grp
	if len(pols) == 0 {
		return nil, fmt.Errorf("%v: there are no valid storage pool. Pls check driver config.\n", DriverName)
	}

	return pols, nil
}

func (d *Driver) InitializeConnection(opt *pb.CreateVolumeAttachmentOpts) (*model.ConnectionInfo, error) {

	if opt.GetAccessProtocol() == ISCSIProtocol || opt.GetAccessProtocol() == FCProtocol {
		log.Infof("%v: trying initialize connection...", DriverName)

		poolId := opt.GetMetadata()["PoolId"]
		initiatorName := opt.GetHostInfo().Initiator
		if initiatorName == "" || opt.GetHostInfo().Ip == "" {
			if opt.GetAccessProtocol() == ISCSIProtocol {
				return nil, fmt.Errorf("%v: pls set initiator IQN and IP address for %v protocol.\n", DriverName, opt.GetAccessProtocol())
			}
			if opt.GetAccessProtocol() == FCProtocol {
				return nil, fmt.Errorf("%v: pls set initiator WWPN for %v protocol.\n", DriverName, opt.GetAccessProtocol())
			}
		}

		storageInitiatorGrpId, err := d.client.GetStorageInitiatorGrpId(poolId, initiatorName)
		if err != nil {
			return nil, err
		}

		// If specified initiator is nothing, register new initiator into default group.
		if storageInitiatorGrpId == "" {
			log.Infof("%v: trying to get default initiator group ID.", DriverName)
			storageInitiatorGrpId, err = d.client.GetDefaultInitiatorGrpId(poolId, opt)
			if err != nil {
				return nil, err
			}

			// Create default initiator group
			if storageInitiatorGrpId == "" {
				log.Infof("%v: trying to create default initiator group..", DriverName)
				respBody, err := d.client.CreateInitiatorDefaultGrp(poolId, opt)
				if err != nil {
					return nil, err
				}
				storageInitiatorGrpId = respBody.Id
			}

			// Register new initiator
			log.Infof("%v: trying to register initiator into default group..", DriverName)
			_, err := d.client.RegisterInitiatorIntoDefaultGrp(poolId, opt, storageInitiatorGrpId)
			if err != nil {
				return nil, err
			}
		}

		// Attach Volume
		attachRespBody, err := d.client.AttachVolume(poolId, opt.GetVolumeId(), storageInitiatorGrpId)
		if err != nil {
			return nil, err
		}

		// Set storage attachment ID
		// TODO: Im not sure how to save attachment ID which is UUID on side of storage
		// This UUID will be needed when terminate attachment
		opt.Metadata[opt.GetId()] = attachRespBody.Id

		// Get Volume Info
		tgtIqnWwn, tgtMgmtIp, err := d.client.GetTargetVolumeInfo(poolId, opt.GetVolumeId())
		if err != nil {
			return nil, err
		}
		if opt.GetAccessProtocol() == ISCSIProtocol {
			log.Infof("%v: attach volume for iSCSI success.", DriverName)

			return &model.ConnectionInfo{
				DriverVolumeType: ISCSIProtocol,
				ConnectionData: map[string]interface{}{
					"targetDiscovered": true,
					"targetIQN":        []string{tgtIqnWwn},
					"targetPortal":     []string{tgtMgmtIp},
					"discard":          false,
					"targetLun":        attachRespBody.Lun,
				},
			}, nil
		}
		if opt.GetAccessProtocol() == FCProtocol {
			log.Infof("%v: attach volume for FC success.", DriverName)
			return &model.ConnectionInfo{
				DriverVolumeType: FCProtocol,
				ConnectionData: map[string]interface{}{
					"targetDiscovered": true,
					"target_wwn":       []string{tgtIqnWwn},
					"volume_id":        opt.GetVolumeId(),
					"description":      "hpe",
					"host_name":        opt.GetHostInfo().Host,
					"targetLun":        attachRespBody.Lun,
				},
			}, nil
		}

	}

	return nil, fmt.Errorf("%v: Only support FC or iSCSI.\n", DriverName)
}

func (d *Driver) TerminateConnection(opt *pb.DeleteVolumeAttachmentOpts) error {
	poolId := opt.GetMetadata()["PoolId"]
	err := d.client.DetachVolume(poolId, opt.GetMetadata()[opt.GetId()])
	if err != nil {
		return err
	}

	// Delete attach OSDS ID <-> storage attach ID from meta data
	// TODO: Im not sure this way is correct. Should review.
	delete(opt.Metadata, opt.GetId())
	log.Infof("%v: detach volume success.", DriverName)
	return nil
}

func (d *Driver) CopyVolume(opt *pb.CreateVolumeOpts, srcid, tgtid string) error {
	return &model.NotImplementError{S: "method initializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) PullSnapshot(snapIdentifier string) (*model.VolumeSnapshotSpec, error) {
	// Not used, do nothing
	return nil, nil
}
func (d *Driver) PullVolume(volIdentifier string) (*model.VolumeSpec, error) {
	// Not used , do nothing
	return nil, nil
}

// The interfaces blow are optional, so implement it or not depends on you.
func (d *Driver) InitializeSnapshotConnection(opt *pb.CreateSnapshotAttachmentOpts) (*model.ConnectionInfo, error) {
	return nil, &model.NotImplementError{S: "method initializeSnapshotConnection has not been implemented yet."}
}

func (d *Driver) TerminateSnapshotConnection(opt *pb.DeleteSnapshotAttachmentOpts) error {
	return &model.NotImplementError{S: "method terminateSnapshotConnection has not been implemented yet."}
}

func (d *Driver) CreateVolumeGroup(opt *pb.CreateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method createVolumeGroup has not been implemented yet"}
}

func (d *Driver) UpdateVolumeGroup(opt *pb.UpdateVolumeGroupOpts) (*model.VolumeGroupSpec, error) {
	return nil, &model.NotImplementError{"method updateVolumeGroup has not been implemented yet"}
}

func (d *Driver) DeleteVolumeGroup(opt *pb.DeleteVolumeGroupOpts) error {
	return &model.NotImplementError{"method deleteVolumeGroup has not been implemented yet"}
}
