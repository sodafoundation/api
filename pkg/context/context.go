// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

// This is self defined context which is stored in context.Input.data.
// It is used to transport data in the pipe line.

package context

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/astaxie/beego/context"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/model"
)

func NewAdminContext() *Context {
	return &Context{
		IsAdmin: true,
	}
}

func NewInternalTenantContext(tenantId, userId string) *Context {
	return &Context{
		TenantId: tenantId,
		UserId:   userId,
		IsAdmin:  true,
	}
}

func NewContextFromJson(s string) *Context {
	ctx := &Context{}
	err := json.Unmarshal([]byte(s), ctx)
	if err != nil {
		log.Errorf("Unmarshal json to context failed, reason: %v", err)
	}
	return ctx
}

func GetContext(httpCtx *context.Context) *Context {
	ctx, _ := httpCtx.Input.GetData("context").(*Context)
	if ctx == nil {
		ctx = &Context{}
	}
	return ctx
}

func UpdateContext(httpCtx *context.Context, param map[string]interface{}) (*Context, error) {

	ctx := GetContext(httpCtx)
	if param == nil || len(param) == 0 {
		log.Warning("Context parameter is empty, nothing to be updated")
		return ctx, nil
	}
	ctxV := reflect.ValueOf(ctx).Elem()
	for key, val := range param {
		field := ctxV.FieldByName(key)
		pv := reflect.ValueOf(val)
		if field.Kind() == pv.Kind() && field.CanSet() {
			field.Set(pv)
		} else {
			log.Errorf("Invalid parameter %s : %v", key, val)
		}
	}

	httpCtx.Input.SetData("context", ctx)
	return ctx, nil
}

type Context struct {
	AuthToken                string   `policy:"true" json:"auth_token"`
	UserId                   string   `policy:"true" json:"user_id"`
	TenantId                 string   `policy:"true" json:"tenant_id"`
	DomainId                 string   `policy:"true" json:"domain_id"`
	UserDomainId             string   `policy:"true" json:"user_domain_id"`
	ProjectDomainId          string   `policy:"true" json:"project_domain_id"`
	IsAdmin                  bool     `policy:"true" json:"is_admin"`
	ReadOnly                 string   `policy:"true" json:"read_only"`
	ShowDeleted              string   `policy:"true" json:"show_deleted"`
	RequestId                string   `policy:"true" json:"request_id"`
	ResourceUuid             string   `policy:"true" json:"resource_uuid"`
	Overwrite                string   `policy:"true" json:"overwrite"`
	Roles                    []string `policy:"true" json:"roles"`
	UserName                 string   `policy:"true" json:"user_name"`
	ProjectName              string   `policy:"true" json:"project_name"`
	DomainName               string   `policy:"true" json:"domain_name"`
	UserDomainName           string   `policy:"true" json:"user_domain_name"`
	ProjectDomainName        string   `policy:"true" json:"project_domain_name"`
	IsAdminProject           bool     `policy:"true" json:"is_admin_project"`
	ServiceToken             string   `policy:"true" json:"service_token"`
	ServiceUserId            string   `policy:"true" json:"service_user_id"`
	ServiceUserName          string   `policy:"true" json:"service_user_name"`
	ServiceUserDomainId      string   `policy:"true" json:"service_user_domain_id"`
	ServiceUserDomainName    string   `policy:"true" json:"service_user_domain_name"`
	ServiceProjectId         string   `policy:"true" json:"service_project_id"`
	ServiceProjectName       string   `policy:"true" json:"service_project_name"`
	ServiceProjectDomainId   string   `policy:"true" json:"service_project_domain_id"`
	ServiceProjectDomainName string   `policy:"true" json:"service_project_domain_name"`
	ServiceRoles             string   `policy:"true" json:"service_roles"`
	Token                    string   `policy:"false" json:"token"`
	Uri                      string   `policy:"false" json:"uri"`
}

func (ctx *Context) ToPolicyValue() map[string]interface{} {
	ctxMap := map[string]interface{}{}
	t := reflect.TypeOf(ctx).Elem()
	v := reflect.ValueOf(ctx).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		name := t.Field(i).Tag.Get("json")
		if t.Field(i).Tag.Get("policy") == "false" {
			continue
		}
		if field.Kind() == reflect.String && field.String() == "" {
			continue
		}
		if field.Kind() == reflect.Slice && field.Len() == 0 {
			continue
		}
		if field.Kind() == reflect.Map && field.Len() == 0 {
			continue
		}
		ctxMap[name] = field.Interface()
	}
	return ctxMap
}

func (ctx *Context) ToJson() string {
	b, err := json.Marshal(ctx)
	if err != nil {
		log.Errorf("context convert to json failed, reason: %v", err)
	}
	return string(b)
}

func HttpError(ctx *context.Context, code int, format string, a ...interface{}) error {
	ctx.Output.SetStatus(code)
	msg := fmt.Sprintf(format, a...)
	ctx.Output.Body(errorStatus(code, msg))
	errInfo := fmt.Errorf("code: %d, reason: %s", code, msg)
	log.Error(errInfo)
	return errInfo
}

func errorStatus(code int, message string) []byte {
	errStatus := &model.ErrorSpec{
		Code:    code,
		Message: message,
	}

	// Mashal the error status.
	body, err := json.Marshal(errStatus)
	if err != nil {
		return []byte("failed to mashal error response: " + err.Error())
	}
	return body
}
