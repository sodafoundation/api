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

/*
This module implements the common data structure.

*/

package model

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/context"
	log "github.com/golang/glog"
)

const (
	// ErrorBadRequest
	ErrorBadRequest = 400
	// ErrorUnauthorized
	ErrorUnauthorized = 401
	// ErrorForbidden
	ErrorForbidden = 403
	// ErrorNotFound
	ErrorNotFound = 404
	// ErrorInternalServer
	ErrorInternalServer = 500
	// ErrorNotImplemented
	ErrorNotImplemented = 501
)

// ErrorSpec describes Detailed HTTP error response, which consists of a HTTP
// status code, and a custom error message unique for each failure case.
type ErrorSpec struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ErrorBadRequestStatus
func ErrorBadRequestStatus(message string) []byte {
	return errorStatus(ErrorBadRequest, message)
}

// ErrorForbiddenStatus
func ErrorForbiddenStatus(message string) []byte {
	return errorStatus(ErrorForbidden, message)
}

// ErrorUnauthorizedStatus
func ErrorUnauthorizedStatus(message string) []byte {
	return errorStatus(ErrorUnauthorized, message)
}

// ErrorNotFoundStatus
func ErrorNotFoundStatus(message string) []byte {
	return errorStatus(ErrorNotFound, message)
}

// ErrorInternalServerStatus
func ErrorInternalServerStatus(message string) []byte {
	return errorStatus(ErrorInternalServer, message)
}

// ErrorNotImplementedStatus
func ErrorNotImplementedStatus(message string) []byte {
	return errorStatus(ErrorNotImplemented, message)
}

func errorStatus(code int, message string) []byte {
	errStatus := &ErrorSpec{
		Code:    code,
		Message: message,
	}

	// Mashal the error status.
	body, err := json.Marshal(errStatus)
	if err != nil {
		return []byte("Failed to mashal error response: " + err.Error())
	}
	return body
}

func HttpError(ctx *context.Context, code int, format string, a ...interface{}) error {
	ctx.Output.SetStatus(code)
	msg := fmt.Sprintf(format, a...)
	ctx.Output.Body(errorStatus(code, msg))
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
