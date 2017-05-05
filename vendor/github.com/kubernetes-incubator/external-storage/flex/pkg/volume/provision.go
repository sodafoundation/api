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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	// Name of the file where an s3fsProvisioner will store its identity
	identityFile = "flex-provisioner.identity"

	// osdslet service entry
	URL_PREFIX string = "http://192.168.0.9:50040"

	// are we allowed to set this? else make up our own
	annCreatedBy = "kubernetes.io/createdby"
	createdBy    = "flex-dynamic-provisioner"

	// A PV annotation for the identity of the s3fsProvisioner that provisioned it
	annProvisionerId = "Provisioner_Id"

	// volume id
	volId = "volumeId"

	//backend driver type, cinder etc
	backendDriver = "backendDriver"

	// file system type, ext4 etc
	fsType = "fsType"

	volumeType           = "volumeType"
	profileName          = "profileName"
	snapshotDeletePolicy = "snapshotDeletePolicy"
)

func NewFlexProvisioner(client kubernetes.Interface, execCommand string) controller.Provisioner {
	return newFlexProvisionerInternal(client, execCommand)
}

func newFlexProvisionerInternal(client kubernetes.Interface, execCommand string) *flexProvisioner {
	var identity types.UID

	provisioner := &flexProvisioner{
		client:      client,
		execCommand: execCommand,
		identity:    identity,
	}

	return provisioner
}

type flexProvisioner struct {
	client      kubernetes.Interface
	execCommand string
	identity    types.UID
}

var _ controller.Provisioner = &flexProvisioner{}

// Provision creates a volume i.e. the storage asset and returns a PV object for
// the volume.
func (p *flexProvisioner) Provision(options controller.VolumeOptions) (*v1.PersistentVolume, error) {
	var (
		backenddriver, fstype, volumetype, profilename string
	)

	fstype = options.PVC.ObjectMeta.Annotations[fsType]
	storageTags := make(map[string]string)
	flexVolumeOptions := make(map[string]string)

	for key, value := range options.Parameters {
		switch key {
		case backendDriver:
			flexVolumeOptions[key] = value
			backenddriver = value
		case volumeType:
			volumetype = value
		case profileName:
			flexVolumeOptions[key] = value
			profilename = value
		default:
			flexVolumeOptions[key] = value
			storageTags[key] = value
		}
	}

	volid, err := p.createVolume(options, backenddriver, volumetype, profilename, storageTags)
	if err != nil {
		return nil, err
	}

	glog.Infof("volume id = %s", volid)
	flexVolumeOptions[volId] = volid

	annotations := make(map[string]string)
	annotations[annCreatedBy] = createdBy

	annotations[annProvisionerId] = string(p.identity)

	/*
		This PV won't work since there's nothing backing it.  the flex script
		is in flex/flex/flex  (that many layers are required for the flex volume plugin)
	*/
	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        options.PVName,
			Labels:      map[string]string{},
			Annotations: annotations,
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{

				FlexVolume: &v1.FlexVolumeSource{
					Driver:   "opensds.io/opensds",
					Options:  flexVolumeOptions,
					FSType:   fstype,
					ReadOnly: false,
				},
			},
		},
	}

	return pv, nil
}

func (p *flexProvisioner) createVolume(
	volumeOptions controller.VolumeOptions,
	backenddriver string,
	volumetype string,
	profilename string,
	storageTags map[string]string,
) (string, error) {
	var storageInfo string = fmt.Sprint(volumeOptions.PVC.Spec.Resources.Requests["storage"])
	storageInfoStr := strings.Fields(storageInfo)
	sizeStr := storageInfoStr[3]

	positive, value, num, denom, suf, err := parseQuantityString(sizeStr)
	glog.Infof("Parse string: value = %s, denom = %s, suf = %s", value, denom, suf)

	if !positive {
		err = errors.New("Volume size cannot be negative!")
		glog.Errorf("Failed to create volume %s, error: %s", volumeOptions, err.Error())
		return "", err
	}

	volSize, err := strconv.Atoi(num)
	if err != nil {
		glog.Errorf("Failed to convert str to int, err: %s", volumeOptions, err.Error())
		return "", err
	}

	url := URL_PREFIX + "/api/v1/volumes"
	vr := &VolumeRequest{
		Schema: &VolumeOperationSchema{
			Name:       volumeOptions.PVName,
			VolumeType: volumetype,
			Size:       int32(volSize),
		},
		Profile: &StorageProfile{
			Name:          profilename,
			BackendDriver: backenddriver,
			StorageTags:   storageTags,
		},
	}

	// Start POST request to create volume
	req := httplib.Post(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		glog.Errorf("Failed to create volume %s, error: %s", volumeOptions, err.Error())
		return "", err
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		glog.Errorf("Failed to create volume %s, error: %s", volumeOptions, err.Error())
		return "", err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to create volume %s, error: %s", volumeOptions, err.Error())
		return "", err
	}

	var volumeResponse = &VolumeResponse{}
	if err = json.Unmarshal(rbody, volumeResponse); err != nil {
		glog.Errorf("Failed to create volume %s, error: %s", volumeOptions, err.Error())
		return "", err
	}
	return volumeResponse.ID, nil
}
