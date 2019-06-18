// Copyright 2018 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package dorado

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/opensds/opensds/pkg/utils"
)

func TestEncodeName(t *testing.T) {
	id := "05935681-8a00-4988-bfd8-90fdb429aecd"
	exspect := "05935681-477ef4d6bb4af7652c1b97"
	result := EncodeName(id)
	if result != exspect {
		t.Error("Test EncodeName failed")
	}
	if len(result) > MaxNameLength {
		t.Error("EncodeName exceed the max name length")
	}
}

func TestEncodeHostName(t *testing.T) {
	normalName := "1234567890ABCabcZz_.-"
	result := EncodeHostName(normalName)
	if result != normalName {
		t.Error("Test EncodeHostName failed")
	}
	if len(result) > MaxNameLength {
		t.Error("EncodeName exceed the max name length")
	}

	longName := "opensds-huawei-dorado-opensds-huawei-dorado"
	result = EncodeHostName(longName)
	if result != "5620c8980c702896b3c719b187c5bfa" {
		t.Error("Test EncodeHostName failed")
	}
	if len(result) > MaxNameLength {
		t.Error("EncodeName exceed the max name length")
	}

	invalidName := "iqn.1993-08.org.debian:01:d1f6c8e930e7"
	result = EncodeHostName(invalidName)
	if result != "7b1d1cdfe7761ae3e7663ff76343ddc" {
		t.Error("Test EncodeHostName failed")
	}
	if len(result) > MaxNameLength {
		t.Error("EncodeName exceed the max name length")
	}
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestTruncateDescription(t *testing.T) {
	normalDescription := "This is huawei dorado driver testing"
	result := TruncateDescription(normalDescription)
	if result != normalDescription {
		t.Error("Test TruncateDescription failed")
	}
	if len(result) > MaxDescriptionLength {
		t.Error("TruncateDescription exceed the max name length")
	}

	longDescription := randSeq(MaxDescriptionLength + 1)
	result = TruncateDescription(longDescription)
	if len(result) > MaxDescriptionLength {
		t.Error("TruncateDescription exceed the max name length")
	}

	longDescription = randSeq(MaxDescriptionLength + 255)
	result = TruncateDescription(longDescription)
	if len(result) > MaxDescriptionLength {
		t.Error("TruncateDescription exceed the max name length")
	}
}

func TestWaitForCondition(t *testing.T) {
	var count = 0
	err := utils.WaitForCondition(func() (bool, error) {
		count++
		time.Sleep(2 * time.Microsecond)
		if count >= 5 {
			return true, nil
		}
		return false, nil
	}, 1*time.Microsecond, 100*time.Second)
	if err != nil {
		t.Errorf("Test WaitForCondition failed, %v", err)
	}

	count = 0
	err = utils.WaitForCondition(func() (bool, error) {
		count++
		time.Sleep(1 * time.Millisecond)
		if count >= 5 {
			return true, nil
		}
		return false, nil
	}, 4*time.Millisecond, 100*time.Millisecond)
	if err != nil {
		t.Errorf("Test WaitForCondition failed, %v", err)
	}

	err = utils.WaitForCondition(func() (bool, error) {
		return true, fmt.Errorf("test error....")
	}, 4*time.Millisecond, 100*time.Millisecond)
	if err == nil {
		t.Errorf("Test WaitForCondition failed, %v", err)
	}

	count = 0
	err = utils.WaitForCondition(func() (bool, error) {
		count++
		time.Sleep(2 * time.Millisecond)
		if count >= 5 {
			return true, nil
		}
		return false, nil
	}, 2*time.Millisecond, 5*time.Millisecond)
	if err == nil {
		t.Errorf("Test WaitForCondition failed, %v", err)
	}
}
