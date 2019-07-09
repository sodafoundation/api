// Copyright 2018 The OpenSDS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package multicloud

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/opensds/opensds/contrib/backup"
	"github.com/opensds/opensds/pkg/utils"
	"gopkg.in/yaml.v2"
)

const (
	ConfFile  = "/etc/opensds/driver/multi-cloud.yaml"
	ChunkSize = 1024 * 1024 * 50
)

func init() {
	backup.RegisterBackupCtor("multi-cloud", NewMultiCloud)
}

func NewMultiCloud() (backup.BackupDriver, error) {
	return &MultiCloud{}, nil
}

type AuthOptions struct {
	Strategy        string `yaml:"Strategy"`
	AuthUrl         string `yaml:"AuthUrl,omitempty"`
	DomainName      string `yaml:"DomainName,omitempty"`
	UserName        string `yaml:"UserName,omitempty"`
	Password        string `yaml:"Password,omitempty"`
	PwdEncrypter    string `yaml:"PwdEncrypter,omitempty"`
	EnableEncrypted bool   `yaml:"EnableEncrypted,omitempty"`
	TenantName      string `yaml:"TenantName,omitempty"`
}

type MultiCloudConf struct {
	Endpoint      string `yaml:"Endpoint,omitempty"`
	UploadTimeout int64  `yaml:"UploadTimeout,omitempty"`
	AuthOptions   `yaml:"AuthOptions,omitempty"`
}

type MultiCloud struct {
	client *Client
	conf   *MultiCloudConf
}

func (m *MultiCloud) loadConf(p string) (*MultiCloudConf, error) {
	conf := &MultiCloudConf{
		Endpoint:      "http://127.0.0.1:8088",
		UploadTimeout: DefaultUploadTimeout,
	}
	confYaml, err := ioutil.ReadFile(p)
	if err != nil {
		glog.Errorf("Read config yaml file (%s) failed, reason:(%v)", p, err)
		return nil, err
	}
	if err = yaml.Unmarshal(confYaml, conf); err != nil {
		glog.Errorf("Parse error: %v", err)
		return nil, err
	}
	return conf, nil
}

func (m *MultiCloud) SetUp() error {
	// Set the default value
	var err error
	if m.conf, err = m.loadConf(ConfFile); err != nil {
		return err
	}

	if m.client, err = NewClient(m.conf.Endpoint, &m.conf.AuthOptions, m.conf.UploadTimeout); err != nil {
		return err
	}

	return nil
}

func (m *MultiCloud) CleanUp() error {
	// Do nothing
	return nil
}

func (m *MultiCloud) Backup(backup *backup.BackupSpec, volFile *os.File) error {
	buf := make([]byte, ChunkSize)
	input := &CompleteMultipartUpload{}

	bucket, ok := backup.Metadata["bucket"]
	if !ok {
		return errors.New("can't find bucket in metadata")
	}
	key := backup.Id
	initResp, err := m.client.InitMultiPartUpload(bucket, key)
	if err != nil {
		glog.Errorf("Init part failed, err:%v", err)
		return err
	}

	defer m.client.AbortMultipartUpload(bucket, key)
	var parts []Part
	for partNum := int64(1); ; partNum++ {
		size, err := volFile.Read(buf)
		glog.Infof("read buf size len:%d", size)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if size == 0 {
			break
		}
		var uploadResp *UploadPartResult
		err = utils.Retry(3, "upload part", false, func(retryIdx int, lastErr error) error {
			var inErr error
			uploadResp, inErr = m.client.UploadPart(bucket, key, partNum, initResp.UploadId, buf[:size], int64(size))
			return inErr
		})
		if err != nil {
			glog.Errorf("upload part failed, err:%v", err)
			return err
		}
		parts = append(parts, Part{PartNumber: partNum, ETag: uploadResp.ETag})
	}
	input.Part = parts
	_, err = m.client.CompleteMultipartUpload(bucket, key, initResp.UploadId, input)
	if err != nil {
		glog.Errorf("complete part failed, err:%v", err)
		return err
	}
	m.client.AbortMultipartUpload(bucket, key)
	glog.Infof("backup success ...")
	return nil
}

func (m *MultiCloud) Restore(backup *backup.BackupSpec, backupId string, volFile *os.File) error {
	bucket, ok := backup.Metadata["bucket"]
	if !ok {
		return errors.New("can't find bucket in metadata")
	}
	var downloadSize = ChunkSize
	// if the size of data of smaller than require download size
	// downloading is completed.
	for offset := int64(0); downloadSize == ChunkSize; offset += ChunkSize {
		var data []byte
		err := utils.Retry(3, "download part", false, func(retryIdx int, lastErr error) error {
			var inErr error
			data, inErr = m.client.DownloadPart(bucket, backupId, offset, ChunkSize)
			return inErr
		})
		if err != nil {
			glog.Errorf("download part failed: %v", err)
			return err
		}
		downloadSize = len(data)
		glog.V(5).Infof("download size: %d\n", downloadSize)
		volFile.Seek(offset, 0)
		size, err := volFile.Write(data)
		if err != nil {
			glog.Errorf("write part failed: %v", err)
			return err
		}
		if size != downloadSize {
			return errors.New("size not equal to download size")
		}
		glog.V(5).Infof("write buf size len:%d", size)
	}
	glog.Infof("restore success ...")
	return nil
}

func (m *MultiCloud) Delete(backup *backup.BackupSpec) error {
	bucket := backup.Metadata["bucket"]
	key := backup.Id
	return m.client.RemoveObject(bucket, key)
}
