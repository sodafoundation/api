// Copyright 2019 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package keystonecredentials_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/opensds/multi-cloud/api/pkg/filters/signature/credentials/keystonecredentials"
	th "github.com/opensds/multi-cloud/testhelper"
	"github.com/opensds/multi-cloud/testhelper/gophercloudclient"
)

const ListOutput = `
{
    "credentials": [
        {
            "user_id": "6f556708d04b4ea6bc72d7df2296b71a",
            "links": {
                "self": "http://identity/v3/credentials/2441494e52ab6d594a34d74586075cb299489bdd1e9389e3ab06467a4f460609"
            },
            "blob": "{\"access\":\"7da79ff0aa364e1396f067e352b9b79a\",\"secret\":\"7a18d68ba8834b799d396f3ff6f1e98c\"}",
            "project_id": "1a1d14690f3c4ec5bf5f321c5fde3c16",
            "type": "ec2",
            "id": "2441494e52ab6d594a34d74586075cb299489bdd1e9389e3ab06467a4f460609"
        },
		{
            "user_id": "bb5476fd12884539b41d5a88f838d773",
            "links": {
                "self": "http://identity/v3/credentials/3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510"
            },
            "blob": "{\"access\":\"access_key\",\"secret\":\"secret_key\"}",
            "project_id": "731fc6f265cd486d900f16e84c5cb594",
            "type": "ec2",
            "id": "3d3367228f9c7665266604462ec60029bcd83ad89614021a80b2eb879c572510"
        }
	],
    "links": {
        "self": "http://identity/v3/credentials",
        "previous": null,
        "next": null
    }
}
`

func HandleListCredentialsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/credentials", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", gophercloudclient.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

func TestRetrieveKeystoneCredentials(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListCredentialsSuccessfully(t)

	c := gophercloudclient.ServiceClient()
	p := &keystonecredentials.KeystoneProvider{c, "access_key"}
	credentials, err := p.Retrieve()
	t.Log(credentials)

	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	if e, a := "access_key", credentials.AccessKeyID; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
	if e, a := "secret_key", credentials.SecretAccessKey; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

	if e, a := "KeystoneProvider", credentials.ProviderName; e != a {
		t.Errorf("expect %v, got %v", e, a)
	}

}

func TestRetrieveKeystoneCredentialsWithError(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListCredentialsSuccessfully(t)

	c := gophercloudclient.ServiceClient()
	p := &keystonecredentials.KeystoneProvider{c, "access"}
	credentials, err := p.Retrieve()
	t.Log(credentials)

	if e, a := "credential is missing", err.Error(); e != a {
		t.Errorf("Expected credential error, %v got %v", e, a)
	}
}
