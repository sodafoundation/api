package api_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/opensds/opensds/api/client/operations"
	"github.com/opensds/opensds/api/models"

	httptransport "github.com/go-openapi/runtime/client"
	apiclient "github.com/opensds/opensds/api/client"
)

func ContainsAll(secret interface{}) (bool, error) {
	value := reflect.ValueOf(secret)
	for i := 0; i < value.NumField(); i++ {
		elem := fmt.Sprintf("%v", value.Field(i))

		switch elem {
		// case "":
		//  	err := errors.New("Some string type properties not supported!")
		//  	return false, err
		case "0":
			err := errors.New("Some integer type properties not supported!")
			return false, err
		default:
			continue
		}
	}
	return true, nil
}

func TestListVersions(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	versions, err := client.Operations.ListVersions(operations.NewListVersionsParams())
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*versions.Payload.Versions[0])
	if !result {
		t.Fatal(err)
	}
}

func TestGetVersionv1(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	version, err := client.Operations.GetVersionv1(operations.NewGetVersionv1Params())
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*version.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestCreateVolume(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	cvp := operations.NewCreateVolumeParams()
	cvp.ResourceType = "cinder"
	cvp.VolumeRequest = &models.VolumeRequest{
		Name: "newvol",
		Size: 1,
	}

	volume, err := client.Operations.CreateVolume(cvp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*volume.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestGetVolume(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	gvp := operations.NewGetVolumeParams()
	gvp.ResourceType = "cinder"
	gvp.ID = "myvol"

	volume, err := client.Operations.GetVolume(gvp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*volume.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestListVolumes(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	lvp := operations.NewListVolumesParams()
	lvp.ResourceType = "cinder"

	volumes, err := client.Operations.ListVolumes(lvp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*volumes.Payload[0])
	if !result {
		t.Fatal(err)
	}
}

func TestDeleteVolume(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	dvp := operations.NewDeleteVolumeParams()
	dvp.ResourceType = "cinder"
	dvp.ID = "myvol"

	resp, err := client.Operations.DeleteVolume(dvp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*resp.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestAttachVolume(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	ovp := operations.NewOperateVolumeParams()
	ovp.ResourceType = "cinder"
	ovp.ID = "myvol"
	ovp.VolumeRequest = &models.VolumeRequest{
		ActionType: "attach",
		Host:       "localhost",
		Device:     "/dev/vdc",
	}

	resp, err := client.Operations.OperateVolume(ovp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*resp.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestDetachVolume(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	ovp := operations.NewOperateVolumeParams()
	ovp.ResourceType = "cinder"
	ovp.ID = "myvol"
	ovp.VolumeRequest = &models.VolumeRequest{
		ActionType: "detach",
		Attachment: "/dev/vdc",
	}

	resp, err := client.Operations.OperateVolume(ovp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*resp.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestMountVolume(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	ovp := operations.NewOperateVolumeParams()
	ovp.VolumeRequest = &models.VolumeRequest{
		ActionType: "mount",
		MountDir:   "/mnt",
		Device:     "/dev/vdc",
		FsType:     "ext4",
	}

	resp, err := client.Operations.OperateVolume(ovp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*resp.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestUnmountVolume(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	ovp := operations.NewOperateVolumeParams()
	ovp.VolumeRequest = &models.VolumeRequest{
		ActionType: "unmount",
		MountDir:   "/mnt",
	}

	resp, err := client.Operations.OperateVolume(ovp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*resp.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestCreateShare(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	csp := operations.NewCreateShareParams()
	csp.ResourceType = "manila"
	csp.ShareRequest = &models.ShareRequest{
		Name:       "newshr",
		ShareType:  "g-nfs",
		ShareProto: "NFS",
		Size:       1,
	}

	share, err := client.Operations.CreateShare(csp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*share.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestGetShare(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	gsp := operations.NewGetShareParams()
	gsp.ResourceType = "manila"
	gsp.ID = "myshr"

	share, err := client.Operations.GetShare(gsp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*share.Payload)
	if !result {
		t.Fatal(err)
	}
}

func TestListShares(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	lsp := operations.NewListSharesParams()
	lsp.ResourceType = "manila"

	shares, err := client.Operations.ListShares(lsp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*shares.Payload[0])
	if !result {
		t.Fatal(err)
	}
}

func TestDeleteShare(t *testing.T) {
	transport := httptransport.New("127.0.0.1:8080", "", nil)

	client := apiclient.New(transport, strfmt.Default)

	dsp := operations.NewDeleteShareParams()
	dsp.ResourceType = "manila"
	dsp.ID = "myshr"

	resp, err := client.Operations.DeleteShare(dsp)
	if err != nil {
		t.Fatal(err)
	}

	result, err := ContainsAll(*resp.Payload)
	if !result {
		t.Fatal(err)
	}
}
