// Copyright 2017 The OpenSDS Authors.
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
	consts "github.com/opensds/opensds/contrib/drivers/utils/constants"
	"os"
	"runtime"
	"strings"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/connector"
	"github.com/opensds/opensds/contrib/drivers"
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
func NewDockDiscoverer(dockType string, driverMgr *drivers.DriverManager) DockDiscoverer {
	switch dockType {
	case model.DockTypeProvioner:
		return &provisionDockDiscoverer{
			DockRegister: NewDockRegister(),
			driverMgr:    driverMgr,
		}
	case model.DockTypeAttacher:
		return &attachDockDiscoverer{
			DockRegister: NewDockRegister(),
		}
	}
	return nil
}

// provisionDockDiscoverer is a struct for exposing some operations of provision
// dock service discovery.
type provisionDockDiscoverer struct {
	*DockRegister
	driverMgr *drivers.DriverManager
	dcks      []*model.DockSpec
	pols      []*model.StoragePoolSpec
}

func (pdd *provisionDockDiscoverer) Init() error {
	// Load resource from specified file
	hostName, err := os.Hostname()
	if err != nil {
		log.Error("When get os hostname:", err)
		return err
	}

	for _, name := range CONF.EnabledBackends {
		bp := CONF.OsdsDock.BackendMap[name]
		if len(bp.DriverName) == 0 || len(bp.Name) == 0 {
			log.Errorf("invalid backend (%s) properties, ignore it", name)
			continue
		}
		if bp.StorageType != string(consts.StorageTypeBlock) && bp.StorageType != (consts.StorageTypeFile) {
			log.Errorf("backend (%s) properties: storage_type must be %v or %v ,ignore it",
				name, consts.StorageTypeBlock, consts.StorageTypeFile)
			continue
		}

		baseName := strings.Join([]string{hostName, bp.DriverName, bp.StorageType}, ":")
		dock := &model.DockSpec{
			BaseModel: &model.BaseModel{
				Id: uuid.NewV5(uuid.NamespaceOID, baseName).String(),
			},
			Name:        bp.Name,
			Description: bp.Description,
			// DriverName:  bp.DriverName, use driver name temporarily
			DriverName:  name,
			BackendName: name,
			Endpoint:    CONF.OsdsDock.ApiEndpoint,
			NodeId:      hostName,
			Type:        model.DockTypeProvioner,
			StorageType: bp.StorageType,
			Metadata:    map[string]string{},
		}
		pdd.dcks = append(pdd.dcks, dock)
	}
	return nil
}

type PoolDriver interface {
	ListPools() ([]*model.StoragePoolSpec, error)
}

func (pdd *provisionDockDiscoverer) discoverPool(backendName, dockId string) ([]*model.StoragePoolSpec, error) {
	d, err := pdd.driverMgr.GetDriverWrapper(consts.DriverTypeProvision, backendName)
	if err != nil {
		log.Error("Get driver failed:", err)
		return nil, err
	}

	pols, err := d.Driver.(PoolDriver).ListPools()
	if err != nil {
		log.Error("Call driver %s to list pools failed:", backendName, err)
		return nil, err
	}

	for _, pol := range pols {
		log.Infof("Backend %s discovered pool %s", backendName, pol.Name)
		pol.DockId = dockId
		if rdn, err := pdd.driverMgr.GetDriverWrapper(consts.DriverTypeReplication, backendName); err == nil {
			pol.ReplicationType = model.ReplicationTypeHost
			if rdn.Name == d.Name {
				pol.ReplicationType = model.ReplicationTypeArray
			}
			pol.ReplicationDriverName = rdn.Name
		}
	}
	return pols, nil
}

func (pdd *provisionDockDiscoverer) Discover() error {
	// Clear existing pool info
	pdd.pols = pdd.pols[:0]
	for _, dck := range pdd.dcks {
		pols, err := pdd.discoverPool(dck.BackendName, dck.Id)
		if err != nil {
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

// attachDockDiscoverer is a struct for exposing some operations of attach
// dock service discovery.
type attachDockDiscoverer struct {
	*DockRegister

	dck *model.DockSpec
}

func (add *attachDockDiscoverer) Init() error { return nil }

func (add *attachDockDiscoverer) Discover() error {
	host, err := os.Hostname()
	if err != nil {
		log.Error("When get os hostname:", err)
		return err
	}

	localIqn, err := connector.NewConnector(connector.IscsiDriver).GetInitiatorInfo()
	if err != nil {
		log.Warning("get initiator failed, ", err)
	}

	bindIp := CONF.BindIp
	if bindIp == "" {
		bindIp = connector.GetHostIP()
	}

	fcInitiator, err := connector.NewConnector(connector.FcDriver).GetInitiatorInfo()
	if err != nil {
		log.Warning("get initiator failed, ", err)
	}

	var wwpns []string
	for _, v := range strings.Split(fcInitiator, ",") {
		if strings.Contains(v, "node_name") {
			wwpns = append(wwpns, strings.Split(v, ":")[1])
		}
	}

	segments := strings.Split(CONF.OsdsDock.ApiEndpoint, ":")
	endpointIp := segments[len(segments)-2]
	add.dck = &model.DockSpec{
		BaseModel: &model.BaseModel{
			Id: uuid.NewV5(uuid.NamespaceOID, host+":"+endpointIp).String(),
		},
		Endpoint: CONF.OsdsDock.ApiEndpoint,
		NodeId:   host,
		Type:     model.DockTypeAttacher,
		Metadata: map[string]string{
			"Platform":  runtime.GOARCH,
			"OsType":    runtime.GOOS,
			"HostIp":    bindIp,
			"Initiator": localIqn,
			"WWPNS":     strings.Join(wwpns, ","),
		},
	}
	return nil
}

func (add *attachDockDiscoverer) Report() error {
	return add.Register(add.dck)
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
