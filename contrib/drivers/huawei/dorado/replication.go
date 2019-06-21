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

package dorado

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils"
	"github.com/opensds/opensds/pkg/utils/config"
)

// ReplicationDriver
type ReplicationDriver struct {
	conf *DoradoConfig
	mgr  *ReplicaPairMgr
}

// Setup
func (r *ReplicationDriver) Setup() (err error) {
	// Read huawei dorado config file
	conf := &DoradoConfig{}
	r.conf = conf
	path := config.CONF.OsdsDock.Backends.HuaweiDorado.ConfigPath

	if "" == path {
		path = defaultConfPath
	}
	Parse(conf, path)
	r.mgr, err = NewReplicaPairMgr(conf)
	if err != nil {
		return err
	}
	return nil
}

// Unset
func (r *ReplicationDriver) Unset() error { return nil }

// CreateReplication
func (r *ReplicationDriver) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {
	log.Info("dorado replication start ...")
	//just be invoked on the primary side.
	if !opt.GetIsPrimary() {
		return &model.ReplicationSpec{}, nil
	}
	pLunId := opt.PrimaryReplicationDriverData[KLunId]
	sLunId := opt.SecondaryReplicationDriverData[KLunId]
	replicationPeriod := strconv.FormatInt(opt.ReplicationPeriod*60, 10)

	replicationMode := ReplicaAsyncMode
	if opt.ReplicationMode == model.ReplicationModeSync {
		replicationMode = ReplicaSyncMode
		replicationPeriod = "0"
	}
	resp, err := r.mgr.CreateReplication(pLunId, sLunId, replicationMode, replicationPeriod)
	if err != nil {
		return nil, err
	}

	return &model.ReplicationSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Metadata: resp,
	}, nil
}

func (r *ReplicationDriver) DeleteReplication(opt *pb.DeleteReplicationOpts) error {
	if !opt.GetIsPrimary() {
		return nil
	}
	pairId, ok := opt.GetMetadata()[KPairId]
	var sLunId string
	if opt.SecondaryVolumeId == "" {
		sLunId = opt.SecondaryReplicationDriverData[KLunId]
	}
	if !ok {
		msg := fmt.Sprintf("Can find pair id in metadata")
		log.Errorf(msg)
		return fmt.Errorf(msg)
	}
	return r.mgr.DeleteReplication(pairId, sLunId)
}

func (r *ReplicationDriver) EnableReplication(opt *pb.EnableReplicationOpts) error {
	if !opt.GetIsPrimary() {
		return nil
	}
	pairId, ok := opt.GetMetadata()[KPairId]
	if !ok {
		msg := fmt.Sprintf("Can find pair id in metadata")
		log.Errorf(msg)
		return fmt.Errorf(msg)
	}
	return r.mgr.localDriver.Enable(pairId, true)
}

func (r *ReplicationDriver) DisableReplication(opt *pb.DisableReplicationOpts) error {
	if !opt.GetIsPrimary() {
		return nil
	}
	pairId, ok := opt.GetMetadata()[KPairId]
	if !ok {
		msg := fmt.Sprintf("Can find pair id in metadata")
		log.Errorf(msg)
		return fmt.Errorf(msg)
	}
	return r.mgr.localDriver.Split(pairId)
}

func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	if !opt.GetIsPrimary() {
		return nil
	}
	pairId, ok := opt.GetMetadata()[KPairId]
	if !ok {
		msg := fmt.Sprintf("Can find pair id in metadata")
		log.Errorf(msg)
		return fmt.Errorf(msg)
	}
	if opt.SecondaryBackendId == model.ReplicationDefaultBackendId {
		return r.mgr.Failover(pairId)
	}
	return r.mgr.Failback(pairId)
}

func NewReplicaPairMgr(conf *DoradoConfig) (r *ReplicaPairMgr, err error) {
	r = &ReplicaPairMgr{}
	r.conf = conf

	r.localClient, err = NewClient(&conf.AuthOptions)
	if err != nil {
		return nil, err
	}
	r.localOp = NewPairOperation(r.localClient)
	r.localDriver = NewReplicaCommonDriver(conf, r.localOp)

	r.remoteClient, err = NewClient(&conf.RemoteAuthOpt)
	if err != nil {
		return nil, err
	}
	r.remoteOp = NewPairOperation(r.remoteClient)
	r.remoteDriver = NewReplicaCommonDriver(conf, r.remoteOp)

	return r, nil
}

type ReplicaPairMgr struct {
	localClient  *DoradoClient
	remoteClient *DoradoClient
	localOp      *PairOperation
	remoteOp     *PairOperation
	localDriver  *ReplicaCommonDriver
	remoteDriver *ReplicaCommonDriver
	conf         *DoradoConfig
}

