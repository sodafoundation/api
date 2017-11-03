This installation document assumes that you use [CEPH](https://github.com/ceph/ceph) as the default driver,and use the [ETCD](https://github.com/coreos/etcd) as the default database.So,
before you start up the OpenSDS,the ceph cluster and etcd should be start up firstly.Then you can excute the followings step by step.

* Download the OpenSDS source code and change it to the development branch.

```bash
git clone https://github.com/opensds/opensds.git -b development
```

* Install the ceph driver dependent packet.

```bash
sudo apt-get install -y librados-dev librbd-dev
```

* Enter the source code directory and build the source code.

```bash
root@opensds-worker-1:~# cd opensds/
root@opensds-worker-1:~/opensds# make
go get github.com/opensds/opensds/cmd/osdslet
go get github.com/opensds/opensds/cmd/osdsdock
mkdir -p  ./build/out/bin/
go build -o ./build/out/bin/osdsdock github.com/opensds/opensds/cmd/osdsdock
mkdir -p  ./build/out/bin/
go build -o ./build/out/bin/osdslet github.com/opensds/opensds/cmd/osdslet
```
* The binary file will be generated to 

```
./build/out/bin
```
* Config the configuration file, you can refer to the following configuration

```
[osdsdock]
enabled_backends = ceph
ceph_config = /etc/opensds/driver/ceph.yaml

[ceph]
name = ceph
description = Ceph Test
endpoint = 127.0.0.1
driver_name = ceph
```
If you want to test ceph driver,you should config the ```/etc/opensds/driver/ceph.yaml``` too, you can refer to the following configuration

```
configFile: /etc/ceph/ceph.conf
pool:
  "rbd":
    diskType: SSD
    iops: 1000
    bandWitdh: 1G
```

* Start up the osdslet and osdsdock. 

```bash
./osdslet

./osdsock
```
* Test the OpenSDS if it is work well by getting the api versions and the ceph storage pools.

```
root@opensds-worker-1:~# curl localhost:50040/api| python -m json.tool
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   110  100   110    0     0  20715      0 --:--:-- --:--:-- --:--:-- 22000
[
    {
        "description": "v1alpha version",
        "name": "v1alpha",
        "status": "CURRENT",
        "updatedAt": "2017-07-10T14:36:58.014Z"
    }
]

root@opensds-worker-1:~# curl localhost:50040/api/v1alpha/block/pools| python -m json.tool
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
