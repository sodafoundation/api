## Pre-configuration
Before you start, some configurations are required:
```shell
export BackendType="sample" # 'sample' is the default option, currently also support 'lvm'

mkdir -p /etc/opensds && sudo cat > /etc/opensds/opensds.conf <<OPENSDS_GLOABL_CONFIG_DOC
[osdslet]
api_endpoint = 0.0.0.0:50040
graceful = True
log_file = /var/log/opensds/osdslet.log
socket_order = inc

[osdsdock]
api_endpoint = 0.0.0.0:50050
log_file = /var/log/opensds/osdsdock.log
# Enabled backend types, such as 'sample', 'lvm', 'ceph', 'cinder', etc.
enabled_backends = ${BackendType}

[sample]
name = sample
description = Sample backend for testing
driver_name = default

[ceph]
name = ceph
description = Ceph Test
driver_name = ceph
config_path = /etc/opensds/driver/ceph.yaml

[cinder]
name = cinder
description = Cinder Test
driver_name = cinder
config_path = /etc/opensds/driver/cinder.yaml

[lvm]
name = lvm
description = LVM Test
driver_name = lvm
config_path = /etc/opensds/driver/lvm.yaml

[database]
# Enabled database types, such as etcd, mysql, fake, etc.
driver = etcd
endpoint = 127.0.0.1:2379,127.0.0.1:2380
OPENSDS_GLOABL_CONFIG_DOC
```
If you choose `lvm` as backend, you need to make sure physical volume and volume group existed. Besides, you need to configure lvm driver.
```
sudo pvdisplay # Check if physical volume existed
sudo vgdisplay # Check if volume group existed

mkdir -p /etc/opensds/driver && sudo cat > /etc/opensds/driver/lvm.yaml <<OPENSDS_DRIVER_CONFIG_DOC
tgtBindIp: 0.0.0.0
pool:
  "vg001":
    diskType: SSD
    AZ: default
OPENSDS_DRIVER_CONFIG_DOC
```

## OpenSDS Service Installation
If you are a lazy one, just like me, you probably want to do this:(`docker-compose` required)
```
wget https://raw.githubusercontent.com/opensds/opensds/master/docker-compose.yml

docker-compose up
```

Or you can do this:
```
docker run -d --net=host -v /usr/share/ca-certificates/:/etc/ssl/certs quay.io/coreos/etcd:latest

docker run -d --net=host -v /etc/opensds:/etc/opensds opensdsio/opensds-controller:latest

docker run -d --net=host --privileged=true -v /etc/opensds:/etc/opensds opensdsio/opensds-dock:latest
```

If you are a smart guy, you probably need to configure your service ip and database endpoint:
```
export HostIP="your_real_ip"
docker run -d --net=host -v /usr/share/ca-certificates/:/etc/ssl/certs quay.io/coreos/etcd:latest \
 -name etcd0 \
 -advertise-client-urls http://${HostIP}:2379,http://${HostIP}:4001 \
 -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 \
 -initial-advertise-peer-urls http://${HostIP}:2380 \
 -listen-peer-urls http://0.0.0.0:2380 \
 -initial-cluster-token etcd-cluster-1 \
 -initial-cluster etcd0=http://${HostIP}:2380 \
 -initial-cluster-state new

docker run -d --net=host -v /etc/opensds:/etc/opensds opensdsio/opensds-controller:latest /usr/bin/osdslet --api-endpoint=0.0.0.0:50040 --db-endpoint=${HostIP}:2379,${HostIP}:2380

docker run -d --net=host --privileged=true -v /etc/opensds:/etc/opensds opensdsio/opensds-dock:latest /usr/bin/osdsdock --api-endpoint=0.0.0.0:50050 --db-endpoint=${HostIP}:2379,${HostIP}:2380
```

## Test

### Download cli tool.
```
curl -sSL https://raw.githubusercontent.com/opensds/opensds/master/osdsctl/bin/osdsctl | mv osdsctl /usr/local/bin/

export OPENSDS_ENDPOINT=http://127.0.0.1:50040
osdsctl pool list
```

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

After this is done, just enjoy it!
