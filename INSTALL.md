This installation document assumes that you use [CEPH](https://github.com/ceph/ceph) as the default driver,and use the [ETCD](https://github.com/coreos/etcd) as the default database.So,
before you start up the OpenSDS,the ceph cluster and etcd should be start up firstly.Then you can excute the followings step by step.

## Pre-configuration

### Download and install Golang
```
wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
tar xvf go1.9.linux-amd64.tar.gz -C /usr/local/
mkdir -p $HOME/gopath/src
mkdir -p $HOME/gopath/bin
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/gopath/bin' >> /etc/profile
echo 'export GOPATH=$HOME/gopath' >> /etc/profile
source /etc/profile
go version (check if go has been installed)
```

### Download the OpenSDS source code and change it to the development branch.
```
mkdir -p $HOME/gopath/src/github.com/opensds
cd $HOME/gopath/src/github.com/opensds
git clone https://github.com/opensds/opensds.git -b development
```

### Install the ceph driver dependent packet.
```bash
sudo apt-get install -y librados-dev librbd-dev ceph-common
```

### Enter the source code directory and build the source code.
```
cd opensds/ && make
```

Then you will find the output below:
```bash
go get github.com/opensds/opensds/cmd/osdslet
go get github.com/opensds/opensds/cmd/osdsdock
go get github.com/opensds/opensds/cmd/osdsctl
mkdir -p  ./build/out/bin/
go build -o ./build/out/bin/osdsdock github.com/opensds/opensds/cmd/osdsdock
mkdir -p  ./build/out/bin/
go build -o ./build/out/bin/osdslet github.com/opensds/opensds/cmd/osdslet
```
Then the binary file will be generated to ```./build/out/bin```

### Config the configuration file, you can refer to the following configuration
```conf
[osdsdock]
enabled_backends = ceph

[ceph]
name = ceph
description = Ceph Test
driver_name = ceph
ceph_config = /etc/opensds/driver/ceph.yaml
```
If you want to test ceph driver,you should config the ```/etc/opensds/driver/ceph.yaml``` too, you can refer to the following configuration:
```yaml
configFile: /etc/ceph/ceph.conf
pool:
  "rbd":
    diskType: SSD
    iops: 1000
    bandWitdh: 1G
```

The configuration process would be as follows:(```Suppose you are under opensds root directory```)
```
vim examples/opensds.conf (modify sample opensds config file)
vim examples/driver/ceph.yaml (modify sample ceph backend config file)
sudo mkdir -p /etc/opensds/driver
cp examples/opensds.conf /etc/opensds/
cp examples/driver/ceph.yaml /etc/opensds/driver/
```

## Run OpenSDS Service

### Start up the osdslet and osdsdock daemon. 
```bash
sudo ./osdslet

sudo ./osdsock
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