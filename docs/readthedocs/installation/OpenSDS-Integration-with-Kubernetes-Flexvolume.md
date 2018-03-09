## Prerequisite ##
### ubuntu
* Version information

	```bash
	root@proxy:~# cat /etc/issue
	Ubuntu 16.04.2 LTS \n \l
	```
### docker
* Version information

	```bash
	root@proxy:~# docker version
	Client:
	 Version:      1.12.6
	 API version:  1.24
	 Go version:   go1.6.2
	 Git commit:   78d1802
	 Built:        Tue Jan 31 23:35:14 2017
	 OS/Arch:      linux/amd64
	
	Server:
	 Version:      1.12.6
	 API version:  1.24
	 Go version:   go1.6.2
	 Git commit:   78d1802
	 Built:        Tue Jan 31 23:35:14 2017
	 OS/Arch:      linux/amd64
	```

### [golang](https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz) 
* Version information

	```bash
	root@proxy:~# go version
	go version go1.9.2 linux/amd64
	```

* You can install golang by executing commands blow:

	```bash
	wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
	tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz
	export PATH=$PATH:/usr/local/go/bin
	export GOPATH=$HOME/gopath
	```

### [kubernetes](https://github.com/kubernetes/kubernetes) local cluster
* Version information
	```bash
	root@proxy:~# kubectl version
	Client Version: version.Info{Major:"1", Minor:"9+", GitVersion:"v1.9.0-beta.0-dirty", GitCommit:"a0fb3baa71f1559fd42d1acd9cbdd8a55ab4dfff", GitTreeState:"dirty", BuildDate:"2017-12-13T09:22:09Z", GoVersion:"go1.9.2", Compiler:"gc", Platform:"linux/amd64"}
	Server Version: version.Info{Major:"1", Minor:"9+", GitVersion:"v1.9.0-beta.0-dirty", GitCommit:"a0fb3baa71f1559fd42d1acd9cbdd8a55ab4dfff", GitTreeState:"dirty", BuildDate:"2017-12-13T09:22:09Z", GoVersion:"go1.9.2", Compiler:"gc", Platform:"linux/amd64"}
	```
* You can startup the k8s local cluster by executing commands blow:

	```bash
	cd $HOME
	git clone https://github.com/kubernetes/kubernetes.git
	cd $HOME/kubernetes
	git checkout v1.9.0
	make
	echo alias kubectl='$HOME/kubernetes/cluster/kubectl.sh' >> /etc/profile
	RUNTIME_CONFIG=settings.k8s.io/v1alpha1=true AUTHORIZATION_MODE=Node,RBAC hack/local-up-cluster.sh -O
	```
**NOTE**:   
<div> Due to opensds using etcd as the database which is same with kubernetes so you should startup kubernetes firstly.
</div>

