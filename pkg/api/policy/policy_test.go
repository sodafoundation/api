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
	"strings"
	"testing"
	"time"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Roles struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Token struct {
	AuditIds  []string      `json:"audit_ids"`
	Catalog   []interface{} `json:"catalog"`
	ID        string        `json:"id"`
	ExpiresAt time.Time     `json:"expires_at"`
	IsAdmin   bool          `json:"is_domain"`
	Project   Project       `json:"project"`
	User      User          `json:"user"`
	Roles     []Roles       `json:"roles"`
}

func TestPolicy(t *testing.T) {
	p := "./testdata/token.json"
	body, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Printf("Read token json file (%s) failed, reason:(%v)\n", p, err)
		return
	}

	var m map[string]interface{}
	err = json.Unmarshal([]byte(body), &m)
	var to Token
	b, err := json.Marshal(m["token"])
	err = json.Unmarshal(b, &to)

	p = "./testdata/policy.json"
	data, err := ioutil.ReadFile(p)
	if err != nil {
		fmt.Printf("Read token json file (%s) failed, reason:(%v)\n", p, err)
		return
	}
	target := map[string]string{"project_id": to.Project.ID}

	var roles []string
	for _, v := range to.Roles {
		roles = append(roles, v.Name)
	}
	cred := map[string]interface{}{
		"roles":      roles,
		"project_id": to.Project.ID,
		"is_admin":   to.IsAdmin,
	}
	// The golang testing framework dosen't invoke init function in linux system,so invoke it.
	registerAll()
	rules, _ := NewRules(data, listRules())
	enforcer := NewEnforcer(false)
	enforcer.Rules = rules.Rules
	expected := map[string]bool{
		"volume:create":  true,
		"volume:delete":  true,
		"volume:get":     false,
		"volume:get_all": true,
	}
	for k, r := range rules.Rules {
		if strings.Contains(k, ":") {
			result := r.Exec(target, cred, *enforcer, "")
			if result != expected[k] {
				t.Errorf("Policy checked failed,\"%s\": \"%s\", expected:%v, got:%v", k, r, expected[k], result)
			}
		}
	}
}
