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
	"github.com/opensds/opensds/contrib/connector"
)

type NFS struct{}

func init() {
	connector.RegisterConnector(connector.NFSDriver, &NFS{})
}

func (n *NFS) Attach(conn map[string]interface{}) (string, error) {
	return connect(conn)
}

func (n *NFS) Detach(conn map[string]interface{}) error {
	return disconnect(conn)
}

// GetInitiatorInfo implementation
func (n *NFS) GetInitiatorInfo() (string, error) {
	return getInitiatorInfo()
}
