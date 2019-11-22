package chubaofs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/golang/glog"
	pb "github.com/opensds/opensds/pkg/model/proto"
)

type RequestType int

func (t RequestType) String() string {
	switch t {
	case createVolumeRequest:
		return "CreateVolume"
	case deleteVolumeRequest:
		return "DeleteVolume"
	default:
	}
	return "N/A"
}

const (
	createVolumeRequest RequestType = iota
	deleteVolumeRequest
)

type clusterInfoResponseData struct {
	LeaderAddr string `json:"LeaderAddr"`
}

type clusterInfoResponse struct {
	Code int                      `json:"code"`
	Msg  string                   `json:"msg"`
	Data *clusterInfoResponseData `json:"data"`
}

// Create and Delete Volume Response
type generalVolumeResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

/*
 * This functions sends http request to the on-premise cluster to
 * get cluster info.
 */
func getClusterInfo(host string) (string, error) {
	url := "http://" + host + "/admin/getCluster"
	log.Infof("chubaofs: GetClusterInfo(%v)", url)

	httpResp, err := http.Get(url)
	if err != nil {
		log.Errorf("chubaofs: failed to GetClusterInfo, url(%v) err(%v)", url, err)
		return "", err
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Errorf("chubaofs: failed to read response, url(%v) err(%v)", url, err)
		return "", err
	}

	resp := &clusterInfoResponse{}
	if err = json.Unmarshal(body, resp); err != nil {
		errmsg := fmt.Sprintf("chubaofs: getClusterInf failed to unmarshal, bodyLen(%d) err(%v)", len(body), err)
		log.Error(errmsg)
		return "", errors.New(errmsg)
	}

	log.Infof("chubaofs: GetClusterInfo, url(%v), resp(%v)", url, resp)

	if resp.Code != 0 {
		errmsg := fmt.Sprintf("chubaofs: GetClusterInfo code NOK, url(%v) code(%v) msg(%v)", url, resp.Code, resp.Msg)
		log.Error(errmsg)
		return "", errors.New(errmsg)
	}

	if resp.Data == nil {
		errmsg := fmt.Sprintf("chubaofs: GetClusterInfo nil data, url(%v) code(%v) msg(%v)", url, resp.Code, resp.Msg)
		log.Error(errmsg)
		return "", errors.New(errmsg)
	}

	return resp.Data.LeaderAddr, nil
}

/*
 * This function sends http request to the on-premise cluster to create
 * or delete a volume according to request type.
 */
func createOrDeleteVolume(req RequestType, leader, name, owner string, size int64) error {
	var url string

	switch req {
	case createVolumeRequest:
		sizeInGB := size
		url = fmt.Sprintf("http://%s/admin/createVol?name=%s&capacity=%v&owner=%v", leader, name, sizeInGB, owner)
	case deleteVolumeRequest:
		key := md5.New()
		if _, err := key.Write([]byte(owner)); err != nil {
			return errors.New(fmt.Sprintf("chubaofs: failed to get md5 sum of owner err(%v)", err))
		}
		url = fmt.Sprintf("http://%s/vol/delete?name=%s&authKey=%v", leader, name, hex.EncodeToString(key.Sum(nil)))
	default:
		return errors.New("chubaofs: request type not recognized! %v")
	}

	log.Infof("chubaofs: %v url(%v)", req, url)

	httpResp, err := http.Get(url)
	if err != nil {
		errmsg := fmt.Sprintf("chubaofs: %v failed, url(%v) err(%v)", req, url, err)
		return errors.New(errmsg)
	}
	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		errmsg := fmt.Sprintf("chubaofs: %v failed to read http response body, bodyLen(%d) err(%v)", req, len(body), err)
		return errors.New(errmsg)
	}

	resp := &generalVolumeResponse{}
	if err := json.Unmarshal(body, resp); err != nil {
		errmsg := fmt.Sprintf("chubaofs: %v failed to unmarshal, url(%v) msg(%v)", req, url, resp.Msg)
		return errors.New(errmsg)
	}

	if resp.Code != 0 {
		errmsg := fmt.Sprintf("chubaofs: %v failed, url(%v) code(%v) msg(%v)", req, url, resp.Code, resp.Msg)
		return errors.New(errmsg)
	}

	log.Infof("chubaofs: %v url(%v) successful!", req, url)
	return nil
}

func doMount(cmdName string, confFile []string) error {
	env := []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	for _, conf := range confFile {
		cmd := exec.Command(cmdName, "-c", conf)
		cmd.Env = append(cmd.Env, env...)
		if msg, err := cmd.CombinedOutput(); err != nil {
			return errors.New(fmt.Sprintf("chubaofs: failed to start client daemon, msg: %v , err: %v", string(msg), err))
		}
	}
	return nil
}

func doUmount(mntPoints []string) error {
	env := []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	for _, mnt := range mntPoints {
		cmd := exec.Command("umount", mnt)
		cmd.Env = append(cmd.Env, env...)
		msg, err := cmd.CombinedOutput()
		if err != nil {
			return errors.New(fmt.Sprintf("chubaofs: failed to umount, msg: %v , err: %v", msg, err))
		}
	}
	return nil
}

