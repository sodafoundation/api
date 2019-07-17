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

package signer_test

import (
	"net/http"
	"testing"

	"github.com/opensds/multi-cloud/api/pkg/filters/signature/credentials"
	"github.com/opensds/multi-cloud/api/pkg/filters/signature/signer"
)

const authorizationStr = "OPENSDS-HMAC-SHA256 Credential=access_key/20190301/us-east-1/s3/sign_request,SignedHeaders=authorization;host;x-auth-date,Signature=472f0a1b7815974847620da53fcdd2fdd53203b5d8d08e7ce81943b260560e26"

type mockProvider struct {
	credentials credentials.Value
	err         error
}

func (m *mockProvider) Retrieve() (credentials.Value, error) {
	m.credentials.ProviderName = "mockProvider"
	return m.credentials, m.err
}

func buildRequest(serviceName, region string) *http.Request {
	endpoint := "https://" + serviceName + "/" + region
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Add("X-Auth-Date", "20190301T220855Z")
	req.Header.Add("Authorization", authorizationStr)
	return req
}

func buildSigner() signer.Signer {
	c := credentials.NewCredentials(&mockProvider{
		credentials: credentials.Value{
			AccessKeyID:     "access_key",
			SecretAccessKey: "secret_key",
		},
	})

	return signer.Signer{
		Credentials: c,
	}
}

func TestSignRequestValidation(t *testing.T) {
	req := buildRequest("s3", "us-east-1")
	signer := buildSigner()
	credentialStr := "access_key/20190301/us-east-1/s3/sign_request"
	calculatedSignature, err := signer.Sign(req, "", "s3", "us-east-1", "20190301T220855Z", "20190301", credentialStr)

	if err != nil {
		return
	}

	expectedDate := "20190301T220855Z"
	expectedSignature := "219270ad1f7a4430b7a6362ab525d94134f5a2b7c3f7e54a6b67629f5262dda4"

	if e, a := expectedSignature, calculatedSignature; e != a {
		t.Errorf("expect \n %v \n actual \n %v \n", e, a)
	}

	q := req.Header

	if e, a := authorizationStr, q.Get("Authorization"); e != a {
		t.Errorf("expect \n %v \n actual \n %v \n", e, a)
	}
	if e, a := expectedDate, q.Get("X-Auth-Date"); e != a {
		t.Errorf("expect \n %v \nactual \n %v \n", e, a)
	}

}
