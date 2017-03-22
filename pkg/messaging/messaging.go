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
This module implements the messaging service.

*/

package messaging

import (
	"log"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

type Client struct {
	cfg       client.Config
	etcd      client.KeysAPI
	watchOpts client.WatcherOptions
}

func (c *Client) Run(url string, action string) (string, error) {
	c.cfg = client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	cli, err := client.New(c.cfg)
	if err != nil {
		log.Println("Client intialized failed:", err)
		return "", err
	}
	c.etcd = client.NewKeysAPI(cli)

	_, err = c.etcd.Set(context.Background(), url, action, nil)
	if err != nil {
		log.Println("Client SET failed:", err)
		return "", err
	}

	c.watchOpts = client.WatcherOptions{AfterIndex: 0, Recursive: true}
	w := c.etcd.Watcher(url, &c.watchOpts)
	r, err := w.Next(context.Background())
	if err != nil {
		log.Println("Client WATCH failed:", err)
		return "", err
	}
	return r.Node.Value, nil
}
