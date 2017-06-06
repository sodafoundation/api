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
This module implements the database operation of data structure
defined in api module.

*/

package db

import (
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/clientv3"
)

var (
	TIME_OUT = 3 * time.Second
)

type DbRequest struct {
	Url        string `json:"url"`
	Content    string `json:"content"`
	NewContent string `json:"newContent"`
}

type DbResponse struct {
	Status  string   `json:"status"`
	Message []string `json:"message"`
	Error   string   `json:"error"`
}

type DbClient struct {
	Client *clientv3.Client `json:"client"`
}

func (dbCli *DbClient) Create(dReq *DbRequest) *DbResponse {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	_, err := dbCli.Client.Put(ctx, dReq.Url, dReq.Content)
	if err != nil {
		log.Println("[Error] When create db request:", err)
		return &DbResponse{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &DbResponse{
		Status:  "Success",
		Message: []string{dReq.Content},
	}
}

func (dbCli *DbClient) Get(dReq *DbRequest) *DbResponse {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	resp, err := dbCli.Client.Get(ctx, dReq.Url)
	if err != nil {
		log.Println("[Error] When get db request:", err)
		return &DbResponse{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	if len(resp.Kvs) == 0 {
		return &DbResponse{
			Status: "Failure",
			Error:  "Wrong volume_id or attachment_id provided!",
		}
	}
	return &DbResponse{
		Status:  "Success",
		Message: []string{string(resp.Kvs[0].Value)},
	}
}

func (dbCli *DbClient) List(dReq *DbRequest) *DbResponse {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	resp, err := dbCli.Client.Get(ctx, dReq.Url, clientv3.WithPrefix())
	if err != nil {
		log.Println("[Error] When get db request:", err)
		return &DbResponse{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	var message = []string{}
	for _, v := range resp.Kvs {
		message = append(message, string(v.Value))
	}
	return &DbResponse{
		Status:  "Success",
		Message: message,
	}
}

func (dbCli *DbClient) Update(dReq *DbRequest) *DbResponse {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	_, err := dbCli.Client.Put(ctx, dReq.Url, dReq.NewContent)
	if err != nil {
		log.Println("[Error] When update db request:", err)
		return &DbResponse{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &DbResponse{
		Status:  "Success",
		Message: []string{dReq.NewContent},
	}
}

func (dbCli *DbClient) Delete(dReq *DbRequest) *DbResponse {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_OUT)
	defer cancel()

	_, err := dbCli.Client.Delete(ctx, dReq.Url)
	if err != nil {
		log.Println("[Error] When delete db request:", err)
		return &DbResponse{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &DbResponse{
		Status: "Success",
	}
}

func (dbCli *DbClient) Destory() {
	dbCli.Client.Close()
}
