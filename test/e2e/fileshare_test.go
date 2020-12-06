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

// +build e2e

package e2e

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/opensds/opensds/test/e2e/utils"
)

func TestFileShare(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "FileShare Suite")
}

func clean_resource(name string){
	var fileshare_id,url string
	ginkgo.It("TC_FS_IT_32: Delete fileshare by specifying fileshare id", func() {
		gomega.Eventually(func() string {
			fileshare_id = utils.Get_file_share_Id_by_name(name)
			return fileshare_id
		}).ShouldNot(gomega.Equal("None"))
		url = "file/shares/" + fileshare_id
		resp, err := utils.DELETE_method_call(url)
		utils.Resp_processing(resp.StatusCode, err)
	})
	ginkgo.It("TC_FS_IT_02: Get profile for file share", func() {
		gomega.Eventually(func() string {
			fileshare_id = utils.Get_file_share_Id_by_name(name)
			return fileshare_id
		}).Should(gomega.Equal("None"))
		url = "file/shares/" + fileshare_id
		resp, err := utils.GET_method_call(url)
		utils.Resource_Not_found(resp.StatusCode, err)
	})
}

func wait_for_available(name string){
	var fileshare_id, url string
	ginkgo.It("TC_FS_IT_02: Get profile for file share", func() {
		gomega.Eventually(func() string {
			fileshare_id = utils.Get_file_share_Id_by_name(name)
			return fileshare_id
		}).ShouldNot(gomega.Equal("None"))
		url = "file/shares/" + fileshare_id
		resp, err := utils.GET_method_call(url)
		utils.Resp_ok(resp.StatusCode, err)
	})
}

