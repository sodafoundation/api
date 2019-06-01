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

package oceanstor

import (
	"fmt"
	"strings"
)

type CIFS struct {
	*Client
}

func (c *CIFS) getShareID(share interface{}) string {
	return share.(*CIFSShareData).ID
}

func (c *CIFS) getShareByID(shareID string) (interface{}, error) {
	url := "/CIFSHARE/" + shareID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var cifsShare CIFSShare
	if err := handleReponse(resp, &cifsShare); err != nil {
		return nil, err
	}

	if cifsShare.Data.ID == "" {
		return nil, nil
	}

	return &cifsShare.Data, nil
}

func (c *CIFS) createShare(shareName, fsId string) (interface{}, error) {
	sharePath := getSharePath(shareName)
	data := map[string]string{
		"SHAREPATH":    sharePath,
		"DESCRIPTION":  "",
		"ABEENABLE":    "false",
		"ENABLENOTIFY": "true",
		"ENABLEOPLOCK": "true",
		"NAME":         strings.Replace(shareName, "-", "_", -1),
		"FSID":         fsId,
		"TENANCYID":    "0",
	}

	url := "/CIFSHARE"

	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, err
	}

	var cifsShare CIFSShare

	if err := handleReponse(resp, &cifsShare); err != nil {
		return nil, err
	}

	return &cifsShare.Data, nil
}

func (c *CIFS) getShare(shareName string) (interface{}, error) {
	url := fmt.Sprintf("/CIFSHARE?filter=NAME:%s&range=[0-100]", strings.Replace(shareName, "-", "_", -1))
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var cifsShareList CIFSShareList
	if err := handleReponse(resp, &cifsShareList); err != nil {
		return nil, err
	}

	if len(cifsShareList.Data) > 0 {
		return &cifsShareList.Data[0], nil
	}

	return nil, nil
}

func (c *CIFS) listShares() ([]CIFSShareData, error) {
	url := "/CIFSHARE"
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var cifsShareList CIFSShareList
	if err := handleReponse(resp, &cifsShareList); err != nil {
		return nil, err
	}

	return cifsShareList.Data, nil
}

func (c *CIFS) deleteShare(shareID string) error {
	url := "/cifshare/" + shareID
	resp, err := c.request(url, "DELETE", nil)
	if err != nil {
		return err
	}

	var errDelete DeleteError
	if err := handleReponse(resp, &errDelete); err != nil {
		return err
	}

	return nil
}

func (c *CIFS) allowAccess(shareID, accessTo, accessLevel string) (interface{}, error) {

	domainType := map[string]string{"local": "2", "ad": "0"}

	sendRest := func(accessTo, domain string) (*CIFSShareClientData, error) {
		url := "/CIFS_SHARE_AUTH_CLIENT"
		data := map[string]string{
			"NAME":       accessTo,
			"PARENTID":   shareID,
			"PERMISSION": accessLevel,
			"DOMAINTYPE": domain,
		}

		resp, err := c.request(url, "POST", data)
		if err != nil {
			return nil, err
		}

		var cifsClient CIFSShareClient
		if err := handleReponse(resp, &cifsClient); err != nil {
			return nil, err
		}

		return &cifsClient.Data, nil
	}

	var data *CIFSShareClientData
	var errRest error

	if !strings.Contains(accessTo, "\\\\") {
		// First, try to add user access
		if data, errRest = sendRest(accessTo, domainType["local"]); errRest != nil {
			// Second, if add user access failed, try to add group access.
			if data, errRest = sendRest("@"+accessTo, domainType["local"]); errRest != nil {
				return nil, errRest
			}
		}

	} else {
		// If add domain user access failed, try to add domain group access.
		if data, errRest = sendRest(accessTo, domainType["ad"]); errRest != nil {
			if data, errRest = sendRest("@"+accessTo, domainType["ad"]); errRest != nil {
				return nil, errRest
			}
		}
	}

	return data, nil
}

// getLocation
func (c *CIFS) getLocation(sharePath, ipAddr string) string {
	path := strings.Replace(sharePath, "-", "_", -1)
	return fmt.Sprintf("\\\\%s\\%s", ipAddr, path)
}

func (c *CIFS) getAccessLevel(accessLevel string) string {
	if accessLevel == AccessLevelRW {
		return AccessCIFSFullControl
	}
	return AccessCIFSRo
}
