.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

http://creativecommons.org/licenses/by/3.0/legalcode

======================================
OpenSDS SourthBound Ceph Driver Design
======================================

Problem description
===================

As a SDS controller, it's essential for OpenSDS to build its eco-system of
sourthbound interface. At the first stage, our strategy is to quickly make up
for our lack with the help of OpenStack(Cinder, Manila). Now it's time to move
to next stage where we should build our own eco-system competitive.

After a careful consideration, we plan to select Ceph as the first OpenSDS
native sourthbound backend driver. The reasons are as follows:

1) Ceph is one of the most popular distributed storage systems in the world
and it holds a large number of users.

2) Ceph has a good performance in IO stream and data high availability.

3) It's open-source and has a large number of active contributors.

This proposal is launched mainly for the design of OpenSDS sourthbound Ceph
driver. With this standalone driver, OpenSDS can directly manage resources in
Ceph cluster and provide these storage resources for bare metals, VMs and
containers. 

Proposed Change
===============

Since OpenSDS repo has OpenSDS-plugins, Controller and Dock these three parts,
most of the work will be done in Dock module. we found that Ceph maintains an
official project "go-ceph", and we can manage resources of Ceph(pools, images
and so on) in Golang.The main jobs are as follows:

1) We need to implement CreateVolume, GetVolume, ListVolumes, DeleteVolume,
AttachVolume, DetachVolume, MountVolume and UnmountVolume in Ceph driver.
Here is the standardized interface:

	type VolumeDriver interface {
		//Any initialization the volume driver does while starting.
		Setup()
		//Any operation the volume driver does while stoping.
		Unset()

		CreateVolume(name string, volType string, size int32) (string, error)

		GetVolume(volID string) (string, error)

		GetAllVolumes(allowDetails bool) (string, error)

		DeleteVolume(volID string) (string, error)

		AttachVolume(volID string) (string, error)

		DetachVolume(device string) (string, error)

		MountVolume(mountDir, device, fsType string) (string, error)

		UnmountVolume(mountDir string) (string, error)
	}
	
2) From step 1 we can find that the return value is string type, which does
not seem to be a standardized description. And at the second step, what we
are going to do is that we will leverage Ceph and Cinder and design a new
unified sourthbound interface.

3) After those two steps, we will start to draft the V0.01 Spec of OpenSDS
sourthbound interface.


Data model impact
-----------------

Add ceph_driver element in Backends description.

REST API impact
---------------

None

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

https://github.com/noahdesu/go-ceph