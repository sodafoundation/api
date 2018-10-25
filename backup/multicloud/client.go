// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package multicloud

import (
	"net/http"
	"net/url"
	"path"

	"encoding/xml"
	"fmt"
	"github.com/astaxie/beego/httplib"
	log "github.com/golang/glog"
	"strconv"
	"time"
)

const (
	DefaultTenantId      = "adminTenantId"
	DefaultTimeout       = 60 // in Seconds
	DefaultUploadTimeout = 30 // in Seconds
	ApiVersion           = "v1"
)

type AuthOptions struct {
	Endpoint string
	UserName string
	Password string
	TenantId string
}

type Client struct {
	endpoint      string
	userName      string
	password      string
	tenantId      string
	version       string
	baseURL       string
	timeout       time.Duration
	uploadTimeout time.Duration
}

func NewClient(opt *AuthOptions, uploadTimeout int64) (*Client, error) {
	u, err := url.Parse(opt.Endpoint)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, ApiVersion)
	baseURL := u.String() + "/"

	c := &Client{
		endpoint:      opt.Endpoint,
		userName:      opt.UserName,
		password:      opt.Password,
		tenantId:      opt.TenantId,
		version:       ApiVersion,
		baseURL:       baseURL,
		timeout:       time.Duration(DefaultTimeout) * time.Minute,
		uploadTimeout: time.Duration(uploadTimeout) * time.Minute,
	}

	return c, nil
}

type ReqSettingCB func(req *httplib.BeegoHTTPRequest) error

func (c *Client) doRequest(method, u string, in interface{}, cb ReqSettingCB) ([]byte, http.Header, error) {
	req := httplib.NewBeegoRequest(u, method)
	req.Header("Content-Type", "application/xml")
	req.SetTimeout(c.timeout, c.timeout)
	if cb != nil {
		if err := cb(req); err != nil {
			return nil, nil, err
		}
	}

	if in != nil {
		var body interface{}
		switch in.(type) {
		case string, []byte:
			body = in
		default:
			body, _ = xml.Marshal(in)
		}
		req.Body(body)
	}

	resp, err := req.Response()
	if err != nil {
		log.Errorf("Do http request failed, method: %s\n url: %s\n error: %v", method, u, err)
		return nil, nil, err
	}

	log.Errorf("%s: %s OK\n", method, u)
	b, err := req.Bytes()
	if err != nil {
		log.Errorf("Get byte[] from response failed, method: %s\n url: %s\n error: %v", method, u, err)
	}
	return b, resp.Header, nil
}

func (c *Client) request(method, p string, in, out interface{}, cb ReqSettingCB) error {
	u, err := url.Parse(p)
	if err != nil {
		return err
	}
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}

	fullUrl := base.ResolveReference(u)
	b, _, err := c.doRequest(method, fullUrl.String(), in, cb)
	if err != nil {
		return err
	}

	if out != nil {
		log.V(5).Infof("Response:\n%s\n", string(b))
		err := xml.Unmarshal(b, out)
		if err != nil {
			log.Errorf("unmarshal error, reason:%v", err)
			return err
		}
	}
	return nil
}

func (c *Client) CreateBucket(backendId string) error {
	return nil
}

func (c *Client) ListBucket() ([]string, error) {
	var bucketList []string
	return bucketList, nil
}

func (c *Client) DeleteBucket(bucketName string) error {

	return nil
}

type Object struct {
	ObjectKey  string `xml:"ObjectKey"`
	BucketName string `xml:"BucketName"`
	Size       uint64 `xml:"Size"`
}

type ListObjectResponse struct {
	ListObjects []Object `xml:"ListObjects"`
}

type InitiateMultipartUploadResult struct {
	Xmlns    string `xml:"xmlns,attr"`
	Bucket   string `xml:"Bucket"`
	Key      string `xml:"Key"`
	UploadId string `xml:"UploadId"`
}

type UploadPartResult struct {
	Xmlns      string `xml:"xmlns,attr"`
	PartNumber int64  `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}

type Part struct {
	PartNumber int64  `xml:"PartNumber"`
	ETag       string `xml:"ETag"`
}

type CompleteMultipartUpload struct {
	Xmlns string `xml:"xmlns,attr"`
	Part  []Part `xml:"Part"`
}

type CompleteMultipartUploadResult struct {
	Xmlns    string `xml:"xmlns,attr"`
	Location string `xml:"Location"`
	Bucket   string `xml:"Bucket"`
	Key      string `xml:"Key"`
	ETag     string `xml:"ETag"`
}

func (c *Client) UploadObject(bucketName, objectKey string, data []byte) error {
	p := path.Join("s3", bucketName, objectKey)
	if err := c.request("PUT", p, data, nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) ListObject(bucketName string) (*ListObjectResponse, error) {
	p := path.Join("s3", bucketName)
	object := &ListObjectResponse{}
	if err := c.request("GET", p, nil, object, nil); err != nil {
		return nil, err
	}
	return object, nil
}

func (c *Client) RemoveObject(bucketName, objectKey string) error {
	p := path.Join("s3", bucketName, objectKey)
	if err := c.request("DELETE", p, nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) InitMultiPartUpload(bucketName, objectKey string) (*InitiateMultipartUploadResult, error) {
	p := path.Join("s3", bucketName, objectKey)
	p += "?uploads"
	out := &InitiateMultipartUploadResult{}
	if err := c.request("PUT", p, nil, out, nil); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) UploadPart(bucketName, objectKey string, partNum int64, uploadId string, data []byte, size int64) (*UploadPartResult, error) {
	log.Infof("upload part buf size:%d", len(data))
	p := path.Join("s3", bucketName, objectKey)
	p += fmt.Sprintf("?partNumber=%d&uploadId=%s", partNum, uploadId)
	out := &UploadPartResult{}
	reqSettingCB := func(req *httplib.BeegoHTTPRequest) error {
		req.Header("Content-Length", strconv.FormatInt(size, 10))
		req.SetTimeout(c.uploadTimeout, c.uploadTimeout)
		return nil
	}
	if err := c.request("PUT", p, data, out, reqSettingCB); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) CompleteMultipartUpload(
	bucketName string,
	objectKey string,
	uploadId string,
	input *CompleteMultipartUpload) (*CompleteMultipartUploadResult, error) {

	p := path.Join("s3", bucketName, objectKey)
	p += fmt.Sprintf("?uploadId=%s", uploadId)
	out := &CompleteMultipartUploadResult{}
	if err := c.request("PUT", p, input, nil, nil); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) AbortMultipartUpload(bucketName, objectKey string) error {
	// TODO: multi-cloud has not implemented it yet. so just comment it.
	//p := path.Join("s3", "AbortMultipartUpload", bucketName, objectKey)
	//if err := c.request("DELETE", p, nil, nil); err != nil {
	//	return err
	//}
	return nil
}
