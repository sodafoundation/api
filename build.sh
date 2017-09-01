#!/bin/bash

go get github.com/opensds/opensds/cmd/osdslet
go get github.com/opensds/opensds/cmd/osdsdock

go install github.com/opensds/opensds/cmd/osdslet
go install github.com/opensds/opensds/cmd/osdsdock
