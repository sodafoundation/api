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
This module implements the grpc server.

*/

package server

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"

	adapterApi "adapter/storageDock/api"
	orchestrationApi "orchestration/api"
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
		log.Fatal(err)
	}
	s.etcd = client.NewKeysAPI(cli)
	resp, err := s.etcd.Set(context.Background(), "opensds/api", "Server start!", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Set is done. Index is %d\n", resp.Index)
	}
	s.watchOpts = client.WatcherOptions{AfterIndex: resp.Index, Recursive: true}
}

func (s *Server) OrchestrationWatch(url string) {
	w := s.etcd.Watcher(url, &s.watchOpts)
	r, err := w.Next(context.Background())
	if err != nil {
		log.Fatal("Error occurred", err)
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
		result, err = orchestrationApi.CreateVolume(resourceType, name, size)
		if err != nil {
			log.Println("Error occured when create volume!")
		}
	case "GetVolume":
		resourceType := tmp[1]
		volID := tmp[2]
		result, err = orchestrationApi.GetVolume(resourceType, volID)
		if err != nil {
			log.Println("Error occured when get volume!")
		}
	case "GetAllVolumes":
		resourceType := tmp[1]
		result, err = orchestrationApi.GetAllVolumes(resourceType)
		if err != nil {
			log.Println("Error occured when get all volumes!")
		}
	case "UpdateVolume":
		resourceType := tmp[1]
		volID := tmp[2]
		name := tmp[3]
		result, err = orchestrationApi.UpdateVolume(resourceType, volID, name)
		if err != nil {
			log.Println("Error occured when update volume!")
		}
	case "DeleteVolume":
		resourceType := tmp[1]
		volID := tmp[2]
		result, err = orchestrationApi.DeleteVolume(resourceType, volID)
		if err != nil {
			log.Println("Error occured when delete volume!")
		}
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
		result = "Error"
	}

	_, err = s.etcd.Set(context.Background(), url, result, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Set is done. URL is %s.\n", url)
	}
}

func (s *Server) AdapterWatch(url string) {
	w := s.etcd.Watcher(url, &s.watchOpts)
	r, err := w.Next(context.Background())
	if err != nil {
		log.Fatal("Error occurred", err)
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
		result, err = adapterApi.CreateVolume(resourceType, name, size)
		if err != nil {
			log.Println("Error occured when create volume!")
		}
	case "GetVolume":
		resourceType := tmp[1]
		volID := tmp[2]
		result, err = adapterApi.GetVolume(resourceType, volID)
		if err != nil {
			log.Println("Error occured when get volume!")
		}
	case "GetAllVolumes":
		resourceType := tmp[1]
		result, err = adapterApi.GetAllVolumes(resourceType)
		if err != nil {
			log.Println("Error occured when get all volumes!")
		}
	case "UpdateVolume":
		resourceType := tmp[1]
		volID := tmp[2]
		name := tmp[3]
		result, err = adapterApi.UpdateVolume(resourceType, volID, name)
		if err != nil {
			log.Println("Error occured when update volume!")
		}
	case "DeleteVolume":
		resourceType := tmp[1]
		volID := tmp[2]
		result, err = adapterApi.DeleteVolume(resourceType, volID)
		if err != nil {
			log.Println("Error occured when delete volume!")
		}
	default:
		log.Printf("Error, no action: %s\n", tmp[0])
		result = "Error"
	}

	_, err = s.etcd.Set(context.Background(), url, result, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Set is done. URL is %s.\n", url)
	}
}
