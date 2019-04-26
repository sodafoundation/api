// Copyright (c) 2019 OpenSDS Authors All Rights Reserved.
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
This module implements the alert data structure.

*/

package model

import "time"

type LabelSet map[string]string

/*
AlertSpec is a data structure that models the alert information
*/
type AlertSpec struct {

	// generator URL
	// Format: uri
	GeneratorURL string `json:"generatorURL,omitempty"`

	// labels
	// Required: true
	Labels LabelSet `json:"labels"`
}

/*
PostableAlertSpec is a data structure that models the alert information to be sent to the Prometheus Alert Manager
*/
type PostableAlertSpec struct {

	// annotations
	Annotations LabelSet `json:"annotations,omitempty"`

	// ends at
	// Format: date-time
	EndsAt time.Time `json:"endsAt,omitempty"`

	// starts at
	// Format: date-time
	StartsAt time.Time `json:"startsAt,omitempty"`

	AlertSpec
}
