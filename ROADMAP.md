#  opensds roadmap

**work in progress**

This document defines a high level roadmap for opensds development.

Currently opensds team is working on an initial **PoC code** which aims to enable Kubernetes to
easily utilize storage provided by OpenStack Cinder, Manila, Swift as well as possibly NVMe
Over Fabric baremetal storage resource.

### Achievement Feb 9th, 2017
- Kubernetes connects to OpenSDS via OpenStack Fuxi Plugin (out of tree)
- OpenSDS connects to OpenStack Cinder and Manila via OpenStack Golang-client
- etcd for MQ and cluster mgmt.

### Achievement Feb 14th, 2017
- REST framework built with [Beego](https://github.com/astaxie/beego)

### To-do (shorter term)
- Develop flexvolume out-of-tree plugin
- Using etcd3 and gRPC

### To-do (short term)
- Plan for collaboration with libStorage
- Plan for involvement with k8s-storage-sig Flex2 development
- NVMe Over Fabric related developments
- Containerization of opensds modules (api,orchestration,adaptor)
