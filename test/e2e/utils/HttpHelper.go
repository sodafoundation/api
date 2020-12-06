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

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func readResponseBody(resp *http.Response) {
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("reading response body failed: %v", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf(bodyString)
	}
}


func ConnectToHTTP(operation, testMgrEndPoint string, payload map[string]interface{}) (*http.Response, error) {
	var httpMethod string

	switch operation {
	case "PUT":
		httpMethod = http.MethodPut

	case "POST":
		httpMethod = http.MethodPost
	case "DELETE":
		httpMethod = http.MethodDelete

	case "GET":
		httpMethod = http.MethodGet
	default:

	}

	respbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("payload marshal failed: %v", err)
	}

	req, err := http.NewRequest(httpMethod, testMgrEndPoint, bytes.NewBuffer(respbytes))
	if err != nil {
		fmt.Printf("error while getting http request: %v", err)

	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("hTTP request is failed :%v", err)

	}
	if resp != nil {
		fmt.Printf("resp is ... :%v", resp)
		io.Copy(os.Stdout, resp.Body)
		defer resp.Body.Close()
		readResponseBody(resp)
	}

	return resp, err
}