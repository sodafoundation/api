package oceanstor

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
)

type AuthOptions struct {
	Username        string `yaml:"username,omitempty"`
	Password        string `yaml:"password,omitempty"`
	PwdEncrypter    string `yaml:"PwdEncrypter,omitempty"`
	EnableEncrypted bool   `yaml:"EnableEncrypted,omitempty"`
	Endpoints       string `yaml:"endpoints,omitempty"`
	Insecure        bool   `yaml:"insecure,omitempty"`
}

type Replication struct {
	RemoteAuthOpt AuthOptions `yaml:"remoteAuthOptions"`
}

type OceanStorConfig struct {
	AuthOptions `yaml:"authOptions"`
	Replication `yaml:"replication"`
	Pool        map[string]PoolProperties `yaml:"pool,flow"`
	TargetIp    string                    `yaml:"targetIp,omitempty"`
}

const UnitGi = 1024 * 1024 * 1024

func EncodeName(id string) string {
	h := md5.New()
	h.Write([]byte(id))
	encodedName := hex.EncodeToString(h.Sum(nil))
	prefix := strings.Split(id, "-")[0] + "-"
	postfix := encodedName[:MaxNameLength-len(prefix)]
	return prefix + postfix
}

func EncodeHostName(name string) string {
	isMatch, _ := regexp.MatchString(`[[:alnum:]-_.]+`, name)
	if len(name) > MaxNameLength || !isMatch {
		h := md5.New()
		h.Write([]byte(name))
		encodedName := hex.EncodeToString(h.Sum(nil))
		return encodedName[:MaxNameLength]
	}
	return name
}

func TruncateDescription(desc string) string {
	if len(desc) > MaxDescriptionLength {
		desc = desc[:MaxDescriptionLength]
	}
	return desc
}

func Sector2Gb(sec string) int64 {
	size, err := strconv.ParseInt(sec, 10, 64)
	if err != nil {
		log.Error("Convert capacity from string to number failed, error:", err)
		return 0
	}
	return size * 512 / UnitGi
}

func Gb2Sector(gb int64) int64 {
	return gb * UnitGi / 512
}
