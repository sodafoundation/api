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

package ceph

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	"golang.org/x/crypto/ssh"
)

type ReplicationDriver struct {
	conf *Config
}

type Replication struct {
	IPaddresshost string `yaml:"hostip,omitempty"`
	Ipaddresspeer string `yaml:"peerip,omitempty"`
	Username      string `yaml:"username,omitempty"`
	Password      string `yaml:"password,omitempty"`
	HostdailIP    string `yaml:"hostdailip,omitempty"`
	PeerdailIP    string `yaml:"peerdailip,omitempty"`
}

type Config struct {
	ConfigFile  string `yaml:"configFile,omitempty"`
	Replication `yaml:"replication"`
}

type Peerinfo struct {
	Mode  string `json:"mode"`
	Peers []struct {
		UUID         json.Number `json:"uuid"`
		Cluster_Name string      `json:"cluster_name"`
		Client_Name  string      `json:"client_name"`
	} `json:"peers"`
}

// Setup
func (r *ReplicationDriver) Setup() error {

	r.conf = &Config{ConfigFile: "/etc/ceph/ceph.conf"}
	p := config.CONF.OsdsDock.Backends.Ceph.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	_, err := Parse(r.conf, p)
	if err != nil {
		logs.Error(err)
	}

	return nil
}

// Unset
func (r *ReplicationDriver) Unset() error { return nil }

// Create replication

