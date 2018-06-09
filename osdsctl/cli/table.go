// This source file has been modified by Huawei Technologies Co., Ltd.
// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.

// Copyright 2017 modood. All rights reserved.
// license that can be found in the LICENSE file.

// Package table produces a string that represents slice of structs data in a text table

package cli

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/opensds/opensds/pkg/model"
)

type bd struct {
	H  rune // BOX DRAWINGS HORIZONTAL
	V  rune // BOX DRAWINGS VERTICAL
	VH rune // BOX DRAWINGS VERTICAL AND HORIZONTAL
	HU rune // BOX DRAWINGS HORIZONTAL AND UP
	HD rune // BOX DRAWINGS HORIZONTAL AND DOWN
	VL rune // BOX DRAWINGS VERTICAL AND LEFT
	VR rune // BOX DRAWINGS VERTICAL AND RIGHT
	DL rune // BOX DRAWINGS DOWN AND LEFT
	DR rune // BOX DRAWINGS DOWN AND RIGHT
	UL rune // BOX DRAWINGS UP AND LEFT
	UR rune // BOX DRAWINGS UP AND RIGHT
}

type FormatterList map[string]func(v interface{}) string
type KeyList []string
type StructElemCb func(name string, value reflect.Value) error

var m = bd{'-', '|', '+', '+', '+', '+', '+', '+', '+', '+', '+'}

func JsonFormatter(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", " ")
	return string(b)
}

// Output formats slice of structs data and writes to standard output.(Using box drawing characters)
func PrintList(slice interface{}, keys KeyList, fmts FormatterList) {
	fmt.Println(TableList(slice, keys, fmts))
}

func PrintDict(u interface{}, keys KeyList, fmts FormatterList) {
	fmt.Println(TableDict(u, keys, fmts))
}

// Table formats slice of structs data and returns the resulting string.(Using box drawing characters)
func TableList(slice interface{}, keys KeyList, fmts FormatterList) string {
	coln, colw, rows := parseList(slice, keys, fmts)
	table := table(coln, colw, rows, m)
	return table
}

// Table formats slice of structs data and returns the resulting string.(Using standard ascii characters)
func TableDict(u interface{}, keys KeyList, fmts FormatterList) string {
	coln, colw, rows := parseDict(u, keys, fmts)
	table := table(coln, colw, rows, m)
	return table
}

func slice2map(slice []string) map[string]int {
	m := make(map[string]int)
	// increment map's value for every key from slice
	for _, s := range slice {
		m[s]++
	}
	return m
}

func visitStructElem(u interface{}, keys KeyList, fn StructElemCb) {
	v := reflect.ValueOf(u)
	t := reflect.TypeOf(u)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("Table: items of slice should be on struct value")
	}

	whiteList := slice2map(keys)
	m := 0 // count of unexported field
	for n := 0; n < v.NumField(); n++ {
		if t.Field(n).PkgPath != "" {
			m++
			continue
		}

		cn := t.Field(n).Name
		if _, ok := whiteList[cn]; !ok {
			m++
			continue
		}
		if err := fn(cn, v); err != nil {
			panic(fmt.Sprintln("Table:", err))
		}
	}
}

func appendRow(rows [][]string, row []string) [][]string {
	maxRowNum := 0
	var items [][]string
	for _, v := range row {
		lines := strings.Split(v, "\n")
		if maxRowNum < len(lines) {
			maxRowNum = len(lines)
		}
		items = append(items, lines)
	}
	for i := 0; i < maxRowNum; i++ {
		var r []string
		for _, v := range items {
			if len(v) <= i {
				r = append(r, "")
			} else {
				r = append(r, v[i])
			}
		}
		rows = append(rows, r)
	}
	return rows
}

func getRow(u interface{}, keys KeyList, fmts FormatterList) (
	row []string, // rows of content
) {
	visitStructElem(u, keys, func(name string, value reflect.Value) error {
		var cv string
		if fn, ok := fmts[name]; ok {
			cv = fn(value.FieldByName(name).Interface())
		} else {
			cv = fmt.Sprintf("%+v", value.FieldByName(name).Interface())
		}
		row = append(row, cv)
		return nil
	})
	return
}

