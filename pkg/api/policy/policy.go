// Copyright 2018 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package policy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	bctx "github.com/astaxie/beego/context"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
)

var enforcer *Enforcer

func init() {
	enforcer = NewEnforcer(false)
	RegisterRules(enforcer)
	enforcer.LoadRules(false)
}

type DefaultRule struct {
	Name     string
	CheckStr string
}

func listRules() []DefaultRule {
	return []DefaultRule{
		{Name: "context_is_admin", CheckStr: "role:admin"},
	}
}

func RegisterRules(e *Enforcer) {
	e.RegisterDefaults(listRules())
}

func NewEnforcer(overWrite bool) *Enforcer {
	return &Enforcer{OverWrite: overWrite}
}

type Enforcer struct {
	Rules        map[string]BaseCheck
	DefaultRules []DefaultRule
	OverWrite    bool
}

func (e *Enforcer) RegisterDefaults(rules []DefaultRule) {
	e.DefaultRules = rules
}

func (e *Enforcer) Enforce(rule string, target map[string]string, cred map[string]interface{}) (bool, error) {
	if err := e.LoadRules(false); err != nil {
		return false, err
	}

	toRule, ok := e.Rules[rule]
	if !ok {
		err := fmt.Errorf("rule [%s] does not exist", rule)
		return false, err
	}
	return check(toRule, target, cred, *e, ""), nil
}

func (e *Enforcer) Authorize(rule string, target map[string]string, cred map[string]interface{}) (bool, error) {
	return e.Enforce(rule, target, cred)
}

func (e *Enforcer) LoadRules(forcedReload bool) error {
	path := config.CONF.OsdsApiServer.PolicyPath
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	// Load all policy files that in the specified path
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
				err := e.LoadPolicyFile(path, forcedReload, false)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return e.LoadPolicyFile(path, forcedReload, e.OverWrite)
}

func (e *Enforcer) UpdateRules(rules map[string]BaseCheck) {
	if e.Rules == nil {
		e.Rules = make(map[string]BaseCheck)
	}
	for k, c := range rules {
		e.Rules[k] = c
	}
}

func (e *Enforcer) LoadPolicyFile(path string, forcedReload bool, overWrite bool) error {
	// if rules is already set or user doesn't want to force reload, return it.
	if e.Rules != nil && !forcedReload {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		msg := fmt.Sprintf("Read policy file (%s) failed, reason:(%v)", path, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	r, err := NewRules(data, e.DefaultRules)
	if err != nil {
		return err
	}

	if overWrite {
		e.Rules = r.Rules
	} else {
		e.UpdateRules(r.Rules)
	}
	return nil
}

func NewRules(data []byte, defaultRule []DefaultRule) (*Rules, error) {
	r := &Rules{}
	err := r.Load(data, defaultRule)
	return r, err
}

type Rules struct {
	Rules map[string]BaseCheck
}

func (r *Rules) Load(data []byte, defaultRules []DefaultRule) error {
	rulesMap := map[string]string{}
	err := json.Unmarshal(data, &rulesMap)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	// add default value
	for _, r := range defaultRules {
		if v, ok := rulesMap[r.Name]; ok {
			log.Warningf("Policy rule (%s:%s) has conflict with default rule(%s:%s),abandon default value\n",
				r.Name, v, r.Name, r.CheckStr)
		} else {
			rulesMap[r.Name] = r.CheckStr
		}
	}

	if r.Rules == nil {
		r.Rules = make(map[string]BaseCheck)
	}
	for k, v := range rulesMap {
		r.Rules[k] = parseRule(v)
	}
	return nil
}

func (r *Rules) String() string {
	b, _ := json.MarshalIndent(r.Rules, "", "  ")
	return string(b)
}

func Authorize(httpCtx *bctx.Context, action string) bool {
	if config.CONF.AuthStrategy != "keystone" {
		return true
	}
	ctx := context.GetContext(httpCtx)
	credentials := ctx.ToPolicyValue()
	TenantId := httpCtx.Input.Param(":tenantId")

	target := map[string]string{
		"tenant_id": TenantId,
	}

	log.V(8).Infof("Action: %v", action)
	log.V(8).Infof("Target: %v", target)
	log.V(8).Infof("Credentials: %v", credentials)
	ok, err := enforcer.Authorize(action, target, credentials)
	if err != nil {
		log.Errorf("Authorize failed, %s", err)
	}
	if !ok {
		context.HttpError(httpCtx, http.StatusForbidden, "operation is not permitted")
	} else {
		ctx.IsAdmin = utils.Contained("admin", ctx.Roles)
	}
	return ok
}
