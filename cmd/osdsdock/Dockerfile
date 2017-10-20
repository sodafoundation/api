# Docker build usage: docker build . -t opensds/opensds-dock:v1alpha
# Docker run usage:	docker run -it --net=host -v /var/log/opensds:/var/log/opensds \
#											 -v /etc/opensds:/etc/opensds \
#											 -v /etc/ceph:/etc/ceph	\
#	                                          opensds/opensds-dock:v1alpha

FROM golang:1.7.1
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY osdsdock /usr/bin

# Install some packages before running command.
RUN apt-get update && apt-get install -y librados-dev librbd-dev ceph-common

# Define default command.
CMD ["/usr/bin/osdsdock"]