func getHead(u interface{}, keys KeyList) (
	n []string, // rows of content
) {
	visitStructElem(u, keys, func(name string, value reflect.Value) error {
		n = append(n, name)
		return nil
	})
	return
}

func mergeStrSlice(s ...[]string) (slice []string) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := mergeStrSlice(s[1:]...)
		slice = make([]string, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}
	return
}

func getColw(head []string, rows [][]string) (colw []int) {
	for _, v := range head {
		colw = append(colw, len(v))
	}
	for _, row := range rows {
		for i, v := range row {
			if colw[i] < len(v) {
				colw[i] = len(v)
			}
		}
	}
	return
}

func parseDict(u interface{}, keys KeyList, fmts FormatterList) (
	coln []string, // name of columns
	colw []int, // width of columns
	rows [][]string, // rows of content
) {
	v := reflect.ValueOf(u)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	bv := v.FieldByName("BaseModel")
	if bv.Kind() != reflect.Invalid {
		bm := bv.Interface().(*model.BaseModel)
		bmHead := getHead(bm, keys)
		bmRow := getRow(bm, keys, fmts)
		for i := 0; i < len(bmHead); i++ {
			rows = appendRow(rows, []string{bmHead[i], bmRow[i]})
		}
	}

	head := getHead(u, keys)
	row := getRow(u, keys, fmts)
	for i := 0; i < len(head); i++ {
		rows = appendRow(rows, []string{head[i], row[i]})
	}
	coln = []string{"Property", "Value"}
	colw = getColw(coln, rows)
	return coln, colw, rows
}

func sliceconv(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		panic("sliceconv: param \"slice\" should be on slice value")
	}
	l := v.Len()
	r := make([]interface{}, l)
	for i := 0; i < l; i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}

func parseList(slice interface{}, keys KeyList, fmts FormatterList) (
	coln []string, // name of columns
	colw []int, // width of columns
	rows [][]string, // rows of content
) {
	for _, u := range sliceconv(slice) {
		v := reflect.ValueOf(u)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		head := getHead(u, keys)
		row := getRow(u, keys, fmts)

		bv := v.FieldByName("BaseModel")
		if bv.Kind() != reflect.Invalid {
			bm := bv.Interface().(*model.BaseModel)
			bmHead := getHead(bm, keys)
			bmRow := getRow(bm, keys, fmts)
			coln = mergeStrSlice(bmHead, head)
			rows = appendRow(rows, mergeStrSlice(bmRow, row))
		} else {
			coln = head
			rows = appendRow(rows, row)
		}
	}
	if len(coln) == 0 {
		coln = keys
	}
	colw = getColw(coln, rows)
	return coln, colw, rows
}

func repeat(time int, char rune) string {
	var s = make([]rune, time)
	for i := range s {
		s[i] = char
	}
	return string(s)
}

func table(coln []string, colw []int, rows [][]string, b bd) (table string) {
	head := [][]rune{[]rune{b.DR}, []rune{b.V}, []rune{b.VR}}
	bttm := []rune{b.UR}
	for i, v := range colw {
		head[0] = append(head[0], []rune(repeat(v+2, b.H)+string(b.HD))...)
		head[1] = append(head[1], []rune(" "+coln[i]+repeat(v-len(coln[i])+1, ' ')+string(b.V))...)
		head[2] = append(head[2], []rune(repeat(v+2, b.H)+string(b.VH))...)
		bttm = append(bttm, []rune(repeat(v+2, b.H)+string(b.HU))...)
	}
	head[0][len(head[0])-1] = b.DL
	head[2][len(head[2])-1] = b.VL
	bttm[len(bttm)-1] = b.UL

	var body [][]rune
	for _, r := range rows {
		row := []rune{b.V}
		for i, v := range colw {
			// handle non-ascii character
			lb := len(r[i])
			lr := len([]rune(r[i]))

			row = append(row, []rune(" "+r[i]+repeat(v-lb+(lb-lr)/2+1, ' ')+string(b.V))...)
		}
		body = append(body, row)
	}

	for _, v := range head {
		table += string(v) + "\n"
	}
	for _, v := range body {
		table += string(v) + "\n"
	}
	table += string(bttm)
	return table
}
