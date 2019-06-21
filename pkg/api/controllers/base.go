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

package controllers

import (
	"net/http"
	"net/url"
	"reflect"

	"github.com/astaxie/beego"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
)

const (
	StatusOK       = http.StatusOK
	StatusAccepted = http.StatusAccepted
)

type BasePortal struct {
	beego.Controller
}

func (b *BasePortal) GetParameters() (map[string][]string, error) {
	u, err := url.Parse(b.Ctx.Request.URL.String())
	if err != nil {
		return nil, err
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Filter some items in spec that no need to transfer to users.
func (b *BasePortal) outputFilter(resp interface{}, whiteList []string) interface{} {
	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Slice {
		var s []map[string]interface{}
		for i := 0; i < v.Len(); i++ {
			m := b.doFilter(v.Index(i).Interface(), whiteList)
			s = append(s, m)
		}
		return s
	} else {
		return b.doFilter(resp, whiteList)
	}
}

func (b *BasePortal) doFilter(resp interface{}, whiteList []string) map[string]interface{} {
	v := reflect.ValueOf(resp).Elem()
	m := map[string]interface{}{}
	for _, name := range whiteList {
		field := v.FieldByName(name)
		if field.IsValid() {
			m[name] = field.Interface()
		}
	}
	return m
}

func (b *BasePortal) ErrorHandle(errType int, errMsg string) {
	var errBody []byte

	switch errType {
	case model.ErrorBadRequest:
		errBody = model.ErrorBadRequestStatus(errMsg)
	case model.ErrorUnauthorized:
		errBody = model.ErrorUnauthorizedStatus(errMsg)
	case model.ErrorForbidden:
		errBody = model.ErrorForbiddenStatus(errMsg)
	case model.ErrorNotFound:
		errBody = model.ErrorNotFoundStatus(errMsg)
	case model.ErrorInternalServer:
		errBody = model.ErrorInternalServerStatus(errMsg)
	default:
		errBody = model.ErrorNotImplementedStatus(errMsg)
	}

	b.Ctx.Output.SetStatus(errType)
	b.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	b.Ctx.Output.Body(errBody)
	log.Error(errMsg)
}

func (b *BasePortal) SuccessHandle(status int, body []byte) {
	b.Ctx.Output.SetStatus(status)
	b.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	if body != nil {
		b.Ctx.Output.Body(body)
	}
}
