#!/bin/bash

go get github.com/opensds/opensds/cmd/osdslet
go get github.com/opensds/opensds/cmd/osdsdock

(cd cmd/osdslet && ./build_docker.sh)
(cd cmd/osdsdock && ./build_docker.sh)