func (r *ReplicaPairMgr) TryGetRemoteWwn() string {
	sys, _ := r.remoteClient.GetArrayInfo()
	return sys.Wwn
}

func (r *ReplicaPairMgr) TryGetRemoteDevByWwn(wwn string) *RemoteDevice {
	devices, _ := r.localClient.ListRemoteDevices()
	for _, d := range *devices {
		if d.Wwn == wwn {
			return &d
		}
	}
	log.Warningln("Not found remote device")
	return nil
}

func (r *ReplicaPairMgr) CheckRemoteAvailable() bool {
	wwn := r.TryGetRemoteWwn()
	if wwn == "" {
		return false
	}
	d := r.TryGetRemoteDevByWwn(wwn)
	if d != nil && d.ArrayType == ArrayTypeReplication &&
		d.HealthStatus == HealthStatusNormal && d.RunningStatus == RunningStatusLinkUp {
		return true
	}
	return false
}

func (r *ReplicaPairMgr) GetRemoteDevInfo() (id, name string) {
	wwn := r.TryGetRemoteWwn()
	if wwn == "" {
		return "", ""
	}
	dev := r.TryGetRemoteDevByWwn(wwn)
	if dev == nil {
		return "", ""
	}
	return dev.Id, dev.Name
}

func (r *ReplicaPairMgr) WaitVolumeOnline(client *DoradoClient, lun *Lun, interval, timeout time.Duration) error {
	if lun.RunningStatus == StatusVolumeReady {
		return nil
	}

	if interval == -1 {
		interval = DefaultReplicaWaitInterval
	}
	if timeout == -1 {
		timeout = DefaultReplicaWaitTimeout
	}

	return utils.WaitForCondition(func() (bool, error) {
		lunInfo, err := client.GetVolume(lun.Id)
		if err != nil {
			log.Errorf("Get lun failed,%v ", err)
			return false, nil
		}
		if lunInfo.RunningStatus == StatusVolumeReady {
			return true, nil
		}
		return false, nil
	}, interval, timeout)
}

func (r *ReplicaPairMgr) DeletePair(id string) error {
	if !r.localClient.CheckPairExist(id) {
		return nil
	}
	if err := r.localDriver.Split(id); err != nil {
		return err
	}

	err := r.localOp.Delete(id)
	return err
}

func (r *ReplicaPairMgr) CreateReplication(localLunId, rmtLunId, replicationMode string, replicaPeriod string) (map[string]string, error) {
	interval := DefaultReplicaWaitInterval
	timeout := DefaultReplicaWaitTimeout
	var respMap = make(map[string]string)

	localLun, err := r.localClient.GetVolume(localLunId)
	if err != nil {
		return nil, err
	}

	err = r.WaitVolumeOnline(r.localClient, localLun, interval, timeout)
	if err != nil {
		return nil, err
	}

	rmtDevId, rmtDevName := r.GetRemoteDevInfo()
	log.Errorf("rmtDevId:%s, rmtDevName:%s", rmtDevId, rmtDevName)
	if rmtDevId == "" || rmtDevName == "" {
		return nil, fmt.Errorf("get remote deivce info failed")
	}

	pair, err := r.localOp.Create(localLun.Id, rmtLunId, rmtDevId, rmtDevName, replicationMode, ReplicaSpeed, replicaPeriod)
	if err != nil {
		return nil, err
	}
	log.Error("start sync ....", pair)
	if err := r.localDriver.Sync(pair.Id, replicationMode == ReplicaSyncMode); err != nil {
		r.DeletePair(pair.Id)
		return nil, err
	}
	respMap[KPairId] = pair.Id
	return respMap, nil
}

func (r *ReplicaPairMgr) DeleteReplication(pairId, rmtLunId string) error {
	if err := r.DeletePair(pairId); err != nil {
		log.Error("Delete pair failed,", err)
		return err
	}
	return nil
}

// Failover volumes back to primary backend.
// The main steps:
// 1. Switch the role of replication pairs.
// 2. Copy the second LUN data back to primary LUN.
// 3. Split replication pairs.
// 4. Switch the role of replication pairs.
// 5. Enable replications.

func (r *ReplicaPairMgr) Failback(pairId string) error {
	r.remoteDriver.Enable(pairId, true)
	r.remoteDriver.WaitReplicaReady(pairId)
	r.localDriver.Enable(pairId, false)
	return nil
}

func (r *ReplicaPairMgr) Failover(pairId string) error {
	return r.remoteDriver.Failover(pairId)
}

