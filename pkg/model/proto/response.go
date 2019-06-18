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

package proto

import (
	"encoding/json"
	"fmt"
)

func GenericResponseResult(resMsg interface{}) *GenericResponse {
	var msg string
	switch resMsg.(type) {
	case nil:
		break
	case string:
		msg = resMsg.(string)
		break
	default:
		msgJSON, _ := json.Marshal(resMsg)
		msg = string(msgJSON)
	}

	return &GenericResponse{
		Reply: &GenericResponse_Result_{
			Result: &GenericResponse_Result{
				Message: msg,
			},
		},
	}

}

func GenericResponseError(errMsg interface{}) *GenericResponse {
	return &GenericResponse{
		Reply: &GenericResponse_Error_{
			Error: &GenericResponse_Error{
				Code:        "400",
				Description: fmt.Sprint(errMsg),
			},
		},
	}
}
