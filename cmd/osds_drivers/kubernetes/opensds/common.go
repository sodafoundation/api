// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/opensds/opensds/cmd/utils"
)

type Result struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Device  string `json:"device,omitempty"`
}

type DefaultOptions struct {
	Action_type string `json:"action_type"`
	MountPath   string `json:"mountPath`
	FsType      string `json:"kubernetes.io/fsType"`
}

// VolumeRequest is a structure for all properties of
// a volume request
type VolumeRequest struct {
	DockId       string `json:"dockId,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size"`
	AllowDetails bool   `json:"allowDetails"`

	ActionType string `json:"actionType"`
	Host       string `json:"host,omitempty"`
	Device     string `json:"device"`
	Attachment string `json:"attachment,omitempty"`
	MountDir   string `json:"mountDir"`
	FsType     string `json:"fsType"`
}

// VolumeResponse is a structure for all properties of
// an volume for a non detailed query
type VolumeResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Status      string              `json:"status"`
	Size        int64               `json:"size"`
	Volume_type string              `json:"volume_type"`
	Attachments []map[string]string `json:"attachments"`
}

// ShareDetailResponse is a structure for all properties of
// an share for a detailed query
type ShareDetailResponse struct {
	Links                    []map[string]string `json:"links"`
	AvailabilityZone         string              `json:"availability_zone,omitempty"`
	ShareNetworkId           string              `json:"share_network_id,omitempty"`
	ExportLocations          []string            `json:"export_locations"`
	ShareServerId            string              `json:"share_server_id,omitempty"`
	SnapshotId               string              `json:"snapshot_id,omitempty"`
	Id                       string              `json:"id,omitempty"`
	Size                     int                 `json:"size"`
	ShareType                string              `json:"share_type,omitempty"`
	ShareTypeName            string              `json:"share_type_name,omitempty"`
	ExportLocation           string              `json:"export_location,omitempty"`
	ConsistencyGroupId       string              `json:"consistency_group_id,omitempty"`
	ProjectId                string              `json:"project_id,omitempty"`
	Metadata                 map[string]string   `json:"metadata"`
	Status                   string              `json:"status,omitempty"`
	AccessRulesStatus        string              `json:"access_rules_status,omitempty"`
	Description              string              `json:"description,omitempty"`
	Host                     string              `json:"host,omitempty"`
	TaskState                string              `json:"task_state,omitempty"`
	IsPublic                 bool                `json:"is_public"`
	SnapshotSupport          bool                `json:"snapshot_support"`
	Name                     string              `json:"name,omitempty"`
	HasReplicas              bool                `json:"has_replicas"`
	ReplicationType          string              `json:"replication_type,omitempty"`
	ShareProto               string              `json:"share_proto,omitempty"`
	VolumeType               string              `json:"volume_type,omitempty"`
	SourceCgsnapshotMemberId string              `json:"source_cgsnapshot_member_id,omitempty"`
}

type FlexVolumePlugin interface {
	NewOptions() interface{}
	Init() Result
	Attach(opt interface{}) Result
	Detach(device string) Result
	Mount(mountDir string, device string, opt interface{}) Result
	Unmount(mountDir string) Result
}

func Succeed(a ...interface{}) Result {
	return Result{
		Status:  "Success",
		Message: fmt.Sprint(a...),
	}
}

func Fail(a ...interface{}) Result {
	return Result{
		Status:  "Failure",
		Message: fmt.Sprint(a...),
	}
}

func finish(result Result) {
	code := 1
	if result.Status == "Success" {
		code = 0
	}
	res, err := json.Marshal(result)
	if err != nil {
		fmt.Println("{\"status\":\"Failure\",\"message\":\"JSON error\"}")
	} else {
		fmt.Println(string(res))
	}
	os.Exit(code)
}

func RunPlugin(plugin FlexVolumePlugin) {
	if len(os.Args) < 2 {
		finish(Fail("Expected at least one argument"))
	}

	switch os.Args[1] {
	case "init":
		finish(plugin.Init())

	case "attach":
		if len(os.Args) != 3 {
			finish(Fail("attach expected exactly 3 arguments; got ", os.Args))
		}

		opt := plugin.NewOptions()
		if err := json.Unmarshal([]byte(os.Args[2]), opt); err != nil {
			finish(Fail("Could not parse options for attach:", err))
		}

		finish(plugin.Attach(opt))

	case "detach":
		if len(os.Args) != 3 {
			finish(Fail("detach expected exactly 3 arguments; got ", os.Args))
		}

		device := os.Args[2]
		finish(plugin.Detach(device))

	case "mount":
		if len(os.Args) != 5 {
			finish(Fail("mount expected exactly 5 argument; got ", os.Args))
		}

		mountDir := os.Args[2]
		device := os.Args[3]

		opt := plugin.NewOptions()
		if err := json.Unmarshal([]byte(os.Args[4]), opt); err != nil {
			finish(Fail("Could not parse options for mount; got ", os.Args[4]))
		}

		finish(plugin.Mount(mountDir, device, opt))

	case "unmount":
		if len(os.Args) != 3 {
			finish(Fail("mount expected exactly 5 argument; got ", os.Args))
		}

		mountDir := os.Args[2]

		finish(plugin.Unmount(mountDir))

	default:
		finish(Fail("Not sure what to do. Called with: ", os.Args))
	}
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

type DockNode struct {
	Id      string   `json:"id"`
	Drivers []string `json:"drivers"`
	Address string   `json:"address"`
}

func getDockNode() (*DockNode, error) {
	var nodePtr = &DockNode{}

	userJSON, err := ioutil.ReadFile("/etc/opensds/dock_node.json")
	if err != nil {
		log.Println("ReadFile json failed:", err)
		return nodePtr, err
	}

	if err = json.Unmarshal(userJSON, nodePtr); err != nil {
		log.Println("Unmarshal json failed:", err)
		return nodePtr, err
	}
	return nodePtr, nil
}

func GetDockId() (string, error) {
	dock, err := getDockNode()
	if err != nil {
		log.Println("Could not get dock routes:", err)
		return "", err
	}

	host, err := utils.GetHostIP()
	if err != nil {
		log.Println("Could not get dock host ip:", err)
		return "", err
	}

	if dock.Address == host {
		return dock.Id, nil
	} else {
		err = errors.New("Could not find dock service!")
		return "", err
	}
}
