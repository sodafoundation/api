// Copyright (c) 2019 OpenSDS Authors. All Rights Reserved.
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
This module implements the entry into operations of storageDock module.

*/

package discovery

import (
	"fmt"
	"os"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/filesharedrivers"
	c "github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/model"
	. "github.com/opensds/opensds/pkg/utils/config"
	"github.com/satori/go.uuid"
)

type Context struct {
	StopChan chan bool
	ErrChan  chan error
	MetaChan chan string
}

func DiscoveryAndReport(dd DockDiscoverer, ctx *Context) {
	for {
		select {
		case <-ctx.StopChan:
			return
		default:
			if err := dd.Discover(); err != nil {
				ctx.ErrChan <- err
			}

			if err := dd.Report(); err != nil {
				ctx.ErrChan <- err
			}
		}

		time.Sleep(60 * time.Second)
	}
}

type DockDiscoverer interface {
	Init() error

	Discover() error

	Report() error
}

// NewDockDiscoverer method creates a new DockDiscoverer.
func NewDockDiscoverer(dockType string) DockDiscoverer {
	switch dockType {
	case model.DockTypeProvioner:
		return &provisionDockDiscoverer{
			DockRegister: NewDockRegister(),
		}
	}
	return nil
}

// provisionDockDiscoverer is a struct for exposing some operations of provision
// dock service discovery.
type provisionDockDiscoverer struct {
	*DockRegister

	dcks []*model.DockSpec
	pols []*model.StoragePoolSpec
}

func (pdd *provisionDockDiscoverer) Init() error {
	// Load resource from specified file
	bm := GetBackendsMap()
	host, err := os.Hostname()
	if err != nil {
		log.Error("When get os hostname:", err)
		return err
	}

	for _, v := range CONF.EnabledBackends {
		b := bm[v]
		if b.Name == "" {
			continue
		}

		dck := &model.DockSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, host+":"+b.DriverName).String(),
			},
			Name:        b.Name,
			Description: b.Description,
			DriverName:  b.DriverName,
			Endpoint:    CONF.OsdsDock.ApiEndpoint,
			NodeId:      host,
			Type:        model.DockTypeProvioner,
			Metadata:    map[string]string{"HostReplicationDriver": CONF.OsdsDock.HostBasedReplicationDriver},
		}
		pdd.dcks = append(pdd.dcks, dck)
	}

	return nil
}

func (pdd *provisionDockDiscoverer) Discover() error {
	// Clear existing pool info
	pdd.pols = pdd.pols[:0]

	for _, dck := range pdd.dcks {
		// Call function of StorageDrivers configured by storage drivers.
		pols, err := filesharedrivers.Init(dck.DriverName).ListPools()
		if err != nil {
			log.Error("Call driver to list pools failed:", err)
			continue
		}

		if len(pols) == 0 {
			log.Warningf("The pool of dock %s is empty!\n", dck.Id)
		}
		pdd.pols = append(pdd.pols, pols...)
	}
	if len(pdd.pols) == 0 {
		return fmt.Errorf("There is no pool can be found.")
	}

	return nil
}

func (pdd *provisionDockDiscoverer) Report() error {
	var err error

	// Store dock resources in database.
	for _, dck := range pdd.dcks {
		if err = pdd.Register(dck); err != nil {
			break
		}
	}

	// Store pool resources in database.
	for _, pol := range pdd.pols {
		if err != nil {
			break
		}
		err = pdd.Register(pol)
	}

	return err
}

func NewDockRegister() *DockRegister {
	return &DockRegister{c: db.C}
}

type DockRegister struct {
	c db.Client
}

func (dr *DockRegister) Register(in interface{}) error {
	ctx := c.NewAdminContext()

	switch in.(type) {
	case *model.DockSpec:
		dck := in.(*model.DockSpec)
		// Call db module to create dock resource.
		if _, err := dr.c.CreateDock(ctx, dck); err != nil {
			log.Errorf("When create dock %s in db: %v\n", dck.Id, err)
			return err
		}
		break
	case *model.StoragePoolSpec:
		pol := in.(*model.StoragePoolSpec)
		// Call db module to create pool resource.
		if _, err := dr.c.CreatePool(ctx, pol); err != nil {
			log.Errorf("When create pool %s in db: %v\n", pol.Id, err)
			return err
		}
		break
	default:
		return fmt.Errorf("Resource type is not supported!")
	}

	return nil
}

func (dr *DockRegister) Unregister(in interface{}) error {
	ctx := c.NewAdminContext()

	switch in.(type) {
	case *model.DockSpec:
		dck := in.(*model.DockSpec)
		// Call db module to delete dock resource.
		if err := dr.c.DeleteDock(ctx, dck.Id); err != nil {
			log.Errorf("When delete dock %s in db: %v\n", dck.Id, err)
			return err
		}
		break
	case *model.StoragePoolSpec:
		pol := in.(*model.StoragePoolSpec)
		// Call db module to delete pool resource.
		if err := dr.c.DeletePool(ctx, pol.Id); err != nil {
			log.Errorf("When delete pool %s in db: %v\n", pol.Id, err)
			return err
		}
		break
	default:
		return fmt.Errorf("Resource type is not supported!")
	}

	return nil
}
