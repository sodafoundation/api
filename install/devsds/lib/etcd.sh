#!/bin/bash

# Copyright 2017 The OpenSDS Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# A set of helpers for starting/running etcd for tests
_XTRACE_ETCD=$(set +o | grep xtrace)
set +o xtrace

osds::etcd::stop() {
    kill "$(cat $ETCD_DIR/etcd.pid)" >/dev/null 2>&1 || :
    wait "$(cat $ETCD_DIR/etcd.pid)" >/dev/null 2>&1 || :
}

osds::etcd::clean_etcd_dir() {
      rm -rf "${ETCD_DIR-}"
}

osds::etcd::download() {
  (
    cd "${OPT_DIR}"
    url="https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz"
    download_file="etcd-v${ETCD_VERSION}-linux-amd64.tar.gz"
    osds::util::download_file "${url}" "${download_file}"
    tar xzf "${download_file}"
    cp etcd-v${ETCD_VERSION}-linux-amd64/etcd bin
    cp etcd-v${ETCD_VERSION}-linux-amd64/etcdctl bin
  )
}

osds::etcd::install() {
    # validate before running
    which etcd >/dev/null || {
    osds::etcd::download
    }

    # Start etcd
    mkdir -p $ETCD_DIR
    nohup etcd --advertise-client-urls http://${ETCD_HOST}:${ETCD_PORT} --listen-client-urls http://${ETCD_HOST}:${ETCD_PORT}\
    --listen-peer-urls http://${ETCD_HOST}:${ETCD_PEER_PORT} --data-dir ${ETCD_DATADIR} --debug 2> "${ETCD_LOGFILE}" >/dev/null &
    echo $! > $ETCD_DIR/etcd.pid

    osds::echo_summary "Waiting for etcd to come up."
    osds::util::wait_for_url "http://${ETCD_HOST}:${ETCD_PORT}/v2/machines" "etcd: " 0.25 80
    curl -fs -X PUT "http://${ETCD_HOST}:${ETCD_PORT}/v2/keys/_test"
}

osds::etcd::cleanup() {
    osds::etcd::stop
    osds::etcd::clean_etcd_dir
}

osds::etcd::uninstall(){
    : # do nothing
}

osds::etcd::uninstall_purge(){
    : # do nothing
}

# Restore xtrace
$_XTRACE_ETCD