func NewReplicaCommonDriver(conf *DoradoConfig, op *PairOperation) *ReplicaCommonDriver {
	return &ReplicaCommonDriver{conf: conf, op: op}
}

type ReplicaCommonDriver struct {
	conf *DoradoConfig
	op   *PairOperation
}

func (r *ReplicaCommonDriver) ProtectSecond(replicaId string) error {
	replica, err := r.op.GetReplicationInfo(replicaId)
	if err != nil {
		return err
	}
	if replica.SecResAccess == ReplicaSecondRo {
		return nil
	}
	r.op.ProtectSecond(replicaId)
	r.WaitSecondAccess(replicaId, ReplicaSecondRo)
	return nil
}

func (r *ReplicaCommonDriver) UnprotectSecond(replicaId string) error {
	replica, err := r.op.GetReplicationInfo(replicaId)
	if err != nil {
		return err
	}
	if replica.SecResAccess == ReplicaSecondRw {
		return nil
	}
	r.op.UnprotectSecond(replicaId)
	r.WaitSecondAccess(replicaId, ReplicaSecondRw)
	return nil
}

func (r *ReplicaCommonDriver) Sync(replicaId string, waitComplete bool) error {
	r.ProtectSecond(replicaId)
	replicaPair, err := r.op.GetReplicationInfo(replicaId)
	if err != nil {
		return err
	}
	expectStatus := []string{
		RunningStatusNormal,
		RunningStatusSync,
		RunningStatusInitialSync,
	}
	if replicaPair.ReplicationMode == ReplicaSyncMode && r.op.isRunningStatus(expectStatus, replicaPair) {
		return nil
	}
	if err := r.op.Sync(replicaId); err != nil {
		return err
	}
	r.WaitExpectState(replicaId, expectStatus, []string{})
	if waitComplete {
		r.WaitReplicaReady(replicaId)
	}
	return nil
}

func (r *ReplicaCommonDriver) Split(replicaId string) error {
	runningStatus := []string{
		RunningStatusSplit,
		RunningStatusInvalid,
		RunningStatusInterrupted,
	}
	replicaPair, err := r.op.GetReplicationInfo(replicaId)
	if err != nil {
		return err
	}
	if r.op.isRunningStatus(runningStatus, replicaPair) {
		return nil
	}
	err = r.op.Split(replicaId)
	if err != nil {
		log.Errorf("Split replication failed, %v", err)
		return err
	}
	err = r.WaitExpectState(replicaId, runningStatus, []string{})
	if err != nil {
		log.Errorf("Split replication failed, %v", err)
		return err
	}
	return nil
}

func (r *ReplicaCommonDriver) Enable(replicaId string, waitSyncComplete bool) error {
	replicaPair, err := r.op.GetReplicationInfo(replicaId)
	if err != nil {
		return err
	}
	if !r.op.isPrimary(replicaPair) {
		r.Switch(replicaId)
	}
	return r.Sync(replicaId, waitSyncComplete)
}

func (r *ReplicaCommonDriver) Switch(replicaId string) error {
	if err := r.Split(replicaId); err != nil {
		return err
	}

	if err := r.UnprotectSecond(replicaId); err != nil {
		return err
	}

	if err := r.op.Switch(replicaId); err != nil {
		return err
	}

	interval := DefaultReplicaWaitInterval
	timeout := DefaultReplicaWaitTimeout
	return utils.WaitForCondition(func() (bool, error) {
		replicaPair, err := r.op.GetReplicationInfo(replicaId)
		if err != nil {
			return false, nil
		}
		return r.op.isPrimary(replicaPair), nil

	}, interval, timeout)
}

func (r *ReplicaCommonDriver) Failover(replicaId string) error {
	replicaPair, err := r.op.GetReplicationInfo(replicaId)
	if err != nil {
		return err
	}
	if r.op.isPrimary(replicaPair) {
		msg := fmt.Sprintf("We should not do switch over on primary array")
		log.Errorf(msg)
		return fmt.Errorf(msg)
	}
	syncStatus := []string{
		RunningStatusSync,
		RunningStatusInitialSync,
	}
	if r.op.isRunningStatus(syncStatus, replicaPair) {
		if err := r.WaitReplicaReady(replicaId); err != nil {
			return err
		}
	}

	if err := r.Split(replicaId); err != nil {
		return err
	}
	err = r.op.UnprotectSecond(replicaId)
	return err
}

