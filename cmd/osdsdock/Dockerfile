# Docker build usage:
# 	docker build . -t opensdsio/opensds-dock:v1beta
# Docker run usage:
# 	docker run -d --privileged --net=host -v /etc/opensds:/etc/opensds opensdsio/opensds-dock:v1beta

FROM golang:1.7.1
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY osdsdock /usr/bin

# Install some packages before running command.
RUN apt-get update && apt-get install -y librados-dev librbd-dev ceph-common lvm2

# Define default command.
CMD ["/usr/bin/osdsdock"]
