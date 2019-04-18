# Docker build usage:
# 	docker build . -t opensdsio/opensds-controller:latest
# Docker run usage:
# 	docker run -d --net=host -v /etc/opensds:/etc/opensds opensdsio/opensds-controller:latest

FROM ubuntu:16.04
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY osdslet /usr/bin

# Define default command.
CMD ["/usr/bin/osdslet"]
