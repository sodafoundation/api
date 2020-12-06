package utils

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega"
	"io/ioutil"
	"log"
	"net/http"
)

func Resp_ok(resp int, err error){
	gomega.Expect(resp).Should(gomega.Equal(200))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	fmt.Println("Resp. Ok....")
}

func Resp_processing(resp int, err error){
	gomega.Expect(resp).Should(gomega.Equal(202))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	fmt.Println("Processing operations....")
}

func Not_allowed_operation(resp int, err error){
	gomega.Expect(resp).Should(gomega.Equal(400))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	fmt.Println("Not allowed operation....")
}

func Resource_Not_found(resp int, err error){
	gomega.Expect(resp).Should(gomega.Equal(404))
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	fmt.Println("No resource....")
}

func POST_method_call(jsonStr map[string]interface{}, url string) (*http.Response, error) {
	inpurl := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/" + url
	resp, err := ConnectToHTTP("POST", inpurl, jsonStr)
	return resp, err
}

func PUT_method_call(jsonStr map[string]interface{}, url string) (*http.Response, error) {
	inpurl := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/" + url
	resp, err := ConnectToHTTP("PUT", inpurl, jsonStr)
	return resp, err
}

func GET_method_call(url string) (*http.Response, error) {
	inpurl := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/" + url
	resp, err := ConnectToHTTP("GET", inpurl, nil)
	return resp, err
}

func DELETE_method_call(url string) (*http.Response, error) {
	inpurl := "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/" + url
	resp, err := ConnectToHTTP("DELETE", inpurl, nil)
	return resp, err
}

func Non_read_onclose_get_method_call(url string)([]map[string]interface{}){
	res, err := http.Get("http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/"+url)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	lstmap := []map[string]interface{}{}
	if err := json.Unmarshal(body, &lstmap); err != nil {
		panic(err)
	}
	return lstmap
}

func Get_profile_id_by_name(name string) string {
	profileslstmap := Non_read_onclose_get_method_call("profiles")
	for _, k := range profileslstmap {
		profilename := fmt.Sprintf("%v", k["name"])
		storagetype := fmt.Sprintf("%v", k["storageType"])
		if profilename == name && storagetype == "file" {
			id := fmt.Sprintf("%v", k["id"])
			return id
		}
	}
	return "None"
}

func Get_file_share_Id_by_name(name string) string {
	filesharelstmap := Non_read_onclose_get_method_call("file/shares")
	for _, k := range filesharelstmap {
		filesharename := fmt.Sprintf("%v", k["name"])
		status := fmt.Sprintf("%v", k["status"])
		if name == filesharename && status == "available" {
			id := fmt.Sprintf("%v", k["id"])
			return id
		}
	}
	return "None"
}

func Get_snapshot_Id_by_name(name string) string {
	filesharelstmap := Non_read_onclose_get_method_call("file/snapshots")
	for _, k := range filesharelstmap {
		snapname := fmt.Sprintf("%v", k["name"])
		status := fmt.Sprintf("%v", k["status"])
		if name == snapname && status == "available" {
			id := fmt.Sprintf("%v", k["id"])
			return id
		}
	}
	return "None"
}

func Get_acl_Id_by_ip(name string) string {
	filesharelstmap := Non_read_onclose_get_method_call("file/acls")
	for _, k := range filesharelstmap {
		ipaddr := fmt.Sprintf("%v", k["ip"])
		status := fmt.Sprintf("%v", k["status"])
		if name == ipaddr && status == "available" {
			id := fmt.Sprintf("%v", k["id"])
			return id
		}
	}
	return "None"
}

func Get_all_file_shares()([]map[string]interface{}){
	filesharelstmap := Non_read_onclose_get_method_call("file/shares")
	return filesharelstmap
}


func Get_capacity_of_pool()string{
	poolslstmap := Non_read_onclose_get_method_call("pools")
	for _, k := range poolslstmap {
		poolname := fmt.Sprintf("%v", k["name"])
		pool_freecapacity := fmt.Sprintf("%v", k["freeCapacity"])
		if poolname == "opensds-files-default"{
			return pool_freecapacity
		}
	}
	return "None"
}

//func Validateresponsecode(respcode int, err error) {
	//	switch respcode{
	//	case 200:
	//		gomega.Expect(respcode).Should(gomega.Equal(200))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//		fmt.Println("Resp. Ok....")
	//	case 202:
	//		gomega.Expect(respcode).Should(gomega.Equal(202))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//		fmt.Println("Processing req....")
	//	case 500:
	//		gomega.Expect(respcode).Should(gomega.Equal(500))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//		fmt.Println("Internal server error....")
	//	case 404:
	//		gomega.Expect(respcode).Should(gomega.Equal(404))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//		fmt.Println("Resource not found....")
	//	case 400:
	//		gomega.Expect(respcode).Should(gomega.Equal(400))
	//		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	//		fmt.Println("Bad req....")
	//	default:
	//		fmt.Println("Test case failed")
	//	}
	//}