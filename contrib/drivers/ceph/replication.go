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
	"github.com/astaxie/beego/logs"
	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	"github.com/opensds/opensds/pkg/utils/config"
	"golang.org/x/crypto/ssh"
	"time"
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
	volumename := opensdsPrefix + opt.PrimaryVolumeId

	cephdemotesession, err := cephclient.NewSession()
	if err != nil {
	}

	demotecmd := "rbd mirror image demote rbd/" + volumename
	if err := cephdemotesession.Run(demotecmd); err != nil {
		log.Error("Failed to run: " + err.Error())
	}

	cephdemotesession.Close()

	time.Sleep(10 * time.Second)

	backuppromotesession, err := backupclient.NewSession()
	if err != nil {
	}
	promotecmd := "rbd mirror image promote rbd/" + volumename + " --cluster remote"
	if err := backuppromotesession.Run(promotecmd); err != nil {
		log.Error("Failed to run: " + err.Error())

	}

	backuppromotesession.Close()

	time.Sleep(10 * time.Second)

	cephdisablesession, err := cephclient.NewSession()
	if err != nil {
	}

	cmd := "rbd mirror image disable rbd/" + volumename + " --force"
	if err := cephdisablesession.Run(cmd); err != nil {
		log.Error("Failed to run: " + err.Error())

	}

	cephdisablesession.Close()
	time.Sleep(10 * time.Second)

	backupsnapshot, err := backupclient.NewSession()
	snapcmd := "rbd snap create rbd/" + volumename + "@" + volumename + " --cluster remote"
	if err := backupsnapshot.Run(snapcmd); err != nil {
		log.Error("Failed to run: " + err.Error())
	}

	backupsnapshot.Close()
	time.Sleep(10 * time.Second)

	return nil
}

func (r *ReplicationDriver) FailoverReplication(opt *pb.FailoverReplicationOpts) error {
	return nil
}
