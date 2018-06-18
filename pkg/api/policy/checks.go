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

package policy

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

func init() {
	registerAll()
}

type NewCheckFunc func(kind string, match string) BaseCheck

var registeredChecks map[string]NewCheckFunc

func register(name string, f NewCheckFunc) {
	registeredChecks[name] = f
}

func registerAll() {
	if registeredChecks == nil {
		registeredChecks = make(map[string]NewCheckFunc)
	}
	register("rule", NewRuleCheck)
	register("role", NewRoleCheck)
	register("generic", NewGenericCheck)
}

type BaseCheck interface {
	String() string
	Exec(target map[string]string, cred map[string]interface{}, enforcer Enforcer, currentRule string) bool
}

func check(rule BaseCheck,
	target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	ret := rule.Exec(target, cred, enforcer, currentRule)
	log.V(8).Infof("check rules:%s -- %v", rule, ret)
	return ret
}

func NewFalseCheck() BaseCheck {
	return &FalseCheck{}
}

type FalseCheck struct{}

func (this *FalseCheck) String() string {
	return "!"
}

func (this *FalseCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	return false
}

func NewTrueCheck() BaseCheck {
	return &TrueCheck{}
}

type TrueCheck struct {
	rule string
}

func (this *TrueCheck) String() string {
	return "@"
}

func (this *TrueCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	return true
}

func NewNotCheck(check BaseCheck) *NotCheck {
	return &NotCheck{check}
}

type NotCheck struct {
	rule BaseCheck
}

func (this *NotCheck) String() string {
	return fmt.Sprintf("not %s", this.rule)
}

func (this *NotCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	return !check(this.rule, target, cred, enforcer, currentRule)
}

func NewAndCheck(check1 BaseCheck, check2 BaseCheck) *AndCheck {
	ac := &AndCheck{}
	ac.AddCheck(check1)
	ac.AddCheck(check2)
	return ac
}

type AndCheck struct {
	rules []BaseCheck
}

func (this *AndCheck) String() string {
	var r []string
	for _, rule := range this.rules {
		r = append(r, rule.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(r, " and "))
}

func (this *AndCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	for _, rule := range this.rules {
		if !check(rule, target, cred, enforcer, currentRule) {
			return false
		}
	}
	return true
}

func (this *AndCheck) AddCheck(rule BaseCheck) *AndCheck {
	this.rules = append(this.rules, rule)
	return this
}

func NewOrCheck(check1 BaseCheck, check2 BaseCheck) *OrCheck {
	oc := &OrCheck{}
	oc.AddCheck(check1)
	oc.AddCheck(check2)
	return oc
}

type OrCheck struct {
	rules []BaseCheck
}

func (this *OrCheck) String() string {
	var r []string
	for _, rule := range this.rules {
		r = append(r, rule.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(r, " or "))
}

func (this *OrCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	for _, rule := range this.rules {
		if check(rule, target, cred, enforcer, currentRule) {
			return true
		}
	}
	return false
}

func (this *OrCheck) AddCheck(rule BaseCheck) *OrCheck {
	this.rules = append(this.rules, rule)
	return this
}

func (this *OrCheck) PopCheck() (*OrCheck, BaseCheck) {
	x := this.rules[len(this.rules)-1]
	this.rules = this.rules[:len(this.rules)-1]
	return this, x
}

func NewRuleCheck(kind string, match string) BaseCheck {
	return &RuleCheck{kind, match}
}

type RuleCheck struct {
	kind  string
	match string
}

func (this *RuleCheck) String() string {
	return fmt.Sprintf("%s:%s", this.kind, this.match)
}

func (this *RuleCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	if len(enforcer.Rules) == 0 {
		return false
	}
	return check(enforcer.Rules[this.match], target, cred, enforcer, currentRule)
}

func keyWorkFormatter(target map[string]string, match string) (string, error) {
	reg := regexp.MustCompile(`%([[:graph:]]+)s`)
	if ms := reg.FindAllString(match, -1); len(ms) == 1 {
		s := ms[0][2 : len(ms[0])-2]
		for key, val := range target {
			if s == key {
				return val, nil
				break
			}
		}
		return "", fmt.Errorf("target key doesn`t match")
	}
	return match, nil
}

func NewRoleCheck(kind string, match string) BaseCheck {
	return &RoleCheck{kind, match}
}

type RoleCheck struct {
	kind  string
	match string
}

func (this *RoleCheck) String() string {
	return fmt.Sprintf("%s:%s", this.kind, this.match)
}

func (this *RoleCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	match, err := keyWorkFormatter(target, this.match)
	if err != nil {
		return false
	}
	if roles, ok := cred["roles"]; ok {
		for _, role := range roles.([]string) {
			if strings.ToLower(match) == strings.ToLower(role) {
				return true
			}
		}
	}
	return false
}

func NewGenericCheck(kind string, match string) BaseCheck {
	return &GenericCheck{kind, match}
}

type GenericCheck struct {
	kind  string
	match string
}

func (this *GenericCheck) String() string {
	return fmt.Sprintf("%s:%s", this.kind, this.match)
}

func (this *GenericCheck) simpleLiteral(expr string) (string, error) {
	s := fmt.Sprintf("%c%c", expr[0], expr[len(expr)-1])
	if len(expr) >= 2 && (s == "\"\"" || s == "''") {
		return expr[1 : len(expr)-1], nil
	}
	if utils.Contained(strings.ToLower(expr), []string{"true", "false"}) {
		return strings.ToLower(expr), nil
	}
	return "", fmt.Errorf("Not support right now")
}

func (this *GenericCheck) findInMap(testVal interface{}, pathSegs []string, match string) bool {
	if len(pathSegs) == 0 {
		switch testVal.(type) {
		case string:
			return strings.ToLower(match) == strings.ToLower(testVal.(string))
		case bool:
			return strings.ToLower(match) == fmt.Sprint(testVal.(bool))
		default:
			return false
		}
	}
	key, pathSegs := pathSegs[0], pathSegs[1:]
	if val, ok := testVal.(map[string]interface{}); ok {
		testVal = val[key]
	} else {
		return false
	}
	if testVal == nil {
		return false
	}

	if reflect.TypeOf(testVal).Kind() == reflect.Slice {
		if vList, ok := testVal.([]interface{}); ok {
			for _, val := range vList {
				if this.findInMap(val, pathSegs, match) {
					return true
				}
			}
		} else {
			for _, val := range testVal.([]string) {
				if this.findInMap(val, pathSegs, match) {
					return true
				}
			}
		}
		return false
	} else {
		return this.findInMap(testVal, pathSegs, match)
	}
}

func (this *GenericCheck) Exec(target map[string]string,
	cred map[string]interface{},
	enforcer Enforcer,
	currentRule string) bool {
	match, err := keyWorkFormatter(target, strings.ToLower(this.match))
	if err != nil {
		return false
	}

	if test_value, err := this.simpleLiteral(this.kind); err == nil {
		return strings.ToLower(match) == test_value
	}
	if len(cred) == 0 {
		return false
	}
	return this.findInMap(cred, strings.Split(this.kind, "."), match)
}
