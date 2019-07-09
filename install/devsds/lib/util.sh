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

osds::util::sortable_date() {
  date "+%Y%m%d-%H%M%S"
}

osds::util::wait_for_url() {
  local url=$1
  local prefix=${2:-}
  local wait=${3:-1}
  local times=${4:-30}

  which curl >/dev/null || {
    osds::echo_summary "curl must be installed"
    exit 1
  }

  local i
  for i in $(seq 1 $times); do
    local out
    if out=$(curl --max-time 1 -gkfs $url 2>/dev/null); then
      osds::echo_summary "On try ${i}, ${prefix}: ${out}"
      return 0
    fi
    sleep ${wait}
  done
  osds::echo_summary "Timed out waiting for ${prefix} to answer at ${url}; tried ${times} waiting ${wait} between each"
  return 1
}

# returns a random port
osds::util::get_random_port() {
  awk -v min=1024 -v max=65535 'BEGIN{srand(); print int(min+rand()*(max-min+1))}'
}

# use netcat to check if the host($1):port($2) is free (return 0 means free, 1 means used)
osds::util::test_host_port_free() {
  local host=$1
  local port=$2
  local success=0
  local fail=1

  which nc >/dev/null || {
    osds::echo_summary "netcat isn't installed, can't verify if ${host}:${port} is free, skipping the check..."
    return ${success}
  }

  if [ ! $(nc -vz "${host}" "${port}") ]; then
    echo "${host}:${port} is free, proceeding..."
    return ${success}
  else
    echo "${host}:${port} is already used"
    return ${fail}
  fi
}

osds::util::download_file() {
  local -r url=$1
  local -r destination_file=$2

  rm  ${destination_file} 2&> /dev/null || true

  for i in $(seq 5)
  do
    if ! curl -fsSL --retry 3 --keepalive-time 2 ${url} -o ${destination_file}; then
      echo "Downloading ${url} failed. $((5-i)) retries left."
      sleep 1
    else
      echo "Downloading ${url} succeed"
      return 0
    fi
  done
  return 1
}

# Prints line number and "message" in warning format
# warn $LINENO "message"
osds::util::warn() {
    local exitcode=$?
    local xtrace
    xtrace=$(set +o | grep xtrace)
    set +o xtrace
    local msg="[WARNING] ${BASH_SOURCE[2]}:$1 $2"
    echo $msg
    $xtrace
    return $exitcode
}

# Prints line number and "message" in error format
# err $LINENO "message"
osds::util::err() {
    local exitcode=$?
    local xtrace
    xtrace=$(set +o | grep xtrace)
    set +o xtrace
    local msg="[ERROR] ${BASH_SOURCE[2]}:$1 $2"
    echo $msg 1>&2;
    $xtrace
    return $exitcode
}

# Prints line number and "message" then exits
# die $LINENO "message"
osds::util::die() {
    local exitcode=$?
    set +o xtrace
    local line=$1; shift
    if [ $exitcode == 0 ]; then
        exitcode=1
    fi
    osds::util::err $line "$*"
    # Give buffers a second to flush
    sleep 1
    exit $exitcode
}

osds::util::get_default_host_ip() {
    local host_ip=$1
    local af=$2
    # Search for an IP unless an explicit is set by ``HOST_IP`` environment variable
    if [ -z "$host_ip" ]; then
        host_ip=""
        # Find the interface used for the default route
        host_ip_iface=${host_ip_iface:-$(ip -f $af route | awk '/default/ {print $5}' | head -1)}
        local host_ips
        host_ips=$(LC_ALL=C ip -f $af addr show ${host_ip_iface} | sed /temporary/d |awk /$af'/ {split($2,parts,"/");  print parts[1]}')
        local ip
        for ip in $host_ips; do
            host_ip=$ip
            break;
        done
    fi
    echo $host_ip
}

osds::util::is_service_enabled() {
    local enabled=1
    local services=$@
    local service
    for service in ${services}; do
        [[ ,${OPENSDS_ENABLED_SERVICES}, =~ ,${service}, ]] && enabled=0
    done
    return $enabled
}

osds::util::serice_operation() {
    local action=$1
    for service in $(echo $SUPPORT_SERVICES|tr ',' ' '); do
            if osds::util::is_service_enabled $service; then
                source $TOP_DIR/lib/$service.sh
                osds::$service::$action
            fi
    done
}
