// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nfs

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/opensds/opensds/contrib/connector"
)

import (
	"fmt"
)

type NFSConnectorInfo struct {
	ExportLocations []string `mapstructure:"exportLocations"`
}

func connect(conn map[string]interface{}) (string, error) {
	exportLocation, err := parseNFSConnectInfo(conn)
	if err != nil {
		return "", err
	}

	ipAddr := strings.Split(exportLocation, ":")[0]
	sharePath := strings.Split(exportLocation, ":")[1]

	fmt.Printf("export locations: %v\n", exportLocation)

	showMountCommand := "showmount"
	_, err = exec.LookPath(showMountCommand)
	if err != nil {
		if err == exec.ErrNotFound {
			return "", fmt.Errorf("%q executable not found in $PATH", showMountCommand)
		}
		return "", err
	}

	cmd := fmt.Sprintf("%s -e %s", showMountCommand, ipAddr)
	res, err := connector.ExecCmd("/bin/bash", "-c", cmd)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(res, "\n") {
		if strings.Contains(line, sharePath) {
			str := strings.TrimSpace(line)
			strArray := strings.Split(str, " ")

			fileShareNameIdx := 0
			if strArray[fileShareNameIdx] == sharePath {
				return exportLocation, nil
			}
		}
	}

	return "", fmt.Errorf("cannot find fileshare path: %s", sharePath)
}

// ParseIscsiConnectInfo decode
func parseNFSConnectInfo(connectInfo map[string]interface{}) (string, error) {
	var con NFSConnectorInfo
	mapstructure.Decode(connectInfo, &con)

	fmt.Printf("connection data : %v\n", con)
	if len(con.ExportLocations) == 0 {
		return "", errors.New("nfs connection data is invalid")
	}

	for _, lo := range con.ExportLocations {
		strs := strings.Split(lo, ":")
		ipIdx := 0
		ip := strs[ipIdx]

		cmd := "ping -c 2 " + ip
		_, err := connector.ExecCmd("/bin/bash", "-c", cmd)
		if err != nil {
			fmt.Printf("ping error: %v\n", err)
		} else {
			return lo, nil
		}
	}

	return "", errors.New("no valid export location can be found")
}

func disconnect(conn map[string]interface{}) error {
	return errors.New("disconnect method of nfs is not implemented")
}

func getInitiatorInfo() (string, error) {
	return "", errors.New("get initiator information method of nfs is not implemented")
}
