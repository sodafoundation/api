.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

http://creativecommons.org/licenses/by/3.0/legalcode

======================================
Kubernetes External Provisioner Design
======================================

Problem description
===================

As we know, Kubernetes provides a framework of storage for container
orchestration layer. When users create a PVC, the controller will retrieve PV
from backends through "Provisioner" interface. But right now each backend has
its own provisioner(nfs, efs, cephfs and so on), and it's unbearable for vendors
to maintain all provisoners if they want to support multiple backends.

As a unified SDS controller, it's necessary for OpenSDS to solve this problem
and thus ease burden of vendors. So we want to create a new OpenSDS provisioner
that shields the detailed implementation of banckends and expose an abstract
description to external storage interface. In that way, all vendors don't need
to care about any backends but just select a backend_type param and assign it
to OpenSDS provisoner.

Proposed Change
===============

From a Kubernetes-incubator project "external-storage", we found that though
several provisoners exist, they all have to follow a standard interface called
"Provison", and this interface doesn't care about what type backends are.
So our proposal is that we could follow this interface and develop a unified
provisoner, and let StorageProfile(resource in OpenSDS) bear the personality of
backends. The main work are as follows:

1) Move the "external-storage" project into OpenSDS repo, cut the redundant
module and develop a new provisioner implementing Provision and Delete method.
And it seems that we should add our VolumeSourceSpec into Kubernetes API.

2) We will design and implemant StorageProfile in OpenSDS controller module.
Profile is an abstract description like StorageClass in Kubernetes, but one
difference is that the former one contains an element "Tags" to identify some
specified backend features. When user creates a PVC, the Provision interface
will pass backend_type and other optional features(IOPS, HA and so on) to
Profile, and then OpenSDS controller service will parse the structure and
schedule the required bankend to do rest work.

3) Since Profile and Tags can only be visible to Admin, then we could add
identity and authentication service in controller if necessary, but we can
ignore it at the early stage.

Data model impact
-----------------

Profile:
    description: An OpenSDS profile is identified by a unique name and UUID.
      Each profile has a set of OpenSDS tags which are desirable features for a 
      class of applications.
    type: object
    required:
      - name
      - tags
    properties:
      name:
        type: string
      uuid:
        type: string
        readOnly: true
	  storageName:
		type: string
		enum:
          - "CIFS"
          - "Ceph"
          - "CephFS"
          - "Dell EMC Isilon"
          - "Dell EMC VMAX"
          - "Dell EMC VNX"
          - "Dell EMC Unity"
          - "GlusterFS"
          - "HDFS"
          - "Hitachi Block Storage"
          - "Hitachi NAS"
          - "Hitachi Hyper Scale-Out Platform"
          - "HPE 3PAR"
          - "Huawei Dorado"
          - "Huawei Fusion Storage"
          - "Huawei OceanStor"
          - "IBM FlashSystem"
          - "IBM GPFS"
          - "IBM Storwize"
          - "LVM"
          - "NetApp Data ONTAP"
          - "NetApp E-Series"
          - "NFS"
          - "solidfire"
          - "Violin Memory 7000"
      tags:
        $ref: "#/definitions/StorageTags"
		
StorageTag:
    description: An OpenSDS tag represents a service provided by the 
      backend storage.
    type: object
    properties:
      dataAvailabilityTag:
        $ref: "#/definitions/DataAvailabilityTag"
      dataDurabilityTag:
        $ref: "#/definitions/DataDurabilityTag"
      dataPersistencyTag:
        $ref: "#/definitions/DataPersistencyTag"
      dataProtectionTag:
        $ref: "#/definitions/DataProtectionTag"
      dataOptimizationTag:
        $ref: "#/definitions/DataOptimizationTag"
      dataReplicationTag:
        $ref: "#/definitions/DataReplicationTag"
      dataSecurityTag:
        $ref: "#/definitions/DataSecurityTag"
      ioTypeTag:
        $ref: "#/definitions/IOTypeTag"
      performanceTag:
        $ref: "#/definitions/PerformanceTag"
      protocolTag:
        $ref: "#/definitions/ProtocolTag"
      qosTag:
        $ref: "#/definitions/QosTag"
    minProperties: 1
    maxProperties: 1

REST API impact
---------------

URL: /apis/opensds.io/v1beta/profiles
Method: POST, GET and DELETE
Description: add create, get, list and delete operations of Profile resource.

Security impact
---------------

None

Other end user impact
---------------------

None

Performance impact
------------------

None

Other deployer impact
---------------------

None

Dependencies
============

None

Testing
=======

None

Documentation Impact
====================

None

References
==========

https://github.com/kubernetes-incubator/external-storage

https://github.com/sodafoundation/api-api-specs/blob/master/northbound-api/v1/openapi-spec/swagger.yaml