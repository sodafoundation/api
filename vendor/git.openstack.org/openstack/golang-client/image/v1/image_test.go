// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
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

// image.go
package v1

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"git.openstack.org/openstack/golang-client/openstack"
	"git.openstack.org/openstack/golang-client/testUtil"
	"git.openstack.org/openstack/golang-client/util"
)

var tokn = "eaaafd18-0fed-4b3a-81b4-663c99ec1cbb"

func TestListImages(t *testing.T) {
	anon := func(imageService *Service) {
		images, err := imageService.Images()
		if err != nil {
			t.Error(err)
		}

		if len(images) != 3 {
			t.Error(errors.New("Incorrect number of images found"))
		}
		expectedImage := Response{
			Name:            "Ubuntu Server 14.04.1 LTS (amd64 20140927) - Partner Image",
			ContainerFormat: "bare",
			DiskFormat:      "qcow2",
			CheckSum:        "6798a7d67ff0b241b6fe165798914d86",
			ID:              "bec3cab5-4722-40b9-a78a-3489218e22fe",
			Size:            255525376}
		// Verify first one matches expected values
		testUtil.Equals(t, expectedImage, images[0])
	}

	testImageServiceAction(t, "images", sampleImagesData, anon)
}

func TestListImageDetails(t *testing.T) {
	anon := func(imageService *Service) {
		images, err := imageService.ImagesDetail()
		if err != nil {
			t.Error(err)
		}

		if len(images) != 2 {
			t.Error(errors.New("Incorrect number of images found"))
		}
		createdAt, _ := util.NewDateTime(`"2014-09-29T14:44:31"`)
		updatedAt, _ := util.NewDateTime(`"2014-09-29T15:33:37"`)
		owner := "10014302369510"
		virtualSize := int64(2525125)
		expectedImageDetail := DetailResponse{
			Status:          "active",
			Name:            "Ubuntu Server 12.04.5 LTS (amd64 20140927) - Partner Image",
			Deleted:         false,
			ContainerFormat: "bare",
			CreatedAt:       createdAt,
			DiskFormat:      "qcow2",
			UpdatedAt:       updatedAt,
			MinDisk:         8,
			Protected:       false,
			ID:              "8ca068c5-6fde-4701-bab8-322b3e7c8d81",
			MinRAM:          0,
			CheckSum:        "de1831ea85702599a27e7e63a9a444c3",
			Owner:           &owner,
			IsPublic:        true,
			DeletedAt:       nil,
			Properties: map[string]string{
				"com.ubuntu.cloud__1__milestone":    "release",
				"com.hp__1__os_distro":              "com.ubuntu",
				"description":                       "Ubuntu Server 12.04.5 LTS (amd64 20140927) for HP Public Cloud. Ubuntu Server is the world's most popular Linux for cloud environments. Updates and patches for Ubuntu 12.04.5 LTS will be available until 2017-04-26. Ubuntu Server is the perfect platform for all workloads from web applications to NoSQL databases and Hadoop. More information regarding Ubuntu Cloud is available from http://www.ubuntu.com/cloud and instructions for using Juju to deploy workloads are available from http://juju.ubuntu.com EULA: http://www.ubuntu.com/about/about-ubuntu/licensing Privacy Policy: http://www.ubuntu.com/privacy-policy",
				"com.ubuntu.cloud__1__suite":        "precise",
				"com.ubuntu.cloud__1__serial":       "20140927",
				"com.hp__1__bootable_volume":        "True",
				"com.hp__1__vendor":                 "Canonical",
				"com.hp__1__image_lifecycle":        "active",
				"com.hp__1__image_type":             "disk",
				"os_version":                        "12.04",
				"architecture":                      "x86_64",
				"os_type":                           "linux-ext4",
				"com.ubuntu.cloud__1__stream":       "server",
				"com.ubuntu.cloud__1__official":     "True",
				"com.ubuntu.cloud__1__published_at": "2014-09-29T15:33:36"},
			Size:        261423616,
			VirtualSize: &virtualSize}
		testUtil.Equals(t, expectedImageDetail, images[0])
	}

	testImageServiceAction(t, "images/detail", sampleImageDetailsData, anon)
}

func TestNameFilterUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?name=CentOS+deprecated",
		QueryParameters{Name: "CentOS deprecated"})
}

func TestStatusUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?status=active",
		QueryParameters{Status: "active"})
}

func TestMinMaxSizeUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?size_max=5300014&size_min=100158",
		QueryParameters{MinSize: 100158, MaxSize: 5300014})
}

func TestMarkerLimitUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?limit=20&marker=bec3cab5-4722-40b9-a78a-3489218e22fe",
		QueryParameters{Marker: "bec3cab5-4722-40b9-a78a-3489218e22fe", Limit: 20})
}

func TestContainerFormatFilterUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?container_format=bare",
		QueryParameters{ContainerFormat: "bare"})
}

func TestSortKeySortUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?sort_key=id",
		QueryParameters{SortKey: "id"})
}

func TestSortDirSortUrlProduced(t *testing.T) {
	testImageQueryParameter(t, "images?sort_dir=asc",
		QueryParameters{SortDirection: Asc})
}

func testImageQueryParameter(t *testing.T, uriEndsWith string, queryParameters QueryParameters) {
	anon := func(imageService *Service) {
		_, _ = imageService.QueryImages(&queryParameters)
	}

	testImageServiceAction(t, uriEndsWith, sampleImagesData, anon)
}

