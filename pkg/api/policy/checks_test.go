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
	"fmt"
	"reflect"
	"testing"
)

func newBoolCheck(result bool) BaseCheck {
	return &boolCheck{false, result}
}

type boolCheck struct {
	called bool
	result bool
}

func (b *boolCheck) String() string {
	return fmt.Sprint(b.result)
}

func (b *boolCheck) Exec(target map[string]string, cred map[string]interface{}, enforcer Enforcer, currentRule string) bool {
	b.called = true
	return b.result
}

func TestRuleCheck(t *testing.T) {
	enforcer := Enforcer{}
	check := NewRuleCheck("rule", "spam")
	target := map[string]string{"target": "fake"}
	if check.Exec(target, make(map[string]interface{}), enforcer, "") {
		t.Error("RuleCheck missing rule test failed")
	}

	enforcer.Rules = make(map[string]BaseCheck)
	enforcer.Rules["spam"] = newBoolCheck(false)
	if check.Exec(target, make(map[string]interface{}), enforcer, "") {
		t.Error("RuleCheck rule false test failed")
	}

	enforcer.Rules = make(map[string]BaseCheck)
	enforcer.Rules["spam"] = newBoolCheck(true)
	if !check.Exec(target, make(map[string]interface{}), enforcer, "") {
		t.Error("RuleCheck rule true test failed")
	}
}

func TestRoleCheck(t *testing.T) {
	enforcer := Enforcer{}
	// Test Case 1
	check := NewRoleCheck("role", "sPaM")
	if !check.Exec(map[string]string{},
		map[string]interface{}{"roles": []string{"SpAm"}},
		enforcer, "") {
		t.Error("RoleCheck role accept test failed")
	}
	//Test case 2
	check = NewRoleCheck("role", "spam")
	if check.Exec(map[string]string{},
		map[string]interface{}{"roles": []string{}},
		enforcer, "") {
		t.Error("RoleCheck role reject test failed")
	}

	//Test case 3
	check = NewRoleCheck("role", "%(target.role.name)s")
	if !check.Exec(map[string]string{"target.role.name": "a"},
		map[string]interface{}{"user": "user", "roles": []string{"a", "b", "c"}},
		enforcer, "") {
		t.Error("RoleCheck format value key exist test failed")
	}

	//Test case 4
	check = NewRoleCheck("role", "%(target.role.name)s")
	if check.Exec(map[string]string{"target.role.name": "d"},
		map[string]interface{}{"user": "user", "roles": []string{"a", "b", "c"}},
		enforcer, "") {
		t.Error("RoleCheck format value key doesn`t exist test failed")
	}

	//Test case 5
	check = NewRoleCheck("role", "%(target.role.name)s")
	if check.Exec(map[string]string{},
		map[string]interface{}{},
		enforcer, "") {
		t.Error("RoleCheck format no roles test failed")
	}
}

