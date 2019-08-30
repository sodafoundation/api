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

func ConnectToHttpPost(operation, testMgrEndPoint string, payload []byte) {
	respbytes, err := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, testMgrEndPoint, bytes.NewBuffer(respbytes))
	if err != nil {
		// handle error
		// common.Failf("Frame HTTP request failed: %v", err)

		fmt.Errorf("Error while getting http request", err)
		// return false
	}

	// client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	// t := time.Now()
	resp, err := http.Post(testMgrEndPoint, "application/json", req.Body)

	if err != nil {
		// handle error
		fmt.Printf("HTTP request is failed :%v", err)
		// return false
	}
	if resp != nil {
		// handle error
		fmt.Printf("Resp is there :%v", resp.Body)

		defer resp.Body.Close()

		readResponseBody(resp)
	}

}
func ConnectToHttpGet(operation, testMgrEndPoint string) {

	req, err := http.NewRequest(http.MethodGet, testMgrEndPoint, nil)
	if err != nil {
		// handle error
		// common.Failf("Frame HTTP request failed: %v", err)

		fmt.Errorf("Error while getting http request", err)
		// return false
	}

	// client := &http.Client{}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// t := time.Now()
	resp, err := http.Get(testMgrEndPoint)

	if err != nil {
		// handle error
		fmt.Printf("HTTP request is failed :%v", err)
		// return false
	}
	if resp != nil {
		// handle error
		fmt.Printf("Resp is there :%v", resp.Body)

		defer resp.Body.Close()

		readResponseBody(resp)
	}
}
func readResponseBody(resp *http.Response) {
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("", err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf(bodyString)
	}
}

//function to handle device addition and deletion.
func ConnectToHTTP(operation, testMgrEndPoint string, payload map[string]interface{}) (*http.Response, error) {
	var httpMethod string
	// var payload dttype.MembershipUpdate
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
		fmt.Errorf("Payload marshal failed", err)
	}

	req, err := http.NewRequest(httpMethod, testMgrEndPoint, bytes.NewBuffer(respbytes))
	if err != nil {
		fmt.Errorf("Error while getting http request", err)

	}

	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	// t := time.Now()
	resp, err := client.Do(req)

	if err != nil {
		// handle error
		fmt.Printf("HTTP request is failed :%v", err)

	}
	if resp != nil {
		// handle error
		fmt.Printf("Resp is ... :%v", resp)
		fmt.Println("response Status:", resp.Status)
		// Print the body to the stdout
		io.Copy(os.Stdout, resp.Body)

		defer resp.Body.Close()

		readResponseBody(resp)
	}
	// common.InfoV6("%s %s %v in %v", req.Method, req.URL, resp.Status, time.Now().Sub(t))
	return resp, err
}
