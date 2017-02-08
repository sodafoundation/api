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

_Please be aware that this code is under heavy development and subject to
change, do use with discreption._

The current PoC code consists of three main components: API, orchestration and
adapter. Those three components communicate with each other through gRPC with
the help of etcd.

** API module manages the request about storage resources, such as volumes, 
databases, file systems, policys and so forth.

** Orchestration module has three roles:

* Handles the request from API module.

* Collects the statistics (connection information, feature and so on) of
   storage resources through adapter module and deliver them to metaData
   module.
   
* Orchestrates storage resources and shows appropriate resources to users
   according to scenarios.

** Adapter module contains a standard storageDock and plugins of cookedStorage
and rawStorage. The cookedStorage contains open source projects (such as
Cinder, Manila, Swift and so on) and enterprise projects (such as
OceanStor DJ). The rawStorage contains raw storage device from Intel and
WD (such as NVMe and NOF).

## POC Installation

# Requirement

* etcd

For easy installation, the download link of etcd binary file is at: https://github.com/coreos/etcd/releases
You can just use the command followed to set up etcd service:
1. curl -L  https://github.com/coreos/etcd/releases/download/v2.0.0-rc.1/etcd-v2.0.0-rc.1-linux-amd64.tar.gz -o etcd-v2.0.0-rc.1-linux-amd64.tar.gz
2. tar xzvf etcd-v2.0.0-rc.1-linux-amd64.tar.gz
3. cd etcd-v2.0.0-rc.1-linux-amd64
4. ./etcd

* OpenStack, OceanStor DJ, CoprHD or ...

As a Software-Defined-Storage controller, OpenSDS must be able to connect to backend storage environment. You can just deploy OpenSDS in any of these environments directly to avoid troublesomes. Since OpenSDS expose its interface through gRPC, there is nothing else to install. Lastly, please attention that POC only support OpenStack and OceanStor DJ right now, but we are working on other backend-storage environments.

* Go environment

To run a Go project, configuring Gopath is indispensable. After downloading the project, you should add "path_to_OpenSDS" in GOPATH environ variable.
