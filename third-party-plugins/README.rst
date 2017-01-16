.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

1. How does OpenSDS work with Docker?

  1.1 Start OpenSDS service:
	etcd
	
	go run goproj/opensds/src/main/service.go
	
	fuxi-server --config-file /etc/fuxi/fuxi.conf --log-dir /var/log/ --log-file fuxi.log
	
  1.2 Show the original volume status:
	cinder list (Openstack node)
	
	docker volume ls (OpenSDS node)
	
  1.3 Create a container with persistent volume:
	docker volume create --driver fuxi --name test-docker --opt size=3
	
	docker run -d -it --name test -v test-docker:/data --volume-driver fuxi ubuntu /bin/bash
	
	docker ps -a
  1.4 Show the volume status in Cinder:
	cinder list

2. How does OpenSDS work with k8s?

  2.1 Start OpenSDS service:
	etcd
	
	go run goproj/opensds/src/main/service.go
	
	fuxi-server --config-file /etc/fuxi/fuxi.conf --log-dir /var/log/ --log-file fuxi.log
	
  2.2 Show the original volume and Pod status:
	cinder list (Openstack node)
	
	kubectl get pods (k8s node)
	
  2.3 Create a Pod with persistent volume:
	vi /opt/k8s/yml/test-fuxi/fuxi-pod4.yml (k8s node)
	
	kubectl create -f /opt/k8s/yml/test-fuxi/fuxi-pod4.yml (k8s node)
	
	docker ps -a (OpenSDS node)
	
  2.4 Show the volume status in Cinder:
	cinder list
