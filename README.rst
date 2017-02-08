.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

# OpenSDS

<img src="https://www.opensds.io/wp-content/uploads/2016/11/logo_opensds.png" width="100">

## Introduction

The [OpenSDS Project](https://opensds.io/) is supported by storage users and vendors, including
Huawei, Fujitsu, HDS, Vodafone and Oregon State University. The project
will also seek to collaborate with other upstream open source communities
such as Cloud Native Computing Foundation, Docker, OpenStack, and Open
Container Initiative. The OpenSDS project is a collaborative under Linux
Foundation.

## Community

The OpenSDS Project is currently running as a technical community which
focus on developing a working PoC code and working on a formal charter
targeted the mid of 2017.

The OpenSDS community welcomes anyone who is interested in software defined
storage and shaping the future of cloud-era storage. If you are a company,
you should consider joining the [OpenSDS Project](https://opensds.io/). If
you are a developer want to be part of the PoC development that is happening
now, please register to the [OpenSDS Mailinglist](https://groups.google.com/forum/#!forum/opensds-dev/) to get involved.

## Contribute

If you're interested in being a contributor and want to get involved in the
OpenSDS PoC code developing, please feel free to fork the code, raise an issue
and submit your contribution via PR. 

## PoC Code Introduction

* Please be aware that this code is under heavy development and subject to
change, do use with discreption.

The current PoC code consists of three main components: API, orchestration and
adapter. Those three components communicate with each other through gRPC with
the help of etcd.

API module manages the request about storage resources, such as volumes,
databases, file systems, policys and so forth.

Orchestration module has three roles:

1. Handles the request from API module.

2. Collects the statistics (connection information, feature and so on) of
   storage resources through adapter module and deliver them to metaData
   module.
   
3. Orchestrates storage resources and shows appropriate resources to users
   according to scenarios.

Adapter module contains a standard storageDock and plugins of cookedStorage
and rawStorage. The cookedStorage contains open source projects (such as
Cinder, Manila, Swift and so on) and enterprise projects (such as
OceanStor DJ). The rawStorage contains raw storage device from Intel and
WD (such as NVMe and NOF).


