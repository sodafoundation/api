# Install OpenSDS on an existing Kubernetes cluster

This tutorial assumes that you already have an existing Kubernetes cluster with
kube-dns service enabled. If there is some DNS problems with your Kubernetes
cluster, please refer to [here](https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/)
for debugging resolution.

## Prepare
Before you start, please make sure you have all stuffs below ready:
- Kubernetes cluster (suggest v1.13.x or later).
- More than 30GB remaining disk.
- Make sure have access to the Internet.

Here is some guide

## Step by Step Installation
### Configuration
Firstly, you need to configure some global files with command below:
```shell
export BackendType="sample" # 'sample' is the default option, currently also support 'lvm'

mkdir -p /etc/opensds && sudo cat > /etc/opensds/opensds.conf <<OPENSDS_GLOABL_CONFIG_DOC
[osdsapiserver]
api_endpoint = 0.0.0.0:50040
dns_endpoint = apiserver.opensds.svc.cluster.local:50040
auth_strategy = noauth
# If https is enabled, the default value of cert file
# is /opt/opensds-security/opensds/opensds-cert.pem,
# and key file is /opt/opensds-security/opensds/opensds-key.pem
https_enabled = False
beego_https_cert_file =
beego_https_key_file =
# Encryption and decryption tool. Default value is aes.
password_decrypt_tool = aes

[osdslet]
api_endpoint = 0.0.0.0:50049
dns_endpoint = controller.opensds.svc.cluster.local:50049

[osdsdock]
api_endpoint = 0.0.0.0:50050
dns_endpoint = dock.opensds.svc.cluster.local:50050
# Choose the type of dock resource, only support 'provisioner' and 'attacher'.
dock_type = provisioner
# Specify which backends should be enabled, sample,ceph,cinder,lvm and so on.
enabled_backends = $BackendType

[sample]
name = sample
description = Sample Test
driver_name = sample

[lvm]
name = lvm
description = LVM Test
driver_name = lvm
config_path = /etc/opensds/driver/lvm.yaml
host_based_replication_driver = DRBD

[database]
dns_endpoint = db.opensds.svc.cluster.local:2379,db.opensds.svc.cluster.local:2380
driver = etcd
OPENSDS_GLOABL_CONFIG_DOC
```

If you choose `lvm` as backend, you need to make sure physical volume and volume group existed. Besides, you need to configure lvm driver.
```
sudo pvdisplay # Check if physical volume existed
sudo vgdisplay # Check if volume group existed

mkdir -p /etc/opensds/driver && sudo cat > /etc/opensds/driver/lvm.yaml <<OPENSDS_DRIVER_CONFIG_DOC
tgtBindIp: 127.0.0.1
tgtConfDir: /etc/tgt/conf.d
pool:
  {{ volume_group_name }}:
    storageType: block
    availabilityZone: default
    extras:
      dataStorage:
        provisioningPolicy: Thin
        isSpaceEfficient: false
      ioConnectivity:
        accessProtocol: iscsi
        maxIOPS: 7000000
        maxBWS: 600
      advanced:
        diskType: SSD
        latency: 5ms
OPENSDS_DRIVER_CONFIG_DOC
```

### OpenSDS service deployment
Thanks to the orchesration feature of Kubernetes, now you can deploy the whole
OpenSDS cluster simply using these commands:
```shell
kubectl create ns opensds
kubectl create -f install/kubernetes/opensds-all.yaml

kubectl get po -n opensds # Check if all pods created
kubectl get svc -n opensds # Check if all services created
```

If everything works well, you will find some info like below:
```shell
root@ubuntu:~# kubectl get po -n opensds
NAME                                 READY   STATUS    RESTARTS   AGE
apiserver-v1beta-5455ddb848-bzxjt    1/1     Running   0          117m
authchecker-v1-7f6866b858-c5w49      1/1     Running   0          117m
controller-v1beta-77c566d4d4-hhdvz   1/1     Running   0          117m
dashboard-v1beta-6bdbc8d5d9-wdg28    1/1     Running   0          117m
db-v1-5f859b7fd9-r55z7               1/1     Running   0          117m
dock-v1beta-77ff5f5d55-c58bl         1/1     Running   0          117m
root@ubuntu:~# kubectl get svc -n opensds
NAME          TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)             AGE
apiserver     ClusterIP   10.0.0.55    <none>        50040/TCP           117m
authchecker   ClusterIP   10.0.0.165   <none>        35357/TCP           117m
controller    ClusterIP   10.0.0.232   <none>        50049/TCP           117m
dashboard     ClusterIP   10.0.0.67    <none>        8088/TCP            117m
db            ClusterIP   10.0.0.37    <none>        2379/TCP,2380/TCP   117m
dock          ClusterIP   10.0.0.111   <none>        50050/TCP           117m
```

## Test work
### Download cli tool.
```
wget https://github.com/opensds/opensds/releases/download/v0.5.1/opensds-hotpot-v0.5.1-linux-amd64.tar.gz 
tar zxvf opensds-hotpot-v0.5.1-linux-amd64.tar.gz 
cp opensds-hotpot-v0.5.1-linux-amd64/bin/* /usr/local/bin
chmod 755 /usr/local/bin/osdsctl

export OPENSDS_ENDPOINT=http://{{ apiserver_cluster_ip }}:50040
export OPENSDS_AUTH_STRATEGY=noauth
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
