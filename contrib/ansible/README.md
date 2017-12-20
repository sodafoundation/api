# opensds-ansible
It's an installation tool of opensds through ansible.

## How to install an opensds local cluster

### Pre-config (Ubuntu 16.04)
Firstly download some ssh packages:
```
sudo apt-get install openssh-server
```
Then config ```/etc/ssh/sshd_config``` file and change one line:
```conf
PermitRootLogin yes
```
Next generate ssh-token:
```bash
ssh-keygen -t rsa
ssh-copy-id -i ~/.ssh/id_rsa.pub <romte_ip>
```

### Download ansible tool
```bash
sudo add-apt-repository ppa:ansible/ansible
sudo apt-get update
sudo apt-get install ansible
```

### Configure opensds cluster variables:
modify ```group_vars/osdsdock.yml```:
```yaml
enabled_backend: lvm # Change it according to your backend, currently support 'lvm', 'ceph'
pv_device: "your_pv_device_path" # Ensure this device existed if you choose lvm
vg_name: "specified_vg_name" # Specify a name randomly
```
modify ```group_vars/lvm/lvm.yaml``` if you specify lvm as your backend:
```yaml
    "vg001" # change pool name same to vg_name, but don't change it if you choose ceph backend
```
modify ```group_vars/ceph/ceph.yaml``` if you specify ceph as your backend:
```yaml
    "rbd" # change pool name same to ceph pool, but don't change it if you choose lvm backend
```

If you choose ceph, you also need to configure two files under ```group_vars/ceph```: all.yml and osds.yml. And here is an example:

modify ```group_vars/ceph/all.yml```:
```yml
ceph_origin: repository
ceph_repository: community
ceph_stable_release: luminous
public_network: "192.168.3.0/24"
cluster_network: "{{ public_network }}"
monitor_interface: eth1
```
modify ```group_vars/ceph/osds.yml```:
```yml
devices:
    - '/dev/sda' # Ensure this device existed if you choose ceph
    - '/dev/sdb' # Ensure this device existed if you choose ceph
osd_scenario: collocated
```

### Check if the hosts could be reached
```bash
ansible all -m ping -i local.hosts
```

### Check if the hosts could be reached
```bash
ansible-playbook site.yml -i local.hosts
```


## How to purge and clean opensds cluster

### Run opensds-ansible playbook to clean the environment
```bash
ansible-playbook clean.yml -i local.hosts
```

### Run ceph-ansible playbook to clean ceph cluster if you deployed ceph
```bash
cd /root/ceph-ansible
ansible-playbook infrastructure-playbooks/purge-cluster.yml -i ceph.hosts
```

### Remove ceph-ansible source code (optionally)
```bash
cd ..
rm -rf ceph-ansible
```