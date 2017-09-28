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

package provider

import (
	_ "encoding/json"
	"log"
	"sync"

	model "github.com/opensds/opensds/contrib/swordfish/proto"
)

var storageServiceLinks = "/redfish/v1/StorageServices(1)/Links/"

type Provider struct {
	rwMutex       sync.RWMutex
	cosCollection *model.ClassOfServiceCollection

	dataProtectionLoS *model.DataProtectionLoSCapabilities
	dataSecurityLoS   *model.DataSecurityLoSCapabilities
	dataStorageLoS    *model.DataStorageLoSCapabilities
	ioConnectivityLoS *model.IOConnectivityLoSCapabilities
	ioPerfermanceLoS  *model.IOPerformanceLoSCapabilities
}

func NewProvider() *Provider {
	var (
		cosc   model.ClassOfServiceCollection
		dpLos  model.DataProtectionLoSCapabilities
		dsLos  model.DataSecurityLoSCapabilities
		dstLos model.DataStorageLoSCapabilities
		icLos  model.IOConnectivityLoSCapabilities
		ipLos  model.IOPerformanceLoSCapabilities
	)

	return &Provider{
		cosCollection:     &cosc,
		dataProtectionLoS: &dpLos,
		dataSecurityLoS:   &dsLos,
		dataStorageLoS:    &dstLos,
		ioConnectivityLoS: &icLos,
		ioPerfermanceLoS:  &ipLos,
	}
}

func (p *Provider) CreateClassOfService(req *model.ClassOfService) (string, error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.cosCollection.Members = append(p.cosCollection.Members, req)

	log.Println(p.cosCollection)

	return storageServiceLinks + "(" + req.GetName() + ")", nil
}

func (p *Provider) GetClassesOfService() (*model.ClassOfServiceCollection, error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	return p.cosCollection, nil
}

func (p *Provider) CreateDataProtectionLineOfService(req *model.DataProtectionLoSCapabilities) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.dataProtectionLoS.SupportedDataProtectionLinesOfService =
		append(p.dataProtectionLoS.SupportedDataProtectionLinesOfService, req.GetSupportedDataProtectionLinesOfService()...)

	log.Println(p.dataProtectionLoS.GetSupportedDataProtectionLinesOfService())

	return nil
}

func (p *Provider) CreateDataSecurityLineOfService(req *model.DataSecurityLoSCapabilities) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.dataSecurityLoS.SupportedDataSecurityLinesOfService =
		append(p.dataSecurityLoS.SupportedDataSecurityLinesOfService, req.GetSupportedDataSecurityLinesOfService()...)

	return nil
}

func (p *Provider) CreateDataStorageLineOfService(req *model.DataStorageLoSCapabilities) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.dataStorageLoS.SupportedDataStorageLinesOfService =
		append(p.dataStorageLoS.SupportedDataStorageLinesOfService, req.GetSupportedDataStorageLinesOfService()...)

	return nil
}

func (p *Provider) CreateIOConnectivityLineOfService(req *model.IOConnectivityLoSCapabilities) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.ioConnectivityLoS.SupportedIOConnectivityLinesOfService =
		append(p.ioConnectivityLoS.SupportedIOConnectivityLinesOfService, req.GetSupportedIOConnectivityLinesOfService()...)

	return nil
}

func (p *Provider) CreateIOPerformanceLineOfService(req *model.IOPerformanceLoSCapabilities) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.ioPerfermanceLoS.SupportedIOPerformanceLinesOfService =
		append(p.ioPerfermanceLoS.SupportedIOPerformanceLinesOfService, req.GetSupportedIOPerformanceLinesOfService()...)

	log.Println(p.ioPerfermanceLoS.GetSupportedIOPerformanceLinesOfService())

	return nil
}

func (p *Provider) GetSupportedDataProtectionLinesOfService() ([]*model.DataProtectionLineOfService, error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	return p.dataProtectionLoS.SupportedDataProtectionLinesOfService, nil
}

func (p *Provider) GetSupportedDataSecurityLinesOfService() ([]*model.DataSecurityLineOfService, error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	return p.dataSecurityLoS.SupportedDataSecurityLinesOfService, nil
}

func (p *Provider) GetSupportedDataStorageLinesOfService() ([]*model.DataStorageLineOfService, error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	return p.dataStorageLoS.SupportedDataStorageLinesOfService, nil
}

func (p *Provider) GetSupportedIOConnectivityLinesOfService() ([]*model.IOConnectivityLineOfService, error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	return p.ioConnectivityLoS.SupportedIOConnectivityLinesOfService, nil
}

func (p *Provider) GetSupportedIOPerformanceLinesOfService() ([]*model.IOPerformanceLineOfService, error) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	return p.ioPerfermanceLoS.SupportedIOPerformanceLinesOfService, nil
}
