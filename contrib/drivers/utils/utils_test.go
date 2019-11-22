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

package utils

import (
	"testing"

	pb "github.com/opensds/opensds/pkg/model/proto"
)

func TestGetInitiatorName(t *testing.T) {

	fakeInitiators := []*pb.Initiator{
		&pb.Initiator{
			PortName: "fake1",
			Protocol: "iscsi",
		},
		&pb.Initiator{
			PortName: "fake2",
			Protocol: "fibre_channel",
		},
	}

	testCases := []struct {
		initiators []*pb.Initiator
		protocol   string
		expected   string
	}{
		{
			initiators: fakeInitiators,
			protocol:   "iscsi",
			expected:   "fake1",
		},
		{
			initiators: fakeInitiators,
			protocol:   "fibre_channel",
			expected:   "fake2",
		},
		{
			initiators: fakeInitiators,
			protocol:   "fake_protocol",
			expected:   "",
		},
	}

	for _, c := range testCases {
		actual := GetInitiatorName(c.initiators, c.protocol)
		if actual != c.expected {
			t.Errorf("Expected %v, get %v", c.expected, actual)
		}
	}

}
