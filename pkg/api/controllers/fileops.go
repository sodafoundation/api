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
	"fmt"
	"net/http"
	"os"

	"github.com/opensds/opensds/pkg/model"

	"github.com/opensds/opensds/pkg/controller/client"
)

var PROMETHEUS_CONF_HOME string
var PROMETHEUS_URL string
var RELOAD_PATH string
var BACKUP_EXTENSION string
var PROMETHEUS_CONF_FILE string

func init() {

	// TODO read these from conf and save to this variable
	PROMETHEUS_CONF_HOME = "/etc/prometheus/"
	PROMETHEUS_URL = "http://localhost:9090"
	RELOAD_PATH = "/-/reload"
	BACKUP_EXTENSION = ".bak"
	PROMETHEUS_CONF_FILE = "prometheus.yml"
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

	m, err := this.GetParameters()
	confType := m["conftype"][0]

	fmt.Println(confType, m, err)

	// get the uploaded file
	f, h, _ := this.GetFile("conf_file")

	switch confType {
	case "prometheus":
		{
			// get the path to save the prometheus configuration
			path := PROMETHEUS_CONF_HOME + h.Filename

			// close the incoming file
			fCloseErr := f.Close()
			if fCloseErr != nil {
				//log.Errorf("error closing uploaded file %s", h.Filename)
				this.ErrorHandle(model.ErrorInternalServer, fCloseErr.Error())
				return
			}

			// backup the current prometheus configuration file
			_, currErr := os.Stat(path)

			// make backup path
			backupPath := path + BACKUP_EXTENSION

			if currErr == nil {
				// current configuration exists, back it up
				fRenameErr := os.Rename(path, backupPath)
				if fRenameErr != nil {
					//log.Errorf("error renaming file %s to %s", path, backupPath)
					this.ErrorHandle(model.ErrorInternalServer, fRenameErr.Error())
					return
				}
			}

			// save file to disk
			fSaveErr := this.SaveToFile("conf_file", path)
			if fSaveErr != nil {
				//log.Errorf("error saving file %s", path)
			} else {
				_, reloadErr := http.Post(PROMETHEUS_URL+RELOAD_PATH, "application/json", nil)
				if reloadErr != nil {
					//log.Errorf("Error on reload of Prometheus configuration %s", reloadErr)
					this.ErrorHandle(model.ErrorInternalServer, reloadErr.Error())
					return
				}
				this.SuccessHandle(StatusOK, nil)
				return
			}
		}
	}
}

func (this *FileOpsPortal) DownloadConfFile() {

	m, err := this.GetParameters()
	confType := m["conftype"][0]

	fmt.Println(confType, m, err)

	switch confType {
	case "prometheus":
		{
			// get the path to save the prometheus configuration
			path := PROMETHEUS_CONF_HOME + PROMETHEUS_CONF_FILE

			this.Ctx.Output.Download(path, path)
		}
	}

}
