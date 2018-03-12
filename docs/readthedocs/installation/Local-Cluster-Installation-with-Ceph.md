This installation document assumes that you use [CEPH](https://github.com/ceph/ceph) as the default driver, and please make sure your ceph cluster deployed.

## Pre-configuration

### Bootstrap
If you have a clean environment (suggest Ubuntu16.04+), please run the script
to install all dependencies of this project (except ceph cluster):
```
curl -sSL https://raw.githubusercontent.com/opensds/opensds/master/script/cluster/bootstrap.sh | sudo bash
```

If everything works well, in the end you will find the output below:
```bash
go get github.com/opensds/opensds/cmd/osdslet
go get github.com/opensds/opensds/cmd/osdsdock
go get github.com/opensds/opensds/cmd/osdsctl
mkdir -p  ./build/out/bin/
go build -o ./build/out/bin/osdsdock github.com/opensds/opensds/cmd/osdsdock
mkdir -p  ./build/out/bin/
go build -o ./build/out/bin/osdslet github.com/opensds/opensds/cmd/osdslet
mkdir -p  ./build/out/bin/
go build -o ./build/out/bin/osdsctl github.com/opensds/opensds/cmd/osdsctl
```
Then the binary file will be generated to ```./build/out/bin```.

### Run etcd daemon in background.
```
cd $HOME/etcd-v3.2.0-linux-amd64
nohup sudo ./etcd > nohup.out 2> nohup.err < /dev/null &
```

### Install the ceph driver dependent packet.
```bash
sudo apt-get install -y librados-dev librbd-dev ceph-common
```

### Config the configuration file, you can refer to the following configuration.
To simplify the configuration, you can just run these two commands below:
```shell
cat > /etc/opensds/opensds.conf << OPENSDS_GLOABL_CONFIG_DOC
[osdslet]
api_endpoint = localhost:50040
graceful = True
log_file = /var/log/opensds/osdslet.log
socket_order = inc

[osdsdock]
api_endpoint = localhost:50050
log_file = /var/log/opensds/osdsdock.log
# Specify which backends should be enabled, sample,ceph,cinder,lvm and so on.
enabled_backends = ceph

[ceph]
name = ceph
description = Ceph E2E Test
driver_name = ceph
config_path = /etc/opensds/driver/ceph.yaml

[database]
endpoint = localhost:2379,localhost:2380
driver = etcd
OPENSDS_GLOABL_CONFIG_DOC

cat > /etc/opensds/driver/ceph.yaml <<OPENSDS_CEPH_DIRVER_CONFIG_DOC
configFile: /etc/ceph/ceph.conf
pool:
  "rbd":
    diskType: SSD
    AZ: default
OPENSDS_CEPH_DIRVER_CONFIG_DOC
```

## Run OpenSDS Service

### Start up the osdslet and osdsdock daemon. 
```bash
cd $GOPATH/src/github.com/opensds/opensds

sudo build/out/bin/osdslet
sudo build/out/bin/osdsdock
```
Or you can run them in background.
```bash
cd $GOPATH/src/github.com/opensds/opensds

sudo build/out/bin/osdslet -daemon
sudo build/out/bin/osdsdock -daemon
```

### Test the OpenSDS if it is work well by getting the api versions and the ceph storage pools.
```bash
sudo cp build/out/bin/osdsctl /usr/local/bin

export OPENSDS_ENDPOINT=http://127.0.0.1:50040
osdsctl pool list
```

## OpenSDS Usage Tutorial

### Create a default profile firstly.
```
osdsctl profile create '{"name": "default", "description": "default policy"}'
```

### Create a volume.
```
osdsctl volume create 1 --name=test-001
```

### List all volumes.
```
osdsctl volume list
```

### Delete the volume.
```
osdsctl volume delete <your_volume_id>
```
