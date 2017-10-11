// Copyright (c) 2017 Huawei Technologies Co., Ltd. All Rights Reserved.
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
This module implements the entry into operations of storageDock module.

*/

package discovery

import (
	"os"

	"github.com/opensds/opensds/pkg/db"
	dockHub "github.com/opensds/opensds/pkg/dock"
	api "github.com/opensds/opensds/pkg/model"
	"github.com/opensds/opensds/pkg/utils"
	. "github.com/opensds/opensds/pkg/utils/config"

	log "github.com/golang/glog"
	"github.com/satori/go.uuid"
)

type Discoverer interface {
	Init() error
	Discovery() error
	Store() error
}

type DockDiscoverer struct {
	dcks []api.DockSpec
	pols []api.StoragePoolSpec

	c db.Client
}

func NewDiscover() Discoverer {
	return &DockDiscoverer{
		c: db.C,
	}
}

func (dd *DockDiscoverer) Init() error {
	// Load resource from specified file
	name2Backend := map[string]BackendProperties{
		"ceph":   BackendProperties(CONF.Ceph),
		"cinder": BackendProperties(CONF.Cinder),
		"sample": BackendProperties(CONF.Sample),
	}

	host, err := os.Hostname()
	if err != nil {
		log.Error("When get os hostname:", err)
		return err
	}

	for _, v := range CONF.EnableBackends {
		b := name2Backend[v]
		if b.Name == "" {
			continue
		}

		dock := api.DockSpec{
			BaseModel: &api.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, host+":"+b.DriverName).String(),
			},
			Name:        b.Name,
			Description: b.Description,
			DriverName:  b.DriverName,
			Endpoint:    CONF.OsdsDock.ApiEndpoint,
		}
		dd.dcks = append(dd.dcks, dock)
	}
	return nil
}

func (dd *DockDiscoverer) Discovery() error {
	var pols *[]api.StoragePoolSpec
	var err error

	for _, dock := range dd.dcks {
		pols, err = dockHub.NewDockHub(dock.GetDriverName()).ListPools()
		if err != nil {
			log.Error("When list pools:", err)
			return err
		}

		if len(*pols) == 0 {
			log.Warningf("The pool of dock %s is empty!\n", dock.GetId())
		}

		for _, pol := range *pols {
			pol.DockId = dock.GetId()
		}
		dd.pols = append(dd.pols, *pols...)
	}

	return err
}

func (dd *DockDiscoverer) Store() error {
	var err error

	// Store dock resources in database.
	for _, dock := range dd.dcks {
		if err = utils.ValidateData(&dock, utils.S); err != nil {
			log.Error("When validate dock structure:", err)
			return err
		}

		// Call db module to create dock resource.
		if _, err = db.C.CreateDock(&dock); err != nil {
			log.Error("When create dock %s in db: %v\n", dock.GetId(), err)
			return err
		}
	}

	// Store pool resources in database.
	for _, pool := range dd.pols {
		if err = utils.ValidateData(&pool, utils.S); err != nil {
			log.Error("When validate pool structure:", err)
			return err
		}

		// Call db module to create pool resource.
		if _, err = db.C.CreatePool(&pool); err != nil {
			log.Error("When create pool %s in db: %v\n", pool.GetId(), err)
			return err
		}
	}

	return err
}

func Discovery(d Discoverer) error {
	var err error

	if err = d.Init(); err != nil {
		return err
	}
	if err = d.Discovery(); err != nil {
		return err
	}
	if err = d.Store(); err != nil {
		return err
	}

	return err
}
