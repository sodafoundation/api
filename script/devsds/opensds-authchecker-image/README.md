# How to build an image of opensdsio/opensds-authchecker?

## Build opensdsio/opensds-authchecker:base from the dockerfile
Execute the following command: docker build -t opensdsio/opensds-authchecker:base ./

## Build opensdsio/opensds-authchecker:latest from the opensdsio/opensds-authchecker:base
1. Execute the following command: docker network create --subnet=173.18.0.0/16 opensds-authchecker-network
1. Execute the following command: docker run -d  --privileged=true  --net opensds-authchecker-network --ip 173.18.0.2 --name opensds-authchecker opensdsio/opensds-authchecker:base "/sbin/init"
2. Execute the following command: docker exec -it opensds-authchecker /bin/bash
3. Execute the following command: sudo bash ./keystone.sh
4. Execute the following command: exit
5. Execute the following command: docker commit opensds-authchecker opensdsio/opensds-authchecker:latest
