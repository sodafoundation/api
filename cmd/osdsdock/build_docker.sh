#!/bin/bash

go build

docker build . -t opensds/opensds-dock:v1alpha
