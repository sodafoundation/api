// Copyright 2017 The OpenSDS Authors.
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
)

const (
	ErrorBadRequest     = 400
	ErrorUnauthorized   = 401
	ErrorForbidden      = 403
	ErrorNotFound       = 404
	ErrorInternalServer = 500
	ErrorNotImplemented = 501
)

// ErrorSpec describes Detailed HTTP error response, which consists of a HTTP
// status code, and a custom error message unique for each failure case.
type ErrorSpec struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func ErrorBadRequestStatus(message string) []byte {
	return errorStatus(ErrorBadRequest, message)
}

func ErrorForbiddenStatus(message string) []byte {
	return errorStatus(ErrorForbidden, message)
}

func ErrorUnauthorizedStatus(message string) []byte {
	return errorStatus(ErrorUnauthorized, message)
}

func ErrorNotFoundStatus(message string) []byte {
	return errorStatus(ErrorNotFound, message)
}

func ErrorInternalServerStatus(message string) []byte {
	return errorStatus(ErrorInternalServer, message)
}

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
