// Copyright 2019 The OpenSDS Authors.
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

package eternus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	"github.com/opensds/opensds/pkg/utils/pwd"
	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
)

// EternusClient :
type EternusClient struct {
	user        string
	password    string
	endpoint    string
	stdin       io.WriteCloser
	stdout      io.Reader
	stderr      io.Reader
	cliConfPath string
}

func NewClient(opt *AuthOptions) (*EternusClient, error) {
	var pwdCiphertext = opt.Password

	if opt.EnableEncrypted {
		// Decrypte the password
		pwdTool := pwd.NewPwdEncrypter(opt.PwdEncrypter)
		password, err := pwdTool.Decrypter(pwdCiphertext)
		if err != nil {
			return nil, err
		}
		pwdCiphertext = password
	}

	c := &EternusClient{
		user:        opt.Username,
		password:    pwdCiphertext,
		endpoint:    opt.Endpoint,
		cliConfPath: defaultCliConfPath,
	}
	return c, nil
}

func NewClientForAdmin(opt *AuthOptions) (*EternusClient, error) {
	var pwdCiphertext = opt.AdminPassword

	if opt.EnableEncrypted {
		// Decrypte the password
		pwdTool := pwd.NewPwdEncrypter(opt.PwdEncrypter)
		password, err := pwdTool.Decrypter(pwdCiphertext)
		if err != nil {
			return nil, err
		}
		pwdCiphertext = password
	}

	c := &EternusClient{
		user:        opt.AdminUsername,
		password:    pwdCiphertext,
		endpoint:    opt.Endpoint,
		cliConfPath: defaultCliConfPath,
	}
	return c, nil
}

func (c *EternusClient) Destroy() error {
	_, err := c.stdin.Write([]byte("exit\n"))
	return err
}

func (c *EternusClient) setConfig() *ssh.ClientConfig {
	var defconfig ssh.Config
	defconfig.SetDefaults()
	cipherOrder := defconfig.Ciphers

	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	config.Ciphers = append(cipherOrder, "3des-cbc")

	return config
}

func (c *EternusClient) createSession(config *ssh.ClientConfig) (*ssh.Session, error) {
	server := c.endpoint
	server = server + ":" + SSHPort

	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Error("failed to dial: " + err.Error())
		return nil, err
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Error("failed to create session: " + err.Error())
		return nil, err
	}

	return session, nil
}

func (c *EternusClient) doRequest(cmd string, param map[string]string) (bytes.Buffer, error) {
	// create command option
	cmdOption := ""
	if param != nil {
		for k, v := range param {
			cmdOption += fmt.Sprintf(" -%s %s ", k, v)
		}
	}
	// execute command
	log.Infof("execute cli. cmd = %s, option = %s", cmd, cmdOption)
	c.stdin.Write([]byte(cmd + cmdOption + "\n"))
	var result bytes.Buffer
	buff := make([]byte, 65535)

	// get command output
	for {
		n, err := c.stdout.Read(buff)
		if err != io.EOF && err != nil {
			return result, err
		}
		n, err = result.Write(buff[:n])

		// ignore first '\r\nCLI>'
		if result.String() == "\r\nCLI> " {
			continue
		}
		// if error occured or suffix is 'CLI> ', break the loop
		if err == io.EOF || strings.HasSuffix(result.String(), "CLI> ") {
			break
		}
	}
	c.stdin.Write([]byte("\n"))
	return result, nil
}

func (c *EternusClient) request(cmd string, param map[string]string) ([]map[string]string, error) {
	var b bytes.Buffer
	var err error
	var resultHash []map[string]string
	success := false
	for i := 0; i < 2; i++ {
		b, err = c.doRequest(cmd, param)
		if err == nil {
			resultArray := c.convStringArray(cmd, &b)
			resultHash, err = c.parseResult(cmd, resultArray)
			if err != nil {
				log.Errorf("failed to execute cmd. err = %s, output = %v", err.Error(), resultArray)
				continue
			}
			success = true
			break
		} else {
			log.Errorf("cmd:%s %s\n param:%v", cmd, c.endpoint, param)
		}
		log.Errorf("request %d times error:%v", i+1, err)
	}
	if success == false {
		return resultHash, err
	}
	return resultHash, nil
}

