// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package config

import (
	gflag "flag"
	"reflect"
	"time"
)

type Flag struct {
	InValue interface{}
	Value   interface{}
}

type FlagSet struct {
	flagMap map[string]*Flag
}

type Value interface {
	Set(string) error
}

func (f *FlagSet) BoolVar(p *bool, name string, defValue bool, usage string) {
	inVal := new(bool)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.BoolVar(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) IntVar(p *int, name string, defValue int, usage string) {
	inVal := new(int)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.IntVar(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) Int64Var(p *int64, name string, defValue int64, usage string) {
	inVal := new(int64)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.Int64Var(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) UintVar(p *uint, name string, defValue uint, usage string) {
	inVal := new(uint)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.UintVar(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) Uint64Var(p *uint64, name string, defValue uint64, usage string) {
	inVal := new(uint64)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.Uint64Var(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) Float64Var(p *float64, name string, defValue float64, usage string) {
	inVal := new(float64)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.Float64Var(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) StringVar(p *string, name string, defValue string, usage string) {
	inVal := new(string)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.StringVar(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) DurationVar(p *time.Duration, name string, defValue time.Duration, usage string) {
	inVal := new(time.Duration)
	flag := &Flag{Value: p, InValue: inVal}
	gflag.DurationVar(inVal, name, defValue, usage)
	f.Add(name, flag)
}

func (f *FlagSet) Add(name string, flag *Flag) {
	if f.flagMap == nil {
		f.flagMap = make(map[string]*Flag)
	}
	f.flagMap[name] = flag
}

func (f *FlagSet) Parse() {
	gflag.Parse()
}

func (f *FlagSet) AssignValue() {
	var actual []string
	gflag.CommandLine.Visit(func(flag *gflag.Flag) {
		actual = append(actual, flag.Name)
	})
	for _, name := range actual {
		if _, ok := f.flagMap[name]; !ok {
			continue
		}
		iv := reflect.ValueOf(f.flagMap[name].InValue).Elem()
		v := reflect.ValueOf(f.flagMap[name].Value).Elem()
		v.Set(iv)
	}
}
