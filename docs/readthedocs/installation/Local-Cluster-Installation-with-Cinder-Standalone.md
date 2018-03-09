Here is a tutorial guiding users and new contributors to get familiar with [OpenSDS](https://github.com/opensds/opensds) by installing a simple local cluster and manage cinder standalone service. You can also use the ansible script to install automatically, see detail in [OpenSDS Local Cluster Installation through ansible](https://github.com/opensds/opensds/blob/master/contrib/ansible/README.md).

## Prepare
Before you start, please make sure you have all stuffs below ready:
- Ubuntu environment (suggest v16.04+).
- More than 30GB remaining disk.
- Make sure have access to the Internet.
- Some tools (```git```, ```docker```, ```docker-compose```) prepared.

## Step by Step Installation
### Install and start cinder standalone service
Install LVM2 and create volume group.
```shell
sudo apt-get install -y lvm2 thin-provisioning-tools
sudo modprobe dm_thin_pool

sudo truncate --size=10G cinder-volumes.img
# Get next available loop device
LD=$(sudo losetup -f)
sudo losetup $LD cinder-volumes.img
sudo sfdisk $LD << EOF
,,8e,,
EOF
sudo pvcreate $LD
sudo vgcreate cinder-volumes $LD
```
Install the cinderclient and brick extensions to do local-attaches.
```shell
git clone https://github.com/openstack/python-cinderclient
cd python-cinderclient
# Tested successfully in this version `ab0185bfc6e8797a35a2274c2a5ee03afb03dd60`
# git checkout -b ab0185bfc6e8797a35a2274c2a5ee03afb03dd60
sudo pip install -e .

git clone https://github.com/openstack/python-brick-cinderclient-ext.git
cd python-brick-cinderclient-ext
# Tested successfully in this version `a281e67bf9c12521ea5433f86cec913854826a33`
# git checkout -b a281e67bf9c12521ea5433f86cec913854826a33
sudo pip install -e .
```
Download source code from cinder repository.
```shell
git clone https://github.com/openstack/cinder

# Tested successfully in this version `7bbc95344d3961d0bf059252723fa40b33d4b3fe`
# cd cinder
# git checkout -b 7bbc95344d3961d0bf059252723fa40b33d4b3fe
```
Build cinder conatiner images.
```shell
cd contrib/block-box
make blockbox
```
Then, you should find these images by `docker images`.
```table
REPOSITORY          TAG         IMAGE ID            CREATED             SIZE
lvm-debian-cinder   latest      14ec2d893e56        37 seconds ago      380MB
debian-cinder       latest      729dc9012984        37 seconds ago      329MB
```
Now, you can start cinder service.
```shell
docker-compose up -d
```
Check whether the cinder service is normal.
```shell
source cinder.rc
cinder get-pools
+----------+--------------------+
| Property | Value              |
+----------+--------------------+
| name     | cinder-lvm@lvm#lvm |
+----------+--------------------+

```

### Install etcd
```shell
wget https://github.com/coreos/etcd/releases/download/v3.2.0/etcd-v3.2.0-linux-arm64.tar.gz

tar -xzf etcd-v3.2.0-linux-amd64.tar.gz
cd etcd-v3.2.0-linux-amd64
cp etcd etcdctl /usr/bin
```

### Install opensds
```shell
git clone -b development https://github.com/opensds/opensds.git

cd opensds
make
cp ./build/out/bin/* /usr/local/bin
```

### Set the opensds configuration file.
To simplify the configuration, you can just run these two commands below:
```shell
cat > /etc/opensds/opensds.conf << EOF
[osdslet]
api_endpoint = localhost:50040
graceful = True
log_file = $OPENSDS_LOG_DIR/osdslet.log
socket_order = inc

[osdsdock]
api_endpoint = localhost:50050
log_file = /var/log/opensds/osdsdock.log
enabled_backends = cinder

[ceph]
name = cinder
description = Cinder standalone service
driver_name = cinder
config_path = /etc/opensds/driver/cinder.yaml

[database]
endpoint = localhost:2379,localhost:2380
driver = etcd
EOF

cat > /etc/opensds/driver/cinder.yaml << EOF
authOptions:
  noAuth: true
  cinderEndpoint: "http://127.0.0.1:8776/v2"
  domainId: "Default"
  domainName: "Default"
  tenantId: "myproject"
  tenantName: "myproject"
pool:
  "cinder-lvm@lvm#lvm":
    AZ: nova
    thin: true
EOF
```

### Start opensds services.
```shell
setsid etcd
setsid osdslet
setsid osdsdock
```

## Testing
### Config osdsctl tool.
```shell
export OPENSDS_ENDPOINT=http://127.0.0.1:50040
osdsctl pool list
```


### Create a default profile firstly.
```shell
osdsctl profile create '{"name": "default", "description": "default policy"}'
```

### Create a volume.
```shell
osdsctl volume create 1 --name=test-001 --az nova
```

### List all volumes.
```shell
osdsctl volume list
```

### Delete the volume.
```shell
osdsctl volume delete <your_volume_id>
```

Hope you could enjoy it, and more suggestions are welcomed!
