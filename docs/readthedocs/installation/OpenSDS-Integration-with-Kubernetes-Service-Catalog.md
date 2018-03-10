In this tutorial, we will show how OpenSDS provide Ceph storage as a service for Kubernetes users through Service Catalog. Now hope you enjoy the trip!

## Pre-configuration

### Check it out about your os (very important)
Please NOTICE that the installation tutorial is tested on Ubuntu16.04, and we SUGGEST you follow our styles and use Ubuntu16.04+.

### Install some package dependencies
```
apt-get install gcc make libc-dev docker.io
```
If the docker command doesn't work, try to restart it:
```
sudo service docker stop
sudo nohup docker daemon -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock &
```

### Local Kubernetes Cluster Setup

* You can download and build k8s local cluster by executing commands blow:
```
mkdir -p /opt && cd /opt
git clone https://github.com/kubernetes/kubernetes.git -b v1.9.0
cd kubernetes && make
```

* Modify one line in `hack/local-up-cluster.sh` to enable PodPreset by appending `PodPreset` to ADMISSION_CONTROL:
```
ADMISSION_CONTROL=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount{security_admission},DefaultStorageClass,DefaultTolerationSeconds,GenericAdmissionWebhook,ResourceQuota,PodPreset
```

* Set up k8s cluster:
```
echo alias kubectl='/opt/kubernetes/cluster/kubectl.sh' >> /etc/profile
RUNTIME_CONFIG=settings.k8s.io/v1alpha1=true AUTHORIZATION_MODE=Node,RBAC hack/local-up-cluster.sh -O
kubectl get pod (check if k8s cluster running)
```

* Configure flexvolume plugin:
```
cd /opt && wget https://github.com/opensds/nbp/releases/download/v0.1.0/opensds-k8s-v0.1.0-linux-amd64.tar.gz
tar zxvf opensds-k8s-v0.1.0-linux-amd64.tar.gz
mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/opensds.io~opensds
cp opensds-k8s-v0.1.0-linux-amd64/flexvolume/opensds /usr/libexec/kubernetes/kubelet-plugins/volume/exec/opensds.io~opensd
```

### Install helm (from scipt)
To avoid some errors, please make sure your helm version is more than v2.7.0:
```
curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh
helm init
kubectl.sh get po -n kube-system (check if the till-deploy pod is running)
```

## Service Catalog Setup

### Tiller permissions

Tiller is the in-cluster server component of Helm. By default, 
`helm init` installs the Tiller pod into the `kube-system` namespace,
and configures Tiller to use the `default` service account.

Tiller will need to be configured with `cluster-admin` access to properly install
Service Catalog:

```console
kubectl create clusterrolebinding tiller-cluster-admin \
    --clusterrole=cluster-admin \
    --serviceaccount=kube-system:default
```

### Helm repository setup

Service Catalog is easily installed via a 
[Helm chart](https://github.com/kubernetes/helm/blob/master/docs/charts.md).

This chart is located in a
[chart repository](https://github.com/kubernetes/helm/blob/master/docs/chart_repository.md)
just for Service Catalog. Add this repository to your local machine:

```console
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
```

Then, ensure that the repository was successfully added:

```console
helm search service-catalog
```

You should see the following output:

```console
NAME           	VERSION	DESCRIPTION
svc-cat/catalog	x,y.z  	service-catalog API server and controller-manag...
```

Now that your cluster and Helm are configured properly, installing Service Catalog is simple:
```console
helm install svc-cat/catalog \
    --name catalog --namespace catalog
```

## OpenSDS Service Broker Configuration

* OpenSDS local cluster installation

For testing purposes you can deploy OpenSDS local cluster referring to the [OpenSDS Cluster Installation through Ansible](https://github.com/opensds/opensds/wiki/OpenSDS-Cluster-Installation-through-Ansible) wiki.

* Service broker

Firstly, please modify the endpoint IP:
```
git clone https://github.com/opensds/nbp.git
cd nbp/service-broker
vim charts/template/broker-deployment
```

Replace `0.0.0.0` with your host ip:
```yaml
containers:
- name: service-broker
  image: "opensdsio/service-broker:latest"
  imagePullPolicy: IfNotPresent
  args:
  - --port
  - ":8080"
  - --endpoint
  - "http://0.0.0.0:50040"
```

Then install it through helm:
```
helm install charts/ --name service-broker --namespace service-broker
kubectl get pod -n service-broker
kubectl get clusterservicebrokers,clusterserviceclasses,serviceinstances,servicebindings
```

## Start to work

* Create service broker

```
kubectl create -f examples/service-broker.yaml
kubectl get clusterservicebrokers,clusterserviceclasses,clusterserviceplans
```

* Create service instance

```
kubectl create ns opensds
kubectl create -f examples/service-instance.yaml -n opensds
kubectl get serviceinstances -n opensds -o yaml
```

* Create service instance binding

```
kubectl create -f examples/service-binding.yaml -n opensds
kubectl get servicebindings -n opensds -o yaml
kubectl describe  servicebindings -n opensdss
kubectl get secrets -n opensds
kubectl get secrets service-binding -o yaml -n opensds
kubectl get secrets service-binding -o yaml -n opensds | grep volumeId | awk  '{print $2}' | base64 -d && echo
```

* Creat service wordpress for testing

Modify the volume name of podpreset:
```
vi examples/podpreset-preset.yaml
kubectl create -f examples/podpreset-preset.yaml 
kubectl get podpreset
```
podpreset-preset.yaml
```yaml
apiVersion: settings.k8s.io/v1alpha1
kind: PodPreset
metadata:
  name: allow-database
spec:
  selector:
    matchLabels:
      role: frontend
  volumeMounts:
    - mountPath: /mnt/wordpress
      name: 938b481a-af44-44e4-9cc0-e13ac5c0cc5e
  volumes:
    - name: 938b481a-af44-44e4-9cc0-e13ac5c0cc5e
      flexVolume:
        driver: "opensds.io/opensds"
        fsType: "ext4"
```

Then create ```wordpress.yaml``` file:
```
kubectl create -f examples/wordpress.yaml

socat tcp-listen:8084,reuseaddr,fork tcp:10.0.0.124:8084
```
wordpress.yaml
```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: wordpress
spec:
  template:
    metadata:
      labels:
        app: wordpress
        role: frontend
    spec:
      containers:
      - name: wordpress
        image: wordpress:latest
        imagePullPolicy: IfNotPresent
        ports:
        - name: wordpress
          containerPort: 8084
---
apiVersion: v1
kind: Service
metadata:
  name: wordpress
spec:
  type: ClusterIP
  ports:
  - name: wordpress
    port: 8084
    targetPort: 8084
    protocol: TCP
  selector:
    app: wordpress
```

After all things done, you can visit your own blog by searching: ```http://service_cluster_ip:8084```!

## Clean it up

```
kubectl delete -f examples/podpreset-preset.yaml 

kubectl delete -f examples/wordpress.yaml

kubectl delete -n opensds  servicebindings service-binding

kubectl delete -n opensds serviceinstances service-instance

kubectl delete clusterservicebrokers service-broker

helm delete --purge service-broker

kubectl delete ns opensds service-broker

helm delete --purge catalog
```

## Ending

That's all the tutorial, thank you for watching it!
