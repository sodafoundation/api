### Requirements

- Clustering : etcd

For easy installation, the download link of etcd binary file is at: https://github.com/coreos/etcd/releases

You can just use the command followed to set up etcd service:
```shell
curl -L https://github.com/coreos/etcd/releases/download/v3.2.0/etcd-v3.2.0-linux-amd64.tar.gz -o $HOME/etcd-v3.2.0-linux-amd64.tar.gz
cd $HOME && tar xzvf etcd-v3.2.0-linux-amd64.tar.gz
cd etcd-v3.2.0-linux-amd64 && nohup ./etcd &>>etcd.log &
```

- Infrastructre : OpenStack, Bare-metal(LVM, Ceph), OceanStor DJ, CoprHD or ...

As a Software-Defined-Storage controller, OpenSDS must be able to connect to backend storage environment. You can just deploy OpenSDS 
in any of these environments directly to avoid troublesomes. Since OpenSDS is writen in Go, the user needs to utilize the golang-client 
for interacting with OpenStack infrastructures. If you have infrastructures other than OpenStack, you could develop the corresponding
plugin for your infrastructure.

- Language : Go environment

To run a Go project, configuring Gopath is indispensable. After downloading the project, you should add "path_to_OpenSDS" in GOPATH environ variable.
