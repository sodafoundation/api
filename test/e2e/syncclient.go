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

package e2e

import (
	"fmt"
	"time"

	"github.com/opensds/opensds/client"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
)

const (
	waitInterval = 500 * time.Millisecond
	waitTimeout  = 10 * time.Second
)

type SyncClient struct {
	*client.Client
}

func NewSyncClient(config *client.Config) (*SyncClient, error) {
	c, err := client.NewClient(config)
	return &SyncClient{Client: c}, err
}

func (s *SyncClient) waitVolumeAvailable(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		vol, err := s.Client.GetVolume(id)
		if err != nil {
			return false, err
		}
		if vol.Status == model.VolumeError {
			return false, fmt.Errorf("volume in error status")
		}
		return vol.Status == model.VolumeAvailable, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitVolumeDeleted(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		vols, err := s.Client.ListVolumes()
		if err != nil {
			return false, err
		}

		for _, vol := range vols {
			if vol.Id == id {
				if vol.Status == model.VolumeErrorDeleting {
					return false, fmt.Errorf("volume in error_deleting status")
				}
				return false, nil
			}
		}
		return true, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitVolumeAttachmentAvailable(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		atc, err := s.Client.GetVolumeAttachment(id)
		if err != nil {
			return false, err
		}
		if atc.Status == model.VolumeAttachError {
			return false, fmt.Errorf("volume attachment in error status")
		}
		return atc.Status == model.VolumeAttachAvailable, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitVolumeAttachmentDeleted(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		atcs, err := s.Client.ListVolumeAttachments()
		if err != nil {
			return false, err
		}

		for _, atc := range atcs {
			if atc.Id == id {
				if atc.Status == model.VolumeAttachErrorDeleting {
					return false, fmt.Errorf("volume attachment in error_deleting status")
				}
				return false, nil
			}
		}
		return true, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitVolumeSnapshotAvailable(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		snp, err := s.Client.GetVolumeSnapshot(id)
		if err != nil {
			return false, err
		}
		if snp.Status == model.VolumeSnapError {
			return false, fmt.Errorf("volume snapshot in error status")
		}
		return snp.Status == model.VolumeSnapAvailable, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitVolumeSnapShotDeleted(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		snps, err := s.Client.ListVolumeSnapshots()
		if err != nil {
			return false, err
		}
		for _, snp := range snps {
			if snp.Id == id {
				if snp.Status == model.VolumeSnapErrorDeleting {
					return false, fmt.Errorf("volume snapshot in error_deleting status")
				}
				return false, nil
			}
		}
		return true, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitVolumeGroupAvailable(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		vg, err := s.Client.GetVolumeGroup(id)
		if err != nil {
			return false, err
		}
		if vg.Status == model.VolumeGroupError {
			return false, fmt.Errorf("volume group in error status")
		}
		return vg.Status == model.VolumeGroupAvailable, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitVolumeGroupDeleted(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		vgs, err := s.Client.ListVolumeGroups()
		if err != nil {
			return false, err
		}
		for _, vg := range vgs {
			if vg.Id == id {
				if vg.Status == model.VolumeGroupErrorDeleting {
					return false, fmt.Errorf("volume group in error_deleting status")
				}
				return false, nil
			}
		}
		return true, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitFileShareAvailable(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		fs, err := s.Client.GetFileShare(id)
		if err != nil {
			return false, err
		}
		if fs.Status == model.FileShareError {
			return false, fmt.Errorf("fileshare in error status")
		}
		return fs.Status == model.FileShareAvailable, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitFileShareDeleted(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		fileshares, err := s.Client.ListFileShares()
		if err != nil {
			return false, err
		}
		for _, fileshare := range fileshares {
			if fileshare.Id == id {
				if fileshare.Status == model.FileShareErrorDeleting {
					return false, fmt.Errorf("volume fileshare in error_deleting status")
				}
				return false, nil
			}
		}
		return true, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) waitFileShareSnapshotAvailable(id string) error {
	return utils.WaitForCondition(func() (bool, error) {
		vol, err := s.Client.GetFileShareSnapshot(id)
		if err != nil {
			return false, err
		}
		return vol.Status == model.FileShareSnapAvailable, nil
	}, waitInterval, waitTimeout)
}

func (s *SyncClient) CreateVolume(body client.VolumeBuilder) (*model.VolumeSpec, error) {
	vol, err := s.Client.CreateVolume(body)
	if err != nil {
		return nil, err
	}
	if err := s.waitVolumeAvailable(vol.Id); err != nil {
		return nil, err
	}
	return vol, nil
}

func (s *SyncClient) ExtendVolume(volID string, body client.ExtendVolumeBuilder) (*model.VolumeSpec, error) {
	vol, err := s.Client.ExtendVolume(volID, body)
	if err != nil {
		return nil, err
	}
	if err := s.waitVolumeAvailable(vol.Id); err != nil {
		return nil, err
	}
	return vol, nil
}

func (s *SyncClient) DeleteVolume(volID string, body client.VolumeBuilder) error {
	if err := s.Client.DeleteVolume(volID, body); err != nil {
		return err
	}
	return s.waitVolumeDeleted(volID)
}

func (s *SyncClient) CreateVolumeSnapshot(body client.VolumeSnapshotBuilder) (*model.VolumeSnapshotSpec, error) {
	snp, err := s.Client.CreateVolumeSnapshot(body)
	if err != nil {
		return nil, err
	}
	if err := s.waitVolumeSnapshotAvailable(snp.Id); err != nil {
		return nil, err
	}
	return snp, nil
}

func (s *SyncClient) DeleteVolumeSnapshot(snpID string, body client.VolumeSnapshotBuilder) error {
	if err := s.Client.DeleteVolumeSnapshot(snpID, body); err != nil {
		return err
	}
	return s.waitVolumeSnapShotDeleted(snpID)
}

func (s *SyncClient) CreateVolumeAttachment(body client.VolumeAttachmentBuilder) (*model.VolumeAttachmentSpec, error) {
	atc, err := s.Client.CreateVolumeAttachment(body)
	if err != nil {
		return nil, err
	}
	if err := s.waitVolumeAttachmentAvailable(atc.Id); err != nil {
		return nil, err
	}
	return atc, nil
}

func (s *SyncClient) DeleteVolumeAttachment(atcID string, body client.VolumeAttachmentBuilder) error {
	if err := s.Client.DeleteVolumeAttachment(atcID, body); err != nil {
		return err
	}
	return s.waitVolumeAttachmentDeleted(atcID)
}

func (s *SyncClient) CreateVolumeGroup(body client.VolumeGroupBuilder) (*model.VolumeGroupSpec, error) {
	vg, err := s.Client.CreateVolumeGroup(body)
	if err != nil {
		return nil, err
	}
	if err := s.waitVolumeGroupAvailable(vg.Id); err != nil {
		return nil, err
	}
	return vg, nil
}

func (s *SyncClient) DeleteVolumeGroup(vgId string, body client.VolumeGroupBuilder) error {
	if err := s.Client.DeleteVolumeGroup(vgId, body); err != nil {
		return err
	}
	return s.waitVolumeGroupDeleted(vgId)
}

func (s *SyncClient) CreateFileShare(body client.FileShareBuilder) (*model.FileShareSpec, error) {
	fs, err := s.Client.CreateFileShare(body)
	if err != nil {
		return nil, err
	}
	if err := s.waitFileShareAvailable(fs.Id); err != nil {
		return nil, err
	}
	return fs, nil
}

func (s *SyncClient) DeleteFileShare(ID string) error {
	if err := s.Client.DeleteFileShare(ID); err != nil {
		return err
	}
	return s.waitFileShareDeleted(ID)
}
