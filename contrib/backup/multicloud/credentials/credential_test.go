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

package credentials

import (
	"testing"
)

type mockProvider struct {
	credentials Value
	err         error
}

func (m *mockProvider) Retrieve() (Value, error) {
	m.credentials.ProviderName = "mockProvider"
	return m.credentials, m.err
}

func TestCredentialsGet(t *testing.T) {
	c := NewCredentials(&mockProvider{
		credentials: Value{
			AccessKeyID:     "access_key",
			SecretAccessKey: "secret_key",
		},
	})

	creds, err := c.Get()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e, a := "access_key", creds.AccessKeyID; e != a {
		t.Errorf("Expect access key ID to match, %v got %v", e, a)
	}
	if e, a := "secret_key", creds.SecretAccessKey; e != a {
		t.Errorf("Expect secret access key to match, %v got %v", e, a)
	}
}

func TestCredentialsGetWithProviderName(t *testing.T) {
	mock := &mockProvider{}

	c := NewCredentials(mock)

	credentials, err := c.Get()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e, a := credentials.ProviderName, "mockProvider"; e != a {
		t.Errorf("Expected provider name to match, %v got %v", e, a)
	}
}

func TestCredentialsGetWithError(t *testing.T) {
	mock := &mockProvider{}

	c := NewCredentials(mock)

	credentials, err := c.Get()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if e, a := credentials.ProviderName, "mockProvider"; e != a {
		t.Errorf("Expected provider name to match, %v got %v", e, a)
	}
}
