#!/bin/bash

# Copyright 2019 The OpenSDS Authors.
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

# Save trace setting
_XTRACE_CEPH=$(set +o | grep xtrace)
set +o xtrace

COMPONENT=("opensds" "nbp")
OPENSDS_CERT_DIR="/opt/opensds-security"
ROOT_CERT_DIR=${ROOT_CERT_DIR:-"${OPENSDS_CERT_DIR}"/ca}

osds::certificate::install(){
    osds::certificate::cleanup
    osds::certificate::check_openssl_installed
    osds::certificate::prepare
    osds::certificate::create_ca_cert
    osds::certificate::create_component_cert
}

osds::certificate::cleanup() {
    osds::certificate::uninstall
}

osds::certificate::uninstall(){
    if [ -d "${OPENSDS_CERT_DIR}" ];then
        rm -rf "${OPENSDS_CERT_DIR}"
    fi
}

osds::certificate::uninstall_purge(){
    osds::certificate::uninstall
}

osds::certificate::check_openssl_installed(){
	openssl version >& /dev/null
    if [ $? -ne 0 ];then
        echo "Failed to run openssl. Please ensure openssl is installed."
        exit 1
    fi
}

osds::certificate::prepare(){
    # Prepare to generate certs
    mkdir -p "${ROOT_CERT_DIR}"
    mkdir -p "${ROOT_CERT_DIR}"/demoCA/
    mkdir -p "${ROOT_CERT_DIR}"/demoCA/newcerts
    touch "${ROOT_CERT_DIR}"/demoCA/index.txt
    echo "01" > "${ROOT_CERT_DIR}"/demoCA/serial
    echo "unique_subject = no" > "${ROOT_CERT_DIR}"/demoCA/index.txt.attr
}

osds::certificate::create_ca_cert(){
    # Create ca cert
    cd "${ROOT_CERT_DIR}"
    openssl genrsa -passout pass:xxxxx -out "${ROOT_CERT_DIR}"/ca-key.pem -aes256 2048
    openssl req -new -x509 -sha256 -key "${ROOT_CERT_DIR}"/ca-key.pem -out "${ROOT_CERT_DIR}"/ca-cert.pem -days 365 -subj "/CN=CA" -passin pass:xxxxx
}

osds::certificate::create_component_cert(){
	# Create component cert
    for com in ${COMPONENT[*]};do
	    openssl genrsa -aes256 -passout pass:xxxxx -out "${ROOT_CERT_DIR}"/"${com}"-key.pem 2048
	    openssl req -new -sha256 -key "${ROOT_CERT_DIR}"/"${com}"-key.pem -out "${ROOT_CERT_DIR}"/"${com}"-csr.pem -days 365 -subj "/CN=${com}" -passin pass:xxxxx
	    openssl ca -batch -in "${ROOT_CERT_DIR}"/"${com}"-csr.pem -cert "${ROOT_CERT_DIR}"/ca-cert.pem -keyfile "${ROOT_CERT_DIR}"/ca-key.pem -out "${ROOT_CERT_DIR}"/"${com}"-cert.pem -md sha256 -days 365 -passin pass:xxxxx
	
	    # Cancel the password for the private key
        openssl rsa -in "${ROOT_CERT_DIR}"/"${com}"-key.pem -out "${ROOT_CERT_DIR}"/"${com}"-key.pem -passin pass:xxxxx
	   
	    mkdir -p "${OPENSDS_CERT_DIR}"/"${com}"
	    mv "${ROOT_CERT_DIR}"/"${com}"-key.pem "${OPENSDS_CERT_DIR}"/"${com}"/
	    mv "${ROOT_CERT_DIR}"/"${com}"-cert.pem "${OPENSDS_CERT_DIR}"/"${com}"/
	    rm -rf "${ROOT_CERT_DIR}"/"${com}"-csr.pem
    done
	
    rm -rf "${ROOT_CERT_DIR}"/demoCA
}

# Restore xtrace
$_XTRACE_CEPH
