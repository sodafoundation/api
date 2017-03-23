#!/bin/bash
#
# Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
#
#    Licensed under the Apache License, Version 2.0 (the "License"); you may
#    not use this file except in compliance with the License. You may obtain
#    a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing, software
#    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#    License for the specific language governing permissions and limitations
#    under the License.
#
# Enables all the examples to execute as a form of acceptance testing.

# Get the directory the examples are in and change into it.
DIR="$(cd $(dirname "$0") && pwd)"
echo "Executing the examples in: $DIR"
cd $DIR

# Run all the tests.
for T in $(ls -1 [0-9][0-9]*.go); do
	if ! [ -x $T ]; then
		CMD="go run $T setup.go"
		echo "$CMD ..."
		if ! $CMD ; then
			echo "Error executing example $T."
			exit 1
		fi
	fi
done