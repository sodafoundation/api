// Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

package fusionstorage

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils/exec"
)

const (
	CmdBin       = "fsc_cli"
	MaxRetryNode = 3
)

type CliError struct {
	Msg  string
	Code string
}

func (c *CliError) Error() string {
	return fmt.Sprintf("msg: %s, code:%s", c.Msg, c.Code)
}

func NewCliErrorBase(msg, code string) *CliError {
	return &CliError{Msg: msg, Code: code}
}

var VolumeNotExist = NewCliErrorBase("volume does not exist", "50150005")
var SnapshotNotExist = NewCliErrorBase("volume does not exist", "50150006")
var PoolNotExist = NewCliErrorBase("pool does not exist", "50151010")

var CliErrorMap = map[string]string{
	"50000001": "DSware error",
	"50150001": "Receive a duplicate request",
	"50150002": "Command type is not supported",
	"50150003": "Command format is error",
	"50150004": "Lost contact with major VBS",
	"50150005": "Volume does not exist",
	"50150006": "Snapshot does not exist",
	"50150007": "Volume already exists or name exists or name duplicates with a snapshot name",
	"50150008": "The snapshot has already existed",
	"50150009": "VBS space is not enough",
	"50150010": "The node type is error",
	"50150011": "Volume and snapshot number is beyond max",
	"50150012": "VBS is not ready",
	"50150013": "The ref num of node is not 0",
	"50150014": "The volume is not in the pre-deletion state.",
	"50150015": "The storage resource pool is faulty",
	"50150016": "VBS handle queue busy",
	"50150017": "VBS handle request timeout",
	"50150020": "VBS metablock is locked",
	"50150021": "VBS pool dose not exist",
	"50150022": "VBS is not ok",
	"50150023": "VBS pool is not ok",
	"50150024": "VBS dose not exist",
	"50150064": "VBS load SCSI-3 lock pr meta failed",
	"50150100": "The disaster recovery relationship exists",
	"50150101": "The DR relationship does not exist",
	"50150102": "Volume has existed mirror",
	"50150103": "The volume does not have a mirror",
	"50150104": "Incorrect volume status",
	"50150105": "The mirror volume already exists",
}

func NewCliError(code string) error {
	switch code {
	case VolumeNotExist.Code:
		return VolumeNotExist
	case SnapshotNotExist.Code:
		return SnapshotNotExist
	case PoolNotExist.Code:
		return PoolNotExist
	default:
		if msg, ok := CliErrorMap[code]; ok {
			return NewCliErrorBase(msg, code)
		}
		return NewCliErrorBase("CLI execute error", code)
	}
}

type Cli struct {
	// FusionStorage manger ip
	fmIp string
	// FusionStorage agent ip list
	fsaIp []string
	// Command executer
	BaseExecuter exec.Executer
	// Command Root exectuer
	RootExecuter exec.Executer
}

var once sync.Once

// for unit testing
var baseExecuter = exec.NewBaseExecuter()
var rootExecuter = exec.NewRootExecuter()

