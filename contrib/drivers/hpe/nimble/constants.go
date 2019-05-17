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

package nimble

const (
	DriverName = "hpe_nimble"
)

const (
	ThickLuntype         = 0
	ThinLuntype          = 1
	MaxNameLength        = 31
	MaxDescriptionLength = 170
	PortNumPerContr      = 2
	PwdExpired           = 3
	PwdReset             = 4
)

// Error Code
const (
	ErrorUnauthorizedToServer = "SM_http_unauthorized"
	ErrorSmVolSizeDecreased   = "SM_vol_size_decreased"
	ErrorSmHttpConflict       = "SM_http_conflict"
)
