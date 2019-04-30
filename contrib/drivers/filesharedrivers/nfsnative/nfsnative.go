// Copyright 2019 The OpenSDS Authors.
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

package nfsnative

import (
	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	"github.com/opensds/opensds/pkg/utils/config"
	. "github.com/opensds/opensds/testutils/collection"
)

const (
	opensdsPrefix   = "opensds-"
	defaultConfPath = "/etc/opensds/driver/nfsnative.yaml"
	Configpath      = "/etc/nfsnative/nfsnative.conf"
)

type NFSNativeConfig struct {
	ConfigFile string                    `yaml:"configFile,omitempty"`
	Pool       map[string]PoolProperties `yaml:"pool,flow"`
}

func EncodeName(id string) string {
	return opensdsPrefix + id
}

type Driver struct {
	conf *NFSNativeConfig
	cli  *Cli
}

func (d *Driver) ListPools() ([]*model.StoragePoolSpec, error) {
	// TODO Sample test pool will load to DB for driver code test
	var pols []*model.StoragePoolSpec
	pols = append(pols, &SamplePools[1])
	return pols, nil
}

func (d *Driver) Setup() error {
	d.conf = &NFSNativeConfig{ConfigFile: Configpath}
	p := config.CONF.OsdsfileDock.FileShareBackends.NFSNative.ConfigPath
	if p == "" {
		p = defaultConfPath
	}
	_, err := Parse(d.conf, p)
	return err
}

func (d *Driver) Unset() error { return nil }

func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {

	// TODO Need to implement, it is dummy driver code
	log.Infof("Create nfs native %s (%s) success.", opt.GetName(), opt.GetId())
	return &SampleFileShares[0], nil
}

func (d *Driver) DeleteFileShare(opt *pb.DeleteFileShareOpts) (*model.FileShareSpec, error) {

	// TODO Need to implement
	return nil, nil
}
