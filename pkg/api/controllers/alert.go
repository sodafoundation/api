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

/*
This module implements a entry into the OpenSDS northbound service to send Alerts to Alert manager

*/

package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/model"
)

func NewAlertPortal() *AlertPortal {
	return &AlertPortal{
		CtrClient: client.NewClient(),
	}
}

type AlertPortal struct {
	BasePortal

	CtrClient client.Client
}

func (v *AlertPortal) CreateAlert() {

	var postableAlert = model.PostableAlertSpec{}

	// Unmarshal the request body
	if err := json.NewDecoder(v.Ctx.Request.Body).Decode(&postableAlert); err != nil {
		errMsg := fmt.Sprintf("parse alert request body failed: %s", err.Error())
		v.ErrorHandle(model.ErrorBadRequest, errMsg)
		return
	}

	alertArr := make([]*model.PostableAlertSpec, 0)
	alertArr = append(alertArr, &postableAlert)

	b, e := json.Marshal(alertArr)
	if e != nil {
		log.Error(e)
		v.ErrorHandle(model.ErrorInternalServer, e.Error())
		return
	}

	body := strings.NewReader(string(b[:]))

	// Alert manager will be co-located on the server, default port is 9093 for the POST API endpoint
	// Raised issue https://github.com/opensds/opensds/issues/691 to make this configurable
	req, err := http.NewRequest("POST", "http://localhost:9093/api/v1/alerts", body)
	if err != nil {
		// handle err
		v.ErrorHandle(model.ErrorInternalServer, e.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		v.ErrorHandle(model.ErrorInternalServer, e.Error())
		return
	}
	defer resp.Body.Close()

	// Marshal the result.
	resBody, _ := json.Marshal(resp)
	v.SuccessHandle(StatusAccepted, resBody)

	return
}