func TestGenericCheck(t *testing.T) {
	enforcer := Enforcer{}
	// Test case 1: no cred check.
	check := NewGenericCheck("name", "%(name)s")
	if check.Exec(map[string]string{"name": "spam"},
		map[string]interface{}{},
		enforcer, "") {
		t.Error("GenericCheck no cred test failed")
	}
	// Test case 2: no cred check.
	if check.Exec(map[string]string{"name": "spam"},
		map[string]interface{}{"name": "ham"},
		enforcer, "") {
		t.Error("GenericCheck cred mismatch test failed")
	}
	// Test case 3: accept.
	if !check.Exec(map[string]string{"name": "spam"},
		map[string]interface{}{"name": "spam"},
		enforcer, "") {
		t.Error("GenericCheck cred mismatch test failed")
	}
	// Test case 4: no key match in target.
	if check.Exec(map[string]string{"name1": "spam"},
		map[string]interface{}{"name": "spam"},
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	// Test case 5: no key match in target.
	check = NewGenericCheck("'spam'", "%(name)s")
	if !check.Exec(map[string]string{"name": "spam"},
		map[string]interface{}{},
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	// Test case 6: constant literal mismatch.
	check = NewGenericCheck("'spam'", "%(name)s")
	if !check.Exec(map[string]string{"name": "spam"},
		map[string]interface{}{},
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	//Test case 7: test_constant_literal_mismatch
	check = NewGenericCheck("True", "%(enabled)s")
	if check.Exec(map[string]string{"enabled": "False"},
		map[string]interface{}{},
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}
	// Test case 8: test_constant_literal_accept
	check = NewGenericCheck("True", "%(enabled)s")
	if !check.Exec(map[string]string{"enabled": "True"},
		map[string]interface{}{},
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	// Test case 9: test_constant_literal_accept
	check = NewGenericCheck("a.b.c.d", "APPLES")
	cred := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": map[string]interface{}{
					"d": "APPLES",
				},
			},
		},
	}
	if !check.Exec(map[string]string{"enabled": "True"},
		cred,
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	cred = map[string]interface{}{
		"a": "APPLES",
		"o": map[string]interface{}{
			"t": "ORANGES",
		},
	}
	// Test case 10: test_missing_credentials_dictionary_lookup
	check = NewGenericCheck("o.t", "ORANGES")
	if !check.Exec(map[string]string{},
		cred,
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	// Test case 11: test_missing_credentials_dictionary_lookup
	check = NewGenericCheck("o.v", "ORANGES")
	if check.Exec(map[string]string{},
		cred,
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}
	// Test case 12: test_missing_credentials_dictionary_lookup
	check = NewGenericCheck("q.v", "ORANGES")
	if check.Exec(map[string]string{},
		cred,
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	// Test case 13: test_single_entry_in_list_accepted
	cred = map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": map[string]interface{}{
					"d": []string{"APPLES"},
				},
			},
		},
	}
	check = NewGenericCheck("a.b.c.d", "APPLES")
	if !check.Exec(map[string]string{},
		cred,
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}
	// Test case 14: test_multiple_entry_in_list_accepted
	cred = map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": map[string]interface{}{
					"d": []string{"Bananas", "APPLES", "Grapes"},
				},
			},
		},
	}
	check = NewGenericCheck("a.b.c.d", "APPLES")
	if !check.Exec(map[string]string{},
		cred,
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}

	//Test case 15: test_multiple_entry_in_list_accepted
	cred = map[string]interface{}{
		"a": map[string]interface{}{
			"b": []interface{}{
				map[string]interface{}{
					"c": map[string]interface{}{
						"d": []string{"Bananas", "APPLES", "Grapes"},
					},
				},
			},
		},
	}
	check = NewGenericCheck("a.b.c.d", "APPLES")
	if !check.Exec(map[string]string{},
		cred,
		enforcer, "") {
		t.Error("GenericCheck no key match in target test failed")
	}
}

func TestFalseCheck(t *testing.T) {
	check := NewFalseCheck()
	if "!" != check.String() {
		t.Errorf("FalseCheck failed.")
	}
	enforcer := Enforcer{}
	if check.Exec(map[string]string{}, map[string]interface{}{}, enforcer, "") {
		t.Errorf("FalseCheck failed.")
	}
}

func TestTrueCheck(t *testing.T) {
	check := NewTrueCheck()
	if "@" != check.String() {
		t.Errorf("TrueCheck failed.")
	}
	enforcer := Enforcer{}
	if !check.Exec(map[string]string{}, map[string]interface{}{}, enforcer, "") {
		t.Errorf("TrueCheck failed.")
	}
}

func TestNotCheck(t *testing.T) {
	enforcer := Enforcer{}
	check := NewNotCheck(NewTrueCheck())
	if "not @" != check.String() {
		t.Errorf("NotCheck string test failed.")
	}
	if check.Exec(map[string]string{}, map[string]interface{}{}, enforcer, "") {
		t.Errorf("NotCheck exeute test failed.")
	}
}

