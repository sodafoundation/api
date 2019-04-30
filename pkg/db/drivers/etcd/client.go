// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package etcd

import (
	"context"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils"
)

var (
	timeOut  = 3 * time.Second
	retryNum = 3
)

// Request
type Request struct {
	Url        string `json:"url"`
	Content    string `json:"content"`
	NewContent string `json:"newContent"`
}

// Response
type Response struct {
	Status  string   `json:"status"`
	Message []string `json:"message"`
	Error   string   `json:"error"`
}

type clientInterface interface {
	Create(req *Request) *Response

	Get(req *Request) *Response

	List(req *Request) *Response

	Update(req *Request) *Response

	Delete(req *Request) *Response
}

// Init
func Init(edps []string) *client {
	var cliv3 *clientv3.Client
	err := utils.Retry(retryNum, "Get etcd client", false, func(retryIdx int, lastErr error) error {
		var err error
		cliv3, err = clientv3.New(clientv3.Config{
			Endpoints:   edps,
			DialTimeout: timeOut,
		})
		return err
	})
	if err != nil {
		panic(err)
	}
	return &client{cli: cliv3}
}

type client struct {
	cli  *clientv3.Client
	lock sync.Mutex
}

func (c *client) Create(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	err := utils.Retry(retryNum, "Etcd put", false, func(retryIdx int, lastErr error) error {
		_, err := c.cli.Put(ctx, req.Url, req.Content)
		return err
	})

	if err != nil {
		log.Error("When create db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &Response{
		Status:  "Success",
		Message: []string{req.Content},
	}
}

func (c *client) Get(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	resp, err := c.cli.Get(ctx, req.Url)
	if err != nil {
		log.Error("When get db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	if len(resp.Kvs) == 0 {
		return &Response{
			Status: "Failure",
			Error:  "Wrong resource uuid provided!",
		}
	}
	return &Response{
		Status:  "Success",
		Message: []string{string(resp.Kvs[0].Value)},
	}
}

func (c *client) List(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	resp, err := c.cli.Get(ctx, req.Url, clientv3.WithPrefix())
	if err != nil {
		log.Error("When get db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	var message = []string{}
	for _, v := range resp.Kvs {
		message = append(message, string(v.Value))
	}
	return &Response{
		Status:  "Success",
		Message: message,
	}
}

func (c *client) Update(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.cli.Put(ctx, req.Url, req.NewContent)
	if err != nil {
		log.Error("When update db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &Response{
		Status:  "Success",
		Message: []string{req.NewContent},
	}
}

func (c *client) Delete(req *Request) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	c.lock.Lock()
	defer c.lock.Unlock()

	_, err := c.cli.Delete(ctx, req.Url)
	if err != nil {
		log.Error("When delete db request:", err)
		return &Response{
			Status: "Failure",
			Error:  err.Error(),
		}
	}

	return &Response{
		Status: "Success",
	}
}
