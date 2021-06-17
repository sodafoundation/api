// Copyright 2021 The SODA Foundation Authors.
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
	"github.com/sodafoundation/api/pkg/api/policy"
)

type AkSkPortal struct {
	BasePortal
}

func NewAkSkPortal() *AkSkPortal {
	return &AkSkPortal{}
}

func (p *AkSkPortal) ListAkSks() {
	if !policy.Authorize(p.Ctx, "AkSk:list") {
		return
	}

}

func (p *AkSkPortal) CreateAkSk() {
	if !policy.Authorize(p.Ctx, "AkSk:create") {
		return
	}

}

func (p *AkSkPortal) GetAkSk() {
	if !policy.Authorize(p.Ctx, "AkSk:get") {
		return
	}
}

func (p *AkSkPortal) DeleteAkSk() {
	if !policy.Authorize(p.Ctx, "AkSk:delete") {
		return
	}

}

