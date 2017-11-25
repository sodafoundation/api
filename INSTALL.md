This installation document assumes that you use [CEPH](https://github.com/ceph/ceph) as the default driver,and use the [ETCD](https://github.com/coreos/etcd) as the default database.

## Pre-configuration

### Bootstrap
If you have a clean environment (suggest Ubuntu16.04+), please run the script
to install all dependencies of this project (except ceph cluster):
```
script/cluster/bootstrap.sh
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
```
echo '
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
description = Ceph Test
driver_name = ceph
config_path = /etc/opensds/driver/ceph.yaml

[database]
endpoint = localhost:2379,localhost:2380
driver = etcd
' >> /etc/opensds/opensds.conf

echo '
configFile: /etc/ceph/ceph.conf
pool:
  "rbd":
    diskType: SSD
    iops: 1000
    bandwidth: 1000
    AZ: default
' >> /etc/opensds/driver/ceph.yaml
```

## Run OpenSDS Service

### Start up the osdslet and osdsdock daemon. 
```bash
cd $GOPATH/src/github.com/opensds/opensds

sudo build/out/bin/osdslet
sudo build/out/bin/osdsock
```
Or you can run them in background.
```bash
cd $GOPATH/src/github.com/opensds/opensds

nohup sudo build/out/bin/osdsdock > nohup.out 2> nohup.err < /dev/null &
nohup sudo build/out/bin/osdslet > nohup.out 2> nohup.err < /dev/null &
```

### Test the OpenSDS if it is work well by getting the api versions and the ceph storage pools.
```bash
root@opensds-worker-1:~# curl localhost:50040| python -m json.tool
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   110  100   110    0     0  20715      0 --:--:-- --:--:-- --:--:-- 22000
[
    {
		"name":        "v1alpha",
		"description": "v1alpha version",
		"status":      "CURRENT",
		"updatedAt":   "2017-07-10T14:36:58.014Z"
    }
]

root@opensds-worker-1:~# curl localhost:50040/v1alpha/block/pools| python -m json.tool
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   317  100   317    0     0  26665      0 --:--:-- --:--:-- --:--:-- 28818
[
    {
        "createAt": "2017-10-17T10:46:29",
        "dockId": "16b435fe-45d9-563b-8c6c-2cc2c438ff7c",
        "freeCapacity": 2,
        "id": "0517f561-85b3-5f6a-a38d-8b5a08bff7df",
        "name": "rbd",
        "parameters": {
            "bandwidth": 0,
            "crushRuleset": "0",
            "diskType": "SSD",
            "iops": 1000,
            "redundancyType": "replicated",
            "replicateSize": "3"
        },
        "totalCapacity": 6,
        "updateAt": ""
    }
]
```

## OpenSDS Usage Tutorial

### Create a default profile firstly.
```
curl -X POST "http://localhost:50040/v1alpha/profiles" -H "Content-Type: application/json" -d '{"name": "default", "description": "default policy", "extra": {}}'
```

### Create a volume.
```
curl -X POST "http://localhost:50040/v1alpha/block/volumes" -H "Content-Type: application/json" -d '{"name": "test001", "description": "this is a test", "size": 1, "profileId": ""}'
```

### List all volumes.
```
curl http://localhost:50040/v1alpha/block/volumes| python -m json.tool
```

### Delete the volume.
```
curl -X DELETE "http://localhost:50040/v1alpha/block/volumes/volume_id" -H "Content-Type: application/json"
```
