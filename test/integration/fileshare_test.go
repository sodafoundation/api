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

// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/opensds/opensds/test/integration/utils"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestFileShare(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "FileShare Suite")
}

var (
	OPERATION_FAILED = "OPERATION_FAILED"
	file_shares      = []string{"hjk1", "", "İnanç Esasları", "#FileShare Code!$!test", "123tmp1", "abqwqwqwggg012345678910gggggggggggggghhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg",
		"mno2", "mno3", "mno4", "mno5", "mno6", "mno7", "mno8", "mno9", "mno10", "hjk1", "mno12", "mno13", "mno14", "mno15", "mno16", "mno17",//6-21
		"mno18", "mno19", "mno20", "mno21", "mno22", "$File$Test!1"}//22-27
	file_share_snapshots = []string{"snap1", "snap2", "snap3", "snap4", "snap5", "snap6", "snap7", "snap8", "snap9", "snap10", "snap11", "snap12", "snap13", "snap14", "snap15", "snap16", "snap17",
		"snap18", "snap19", "snap20", "snap21", "snap22", "snap23"}
)

func get_profile_id()string{
	res, err := http.Get("http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/profiles")
	if err != nil{
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		log.Fatalln(err)
	}

	filesharelstmap := []map[string]interface{}{}
	if err := json.Unmarshal(body, &filesharelstmap); err != nil {
		panic(err)
	}
	for _,k := range filesharelstmap{
		profilename := fmt.Sprintf("%v", k["name"])
		fmt.Println("+++++++++++++++++++++++++++++++++++++++",profilename)
		storagetype := fmt.Sprintf("%v",k["storageType"])
		if profilename == "mpolicy" && storagetype == "file"{
			id := fmt.Sprintf("%v", k["id"])
			return id
		}
	}
	return "None"
}

func get_all_file_share_snapshots()[]map[string]interface{}{
	res, err := http.Get("http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots")
	if err != nil{
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		log.Fatalln(err)
	}

	filesharesnaplstmap := []map[string]interface{}{}
	if err := json.Unmarshal(body, &filesharesnaplstmap); err != nil {
		panic(err)
	}
	return filesharesnaplstmap
}

func get_file_share_snapshots_Id_by_name(name string) string{
	filesharelstmap := get_all_file_share_snapshots()
	for _,k := range filesharelstmap{
		filesharename := fmt.Sprintf("%v", k["name"])
		status := fmt.Sprintf("%v", k["status"])
		if name == filesharename && status == "available"{
			id := fmt.Sprintf("%v", k["id"])
			return id
		}
	}
	return "None"
}

func get_all_file_shares()[]map[string]interface{}{
	res, err := http.Get("http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares")
	if err != nil{
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		log.Fatalln(err)
	}

	filesharelstmap := []map[string]interface{}{}
	if err := json.Unmarshal(body, &filesharelstmap); err != nil {
		panic(err)
	}
	return filesharelstmap
}

func get_file_share_Id_by_name(name string) string{
	filesharelstmap := get_all_file_shares()
	for _,k := range filesharelstmap{
		filesharename := fmt.Sprintf("%v", k["name"])
		status := fmt.Sprintf("%v", k["status"])
		if name == filesharename && status == "available"{
			id := fmt.Sprintf("%v", k["id"])
			return id
		}
	}
	return "None"
}

func hello(done chan bool) {
	fmt.Println("hello go routine is going to sleep")
	////gomega.Eventually(get_profile_id(),(80 * time.Second), (time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//profile_id := get_profile_id()
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>",profile_id)
	//var jsonStr= map[string]interface{}{"name": file_shares[0], "description": "This is just for TCFSIT01", "size": 1,"profileId":profile_id}
	//url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
	//methodName := "POST"
	//resp, err := utils.ConnectToHTTP(methodName, url, jsonStr)
	//fmt.Println(resp.Body)
	//gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//gomega.Expect(err).NotTo(gomega.HaveOccurred())

	var jsonStr= map[string]interface{}{"name": "mpolicy", "description": "This is just to test all file share test cases", "storageType":"file"}
	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/profiles"
	methodName := "POST"
	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr)

	fmt.Println("End of profile creation============================ ")//, body)
	gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	time.Sleep(60 * time.Second)
	fmt.Println("hello go routine awake and going to write to done")
	done <- true
}

