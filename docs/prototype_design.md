## OpenSDS PoC Development (master branch)

_Please be aware that this code is under heavy development and subject to
change, do use with discreption._

### Purpose

The purpose of the opensds poc is to verify the concept we envisioned for opensds.
The PoC scenario is to have Kubernetes use OpenSDS as its storage provider via a
plugin, and OpenSDS use OpenStack Cinder and Manila for its storage resource infra-
structure.

The goal is to demonstrate the capability of having OpenSDS provision block and file
services that are provided by OpenStack for Kubernetes. OpenSDS's api will provide
a single entry for Kubernetes to talk to different OpenStack storage services.

### Structure
The current PoC code consists of three main components: api, orchestration and
adapter. Those three components communicate with each other with
the help of etcd.

* API module manages the request about storage resources, such as volumes, databases, file systems, policys and so forth. 

* Orchestration module has three roles: 

- Handles the request from API module.

- Collects the statistics (connection information, feature and so on) of
   storage resources through adapter module and deliver them to metaData
   module.
   
- Orchestrates storage resources and shows appropriate resources to users
   according to scenarios.

* Adapter module contains a standard Dock and plugins of different storage backends which contains both open source projects (such as Cinder, Manila, Swift and so on) and enterprise projects (such as Vipr, OceanStor DJ). Raw storage device will also be supported later



