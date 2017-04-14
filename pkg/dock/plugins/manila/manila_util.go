// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

/*
This module implements manila plugin for OpenSDS. Manila plugin will pass these
operation requests about share to OpenStack go-client module.

*/

package manila

import (
	"errors"
	"log"
	// "time"
)

func AttachShareToHost(plugin *ManilaPlugin, shrID string) (string, error) {
	shareService, err := plugin.getShareService()
	if err != nil {
		log.Println("Cannot access share service:", err)
		return "", err
	}
	shr, err := shareService.Show(shrID)
	if err != nil {
		log.Println("Cannot get share:", err)
		return "", err
	}

	if shr.Export_location == "" {
		err = errors.New("Share not exported!")
		return "", err
	}

	return shr.Export_location, nil
}

func DetachShareFromHost(plugin *ManilaPlugin, device string) (string, error) {
	return "Detach share success!", nil
}
