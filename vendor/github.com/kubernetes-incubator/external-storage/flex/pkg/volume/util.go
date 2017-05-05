/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package volume

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// StorageProfile is a structure for all properties of
// profile configured by admin
type StorageProfile struct {
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	BackendDriver string            `json:"backend"`
	StorageTags   map[string]string `json:"tags"`
}

// VolumeOperationSchema is a structure for all properties of
// volume operation
type VolumeOperationSchema struct {
	// Some properties related to basic operation of volumes
	DockId       string `json:"dockId,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	VolumeType   string `json:"volumeType"`
	Size         int32  `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	// Some properties related to attach and mount operation of volumes
	Device   string `json:"device,omitempty"`
	MountDir string `json:"mountDir,omitempty"`
	FsType   string `json:"fsType,omitempty"`

	// Some properties related to basic operation of volume snapshots
	SnapshotId      string `json:"snapshotId,omitempty"`
	SnapshotName    string `json:"snapshotName,omitempty"`
	Description     string `json:"description,omitempty"`
	ForceSnapshoted bool   `json:"forceSnapshoted,omitempty"`
}

// VolumeRequest is a structure for all properties of
// a volume request
type VolumeRequest struct {
	Schema  *VolumeOperationSchema `json:"schema"`
	Profile *StorageProfile        `json:"profile"`
}

// VolumeResponse is a structure for all properties of
// a volume for a non detailed query
type VolumeResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Status      string              `json:"status"`
	Size        int                 `json:"size"`
	VolumeType  string              `json:"volume_type"`
	Attachments []map[string]string `json:"attachments"`
}

// VolumeDeleteResponse is a structure for deleting
// a volume
type VolumeDeleteResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// generateId method generates a unique exportId to assign an export
func generateId(mutex *sync.Mutex, ids map[uint16]bool) uint16 {
	mutex.Lock()
	id := uint16(1)
	for ; id <= math.MaxUint16; id++ {
		if _, ok := ids[id]; !ok {
			break
		}
	}
	ids[id] = true
	mutex.Unlock()
	return id
}

func deleteId(mutex *sync.Mutex, ids map[uint16]bool, id uint16) {
	mutex.Lock()
	delete(ids, id)
	mutex.Unlock()
}

// getExistingIds populates a map with existing ids found in the given config
// file using the given regexp. Regexp must have a "digits" submatch.
func getExistingIds(config string, re *regexp.Regexp) (map[uint16]bool, error) {
	ids := map[uint16]bool{}

	digitsRe := "([0-9]+)"
	if !strings.Contains(re.String(), digitsRe) {
		return ids, fmt.Errorf("regexp %s doesn't contain digits submatch %s", re.String(), digitsRe)
	}

	read, err := ioutil.ReadFile(config)
	if err != nil {
		return ids, err
	}

	allMatches := re.FindAllSubmatch(read, -1)
	for _, match := range allMatches {
		digits := match[1]
		if id, err := strconv.ParseUint(string(digits), 10, 16); err == nil {
			ids[uint16(id)] = true
		}
	}
	return ids, nil
}

// CheckHTTPResponseStatusCode compares http response header StatusCode against expected
// statuses. Primary function is to ensure StatusCode is in the 20x (return nil).
// Ok: 200. Created: 201. Accepted: 202. No Content: 204. Partial Content: 206.
// Otherwise return error message.
func CheckHTTPResponseStatusCode(resp *http.Response) error {
	switch resp.StatusCode {
	case 200, 201, 202, 204, 206:
		return nil
	case 400:
		return errors.New("Error: response == 400 bad request")
	case 401:
		return errors.New("Error: response == 401 unauthorised")
	case 403:
		return errors.New("Error: response == 403 forbidden")
	case 404:
		return errors.New("Error: response == 404 not found")
	case 405:
		return errors.New("Error: response == 405 method not allowed")
	case 409:
		return errors.New("Error: response == 409 conflict")
	case 413:
		return errors.New("Error: response == 413 over limit")
	case 415:
		return errors.New("Error: response == 415 bad media type")
	case 422:
		return errors.New("Error: response == 422 unprocessable")
	case 429:
		return errors.New("Error: response == 429 too many request")
	case 500:
		return errors.New("Error: response == 500 instance fault / server err")
	case 501:
		return errors.New("Error: response == 501 not implemented")
	case 503:
		return errors.New("Error: response == 503 service unavailable")
	}
	return errors.New("Error: unexpected response status code")
}

// parseQuantityString is a fast scanner for quantity values.
func parseQuantityString(str string) (positive bool, value, num, denom, suffix string, err error) {
	positive = true
	pos := 0
	end := len(str)

	// handle leading sign
	if pos < end {
		switch str[0] {
		case '-':
			positive = false
			pos++
		case '+':
			pos++
		}
	}

	// strip leading zeros
Zeroes:
	for i := pos; ; i++ {
		if i >= end {
			num = "0"
			value = num
			return
		}
		switch str[i] {
		case '0':
			pos++
		default:
			break Zeroes
		}
	}

	// extract the numerator
Num:
	for i := pos; ; i++ {
		if i >= end {
			num = str[pos:end]
			value = str[0:end]
			return
		}
		switch str[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			num = str[pos:i]
			pos = i
			break Num
		}
	}

	// if we stripped all numerator positions, always return 0
	if len(num) == 0 {
		num = "0"
	}

	// handle a denominator
	if pos < end && str[pos] == '.' {
		pos++
	Denom:
		for i := pos; ; i++ {
			if i >= end {
				denom = str[pos:end]
				value = str[0:end]
				return
			}
			switch str[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				denom = str[pos:i]
				pos = i
				break Denom
			}
		}
		// TODO: we currently allow 1.G, but we may not want to in the future.
		// if len(denom) == 0 {
		// 	err = ErrFormatWrong
		// 	return
		// }
	}
	value = str[0:pos]

	// grab the elements of the suffix
	suffixStart := pos
	for i := pos; ; i++ {
		if i >= end {
			suffix = str[suffixStart:end]
			return
		}
		if !strings.ContainsAny(str[i:i+1], "eEinumkKMGTP") {
			pos = i
			break
		}
	}
	if pos < end {
		switch str[pos] {
		case '-', '+':
			pos++
		}
	}
Suffix:
	for i := pos; ; i++ {
		if i >= end {
			suffix = str[suffixStart:end]
			return
		}
		switch str[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			break Suffix
		}
	}
	// we encountered a non decimal in the Suffix loop, but the last character
	// was not a valid exponent
	err = errors.New("quantities must match the regular expression!")
	return
}
