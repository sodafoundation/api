# Local-Cluster-Installation-through-Ansible

## 1. How to install an opensds local cluster
### Pre-config (Ubuntu 16.04)
First download some system packages:
```
sudo apt-get install -y openssh-server git make gcc
```
Then config ```/etc/ssh/sshd_config``` file and change one line:
```conf
PermitRootLogin yes
```
Next generate ssh-token:
```bash
ssh-keygen -t rsa
ssh-copy-id -i ~/.ssh/id_rsa.pub <ip_address> # IP address of the target machine of the installation
```

### Install docker
If use a standalone cinder as backend, you also need to install docker to run cinder service. Please see the [docker installation document](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/) for details.

### Install ansible tool
To install ansible, you can run `install_ansible.sh` directly or input these commands below:
```bash
sudo add-apt-repository ppa:ansible/ansible # This step is needed to upgrade ansible to version 2.4.2 which is required for the ceph backend.
sudo apt-get update
sudo apt-get install ansible
ansible --version # Ansible version 2.4.2 or higher is required for ceph; 2.0.0.2 or higher is needed for other backends.
```

### Download opensds-ansible release file
```bash
export OPENSDS_RELEASE=v0.1.4 # The version MUST be v0.1.4 at least
wget https://github.com/opensds/opensds/releases/download/$OPENSDS_RELEASE/opensds-ansible-$OPENSDS_RELEASE-linux-amd64.tar.gz
tar zxvf opensds-ansible-$OPENSDS_RELEASE-linux-amd64.tar.gz
cd opensds-ansible-$OPENSDS_RELEASE-linux-amd64
```

### Configure opensds cluster variables:
##### System environment:
Configure these variables below in `group_vars/common.yml`:
```yaml
opensds_release: v0.1.4 # The version should be at least v0.1.4.
nbp_release: v0.1.0 # The version should be at least v0.1.0.

container_enabled: <false_or_true>
```

If you want to integrate OpenSDS with cloud platform (for example k8s), please modify `nbp_plugin_type` variable in `group_vars/common.yml`:
```yaml
nbp_plugin_type: standalone # standalone is the default integration way, but you can change it to 'csi', 'flexvolume'
```

#### Database configuration
Currently OpenSDS adopts `etcd` as database backend, and the default db endpoint is `localhost:2379,localhost:2380`. But to avoid some conflicts with existing environment (k8s local cluster), we suggest you change the port of etcd cluster in `group_vars/osdsdb.yml`:
```yaml
db_endpoint: localhost:62379,localhost:62380

etcd_host: 127.0.0.1
etcd_port: 62379
etcd_peer_port: 62380
```

##### LVM
If `lvm` is chosen as storage backend, modify `group_vars/osdsdock.yml`:
```yaml
enabled_backend: lvm # Change it according to the chosen backend. Supported backends include 'lvm', 'ceph', and 'cinder'
pv_devices: # Specify block devices and ensure them existed if you choose lvm
  #- /dev/sdc
  #- /dev/sdd
vg_name: "specified_vg_name" # Specify a name for VG if choosing lvm
```

Modify ```group_vars/lvm/lvm.yaml```, change pool name to be the same as `vg_name` above:
```yaml
tgtBindIp: 127.0.0.1 # change tgtBindIp to your real host ip, run 'ifconfig' to check
pool:
  "vg001" # change pool name to be the same as vg_name
```

Besides, Change `tgtBindIp` variable in `group_vars/lvm/lvm.yaml` to your real host ip.

##### Ceph
If `ceph` is chosen as storage backend, modify `group_vars/osdsdock.yml`:
```yaml
enabled_backend: ceph # Change it according to the chosen backend. Supported backends include 'lvm', 'ceph', and 'cinder'.
ceph_pools: # Specify pool name randomly if choosing ceph
  - rbd
  #- ssd
  #- sas
```

Modify ```group_vars/ceph/ceph.yaml```, change pool name to be the same as `ceph_pool_name`. But if you enable multiple pools, please append the current pool format:
```yaml
"rbd" # change pool name to be the same as ceph pool
```

Configure two files under ```group_vars/ceph```: `all.yml` and `osds.yml`. Here is an example:

```group_vars/ceph/all.yml```:
```yml
ceph_origin: repository
ceph_repository: community
ceph_stable_release: luminous # Choose luminous as default version
public_network: "192.168.3.0/24" # Run 'ip -4 address' to check the ip address
cluster_network: "{{ public_network }}"
monitor_interface: eth1 # Change to the network interface on the target machine
```
```group_vars/ceph/osds.yml```:
```yml
devices: # For ceph devices, append ONE or MULTIPLE devices like the example below:
    - '/dev/sda' # Ensure this device exists and available if ceph is chosen
    - '/dev/sdb' # Ensure this device exists and available if ceph is chosen
osd_scenario: collocated
```

##### Cinder
If `cinder` is chosen as storage backend, modify `group_vars/osdsdock.yml`:
```yaml
enabled_backend: cinder # Change it according to the chosen backend. Supported backends include 'lvm', 'ceph', and 'cinder'

# Use block-box install cinder_standalone if true, see details in:
use_cinder_standalone: true
# If true, you can configure cinder_container_platform,  cinder_image_tag,
# cinder_volume_group.

# Default: debian:stretch, and ubuntu:xenial, centos:7 is also supported.
cinder_container_platform: debian:stretch
# The image tag can be arbitrarily modified, as long as follow the image naming
# conventions, default: debian-cinder
cinder_image_tag: debian-cinder
# The cinder standalone use lvm driver as default driver, therefore `volume_group`
# should be configured, the default is: cinder-volumes. The volume group will be
# removed when use ansible script clean environment.
cinder_volume_group: cinder-volumes
```

Configure the auth and pool options to access cinder in `group_vars/cinder/cinder.yaml`. Do not need to make additional configure changes if using cinder standalone.

### Check if the hosts can be reached
```bash
sudo ansible all -m ping -i local.hosts
```

### Run opensds-ansible playbook to start deploy
```bash
sudo ansible-playbook site.yml -i local.hosts
```

## 2. How to test opensds cluster

### Configure opensds CLI tool
```bash
sudo cp $GOPATH/src/github.com/opensds/opensds/build/out/bin/osdsctl /usr/local/bin
export OPENSDS_ENDPOINT=http://127.0.0.1:50040
export OPENSDS_AUTH_STRATEGY=noauth
osdsctl pool list # Check if the pool resource is available
```

### Create a default profile first.
```
osdsctl profile create '{"name": "default", "description": "default policy"}'
```

### Create a volume.
```
osdsctl volume create 1 --name=test-001
```
For cinder, az needs to be specified.
```
osdsctl volume create 1 --name=test-001 --az nova
```

### List all volumes.
```
osdsctl volume list
```

### Delete the volume.
```
osdsctl volume delete <your_volume_id>
```


## 3. How to purge and clean opensds cluster

### Run opensds-ansible playbook to clean the environment
```bash
sudo ansible-playbook clean.yml -i local.hosts
```

### Run ceph-ansible playbook to clean ceph cluster if ceph is deployed
```bash
cd /opt/ceph-ansible
sudo ansible-playbook infrastructure-playbooks/purge-cluster.yml -i ceph.hosts
```

In addition, clean up the logical partition on the physical block device used by ceph, using the ```fdisk``` tool.

### Remove ceph-ansible source code (optional)
```bash
cd ..
sudo rm -rf /opt/ceph-ansible
```
