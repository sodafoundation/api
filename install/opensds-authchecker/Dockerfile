# Docker build usage:
#   docker build . -t opensdsio/opensds-authchecker:base
#   docker run -d --privileged=true --name opensds-authchecker-base opensdsio/opensds-authchecker:base "/sbin/init"
#   docker exec -it opensds-authchecker-base /keystone.sh
#   docker commit opensds-authchecker-base opensdsio/opensds-authchecker:latest
#   docker rm -f opensds-authchecker-base
# Docker run usage:
#   docker run -d --privileged=true --name opensds-authchecker opensdsio/opensds-authchecker:latest

FROM ubuntu:16.04
MAINTAINER Leon Wang <wanghui71leon@gmail.com>

COPY keystone.sh /keystone.sh
COPY entrypoint.sh /entrypoint.sh

# Install some packages before running command.
RUN apt-get update && apt-get install -y \
    sudo nano git telnet net-tools iptables gnutls-bin ca-certificates && \
    mkdir -p /opt/stack/

RUN ["chmod", "+x", "/keystone.sh", "/entrypoint.sh"]
