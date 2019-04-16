# Docker build usage:
# 	docker build . -t opensdsio/opensds-authchecker:latest
# Docker run usage:
# 	docker run -d --privileged=true --name opensds-authchecker opensdsio/opensds-authchecker:latest

FROM ubuntu:16.04
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY keystone.sh /keystone.sh
COPY entrypoint.sh /entrypoint.sh

# Install some packages before running command.
RUN apt-get update && apt-get install -y --no-install-recommends \
    sudo git telnet net-tools iptables gnutls-bin && \
    rm -rf /var/lib/apt/lists/*

# Build keystone service.
RUN sudo bash /keystone.sh

# Define default command.
ENTRYPOINT ["/entrypoint.sh"]
