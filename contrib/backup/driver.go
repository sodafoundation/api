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

package backup

import (
	"fmt"
	"os"
)

type BackupSpec struct {
	Id       string
	Name     string
	Metadata map[string]string
}

type BackupDriver interface {
	SetUp() error
	Backup(backup *BackupSpec, volumeFile *os.File) error
	Restore(backup *BackupSpec, backupId string, volFile *os.File) error
	Delete(backup *BackupSpec) error
	CleanUp() error
}

type ctorFun func() (BackupDriver, error)

var ctorFunMap = map[string]ctorFun{}

func NewBackup(backupDriverName string) (BackupDriver, error) {
	fun, exist := ctorFunMap[backupDriverName]
	if !exist {
		return nil, fmt.Errorf("specified backup driver does not exist")
	}

	drv, err := fun()
	if err != nil {
		return nil, err
	}
	return drv, nil
}

func RegisterBackupCtor(bType string, fun ctorFun) error {
	if _, exist := ctorFunMap[bType]; exist {
		return fmt.Errorf("backup driver construct function %s already exist", bType)
	}
	ctorFunMap[bType] = fun
	return nil
}

func UnregisterBackupCtor(cType string) {
	if _, exist := ctorFunMap[cType]; !exist {
		return
	}

	delete(ctorFunMap, cType)
	return
}
