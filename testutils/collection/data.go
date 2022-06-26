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

This package includes a collection of fake stuffs for testing work.
*/

package collection

import (
	"github.com/sodafoundation/api/pkg/model"
)

var (
	SampleProfiles = []model.ProfileSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:             "default",
			Description:      "default policy",
			StorageType:      "block",
			CustomProperties: model.CustomPropertiesSpec{},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "silver",
			Description: "silver policy",
			StorageType: "block",
			CustomProperties: model.CustomPropertiesSpec{
				"dataStorage": map[string]interface{}{
					"provisioningPolicy": "Thin",
					"compression":        true,
					"deduplication":      true,
				},
				"ioConnectivity": map[string]interface{}{
					"accessProtocol": "rbd",
					"maxIOPS":        float64(5000000),
					"maxBWS":         float64(500),
					"minIOPS":        float64(1000000),
					"minBWS":         float64(100),
					"latency":        float64(100),
				},
			},
		},
	}

	SampleFileShareProfiles = []model.ProfileSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "1106b972-66ef-11e7-b172-db03f3689c9c",
			},
			Name:             "default",
			Description:      "default policy",
			StorageType:      "file",
			CustomProperties: model.CustomPropertiesSpec{},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "2f9c0a04-66ef-11e7-ade2-43158893e017",
			},
			Name:        "silver",
			Description: "silver policy",
			StorageType: "file",
			CustomProperties: model.CustomPropertiesSpec{
				"dataStorage": map[string]interface{}{
					"provisioningPolicy": "Thin",
					"compression":        true,
					"deduplication":      false,
				},
				"ioConnectivity": map[string]interface{}{
					"accessProtocol": "NFS",
					"maxIOPS":        float64(5000000),
					"maxBWS":         float64(500),
					"minIOPS":        float64(1000000),
					"minBWS":         float64(100),
					"latency":        float64(100),
				},
			},
		},
	}

	SampleCustomProperties = model.CustomPropertiesSpec{
		"dataStorage": map[string]interface{}{
			"provisioningPolicy": "Thin",
			"compression":        true,
			"deduplication":      true,
		},
		"ioConnectivity": map[string]interface{}{
			"accessProtocol": "rbd",
			"maxIOPS":        float64(5000000),
			"maxBWS":         float64(500),
			"minIOPS":        float64(1000000),
			"minBWS":         float64(100),
			"latency":        float64(100),
		},
	}

	SampleDocks = []model.DockSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "sample",
			Description: "sample backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "sample",
			Type:        model.DockTypeProvioner,
		},
	}

	SampleMultiDocks = []model.DockSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			},
			Name:        "sample-dock-01",
			Description: "sample backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "sample",
			Type:        model.DockTypeProvioner,
		},
		{
			BaseModel: &model.BaseModel{
				Id: "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			},
			Name:        "sample-dock-03",
			Description: "sample backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "cinder",
			Type:        model.DockTypeProvioner,
		},
		{
			BaseModel: &model.BaseModel{
				Id: "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg",
			},
			Name:        "sample-dock-02",
			Description: "sample backend service",
			Endpoint:    "localhost:50050",
			DriverName:  "ceph",
			Type:        model.DockTypeProvioner,
		},
	}

	SamplePools = []model.StoragePoolSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "084bf71e-a102-11e7-88a8-e31fe6d52248",
			},
			Name:             "sample-pool-01",
			Description:      "This is the first sample storage pool for testing",
			StorageType:      "block",
			TotalCapacity:    int64(100),
			FreeCapacity:     int64(90),
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			AvailabilityZone: "default",
			MultiAttach:      true,
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					Compression:        true,
					Deduplication:      false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "rbd",
					MaxIOPS:        8000000,
					MaxBWS:         700,
					MinIOPS:        1000000,
					MinBWS:         100,
					Latency:        100,
				},
				Advanced: map[string]interface{}{
					"diskType": "SSD",
					"latency":  "3ms",
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			},
			Name:             "sample-pool-02",
			Description:      "This is the second sample storage pool for testing",
			StorageType:      "block",
			TotalCapacity:    int64(200),
			FreeCapacity:     int64(170),
			AvailabilityZone: "default",
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy: "Thin",
					Compression:        true,
					Deduplication:      false,
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "rbd",
					MaxIOPS:        3000000,
					MaxBWS:         350,
					MinIOPS:        1000000,
					MinBWS:         100,
					Latency:        100,
				},
				Advanced: map[string]interface{}{
					"diskType": "SAS",
					"latency":  "500ms",
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "bdd44c8e-b8a9-488a-89c0-d1e5beb902dg",
			},
			Name:             "opensds-files-default",
			Description:      "This is the first file sample storage pool for testing",
			StorageType:      "file",
			TotalCapacity:    int64(200),
			FreeCapacity:     int64(170),
			AvailabilityZone: "default",
			DockId:           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			Extras: model.StoragePoolExtraSpec{
				DataStorage: model.DataStorageLoS{
					ProvisioningPolicy:      "Thin",
					Compression:             false,
					Deduplication:           false,
					StorageAccessCapability: []string{"Read", "Write", "Execute"},
				},
				IOConnectivity: model.IOConnectivityLoS{
					AccessProtocol: "nfs",
					MaxIOPS:        7000000,
					MaxBWS:         600,
					MinIOPS:        1000000,
					MinBWS:         100,
					Latency:        100,
				},
				Advanced: map[string]interface{}{
					"diskType": "SSD",
					"latency":  "5ms",
				},
			},
		},
	}

	SampleAvailabilityZones = []string{"default"}

	SampleFileShares = []model.FileShareSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "d2975ebe-d82c-430f-b28e-f373746a71ca",
			},
			Name:             "sample-fileshare-01",
			Description:      "This is first sample fileshare for testing",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "a5965ebe-dg2c-434t-b28e-f373746a71ca",
			ProfileId:        "b3585ebe-c42c-120g-b28e-f373746a71ca",
			SnapshotId:       "b7602e18-771e-11e7-8f38-dbd6d291f4eg",
			AvailabilityZone: "default",
			ExportLocations:  []string{"192.168.100.100"},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "1e643aca-4922-4b1a-bb98-4245054aeff4",
			},
			Name:             "sample-fileshare-2",
			Description:      "This is second sample fileshare for testing",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "d5f65ebe-ag2c-341s-a25e-f373746a71dr",
			ProfileId:        "1e643aca-4922-4b1a-bb98-4245054aeff4",
			SnapshotId:       "a5965ebe-dg2c-434t-b28e-f373746a71ca",
			AvailabilityZone: "default",
			ExportLocations:  []string{"192.168.100.101"},
		},
	}

	SampleFileSharesAcl = []model.FileShareAclSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "d2975ebe-d82c-430f-b28e-f373746a71ca",
			},
			Description: "This is a sample Acl for testing",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "1e643aca-4922-4b1a-bb98-4245054aeff4",
			},
			Description: "This is a sample Acl for testing",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "6ad25d59-a160-45b2-8920-211be282e2df",
			},
			Description:      "This is a sample Acl for testing",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
			Type:             "ip",
			AccessCapability: []string{"Read", "Write"},
			AccessTo:         "10.32.109.15",
			FileShareId:      "d2975ebe-d82c-430f-b28e-f373746a71ca",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "ad25d59-a160-45b2-8920-211be282e2dfh",
			},
			Description:      "This is a sample Acl for testing",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
			Type:             "ip",
			AccessCapability: []string{"Read", "Write"},
			AccessTo:         "10.32.109.151",
			FileShareId:      "d2975ebe-d82c-430f-b28e-f373746a71ca",
		},
	}

	SampleFileShareSnapshots = []model.FileShareSnapshotSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Name:         "sample-snapshot-01",
			Description:  "This is the first sample snapshot for testing",
			SnapshotSize: int64(1),
			FileShareId:  "d2975ebe-d82c-430f-b28e-f373746a71ca",
			Status:       "available",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			},
			Name:         "sample-snapshot-02",
			Description:  "This is the second sample snapshot for testing",
			SnapshotSize: int64(1),
			Status:       "available",
		},
	}

	SampleVolumeNames = []string{}

	SampleVolumes = []model.VolumeSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:             "sample-volume",
			Description:      "This is a sample volume for testing",
			AvailabilityZone: "default",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
			SnapshotId:       "",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:             "sample-volume",
			AvailabilityZone: "default",
			Description:      "This is a sample volume for testing",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
			SnapshotId:       "3769855c-a102-11e7-b772-17b880d2f537",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:             "sample-volume",
			Description:      "This is a sample volume for testing",
			AvailabilityZone: "default",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
	}

	SampleMultiVolumes = []model.VolumeSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:             "sample-volume-01",
			Description:      "This is a sample volume for testing",
			AvailabilityZone: "default",
			Size:             int64(2),
			Status:           "available",
			PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
			SnapshotId:       "",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:             "sample-volume-02",
			AvailabilityZone: "default",
			Description:      "This is a sample volume for testing",
			Size:             int64(1),
			Status:           "available",
			PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
			SnapshotId:       "3769855c-a102-11e7-b772-17b880d2f537",
		},
	}
	SampleVolumeWithDurableName = []model.VolumeSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:             "sample-volume-01",
			Description:      "This is a sample volume for testing",
			AvailabilityZone: "default",
			Size:             int64(2),
			Status:           "available",
			PoolId:           "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:        "1106b972-66ef-11e7-b172-db03f3689c9c",
			SnapshotId:       "",
			Identifier:       &model.Identifier{DurableName: "6216b2326e974b5fb0b3d2af5cd6b25b", DurableNameFormat: "NAA"},
		},
	}

	SampleShareNames = []string{}

	SampleShares = []model.FileShareSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:        "sample-fileshare",
			Description: "This is a sample fileshare for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			SnapshotId:  "",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			},
			Name:        "sample-fileshare",
			Description: "This is a sample fileshare for testing",
			Size:        int64(1),
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
			SnapshotId:  "3769855c-a102-11e7-b772-17b880d2f537",
		},
	}

	SampleConnection = model.ConnectionInfo{
		DriverVolumeType: "iscsi",
		ConnectionData: map[string]interface{}{
			"targetDiscovered": true,
			"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
			"targetPortal":     "127.0.0.0.1:3260",
			"discard":          false,
		},
	}

	SampleAttachments = []model.VolumeAttachmentSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			},
			Status:   "available",
			VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			HostId:   "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			ConnectionInfo: model.ConnectionInfo{
				DriverVolumeType: "iscsi",
				ConnectionData: map[string]interface{}{
					"targetDiscovered": true,
					"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal":     "127.0.0.0.1:3260",
					"discard":          false,
				},
			},
		},
	}

	SampleMultiAttachments = []model.VolumeAttachmentSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			},
			Status:   "available",
			VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			HostId:   "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			ConnectionInfo: model.ConnectionInfo{
				DriverVolumeType: "iscsi",
				ConnectionData: map[string]interface{}{
					"targetDiscovered": true,
					"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal":     "127.0.0.0.1:3260",
					"discard":          false,
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Status:   "attached",
			VolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d9",
			HostId:   "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			ConnectionInfo: model.ConnectionInfo{
				DriverVolumeType: "iscsi",
				ConnectionData: map[string]interface{}{
					"targetDiscovered": true,
					"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal":     "127.0.0.0.1:3260",
					"discard":          false,
				},
			},
		},
	}

	SampleSnapshots = []model.VolumeSnapshotSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Name:        "sample-snapshot-01",
			Description: "This is the first sample snapshot for testing",
			Size:        int64(1),
			Status:      "available",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			},
			Name:        "sample-snapshot-02",
			Description: "This is the second sample snapshot for testing",
			Size:        int64(1),
			Status:      "available",
			VolumeId:    "bd5b12a8-a101-11e7-941e-d77981b584d8",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
	}

	SampleShareSnapshots = []model.FileShareSnapshotSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f537",
			},
			Name:        "sample-snapshot-01",
			Description: "This is the first sample snapshot for testing",
			ShareSize:   int64(1),
			Status:      "available",
			FileShareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			},
			Name:        "sample-snapshot-02",
			Description: "This is the second sample snapshot for testing",
			ShareSize:   int64(1),
			Status:      "available",
			FileShareId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			ProfileId:   "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
	}

	SampleReplications = []model.ReplicationSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "c299a978-4f3e-11e8-8a5c-977218a83359",
			},
			PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
			SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			Name:              "sample-replication-01",
			Description:       "This is a sample replication for testing",
			PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "73bfdd58-4f3f-11e8-91c0-d39a05f391ee",
			},
			PrimaryVolumeId:   "bd5b12a8-a101-11e7-941e-d77981b584d8",
			SecondaryVolumeId: "bd5b12a8-a101-11e7-941e-d77981b584d8",
			Name:              "sample-replication-02",
			Description:       "This is a sample replication for testing",
			PoolId:            "084bf71e-a102-11e7-88a8-e31fe6d52248",
			ProfileId:         "1106b972-66ef-11e7-b172-db03f3689c9c",
		},
	}

	SampleVolumeGroups = []model.VolumeGroupSpec{
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f555",
			},
			Name:        "sample-group-01",
			Description: "This is the first sample group for testing",
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
		},
		{
			BaseModel: &model.BaseModel{
				Id: "3769855c-a102-11e7-b772-17b880d2f555",
			},
			Name:        "sample-group-01",
			Description: "This is the first sample group for testing",
			Status:      "available",
			PoolId:      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			Profiles: []string{"1106b972-66ef-11e7-b172-db03f3689c9c"},
		},
	}

	SampleHosts = []model.HostSpec{
		{
			BaseModel: &model.BaseModel{
				Id:        "202964b5-8e73-46fd-b41b-a8e403f3c30b",
				CreatedAt: "2019-11-11T11:01:33",
			},
			TenantId:          "x",
			AccessMode:        "agentless",
			HostName:          "sap1",
			IP:                "192.168.56.12",
			AvailabilityZones: []string{"default", "az2"},
			Initiators: []*model.Initiator{
				&model.Initiator{
					PortName: "20000024ff5bb888",
					Protocol: "iscsi",
				},
				&model.Initiator{
					PortName: "20000024ff5bc999",
					Protocol: "iscsi",
				},
			},
		},
		{
			BaseModel: &model.BaseModel{
				Id:        "eb73e59a-8b0f-4517-8b95-023ec134aec9",
				CreatedAt: "2019-11-11T11:13:57",
			},
			TenantId:          "x",
			AccessMode:        "agentless",
			HostName:          "sap2",
			IP:                "192.168.56.13",
			AvailabilityZones: []string{"default", "az2"},
			Initiators: []*model.Initiator{
				&model.Initiator{
					PortName: "20012324ff5ac132",
					Protocol: "iscsi",
				},
			},
		},
	}
)

