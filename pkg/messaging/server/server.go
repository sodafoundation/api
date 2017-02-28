// Copyright (c) 2016 Huawei Technologies Co., Ltd. All Rights Reserved.
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

/*
This module implements the messaging server service.

*/

package server

import (
	_ "encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	adapterApi "github.com/opensds/opensds/pkg/adapter/dock/api"
	orchestrationApi "github.com/opensds/opensds/pkg/orchestration/api"
)

type Server struct {
	cfg       client.Config
	etcd      client.KeysAPI
	watchOpts client.WatcherOptions
}

func (s *Server) Init() {
	s.cfg = client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	cli, err := client.New(s.cfg)
	if err != nil {
		log.Fatalln("Server initialized failed:", err)
	}
	s.etcd = client.NewKeysAPI(cli)
	s.watchOpts = client.WatcherOptions{AfterIndex: 0, Recursive: true}
	log.Println("Server initialized success!")
}

func (s *Server) OrchestrationWatch(url string) {
	for {
		w := s.etcd.Watcher(url, &s.watchOpts)
		r, err := w.Next(context.Background())
		if err != nil {
			log.Fatalln("Orchestration modlue server WATCH failed:", err)
		}

		value := r.Node.Value
		tmp := make([]string, 5, 10)
		tmp = strings.Split(value, ",")
		var result string

		switch tmp[0] {
		case "CreateVolume":
			resourceType := tmp[1]
			name := tmp[2]
			size, _ := strconv.Atoi(tmp[3])
			result, _ = orchestrationApi.CreateVolume(resourceType, name, size)
		case "GetVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			result, _ = orchestrationApi.GetVolume(resourceType, volID)
		case "GetAllVolumes":
			resourceType := tmp[1]
			allowDetails, _ := strconv.ParseBool(tmp[2])
			result, _ = orchestrationApi.GetAllVolumes(resourceType, allowDetails)
		case "UpdateVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			name := tmp[3]
			result, _ = orchestrationApi.UpdateVolume(resourceType, volID, name)
		case "DeleteVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			result, _ = orchestrationApi.DeleteVolume(resourceType, volID)
		case "AttachVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			host := tmp[3]
			device := tmp[4]
			result, _ = orchestrationApi.AttachVolume(resourceType, volID, host, device)
		case "DetachVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			attachment := tmp[3]
			result, _ = orchestrationApi.DetachVolume(resourceType, volID, attachment)
		case "MountVolume":
			mountDir := tmp[1]
			device := tmp[2]
			fsType := tmp[4]
			result, _ = orchestrationApi.MountVolume(mountDir, device, fsType)
		case "UnmountVolume":
			mountDir := tmp[1]
			result, _ = orchestrationApi.UnmountVolume(mountDir)
		case "CreateShare":
			resourceType := tmp[1]
			name := tmp[2]
			shrType := tmp[3]
			shrProto := tmp[4]
			size, _ := strconv.Atoi(tmp[5])
			result, _ = orchestrationApi.CreateShare(resourceType, name, shrType, shrProto, size)
		case "GetShare":
			resourceType := tmp[1]
			shrID := tmp[2]
			result, _ = orchestrationApi.GetShare(resourceType, shrID)
		case "GetAllShares":
			resourceType := tmp[1]
			allowDetails, _ := strconv.ParseBool(tmp[2])
			result, _ = orchestrationApi.GetAllShares(resourceType, allowDetails)
		case "UpdateShare":
			resourceType := tmp[1]
			shrID := tmp[2]
			name := tmp[3]
			result, _ = orchestrationApi.UpdateShare(resourceType, shrID, name)
		case "DeleteShare":
			resourceType := tmp[1]
			shrID := tmp[2]
			result, _ = orchestrationApi.DeleteShare(resourceType, shrID)
		case "CreateDatabase":
			name := tmp[1]
			size, _ := strconv.Atoi(tmp[2])
			result, err = orchestrationApi.CreateDatabase(name, size)
			if err != nil {
				log.Println("Error occured when create database!")
			}
		case "GetDatabase":
			id, _ := strconv.Atoi(tmp[1])
			name := tmp[2]
			result, err = orchestrationApi.GetDatabase(id, name)
			if err != nil {
				log.Println("Error occured when get database!")
			}
		case "GetAllDatabases":
			result, err = orchestrationApi.GetAllDatabases()
			if err != nil {
				log.Println("Error occured when get all databases!")
			}
		case "UpdateDatabase":
			id, _ := strconv.Atoi(tmp[1])
			size, _ := strconv.Atoi(tmp[2])
			name := tmp[3]
			result, err = orchestrationApi.UpdateDatabase(id, size, name)
			if err != nil {
				log.Println("Error occured when update database!")
			}
		case "DeleteDatabase":
			id, _ := strconv.Atoi(tmp[1])
			name := tmp[2]
			cascade, _ := strconv.ParseBool(tmp[3])
			result, err = orchestrationApi.DeleteDatabase(id, name, cascade)
			if err != nil {
				log.Println("Error occured when delete database!")
			}
		case "CreateFileSystem":
			name := tmp[1]
			size, _ := strconv.Atoi(tmp[2])
			result, err = orchestrationApi.CreateFileSystem(name, size)
			if err != nil {
				log.Println("Error occured when create file system!")
			}
		case "GetFileSystem":
			id, _ := strconv.Atoi(tmp[1])
			name := tmp[2]
			result, err = orchestrationApi.GetFileSystem(id, name)
			if err != nil {
				log.Println("Error occured when get file system!")
			}
		case "GetAllFileSystems":
			result, err = orchestrationApi.GetAllFileSystems()
			if err != nil {
				log.Println("Error occured when get all file systems!")
			}
		case "UpdateFileSystem":
			id, _ := strconv.Atoi(tmp[1])
			size, _ := strconv.Atoi(tmp[2])
			name := tmp[3]
			result, err = orchestrationApi.UpdateFileSystem(id, size, name)
			if err != nil {
				log.Println("Error occured when update file system!")
			}
		case "DeleteFileSystem":
			id, _ := strconv.Atoi(tmp[1])
			name := tmp[2]
			cascade, _ := strconv.ParseBool(tmp[3])
			result, err = orchestrationApi.DeleteFileSystem(id, name, cascade)
			if err != nil {
				log.Println("Error occured when delete file system!")
			}
		default:
			log.Printf("Error, no action: %s\n", tmp[0])
			result = ""
		}

		_, err = s.etcd.Set(context.Background(), url, result, nil)
		if err != nil {
			log.Fatalln("Orchesration modlue server SET failed:", err)
		}
	}
}