var _ = ginkgo.Describe("FileShare Testing", func() {

	var  url, profile_id string //fileshare_id, snap_id, acl_id string
	var jsonStr map[string]interface{}
	//big_name := "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"
	ginkgo.Context("Profile for FileShare Scenarios", func() {
		//url = "profiles"
		////ginkgo.It("TC_FS_IT_01: Create profile for file share", func() {
		////	jsonStr = map[string]interface{}{"name": "temp_profile", "description": "This is for TC_FS_IT_01", "storageType": "file"}
		////	resp, err := utils.POST_method_call(jsonStr, url)
		////	utils.Resp_ok(resp.StatusCode, err)
		////})
		//ginkgo.It("TC_FS_IT_02: Get profile for file share", func() {
		//	profile_id = utils.Get_profile_id_by_name("default_file")
		//	url = "profiles/" + profile_id
		//	resp, err := utils.GET_method_call(url)
		//	utils.Resp_ok(resp.StatusCode, err)
		//})
	})
	ginkgo.Context("Create FileShare Scenarios", func() {
		url = "file/shares"
		ginkgo.It("TC_FS_IT_03: Create fileshare with name 'filesharetest1' and size 2", func() {
			jsonStr = map[string]interface{}{"name": "filesharetest45", "description": "This is for TCFSIT03", "size": 2}
			resp, err := utils.POST_method_call(jsonStr, url)
			utils.Resp_processing(resp.StatusCode, err)
		})
		wait_for_available("filesharetest45")
		clean_resource("filesharetest45")
		ginkgo.It("TC_FS_IT_04: Create fileshare name with empty string", func() {
			jsonStr = map[string]interface{}{"name": "", "description": "This is for TCFSIT04", "size": 1}
			resp, err := utils.POST_method_call(jsonStr, url)
			utils.Not_allowed_operation(resp.StatusCode, err)
		})
		//ginkgo.It("TC_FS_IT_05: Create file share name with non utf-8 character codes)", func() {
		//	jsonStr = map[string]interface{}{"name": "İnanç Esasları", "description": "This is for TCFSIT05", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("İnanç Esasları")
		//clean_resource("İnanç Esasları")
		//ginkgo.It("TC_FS_IT_06: Create file share name having special characters)", func() {
		//	jsonStr = map[string]interface{}{"name": "#Share !$!test", "description": "This is for TCFSIT06", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("#Share !$!test")
		//clean_resource("#Share !$!test")
		//ginkgo.It("TC_FS_IT_07: Create file share name starts with numbers)", func() {
		//	jsonStr = map[string]interface{}{"name": "23tmp", "description": "This is for TCFSIT07", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("23tmp")
		//clean_resource("23tmp")
		//ginkgo.It("TC_FS_IT_08: Create file share name length more than 255 characters)", func() {
		//	jsonStr = map[string]interface{}{"name": big_name, "description": "This is for TCFSIT08", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available(big_name)
		//clean_resource(big_name)
		//ginkgo.It("TC_FS_IT_09: Create file share description with empty string)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest2", "description": "", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Resp_processing(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest2")
		//clean_resource("filesharetest2")
		//ginkgo.It("TC_FS_IT_10: Create file share with description having special characters)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest3", "description": "#FileShare Code!$!test TC_FS_IT_10", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest3")
		//clean_resource("filesharetest3")
		//ginkgo.It("TC_FS_IT_11: Create file share name length more than 255 characters)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest4", "decription": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest4")
		//clean_resource("filesharetest4")
		//ginkgo.It("TC_FS_IT_12: Create file share without required parameters like fileshare name)", func() {
		//	jsonStr = map[string]interface{}{"description": "This is for TCFSIT12", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//ginkgo.It("TC_FS_IT_13: Create file share with size with -ve number)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest6", "description": "This is for TCFSIT13", "size": -1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest6")
		//clean_resource("filesharetest6")
		//ginkgo.It("TC_FS_IT_14: Create file share without required parameters like fileshare size)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest7", "description": "This is for TCFSIT14"}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest7")
		//clean_resource("filesharetest7")
		//ginkgo.It("TC_FS_IT_15: Create file share with size with +ve number)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest8", "description": "This is for TCFSIT15", "size": 2}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Resp_processing(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest8")
		//clean_resource("filesharetest8")
		//ginkgo.It("TC_FS_IT_16: Create file share with size with +ve number)", func() {
		//	cap := utils.Get_capacity_of_pool()
		//	free_capacity, err := strconv.Atoi(cap)
		//	jsonStr = map[string]interface{}{"name": "filesharetest9", "description": "This is for TCFSIT16", "size": free_capacity + 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest9")
		//clean_resource("filesharetest9")
		//ginkgo.It("TC_FS_IT_17: Create file share with size 0)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest10", "description": "This is for TCFSIT17", "size": 0}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest10")
		//clean_resource("filesharetest10")
		//ginkgo.It("TC_FS_IT_18: Create file share by specifying fileshare id)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest11", "description": "This is for TCFSIT18", "size": 1, "id": "5ce2fead-d404-47a6-b6b1-e2b069129240"}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest11")
		//clean_resource("filesharetest11")
		//ginkgo.It("TC_FS_IT_19: Create file share by specifying profile id)", func() {
		//	profile_id = utils.Get_profile_id_by_name("temp_profile")
		//	jsonStr = map[string]interface{}{"name": "filesharetest12", "description": "This is for TCFSIT19", "size": 1, "profileId": profile_id}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Resp_processing(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest12")
		//clean_resource("filesharetest12")
		//ginkgo.It("TC_FS_IT_20: Create file share by without specifying profile id)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest13", "description": "This is for TCFSIT20", "size": 1}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Resp_processing(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest13")
		//clean_resource("filesharetest13")
		//ginkgo.It("TC_FS_IT_21: Create file share by specifying wrong profile id)", func() {
		//	profile_id = utils.Get_profile_id_by_name("temp_profile")
		//	jsonStr = map[string]interface{}{"name": "filesharetest14", "description": "This is for TCFSIT21", "size": 1, "profileId": profile_id + "2"}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest14")
		//clean_resource("filesharetest14")
		//ginkgo.It("TC_FS_IT_22: Create file share by specifying availability zone)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest15", "description": "This is for TCFSIT22", "size": 1, "availabilityZone": "default"}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Resp_processing(resp.StatusCode, err)
		//})
		//wait_for_available("filesharetest15")
		//clean_resource("filesharetest15")
		//ginkgo.It("TC_FS_IT_23: Create file share by specifying wrong availability zone)", func() {
		//	jsonStr = map[string]interface{}{"name": "filesharetest16", "description": "This is for TCFSIT23", "size": 1, "availabilityZone": "ndefault"}
		//	resp, err := utils.POST_method_call(jsonStr, url)
		//	utils.Not_allowed_operation(resp.StatusCode, err)
		//})
	})
	//ginkgo.Context("Prepare for update fileshare scenarios", func() {
	//	url = "file/shares"
	//	//ginkgo.It("Dep_TC_FS_IT_24: Create file share for TC_FS_IT_24)", func() {
	//	//	jsonStr = map[string]interface{}{"name": "testshare1", "description": "This is for DepTCFSIT24", "size": 1}
	//	//	resp, err := utils.POST_method_call(jsonStr, url)
	//	//	utils.Resp_processing(resp.StatusCode, err)
	//	//})
	//	ginkgo.It("Dep_TC_FS_IT_25: Create file share for TC_FS_IT_25)", func() {
	//		jsonStr = map[string]interface{}{"name": "testshare2", "description": "This is for TCFSIT25", "size": 1}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_26: Create file share for TC_FS_IT_26)", func() {
	//		jsonStr = map[string]interface{}{"name": "testshare3", "description": "This is for TCFSIT26", "size": 1}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_27: Create file share for TC_FS_IT_27)", func() {
	//		jsonStr = map[string]interface{}{"name": "testshare4", "description": "This is for TCFSIT27", "size": 1}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	//ginkgo.It("Dep_TC_FS_IT_28: Create file share for TC_FS_IT_28)", func() {
	//	//	jsonStr = map[string]interface{}{"name": "testshare5", "description": "This is for TCFSIT28", "size": 1}
	//	//	resp, err := utils.POST_method_call(jsonStr, url)
	//	//	utils.Resp_processing(resp.StatusCode, err)
	//	//})
	//	//ginkgo.It("Dep_TC_FS_IT_29: Create file share for TC_FS_IT_29)", func() {
	//	//	jsonStr = map[string]interface{}{"name": "testshare6", "description": "This is for TCFSIT29", "size": 1}
	//	//	resp, err := utils.POST_method_call(jsonStr, url)
	//	//	utils.Resp_processing(resp.StatusCode, err)
	//	//})
	//	//ginkgo.It("Dep_TC_FS_IT_30: Create file share for TC_FS_IT_30)", func() {
	//	//	jsonStr = map[string]interface{}{"name": "testshare7", "description": "This is for TCFSIT30", "size": 1}
	//	//	resp, err := utils.POST_method_call(jsonStr, url)
	//	//	utils.Resp_processing(resp.StatusCode, err)
	//	//})
	//	//ginkgo.It("Dep_TC_FS_IT_31: Create file share for TC_FS_IT_31)", func() {
	//	//	jsonStr = map[string]interface{}{"name": "testshare8", "description": "This is for TCFSIT31", "size": 1}
	//	//	resp, err := utils.POST_method_call(jsonStr, url)
	//	//	utils.Resp_processing(resp.StatusCode, err)
	//	//})
	//})
	//ginkgo.Context("Update FileShare Scenarios", func() {
	//
	//	ginkgo.It("TC_FS_IT_24: Update fileshare with description empty/'' ", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("filesharetest13")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"description": " "}
	//		fileshare_id = utils.Get_file_share_Id_by_name("filesharetest1")
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//
	//	ginkgo.It("TC_FS_IT_25: Update fileshare with name empty/'' ", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("filesharetest8")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": ""}
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_26: Update fileshare name within 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("filesharetest2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "testshare3_renamed"}
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_27: Update fileshare name with special characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("filesharetest1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "#FileShare Code!$!test"}
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_28: Update fileshare name length greater than 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("testshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"}
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_29: Update fileshare description with special characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("testshare3")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"description": "#FileShare Code!$!test"}
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_30: Update fileshare description length greater than 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("testshare4")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"description": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"}
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_31: Update fileshare name with wrong resource id", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("testshare8")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "wrong_resource_update", "id": "fb5357a9-40cc-45a0-bdc2-14759e9c7d7a"}
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Delete FileShare Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_32: Delete fileshare by specifying fileshare id", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("filesharetest13")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.DELETE_method_call(url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_33: Delete fileshare by specifying wrong fileshare id", func() {
	//		url = "file/shares/" + "fb5357a9-40cc-45a0-bdc2-14759e9c7d7a"
	//		resp, err := utils.DELETE_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Get FileShare Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_34: Get fileshare by specifying fileshare id", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("testshare3_renamed")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		url = "file/shares/" + fileshare_id
	//		resp, err := utils.GET_method_call(url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_35: Get fileshare by specifying wrong fileshare id", func() {
	//		url = "file/shares/" + "901ce29a-24d6-4d52-be3e-02a1ac5a7b46"
	//		resp, err := utils.GET_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("List all fileShares", func() {
	//	ginkgo.It("TC_FS_IT_36: List all file shares", func() {
	//		url = "file/shares"
	//		resp, err := utils.GET_method_call(url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Prepare file share for snapshots", func() {
	//	url = "file/shares"
	//	ginkgo.It("Dep_TC_FS_IT_38: Create file share for TC_FS_IT_38)", func() {
	//		jsonStr = map[string]interface{}{"name": "snaptestshare1", "description": "This is for TCFSIT38", "size": 1}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_50: Create file share for TC_FS_IT_50)", func() {
	//		jsonStr = map[string]interface{}{"name": "snaptestshare2", "description": "This is for TCFSIT50", "size": 1}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Create file shares snapshot scenarios", func() {
	//	url = "file/snapshots"
	//	ginkgo.It("TC_FS_IT_38:Create file shares with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snap1", "description": "This is for TCFSIT38", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_39:Create file share snapshot name with empty string name", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "", "description": "This is for TCFSIT39", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_40:Create file share snapshot name with non utf-8 encoding characters)", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "İnanç Esasları", "description": "This is for TCFSIT40", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_41:Create file share snapshot name having special characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "Snapshot!$! test", "description": "This is for TCFSIT41", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_42:Create file share snapshot name starts with numbers", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "23snap", "description": "This is for TCFSIT42", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_43:Create file share snapshot name length more than 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", "description": "This is for TCFSIT43", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_44:Create file share snapshot description with empty string", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snap2", "description": "", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_45:Create file share snapshot with description having special characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snap3", "description": "FileShare Code!$!test2", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_46:Create file share snapshot name length more than 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snap4", "description": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_47:Create file share snapshot without required parameters like fileshare Id", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snap5", "description": "This is for TCFSIT47"}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_48:Create file share snapshot by specifying the profile id", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		gomega.Eventually(func() string {
	//			profile_id = utils.Get_profile_id_by_name("temp_profile")
	//			return profile_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snap5", "description": "This is for TCFSIT48", "fileshareId": fileshare_id, "profileId": profile_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_49:Create file share snapshot by specifying wrong/non-existing profile Id", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snap6", "description": "This is for TCFSIT22", "fileshareId": fileshare_id, "profileId": "5bc5d37e-da03-48be-b8b7-624d6c117cfb"}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Prepare for update file shares snapshot scenarios", func() {
	//	url = "file/snapshots"
	//	ginkgo.It("Dep_TC_FS_IT_50:Create file share snapshot with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest1", "description": "This is for TCFSIT50", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_51:Create file share snapshot with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest2", "description": "This is for TCFSIT51", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_52:Create file share snapshot with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest3", "description": "This is for TCFSIT52", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_53:Create file share snapshot with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest4", "description": "This is for TCFSIT53", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_54:Create file share snapshot with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest5", "description": "This is for TCFSIT54", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_55:Create file share snapshot with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest6", "description": "This is for TCFSIT55", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_56:Create file share snapshot with name, description and fileshareId", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("snaptestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest7", "description": "This is for TCFSIT56", "fileshareId": fileshare_id}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Update FileShare Snapshot Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_50: Update fileshare snapshot with description empty/'' ", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest1")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"description": " "}
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//
	//	ginkgo.It("TC_FS_IT_51: Update fileshare snapshot with name empty/'' ", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest2")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": " "}
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_52: Update fileshare snapshot name within 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest3")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "snaptest3_renamed"}
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_53: Update fileshare snapshot name with special characters", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest4")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "#FileShare Snap!$!test"}
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_54: Update fileshare snapshot name length greater than 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest5")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"name": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"}
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_55: Update fileshare snapshot description with special characters", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest6")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"description": "#FileShare Code!$!test"}
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_56: Update fileshare snapshot description length greater than 255 characters", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest7")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"description": "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhyutyuytuytututututututututututututututututututyuhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"}
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_57: Update fileshare snapshot name with wrong resource id", func() {
	//		jsonStr = map[string]interface{}{"name": "wrong_resource_update"}
	//		url = "file/snapshots/fb5357a9-40cc-45a0-bdc2-14759e9c7d7a"
	//		resp, err := utils.PUT_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Delete FileShare Snapshot Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_58: Delete fileshare snapshot by specifying fileshare snapshot id", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest1")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		url = "file/snapshots/" + snap_id
	//		resp, err := utils.DELETE_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_59: Delete fileshare snapshot by specifying wrong fileshare snapshot id", func() {
	//		url = "file/snapshots/fb5357a9-40cc-45a0-bdc2-14759e9c7d7a"
	//		resp, err := utils.DELETE_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Get FileShare Snapshot Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_60: Get fileshare snapshot by specifying fileshare snapshot id", func() {
	//		gomega.Eventually(func() string {
	//			snap_id = utils.Get_snapshot_Id_by_name("snaptest3_renamed")
	//			return snap_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		url = "file/snapshots" + snap_id
	//		resp, err := utils.GET_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_61: Get fileshare snapshot by specifying wrong fileshare snapshot id", func() {
	//		url = "file/snapshots/901ce29a-24d6-4d52-be3e-02a1ac5a7b46"
	//		resp, err := utils.GET_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("List all fileshare snapshots", func() {
	//	ginkgo.It("TC_FS_IT_62: List all fileshare snapshots", func() {
	//		url = "file/snapshots"
	//		resp, err := utils.GET_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Prepare for create file share permission set scenarios", func() {
	//	ginkgo.It("Dep_TC_FS_IT_63: Create file share for TC_FS_IT_63)", func() {
	//		jsonStr = map[string]interface{}{"name": "acltestshare1", "description": "This is for TCFSIT63", "size": 1}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_64: Create file share for TC_FS_IT_24)", func() {
	//		jsonStr = map[string]interface{}{"name": "acltestshare2", "description": "This is for TCFSIT24", "size": 1}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("Dep_TC_FS_IT_69:Create the access to a valid IP for TC_FS_IT_69", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("acltestshare2")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"fileshareId": fileshare_id, "description": "file share acl test", "type": "ip", "accessTo": "10.32.104.4", "accessCapability": []string{"Read", "Write"}}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Create file shares acls scenarios", func() {
	//	url = "file/acls"
	//	ginkgo.It("TC_FS_IT_63:Create the access to a valid IP", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("acltestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"fileshareId": fileshare_id, "description": "file share acl test", "type": "ip", "accessTo": "10.32.104.2", "accessCapability": []string{"Read", "Write"}}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_64:Create the access to a invalid IP", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("acltestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"fileshareId": fileshare_id, "description": "file share acl test", "type": "ip", "accessTo": "10.32.104", "accessCapability": []string{"Read", "Write"}}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_65:Create the access to a valid IP segment, EX: 10.0.0.0/10", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("acltestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"fileshareId": fileshare_id, "description": "file share acl test", "type": "ip", "accessTo": "10.32.104.8/20", "accessCapability": []string{"Read", "Write"}}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_66:Create the access to a invalid IP segment, EX: 10.0.0.0/10.2", func() {
	//		gomega.Eventually(func() string {
	//			fileshare_id = utils.Get_file_share_Id_by_name("acltestshare1")
	//			return fileshare_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		jsonStr = map[string]interface{}{"fileshareId": fileshare_id, "description": "file share acl test", "type": "ip", "accessTo": "10.32.104.8/20.2", "accessCapability": []string{"Read", "Write"}}
	//		resp, err := utils.POST_method_call(jsonStr, url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Delete FileShare Acl Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_67: Delete fileshare acl by specifying fileshare acl id", func() {
	//		gomega.Eventually(func() string {
	//			acl_id = utils.Get_acl_Id_by_ip("10.32.104.2")
	//			return acl_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		url = "file/acls/" + acl_id
	//		resp, err := utils.DELETE_method_call(url)
	//		utils.Resp_processing(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_68: Delete fileshare acl by specifying wrong fileshare acl id", func() {
	//		url = "file/acls/" + "fb5357a9-40cc-45a0-bdc2-14759e9c7d7a"
	//		resp, err := utils.DELETE_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("Get FileShare Acl Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_69: Get fileshare acl by specifying fileshare acl id", func() {
	//		gomega.Eventually(func() string {
	//			acl_id = utils.Get_acl_Id_by_ip("10.32.104.4")
	//			return acl_id
	//		}).ShouldNot(gomega.Equal("None"))
	//		url = "file/acls" + acl_id
	//		resp, err := utils.GET_method_call(url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//	ginkgo.It("TC_FS_IT_70: Get fileshare acl by specifying wrong fileshare acl id", func() {
	//		url = "file/acls/901ce29a-24d6-4d52-be3e-02a1ac5a7b46"
	//		resp, err := utils.GET_method_call(url)
	//		utils.Not_allowed_operation(resp.StatusCode, err)
	//	})
	//})
	//ginkgo.Context("List all fileshare acls", func() {
	//	ginkgo.It("TC_FS_IT_71: List all fileshare acls", func() {
	//		url = "file/acls"
	//		resp, err := utils.GET_method_call(url)
	//		utils.Resp_ok(resp.StatusCode, err)
	//	})
	//})
})