func TestAndCheck(t *testing.T) {
	// Test case 1
	c1 := NewTrueCheck()
	c2 := NewTrueCheck()
	c3 := NewTrueCheck()
	check := NewAndCheck(c1, c2)
	if !reflect.DeepEqual(check.rules, []BaseCheck{c1, c2}) {
		t.Errorf("AndCheck new test failed")
	}
	check.AddCheck(c3)
	if !reflect.DeepEqual(check.rules, []BaseCheck{c1, c2, c3}) {
		t.Errorf("AndCheck add check test failed")
	}
	if check.String() != "(@ and @ and @)" {
		t.Errorf("AndCheck string test failed")
	}
	//first true
	b1 := newBoolCheck(true)
	b2 := newBoolCheck(false)
	check = NewAndCheck(b1, b2)
	if check.Exec(map[string]string{}, map[string]interface{}{}, Enforcer{}, "") {
		t.Errorf("AndCheck call first true test failed")
	}
	if !(check.rules[0].(*boolCheck).called && check.rules[1].(*boolCheck).called) {
		t.Errorf("AndCheck call first true test failed")
	}

	// second true
	b1 = newBoolCheck(false)
	b2 = newBoolCheck(true)
	check = NewAndCheck(b1, b2)
	if check.Exec(map[string]string{}, map[string]interface{}{}, Enforcer{}, "") {
		t.Errorf("AndCheck call second true test failed")
	}
	if !(check.rules[0].(*boolCheck).called && !check.rules[1].(*boolCheck).called) {
		t.Errorf("AndCheck call second true test failed")
	}
}

func TestOrCheck(t *testing.T) {
	// Test case 1
	c1 := NewTrueCheck()
	c2 := NewTrueCheck()
	c3 := NewTrueCheck()
	check := NewOrCheck(c1, c2)
	if !reflect.DeepEqual(check.rules, []BaseCheck{c1, c2}) {
		t.Errorf("OrCheck new test failed")
	}
	check.AddCheck(c3)
	if !reflect.DeepEqual(check.rules, []BaseCheck{c1, c2, c3}) {
		t.Errorf("OrCheck add check test failed")
	}
	if check.String() != "(@ or @ or @)" {
		t.Errorf("OrCheck string test failed")
	}
	_, check1 := check.PopCheck()
	if !reflect.DeepEqual(check.rules, []BaseCheck{c1, c2}) {
		t.Errorf("OrCheck pop check test failed")
	}
	if !reflect.DeepEqual(check1, c3) {
		t.Errorf("OrCheck pop check test failed")
	}
	// all false
	check = NewOrCheck(newBoolCheck(false), newBoolCheck(false))
	if check.Exec(map[string]string{}, map[string]interface{}{}, Enforcer{}, "") {
		t.Errorf("OrCheck call all false test failed")
	}
	if !(check.rules[0].(*boolCheck).called && check.rules[1].(*boolCheck).called) {
		t.Errorf("OrCheck call all false test failed")
	}

	// first false
	check = NewOrCheck(newBoolCheck(false), newBoolCheck(true))
	if !check.Exec(map[string]string{}, map[string]interface{}{}, Enforcer{}, "") {
		t.Errorf("OrCheck call first false test failed")
	}
	if !(check.rules[0].(*boolCheck).called && check.rules[1].(*boolCheck).called) {
		t.Errorf("OrCheck call first false test failed")
	}

	// second false
	check = NewOrCheck(newBoolCheck(true), newBoolCheck(false))
	if !check.Exec(map[string]string{}, map[string]interface{}{}, Enforcer{}, "") {
		t.Errorf("OrCheck call second false test failed")
	}
	if !(check.rules[0].(*boolCheck).called && !check.rules[1].(*boolCheck).called) {
		t.Errorf("OrCheck call second false test failed")
	}
}
