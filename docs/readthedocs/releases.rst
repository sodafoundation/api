Official Releases
=================

Bali
-----

The OpenSDS Bali release completed December 19, 2018.

Features
>>>>>>>>

The Bali release adds the following functionality:

* Introduced management of multiple OpenStack deployments
* Multi-Cloud support
   - S3 API support with AWS, Azure, Huawei Cloud, and local Ceph object store and Fusion Storage backends.
   - Manual and basic policy based migration for AWS, Azure, Huawei Cloud, Ceph, and Fusion Storage.
* Dashboard UI interface has been broken out into separately installed component and integrated with multicloud support.
* Added ability to upload/download snapshot to/from cloud storage
* Updated support for the Container Storage Interface (CSI) v1.0 specification
   - Added support to create snapshot and create volume from snapshot
   - Added support for NodeStageVolume and NodeUnstageVolume.
* Support to provision replicated volumes using OpenSDS CSI plugin
* CSI plugin refactoring and FC support
* Southbound Fusion Storage and OceanStor V3/V5 volume drivers
* Integrated profiles properties definition and selector filtering.
* Support external volumes for VMs or baremetal.
* Add API support for AvailabilityZone.
* Installation with Helm (tested with LVM)

The OpenSDS controller (Hotpot), the north bound plugins (Sushi), the multiple cloud(Gelato), 
the opensds dashboard and the installer can be downloaded from here:

`Hotpot <https://github.com/sodafoundation/api/releases/tag/v0.4.0>`__
`Sushi <https://github.com/opensds/nbp/releases/tag/v0.4.0>`__
`Gelato <https://github.com/opensds/multi-cloud/releases/tag/v0.4.0>`__
`Dashboard <https://github.com/sodafoundation/api-dashboard/releases/tag/v0.4.0>`__
`Installer <https://github.com/sodafoundation/api-installer/releases/tag/v0.4.0>`__


Aruba
-----

The OpenSDS Aruba release completed June 30, 2018.

Features
>>>>>>>>

The Aruba release adds the following functionality:

* Array-based replication
* Cinder compatible API
* Containerized deployment
* Controller API request filter
* Create volume from snapshot
* Dashboard UI interface
* Extend volume support
* Fibre channel protocol support
* Host-based replication
* Multi-tenancy support in the API
* OpenStack Keystone authentication
* Storage backend capabilities reporting
* Storage pool capability reporting
* Volume groups

The OpenSDS controller (Hotpot), the north bound plugins (Sushi), and the
installer can be downloaded from here:

`Hotpot <https://github.com/sodafoundation/api/releases/tag/v0.2.0>`__
`Sushi <https://github.com/opensds/nbp/releases/tag/v0.2.0>`__
`Installer <https://github.com/sodafoundation/api-installer/releases/tag/v0.2.0>`__
