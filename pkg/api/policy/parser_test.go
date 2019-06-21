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
	"reflect"
	"testing"
)

func TestParseCheck(t *testing.T) {
	check := parseCheck("!")
	if reflect.TypeOf(check) != reflect.TypeOf(&FalseCheck{}) {
		t.Errorf("Parse check \"!\" failed")
	}

	check = parseCheck("@")
	if reflect.TypeOf(check) != reflect.TypeOf(&TrueCheck{}) {
		t.Errorf("Parse check \"@\" failed")
	}

	check = parseCheck("rule:handler")
	if reflect.TypeOf(check) != reflect.TypeOf(&RuleCheck{}) {
		t.Errorf("Parse check rule failed")
	}

	check = parseCheck("role:handler")
	if reflect.TypeOf(check) != reflect.TypeOf(&RoleCheck{}) {
		t.Errorf("Parse check role failed")
	}
	check = parseCheck("no:handler")
	if reflect.TypeOf(check) != reflect.TypeOf(&GenericCheck{}) {
		t.Errorf("Parse check generic failed")
	}
	check = parseCheck("foobar")
	if reflect.TypeOf(check) != reflect.TypeOf(&FalseCheck{}) {
		t.Errorf("Parse check bad rule failed")
	}
	delete(registeredChecks, "Generic")
	check = parseCheck("foobar")
	if reflect.TypeOf(check) != reflect.TypeOf(&FalseCheck{}) {
		t.Errorf("Parse check bad rule failed")
	}
}
func TestParseTokenize(t *testing.T) {
	exemplar := ("(( ( ((() And)) or ) (check:%(miss)s) not)) 'a-string' \"another-string\"")
	results := parseTokenize(exemplar)
	expected := []TokenPair{
		{"(", "("}, {"(", "("}, {"(", "("}, {"(", "("},
		{"(", "("}, {"(", "("}, {")", ")"}, {"and", "And"},
		{")", ")"}, {")", ")"}, {"or", "or"}, {")", ")"},
		{"(", "("}, {"check", "check:%(miss)s"}, {")", ")"},
		{"not", "not"}, {")", ")"}, {")", ")"},
		{"string", "a-string"}, {"string", "another-string"}}
	// please synchronized update the index when if modified the unit test.
	results[13].value = "check:%(miss)s"
	for i := range results {
		if results[i].token != expected[i].token || results[i].value != expected[i].value {
			t.Errorf("Test parseTokenize failed, results:%v, expected:%v", results[i], expected[i])
		}
	}
}

func TestParseState(t *testing.T) {
	state := NewParseState()
	state.tokens = []string{"tok2"}
	state.values = []interface{}{"val2"}
	state.reduce()
}

