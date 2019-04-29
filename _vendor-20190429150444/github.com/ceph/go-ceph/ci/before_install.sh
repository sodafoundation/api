#!/bin/bash

set -e
set -x

sudo apt-get install -y python-virtualenv

# ceph-deploy and ceph

WORKDIR=$HOME/workdir
mkdir $WORKDIR
pushd $WORKDIR

ssh-keygen -f $HOME/.ssh/id_rsa -t rsa -N ''
cat ~/.ssh/id_rsa.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

git clone git://github.com/ceph/ceph-deploy
pushd ceph-deploy
./bootstrap
./ceph-deploy install --release ${CEPH_RELEASE} `hostname`
./ceph-deploy pkg --install librados-dev `hostname`
./ceph-deploy pkg --install librbd-dev `hostname`
./ceph-deploy pkg --install libcephfs-dev `hostname`
popd # ceph-deploy

popd # workdir
