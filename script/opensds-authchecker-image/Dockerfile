# Pull base image.
FROM ubuntu:16.04

COPY keystone.sh /keystone.sh
COPY sources.list /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y sudo vim git telnet net-tools && \ 
    mkdir -p /opt/stack/
 