func NewCli(fmIp string, fsaIP []string) (*Cli, error) {
	if len(fmIp) == 0 || len(fsaIP) == 0 {
		return nil, fmt.Errorf("new cli failed, FM ip or FSA ip can not be set to empty")
	}

	c := &Cli{
		fmIp:  fmIp,
		fsaIp: fsaIP,
	}
	c.BaseExecuter = baseExecuter
	c.RootExecuter = rootExecuter

	//  only execute once
	var err error
	once.Do(func() {
		err = c.StartServer()
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Cli) StartServer() error {
	_, err := c.RootExecuter.Run(CmdBin, "--op", "startServer")
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	log.Info("FSC CLI server start successfully")
	return nil
}

func (c *Cli) doRunCmd(args ...string) (string, error) {

	fsaIp := c.fsaIp
	for i := range fsaIp {
		j := rand.Intn(i + 1)
		fsaIp[i], fsaIp[j] = fsaIp[j], fsaIp[i]
	}
	if len(fsaIp) > MaxRetryNode {
		fsaIp = fsaIp[:3]
	}

	var err error
	var out string
	args = append(args, "--manage_ip", c.fmIp, "--ip", "")
	for key, ip := range fsaIp {
		args[len(args)-1] = ip
		out, err = c.RootExecuter.Run(CmdBin, args...)
		if _, ok := err.(*CliError); ok {
			return out, err
		}
		if err != nil {
			log.Errorf("Run command failed:%v ,retry %d times", err, key+1)
			continue
		}

		return out, err
	}

	return "", err
}

func (c *Cli) RunCmd(args ...string) (string, error) {
	cmdOut, err := c.doRunCmd(args...)
	if err != nil {
		return "", err
	}

	var result string
	const resultPrefix = "result="
	lines := strings.Split(strings.TrimSpace(cmdOut), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, resultPrefix) {
			result = line[len(resultPrefix):]
		}
	}
	if result != "0" {
		return "", NewCliError(result)
	}
	if len(lines) == 1 {
		// nothing to return
		return "", nil
	}
	return strings.Join(lines[:len(lines)-1], "\n"), nil
}

func (c *Cli) QueryAllPoolInfo() ([]*PoolResp, error) {
	args := []string{
		"--op", "queryAllPoolInfo",
	}
	out, err := c.RunCmd(args...)
	if err != nil {
		log.Errorf("Create volume failed: %v", err)
		return nil, err
	}
	pools := []*PoolResp{}
	Unmarshal([]byte(out), &pools)
	return pools, nil
}

func (c *Cli) CreateVolume(name string, size int64, isThin bool, poolId string, encryptOpt *EncryptOpts) error {
	thinFlag := "1"
	if !isThin {
		thinFlag = "0"
	}
	args := []string{
		"--op", "createVolume",
		"--volName", name,
		"--poolId", poolId,
		"--volSize", strconv.FormatInt(size, 10),
		"--thinFlag", thinFlag,
	}
	if encryptOpt != nil {
		encryptArgs := []string{
			"--encrypted", "1",
			"--cmkId", encryptOpt.cmkId,
			"--authCredentials", encryptOpt.authToken,
		}
		args = append(args, encryptArgs...)
	}
	if _, err := c.RunCmd(args...); err != nil {
		log.Errorf("Create volume failed: %v", err)
		return err
	}
	return nil
}

func (c *Cli) DeleteVolume(name string) error {
	args := []string{
		"--op", "deleteVolume",
		"--volName", name,
	}
	_, err := c.RunCmd(args...)
	if err == VolumeNotExist {
		return nil
	}
	return err
}

func (c *Cli) ExtendVolume(name string, newSize int64) error {
	args := []string{
		"--op", "expandVolume",
		"--volName", name,
		"--volSize", strconv.FormatInt(newSize, 10),
	}
	_, err := c.RunCmd(args...)
	return err
}

func (c *Cli) CreateVolumeFromSnapshot(name string, size int64, snapName string) error {
	args := []string{
		"--op", "createVolumeFromSnap",
		"--volName", name,
		"--volSize", strconv.FormatInt(size, 10),
		"--snapNameSrc", snapName,
	}
	_, err := c.RunCmd(args...)
	return err
}

func (c *Cli) QueryVolume(name string) error {
	args := []string{
		"--op", "queryVolume",
		"--volName", name,
	}
	out, err := c.RunCmd(args...)
	if err != nil {
		return err
	}

	snap := &SnapshotResp{}
	return Unmarshal([]byte(out), snap)
}

func (c *Cli) CreateSnapshot(name string, volName string, smartFlag bool) error {
	smartFlagStr := "0"
	if smartFlag {
		smartFlagStr = "1"
	}

	args := []string{
		"--op", "createSnapshot",
		"--volName", volName,
		"--snapName", name,
		"--smartFlag", smartFlagStr,
	}
	_, err := c.RunCmd(args...)
	return err
}

func (c *Cli) QuerySnapshot(name string) error {
	args := []string{
		"--op", "querySnapshot",
		"--snapName", name,
	}
	out, err := c.RunCmd(args...)
	if err != nil {
		return err
	}

	snap := &SnapshotResp{}
	return Unmarshal([]byte(out), snap)
}

func (c *Cli) DeleteSnapshot(name string) error {

	args := []string{
		"--op", "deleteSnapshot",
		"--snapName", name,
	}
	_, err := c.RunCmd(args...)
	if err == SnapshotNotExist {
		return nil
	}
	return err
}
