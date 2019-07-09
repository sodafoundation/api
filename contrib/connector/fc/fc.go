// Copyright 2018 The OpenSDS Authors.
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

package fc

import (
	"github.com/opensds/opensds/contrib/connector"
)

// FC struct
type FC struct{}

// init ...
func init() {
	connector.RegisterConnector(connector.FcDriver, &FC{})
}

// Attach ...
func (f *FC) Attach(conn map[string]interface{}) (string, error) {
	deviceInfo, err := connectVolume(conn)
	if err != nil {
		return "", err
	}
	return deviceInfo["path"], nil
}

// Detach ...
func (f *FC) Detach(conn map[string]interface{}) error {
	return disconnectVolume(conn)
}

// GetInitiatorInfo ...
func (f *FC) GetInitiatorInfo() (string, error) {
	return getInitiatorInfo()
}