// requestForadmin is temporary function for snapshot
// Do not use the function except snapshot
func (c *EternusClient) requestForAdmin(cmd string, param map[string]string) (bytes.Buffer, error) {
	var b bytes.Buffer
	var err error
	success := false
	for i := 0; i < 2; i++ {
		b, err = c.doRequest(cmd, param)
		if err == nil {
			success = true
			break
		} else {
			log.Errorf("cmd:%s %s\n param:%v", cmd, c.endpoint, param)
		}
		log.Errorf("request %d times error:%v", i+1, err)
	}

	if success == false {
		return b, err
	}

	for _, s := range strings.Split(b.String(), "\r\n") {
		// ignore empty line(first elem)
		if s == "" {
			continue
		}
		// ignore echo back string
		if strings.HasPrefix(s, "CLI> "+cmd) {
			continue
		}
		// ignore last line and stop parse
		if s == "CLI> " {
			break
		}
		// check error
		if strings.HasPrefix(s, "Error: ") {
			errMsg := fmt.Sprintf("failed to command output = %s", s)
			log.Error(errMsg)
			return b, errors.New(s)
		}
	}
	return b, nil
}

func (c *EternusClient) login() error {
	config := c.setConfig()
	session, err := c.createSession(config)
	if err != nil {
		log.Error("failed to get session: " + err.Error())
		return err
	}

	c.stdin, err = session.StdinPipe()
	if err != nil {
		log.Error("failed to get StdinPipe: " + err.Error())
		return err
	}

	c.stdout, err = session.StdoutPipe()
	if err != nil {
		log.Error("failed to get StdoutPipe: " + err.Error())
		return err
	}

	c.stderr, err = session.StderrPipe()
	if err != nil {
		log.Error("failed to get StderrPipe: " + err.Error())
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	err = session.RequestPty("xterm", 80, 1024, modes)
	if err != nil {
		log.Error("failed to request pty: " + err.Error())
		return err
	}

	err = session.Shell()
	if err != nil {
		log.Error("failed to get shell: " + err.Error())
		return err
	}
	return nil
}

func (c *EternusClient) parseResult(cmd string, resultArray [][]string) ([]map[string]string, error) {
	// read cli config file
	yamlConfig, err := ioutil.ReadFile(c.cliConfPath)
	if err != nil {
		log.Error("failed to read cli_response.yml: " + err.Error())
		return nil, err
	}
	// parse yaml
	var config map[string]map[string]([]map[string]interface{})
	err = yaml.Unmarshal(yamlConfig, &config)

	// get config for specified cmd
	cmdConfig := config["cli"][strings.Replace(cmd, " ", "_", -1)]

	// parse resultArray
	var ret []map[string]string
	resultIndex := 0
	var dataCount int
	for _, v := range cmdConfig {
		switch v["type"] {
		case "status":
			var status int
			// check for correct response
			if len(resultArray) > resultIndex {
				// get response status
				if len(resultArray[resultIndex]) == 1 {
					status, _ = strconv.Atoi(resultArray[resultIndex][0])
				} else {
					return nil, errors.New("error response. Failed to get status")
				}
			} else {
				return nil, errors.New("error response. Failed to get status")
			}
			// check response status
			if status != 0 {
				if len(resultArray) > (resultIndex+1) &&
					len(resultArray[resultIndex+1]) == 1 {
					errorCode := map[string]string{
						"errorCode": resultArray[resultIndex+1][0],
					}
					ret = append(ret, errorCode)
				}
				return ret, errors.New("error response. Command execute error")
			}
			resultIndex++
		case "count":
			// check for correct response
			if len(resultArray) > resultIndex {
				// get data count
				if len(resultArray[resultIndex]) == 1 {
					tmpCount, _ := strconv.ParseInt(resultArray[resultIndex][0], 16, 64)
					dataCount = int(tmpCount)
				} else {
					fmt.Println(resultArray[resultIndex])
					return nil, errors.New("error response. Failed to get count")
				}
			} else {
				return nil, errors.New("error response. Failed to get count")
			}
			if v["if_zero_skip_all"] == true && dataCount == 0 {
				break
			}
			resultIndex++
		case "data":
			// check drop flag
			if v["drop"] == true {
				resultIndex++
				continue
			}
			// check for correct response
			if len(resultArray) > resultIndex {
				// get single data
				datas := v["data"].([]interface{})
				result, err := c.getData(datas, resultArray[resultIndex])
				if err != nil {
					return nil, err
				}
				ret = append(ret, result)
				resultIndex++
			} else {
				return nil, errors.New("error response. Failed to get data")
			}
		case "multiple_data":
			// get multiple data, data count = dataCount variable
			datas := v["data"].([]interface{})
			for j := 0; j < dataCount; j++ {
				// check drop flag
				if v["drop"] == true {
					resultIndex++
					continue
				}
				if len(resultArray) > resultIndex {
					result, err := c.getData(datas, resultArray[resultIndex])
					if err != nil {
						return nil, err
					}
					ret = append(ret, result)
					resultIndex++
				} else {
					return nil, errors.New("error response. Failed to get multiple_data")
				}
			}
		}
	}
	return ret, nil
}

func (c *EternusClient) getData(datas []interface{}, result []string) (map[string]string, error) {
	data := map[string]string{}
	for i, v := range datas {
		// store each param
		paramName := v.(map[interface{}]interface{})["name"].(string)
		if len(result) > i {
			data[paramName] = result[i]
		} else {
			errMsg := "the response is inconsistent with the response def"
			return nil, errors.New(errMsg)
		}
	}
	return data, nil
}

func (c *EternusClient) convStringArray(cmd string, result *bytes.Buffer) [][]string {
	output := [][]string{}
	for _, s := range strings.Split(result.String(), "\r\n") {
		// ignore empty line(first elem)
		if s == "" {
			continue
		}
		// ignore echo back string
		if strings.HasPrefix(s, "CLI> "+cmd) {
			continue
		}
		// ignore last line and stop parse
		if s == "CLI> " {
			break
		}
		output = append(output, strings.Split(s, "\t"))
	}
	return output
}

// ListStoragePools :
func (c *EternusClient) ListStoragePools() ([]StoragePool, error) {
	var pools []StoragePool
	ret, err := c.request("show thin-pro-pools", nil)
	if err != nil {
		return pools, err
	}
	for _, v := range ret {
		id, _ := strconv.ParseInt(v["tpp_number"], 16, 64)
		// calculate free capacity
		// cut off after the decimal point
		totalCapacity, _ := strconv.ParseInt(v["total_capacity"], 16, 64)
		usedCapacity, _ := strconv.ParseInt(v["used_capacity"], 16, 64)
		totalCapacity = totalCapacity / LBASize
		usedCapacity = usedCapacity / LBASize
		freeCapacity := totalCapacity - usedCapacity

		pool := StoragePool{}
		pool.Id = strconv.FormatInt(id, 10)
		pool.Name = v["tpp_name"]
		pool.TotalCapacity = totalCapacity
		pool.FreeCapacity = freeCapacity
		pools = append(pools, pool)
	}
	return pools, nil
}

// ListAllStoragePools :
func (c *EternusClient) ListAllStoragePools() ([]StoragePool, error) {
	return c.ListStoragePools()
}

// GetVolume :
func (c *EternusClient) GetVolume(lunID string) (Volume, error) {
	param := map[string]string{
		"volume-number": lunID,
	}
	return c.execGetVolume(param)
}

// GetVolumeByName :
func (c *EternusClient) GetVolumeByName(name string) (Volume, error) {
	param := map[string]string{
		"volume-name": name,
	}
	return c.execGetVolume(param)
}

func (c *EternusClient) execGetVolume(param map[string]string) (Volume, error) {
	var vol Volume
	ret, err := c.request("show volumes", param)
	if err != nil {
		log.Error("failed to get volume information: " + err.Error())
		return vol, err
	}
	v := ret[0]
	id, _ := strconv.ParseInt(v["volume_number"], 16, 64)
	poolID, _ := strconv.ParseInt(v["pool_number"], 16, 64)
	size, _ := strconv.ParseInt(v["size"], 16, 64)
	size = size / LBASize

	vol.Id = strconv.FormatInt(id, 10)
	vol.Name = v["volume_name"]
	vol.Size = size
	vol.Status = v["status"]
	vol.PoolName = v["pool_name"]
	vol.PoolId = strconv.FormatInt(poolID, 10)

	return vol, nil
}

// CreateVolume :
func (c *EternusClient) CreateVolume(id string, size int64, desc string,
	poolName string, provPolicy string) (Volume, error) {

	// use hash value because eternus has limitation of name length
	name := GetFnvHash(id)
	sizeGB := fmt.Sprintf("%dgb", size)
	allocation := "thin"
	if provPolicy != "Thin" {
		allocation = "thick"
	}
	param := map[string]string{
		"name":       name,
		"size":       sizeGB,
		"pool-name":  poolName,
		"type":       "tpv",
		"allocation": allocation,
	}
	var vol Volume
	ret, err := c.request(CreateVolume, param)
	if err != nil {
		log.Error("failed to create volume: " + err.Error())
		return vol, err
	}

	v := ret[0]
	convID, _ := strconv.ParseInt(v["volume_number"], 16, 64)
	vol.Id = strconv.FormatInt(convID, 10)
	vol.Name = name
	vol.Size = size

	return vol, nil
}

// DeleteVolume :
func (c *EternusClient) DeleteVolume(volumeNumber string) error {
	param := map[string]string{
		"volume-number": volumeNumber,
	}
	_, err := c.request("delete volume", param)
	if err != nil {
		log.Error("failed to delete volume: " + err.Error())
		return err
	}
	return nil
}

// ExtendVolume :
func (c *EternusClient) ExtendVolume(volumeNumber string, size int64) error {
	param := map[string]string{
		"volume-number": volumeNumber,
		"size":          strconv.FormatInt(size, 10) + "gb",
	}

	_, err := c.request("expand volume", param)
	if err != nil {
		log.Error("failed to expand volume: " + err.Error())
		return err
	}
	return nil
}

// AddIscsiHostWithCheck :
func (c *EternusClient) AddIscsiHostWithCheck(name string, iscsiName string, ipAddr string) (string, bool, error) {
	// check duplicate host. if already exists, retuen exist host id.
	ret, err := c.request("show host-iscsi-names", nil)
	if err != nil {
		log.Error("failed to show host-iscsi-names: " + err.Error())
		return "", false, err
	}
	for _, v := range ret {
		ipStr := ""
		if v["ip_version"] == "00" {
			ipStr = ParseIPv4(v["host_iscsi_ip_address"])
		} else {
			ipStr = ParseIPv6(v["host_iscsi_ip_address"])
		}
		if v["host_iscsi_name"] == iscsiName && EqualIP(ipStr, ipAddr) {
			hostNumber, _ := strconv.ParseInt(v["host_iscsi_number"], 16, 64)
			return strconv.FormatInt(hostNumber, 10), true, nil
		}
	}

	// create new host
	ipVersion := "ipv4"
	if !IsIPv4(ipAddr) {
		ipVersion = "ipv6"
	}
	param := map[string]string{
		"iscsi-name": iscsiName,
		"ip-version": ipVersion,
		"name":       name,
	}
	if ipAddr != "" {
		param["ip"] = ipAddr
	}
	ret, err = c.request("create host-iscsi-name", param)
	if err != nil {
		log.Error("failed to create host-iscsi-name: " + err.Error())
		return "", false, err
	}
	v := ret[0]
	hostNumber, _ := strconv.ParseInt(v["host_number"], 16, 64)
	return strconv.FormatInt(hostNumber, 10), false, nil
}

// DeleteIscsiHost :
func (c *EternusClient) DeleteIscsiHostByName(name string) error {
	param := map[string]string{
		"host-name": name,
	}
	ret, err := c.request("delete host-iscsi-name", param)
	if err != nil {
		if len(ret) == 1 && ret[0]["errorCode"] == NotFound {
			log.Info("target iscsi host already deleted")
			return nil
		}
		log.Error("failed to delete host-iscsi-name: " + err.Error())
		return err
	}
	return nil
}

// GetLunGroupByName :
func (c *EternusClient) GetLunGroupByName(name string) (LunGroup, error) {
	lunGroup := LunGroup{}
	param := map[string]string{
		"lg-name": name,
	}
	ret, err := c.request("show lun-group", param)
	if err != nil {
		log.Error("failed to show lun-group: " + err.Error())
		return lunGroup, err
	}
	lunGroupVolumes := []LunGroupVolume{}
	for _, v := range ret {
		vol := LunGroupVolume{}
		volID, _ := strconv.ParseInt(v["volume_no"], 16, 64)
		hostLunID, _ := strconv.ParseInt(v["lun"], 16, 64)
		tmpSize, _ := strconv.ParseInt(v["total_capacity"], 16, 64)
		size := tmpSize / LBASize

		vol.Id = strconv.FormatInt(volID, 10)
		vol.Name = v["volume_name"]
		vol.RawStatus = v["volume_raw_status"]
		vol.RoundStatus = v["volume_round_status"]
		vol.Size = size
		vol.Uid = v["uid"]
		vol.Lun = strconv.FormatInt(hostLunID, 10)

		lunGroupVolumes = append(lunGroupVolumes, vol)
	}
	lunGroup.Volumes = lunGroupVolumes
	return lunGroup, nil
}

// AddLunGroupWithCheck :
func (c *EternusClient) AddLunGroupWithCheck(lgName string, lunID string) (string, error) {
	// check lunGrp
	ret, err := c.request("show lun-groups", nil)
	if err != nil {
		log.Error("failed to show lun-groups: " + err.Error())
		return "", err
	}
	lgNumberStr := ""
	for _, v := range ret {
		if v["lun_group_name"] == lgName {
			lgNumber, _ := strconv.ParseInt(v["lun_group_no"], 10, 64)
			lgNumberStr = strconv.FormatInt(lgNumber, 10)
			break
		}
	}
	// if already exists for the target host, add volume to the lunGrp.
	if lgNumberStr != "" {
		param := map[string]string{
			"volume-number": lunID,
			"lg-number":     lgNumberStr,
		}
		ret, err = c.request("set lun-group", param)
		if err != nil {
			log.Error("failed to set lun-group: " + err.Error())
			return "", err
		}
		return lgNumberStr, nil
	}

	// if does not exists for the target host, create new lunGrp.
	lun := "0"
	param := map[string]string{
		"name":          lgName,
		"volume-number": lunID,
		"lun":           lun,
	}
	ret, err = c.request("create lun-group", param)
	if err != nil {
		log.Error("failed to create lun-group: " + err.Error())
		return "", err
	}
	v := ret[0]
	lunNumber, _ := strconv.ParseInt(v["lun_group_number"], 16, 64)
	return strconv.FormatInt(lunNumber, 10), nil
}

// RemoveVolumeFromLunGroup :
func (c *EternusClient) RemoveVolumeFromLunGroup(lunID string, lgName string) error {
	param := map[string]string{
		"lg-name": lgName,
		"lun":     lunID,
	}
	_, err := c.request("delete lun-group", param)
	if err != nil {
		log.Error("failed to remove volume from lun-group: " + err.Error())
		return err
	}
	return nil
}

// DeleteLunGroup :
func (c *EternusClient) DeleteLunGroupByName(lgName string) error {
	param := map[string]string{
		"lg-name": lgName,
	}
	ret, err := c.request("delete lun-group", param)
	if err != nil {
		if len(ret) == 1 && ret[0]["errorCode"] == NotFound {
			log.Info("target lun group already deleted")
			return nil
		}
		log.Error("failed to delete lun-group: " + err.Error())
		return err
	}
	return nil
}

// GetIscsiPortInfo :
func (c *EternusClient) GetIscsiPortInfo(ceSupport bool, needHostAffinity bool) (IscsiPortInfo, error) {
	portInfo := IscsiPortInfo{}
	// select port
	ret, err := c.request("show iscsi-parameters", nil)
	if err != nil {
		log.Error("failed to get iscsi-parameters: " + err.Error())
		return portInfo, err
	}

	usePort, portNumber := c.getConnectionPort(ret, ceSupport, needHostAffinity)
	if portNumber == "" {
		msg := "there is no iscsi port."
		log.Error(msg)
		return portInfo, errors.New(msg)
	}

	tcpPort, _ := strconv.ParseInt(usePort["tcp_port_number"], 16, 64)
	isnsPort, _ := strconv.ParseInt(usePort["isns_server_port"], 16, 64)
	portInfo.PortNumber = portNumber
	portInfo.IscsiName = usePort["iscsi_name"]
	portInfo.Ip = usePort["ip_address"]
	portInfo.TcpPort = int(tcpPort)
	portInfo.IsnsServerIp = usePort["isns_server_ip"]
	portInfo.IsnsServerPort = int(isnsPort)

	return portInfo, nil
}

// GetFcPortInfo :
func (c *EternusClient) GetFcPortInfo(ceSupport bool, needHostAffinity bool) (FcPortInfo, error) {
	portInfo := FcPortInfo{}
	// select port
	ret, err := c.request("show fc-parameters", nil)
	if err != nil {
		log.Error("failed to get fc-parameters: " + err.Error())
		return portInfo, err
	}
	usePort, portNumber := c.getConnectionPort(ret, ceSupport, needHostAffinity)

	if portNumber == "" {
		msg := "there is no fc port."
		log.Error(msg)
		return portInfo, errors.New(msg)
	}
	portInfo.PortNumber = portNumber
	portInfo.Wwpn = usePort["wwpn"]

	return portInfo, nil
}

func (c *EternusClient) getConnectionPort(portList []map[string]string,
	ceSupport bool, needHostAffinity bool) (map[string]string, string) {
	port := ""
	usePort := map[string]string{}
	for _, v := range portList {
		// if port_mode is not "CA" and "CA/RA", skip
		if v["port_mode"] != "00" && v["port_mode"] != "04" {
			continue
		}
		if v["host_affinity"] == "00" && needHostAffinity {
			usePort = v
			break
		} else if v["host_affinity"] != "00" && !needHostAffinity {
			usePort = v
			break
		}
	}
	if len(usePort) == 0 {
		return usePort, port
	}
	if ceSupport {
		port = GetPortNumberV2(usePort["ca_module_id"], usePort["port_number"])
	} else {
		port = GetPortNumber(usePort["ca_module_id"], usePort["port_number"])
	}
	return usePort, port
}

// AddHostAffinity :
func (c *EternusClient) AddHostAffinity(lunGrpID string, hostID string, iscsiPort string) (string, error) {
	// create new host affinity
	param := map[string]string{
		"port":        iscsiPort,
		"lg-number":   lunGrpID,
		"host-number": hostID,
	}
	ret, err := c.request("set host-affinity", param)
	if err != nil {
		log.Error("failed to set host-affinity: " + err.Error())
		return "", err
	}
	v := ret[0]
	lunMaskGroupNo, _ := strconv.ParseInt(v["lun_mask_group_no"], 16, 64)
	return strconv.FormatInt(lunMaskGroupNo, 10), nil
}

// DeleteHostAffinity :
func (c *EternusClient) DeleteHostAffinity(portNumber string, hostname string) error {
	param := map[string]string{
		"port":      portNumber,
		"host-name": hostname,
		"mode":      "all",
	}
	ret, err := c.request("release host-affinity", param)
	if err != nil {
		if len(ret) == 1 && ret[0]["errorCode"] == NotFound {
			log.Info("target host affinity already deleted")
			return nil
		}
		log.Error("failed to release host-affinity: " + err.Error())
		return err
	}
	return nil
}

// GetHostLunID :
func (c *EternusClient) GetHostLunID(lunGrpID string, lunID string) (string, error) {
	param := map[string]string{
		"lg-number": lunGrpID,
	}
	ret, err := c.request("show lun-group", param)
	if err != nil {
		log.Error("failed to get lun-group: " + err.Error())
		return "", err
	}
	var hostLunID int64
	for _, v := range ret {
		volID, _ := strconv.ParseInt(v["volume_no"], 16, 64)
		if strconv.FormatInt(volID, 10) == lunID {
			hostLunID, _ = strconv.ParseInt(v["lun"], 16, 64)
		}
	}
	return strconv.FormatInt(hostLunID, 10), nil
}

// AddFcHostWithCheck :
func (c *EternusClient) AddFcHostWithCheck(name string, wwnName string) (string, bool, error) {
	// check duplicate host. if already exists, retuen exist host id.
	ret, err := c.request("show host-wwn-names", nil)
	if err != nil {
		log.Error("failed to show host-wwn-names: " + err.Error())
		return "", false, err
	}
	for _, v := range ret {
		if strings.ToUpper(v["host_wwn_name"]) == strings.ToUpper(wwnName) {
			hostNumber, _ := strconv.ParseInt(v["host_wwn_no"], 16, 64)
			return strconv.FormatInt(hostNumber, 10), true, nil
		}
	}

	// create new host
	param := map[string]string{
		"wwn":  wwnName,
		"name": name,
	}
	ret, err = c.request("create host-wwn-name", param)
	if err != nil {
		log.Error("failed to create host-wwn-name: " + err.Error())
		return "", true, err
	}
	v := ret[0]
	hostNumber, _ := strconv.ParseInt(v["host_number"], 16, 64)
	return strconv.FormatInt(hostNumber, 10), false, nil
}

// DeleteFcHost :
func (c *EternusClient) DeleteFcHostByName(name string) error {
	param := map[string]string{
		"host-name": name,
	}
	_, err := c.request("delete host-wwn-name", param)
	if err != nil {
		log.Error("failed to delete host-wwn-name: " + err.Error())
		return err
	}
	return nil
}

// ListMapping :
func (c *EternusClient) ListMapping(port string) ([]Mapping, error) {
	mappings := []Mapping{}
	param := map[string]string{
		"port": port,
	}
	ret, err := c.request("show mapping", param)
	if err != nil {
		log.Error("failed to show mapping: " + err.Error())
		return nil, err
	}

	for _, v := range ret {
		lun, _ := strconv.ParseInt(v["lun"], 16, 64)
		volID, _ := strconv.ParseInt(v["volume_number"], 16, 64)
		tmpSize, _ := strconv.ParseInt(v["volume_size"], 16, 64)
		size := tmpSize / LBASize
		tmpMap := Mapping{}
		tmpMap.Lun = strconv.FormatInt(lun, 10)
		tmpMap.VolumeNumber = strconv.FormatInt(volID, 10)
		tmpMap.VolumeName = v["volume_name"]
		tmpMap.VolumeRawStatus = v["volume_raw_status"]
		tmpMap.VolumeRoundStatus = v["volume_round_status"]
		tmpMap.VolumeSize = size
		mappings = append(mappings, tmpMap)
	}
	return mappings, nil
}

// AddMapping :
func (c *EternusClient) AddMapping(lunID string, hostLunID string, port string) error {
	param := map[string]string{
		"port":          port,
		"volume-number": lunID,
		"lun":           hostLunID,
	}
	_, err := c.request("set mapping", param)
	if err != nil {
		log.Error("failed to set mapping: " + err.Error())
		return err
	}
	return nil
}

// DeleteMapping :
func (c *EternusClient) DeleteMapping(hostLunID string, Port string) error {
	param := map[string]string{
		"port": Port,
		"lun":  hostLunID,
	}
	_, err := c.request("release mapping", param)
	if err != nil {
		log.Error("failed to release mapping: " + err.Error())
		return err
	}
	return nil
}

// CreateSnapshot is for admin role
func (c *EternusClient) CreateSnapshot(srcLunID string, destLunID string) error {
	param := map[string]string{
		"source-volume-number":      srcLunID,
		"destination-volume-number": destLunID,
	}
	_, err := c.request("start advanced-copy", param)
	if err != nil {
		log.Error("failed to start advanced-copy: " + err.Error())
		return err
	}
	return nil
}

// ListSnapshot is for admin role
func (c *EternusClient) ListSnapshot() ([]SnapShot, error) {
	param := map[string]string{
		"type": "sopc+",
	}
	cmd := "show advanced-copy-sessions"
	ret, err := c.requestForAdmin(cmd, param)
	if err != nil {
		log.Error("failed to show advanced-copy-sessions: " + err.Error())
		return nil, err
	}
	output := [][]string{}
	for i, s := range strings.Split(ret.String(), "\r\n") {
		// ignore empty line(first elem)
		if i < 5 {
			continue
		}
		// ignore last line and stop parse
		if s == "CLI> " {
			break
		}
		output = append(output, strings.Split(s, " "))
	}
	snapshotList := []SnapShot{}
	for _, v := range output {
		sp := []string{}
		snapshot := SnapShot{}
		for _, e := range v {
			if e != "" {
				sp = append(sp, e)
			}
		}
		snapshot.Sid = sp[0]
		snapshot.Gen = sp[1]
		snapshot.GenTotal = sp[2]
		snapshot.Type = sp[3]
		snapshot.VolumeType = sp[4]
		snapshot.SrcNo = sp[5]
		snapshot.SrcName = sp[6]
		snapshot.DestNo = sp[7]
		snapshot.DestName = sp[8]
		snapshot.Status = sp[9]
		snapshot.Phase = sp[10]
		snapshot.ErrorCode = sp[11]
		snapshot.Requestor = sp[12]
		snapshotList = append(snapshotList, snapshot)
	}
	return snapshotList, nil
}

// DeleteSnapshot is for admin role
func (c *EternusClient) DeleteSnapshot(sid string) error {
	param := map[string]string{
		"session-id": sid,
	}
	_, err := c.requestForAdmin("stop advanced-copy", param)
	if err != nil {
		log.Error("failed to stop advanced-copy: " + err.Error())
		errID := strings.Split(err.Error(), " ")[1]
		if errID == ("E" + NotFound) {
			log.Info("target snapshot session already deleted. Ignore the error.")
			return nil
		}
		return err
	}
	return nil
}
