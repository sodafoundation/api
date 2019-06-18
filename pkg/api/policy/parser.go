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
	"strings"

	"github.com/opensds/opensds/pkg/utils"
)

type TokenPair struct {
	token string
	value interface{}
}

func parseCheck(rule string) BaseCheck {
	if rule == "!" {
		return &FalseCheck{}
	} else if rule == "@" {
		return &TrueCheck{}
	}
	items := strings.SplitN(rule, ":", 2)
	if len(items) != 2 {
		return &FalseCheck{}
	}
	kind, match := items[0], items[1]
	if check, ok := registeredChecks[kind]; ok {
		return check(kind, match)
	} else if check, ok := registeredChecks["generic"]; ok {
		return check(kind, match)
	} else {
		return &FalseCheck{}
	}
}

func parseTokenize(rule string) []TokenPair {
	var tokPairs []TokenPair
	for _, tok := range strings.Fields(rule) {
		if tok == "" {
			continue
		}

		clean := strings.TrimLeft(tok, "(")
		for i := 0; i < len(tok)-len(clean); i++ {
			tokPairs = append(tokPairs, TokenPair{"(", "("})
		}

		// If it was only parentheses, continue
		if clean == "" {
			continue
		}

		tok = clean
		// Handle trailing parens on the token
		clean = strings.TrimRight(tok, ")")
		trail := len(tok) - len(clean)
		lowered := strings.ToLower(clean)

		if utils.Contained(lowered, []string{"and", "or", "not"}) {
			tokPairs = append(tokPairs, TokenPair{lowered, clean})
		} else if clean != "" {
			s := fmt.Sprintf("%c%c", tok[0], tok[len(tok)-1])
			if len(tok) >= 2 && (s == "\"\"" || s == "''") {
				tokPairs = append(tokPairs, TokenPair{"string", tok[1 : len(tok)-1]})
			} else {
				tokPairs = append(tokPairs, TokenPair{"check", parseCheck(clean)})
			}
		}

		for i := 0; i < trail; i++ {
			tokPairs = append(tokPairs, TokenPair{")", ")"})
		}
	}

	return tokPairs
}

func parseRule(rule string) BaseCheck {
	if rule == "" {
		return &TrueCheck{}
	}
	state := NewParseState()
	tokPairs := parseTokenize(rule)
	for _, tp := range tokPairs {
		state.Shift(tp.token, tp.value)
	}
	if result, err := state.Result(); err == nil {
		return result.(BaseCheck)
	}
	return &FalseCheck{}
}

var ReduceFuncMap = map[string]ReduceFunc{
	"(,check,)":          wrapCheck,
	"(,and_expr,)":       wrapCheck,
	"(,or_expr,)":        wrapCheck,
	"check,and,check":    makeAndExpr,
	"or_expr,and,check":  mixOrAndExpr,
	"and_expr,and,check": extendAndExpr,
	"check,or,check":     makeOrExpr,
	"and_expr,or,check":  makeOrExpr,
	"or_expr,or,check":   extendOrExpr,
	"not,check":          makeNotExpr,
}

func NewParseState() *ParseState {
	return &ParseState{}
}

type ParseState struct {
	tokens []string
	values []interface{}
}

type ReduceFunc func(args ...interface{}) []TokenPair

func (p *ParseState) reduce() {
	tokenStr := strings.Join(p.tokens, ",")
	for key, fun := range ReduceFuncMap {
		if strings.HasSuffix(tokenStr, key) {
			argNum := strings.Count(key, ",") + 1
			argIdx := len(p.values) - argNum
			args := p.values[argIdx:]
			results := fun(args...)
			p.tokens = append(p.tokens[:argIdx], results[0].token)
			p.values = append(p.values[:argIdx], results[0].value)
			p.reduce()
		}
	}
}

func (p *ParseState) Shift(tok string, val interface{}) {
	p.tokens = append(p.tokens, tok)
	p.values = append(p.values, val)
	p.reduce()
}

func (p *ParseState) Result() (interface{}, error) {
	if len(p.values) != 1 {
		return nil, fmt.Errorf("Could not parse rule")
	}
	return p.values[0], nil
}

func wrapCheck(args ...interface{}) []TokenPair {
	check := args[1].(BaseCheck)
	return []TokenPair{{"check", check}}
}

func makeAndExpr(args ...interface{}) []TokenPair {
	check1 := args[0].(BaseCheck)
	check2 := args[2].(BaseCheck)
	return []TokenPair{{"and_expr", NewAndCheck(check1, check2)}}
}

func mixOrAndExpr(args ...interface{}) []TokenPair {
	orExpr := args[0].(*OrCheck)
	check := args[2].(BaseCheck)
	var andExpr *AndCheck
	orExpr, check1 := orExpr.PopCheck()
	if v, ok := check1.(*AndCheck); ok {
		andExpr = v
		andExpr.AddCheck(check)
	} else {
		andExpr = NewAndCheck(check1, check)
	}
	return []TokenPair{{"or_expr", orExpr.AddCheck(andExpr)}}
}

func extendAndExpr(args ...interface{}) []TokenPair {
	andExpr := args[0].(*AndCheck)
	check2 := args[2].(BaseCheck)
	return []TokenPair{{"and_expr", andExpr.AddCheck(check2)}}
}

func makeOrExpr(args ...interface{}) []TokenPair {
	check1 := args[0].(BaseCheck)
	check2 := args[2].(BaseCheck)
	return []TokenPair{{"or_expr", NewOrCheck(check1, check2)}}
}

func extendOrExpr(args ...interface{}) []TokenPair {
	orExpr := args[0].(*OrCheck)
	check := args[2].(BaseCheck)
	return []TokenPair{{"or_expr", orExpr.AddCheck(check)}}
}

func makeNotExpr(args ...interface{}) []TokenPair {
	return []TokenPair{{"check", NewNotCheck(args[1].(BaseCheck))}}
}