func (s *Server) AdapterWatch(url string) {
	for {
		w := s.etcd.Watcher(url, &s.watchOpts)
		r, err := w.Next(context.Background())
		if err != nil {
			log.Fatalln("Adapter modlue server WATCH failed:", err)
		}

		value := r.Node.Value
		tmp := make([]string, 5, 10)
		tmp = strings.Split(value, ",")
		var result string

		switch tmp[0] {
		case "CreateVolume":
			resourceType := tmp[1]
			name := tmp[2]
			size, _ := strconv.Atoi(tmp[3])
			result, _ = adapterApi.CreateVolume(resourceType, name, size)
		case "GetVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			result, _ = adapterApi.GetVolume(resourceType, volID)
		case "GetAllVolumes":
			resourceType := tmp[1]
			allowDetails, _ := strconv.ParseBool(tmp[2])
			result, _ = adapterApi.GetAllVolumes(resourceType, allowDetails)
		case "UpdateVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			name := tmp[3]
			result, _ = adapterApi.UpdateVolume(resourceType, volID, name)
		case "DeleteVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			result, _ = adapterApi.DeleteVolume(resourceType, volID)
		case "AttachVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			host := tmp[3]
			device := tmp[4]
			result, _ = adapterApi.AttachVolume(resourceType, volID, host, device)
		case "DetachVolume":
			resourceType := tmp[1]
			volID := tmp[2]
			attachment := tmp[3]
			result, _ = adapterApi.DetachVolume(resourceType, volID, attachment)
		case "MountVolume":
			mountDir := tmp[1]
			device := tmp[2]
			fsType := tmp[3]
			result, _ = adapterApi.MountVolume(mountDir, device, fsType)
		case "UnmountVolume":
			mountDir := tmp[1]
			result, _ = adapterApi.UnmountVolume(mountDir)
		case "CreateShare":
			resourceType := tmp[1]
			name := tmp[2]
			shrType := tmp[3]
			shrProto := tmp[4]
			size, _ := strconv.Atoi(tmp[5])
			result, _ = adapterApi.CreateShare(resourceType, name, shrType, shrProto, size)
		case "GetShare":
			resourceType := tmp[1]
			shrID := tmp[2]
			result, _ = adapterApi.GetShare(resourceType, shrID)
		case "GetAllShares":
			resourceType := tmp[1]
			allowDetails, _ := strconv.ParseBool(tmp[2])
			result, _ = adapterApi.GetAllShares(resourceType, allowDetails)
		case "UpdateShare":
			resourceType := tmp[1]
			shrID := tmp[2]
			name := tmp[3]
			result, _ = adapterApi.UpdateShare(resourceType, shrID, name)
		case "DeleteShare":
			resourceType := tmp[1]
			shrID := tmp[2]
			result, _ = adapterApi.DeleteShare(resourceType, shrID)
		default:
			log.Printf("Error, no action: %s\n", tmp[0])
			result = ""
		}

		_, err = s.etcd.Set(context.Background(), url, result, nil)
		if err != nil {
			log.Fatalln("Adapter modlue server SET failed:", err)
		}
	}
}
