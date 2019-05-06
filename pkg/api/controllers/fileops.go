// Copyright (c) 2019 OpenSDS Authors All Rights Reserved.
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

/*
This module implements a entry into the OpenSDS northbound service to upload and download configuration files

*/

package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/model"
)

// prometheus constants
var PrometheusConfHome string
var PrometheusUrl string
var PrometheusConfFile string

// alert manager constants
var AlertmgrConfHome string
var AlertmgrUrl string
var AlertmgrConfFile string

var GrafanaConfHome string
var GrafanaRestartCmd string
var GrafanaConfFile string

var ReloadPath string
var BackupExtension string

func init() {

	// TODO Prakash read these from conf and save to these variables
	ReloadPath = "/-/reload"
	BackupExtension = ".bak"

	PrometheusConfHome = "/etc/prometheus/"
	PrometheusUrl = "http://localhost:9090"
	PrometheusConfFile = "prometheus.yml"

	AlertmgrConfHome = "/root/alertmanager-0.16.2.linux-amd64/"
	AlertmgrUrl = "http://localhost:9093"
	AlertmgrConfFile = "alertmanager.yml"

	GrafanaConfHome = "/etc/grafana/"
	GrafanaRestartCmd = "grafana-server"
	GrafanaConfFile = "grafana.ini"
}

func NewFileOpsPortal() *FileOpsPortal {
	return &FileOpsPortal{
		CtrClient: client.NewClient(),
	}
}

type FileOpsPortal struct {
	BasePortal

	CtrClient client.Client
}

func (f *FileOpsPortal) UploadConfFile() {

	m, _ := f.GetParameters()
	confType := m["conftype"][0]

	switch confType {
	case "prometheus":
		DoUpload(f, PrometheusConfHome, PrometheusUrl, ReloadPath, true)
	case "alertmanager":
		DoUpload(f, AlertmgrConfHome, AlertmgrUrl, ReloadPath, true)
	case "grafana":
		// for grafana, there is no reload endpoint to call
		DoUpload(f, GrafanaConfHome, "", "", false)
		// to reload the configuration, run the reload command for grafana
		cmd := exec.Command("systemctl", "restart", GrafanaRestartCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Fatalf("restart grafana failed with %s\n", err)
		}
		return

	}
}

func DoUpload(fileOpsPortal *FileOpsPortal, confHome string, url string, reloadPath string, toCallReloadEndpoint bool) {

	// get the uploaded file
	f, h, _ := fileOpsPortal.GetFile("conf_file")

	// get the path to save the configuration
	path := confHome + h.Filename

	// close the incoming file
	fCloseErr := f.Close()
	if fCloseErr != nil {
		log.Errorf("error closing uploaded file %s", h.Filename)
		fileOpsPortal.ErrorHandle(model.ErrorInternalServer, fCloseErr.Error())
		return
	}

	// backup the current configuration file
	_, currErr := os.Stat(path)

	// make backup path
	backupPath := path + BackupExtension

	if currErr == nil {
		// current configuration exists, back it up
		fRenameErr := os.Rename(path, backupPath)
		if fRenameErr != nil {
			log.Errorf("error renaming file %s to %s", path, backupPath)
			fileOpsPortal.ErrorHandle(model.ErrorInternalServer, fRenameErr.Error())
			return
		}
	}

	// save file to disk
	fSaveErr := fileOpsPortal.SaveToFile("conf_file", path)
	if fSaveErr != nil {
		log.Errorf("error saving file %s", path)
	} else {
		if toCallReloadEndpoint == true {
			reloadResp, reloadErr := http.Post(url+reloadPath, "application/json", nil)
			if reloadErr != nil {
				log.Errorf("Error on reload of configuration %s", reloadErr)
				fileOpsPortal.ErrorHandle(model.ErrorInternalServer, reloadErr.Error())
				return
			}
			respBody, readBodyErr := ioutil.ReadAll(reloadResp.Body)
			if readBodyErr != nil {
				log.Errorf("Error on reload of configuration %s", reloadErr)
				fileOpsPortal.ErrorHandle(model.ErrorInternalServer, readBodyErr.Error())
				return
			}
			fileOpsPortal.SuccessHandle(StatusOK, respBody)
			return
		}
		fileOpsPortal.SuccessHandle(StatusOK, nil)
		return
	}
}

func (f *FileOpsPortal) DownloadConfFile() {

	m, _ := f.GetParameters()
	confType := m["conftype"][0]

	switch confType {
	case "prometheus":
		DoDownload(f, PrometheusConfHome, PrometheusConfFile)
	case "alertmanager":
		DoDownload(f, AlertmgrConfHome, AlertmgrConfFile)
	case "grafana":
		DoDownload(f, GrafanaConfHome, GrafanaConfFile)
	}
}

func DoDownload(fileOpsPortal *FileOpsPortal, confHome string, confFile string) {
	// get the path to the configuration file
	path := confHome + confFile
	// check, if file exists
	_, currErr := os.Stat(path)
	if currErr != nil && os.IsNotExist(currErr) {
		log.Errorf("file %s not found", path)
		fileOpsPortal.ErrorHandle(model.ErrorNotFound, currErr.Error())
		return
	}
	// file exists, download it
	fileOpsPortal.Ctx.Output.Download(path, path)
}
