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
	"hash/fnv"
	"net"
	"strconv"
	"strings"

	. "github.com/opensds/opensds/contrib/drivers/utils/config"
)

type AuthOptions struct {
	Username        string `yaml:"username,omitempty"`
	Password        string `yaml:"password,omitempty"`
	AdminUsername   string `yaml:"adminUsername,omitempty"`
	AdminPassword   string `yaml:"adminPassword,omitempty"`
	PwdEncrypter    string `yaml:"PwdEncrypter,omitempty"`
	EnableEncrypted bool   `yaml:"EnableEncrypted,omitempty"`
	Endpoint        string `yaml:"endpoint,omitempty"`
	Insecure        bool   `yaml:"insecure,omitempty"`
}

type Replication struct {
	RemoteAuthOpt AuthOptions `yaml:"remoteAuthOptions"`
}

type EternusConfig struct {
	AuthOptions `yaml:"authOptions"`
	Replication `yaml:"replication"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
	TargetIp    string                    `yaml:"targetIp,omitempty"`
	CeSupport   bool                      `yaml:"ceSupport,omitempty"`
}

func IsIPv4(ip string) bool {
	return true
}

// GetPortNumber
func GetPortNumber(caModuleId string, portNumber string) string {
	caInt, _ := strconv.ParseInt(caModuleId[0:1], 16, 64)
	cmInt, _ := strconv.ParseInt(caModuleId[1:2], 16, 64)
	portInt, _ := strconv.ParseInt(portNumber[1:2], 16, 64)
	cm := cmInt % 8
	ca := caInt % 4
	port := (caInt/4/2)*2 + portInt%2
	ret := strconv.FormatInt(cm, 16) + strconv.FormatInt(ca, 16) + strconv.FormatInt(port, 16)
	return ret
}

// GetPortNumberV2
func GetPortNumberV2(caModuleId string, portNumber string) string {
	var base int64 = 4
	caInt, _ := strconv.ParseInt(caModuleId[0:1], 16, 64)
	cmInt, _ := strconv.ParseInt(caModuleId[1:2], 16, 64)
	portInt, _ := strconv.ParseInt(portNumber[1:2], 16, 64)
	ce := ((caInt % 8) / 2) + (((cmInt - base) / 4) * 4)
	cm := cmInt % 8
	ca := caInt % 4
	port := (caInt/4/2)*2 + portInt%2
	ret := strconv.FormatInt(ce, 16) + strconv.FormatInt(cm, 16) + strconv.FormatInt(ca, 16) + strconv.FormatInt(port, 16)
	return ret
}

// ParseIPv4 : convert hex string to ip address
func ParseIPv4(ip string) string {
	ipStr := []string{}
	for i := 0; i < len(ip); i += 2 {
		tmpIP, _ := strconv.ParseInt(ip[i:i+2], 16, 64)
		ipStr = append(ipStr, strconv.FormatInt(tmpIP, 10))
	}
	return strings.Join(ipStr, ".")
}

// ParseIPv6 : convert hex string to ip address
func ParseIPv6(ip string) string {
	ipStr := ip[:4] + ":" + ip[4:8] + ":" + ip[8:12] + ":" + ip[12:16] + ":"
	ipStr += ip[16:20] + ":" + ip[20:24] + ":" + ip[24:28] + ":" + ip[28:32]
	return ipStr
}

// EqualIP : if ip address is eauql return ture
func EqualIP(ip1 string, ip2 string) bool {
	return net.ParseIP(ip1).Equal(net.ParseIP(ip2))
}

// GetFnvHash :
func GetFnvHash(str string) string {
	hash := fnv.New64()
	hash.Write([]byte(str))
	val := hash.Sum64()
	return strconv.FormatUint(uint64(val), 16)
}
