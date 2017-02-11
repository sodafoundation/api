// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package util

import (
	"time"
)

// RFC8601DateTime is a type for decoding and encoding json
// date times that follow RFC 8601 format. The type currently
// decodes and encodes with exactly precision to seconds. If more
// formats of RFC8601 need to be supported additional work
// will be needed.
type RFC8601DateTime struct {
	time.Time
}

// NewDateTime creates a new RFC8601DateTime taking a string as input.
// It must follow the "2006-01-02T15:04:05" pattern.
func NewDateTime(input string) (val RFC8601DateTime, err error) {
	val = RFC8601DateTime{}
	err = val.formatValue(input)
	return val, err
}

// UnmarshalJSON converts the bytes give to a RFC8601DateTime
// Errors will occur if the bytes when converted to a string
// don't match the format "2006-01-02T15:04:05".
func (r *RFC8601DateTime) UnmarshalJSON(data []byte) error {
	return r.formatValue(string(data))
}

// MarshalJSON converts a RFC8601DateTime to a []byte.
func (r RFC8601DateTime) MarshalJSON() ([]byte, error) {
	val := r.Time.Format(format)
	return []byte(val), nil
}

func (r *RFC8601DateTime) formatValue(input string) (err error) {
	timeVal, err := time.Parse(format, input)
	r.Time = timeVal
	return
}

const format = `"2006-01-02T15:04:05"`
