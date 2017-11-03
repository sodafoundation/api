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
This module implements the common data structure.

*/

package model

import (
	"errors"
)

type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (r *Response) GetStatus() string {
	return r.Status
}

func (r *Response) GetError() string {
	return r.Error
}

func (r *Response) ToError() error {
	if r.Error != "" {
		return errors.New(r.Error)
	}

	return nil
}

func (r *Response) GetMessage() string {
	return r.Message
}

func (r *Response) SetStatus(stat string) {
	r.Status = stat
}

func (r *Response) SetError(err string) {
	r.Error = err
}

func (r *Response) SetMessage(mesg string) {
	r.Message = mesg
}