func (r *ReplicationDriver) CreateReplication(opt *pb.CreateReplicationOpts) (*model.ReplicationSpec, error) {

	volumename := opensdsPrefix + opt.PrimaryVolumeId

	cephconfig := &ssh.ClientConfig{
		User: r.conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cephclient, err := ssh.Dial("tcp", r.conf.HostdailIP, cephconfig)
	if err != nil {
		log.Error("Failed to dial: " + err.Error())

	}

	backupconfig := &ssh.ClientConfig{
		User: r.conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	backupclient, error := ssh.Dial("tcp", r.conf.PeerdailIP, backupconfig)

	if error != nil {
		log.Error("Failed to dial: " + error.Error())

	}

	cephenablesession, err := cephclient.NewSession()
	if err != nil {
	}

	cephenablesession.Run("rbd  feature enable rbd/" + volumename + " exclusive-lock,journaling")
	cephenablesession.Close()

	cephrbdinstall, err := cephclient.NewSession()
	if err := cephrbdinstall.Run("apt install -y rbd-mirror"); err != nil {
		log.Error("Failed to run: " + err.Error())

	}
	cephrbdinstall.Close()

	backupsession6, err := backupclient.NewSession()

	if err := backupsession6.Run("apt install -y rbd-mirror"); err != nil {
		log.Error("Failed to run: " + err.Error())

	}
	backupsession6.Close()

	cephauthsession, err := cephclient.NewSession()
	if err != nil {
		log.Error("Failed to create session: " + err.Error())

	}

	if err := cephauthsession.Run("ceph auth get-or-create client.ceph mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=rbd' -o /etc/ceph/ceph.client.ceph.keyring"); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	cephauthsession.Close()

	backupauthsession, error := backupclient.NewSession()
	if error != nil {
		log.Error("Failed to create session: " + err.Error())

	}

	if err := backupauthsession.Run("ceph --cluster remote auth get-or-create client.remote mon 'allow r' osd 'allow class-read object_prefix rbd_children, allow rwx pool=rbd' -o /etc/ceph/remote.client.remote.keyring"); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	backupauthsession.Close()

	cephenablesession, err = cephclient.NewSession()
	if err != nil {

	}

	if err := cephenablesession.Run("rbd mirror pool enable rbd image"); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	cephenablesession.Close()

	backupenablesession, error := backupclient.NewSession()
	if error != nil {
	}

	if err := backupenablesession.Run("rbd --cluster remote mirror pool enable rbd image"); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	backupenablesession.Close()

	cephscpsession, err := cephclient.NewSession()
	if err != nil {
		log.Error("Failed to create session: " + err.Error())

	}
	cmd := "scp /etc/ceph/ceph.client.ceph.keyring /etc/ceph/ceph.conf root@" + r.conf.Ipaddresspeer + ":/etc/ceph/"
	if err := cephscpsession.Run(cmd); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	cephscpsession.Close()

	backupscpsession, error := backupclient.NewSession()
	if error != nil {
	}
	peercmd := "scp /etc/ceph/remote.client.remote.keyring /etc/ceph/remote.conf root@" + r.conf.IPaddresshost + ":/etc/ceph/"

	if err := backupscpsession.Run(peercmd); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	backupscpsession.Close()

	cephrbdmirrorsession, err := cephclient.NewSession()
	if err != nil {
		log.Error("Failed to create session: " + err.Error())

	}

	if err := cephrbdmirrorsession.Run("systemctl start ceph-rbd-mirror@ceph"); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	cephrbdmirrorsession.Close()

	backuprbdmirrorsession, error := backupclient.NewSession()
	if error != nil {
	}

	if err := backuprbdmirrorsession.Run("systemctl start ceph-rbd-mirror@remote"); err != nil {
		log.Error("Failed to run: " + err.Error())
	}
	backuprbdmirrorsession.Close()

	cephpeeraddsession, err := cephclient.NewSession()

	err = cephpeeraddsession.Run("rbd mirror pool peer add rbd client.remote@remote")
	if err != nil {
	}

	cephpeeraddsession.Close()

	backuppeersession, err := backupclient.NewSession()

	if err := backuppeersession.Run("rbd --cluster remote mirror pool peer add rbd client.ceph@ceph"); err != nil {
	}
	backuppeersession.Close()
	PrimaryVolumeId := opensdsPrefix + opt.PrimaryVolumeId
	SecondaryVolumeId := opensdsPrefix + opt.SecondaryVolumeId
	PoolId := opt.PoolId
	AvailbilityZone := opt.AvailabilityZone
	profileid := opt.ProfileId

	additionalCephData := map[string]string{
		"PrimaryIP": r.conf.IPaddresshost,
	}

	additionalBackupData := map[string]string{
		"RemoteIP": r.conf.Ipaddresspeer,
	}

	return &model.ReplicationSpec{
		PrimaryVolumeId:                PrimaryVolumeId,
		SecondaryVolumeId:              SecondaryVolumeId,
		PoolId:                         PoolId,
		ProfileId:                      profileid,
		AvailabilityZone:               AvailbilityZone,
		PrimaryReplicationDriverData:   additionalCephData,
		SecondaryReplicationDriverData: additionalBackupData,
	}, nil

}

// Delete replication
func (r *ReplicationDriver) DeleteReplication(opt *pb.DeleteReplicationOpts) error {

	cephconfig := &ssh.ClientConfig{
		User: r.conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cephclient, err := ssh.Dial("tcp", r.conf.HostdailIP, cephconfig)
	if err != nil {
		log.Error("Failed to dial: " + err.Error())

	}
	backupconfig := &ssh.ClientConfig{
		User: r.conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	backupclient, error := ssh.Dial("tcp", r.conf.PeerdailIP, backupconfig)

	if error != nil {
		log.Error("Failed to dial: " + error.Error())

	}

	cephsession2, err := cephclient.NewSession()
	var peerremove bytes.Buffer
	cephsession2.Stdout = &peerremove
	cephsession2.Run("rbd mirror pool info --format json")
	cephsession2.Close()

	byteValue := []byte(peerremove.String())
	peerinfo := &Peerinfo{}
	if err := json.Unmarshal(byteValue, peerinfo); err != nil {
		logs.Error("unmarshal error: %v", err)
	}

	uuid := peerinfo.Peers[0].UUID.String()

	ceph2, err := cephclient.NewSession()
	if err != nil {
	}
	cmd := "rbd mirror pool peer remove rbd " + uuid
	ceph2.Run(cmd)
	ceph2.Close()

	backupsession2, err := backupclient.NewSession()
	var peerremoveremote bytes.Buffer
	backupsession2.Stdout = &peerremoveremote
	backupsession2.Run("rbd --cluster remote mirror pool info --format json")
	cephsession2.Close()
	remotebytevalue := []byte(peerremoveremote.String())
	if err := json.Unmarshal(remotebytevalue, peerinfo); err != nil {
		logs.Error("unmarshal error: %v", err)
	}
	remoteuuid := peerinfo.Peers[0].UUID.String()

	backupsession3, err := backupclient.NewSession()
	backupcmd := "rbd mirror pool peer remove rbd " + remoteuuid + " --cluster remote"
	backupsession3.Run(backupcmd)
	backupsession3.Close()

	cephsession, err := cephclient.NewSession()

	if err := cephsession.Run("systemctl stop ceph-rbd-mirror@ceph"); err != nil {
		logs.Error(err)
	}
	cephsession.Close()

	backupsession, err := backupclient.NewSession()

	if err := backupsession.Run("systemctl stop ceph-rbd-mirror@remote"); err != nil {
		logs.Error(err)
	}

	backupsession.Close()

	return nil
}

// Start Replication
func (r *ReplicationDriver) EnableReplication(opt *pb.EnableReplicationOpts) error {

	cephconfig := &ssh.ClientConfig{
		User: r.conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cephclient, err := ssh.Dial("tcp", r.conf.HostdailIP, cephconfig)
	if err != nil {
		log.Error("Failed to dial: " + err.Error())

	}

	cephsession, err := cephclient.NewSession()
	if err != nil {
		log.Error("Failed to create session: " + err.Error())

	}

	volumename := opensdsPrefix + opt.PrimaryVolumeId

	cmd := "rbd mirror image enable rbd/" + volumename + " --pool rbd"

	if err := cephsession.Run(cmd); err != nil {
		log.Error("Failed to run: " + err.Error())

	}

	cephsession.Close()

	return nil
}

// Stop Replication
func (r *ReplicationDriver) DisableReplication(opt *pb.DisableReplicationOpts) error {

	cephconfig := &ssh.ClientConfig{
		User: r.conf.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.conf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	cephclient, err := ssh.Dial("tcp", r.conf.HostdailIP, cephconfig)
	if err != nil {
		log.Error("Failed to dial: " + err.Error())

	}
	cephdisablesession, err := cephclient.NewSession()
	if err != nil {
	}
	volumename := opensdsPrefix + opt.PrimaryVolumeId

	cmd := "rbd mirror image disable rbd/" + volumename + " --pool rbd"
	if err := cephdisablesession.Run(cmd); err != nil {
		log.Error("Failed to run: " + err.Error())

	}

	cephdisablesession.Close()

	return nil
}

func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	return nil
}
