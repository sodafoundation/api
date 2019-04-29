FROM ubuntu:xenial

RUN apt-get update && apt-get install -y \
  apt-transport-https \
  git \
  software-properties-common \
  uuid-runtime \
  wget

ARG CEPH_REPO_URL=https://download.ceph.com/debian-luminous/
RUN wget -q -O- 'https://download.ceph.com/keys/release.asc' | apt-key add -
RUN apt-add-repository "deb ${CEPH_REPO_URL} xenial main"

RUN add-apt-repository ppa:gophers/archive

RUN apt-get update && apt-get install -y \
  ceph \
  libcephfs-dev \
  librados-dev \
  librbd-dev \
  golang-1.10-go

# add user account to test permissions
RUN groupadd -g 1010 bob
RUN useradd -u 1010 -g bob -M bob

ENV GOPATH /go
WORKDIR /go/src/github.com/ceph/go-ceph
VOLUME /go/src/github.com/ceph/go-ceph

COPY micro-osd.sh /
COPY entrypoint.sh /
ENTRYPOINT /entrypoint.sh
