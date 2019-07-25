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

type NFS struct {
	*Client
}

func (c *NFS) getShareID(share interface{}) string {
	return share.(*NFSShareData).ID
}

func (c *NFS) createShare(shareName, fsID string) (interface{}, error) {
	sharePath := getSharePath(shareName)
	data := map[string]string{
		"DESCRIPTION": "",
		"FSID":        fsID,
		"SHAREPATH":   sharePath,
	}

	url := "/NFSHARE"

	resp, err := c.request(url, "POST", data)
	if err != nil {
		return nil, fmt.Errorf("create nfs share %s failed: %v", sharePath, err)
	}

	var nfsShare NFSShare
	if err := handleReponse(resp, &nfsShare); err != nil {
		return nil, fmt.Errorf("create nfs share %s failed: %v", sharePath, err)
	}

	return &nfsShare.Data, nil
}

func (c *NFS) getShare(shareName string) (interface{}, error) {
	url := fmt.Sprintf("/NFSHARE?filter=SHAREPATH::%s&range=[0-100]", getSharePath(shareName))
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var nfsShareList NFSShareList
	if err := handleReponse(resp, &nfsShareList); err != nil {
		return nil, err
	}

	if len(nfsShareList.Data) > 0 {
		return &nfsShareList.Data[0], nil
	}

	return nil, nil
}

func (c *NFS) deleteShare(shareID string) error {
	url := "/nfshare/" + shareID
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

func (c *NFS) getShareByID(shareID string) (interface{}, error) {
	url := "/NFSHARE/" + shareID
	resp, err := c.request(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	var nfsShare NFSShare
	if err := handleReponse(resp, &nfsShare); err != nil {
		return nil, err
	}

	if nfsShare.Data.ID == "" {
		return nil, nil
	}

	return &nfsShare.Data, nil
}

func (c *NFS) allowAccess(shareID, accessTo, accessLevel string) (interface{}, error) {
	url := "/NFS_SHARE_AUTH_CLIENT"
	data := map[string]string{
		"TYPE":       "16409",
		"NAME":       accessTo,
		"PARENTID":   shareID,
		"ACCESSVAL":  accessLevel,
		"SYNC":       "0",
		"ALLSQUASH":  "1",
		"ROOTSQUASH": "0",
	}

	resp, err := c.request(url, "Post", data)
	if err != nil {
		return nil, err
	}

	var nfsClient NFSShareClient
	if err := handleReponse(resp, &nfsClient); err != nil {
		return nil, err
	}

	return &nfsClient, nil
}

// getLocation
func (c *NFS) getLocation(sharePath, ipAddr string) string {
	path := strings.Replace(sharePath, "-", "_", -1)
	return fmt.Sprintf("%s:/%s", ipAddr, path)
}

func (c *NFS) getAccessLevel(accessLevel string) string {
	if accessLevel == AccessLevelRW {
		return AccessNFSRw
	}
	return AccessNFSRo
}
