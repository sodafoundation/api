# Docker build usage:
# 	docker build . -t opensdsio/opensds-controller:v1beta
# Docker run usage:
# 	docker run -d --privileged --net=host -v /etc/opensds:/etc/opensds opensdsio/opensds-controller:v1beta

FROM ubuntu
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY osdslet /usr/bin

# Define default command.
CMD ["/usr/bin/osdslet"]