### [opensds](https://github.com/opensds/opensds) local cluster
* For testing purposes you can deploy OpenSDS referring the [Local Cluster Installation with LVM](https://github.com/opensds/opensds/wiki/Local-Cluster-Installation-with-LVM) wiki.

## Testing steps ##
* Load some ENVs which is setted before.

    ```bash
    source /etc/profile
    ```
* Download nbp source code.

    using git clone  
	```bash
	git clone https://github.com/opensds/nbp.git  $GOPATH/src/github.com/opensds/nbp
	```
	
	or using go get  
	```bash
	go get -v  github.com/opensds/nbp/...
	```  

* Build the FlexVolume.

	```bash
	cd $GOPATH/src/github.com/opensds/nbp/flexvolume
	go build -o opensds ./cmd/flex-plugin/
	```
	
    FlexVolume plugin binary is on the current directory.  


* Copy the OpenSDS FlexVolume binary file to k8s kubelet `volume-plugin-dir`.  
	if you don't specify the `volume-plugin-dir`, you can execute commands blow:

	```bash
	mkdir -p /usr/libexec/kubernetes/kubelet-plugins/volume/exec/opensds.io~opensds/
	cp $GOPATH/src/github.com/opensds/nbp/flexvolume/opensds /usr/libexec/kubernetes/kubelet-plugins/volume/exec/opensds.io~opensds/
	```  
	
	**NOTE**: 
	<div>
	OpenSDS FlexVolume will get the opensds api endpoint from the environment variable `OPENSDS_ENDPOINT`, if you don't specify it, the FlexVolume will use the default vaule: `http://127.0.0.1:50040`. if you want to specify the `OPENSDS_ENDPOINT` executing command `export OPENSDS_ENDPOINT=http://ip:50040` and restart the k8s local cluster.
</div>

* Build the provisioner docker image.

	```bash
	cd $GOPATH/src/github.com/opensds/nbp/opensds-provisioner
	make container
	```

* Create service account, role and bind them.
	```bash
	cd $GOPATH/src/github.com/opensds/nbp/opensds-provisioner/examples
	kubectl create -f serviceaccount.yaml
	kubectl create -f clusterrole.yaml
	kubectl create -f clusterrolebinding.yaml
	```

* Change the opensds endpoint IP in pod-provisioner.yaml   
The IP (192.168.56.106) should be replaced with the OpenSDS osdslet actual endpoint IP.
    ```yaml
    kind: Pod
    apiVersion: v1
    metadata:
      name: opensds-provisioner
    spec:
      serviceAccount: opensds-provisioner
      containers:
        - name: opensds-provisioner
          image: opensdsio/opensds-provisioner
          securityContext:
          args:
            - "-endpoint=http://192.168.56.106:50040" # should be replaced
          imagePullPolicy: "IfNotPresent"
    ```

* Create provisioner pod.
	```bash
	kubectl create -f pod-provisioner.yaml
	```
	
    Execute `kubectl get pod` to check if the opensds-provisioner is ok.
    ```bash
    root@nbp:~/go/src/github.com/opensds/nbp/opensds-provisioner/examples# kubectl get pod
    NAME                  READY     STATUS    RESTARTS   AGE
    opensds-provisioner   1/1       Running   0          42m
    ```
* You can use the following cammands to test the OpenSDS FlexVolume and Proversioner functions.

    Create storage class.
	```bash
	kubectl create -f sc.yaml              # Create StorageClass
	```
	Execute `kubectl get sc` to check if the storage class is ok. 
	```bash
	root@nbp:~/go/src/github.com/opensds/nbp/opensds-provisioner/examples# kubectl get sc
    NAME                 PROVISIONER               AGE
    opensds              opensds/nbp-provisioner   46m
    standard (default)   kubernetes.io/host-path   49m
	```
	Create PVC.
	```bash
	kubectl create -f pvc.yaml             # Create PVC
	```
	Execute `kubectl get pvc` to check if the pvc is ok. 
	```bash
	root@nbp:~/go/src/github.com/opensds/nbp/opensds-provisioner/examples# kubectl get pvc
    NAME          STATUS    VOLUME                                 CAPACITY   ACCESS MODES   STORAGECLASS   AGE
    opensds-pvc   Bound     731da41e-c9ee-4180-8fb3-d1f6c7f65378   1Gi        RWO            opensds        48m

	```
	Create busybox pod.
	
	```bash
	kubectl create -f pod-application.yaml # Create busybox pod and mount the block storage.
	```
	Execute `kubectl get pod` to check if the busybox pod is ok. 
    ```bash
    root@nbp:~/go/src/github.com/opensds/nbp/opensds-provisioner/examples# kubectl get pod
    NAME                  READY     STATUS    RESTARTS   AGE
    busy-pod              1/1       Running   0          49m
    opensds-provisioner   1/1       Running   0          50m
    ```
	Execute the `findmnt|grep opensds` to confirm whether the volume has been provided.
	If there is some thing that goes wrong, you can check the log files in directory `/var/log/opensds`.

## Clean up steps ##

```
kubectl delete -f pod-application.yaml
kubectl delete -f pvc.yaml
kubectl delete -f sc.yaml

kubectl delete -f pod-provisioner.yaml
kubectl delete -f clusterrolebinding.yaml
kubectl delete -f clusterrole.yaml
kubectl delete -f serviceaccount.yaml
```
