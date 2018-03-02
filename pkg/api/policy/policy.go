// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package policy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/api/filter"
	"github.com/opensds/opensds/pkg/utils/config"
)

var enforcer *Enforcer

func init() {
	enforcer = NewEnforcer(false)
}
func NewEnforcer(overWrite bool) *Enforcer {
	return &Enforcer{OverWrite: overWrite}
}

type Enforcer struct {
	Rules       map[string]BaseCheck
	DefaultRule string
	OverWrite   bool
}

func (e *Enforcer) Enforce(rule string, target map[string]string, cred map[string]interface{}) (bool, error) {
	if err := e.LoadRules(false); err != nil {
		return false, err
	}

	toRule, ok := e.Rules[rule]
	if !ok {
		err := fmt.Errorf("Rule [%s] does not exist", rule)
		return false, err
	}
	return check(toRule, target, cred, *e, ""), nil
}

func (e *Enforcer) Authorize(rule string, target map[string]string, cred map[string]interface{}) (bool, error) {
	return e.Enforce(rule, target, cred)
}

func (e *Enforcer) LoadRules(forcedReload bool) error {
	path := config.CONF.OsdsLet.PolicyPath
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
	} else {
		return e.LoadPolicyFile(path, forcedReload, e.OverWrite)
	}
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
		msg := fmt.Sprintf("Read policy file file (%s) failed, reason:(%v)", path, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

	rules, err := NewRules(data, e.DefaultRule)
	if err != nil {
		return err
	}

	if overWrite {
		e.Rules = rules.Rules
	} else {
		e.UpdateRules(rules.Rules)
	}
	return nil
}

func NewRules(data []byte, defaultRule string) (*Rules, error) {
	rules := &Rules{}
	err := rules.Load(data, defaultRule)
	return rules, err
}

type Rules struct {
	Rules       map[string]BaseCheck
	DefaultRule string
}

func (r *Rules) Load(data []byte, defaultRule string) error {
	rulesMap := map[string]string{}
	err := json.Unmarshal(data, &rulesMap)
	if err != nil {
		err := fmt.Errorf("Json unmarshal failed:", err)
		log.Errorf(err.Error())
		return err
	}

	if r.Rules == nil {
		r.Rules = make(map[string]BaseCheck)
	}
	for k, v := range rulesMap {
		r.Rules[k] = parseRule(v)
	}
	r.DefaultRule = defaultRule
	return nil
}

func (r *Rules) String() string {
	b, _ := json.MarshalIndent(r.Rules, "", "  ")
	return string(b)
}

func Authorize(ctx *filter.Context, action string) bool {
	credentials := ctx.ToPolicyValue()
	target := map[string]string{
		"project_id": ctx.ProjectId,
		"user_id":    ctx.UserId,
	}
	log.V(8).Infof("Action: %v", action)
	log.V(8).Infof("Target: %v", target)
	log.V(8).Infof("credentials: %v", credentials)
	ok, err := enforcer.Authorize(action, target, credentials)
	if err != nil {
		log.Error("Authorize failed, %s", err)
	}
	return ok
}
