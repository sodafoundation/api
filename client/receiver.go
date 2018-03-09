// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego/httplib"
	"github.com/opensds/opensds/pkg/model"
)

func checkHTTPResponseStatusCode(resp *http.Response) error {
	if 400 <= resp.StatusCode && resp.StatusCode <= 599 {
		return fmt.Errorf("response == %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return nil
}

func NewHttpError(code int, msg string) error {
	return &HttpError{Code: code, Msg: msg}
}

type HttpError struct {
	Code int
	Msg  string
}

func (e *HttpError) Error() string {
	errSpec := model.ErrorSpec{}
	err := json.Unmarshal([]byte(e.Msg), &errSpec)
	if err != nil {
		return fmt.Sprintf("Code: %v, Desc: %s, Msg: %v", e.Code, http.StatusText(e.Code), e.Msg)
	}
	return fmt.Sprintf("Code: %v, Desc: %s, Msg: %v", e.Code, http.StatusText(e.Code), errSpec.Message)
}

// ParamOption
type HeaderOption map[string]string

// Receiver
type Receiver interface {
	Recv(url string, method string, input interface{}, output interface{}) error
}

// NewReceiver
func NewReceiver() Receiver {
	return &receiver{}
}

func request(url string, method string, headers HeaderOption, input interface{}, output interface{}) error {
	req := httplib.NewBeegoRequest(url, strings.ToUpper(method))
	// init body
	if input != nil {
		req.JSONBody(input)
	}
	//init header
	if headers != nil {
		for k, v := range headers {
			req.Header(k, v)
		}
	}
	// Get http response.
	resp, err := req.Response()
	if err != nil {
		return err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if 400 <= resp.StatusCode && resp.StatusCode <= 599 {
		return NewHttpError(resp.StatusCode, string(rbody))
	}

	// If the format of output is nil, skip unmarshaling the result.
	if output == nil {
		return nil
	}
	if err = json.Unmarshal(rbody, output); err != nil {
		return fmt.Errorf("failed to unmarshal result message: %v", err)
	}
	return nil
}

type receiver struct{}

func (*receiver) Recv(url string, method string, input interface{}, output interface{}) error {
	return request(url, method, nil, input, output)
}