func testImageServiceAction(t *testing.T, uriEndsWith string, testData string, imageServiceAction func(*Service)) {
	anon := func(req *http.Request) {
		reqURL := req.URL.String()
		if !strings.HasSuffix(reqURL, uriEndsWith) {
			t.Error(errors.New("Incorrect url created, expected:" + uriEndsWith + " at the end, actual url:" + reqURL))
		}
	}
	apiServer := testUtil.CreateGetJSONTestRequestServer(t, tokn, testData, anon)
	defer apiServer.Close()

	auth := openstack.AuthToken{
		Access: openstack.AccessType{
			Token: openstack.Token{
				ID: tokn,
			},
		},
	}
	sess, _ := openstack.NewSession(http.DefaultClient, auth, nil)
	imageService := Service{
		Session: *sess,
		URL:     apiServer.URL,
	}
	imageServiceAction(&imageService)
}

var sampleImagesData = `{
   "images":[
      {
         "name":"Ubuntu Server 14.04.1 LTS (amd64 20140927) - Partner Image",
         "container_format":"bare",
         "disk_format":"qcow2",
         "checksum":"6798a7d67ff0b241b6fe165798914d86",
         "id":"bec3cab5-4722-40b9-a78a-3489218e22fe",
         "size":255525376
      },
      {
         "name":"Ubuntu Server 12.04.5 LTS (amd64 20140927) - Partner Image",
         "container_format":"bare",
         "disk_format":"qcow2",
         "checksum":"de1831ea85702599a27e7e63a9a444c3",
         "id":"8ca068c5-6fde-4701-bab8-322b3e7c8d81",
         "size":261423616
      },
      {
         "name":"HP_LR-PC_Load_Generator_12-02_Windows-2008R2x64",
         "container_format":"bare",
         "disk_format":"qcow2",
         "checksum":"052d70c2b4d4988a8816197381e9083a",
         "id":"12b9c19b-8823-4f40-9531-0f05fb0933f2",
         "size":14012055552
      }
   ]
}`

var sampleImageDetailsData = `{
   "images":[
      {
         "status":"active",
         "name":"Ubuntu Server 12.04.5 LTS (amd64 20140927) - Partner Image",
         "deleted":false,
         "container_format":"bare",
         "created_at":"2014-09-29T14:44:31",
         "disk_format":"qcow2",
         "updated_at":"2014-09-29T15:33:37",
         "min_disk":8,
         "protected":false,
         "id":"8ca068c5-6fde-4701-bab8-322b3e7c8d81",
         "min_ram":0,
         "checksum":"de1831ea85702599a27e7e63a9a444c3",
         "owner":"10014302369510",
         "is_public":true,
         "deleted_at":null,
         "properties":{
            "com.ubuntu.cloud__1__milestone":"release",
            "com.hp__1__os_distro":"com.ubuntu",
            "description":"Ubuntu Server 12.04.5 LTS (amd64 20140927) for HP Public Cloud. Ubuntu Server is the world's most popular Linux for cloud environments. Updates and patches for Ubuntu 12.04.5 LTS will be available until 2017-04-26. Ubuntu Server is the perfect platform for all workloads from web applications to NoSQL databases and Hadoop. More information regarding Ubuntu Cloud is available from http://www.ubuntu.com/cloud and instructions for using Juju to deploy workloads are available from http://juju.ubuntu.com EULA: http://www.ubuntu.com/about/about-ubuntu/licensing Privacy Policy: http://www.ubuntu.com/privacy-policy",
            "com.ubuntu.cloud__1__suite":"precise",
            "com.ubuntu.cloud__1__serial":"20140927",
            "com.hp__1__bootable_volume":"True",
            "com.hp__1__vendor":"Canonical",
            "com.hp__1__image_lifecycle":"active",
            "com.hp__1__image_type":"disk",
            "os_version":"12.04",
            "architecture":"x86_64",
            "os_type":"linux-ext4",
            "com.ubuntu.cloud__1__stream":"server",
            "com.ubuntu.cloud__1__official":"True",
            "com.ubuntu.cloud__1__published_at":"2014-09-29T15:33:36"
         },
         "size":261423616,
		 "virtual_size":2525125
      },
      {
         "status":"active",
         "name":"Windows Server 2008 Enterprise SP2 x64 Volume License 20140415 (b)",
         "deleted":false,
         "container_format":"bare",
         "created_at":"2014-04-25T19:53:24",
         "disk_format":"qcow2",
         "updated_at":"2014-04-25T19:57:11",
         "min_disk":30,
         "protected":true,
         "id":"1294610e-fdc4-579b-829b-d0c9f5c0a612",
         "min_ram":0,
         "checksum":"37208aa6d49929f12132235c5f834f2d",
         "owner":null,
         "is_public":true,
         "deleted_at":null,
         "properties":{
            "hp_image_license":"1002",
            "com.hp__1__os_distro":"com.microsoft.server",
            "com.hp__1__image_lifecycle":"active",
            "com.hp__1__image_type":"disk",
            "architecture":"x86_64",
            "com.hp__1__license_os":"1002",
            "com.hp__1__bootable_volume":"true"
         },
         "size":6932856832,
		"virtual_size":null
      }
   ]
}`
