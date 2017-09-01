# Docker build usage: docker build . -t opensds/opensds-controller:v1alpha
# Docker run usage:	docker run -it --net=host -v /var/log/opensds:/var/log/opensds
#                                             opensds/opensds-controller:v1alpha

FROM ubuntu
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY osdslet /usr/bin

# Define default command.
CMD ["/usr/bin/osdslet"]