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

package drivers

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/contrib/drivers/factory"
	. "github.com/opensds/opensds/contrib/drivers/utils/constants"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
)

func NewDriverWrapper(driver factory.Driver, name string) *DriverWrapper {
	dw := &DriverWrapper{Driver: driver, Name: name}
	dw.SetInitialized(false)
	return dw
}

type DriverWrapper struct {
	// driver name
	Name string
	// Base Driver
	factory.Driver
	// whether driver initialized used by atomic
	initialized int32
}

func (b *DriverWrapper) SetInitialized(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.initialized), int32(i))
}

func (b *DriverWrapper) GetInitialized() bool {
	if atomic.LoadInt32(&(b.initialized)) != 0 {
		return true
	}
	return false
}

func NewBackend(properties *config.BackendProperties) *Backend {
	return &Backend{properties: properties}
}

type Backend struct {
	provisionDriver   *DriverWrapper
	replicationDriver *DriverWrapper
	metricDriver      *DriverWrapper
	properties        *config.BackendProperties
}

func (b *Backend) SetProvisionDriver(d *DriverWrapper) {
	b.provisionDriver = d
}

func (b *Backend) SetReplicationDriver(d *DriverWrapper) {
	b.replicationDriver = d
}

func (b *Backend) SetMetricDriver(d *DriverWrapper) {
	b.metricDriver = d
}

func (b *Backend) GetProvisionDriver() *DriverWrapper {
	return b.provisionDriver
}

func (b *Backend) GetReplicationDriver() *DriverWrapper {
	return b.replicationDriver
}
func (b *Backend) GetMetricDriver() *DriverWrapper {
	return b.metricDriver
}

func (b *Backend) GetDriverByType(t DriverType) *DriverWrapper {
	// please don't change the order
	return [DriverTypeNum]*DriverWrapper{b.provisionDriver, b.replicationDriver, b.metricDriver}[t]
}

func NewDriverManager() *DriverManager {
	return &DriverManager{
		backendMap:    map[string]*Backend{},
		driverFactory: NewDriverFactory(),
	}
}

type DriverManager struct {
	// drivers store all enabled driver
	backendMap map[string]*Backend
	// Driver lock for driver initialize
	driverLock sync.RWMutex
	//
	driverFactory factory.DriverFactory
}

func (dm *DriverManager) SetupDriver(d *DriverWrapper) {
	retryTimes := 1
	_ = utils.WaitForCondition(func() (bool, error) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("driver '%s' setup run time panic: %v, retry %d time(s)", d.Name, err, retryTimes)
			}
			retryTimes++
		}()
		err := d.Setup()
		if err != nil {
			log.Errorf("driver '%s' setup failed: %v, retry %d time(s)", d.Name, err, retryTimes)
			return false, nil // setup failed, retry it.
		}
		// set driver status
		d.SetInitialized(true)
		log.Infof("driver '%s' setup successfully", d.Name)
		return true, nil
	}, 10*time.Second, 600*time.Second)
}

func (dm *DriverManager) doLoadDriver(backendName string, bp config.BackendProperties) error {
	b := NewBackend(&bp)
	// provision driver
	driver, err := dm.driverFactory.GetDriver(StorageType(bp.StorageType), DriverTypeProvision, bp)
	if err != nil {
		log.Errorf("get provision driver '%s'failed: %v", bp.DriverName, err)
		return err
	}
	b.SetProvisionDriver(NewDriverWrapper(driver, bp.DriverName))
	go dm.SetupDriver(b.GetProvisionDriver())

	// replication driver
	if len(bp.ReplicationDriverName) != 0 {
		driver, err := dm.driverFactory.GetDriver(StorageType(bp.StorageType), DriverTypeReplication, bp)
		if err != nil {
			log.Errorf("get replication driver '%s'failed: %v", bp.ReplicationDriverName, err)
			return err
		}
		b.SetReplicationDriver(NewDriverWrapper(driver, bp.ReplicationDriverName))
		go dm.SetupDriver(b.GetProvisionDriver())
	}

	// metric driver
	if len(bp.MetricDriverName) != 0 {
		driver, err := dm.driverFactory.GetDriver(StorageType(bp.StorageType), DriverTypeMetric, bp)
		if err != nil {
			log.Errorf("get metric driver '%s'failed: %v", bp.MetricDriverName, err)
			return err
		}
		b.SetReplicationDriver(NewDriverWrapper(driver, bp.MetricDriverName))
		go dm.SetupDriver(b.GetMetricDriver())
	}

	dm.backendMap[backendName] = b
	return nil
}

func (dm *DriverManager) LoadDriver() error {
	// volume
	for n, b := range config.GetBackendMap() {
		if err := dm.doLoadDriver(n, b); err != nil {
			log.Errorf("load block driver %s failed:%s", b.DriverName, err)
			return err
		}
	}
	return nil
}

func (dm *DriverManager) GetDriver(driverType DriverType, backendName string) (factory.Driver, error) {
	dw, err := dm.GetDriverWrapper(driverType, backendName)
	if err != nil {
		return nil, err
	}
	return dw.Driver, nil
}

func (dm *DriverManager) GetDriverWrapper(driverType DriverType, backendName string) (*DriverWrapper, error) {

	log.Infof("%v, %v, %v", driverType, backendName)
	backend, ok := dm.backendMap[backendName]
	if !ok {
		return nil, fmt.Errorf("specified backend '%s' doesn't exist", backendName)
	}
	dw := backend.GetDriverByType(driverType)
	if !dw.GetInitialized() {
		return nil, fmt.Errorf("specified backend '%s' is not initialized", backendName)
	}
	return dw, nil
}

func (dm *DriverManager) GetBackend(name string) (*Backend, error) {
	return dm.backendMap[name], nil
}
