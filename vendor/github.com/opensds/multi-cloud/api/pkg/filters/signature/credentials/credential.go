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

package credentials

import (
	"sync"
)

// A Value is the OpenSDS multi-cloud credentials value.
type Value struct {
	// OpenSDS multi-cloud Access key ID
	AccessKeyID string

	// OpenSDS multi-cloud Secret Access Key
	SecretAccessKey string

	// Provider used to get credentials
	ProviderName string
}

// A Provider is the interface to provide credentials Value.
type Provider interface {
	// Retrieve returns credential value if it successfully retrieved the value,
	// otherwise returns Error.
	Retrieve() (Value, error)
}

// An ErrorProvider is a credentials provider that always returns an error
type ErrorProvider struct {
	// The error to be returned from Retrieve
	Err error

	// The provider name to set on the Retrieved returned Value
	ProviderName string
}

// Retrieve will always return the error that the ErrorProvider was created with.
func (p ErrorProvider) Retrieve() (Value, error) {
	return Value{ProviderName: p.ProviderName}, p.Err
}

// A Credentials provides concurrency safe retrieval of OpenSDS multi-cloud credentials Value.
type Credentials struct {
	credentials Value

	m sync.RWMutex

	provider Provider
}

// NewCredentials returns a pointer to a new Credentials with the provider set.
func NewCredentials(provider Provider) *Credentials {
	return &Credentials{
		provider: provider,
	}
}

// Get returns the credentials value, or error on failed retrieval
func (c *Credentials) Get() (Value, error) {

	c.m.RLock()
	credentials, err := c.provider.Retrieve()
	if err != nil {
		return Value{}, err
	}
	c.credentials = credentials
	c.m.RUnlock()

	return c.credentials, nil
}
