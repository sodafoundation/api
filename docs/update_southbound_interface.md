.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

http://creativecommons.org/licenses/by/3.0/legalcode

====================================
OpenSDS SourthBound Interface Design
====================================

Problem description
===================

Currently OpenSDS SourthBound interface is as follow:

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
	
From the interface shown above, we came to a conclusion that MountVolume and
UnmountVolume methods can not be handled by drivers, thus they should be
from here. The reasons are as follows:

1) From the point of semantics, mount and unmount operations belong to host
rather than storage backends, so there is no need for backends to receive
these two requests.

2) From the point of architecture design, OpenSDS only contains volume and
share resources. So when users want to mount a resource, just tell OpenSDS
which type of resource it should choose and let OpenSDS do the remaining
work.

3) From the point of implementation, if we move these two operations to
dock module, a lot of redundant code will be removed.

Proposed Change
===============

The main changes are as follows:

1) Remove mount and unmount operation in VolumeDriver and ShareDriver interface.

2) Create two packages(volume and share) in dock module, and move the two
operations above to new files(such as "volume_mount.go").

3) Remove code of these two operation in all backend drivers.

Data model impact
-----------------

After changed, the interface will be like this:

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
	}

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

None
