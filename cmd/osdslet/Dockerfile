# Docker build usage: docker build . -t opensds/opensds-controller:v1beta
# Docker run usage:	docker run -it --net=host -v /var/log/opensds:/var/log/opensds
#                                             opensds/opensds-controller:v1beta

FROM ubuntu
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY osdslet /usr/bin

# Define default command.
CMD ["/usr/bin/osdslet"]