# Docker build usage: docker build . -t leonwanghui/opensds-dock:v1alpha1
# Docker run usage:	docker run -it --net=host -v /var/log/opensds:/var/log/opensds \
#											 -v /etc/opensds:/etc/opensds \
#											 -v /etc/ceph:/etc/ceph	\
#	                                          leonwanghui/opensds-dock:v1alpha1

FROM golang:1.7.1
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

ENV CEPH_VERSION jewel

COPY osdsdock /usr/bin

RUN echo deb http://download.ceph.com/debian-$CEPH_VERSION/ jessie main | tee /etc/apt/sources.list.d/ceph-$CEPH_VERSION.list

# Running wget with no certificate checks, alternatively ssl-cert package should be installed
RUN wget --no-check-certificate -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | apt-key add - \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        ceph \
        librados-dev \
        librbd-dev \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Define default command.
CMD ["/usr/bin/osdsdock"]