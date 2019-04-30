# Docker build usage:
# 	docker build . -t opensdsio/opensds-filedock:latest
# Docker run usage:
# 	docker run -d --privileged=true --net=host -v /etc/opensds:/etc/opensds opensdsio/opensds-filedock:latest

FROM ubuntu:16.04

COPY osdsfiledock /usr/bin

# Install some packages before running command.
RUN apt-get update && apt-get install -y \
 librados-dev librbd-dev ceph-common lvm2 udev tgt \
 && rm -rf /var/lib/apt
RUN sed -i -e 's/udev_sync = 1/udev_sync = 0/g' /etc/lvm/lvm.conf \
 && sed -i -e 's/udev_rules = 1/udev_rules = 0/g' /etc/lvm/lvm.conf \
 && sed -i -e 's/use_lvmetad = 0/use_lvmetad =1/g' /etc/lvm/lvm.conf

# Define default command.
CMD ["/usr/bin/osdsfiledock"]
