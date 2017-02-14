# OpenSDS

[![Go Report Card](https://goreportcard.com/badge/github.com/opensds/opensds)](https://goreportcard.com/report/github.com/opensds/opensds)
[![Build Status](https://travis-ci.org/opensds/opensds.svg?branch=master)](https://travis-ci.org/opensds/opensds)

<img src="https://www.opensds.io/wp-content/uploads/2016/11/logo_opensds.png" width="100">

## Introduction

The [OpenSDS Project](https://opensds.io/) is a collaborative project under Linux
Foundation supported by storage users and vendors, including
EMC, Intel, Huawei, Fujitsu, Western Digital, Vodafone and Oregon State University. The project
will also seek to collaborate with other upstream open source communities
such as Cloud Native Computing Foundation, Docker, OpenStack, and Open
Container Initiative. 

### Community

The OpenSDS Project is currently running as a technical community which
focus on developing a working PoC code and working on a formal charter
targeted the mid of 2017.

The OpenSDS community welcomes anyone who is interested in software defined
storage and shaping the future of cloud-era storage. If you are a company,
you should consider joining the [OpenSDS Project](https://opensds.io/). 
If you are a developer want to be part of the PoC development that is happening
now, please refer to the Contribute sections below.

The current opensds team that is developing the PoC prototype comes from Huawei,Intel,
EMC and Wetern Digital.

### Collaborative Testing

* [CNCF Cluster](https://github.com/cncf/cluster/issues/30)
* OpenStack OISC (submitted)

### Contact

- Mailing list: [opensds-dev](https://groups.google.com/forum/?hl=en#!forum/opensds-dev)
- slack: #[opensds](https://opensds.slack.com)
- Planning/Roadmap: [milestones](https://github.com/opensds/opensds/milestones), [roadmap](./ROADMAP.md)
- Ideas/Bugs: [issues](https://github.com/opensds/opensds/issues)

### Contribute

If you're interested in being a contributor and want to get involved in the
OpenSDS PoC code developing, please see [CONTRIBUTING](CONTRIBUTING.md) for 
details on submitting patches and the contribution workflow.

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
adapter. Those three components communicate with each other through gRPC with
the help of etcd.

* API module manages the request about storage resources, such as volumes, databases, file systems, policys and so forth. 

* Orchestration module has three roles: 

**Handles the request from API module.**

**Collects the statistics (connection information, feature and so on) of
   storage resources through adapter module and deliver them to metaData
   module.**
   
**Orchestrates storage resources and shows appropriate resources to users
   according to scenarios.**

* Adapter module contains a standard Dock and plugins of different storage backends which contains both open source projects (such as Cinder, Manila, Swift and so on) and enterprise projects (such as Vipr, OceanStor DJ). Raw storage device will also be supported later

### Usage/Hacking

#### Requirements

* Clustering : etcd

For easy installation, the download link of etcd binary file is at: https://github.com/coreos/etcd/releases

You can just use the command followed to set up etcd service:

1. curl -L  https://github.com/coreos/etcd/releases/download/v2.0.0-rc.1/etcd-v2.0.0-rc.1-linux-amd64.tar.gz -o etcd-v2.0.0-rc.1-linux-amd64.tar.gz

2. tar xzvf etcd-v2.0.0-rc.1-linux-amd64.tar.gz

3. cd etcd-v2.0.0-rc.1-linux-amd64

4. ./etcd

* Infrastructre : OpenStack, OceanStor DJ, CoprHD or ...

As a Software-Defined-Storage controller, OpenSDS must be able to connect to backend storage environment. You can just deploy OpenSDS in any of these environments directly to avoid troublesomes. Since OpenSDS expose its interface through gRPC, the user needs to utilize the golang-client for interacting with OpenStack infrastructures. If you have infrastructures other than OpenStack, you could develop the corresponding plugin for your infrastructure.

* Language : Go environment

To run a Go project, configuring Gopath is indispensable. After downloading the project, you should add "path_to_OpenSDS" in GOPATH environ variable.

#### Build

1. export GOPATH=$HOME/gopath

   export PATH=$HOME/gopath/bin:$PATH
   
   mkdir -p $HOME/gopath/src/github.com/opensds/
   
   cd $HOME/gopath/src/github.com/opensds
   
2. git clone https://github.com/opensds/opensds.git $HOME/gopath/src/github.com/opensds/

3. cd opensds (import necessary packages)

   go get github.com/spf13/cobra

   go get github.com/astaxie/beego

   go get github.com/coreos/etcd/client
   
4. cd cmd/sdslet

   go build
   
5. cd cmd/sdsctl

   go buld
   
6. cp cmd/sdslet/sdslet /usr/local/bin

   cp cmd/sdsctl/sdsctl /usr/local/bin

7. vim examples/config.json (config backend storage credential information)

   sudo mkdir /etc/opensds

   sudo cp examples/config.json /etc/opensds/

8. sudo touch /var/log/opensds.log (create OpenSDS logging file)

#### Run

* Make sure **etcd** is up

```sh
./bin/etcd
```

* Start **sdslet** with root access (for logging purpose)

```sh
sudo sdslet //suppose the user has copied the compiled binary to /usr/local/bin
```

* Run **sdsctl** for operations you want to perform. 

```sh
sdsctl --help //see what you can do with opensds
```

Currently sdsctl supports all the basic Cinder/Manila operations, for example if you want to 
create a 1GB volume from a Dell-EMC VMAX, which is connected to the OpenSDS underlay infra - 
OpenStack Cinder via its in-tree vmax cinder driver, using OpenSDS for an easy access:

```sh
sdsctl volume create 1 -n cinder-vmax-volume -b cinder
```
Viola !
