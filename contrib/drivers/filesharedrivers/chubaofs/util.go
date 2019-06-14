package chubaofs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	log "github.com/golang/glog"
)

type clusterInfoResponseData struct {
	LeaderAddr string `json:"LeaderAddr"`
}

type clusterInfoResponse struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data *clusterInfoResponseData `json:"data"`
}

type createVolumeResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func getClusterInfo(host string) (string, error) {
	url := "http://" + host + "/admin/getCluster"
	log.Infof("chubaofs: GetClusterInfo(%v)", url)

	httpResp, err := http.Get(url)
	if err != nil {
		log.Errorf("chubaofs: failed to GetClusterInfo, url(%v) err(%v)", url, err)
		return "", err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Errorf("chubaofs: failed to read response, url(%v) err(%v)", url, err)
		return "", err
	}

	resp := &clusterInfoResponse{}
	if err = json.Unmarshal(body, resp); err != nil {
		errmsg := fmt.Sprintf("chubaofs: getClusterInf failed to unmarshal, bodyLen(%d) err(%v)", len(body), err)
		log.Error(errmsg)
		return "", errors.New(errmsg)
	}

	log.Infof("chubaofs: GetClusterInfo, url(%v), resp(%v)", url, resp)

	if resp.Code != 0 {
		errmsg := fmt.Sprintf("chubaofs: GetClusterInfo code NOK, url(%v) code(%v) msg(%v)", url, resp.Code, resp.Msg)
		log.Error(errmsg)
		return "", errors.New(errmsg)
	}

	if resp.Data == nil {
		errmsg := fmt.Sprintf("chubaofs: GetClusterInfo nil data, url(%v) msg(%v)", url, resp.Code, resp.Msg)
		log.Error(errmsg)
		return "", errors.New(errmsg)
	}

	return resp.Data.LeaderAddr, nil
}

func createVolume(leader string, name string, size int64) error {
	// Round up to Giga bytes
	sizeInGB := (size + (1 << 30) - 1) >> 30
	url := fmt.Sprintf("http://%s/admin/createVol?name=%s&capacity=%v&owner=chubaofs", leader, name, sizeInGB)
	log.Infof("chubaofs: CreateVolume url(%v)", url)

	httpResp, err := http.Get(url)
	if err != nil {
		errmsg := fmt.Sprintf("chubaofs: CreateVolume failed, url(%v) err(%v)", url, err)
		log.Error(errmsg)
		return errors.New(errmsg)
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		errmsg := fmt.Sprintf("chubaofs: CreateVolume failed to unmarshal, bodyLen(%d) err(%v)", len(body), err)
		log.Error(errmsg)
		return errors.New(errmsg)
	}

	resp := &createVolumeResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		errmsg := fmt.Sprintf("chubaofs: GetClusterInfo code NOK, url(%v) code(%v) msg(%v)", url, resp.Code, resp.Msg)
		log.Error(errmsg)
		return errors.New(errmsg)
	}

	log.Infof("chubaofs: CreateVolume url(%v) resp(%v)", url, resp)

	if resp.Code != 0 {
		if resp.Code == 1 {
			log.Warning("chubaofs: CreateVolume volume exist, url(%v) msg(%v)", url, resp.Msg)
		} else {
			errmsg := fmt.Sprintf("chubaofs: CreateVolume failed, url(%v) code(%v) msg(%v)", url, resp.Code, resp.Msg)
			log.Error(errmsg)
			return errors.New(errmsg)
		}
	}

	return nil
}

func generateFile(filePath string, data []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(data)
}
