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
This module implements the entry into operations of storageDock module.

*/

package app

import (
	"errors"
	log "github.com/golang/glog"

	"github.com/opensds/opensds/pkg/db"
	"github.com/opensds/opensds/pkg/dock"
	"github.com/opensds/opensds/pkg/utils"
)

func ResourceDiscovery() error {
	dcks, err := ListDocks()
	if err != nil {
		log.Error("When list docks:", err)
		return err
	}
	if len(*dcks) == 0 {
		return errors.New("The dock resource is empty")
	}

	for _, dck := range *dcks {
		// If dock uuid is null, generate it randomly.
		if dck.GetId() == "" {
			if ok := utils.NewSetter().SetUuid(dck); ok != nil {
				log.Error("When set dock uuid:", ok)
				return ok
			}
		}

		// Set dock created time.
		if ok := utils.NewSetter().SetCreatedTimeStamp(dck); ok != nil {
			log.Error("When set dock created time:", ok)
			return ok
		}

		// Call db module to create dock resource.
		if _, err = db.C.CreateDock(&dck); err != nil {
			log.Error("When create dock %s in db: %v\n", dck.Id, err)
			return err
		}

		if err := poolResourceDiscovery(dck.GetDriverName()); err != nil {
			log.Error("When discovery pool in dock %s: %v\n", dck.Id, err)
			return err
		}
	}
	return nil
}

func poolResourceDiscovery(driver string) error {
	pols, err := dock.NewDockHub(driver).ListPools()
	if err != nil {
		log.Error("When list pools:", err)
		return err
	}
	if len(*pols) == 0 {
		return errors.New("The pool resource is empty")
	}

	for _, pol := range *pols {
		// If pool uuid is null, generate it randomly.
		if pol.GetId() == "" {
			if ok := utils.NewSetter().SetUuid(pol); ok != nil {
				log.Error("When set pool uuid:", ok)
				return ok
			}
		}

		// Set pool created time.
		if ok := utils.NewSetter().SetCreatedTimeStamp(pol); ok != nil {
			log.Error("When set pool created time:", ok)
			return ok
		}

		// Call db module to create pool resource.
		if _, err = db.C.CreatePool(&pol); err != nil {
			log.Error("When create pool %s in db: %v\n", pol.Id, err)
			return err
		}
	}
	return nil
}