var _ = ginkgo.Describe("FileShare Testing", func() {
	ginkgo.Context("Create FileShare Scenarios", func() {

		var profile_id string
		ginkgo.BeforeEach(func() {
			//profile_id = get_profile_id()
			//fmt.Println(profile_id)

		})

		ginkgo.It("TC_FS_IT_01: Create profile for file share", func() {
			done := make(chan bool)
			fmt.Println("Main going to call hello go goroutine")
			go hello(done)
			<-done
			profile_id = get_profile_id()
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>",profile_id)
			fId := "v1beta/file/shares/e93b4c0934da416eb9c8d120c5d04d96/"+profile_id
			ret := utils.GetValueByKeyFromDB(fId)
			gomega.Expect(ret).ShouldNot(gomega.Equal(OPERATION_FAILED))
			fmt.Println("Main received data")
		})
		ginkgo.It("TC_FS_IT_02: Get profile for file share", func() {
			//var jsonStr= map[string]interface{}{"name": "policy", "description": "This is just to test all file share test cases", "storageType":"file"}
			profile_id = get_profile_id()
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>",profile_id)
			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/profiles/"+profile_id
			methodName := "GET"
			resp, err := utils.ConnectToHTTP(methodName, url, nil)
			fmt.Println(resp.Body)
			gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
		//ginkgo.It("TC_FS_IT_01: Create fileshare with name input", func() {

		//
		//})
		//ginkgo.It("TC_FS_IT_02: Create fileshare with empty file share name", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[1], "description": "This is just for TCFSIT02", "size": 2}
		//	url := "http://localhost:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, _ := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//})
		//ginkgo.It("TC_FS_IT_03: Create file share name with other encoding characters(except utf-8)", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[2], "description": "This is just for TCFSIT03", "size": 2}
		//	url := "http://localhost:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_04: Create file share name having special characters", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[3], "description": "This is just for TCFSIT04", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, _ := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//})
		//ginkgo.It("TC_FS_IT_05: Create file share name starts with numbers", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[4], "description": "This is just for TCFSIT05", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_06: Create file share name length more than 255 characters", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[5], "description": "This is just for TCFSIT06", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_08: Create file share description with empty string", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[6], "description": "", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_09: Create file share with description having special characters", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[7], "description": "#FileShare Code!$!test", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_10: Create file share with description with more than 255 characters", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[8], "description": "abqwqwqwggg012345678910ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggghhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg",
		//		"size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_11: Create file share without required parameters like fileshare name, size", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[9], "description": " "}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_12: Create file share with size with -ve number", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[10], "description": "This is just for TCFSIT12", "size": -2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_13: Create file share with size with +ve number", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[11], "description": "This is just for TCFSIT13", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_14: Create file share with size with more than total capacity", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[12], "description": "This is just for TCFSIT14", "size": 50}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_15: Create file share Size with 0", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[13], "description": "This is just for TCFSIT15", "size": 0}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_16: Create file share by specifying the File Share Id", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[14], "description": "This is just for TCFSIT16", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		////profile_id := get_profile_id()
		////fmt.Println(profile_id)
		//ginkgo.It("TC_FS_IT_17: Create file share by specifying the Profile id", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[15], "description": "This is just for TCFSIT17", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_18: Create file share by without specifying Profile Id", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[16], "description": "This is just for TCFSIT18", "size": 2}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_19: Create file share by specifying wrong profile Id", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[17], "description": "This is just for TCFSIT19", "size": 2, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_20: Create file share by specifying Availability zone name", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[18], "description": "This is just for TCFSIT20", "size": 2, "availabillityZone": "default"}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
		//ginkgo.It("TC_FS_IT_21: Create file share by specifying wrong Availability zone name", func() {
		//	var jsonStr2= map[string]interface{}{"name": file_shares[19], "description": "This is just for TCFSIT21", "size": 2, "availabillityZone": "default1"}
		//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
		//	methodName := "POST"
		//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})
	})
	//ginkgo.Context("Create FileShare Snapshots Scenarios", func() {
	//	var fileshareid string
	//	filesharename := file_shares[15]
	//	fmt.Println(filesharename)
	//	ginkgo.BeforeEach(func() {
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//	})
	//	ginkgo.It("TC_FS_IT_22: Create fileshare snapshot with name input", func() {
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second), (time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr = map[string]interface{}{"name": file_share_snapshots[0], "description": "This is just for TCFSIT22"}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_23: Create fileshare snapshot with empty file share snapshot name", func() {
	//		filesharename = file_shares[16]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": " ", "description": "This is just for TCFSIT23", "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, _ := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	})
	//	ginkgo.It("TC_FS_IT_24: Create file share snapshot name with other encoding characters(except utf-8)", func() {
	//		filesharename = file_shares[18]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_shares[2], "description": "This is just for TCFSIT24", "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_25: Create file share snapshot name having special characters", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_shares[3], "description": "This is just for TCFSIT25", "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, _ := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//})
	//	ginkgo.It("TC_FS_IT_26: Create file share snapshot name starts with numbers", func() {
	//      filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_shares[4], "description": "This is just for TCFSIT26", "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_27: Create file share snapshot name length more than 255 characters", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_shares[5], "description": "This is just for TCFSIT27", " "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_28: Create file share snapshot description with empty string", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[6], "description": "This is just for TCFSIT28",  "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_29: Create file share snapshot with description having special characters", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[7], "description": "#FileShare Code!$!test",  "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_30: Create file share snapshot with description with more than 255 characters", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[8], "description": "abqwqwqwggg012345678910ggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggghhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg",
	//		 "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_31: Create file share snapshot without required parameters like fileshare snapshot name", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"description": "This is just for TCFSIT31",  "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_32: Create file share snapshot by specifying the File Share snapshot Id", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[14], "description": "This is just for TCFSIT32",  "fileshareId":fileshareid, "id":"04f46c79-9498-4d29-a4d3-7996386723cc"}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	profile_id := get_profile_id()
	//	fmt.Println(profile_id)
	//	ginkgo.It("TC_FS_IT_33: Create file share snapshot by specifying the Profile id", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[15], "description": "This is just for TCFSIT33",  "fileshareId":fileshareid, "profileId": profile_id}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_34: Create file share snapshot by without specifying Profile Id", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[16], "description": "This is just for TCFSIT34",  "fileshareId":fileshareid}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_35: Create file share snapshot by specifying wrong profile Id", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[17], "description": "This is just for TCFSIT35",  "fileshareId":fileshareid, "profileId": "df40af1a-17b5-48e5-899f-fa098b0bd5da"}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_36: Create file share snapshot by specifying Availability zone name", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[18], "description": "This is just for TCFSIT36",  "fileshareId":fileshareid, "availabillityZone": "default"}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_37: Create file share snapshot by specifying wrong Availability zone name", func() {
	//		filesharename = file_shares[11]
	//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid = get_file_share_Id_by_name(filesharename)
	//		var jsonStr2= map[string]interface{}{"name": file_share_snapshots[19], "description": "This is just for TCFSIT37",  "fileshareId":fileshareid, "availabillityZone": "default1"}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "POST"
	//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//})
	//ginkgo.Context("Create FileShare Access Permission Scenarios", func() {
		//	ginkgo.It("TC_FS_IT_38: Create file share access permission by providing valid IP", func() {
		//		filesharename = file_shares[11]
		//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
		//		fileshareid = get_file_share_Id_by_name(filesharename)
		//		var jsonStr2= map[string]interface{}{"type":"ip","accessCapability":[]string{"read", "write"}, "accessTo":"10.32.105.191", "description": "This is just for TCFSIT38",  "shareId":fileshareid}
		//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/acls"
		//		methodName := "POST"
		//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//	})
		//	ginkgo.It("TC_FS_IT_39: Create file share access permission by providing Invalid IP", func() {
		//		filesharename = file_shares[11]
		//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
		//		fileshareid = get_file_share_Id_by_name(filesharename)
		//		var jsonStr2= map[string]interface{}{"type":"ip","accessCapability":[]string{"read", "write"}, "accessTo":"10.32.105", "description": "This is just for TCFSIT39",  "shareId":fileshareid}
		//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/acls"
		//		methodName := "POST"
		//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//	ginkgo.It("TC_FS_IT_40: Create the access to a valid IP segment, EX: 10.0.0.0/10", func() {
		//		filesharename = file_shares[11]
		//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
		//		fileshareid = get_file_share_Id_by_name(filesharename)
		//		var jsonStr2= map[string]interface{}{"type":"ip","accessCapability":[]string{"read", "write"}, "accessTo":"10.32.105.10/10", "description": "This is just for TCFSIT40",  "shareId":fileshareid}
		//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/acls"
		//		methodName := "POST"
		//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//	ginkgo.It("TC_FS_IT_41: Create the access to a Invalid IP segment, EX: 10.0.0.0/10", func() {
		//		filesharename = file_shares[11]
		//		gomega.Eventually(get_file_share_Id_by_name(filesharename),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
		//		fileshareid = get_file_share_Id_by_name(filesharename)
		//		var jsonStr2= map[string]interface{}{"type":"ip","accessCapability":[]string{"read", "write"}, "accessTo":"10.32.105.10/10.10", "description": "This is just for TCFSIT41",  "shareId":fileshareid}
		//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/acls"
		//		methodName := "POST"
		//		resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
		//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(400))
		//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
//	})
	//})
	//ginkgo.Context("Get FileShare Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_42: Get file share by specifying FileShareID", func() {
	//		gomega.Eventually(get_file_share_Id_by_name("hjk1"),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		fileshareid := get_file_share_Id_by_name("hjk1")
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/"+fileshareid
	//		methodName := "GET"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//	ginkgo.It("TC_FS_IT_43: Get file share by specifying wrong/non-existing FileShareID", func() {
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/"+fileshareid+"1"
	//		methodName := "GET"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(404))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	})
	//ginkgo.Context("Get FileShare Snapshots Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_44: Get file share snapshot by specifying SnapID", func() {
	//		gomega.Eventually(get_file_share_snapshots_Id_by_name("snap1"),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
	//		snapid = get_file_share_snapshots_Id_by_name("snap1")
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots/"+snapid
	//		methodName := "GET"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//  snapid = get_file_share_snapshots_Id_by_name("snap1")
	//	ginkgo.It("TC_FS_IT_45: Get file share by specifying wrong/non-existing FileShareID", func() {
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots/"+snapid+"1"
	//		methodName := "GET"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(404))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//})
	//ginkgo.Context("Get FileShare Access Permissions Scenarios", func() {
		//	ginkgo.It("TC_FS_IT_46: Get file share acl by specifying aclID", func() {
		//		gomega.Eventually(get_file_share_acls_Id_by_name("snap1"),(60 * time.Second)).ShouldNot(gomega.BeEquivalentTo("None"))
		//		aclid = get_file_share_snapshots_Id_by_name("snap1")
		//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/acls/"+aclid
		//		methodName := "GET"
		//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
		//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
		//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//	})
		//  aclid = get_file_share_acls_Id_by_name("snap1")
		//	ginkgo.It("TC_FS_IT_47: Get file share acl by specifying wrong/non-existing aclid", func() {
		//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/acls/"+aclid+"1"
		//		methodName := "GET"
		//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
		//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(404))
		//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//	})
	//})
	//ginkgo.Context("List FileShare Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_47: List all file shares", func() {
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares"
	//		methodName := "GET"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//})
	//ginkgo.Context("List FileShare Snapshot Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_47: List all file share snapshots", func() {
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots"
	//		methodName := "GET"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//})
	//ginkgo.Context("List FileShare Acccess Permissions Scenarios", func() {
	//	ginkgo.It("TC_FS_IT_47: List all file share acls", func() {
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/acls"
	//		methodName := "GET"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(200))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	})
	//})
	//ginkgo.Context("Update FileShare Scenarios", func() {
	//	//ginkgo.It("TC_FS_IT_22: Update file share name with empty string", func() {
	//	//	var jsonStr2 = map[string]interface{}{"name": file_shares[1], "description": "This is test for case TC_FS_IT_22"}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/add61d3c-5248-4dec-8d0d-9e84fd6fbbb6"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	//ginkgo.It("TC_FS_IT_23: Update file share name with special character string", func() {
	//	//	var jsonStr2 = map[string]interface{}{"name":file_shares[27], "description": "This is test for case TC_FS_IT_23"}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96 /file/shares/e3885ace-23d4-4a52-a47f-0efc87b82821"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	//ginkgo.It("TC_FS_IT_24: Update file share name length greater than 255 characters", func() {
	//	//	var jsonStr2 = map[string]interface{}{"name": file_shares[5]+"ffffffgfgfgfgffffffffffgffffffffffffffgfgfggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg", "description": "This is test for case TC_FS_IT_24"}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/ea5b77d1-9cb1-4ccc-b5a5-7126e1f6fffb"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	//ginkgo.It("TC_FS_IT_25: Update file share Description with empty string", func() {
	//	//	var jsonStr2 = map[string]interface{}{"description": " "}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/de3c2b70-ae83-4a90-b5be-c792d22c7930"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	//ginkgo.It("TC_FS_IT_26: Update file share Description with special characters", func() {
	//	//	var jsonStr2 = map[string]interface{}{"description": " "}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/de3c2b70-ae83-4a90-b5be-c792d22c7930"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	//ginkgo.It("TC_FS_IT_27: Update file share Description length more than 255 characters", func() {
	//	//	var jsonStr2 = map[string]interface{}{"description": " "}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/de3c2b70-ae83-4a90-b5be-c792d22c7930"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	//ginkgo.It("TC_FS_IT_28: Update file share Description with non utf-8 code characters", func() {
	//	//	var jsonStr2 = map[string]interface{}{"description": " "}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/de3c2b70-ae83-4a90-b5be-c792d22c7930"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	//ginkgo.It("TC_FS_IT_29: Update file share with wrong File Share Id", func() {
	//	//	var jsonStr2 = map[string]interface{}{"description": " "}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/de3c2b70-ae83-4a90-b5be-c792d22c7930"
	//	//	methodName := "PUT"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, jsonStr2)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//})
	//ginkgo.Context("Update FileShare Snapshots Scenarios", func() {
	//
	//})
	//ginkgo.Context("Update FileShare Access Permissions Scenarios", func() {
	//
	//})
	//ginkgo.Context("Delete FileShare Access Permissions Scenarios", func() {
	//
	//})
	//ginkgo.Context("Delete FileShare Snapshots Scenarios", func() {
	//	filesharelstmap := get_all_file_share_snapshots()
	//	for _,k := range filesharelstmap{
	//		id := fmt.Sprintf("%v", k["id"])
	//		ginkgo.It("TC_FS_IT_32: Delete all file shares", func() {
	//			//var jsonStr2 = map[string]interface{}{"name": file_shares[19], "description": " ", "size": 2, "availabillityZone": "default1"}
	//			url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/snapshots/"+id
	//			methodName := "DELETE"
	//			resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//			gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//			gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//		})
	//	}
	//})
	//ginkgo.Context("Delete FileShare Scenarios", func() {
	//	fmt.Println(file_shares[11])
	//	fileshareid := get_file_share_Id_by_name("xyz11")
	//	fmt.Println(fileshareid)
	//	//ginkgo.It("TC_FS_IT_31: Delete file share by specifying FileShare ID", func() {
	//	//	//var jsonStr2 = map[string]interface{}{"name": file_shares[19], "description": " ", "size": 2, "availabillityZone": "default1"}
	//	//	url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/"+fileshareid
	//	//	methodName := "DELETE"
	//	//	resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//	//	gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
	//	//	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//	//})
	//	ginkgo.It("TC_FS_IT_31: Delete file share by specifying wrong/non existing FileShare ID", func() {
	//		//var jsonStr2 = map[string]interface{}{"name": file_shares[19], "description": " ", "size": 2, "availabillityZone": "default1"}
	//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/"+fileshareid+"2"
	//		methodName := "DELETE"
	//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
	//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(404))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//})

		//filesharelstmap := get_all_file_shares()
		//for _,k := range filesharelstmap{
		//	id := fmt.Sprintf("%v", k["id"])
		//	ginkgo.It("TC_FS_IT_32: Delete all file shares", func() {
		//		//var jsonStr2 = map[string]interface{}{"name": file_shares[19], "description": " ", "size": 2, "availabillityZone": "default1"}
		//		url := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/"+id
		//		methodName := "DELETE"
		//		resp, err := utils.ConnectToHTTP(methodName, url, nil)
		//		gomega.Expect(resp.StatusCode).Should(gomega.Equal(202))
		//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//	})
		//}
	//})
//})