func TestParseRule(t *testing.T) {
	// test case 1: "a or b or c".
	result := parseRule("@ or ! or @")
	if result.String() != "(@ or ! or @)" {
		t.Error("Parse rule in 'a or b or c' case failed")
	}

	// test case 2: "a or b and c".
	result = parseRule("@ or ! and @")
	if result.String() != "(@ or (! and @))" {
		t.Error("Parse rule in 'a or b and c' case failed")
	}

	// test case 3: "a and b or c".
	result = parseRule("@ and ! or @")
	if result.String() != "((@ and !) or @)" {
		t.Error("Parse rule in 'a and b or c' case failed")
	}

	// test case 4: "a and b and c".
	result = parseRule("@ and ! and @")
	if result.String() != "(@ and ! and @)" {
		t.Error("Parse rule in 'a and b and c' case failed")
	}

	// test case 5: "a or b or c or d" .
	result = parseRule("@ or ! or @ or !")
	if result.String() != "(@ or ! or @ or !)" {
		t.Error("Parse rule in 'a or b or c or d' case failed")
	}

	// test case 6: "a or b or c and d" .
	result = parseRule("@ or ! or @ and !")
	if result.String() != "(@ or ! or (@ and !))" {
		t.Error("Parse rule in 'a or b or c and d' case failed")
	}

	// test case 7: "a or b and c or d" .
	result = parseRule("@ or ! and @ or !")
	if result.String() != "(@ or (! and @) or !)" {
		t.Error("Parse rule in 'a or b and c or d' case failed")
	}

	// test case 8: "a or b and c and d" .
	result = parseRule("@ or ! and @ and !")
	if result.String() != "(@ or (! and @ and !))" {
		t.Error("Parse rule in 'a or b and c and d' case failed")
	}

	// test case 9: "a and b or c or d" .
	result = parseRule("@ and ! or @ or !")
	if result.String() != "((@ and !) or @ or !)" {
		t.Error("Parse rule in 'a and b or c or d' case failed")
	}

	// test case 10: "a and b or c and d" .
	result = parseRule("@ and ! or @ and !")
	if result.String() != "((@ and !) or (@ and !))" {
		t.Error("Parse rule in 'a and b or c and d' case failed")
	}

	// test case 11: "a and b and c or d" .
	result = parseRule("@ and ! and @ or !")
	if result.String() != "((@ and ! and @) or !)" {
		t.Error("Parse rule in 'a and b and c or d' case failed")
	}

	// test case 12: "a and b and c and d" .
	result = parseRule("@ and ! and @ and !")
	if result.String() != "(@ and ! and @ and !)" {
		t.Error("Parse rule in 'a and b and c and d' case failed")
	}

	// test case 13: "a and b or with not 1" .
	result = parseRule("not @ and ! or @")
	if result.String() != "((not @ and !) or @)" {
		t.Error("Parse rule in 'a and b or with not 1' case failed")
	}

	// test case 14: "a and b or with not 2" .
	result = parseRule("@ and not ! or @")
	if result.String() != "((@ and not !) or @)" {
		t.Error("Parse rule in 'a and b or with not 2' case failed")
	}

	// test case 15: "a and b or with not 3" .
	result = parseRule("@ and ! or not @")
	if result.String() != "((@ and !) or not @)" {
		t.Error("Parse rule in 'a and b or with not 3' case failed")
	}

	// test case 16: "a and b and c with group" .
	rules := []string{
		"@ and ( ! ) or @",
		"@ and ! or ( @ )",
		"( @ ) and ! or ( @ )",
		"@ and ( ! ) or ( @ )",
		"( @ ) and ( ! ) or ( @ )",
		"( @ and ! ) or @",
		"( ( @ ) and ! ) or @",
		"( @ and ( ! ) ) or @",
		"( ( @ and ! ) ) or @",
		"( @ and ! or @ )"}
	for _, r := range rules {
		result = parseRule(r)
		if result.String() != "((@ and !) or @)" {
			t.Error("Parse rule in 'a and b and c with group' case failed")
		}
	}

	// test case 17: "a and b or c with group and not" .
	rules = []string{
		"not ( @ ) and ! or @",
		"not @ and ( ! ) or @",
		"not @ and ! or ( @ )",
		"( not @ ) and ! or @",
		"( not @ and ! ) or @",
		"( not @ and ! or @ )"}

	for _, r := range rules {
		result = parseRule(r)
		if result.String() != "((not @ and !) or @)" {
			t.Error("Parse rule in 'a and b and c with group and not'case failed")
		}
	}

	// test case 18: "a and b or c with group and not 2" .
	result = parseRule("not @ and ( ! or @ )")
	if result.String() != "(not @ and (! or @))" {
		t.Error("Parse rule in 'a and b or c with group and not 2' case failed")
	}

	// test case 19: "a and b or c with group and not 3" .
	result = parseRule("not ( @ and ! or @ )")
	if result.String() != "not ((@ and !) or @)" {
		t.Error("Parse rule in 'a and b or c with group and not 3' case failed")
	}

	// test case 20: "a and b or c with group and not 4" .
	rules = []string{
		"( @ ) and not ! or @",
		"@ and ( not ! ) or @",
		"@ and not ( ! ) or @",
		"@ and not ! or ( @ )",
		"( @ and not ! ) or @",
		"( @ and not ! or @ )"}

	for _, r := range rules {
		result = parseRule(r)
		if result.String() != "((@ and not !) or @)" {
			t.Error("Parse rule in 'a and b and c with group and not 4'case failed")
		}
	}

	// test case 21: "a and b or c with group and not 5" .
	result = parseRule("@ and ( not ! or @ )")
	if result.String() != "(@ and (not ! or @))" {
		t.Error("Parse rule in 'a and b or c with group and not 5' case failed")
	}

	// test case 22: "a and b or c with group and not 6" .
	result = parseRule("@ and not ( ! or @ )")
	if result.String() != "(@ and not (! or @))" {
		t.Error("Parse rule in 'a and b or c with group and not 6' case failed")
	}

	// test case 23: "a and b or c with group and not 7" .
	rules = []string{
		"( @ ) and ! or not @",
		"@ and ( ! ) or not @",
		"@ and ! or not ( @ )",
		"@ and ! or ( not @ )",
		"( @ and ! ) or not @",
		"( @ and ! or not @ )"}

	for _, r := range rules {
		result = parseRule(r)
		if result.String() != "((@ and !) or not @)" {
			t.Error("Parse rule in 'a and b and c with group and not 7'case failed")
		}
	}

	// test case 24: "a and b or c with group and not 8" .
	result = parseRule("@ and ( ! or not @ )")
	if result.String() != "(@ and (! or not @))" {
		t.Error("Parse rule in 'a and b or c with group and not 8' case failed")
	}
}