// The Byte*** variable here is designed for unit test in client package.
// For how to ultilize these pre-assigned variables, please refer to
// (github.com/sodafoundation/api/client/dock_test.go).
var (
	ByteProfile = `{
		"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
		"name": "default",
		"description": "default policy",
		"storageType": "block"
	}`

	ByteProfiles = `[
		{
			"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
			"name": "default",
			"description": "default policy",
			"storageType": "block"
		},
		{
			"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
			"name": "silver",
			"description": "silver policy",
			"customProperties": {
				"dataStorage": {
					"provisioningPolicy": "Thin",
					"compression":   true,
					"deduplication": true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        5000000,
					"maxBWS":         500,
					"minIOPS": 	  1000000,
					"minBWS": 	  100,
					"latency":	  100
				}
			}
		}
	]`

	ByteCustomProperties = `{
		"dataStorage": {
			"provisioningPolicy": "Thin",
			"compression":   true,
			"deduplication": true
		},
		"ioConnectivity": {
			"accessProtocol": "rbd",
			"maxIOPS":        5000000,
			"maxBWS":         500,
			"minIOPS": 	  1000000,
			"minBWS": 	  100,
			"latency":	  100
		}
	}`

	ByteDock = `{
		"id": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		"name":        "sample",
		"description": "sample backend service",
		"endpoint":    "localhost:50050",
		"driverName":  "sample"
	}`

	ByteDocks = `[
		{
			"id": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"name":        "sample",
			"description": "sample backend service",
			"endpoint":    "localhost:50050",
			"driverName":  "sample"
		}
	]`

	BytePool = `{
		"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"name": "sample-pool-01",
		"description": "This is the first sample storage pool for testing",
		"storageType": "block",
		"totalCapacity": 100,
		"freeCapacity": 90,
		"dockId": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
		"extras": {
			"dataStorage": {
					"provisioningPolicy": "Thin",
					"compression":   true,
					"deduplication": true
				},
			"ioConnectivity": {
				"accessProtocol": "rbd",
				"maxIOPS":        1000
			},
			"advanced": {
				"diskType":   "SSD",
				"throughput": 1000
			}
		}
	}`

	BytePools = `[
		{
			"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"name": "sample-pool-01",
			"description": "This is the first sample storage pool for testing",
			"storageType": "block",
			"totalCapacity": 100,
			"freeCapacity": 90,
			"dockId": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"extras": {
				"dataStorage": {
					"provisioningPolicy": "Thin",
					"compression":   true,
					"deduplication": true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        8000000,
					"maxBWS": 	  700,
					"minIOPS": 	  1000000,
					"minBWS": 	  100,
					"latency":	  100
				},
				"advanced": {
					"diskType": "SSD",
					"latency":  "3ms"
				}
			}
		},
		{
			"id": "a594b8ac-a103-11e7-985f-d723bcf01b5f",
			"name": "sample-pool-02",
			"description": "This is the second sample storage pool for testing",
			"storageType": "block",
			"totalCapacity": 200,
			"freeCapacity": 170,
			"dockId": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"extras": {
				"dataStorage": {
					"provisioningPolicy": "Thin",
					"compression":   true,
					"deduplication": true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        3000000,
					"maxBWS": 	  350,
					"minIOPS": 	  1000000,
					"minBWS": 	  100,
					"latency":	  100
				},
				"advanced": {
					"diskType": "SAS",
					"latency":  "500ms"
				}
			}
		}
	]`

	ByteFileShare = `{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name": "sample-fileshare",
		"description": "This is a sample fileshare for testing",
		"size": 1,
		"status": "available",
		"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`

	ByteFileShares = `[
		{
			"id": "d2975ebe-d82c-430f-b28e-f373746a71ca",
			"createdAt": "",
			"updatedAt": "",
			"name": "sample-fileshare-01",
			"description": "This is first sample fileshare for testing",
			"size": 1,
			"availabilityZone": "default",
			"status": "available",
			"poolId": "a5965ebe-dg2c-434t-b28e-f373746a71ca",
			"profileId": "b3585ebe-c42c-120g-b28e-f373746a71ca",
			"snapshotId": "b7602e18-771e-11e7-8f38-dbd6d291f4eg",
			"exportLocations": [
				"192.168.100.100"
			]
		},
		{
			"id": "1e643aca-4922-4b1a-bb98-4245054aeff4",
			"createdAt": "",
			"updatedAt": "",
			"name": "sample-fileshare-2",
			"description": "This is second sample fileshare for testing",
			"size": 1,
			"availabilityZone": "default",
			"status": "available",
			"poolId": "d5f65ebe-ag2c-341s-a25e-f373746a71dr",
			"profileId": "1e643aca-4922-4b1a-bb98-4245054aeff4",
			"snapshotId": "a5965ebe-dg2c-434t-b28e-f373746a71ca",
			"exportLocations": [
				"192.168.100.101"
			]
		}
	]`

	ByteFileShareSnapshot = `{
		"id": "3769855c-a102-11e7-b772-17b880d2f537",
		"name": "sample-snapshot-01",
		"description": "This is the first sample snapshot for testing",
		"sharesize": 1,
		"status": "available",
		"fileshareId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`

	ByteFileShareSnapshots = `[
		{
			"id": "3769855c-a102-11e7-b772-17b880d2f537",
			"createdAt": "",
			"updatedAt": "",
			"name": "sample-snapshot-01",
			"description": "This is the first sample snapshot for testing",
			"snapshotSize": 1,
			"status": "available"
		},
		{
			"id": "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			"createdAt": "",
			"updatedAt": "",
			"name": "sample-snapshot-02",
			"description": "This is the second sample snapshot for testing",
			"snapshotSize": 1,
			"status": "available"
		}
	]`

	ByteFileShareAcl = `{
		"id": "d2975ebe-d82c-430f-b28e-f373746a71ca",
		"description": "This is a sample Acl for testing"
    }`

	ByteFileSharesAcls = `[
		{
			"id": "d2975ebe-d82c-430f-b28e-f373746a71ca",
			"createdAt": "",
			"updatedAt": "",
			"description": "This is a sample Acl for testing"
		},
		{
			"id": "1e643aca-4922-4b1a-bb98-4245054aeff4",
			"createdAt": "",
			"updatedAt": "",
			"description": "This is a sample Acl for testing"
		}
	]`

	ByteVolume = `{
		"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"name": "sample-volume",
		"description": "This is a sample volume for testing",
		"size": 1,
        "availabilityZone": "default",
		"status": "available",
		"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`

	ByteVolumes = `[
		{
			"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name": "sample-volume",
			"description": "This is a sample volume for testing",
			"size": 1,
            "availabilityZone": "default",
			"status": "available",
			"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
		}
	]`

	ByteAttachment = `{
		"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
		"name": "sample-volume-attachment",
		"description": "This is a sample volume attachment for testing",
		"status": "available",
		"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"hostId": "202964b5-8e73-46fd-b41b-a8e403f3c30b",
		"connectionInfo": {
			"driverVolumeType": "iscsi",
			"data": {
				"targetDiscovered": true,
				"targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
				"targetPortal": "127.0.0.0.1:3260",
				"discard": false
			}
		}
	}`

	ByteAttachments = `[
		{
			"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"name": "sample-volume-attachment",
			"description": "This is a sample volume attachment for testing",
			"status": "available",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"hostId": "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			"connectionInfo": {
				"driverVolumeType": "iscsi",
				"data": {
					"targetDiscovered": true,
					"targetIqn": "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal": "127.0.0.0.1:3260",
					"discard": false
				}
			}
		}
	]`

	ByteSnapshot = `{
		"id": "3769855c-a102-11e7-b772-17b880d2f537",
		"name": "sample-snapshot-01",
		"description": "This is the first sample snapshot for testing",
		"size": 1,
		"status": "available",
		"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
		"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`

	ByteVolumeGroup = `{
		"id": "3769855c-a102-11e7-b772-17b880d2f555",
		"name": "sample-group-01",
		"description": "This is the first sample group for testing",
		"status": "creating"
	}`

	ByteVolumeGroups = `[
		{
			"id": "3769855c-a102-11e7-b772-17b880d2f555",
			"name": "sample-group-01",
			"description": "This is the first sample group for testing",
			"status": "creating"
		}
	]`

	ByteSnapshots = `[
		{
			"id": "3769855c-a102-11e7-b772-17b880d2f537",
			"name": "sample-snapshot-01",
			"description": "This is the first sample snapshot for testing",
			"size": 1,
			"status": "available",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8"
		},
		{
			"id": "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			"name": "sample-snapshot-02",
			"description": "This is the second sample snapshot for testing",
			"size": 1,
			"status": "available",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8"
		}
	]`

	ByteReplication = `{
			"id": "c299a978-4f3e-11e8-8a5c-977218a83359",
			"primaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"secondaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name": "sample-replication-01",
			"description": "This is a sample replication for testing",
			"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
	}`

	ByteReplications = `[
		{
			"id": "c299a978-4f3e-11e8-8a5c-977218a83359",
			"primaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"secondaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name": "sample-replication-01",
			"description": "This is a sample replication for testing",
			"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
		},
		{
			"id": "73bfdd58-4f3f-11e8-91c0-d39a05f391ee",
			"primaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"secondaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name": "sample-replication-02",
			"description": "This is a sample replication for testing",
			"poolId": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId": "1106b972-66ef-11e7-b172-db03f3689c9c"
		}
	]`

	ByteVersion = `{
		"name": "v1beta",
		"status": "SUPPORTED",
		"updatedAt": "2017-04-10T14:36:58.014Z"
	}`

	ByteVersions = `[
		{
			"name": "v1beta",
			"status": "CURRENT",
			"updatedAt": "2017-07-10T14:36:58.014Z"
		}
	]`

	ByteHost = `{
        "id": "202964b5-8e73-46fd-b41b-a8e403f3c30b",
        "accessMode": "agentless",
        "createdAt": "2019-11-11T11:01:33",
        "tenantId": "x",
        "hostName": "sap1",
        "ip": "192.168.56.12",
        "availabilityZones": [
            "default",
            "az2"
        ],
        "initiators": [
            {
                "portName": "20000024ff5bb888",
                "protocol": "iscsi"
            },
            {
                "portName": "20000024ff5bc999",
                "protocol": "iscsi"
            }
        ]
    }`

	ByteHosts = `[
            {
                "id": "202964b5-8e73-46fd-b41b-a8e403f3c30b",
                "accessMode": "agentless",
                "createdAt": "2019-11-11T11:01:33",
                "tenantId": "x",
                "hostName": "sap1",
                "ip": "192.168.56.12",
                "availabilityZones": [
                    "default",
                    "az2"
                ],
                "initiators": [
                    {
                        "portName": "20000024ff5bb888",
                        "protocol": "iscsi"
                    },
                    {
                        "portName": "20000024ff5bc999",
                        "protocol": "iscsi"
                    }
                ]
            },
            {
                "id": "eb73e59a-8b0f-4517-8b95-023ec134aec9",
                "accessMode": "agentless",
                "createdAt": "2019-11-11T11:13:57",
                "tenantId": "x",
                "hostName": "sap2",
                "ip": "192.168.56.13",
                "availabilityZones": [
                    "default",
                    "az2"
                ],
                "initiators": [
                    {
                        "portName": "20012324ff5ac132",
                        "protocol": "iscsi"
                    }
                ]
            }
    	]`
)

