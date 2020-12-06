package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"testing"
)

func TestClean_mount_paths(t *testing.T){
	ids := ClientPathsvalidation()
	for _,strout := range ids{
		client := &http.Client{}
		fmt.Println("http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/"+strout)
		// Create request
		req, err := http.NewRequest("DELETE", "http://127.0.0.1:50040/v1beta/e93b4c0934da416eb9c8d120c5d04d96/file/shares/"+strout, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		fmt.Println("response Status : ", resp.Status)
	}
}

func ClientPathsvalidation()[]string{
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

	if len(filesharelstmap) == 0{
		fmt.Println("empty file share paths")
	}
	
	df := exec.Command("df")
	grep := exec.Command("grep", "-w", "/mnt")
	cut := exec.Command("awk", "{print $6}")
	result, _, err := Pipline(df, grep, cut)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	out := string(result)
	result1 := strings.Split(out, "\n")
	var ids []string
	for _,k := range filesharelstmap{
		exportlocations := fmt.Sprintf("%v", k["exportLocations"])
		lstlocations := strings.Split(exportlocations, ":")
		mountpaths := TrimSuffix(lstlocations[1], "]")
		id := fmt.Sprintf("%v", k["id"])
		strout := strings.Replace(mountpaths, "var", "mnt", -1)
		bool1 := contains(result1, strout)
		if !bool1 {
			fmt.Sprintf("this path doesn't exists %s", strout)
		}else{
			ids = append(ids,id)
		}
	}
	return ids
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


func Pipline(cmds ...*exec.Cmd) ([]byte, []byte, error) {
	// credits to https://studygolang.com/articles/5387

	// At least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	var output bytes.Buffer
	var stderr bytes.Buffer
	var err error
	maxindex := len(cmds) - 1
	cmds[maxindex].Stdout = &output
	cmds[maxindex].Stderr = &stderr

	for i, cmd := range cmds[:maxindex] {
		if i == maxindex {
			break
		}

		cmds[i+1].Stdin, err = cmd.StdoutPipe()
		if err != nil {
			return nil, nil, err
		}
	}

	// Start each command
	for _, cmd := range cmds {
		err := cmd.Start()
		if err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		err := cmd.Wait()
		if err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	return output.Bytes(), stderr.Bytes(), nil
}

