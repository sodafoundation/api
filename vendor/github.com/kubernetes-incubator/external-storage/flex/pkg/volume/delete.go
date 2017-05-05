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
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"k8s.io/client-go/pkg/api/v1"
)

func (p *flexProvisioner) Delete(volume *v1.PersistentVolume) error {
	var (
		backenddriver, profilename, volid string
	)
	glog.Infof("Delete called for volume:", volume.Name)

	provisioned, err := p.provisioned(volume)
	if err != nil {
		return fmt.Errorf("error determining if this provisioner was the one to provision volume %q: %v", volume.Name, err)
	}
	if !provisioned {
		strerr := fmt.Sprintf("this provisioner id %s didn't provision volume %q and so can't delete it; id %s did & can", p.identity, volume.Name, volume.Annotations[annProvisionerId])
		return &controller.IgnoredError{strerr}
	}

	flexVolumeOptions := volume.Spec.PersistentVolumeSource.FlexVolume.Options
	storageTags := make(map[string]string)

	for key, value := range flexVolumeOptions {
		switch key {
		case backendDriver:
			backenddriver = value
		case profileName:
			profilename = value
		case volId:
			volid = value
		default:
			storageTags[key] = value
		}
	}

	url := URL_PREFIX + "/api/v1/volumes/" + volid
	vr := &VolumeRequest{
		Schema: &VolumeOperationSchema{},
		Profile: &StorageProfile{
			Name:          profilename,
			BackendDriver: backenddriver,
			StorageTags:   storageTags,
		},
	}

	// Start DELETE request to delete volume
	req := httplib.Delete(url).SetTimeout(100*time.Second, 50*time.Second)
	req.JSONBody(vr)

	resp, err := req.Response()
	if err != nil {
		glog.Errorf("Failed to delete volume %s, error: %s", volume.Name, err.Error())
		return err
	}

	err = CheckHTTPResponseStatusCode(resp)
	if err != nil {
		glog.Errorf("Failed to delete volume %s, error: %s", volume.Name, err.Error())
		return err
	}

	rbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to delete volume %s, error: %s", volume.Name, err.Error())
		return err
	}

	var volumeDeleteResponse = &VolumeDeleteResponse{}
	if err = json.Unmarshal(rbody, volumeDeleteResponse); err != nil {
		glog.Errorf("Failed to delete volume %s, error: %s", volume.Name, err.Error())
		return err
	}
	if !strings.Contains(string(rbody), "Success") {
		glog.Errorf("Failed to delete volume %s, error: %s", volume.Name, err.Error())
		return err
	}
	return nil
}

func (p *flexProvisioner) provisioned(volume *v1.PersistentVolume) (bool, error) {
	provisionerId, ok := volume.Annotations[annProvisionerId]
	if !ok {
		return false, fmt.Errorf("PV doesn't have an annotation %s", annProvisionerId)
	}

	return provisionerId == string(p.identity), nil
}

func getBlockAndId(volume *v1.PersistentVolume, annBlock, annId string) (string, uint16, error) {
	block, ok := volume.Annotations[annBlock]
	if !ok {
		return "", 0, fmt.Errorf("PV doesnot have an annotation with key %s", annBlock)
	}

	idStr, ok := volume.Annotations[annId]
	if !ok {
		return "", 0, fmt.Errorf("PV doesn't have an annotation %s", annId)
	}

	id, _ := strconv.ParseUint(idStr, 10, 16)
	return block, uint16(id), nil
}
