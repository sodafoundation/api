#  opensds roadmap

**work in progress**

This document defines a high level roadmap for opensds development.

Currently opensds team is working on an initial PoC code which aims to enable Kubernetes to
easily utilize storage provided by OpenStack Cinder, Manila, Swift as well as possibly NVMe
Over Fabric baremetal storage resource.

### Achievement Feb 9th, 2017
- Kubernetes connects to OpenSDS via OpenStack Fuxi Plugin (out of tree)
- OpenSDS connects to OpenStack Cinder and Manila via OpenStack Golang-client
- etcd for MQ and cluster mgmt.
- gRPC for communication between modules
