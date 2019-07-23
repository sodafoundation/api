// Copyright 2019 The OpenSDS Authors.
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

/*
This module implements the common data structure.

*/

package model

import (
	"errors"
	"fmt"

	"github.com/emicklei/go-restful"
	log "github.com/golang/glog"
)

const (
	// ErrorBadRequest
	ErrorBadRequest = 400
	ErrorNotFound   = 404
	// ErrorInternalServer
	ErrorInternalServer = 500
)

// ErrorSpec describes Detailed HTTP error response, which consists of a HTTP
// status code, and a custom error message unique for each failure case.
type ErrorSpec struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func HttpError(res *restful.Response, code int, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	res.WriteError(code, errors.New(msg))
	errInfo := fmt.Sprintf("Code:%d, Reason:%s", code, msg)
	log.Error(errInfo)
	return fmt.Errorf(errInfo)
}

// Volume group error
type NotImplementError struct {
	S string
}

func (e *NotImplementError) Error() string {
	return e.S
}

type NotFoundError struct {
	S string
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{S: msg}
}

func (e *NotFoundError) Error() string {
	return e.S
}