// The StringSlice*** variable here is designed for unit test in etcd package.
// For how to ultilize these pre-assigned variables, please refer to
// (github.com/sodafoundation/api/pkg/db/drivers/etcd/etcd_test.go).
var (
	StringSliceProfiles = []string{
		`{
			"id": "1106b972-66ef-11e7-b172-db03f3689c9c",
			"name":        "default",
			"description": "default policy",
			"storageType": "block",
			"customProperties": {}
		}`,
		`{
			"id": "2f9c0a04-66ef-11e7-ade2-43158893e017",
			"name":        "silver",
			"description": "silver policy",
			"storageType": "block",
			"customProperties": {
				"dataStorage": {
					"provisioningPolicy": "Thin",
					"compression":   true,
					"deduplication": true
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        5000000,
					"maxBWS":         500,
					"minIOPS": 	  1000000,
					"minBWS": 	  100,
					"latency":	  100
				}
			}
		}`,
	}

	StringSliceDocks = []string{
		`{
			"id": "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"name":        "sample",
			"description": "sample backend service",
			"endpoint":    "localhost:50050",
			"driverName":  "sample",
			"type":        "provisioner"
		}`,
	}

	StringSlicePools = []string{
		`{
			"id": "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"name":             "sample-pool-01",
			"description":      "This is the first sample storage pool for testing",
			"storageType":		"block",
			"totalCapacity":    100,
			"freeCapacity":     90,
			"dockId":           "b7602e18-771e-11e7-8f38-dbd6d291f4e0",
			"availabilityZone": "default",
			"multiAttach": true,
			"extras": {
				"dataStorage": {
					"provisioningPolicy": "Thin",
					"compression":   true,
					"deduplication": false
				},
				"ioConnectivity": {
					"accessProtocol": "rbd",
					"maxIOPS":        8000000,
					"maxBWS": 	  700,
					"minIOPS": 	  1000000,
					"minBWS": 	  100,
					"latency":	  100
				},
				"advanced": {
					"diskType": "SSD",
					"latency":  "3ms"
				}
			}
		}`,
	}

	StringSliceVolumes = []string{
		`{
			"id": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name":        "sample-volume",
			"description": "This is a sample volume for testing",
			"size":        1,
            "availabilityZone": "default",
			"status":      "available",
			"poolId":      "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId":   "1106b972-66ef-11e7-b172-db03f3689c9c"
		}`,
	}

	StringSliceAttachments = []string{
		`{
			"id": "f2dda3d2-bf79-11e7-8665-f750b088f63e",
			"status":   "available",
			"volumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"hostId": "202964b5-8e73-46fd-b41b-a8e403f3c30b",
			"connectionInfo": {
				"driverVolumeType": "iscsi",
				"data": {
					"targetDiscovered": true,
					"targetIqn":        "iqn.2017-10.io.opensds:volume:00000001",
					"targetPortal":     "127.0.0.0.1:3260",
					"discard":          false
				}
			}
		}`,
	}

	StringSliceSnapshots = []string{
		`{
			"id": "3769855c-a102-11e7-b772-17b880d2f537",
			"name":        "sample-snapshot-01",
			"description": "This is the first sample snapshot for testing",
			"size":        1,
			"status":      "available",
			"volumeId":    "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"profileId":   "1106b972-66ef-11e7-b172-db03f3689c9c"
		}`,
		`{
			"id": "3bfaf2cc-a102-11e7-8ecb-63aea739d755",
			"name":        "sample-snapshot-02",
			"description": "This is the second sample snapshot for testing",
			"size":        1,
			"status":      "available",
			"volumeId":    "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"profileId":   "1106b972-66ef-11e7-b172-db03f3689c9c"
		}`,
	}

	StringSliceReplications = []string{
		`{
			"id":                "c299a978-4f3e-11e8-8a5c-977218a83359",
			"primaryVolumeId":   "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"secondaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name":              "sample-replication-01",
			"description":       "This is a sample replication for testing",
			"poolId":            "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId":         "1106b972-66ef-11e7-b172-db03f3689c9c"
		}`,
		`{
			"id":                "73bfdd58-4f3f-11e8-91c0-d39a05f391ee",
			"primaryVolumeId":   "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"secondaryVolumeId": "bd5b12a8-a101-11e7-941e-d77981b584d8",
			"name":              "sample-replication-02",
			"description":       "This is a sample replication for testing",
			"poolId":            "084bf71e-a102-11e7-88a8-e31fe6d52248",
			"profileId":         "1106b972-66ef-11e7-b172-db03f3689c9c"
		}`,
	}
)
