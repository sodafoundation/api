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
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/constants"
)

func NewHttpError(code int, msg string) error {
	return &HttpError{Code: code, Msg: msg}
}

type HttpError struct {
	Code int
	Msg  string
}

func (e *HttpError) Decode() {
	errSpec := model.ErrorSpec{}
	err := json.Unmarshal([]byte(e.Msg), &errSpec)
	if err == nil {
		e.Msg = errSpec.Message
	}
}

func (e *HttpError) Error() string {
	e.Decode()
	return fmt.Sprintf("Code: %v, Desc: %s, Msg: %v", e.Code, http.StatusText(e.Code), e.Msg)
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

func NewKeystoneReciver(auth *KeystoneAuthOptions) Receiver {
	k := &KeystoneReciver{auth: auth}
	k.GetToken()
	return k
}

type KeystoneReciver struct {
	auth *KeystoneAuthOptions
}

func (k *KeystoneReciver) GetToken() error {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: k.auth.IdentityEndpoint,
		Username:         k.auth.Username,
		UserID:           k.auth.UserID,
		Password:         k.auth.Password,
		DomainID:         k.auth.DomainID,
		DomainName:       k.auth.DomainName,
		TenantID:         k.auth.TenantID,
		TenantName:       k.auth.TenantName,
		AllowReauth:      k.auth.AllowReauth,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return fmt.Errorf("When get auth client: %v", err)
	}

	// Only support keystone v3
	identity, err := openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{})
	if err != nil {
		return fmt.Errorf("When get identity session: %v", err)
	}
	r := tokens.Create(identity, &opts)
	token, err := r.ExtractToken()
	if err != nil {
		return fmt.Errorf("When get extract token session: %v", err)
	}
	project, err := r.ExtractProject()
	if err != nil {
		return fmt.Errorf("When get extract project session: %v", err)
	}
	k.auth.TenantID = project.ID
	k.auth.TokenID = token.ID
	return nil
}

func (k *KeystoneReciver) Recv(url string, method string, body interface{}, output interface{}) error {
	desc := fmt.Sprintf("%s %s", method, url)
	return utils.Retry(2, desc, true, func(retryIdx int, lastErr error) error {
		if retryIdx > 0 {
			err, ok := lastErr.(*HttpError)
			if ok && err.Code == http.StatusUnauthorized {
				k.GetToken()
			} else {
				return lastErr
			}
		}

		headers := HeaderOption{}
		headers[constants.AuthTokenHeader] = k.auth.TokenID
		return request(url, method, headers, body, output)
	})
}

func checkHTTPResponseStatusCode(resp *http.Response) error {
	if 400 <= resp.StatusCode && resp.StatusCode <= 599 {
		return fmt.Errorf("response == %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return nil
}
