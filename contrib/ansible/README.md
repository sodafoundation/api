# opensds-ansible
It's an installation tool of opensds through ansible.

## 1. How to install an opensds local cluster
This installation document assumes you have a clean Ubuntu16.04 environment, so if you have installed golang environment, please make sure some parameters configured at ```/etc/profile```:
```conf
export GOROOT=/usr/local/go
export GOPATH=$HOME/gopath
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

### Pre-config (Ubuntu 16.04)
Firstly download some system packages:
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
ssh-copy-id -i ~/.ssh/id_rsa.pub <remote_ip> # Choose an ip address which can login your target machine
```

### Download ansible tool
```bash
sudo add-apt-repository ppa:ansible/ansible # It doesn't matter if failed
sudo apt-get update
sudo apt-get install ansible
ansible --version # Make sure your ansible version is 2.x.x
```

### Download opensds source code
```bash
mkdir -p $GOPATH/src/github.com/opensds && cd $GOPATH/src/github.com/opensds
git clone https://github.com/opensds/opensds.git
cd opensds/contrib/ansible
```

### Configure opensds cluster variables:
##### System environment:
Since it's hard to configure your username, you need to configure your ```workplace``` at `group_vars/common.yml`:
```yaml
workplace: /home/your_username # Change this field according to your username, if you login as root, please configure this parameter to '/root'
```

##### LVM
If choose `lvm` as storage backend, you should modify `group_vars/osdsdock.yml`:
```yaml
enabled_backend: lvm # Change it according to your backend, currently support 'lvm', 'ceph', 'cinder'
pv_device: "your_pv_device_path" # Ensure this device existed and available if you choose lvm
vg_name: "specified_vg_name" # Specify a name randomly, but don't change it if you choose other backends
```
And modify ```group_vars/lvm/lvm.yaml```, change pool name same to `vg_name`:
```yaml
"vg001" # change pool name same to vg_name
```
##### Ceph
If choose `ceph` as storage backend, you should modify `group_vars/osdsdock.yml`:
```yaml
enabled_backend: ceph # Change it according to your backend, currently support 'lvm', 'ceph', 'cinder'
ceph_pool_name: "specified_pool_name" # Specify a name randomly, but don't change it if you choose other backends
```
And modify ```group_vars/ceph/ceph.yaml```, change pool name same to `ceph_pool_name`:
```yaml
"rbd" # change pool name same to ceph pool
```
Then you also need to configure two files under ```group_vars/ceph```: `all.yml` and `osds.yml`. And here is an example:

```group_vars/ceph/all.yml```:
```yml
ceph_origin: repository
ceph_repository: community
ceph_stable_release: luminous # Choose luminous as default version
public_network: "192.168.3.0/24" # Run 'ip -4 address' to check the ip address
cluster_network: "{{ public_network }}"
monitor_interface: eth1 # Change to your own network interface
```
```group_vars/ceph/osds.yml```:
```yml
devices:
    - '/dev/sda' # Ensure this device existed and available if you choose ceph
    - '/dev/sdb' # Ensure this device existed and available if you choose ceph
osd_scenario: collocated
```

##### Cinder
If choose `cinder` as storage backend, you should modify `group_vars/osdsdock.yml`:
```yaml
enabled_backend: cinder # Change it according to your backend, currently support 'lvm', 'ceph', 'cinder'
use_cinder_standalone: true # if use cinder standalone
```

And configure the auth and pool options to access cinder in `group_vars/cinder/cinder.yaml`. Don't need to do anything if use cinder standalone.

### Check if the hosts could be reached
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
cp $GOPATH/src/github.com/opensds/opensds/build/out/bin/osdsctl /usr/local/bin
export OPENSDS_ENDPOINT=http://127.0.0.1:50040
osdsctl pool list # Check if the pool resource is available
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


## 3. How to purge and clean opensds cluster

### Run opensds-ansible playbook to clean the environment
```bash
sudo ansible-playbook clean.yml -i local.hosts
```

### Run ceph-ansible playbook to clean ceph cluster if you deployed ceph
```bash
cd /tmp/ceph-ansible
sudo ansible-playbook infrastructure-playbooks/purge-cluster.yml -i ceph.hosts
```

Besides, you will also need to clean the logical partition on the physical block device, suggest using ```fdisk``` tool.

### Remove ceph-ansible source code (optionally)
```bash
cd ..
sudo rm -rf /tmp/ceph-ansible
```
