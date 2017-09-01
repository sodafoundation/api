#!/bin/bash

go build

docker build . -t opensds/opensds-controller:v1alpha
