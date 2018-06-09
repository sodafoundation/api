// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

const (
	fcDriver = "fc"
)

type FC struct {
	self *fibreChannel
}

func init() {
	connector.RegisterConnector(fcDriver,
		&FC{
			self: &fibreChannel{
				helper: &linuxfc{},
			},
		})
}

func (f *FC) Attach(conn map[string]interface{}) (string, error) {
	deviceInfo, err := f.self.connectVolume(conn)
	if err != nil {
		return "", err
	}
	return deviceInfo["path"], nil
}

func (f *FC) Detach(conn map[string]interface{}) error {
	return f.self.disconnectVolume(conn)
}
