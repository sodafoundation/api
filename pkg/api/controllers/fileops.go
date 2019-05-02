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

	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/controller/client"
	"github.com/opensds/opensds/pkg/model"
)

// prometheus constants
var PROMETHEUS_CONF_HOME string
var PROMETHEUS_URL string
var PROMETHEUS_CONF_FILE string

// alert manager constants
var ALERTMGR_CONF_HOME string
var ALERTMGR_URL string
var ALERTMGR_CONF_FILE string

var RELOAD_PATH string
var BACKUP_EXTENSION string

func init() {

	// TODO Prakash read these from conf and save to these variables
	RELOAD_PATH = "/-/reload"
	BACKUP_EXTENSION = ".bak"

	PROMETHEUS_CONF_HOME = "/etc/prometheus/"
	PROMETHEUS_URL = "http://localhost:9090"
	PROMETHEUS_CONF_FILE = "prometheus.yml"

	ALERTMGR_CONF_HOME = "/root/alertmanager-0.16.2.linux-amd64/"
	ALERTMGR_URL = "http://localhost:9093"
	ALERTMGR_CONF_FILE = "alertmanager.yml"
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

func (this *FileOpsPortal) UploadConfFile() {

	m, _ := this.GetParameters()
	confType := m["conftype"][0]

	switch confType {
	case "prometheus":
		{
			DoUpload(this, PROMETHEUS_CONF_HOME, PROMETHEUS_URL, RELOAD_PATH)
		}
	case "alertmanager":
		{
			DoUpload(this, ALERTMGR_CONF_HOME, ALERTMGR_URL, RELOAD_PATH)
		}
	}
}

func DoUpload(this *FileOpsPortal, confHome string, url string, reloadPath string) {

	// get the uploaded file
	f, h, _ := this.GetFile("conf_file")

	// get the path to save the configuration
	path := confHome + h.Filename

	// close the incoming file
	fCloseErr := f.Close()
	if fCloseErr != nil {
		log.Errorf("error closing uploaded file %s", h.Filename)
		this.ErrorHandle(model.ErrorInternalServer, fCloseErr.Error())
		return
	}

	// backup the current configuration file
	_, currErr := os.Stat(path)

	// make backup path
	backupPath := path + BACKUP_EXTENSION

	if currErr == nil {
		// current configuration exists, back it up
		fRenameErr := os.Rename(path, backupPath)
		if fRenameErr != nil {
			log.Errorf("error renaming file %s to %s", path, backupPath)
			this.ErrorHandle(model.ErrorInternalServer, fRenameErr.Error())
			return
		}
	}

	// save file to disk
	fSaveErr := this.SaveToFile("conf_file", path)
	if fSaveErr != nil {
		log.Errorf("error saving file %s", path)
	} else {
		reloadResp, reloadErr := http.Post(url+reloadPath, "application/json", nil)
		if reloadErr != nil {
			log.Errorf("Error on reload of configuration %s", reloadErr)
			this.ErrorHandle(model.ErrorInternalServer, reloadErr.Error())
			return
		}
		respBody, readBodyErr := ioutil.ReadAll(reloadResp.Body)
		if readBodyErr != nil {
			log.Errorf("Error on reload of configuration %s", reloadErr)
			this.ErrorHandle(model.ErrorInternalServer, readBodyErr.Error())
			return
		}
		this.SuccessHandle(StatusOK, respBody)
		return
	}
}

func (this *FileOpsPortal) DownloadConfFile() {

	m, _ := this.GetParameters()
	confType := m["conftype"][0]

	switch confType {
	case "prometheus":
		{
			DoDownload(this, PROMETHEUS_CONF_HOME, PROMETHEUS_CONF_FILE)
		}
	case "alertmanager":
		{
			DoDownload(this, ALERTMGR_CONF_HOME, ALERTMGR_CONF_FILE)
		}
	}
}

func DoDownload(this *FileOpsPortal, confHome string, confFile string) {
	// get the path to the configuration file
	path := confHome + confFile
	// check, if file exists
	_, currErr := os.Stat(path)
	if currErr != nil && os.IsNotExist(currErr) {
		log.Errorf("file %s not found", path)
		this.ErrorHandle(model.ErrorNotFound, currErr.Error())
		return
	}
	// file exists, download it
	this.Ctx.Output.Download(path, path)
}
