### Requirements

* Clustering : etcd

For easy installation, the download link of etcd binary file is at: https://github.com/coreos/etcd/releases

You can just use the command followed to set up etcd service:

1. curl -L  https://github.com/coreos/etcd/releases/download/v2.0.0-rc.1/etcd-v2.0.0-rc.1-linux-amd64.tar.gz -o etcd-v2.0.0-rc.1-linux-amd64.tar.gz

2. tar xzvf etcd-v2.0.0-rc.1-linux-amd64.tar.gz

3. cd etcd-v2.0.0-rc.1-linux-amd64

4. ./etcd

* Infrastructre : OpenStack, OceanStor DJ, CoprHD or ...

As a Software-Defined-Storage controller, OpenSDS must be able to connect to backend storage environment. You can just deploy OpenSDS 
in any of these environments directly to avoid troublesomes. Since OpenSDS is writen in Go, the user needs to utilize the golang-client 
for interacting with OpenStack infrastructures. If you have infrastructures other than OpenStack, you could develop the corresponding
plugin for your infrastructure.

* Language : Go environment

To run a Go project, configuring Gopath is indispensable. After downloading the project, you should add "path_to_OpenSDS" in GOPATH
environ variable.
