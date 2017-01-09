.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

This project consists of three main components: API, orchestration and
adapter. Those three components communicate with each other through gRPC
mechanism (using etcd).

API module manages the request about storage resources, such as volumes,
databases, file systems, policys and so forth.

Orchestration module has three roles:

1. Handles the request from API module.

2. Collects the statistics (connection information, feature and so on) of
   storage resources through adapter module and deliver them to metaData
   module. (not achieved now)

3. Orchestrates storage resources and shows appropriate resources to users
   according to scenarios.

Adapter module contains a standard storageDock and plugins of cookedStorage
and rawStorage. The cookedStorage contains open source projects (such as
Cinder, Manila, Swift and so on) and enterprise projects (such as
OceanStor DJ). The rawStorage contains raw storage device from Intel and
WD (such as NVMe and NOF).
