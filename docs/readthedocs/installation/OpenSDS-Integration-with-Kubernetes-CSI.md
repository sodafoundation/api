# Prerequisite ##

### ubuntu

* Version information

	```
	root@proxy:~# cat /etc/issue
	Ubuntu 16.04.2 LTS \n \l
	```

### docker

* Version information

	```
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

* Check golang version information

	```
	root@proxy:~# go version
	go version go1.9.2 linux/amd64
	```

* You can install golang by executing commands blow:

	```
	wget https://storage.googleapis.com/golang/go1.9.2.linux-amd64.tar.gz
	tar -C /usr/local -xzf go1.9.2.linux-amd64.tar.gz
	echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
	echo 'export GOPATH=$HOME/gopath' >> /etc/profile
	source /etc/profile
	```

### [Etcd](https://github.com/coreos/etcd.git) 
* You can install etcd by executing commands blow:
	```
	cd $HOME
	wget https://github.com/coreos/etcd/releases/download/v3.3.0/etcd-v3.3.0-linux-amd64.tar.gz
	tar -xzf etcd-v3.3.0-linux-amd64.tar.gz
	cd etcd-v3.3.0-linux-amd64
	sudo cp -f etcd etcdctl /usr/local/bin/
	```

### [kubernetes](https://github.com/kubernetes/kubernetes) local cluster

* You can startup the lastest k8s local cluster by executing commands blow:

	```
	cd $HOME
	git clone https://github.com/kubernetes/kubernetes.git
	cd $HOME/kubernetes
	git checkout v1.9.0
	make
	echo alias kubectl='$HOME/kubernetes/cluster/kubectl.sh' >> /etc/profile
	ALLOW_PRIVILEGED=true FEATURE_GATES=CSIPersistentVolume=true,MountPropagation=true RUNTIME_CONFIG="storage.k8s.io/v1alpha1=true" LOG_LEVEL=5 hack/local-up-cluster.sh
	```

	
### [opensds](https://github.com/opensds/opensds) local cluster

* For testing purposes you can deploy OpenSDS referring the [Local Cluster Installation with LVM](https://github.com/opensds/opensds/wiki/Local-Cluster-Installation-with-LVM) wiki.

## Testing steps ##

* Load some ENVs which is setted before.

	```
	source /etc/profile
	```

* Download nbp source code.
	
	using git clone
	```
	git clone https://github.com/opensds/nbp.git  $GOPATH/src/github.com/opensds/nbp
	```
	or using go get
	```
	go get -v  github.com/opensds/nbp/...
	```  

* Build opensds CSI plug-in docker image.

	```
	cd $GOPATH/src/github.com/opensds/nbp/
	make docker
	```
	
	CSI plug-in image named ```opensdsio/csiplugin``` can be found by ```docker images```.

* Configure opensds endpoint IP

	```
	vi csi/server/deploy/kubernetes/csi-configmap-opensdsplugin.yaml
	```

	The IP (127.0.0.1) should be replaced with the opensds actual endpoint IP.
	```yaml
	kind: ConfigMap
	apiVersion: v1
	    metadata:
	name: csi-configmap-opensdsplugin
	    data:
	    opensdsendpoint: http://127.0.0.1:50040
	```

* Create opensds CSI pods.

	```
	kubectl create -f csi/server/deploy/kubernetes
	```

	After this three pods can be found by ```kubectl get pods``` like below:

	- csi-provisioner-opensdsplugin
	- csi-attacher-opensdsplugin
	- csi-nodeplugin-opensdsplugin

	You can find more design details from
    [CSI Volume Plugins in Kubernetes Design Doc](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/storage/container-storage-interface.md)

* Create example nginx application

	```
	kubectl create -f csi/server/examples/kubernetes/nginx.yaml
	```

	This example will mount a opensds volume into ```/var/lib/www/html```.

	You can use the following command to inspect into nginx container to verify it.

	```
	docker exec -it <nginx container id> /bin/bash
	```

## Clean up steps ##

Clean up example nginx application and opensds CSI pods by the following commands.

```
kubectl delete -f csi/server/examples/kubernetes/nginx.yaml
kubectl delete -f csi/server/deploy/kubernetes
```
