// Copyright 2017 The OpenSDS Authors.
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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/astaxie/beego/httplib"
	"github.com/opensds/opensds/pkg/model"
)

type reqFunc func(string, string, interface{}) *httplib.BeegoHTTPRequest

type Receiver interface {
	Recv(reqFunc, string, string, interface{}, interface{}) error
}

func NewReceiver() Receiver {
	return &receiver{}
}

type receiver struct{}

func (*receiver) Recv(
	f reqFunc,
	url string,
	method string,
	input interface{},
	output interface{},
) error {
	req := f(url, method, input)

	// Get http response.
	resp, err := req.Response()
	if err != nil {
		return err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = checkHTTPResponseStatusCode(resp); err != nil {
		var errorMsg model.ErrorSpec
		if err = json.Unmarshal(rbody, &errorMsg); err != nil {
			return fmt.Errorf("failed to unmarshal error message: %v", err)
		}
		return fmt.Errorf("failed to exec this operation, code: %v, message: %v",
			errorMsg.Code, errorMsg.Message)
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

type ParamOption map[string]string

func request(
	url string,
	method string,
	input interface{},
) *httplib.BeegoHTTPRequest {
	req := httplib.NewBeegoRequest(url, strings.ToUpper(method))

	// If input is empty, consider as user doesn't specify any parameter or
	// request body.
	if input == nil {
		return req
	}
	// If input type is a map, it means user tends to send a GET request
	// and specifiy some parameters.
	p, ok := input.(ParamOption)
	if ok {
		for k := range p {
			req.Param(k, p[k])
		}
		return req
	}
	// The other condition will be handled in httplib package.
	req.JSONBody(input)
	return req
}

// CheckHTTPResponseStatusCode compares http response header StatusCode against expected
// statuses. Primary function is to ensure StatusCode is in the 20x (return nil).
// Ok: 200. Created: 201. Accepted: 202. No Content: 204. Partial Content: 206.
// Otherwise return error message.
func checkHTTPResponseStatusCode(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 202, 204, 206:
		return nil
	case 400:
		return errors.New("Error: response == 400 bad request")
	case 401:
		return errors.New("Error: response == 401 unauthorised")
	case 403:
		return errors.New("Error: response == 403 forbidden")
	case 404:
		return errors.New("Error: response == 404 not found")
	case 405:
		return errors.New("Error: response == 405 method not allowed")
	case 409:
		return errors.New("Error: response == 409 conflict")
	case 413:
		return errors.New("Error: response == 413 over limit")
	case 415:
		return errors.New("Error: response == 415 bad media type")
	case 422:
		return errors.New("Error: response == 422 unprocessable")
	case 429:
		return errors.New("Error: response == 429 too many request")
	case 500:
		return errors.New("Error: response == 500 instance fault / server err")
	case 501:
		return errors.New("Error: response == 501 not implemented")
	case 503:
		return errors.New("Error: response == 503 service unavailable")
	}
	return errors.New("Error: unexpected response status code")
}