/*
 * This function creates client config files according to export locations.
 */
func prepareConfigFiles(d *Driver, opt *pb.CreateFileShareOpts) (configFiles, fsMntPoints []string, owner string, err error) {
	volName := opt.GetId()
	configFiles = make([]string, 0)
	fsMntPoints = make([]string, 0)

	/*
	 * Check client runtime path.
	 */
	fi, err := os.Stat(d.conf.ClientPath)
	if err != nil || !fi.Mode().IsDir() {
		err = errors.New(fmt.Sprintf("chubaofs: invalid client runtime path, path: %v, err: %v", d.conf.ClientPath, err))
		return
	}

	clientConf := path.Join(d.conf.ClientPath, volName, "conf")
	clientLog := path.Join(d.conf.ClientPath, volName, "log")
	clientWarnLog := path.Join(d.conf.ClientPath, volName, "warnlog")

	if err = os.MkdirAll(clientConf, os.ModeDir); err != nil {
		err = errors.New(fmt.Sprintf("chubaofs: failed to create client config dir, path: %v , err: %v", clientConf, err))
		return
	}
	defer func() {
		if err != nil {
			log.Warningf("chubaofs: cleaning config dir, %v", clientConf)
			os.RemoveAll(clientConf)
		}
	}()

	if err = os.MkdirAll(clientLog, os.ModeDir); err != nil {
		err = errors.New(fmt.Sprintf("chubaofs: failed to create client log dir, path: %v", clientLog))
		return
	}
	defer func() {
		if err != nil {
			log.Warningf("chubaofs: cleaning log dir, %v", clientLog)
			os.RemoveAll(clientLog)
		}
	}()

	if err = os.MkdirAll(clientWarnLog, os.ModeDir); err != nil {
		err = errors.New(fmt.Sprintf("chubaofs: failed to create client warn log dir, path: %v , err: %v", clientWarnLog, err))
		return
	}
	defer func() {
		if err != nil {
			log.Warningf("chubaofs: cleaning warn log dir, %v", clientWarnLog)
			os.RemoveAll(clientWarnLog)
		}
	}()

	/*
	 * Check and create mount point directory.
	 * Mount point dir has to be newly created.
	 */
	locations := make([]string, 0)
	if len(opt.ExportLocations) == 0 {
		locations = append(locations, path.Join(d.conf.MntPoint, volName))
	} else {
		locations = append(locations, opt.ExportLocations...)
	}

	/*
	 * Mount point has to be absolute path to avoid umounting the current
	 * working directory.
	 */
	fsMntPoints, err = createAbsMntPoints(locations)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			for _, mnt := range fsMntPoints {
				os.RemoveAll(mnt)
			}
		}
	}()

	/*
	 * Generate client mount config file.
	 */
	mntConfig := make(map[string]interface{})
	mntConfig[KVolumeName] = volName
	mntConfig[KMasterAddr] = strings.Join(d.conf.MasterAddr, ",")
	mntConfig[KLogDir] = clientLog
	mntConfig[KWarnLogDir] = clientWarnLog
	if d.conf.LogLevel != "" {
		mntConfig[KLogLevel] = d.conf.LogLevel
	} else {
		mntConfig[KLogLevel] = defaultLogLevel
	}
	if d.conf.Owner != "" {
		mntConfig[KOwner] = d.conf.Owner
	} else {
		mntConfig[KOwner] = defaultOwner
	}
	if d.conf.ProfPort != "" {
		mntConfig[KProfPort] = d.conf.ProfPort
	} else {
		mntConfig[KProfPort] = defaultProfPort
	}

	owner = mntConfig[KOwner].(string)

	for i, mnt := range fsMntPoints {
		mntConfig[KMountPoint] = mnt
		data, e := json.MarshalIndent(mntConfig, "", "    ")
		if e != nil {
			err = errors.New(fmt.Sprintf("chubaofs: failed to generate client config file, err(%v)", e))
			return
		}
		filePath := path.Join(clientConf, strconv.Itoa(i), clientConfigFileName)
		_, e = generateFile(filePath, data)
		if e != nil {
			err = errors.New(fmt.Sprintf("chubaofs: failed to generate client config file, err(%v)", e))
			return
		}
		configFiles = append(configFiles, filePath)
	}

	return
}

/*
 * This function creates mount points according to the specified paths,
 * and returns the absolute paths.
 */
func createAbsMntPoints(locations []string) (mntPoints []string, err error) {
	mntPoints = make([]string, 0)
	for _, loc := range locations {
		mnt, e := filepath.Abs(loc)
		if e != nil {
			err = errors.New(fmt.Sprintf("chubaofs: failed to get absolute path of export locations, loc: %v , err: %v", loc, e))
			return
		}
		if e = os.MkdirAll(mnt, os.ModeDir); e != nil {
			err = errors.New(fmt.Sprintf("chubaofs: failed to create mount point dir, mnt: %v , err: %v", mnt, e))
			return
		}
		mntPoints = append(mntPoints, mnt)
	}
	return
}

/*
 * This function generates the target file with specified path, and writes data.
 */
func generateFile(filePath string, data []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(data)
}
