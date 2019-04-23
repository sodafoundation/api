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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

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

func customVerify(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	roots := x509.NewCertPool()
	caStr, err := ioutil.ReadFile(constants.OpensdsCaCertFile)
	if err != nil {
		msg := fmt.Sprintf("Read certfile failed,error:%v ", err)
		log.Println(msg)
		return err
	}

	roots.AppendCertsFromPEM(caStr)

	for _, rawCert := range rawCerts {
		cert, _ := x509.ParseCertificate(rawCert)
		opts := x509.VerifyOptions{
			Roots: roots,
		}
		_, err := cert.Verify(opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func request(urlStr string, method string, headers HeaderOption, input interface{}, output interface{}) error {
	req := httplib.NewBeegoRequest(urlStr, strings.ToUpper(method))

	u, _ := url.Parse(urlStr)
	if u.Scheme == "https" && cacert != "" {
		log.Println("Https mode.")
		req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true, VerifyPeerCertificate: customVerify})
	}

	// Set the request timeout a little bit longer upload snapshot to cloud temporarily.
	req.SetTimeout(time.Minute*6, time.Minute*6)
	// init body
	if input != nil {
		body, err := json.MarshalIndent(input, "", "  ")
		if err != nil {
			return err
		}
		req.Body(body)
	}

	//init header
	if headers != nil {
		for k, v := range headers {
			req.Header(k, v)
		}
	}
	// Get http response.
	resp, err := req.Response()
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("\nStatusCode: %s\nResponse Body:\n%s\n", resp.Status, string(rbody))
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

func NewKeystoneReceiver(auth *KeystoneAuthOptions) (Receiver, error) {
	k := &KeystoneReceiver{Auth: auth}
	if err := k.GetToken(); err != nil {
		log.Printf("get token failed, %v", err)
		return nil, err
	}
	return k, nil
}

type KeystoneReceiver struct {
	Auth *KeystoneAuthOptions
}

func (k *KeystoneReceiver) GetToken() error {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: k.Auth.IdentityEndpoint,
		Username:         k.Auth.Username,
		UserID:           k.Auth.UserID,
		Password:         k.Auth.Password,
		DomainID:         k.Auth.DomainID,
		DomainName:       k.Auth.DomainName,
		TenantID:         k.Auth.TenantID,
		TenantName:       k.Auth.TenantName,
		AllowReauth:      k.Auth.AllowReauth,
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
	k.Auth.TenantID = project.ID
	k.Auth.TokenID = token.ID
	return nil
}

func (k *KeystoneReceiver) Recv(url string, method string, body interface{}, output interface{}) error {
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
		headers[constants.AuthTokenHeader] = k.Auth.TokenID
		return request(url, method, headers, body, output)
	})
}

func checkHTTPResponseStatusCode(resp *http.Response) error {
	if 400 <= resp.StatusCode && resp.StatusCode <= 599 {
		return fmt.Errorf("response == %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return nil
}