func (r *ReplicaCommonDriver) WaitReplicaReady(replicaId string) error {
	log.Info("Wait synchronize complete.")
	runningNormal := []string{
		RunningStatusNormal,
		RunningStatusSynced,
	}
	runningSync := []string{
		RunningStatusSync,
		RunningStatusInitialSync,
	}
	healthNormal := []string{
		ReplicaHealthStatusNormal,
	}
	interval := DefaultReplicaWaitInterval
	timeout := DefaultReplicaWaitTimeout
	return utils.WaitForCondition(func() (bool, error) {
		replicaPair, err := r.op.GetReplicationInfo(replicaId)
		if err != nil {
			return false, nil
		}
		if r.op.isRunningStatus(runningNormal, replicaPair) && r.op.isHealthStatus(healthNormal, replicaPair) {
			return true, nil
		}
		if !r.op.isRunningStatus(runningSync, replicaPair) {
			msg := fmt.Sprintf("wait synchronize failed, running status:%s", replicaPair.RunningStatus)
			return false, fmt.Errorf(msg)
		}
		return false, nil
	}, interval, timeout)
}

func (r *ReplicaCommonDriver) WaitSecondAccess(replicaId string, accessMode string) {
	interval := DefaultReplicaWaitInterval
	timeout := DefaultReplicaWaitTimeout
	utils.WaitForCondition(func() (bool, error) {
		replicaPair, err := r.op.GetReplicationInfo(replicaId)
		if err != nil {
			return false, nil
		}
		return replicaPair.SecResAccess == accessMode, nil
	}, interval, timeout)
}

func (r *ReplicaCommonDriver) WaitExpectState(replicaId string, runningStatus, healthStatus []string) error {
	interval := DefaultReplicaWaitInterval
	timeout := DefaultReplicaWaitTimeout
	return utils.WaitForCondition(func() (bool, error) {
		replicaPair, err := r.op.GetReplicationInfo(replicaId)
		if err != nil {
			return false, nil
		}
		if r.op.isRunningStatus(runningStatus, replicaPair) {
			if len(healthStatus) == 0 || r.op.isHealthStatus(healthStatus, replicaPair) {
				return true, nil
			}
		}
		return false, nil
	}, interval, timeout)
}

func NewPairOperation(client *DoradoClient) *PairOperation {
	return &PairOperation{client: client}
}

type PairOperation struct {
	client *DoradoClient
}

func (p *PairOperation) isPrimary(replicaPair *ReplicationPair) bool {
	return strings.ToLower(replicaPair.IsPrimary) == "true"
}

func (p *PairOperation) isRunningStatus(status []string, replicaPair *ReplicationPair) bool {
	return utils.Contained(replicaPair.RunningStatus, status)
}

func (p *PairOperation) isHealthStatus(status []string, replicaPair *ReplicationPair) bool {
	return utils.Contained(replicaPair.HealthStatus, status)
}

func (p *PairOperation) Create(localLunId, rmtLunId, rmtDevId, rmtDevName,
	replicationMode, speed, period string) (*ReplicationPair, error) {
	params := map[string]interface{}{
		"LOCALRESID":       localLunId,
		"LOCALRESTYPE":     ObjectTypeLun,
		"REMOTEDEVICEID":   rmtDevId,
		"REMOTEDEVICENAME": rmtDevName,
		"REMOTERESID":      rmtLunId,
		"REPLICATIONMODEL": replicationMode,
		// recovery policy. 1: auto, 2: manual
		"RECOVERYPOLICY": "1",
		"SPEED":          speed,
	}

	if replicationMode == ReplicaAsyncMode {
		// Synchronize type values:
		// 1, manual
		// 2, timed wait when synchronization begins
		// 3, timed wait when synchronization ends
		params["SYNCHRONIZETYPE"] = "2"
		params["TIMINGVAL"] = period
	}
	log.Error(params)
	pair, err := p.client.CreatePair(params)
	if err != nil {
		log.Errorf("Create pair failed,%v", err)
		return nil, err
	}
	return pair, nil
}

func (p *PairOperation) Split(id string) error {
	return p.client.SplitPair(id)
}

func (p *PairOperation) Delete(id string) error {
	return p.client.DeletePair(id)
}

func (p *PairOperation) ProtectSecond(id string) error {
	return p.client.SetPairSecondAccess(id, ReplicaSecondRo)
}

func (p *PairOperation) UnprotectSecond(id string) error {
	return p.client.SetPairSecondAccess(id, ReplicaSecondRw)

}

func (p *PairOperation) Sync(id string) error {
	return p.client.SyncPair(id)
}

func (p *PairOperation) Switch(id string) error {
	return p.client.SwitchPair(id)
}

func (p *PairOperation) GetReplicationInfo(id string) (*ReplicationPair, error) {
	return p.client.GetPair(id)
}
