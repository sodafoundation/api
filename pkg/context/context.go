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
	"github.com/astaxie/beego/context"
	"reflect"
)

type Context struct {
	AuthToken                string   `json:"auth_token"`
	UserId                   string   `json:"user_id"`
	ProjectId                string   `json:"project_id"`
	DomainId                 string   `json:"domain_id"`
	UserDomainId             string   `json:"user_domain_id"`
	ProjectDomainId          string   `json:"project_domain_id"`
	IsAdmin                  bool     `json:"is_admin"`
	ReadOnly                 string   `json:"read_only"`
	ShowDeleted              string   `json:"show_deleted"`
	RequestId                string   `json:"request_id"`
	ResourceUuid             string   `json:"resource_uuid"`
	Overwrite                string   `json:"overwrite"`
	Roles                    []string `json:"roles"`
	UserName                 string   `json:"user_name"`
	ProjectName              string   `json:"project_name"`
	DomainName               string   `json:"domain_name"`
	UserDomainName           string   `json:"user_domain_name"`
	ProjectDomainName        string   `json:"project_domain_name"`
	IsAdminProject           bool     `json:"is_admin_project"`
	ServiceToken             string   `json:"service_token"`
	ServiceUserId            string   `json:"service_user_id"`
	ServiceUserName          string   `json:"service_user_name"`
	ServiceUserDomainId      string   `json:"service_user_domain_id"`
	ServiceUserDomainName    string   `json:"service_user_domain_name"`
	ServiceProjectId         string   `json:"service_project_id"`
	ServiceProjectName       string   `json:"service_project_name"`
	ServiceProjectDomainId   string   `json:"service_project_domain_id"`
	ServiceProjectDomainName string   `json:"service_project_domain_name"`
	ServiceRoles             string   `json:"service_roles"`
}

func (ctx *Context) ToPolicyValue() map[string]interface{} {
	ctxMap := map[string]interface{}{}
	t := reflect.TypeOf(ctx).Elem()
	v := reflect.ValueOf(ctx).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		name := t.Field(i).Tag.Get("json")
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

func CreateAdminContext() *Context {
	return &Context{
		IsAdmin: true,
	}
}

func CreateInternalTenantContext(projectId, userId string) *Context {
	return &Context{
		ProjectId: projectId,
		UserId:    userId,
		IsAdmin:   true,
	}
}

func GetContext(httpCtx *context.Context) *Context {
	return httpCtx.Input.GetData("context").(*Context)
}
