package rados_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	//"net"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/ceph/go-ceph/rados"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RadosTestSuite struct {
	suite.Suite
	conn  *rados.Conn
	ioctx *rados.IOContext
	pool  string
	count int
}

func (suite *RadosTestSuite) SetupSuite() {
	conn, err := rados.NewConn()
	require.NoError(suite.T(), err)
	defer conn.Shutdown()

	conn.ReadDefaultConfigFile()

	if err = conn.Connect(); assert.NoError(suite.T(), err) {
		pool := uuid.Must(uuid.NewV4()).String()
		if err = conn.MakePool(pool); assert.NoError(suite.T(), err) {
			suite.pool = pool
			return
		}
	}

	suite.T().FailNow()
}

func (suite *RadosTestSuite) SetupTest() {
	suite.conn = nil
	suite.ioctx = nil
	suite.count = 0

	conn, err := rados.NewConn()
	require.NoError(suite.T(), err)
	suite.conn = conn
	suite.conn.ReadDefaultConfigFile()
}

func (suite *RadosTestSuite) SetupConnection() {
	if err := suite.conn.Connect(); assert.NoError(suite.T(), err) {
		ioctx, err := suite.conn.OpenIOContext(suite.pool)
		if assert.NoError(suite.T(), err) {
			suite.ioctx = ioctx
			return
		}
	}
	suite.conn.Shutdown()
	suite.T().FailNow()
}

func (suite *RadosTestSuite) GenObjectName() string {
	name := fmt.Sprintf("%s_%d", suite.T().Name(), suite.count)
	suite.count++
	return name
}

func (suite *RadosTestSuite) RandomBytes(size int) []byte {
	bytes := make([]byte, size)
	n, err := rand.Read(bytes)
	require.Equal(suite.T(), n, size)
	require.NoError(suite.T(), err)
	return bytes
}

func (suite *RadosTestSuite) TearDownTest() {
	if suite.ioctx != nil {
		suite.ioctx.Destroy()
	}
	suite.conn.Shutdown()
}

func (suite *RadosTestSuite) TearDownSuite() {
	conn, err := rados.NewConn()
	require.NoError(suite.T(), err)
	defer conn.Shutdown()

	conn.ReadDefaultConfigFile()

	if err = conn.Connect(); assert.NoError(suite.T(), err) {
		err = conn.DeletePool(suite.pool)
		assert.NoError(suite.T(), err)
	}
}

func TestVersion(t *testing.T) {
	var major, minor, patch = rados.Version()
	assert.False(t, major < 0 || major > 1000, "invalid major")
	assert.False(t, minor < 0 || minor > 1000, "invalid minor")
	assert.False(t, patch < 0 || patch > 1000, "invalid patch")
}

func (suite *RadosTestSuite) TestGetFSID() {
	fsid, err := suite.conn.GetFSID()
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), fsid, "")
}

func (suite *RadosTestSuite) TestGetSetConfigOption() {
	// rejects invalid options
	err := suite.conn.SetConfigOption("___dne___", "value")
	assert.Error(suite.T(), err, "Invalid option")

	// verify SetConfigOption changes a values
	prev_val, err := suite.conn.GetConfigOption("log_file")
	assert.NoError(suite.T(), err, "Invalid option")

	err = suite.conn.SetConfigOption("log_file", "/dev/null")
	assert.NoError(suite.T(), err, "Invalid option")

	curr_val, err := suite.conn.GetConfigOption("log_file")
	assert.NoError(suite.T(), err, "Invalid option")

	assert.NotEqual(suite.T(), prev_val, "/dev/null")
	assert.Equal(suite.T(), curr_val, "/dev/null")
}

func (suite *RadosTestSuite) TestParseDefaultConfigEnv() {
	prev_val, err := suite.conn.GetConfigOption("log_file")
	assert.NoError(suite.T(), err, "Invalid option")

	err = os.Setenv("CEPH_ARGS", "--log-file /dev/null")
	assert.NoError(suite.T(), err)

	err = suite.conn.ParseDefaultConfigEnv()
	assert.NoError(suite.T(), err)

	curr_val, err := suite.conn.GetConfigOption("log_file")
	assert.NoError(suite.T(), err, "Invalid option")

	assert.NotEqual(suite.T(), prev_val, "/dev/null")
	assert.Equal(suite.T(), curr_val, "/dev/null")
}

func (suite *RadosTestSuite) TestParseCmdLineArgs() {
	prev_val, err := suite.conn.GetConfigOption("log_file")
	assert.NoError(suite.T(), err, "Invalid option")

	args := []string{"--log_file", "/dev/null"}
	err = suite.conn.ParseCmdLineArgs(args)
	assert.NoError(suite.T(), err)

	curr_val, err := suite.conn.GetConfigOption("log_file")
	assert.NoError(suite.T(), err, "Invalid option")

	assert.NotEqual(suite.T(), prev_val, "/dev/null")
	assert.Equal(suite.T(), curr_val, "/dev/null")
}

func (suite *RadosTestSuite) TestReadConfigFile() {
	// check current log_file value
	prev_str, err := suite.conn.GetConfigOption("log_max_new")
	assert.NoError(suite.T(), err)

	prev_val, err := strconv.Atoi(prev_str)
	assert.NoError(suite.T(), err)

	// create conf file that changes log_file conf option
	file, err := ioutil.TempFile("/tmp", "go-rados")
	assert.NoError(suite.T(), err)

	next_val := prev_val + 1
	conf := fmt.Sprintf("[global]\nlog_max_new = %d\n", next_val)
	_, err = io.WriteString(file, conf)
	assert.NoError(suite.T(), err)

	// parse the config file
	err = suite.conn.ReadConfigFile(file.Name())
	assert.NoError(suite.T(), err)

	// check current log_file value
	curr_str, err := suite.conn.GetConfigOption("log_max_new")
	assert.NoError(suite.T(), err)

	curr_val, err := strconv.Atoi(curr_str)
	assert.NoError(suite.T(), err)

	assert.NotEqual(suite.T(), prev_str, curr_str)
	assert.Equal(suite.T(), curr_val, prev_val+1)

	file.Close()
	os.Remove(file.Name())
}

func (suite *RadosTestSuite) TestGetClusterStats() {
	suite.SetupConnection()

	// grab current stats
	prev_stat, err := suite.conn.GetClusterStats()
	fmt.Printf("prev_stat: %+v\n", prev_stat)
	assert.NoError(suite.T(), err)

	// make some changes to the cluster
	buf := make([]byte, 1<<20)
	for i := 0; i < 10; i++ {
		objname := suite.GenObjectName()
		suite.ioctx.Write(objname, buf, 0)
	}

	// wait a while for the stats to change
	for i := 0; i < 30; i++ {
		stat, err := suite.conn.GetClusterStats()
		assert.NoError(suite.T(), err)

		// wait for something to change
		if stat == prev_stat {
			fmt.Printf("curr_stat: %+v (trying again...)\n", stat)
			time.Sleep(time.Second)
		} else {
			// success
			fmt.Printf("curr_stat: %+v (change detected)\n", stat)
			return
		}
	}

	suite.T().Error("Cluster stats aren't changing")
}

func (suite *RadosTestSuite) TestGetInstanceID() {
	suite.SetupConnection()

	id := suite.conn.GetInstanceID()
	assert.NotEqual(suite.T(), id, 0)
}

func (suite *RadosTestSuite) TestMakeDeletePool() {
	suite.SetupConnection()

	// get current list of pool
	pools, err := suite.conn.ListPools()
	assert.NoError(suite.T(), err)

	// check that new pool name is unique
	new_name := uuid.Must(uuid.NewV4()).String()
	for _, poolname := range pools {
		if new_name == poolname {
			suite.T().Error("Random pool name exists!")
			return
		}
	}

	// create pool
	err = suite.conn.MakePool(new_name)
	assert.NoError(suite.T(), err)

	// verify that the new pool name exists
	pools, err = suite.conn.ListPools()
	assert.NoError(suite.T(), err)

	found := false
	for _, poolname := range pools {
		if new_name == poolname {
			found = true
		}
	}

	if !found {
		suite.T().Error("Cannot find newly created pool")
	}

	// delete the pool
	err = suite.conn.DeletePool(new_name)
	assert.NoError(suite.T(), err)

	// verify that it is gone
	pools, err = suite.conn.ListPools()
	assert.NoError(suite.T(), err)

	found = false
	for _, poolname := range pools {
		if new_name == poolname {
			found = true
		}
	}

	if found {
		suite.T().Error("Deleted pool still exists")
	}
}

func (suite *RadosTestSuite) TestPingMonitor() {
	suite.SetupConnection()

	// mon id that should work with vstart.sh
	reply, err := suite.conn.PingMonitor("a")
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), reply, "")
}

func (suite *RadosTestSuite) TestWaitForLatestOSDMap() {
	suite.SetupConnection()

	err := suite.conn.WaitForLatestOSDMap()
	assert.NoError(suite.T(), err)
}

func (suite *RadosTestSuite) TestReadWrite() {
	suite.SetupConnection()

	bytes_in := []byte("input data")
	err := suite.ioctx.Write("obj", bytes_in, 0)
	assert.NoError(suite.T(), err)

	bytes_out := make([]byte, len(bytes_in))
	n_out, err := suite.ioctx.Read("obj", bytes_out, 0)

	assert.Equal(suite.T(), n_out, len(bytes_in))
	assert.Equal(suite.T(), bytes_in, bytes_out)

	bytes_in = []byte("input another data")
	err = suite.ioctx.WriteFull("obj", bytes_in)
	assert.NoError(suite.T(), err)

	bytes_out = make([]byte, len(bytes_in))
	n_out, err = suite.ioctx.Read("obj", bytes_out, 0)

	assert.Equal(suite.T(), n_out, len(bytes_in))
	assert.Equal(suite.T(), bytes_in, bytes_out)
}

func (suite *RadosTestSuite) TestAppend() {
	suite.SetupConnection()

	mirror := []byte{}
	oid := suite.GenObjectName()
	for i := 0; i < 3; i++ {
		// append random bytes
		bytes := suite.RandomBytes(33)
		err := suite.ioctx.Append(oid, bytes)
		assert.NoError(suite.T(), err)

		// what the object should contain
		mirror = append(mirror, bytes...)

		// check object contains what we expect
		buf := make([]byte, len(mirror))
		n, err := suite.ioctx.Read(oid, buf, 0)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), n, len(buf))
		assert.Equal(suite.T(), buf, mirror)
	}
}

func (suite *RadosTestSuite) TestReadNotFound() {
	suite.SetupConnection()

	var bytes []byte
	oid := suite.GenObjectName()
	_, err := suite.ioctx.Read(oid, bytes, 0)
	assert.Equal(suite.T(), err, rados.RadosErrorNotFound)
}

func (suite *RadosTestSuite) TestDeleteNotFound() {
	suite.SetupConnection()

	oid := suite.GenObjectName()
	err := suite.ioctx.Delete(oid)
	assert.Equal(suite.T(), err, rados.RadosErrorNotFound)
}

func (suite *RadosTestSuite) TestStatNotFound() {
	suite.SetupConnection()

	oid := suite.GenObjectName()
	_, err := suite.ioctx.Stat(oid)
	assert.Equal(suite.T(), err, rados.RadosErrorNotFound)
}

func (suite *RadosTestSuite) TestObjectStat() {
	suite.SetupConnection()

	oid := suite.GenObjectName()
	bytes := suite.RandomBytes(234)
	err := suite.ioctx.Write(oid, bytes, 0)
	assert.NoError(suite.T(), err)

	stat, err := suite.ioctx.Stat(oid)
	assert.Equal(suite.T(), uint64(len(bytes)), stat.Size)
	assert.NotNil(suite.T(), stat.ModTime)
}

func (suite *RadosTestSuite) TestGetPoolStats() {
	suite.SetupConnection()

	// grab current stats
	prev_stat, err := suite.ioctx.GetPoolStats()
	fmt.Printf("prev_stat: %+v\n", prev_stat)
	assert.NoError(suite.T(), err)

	// make some changes to the cluster
	buf := make([]byte, 1<<20)
	for i := 0; i < 10; i++ {
		oid := suite.GenObjectName()
		suite.ioctx.Write(oid, buf, 0)
	}

	// wait a while for the stats to change
	for i := 0; i < 30; i++ {
		stat, err := suite.ioctx.GetPoolStats()
		assert.NoError(suite.T(), err)

		// wait for something to change
		if stat == prev_stat {
			fmt.Printf("curr_stat: %+v (trying again...)\n", stat)
			time.Sleep(time.Second)
		} else {
			// success
			fmt.Printf("curr_stat: %+v (change detected)\n", stat)
			return
		}
	}

	suite.T().Error("Pool stats aren't changing")
}

func (suite *RadosTestSuite) TestGetPoolName() {
	suite.SetupConnection()

	name, err := suite.ioctx.GetPoolName()
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), name, suite.pool)
}

func (suite *RadosTestSuite) TestMonCommand() {
	suite.SetupConnection()

	command, err := json.Marshal(
		map[string]string{"prefix": "df", "format": "json"})
	assert.NoError(suite.T(), err)

	buf, info, err := suite.conn.MonCommand(command)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), info, "")

	var message map[string]interface{}
	err = json.Unmarshal(buf, &message)
	assert.NoError(suite.T(), err)
}

func (suite *RadosTestSuite) TestMonCommandWithInputBuffer() {
	suite.SetupConnection()

	entity := fmt.Sprintf("client.testMonCmdUser%d", time.Now().UnixNano())

	// first add the new test user, specifying its key in the input buffer
	command, err := json.Marshal(map[string]interface{}{
		"prefix": "auth add",
		"format": "json",
		"entity": entity,
	})
	assert.NoError(suite.T(), err)

	client_key := fmt.Sprintf(`
	  [%s]
	  key = AQD4PGNXBZJNHhAA582iUgxe9DsN+MqFN4Z6Jw==
	`, entity)

	inbuf := []byte(client_key)

	buf, info, err := suite.conn.MonCommandWithInputBuffer(command, inbuf)
	assert.NoError(suite.T(), err)
	expected_info := fmt.Sprintf("added key for %s", entity)
	assert.Equal(suite.T(), expected_info, info)
	assert.Equal(suite.T(), "", string(buf[:]))

	// get the key and verify that it's what we previously set
	command, err = json.Marshal(map[string]interface{}{
		"prefix": "auth get-key",
		"format": "json",
		"entity": entity,
	})
	assert.NoError(suite.T(), err)

	buf, info, err = suite.conn.MonCommand(command)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "", info)
	assert.Equal(suite.T(),
		`{"key":"AQD4PGNXBZJNHhAA582iUgxe9DsN+MqFN4Z6Jw=="}`,
		string(buf[:]))
}

func (suite *RadosTestSuite) TestObjectListObjects() {
	suite.SetupConnection()

	// objects currently in pool
	prevObjectList := []string{}
	err := suite.ioctx.ListObjects(func(oid string) {
		prevObjectList = append(prevObjectList, oid)
	})
	assert.NoError(suite.T(), err)

	// create some objects
	createdList := []string{}
	for i := 0; i < 10; i++ {
		oid := suite.GenObjectName()
		bytes := []byte("input data")
		err := suite.ioctx.Write(oid, bytes, 0)
		assert.NoError(suite.T(), err)
		createdList = append(createdList, oid)
	}

	// join the lists of objects
	expectedObjectList := prevObjectList
	expectedObjectList = append(expectedObjectList, createdList...)

	// now list the current set of objects in the pool
	currObjectList := []string{}
	err = suite.ioctx.ListObjects(func(oid string) {
		currObjectList = append(currObjectList, oid)
	})
	assert.NoError(suite.T(), err)

	// lists should be equal
	sort.Strings(currObjectList)
	sort.Strings(expectedObjectList)
	assert.Equal(suite.T(), currObjectList, expectedObjectList)
}

func (suite *RadosTestSuite) TestObjectIterator() {
	suite.SetupConnection()

	// current objs in default namespace
	prevObjectList := []string{}
	iter, err := suite.ioctx.Iter()
	assert.NoError(suite.T(), err)
	for iter.Next() {
		prevObjectList = append(prevObjectList, iter.Value())
	}
	iter.Close()
	assert.NoError(suite.T(), iter.Err())

	// create an object in a different namespace to verify that
	// iteration within a namespace does not return it
	suite.ioctx.SetNamespace("ns1")
	bytes_in := []byte("input data")
	err = suite.ioctx.Write(suite.GenObjectName(), bytes_in, 0)
	assert.NoError(suite.T(), err)

	// create some objects in default namespace
	suite.ioctx.SetNamespace("")
	createdList := []string{}
	for i := 0; i < 10; i++ {
		oid := suite.GenObjectName()
		bytes_in := []byte("input data")
		err = suite.ioctx.Write(oid, bytes_in, 0)
		assert.NoError(suite.T(), err)
		createdList = append(createdList, oid)
	}

	// prev list plus new oids
	expectedObjectList := prevObjectList
	expectedObjectList = append(expectedObjectList, createdList...)

	currObjectList := []string{}
	iter, err = suite.ioctx.Iter()
	assert.NoError(suite.T(), err)
	for iter.Next() {
		currObjectList = append(currObjectList, iter.Value())
	}
	iter.Close()
	assert.NoError(suite.T(), iter.Err())

	// curr list doesn't include the obj in ns1
	sort.Strings(expectedObjectList)
	sort.Strings(currObjectList)
	assert.Equal(suite.T(), currObjectList, expectedObjectList)
}

func (suite *RadosTestSuite) TestObjectIteratorAcrossNamespaces() {
	suite.SetupConnection()

	const perNamespace = 100

	// tests use a shared pool so namespaces need to be unique across tests.
	// below ns1=nsX and ns2=nsY. ns1 is used elsewhere.
	objectListNS1 := []string{}
	objectListNS2 := []string{}

	// populate list of current objects
	suite.ioctx.SetNamespace(rados.RadosAllNamespaces)
	existingList := []string{}
	iter, err := suite.ioctx.Iter()
	assert.NoError(suite.T(), err)
	for iter.Next() {
		existingList = append(existingList, iter.Value())
	}
	iter.Close()
	assert.NoError(suite.T(), iter.Err())

	// create some new objects in namespace: nsX
	createdList := []string{}
	suite.ioctx.SetNamespace("nsX")
	for i := 0; i < 10; i++ {
		oid := suite.GenObjectName()
		bytes_in := []byte("input data")
		err = suite.ioctx.Write(oid, bytes_in, 0)
		assert.NoError(suite.T(), err)
		createdList = append(createdList, oid)
	}
	assert.True(suite.T(), len(createdList) == 10)

	// create some new objects in namespace: nsY
	suite.ioctx.SetNamespace("nsY")
	for i := 0; i < 10; i++ {
		oid := suite.GenObjectName()
		bytes_in := []byte("input data")
		err = suite.ioctx.Write(oid, bytes_in, 0)
		assert.NoError(suite.T(), err)
		createdList = append(createdList, oid)
	}
	assert.True(suite.T(), len(createdList) == 20)

	suite.ioctx.SetNamespace(rados.RadosAllNamespaces)
	iter, err = suite.ioctx.Iter()
	assert.NoError(suite.T(), err)
	rogueList := []string{}
	for iter.Next() {
		if iter.Namespace() == "nsX" {
			objectListNS1 = append(objectListNS1, iter.Value())
		} else if iter.Namespace() == "nsY" {
			objectListNS2 = append(objectListNS2, iter.Value())
		} else {
			rogueList = append(rogueList, iter.Value())
		}
	}
	iter.Close()
	assert.NoError(suite.T(), iter.Err())

	assert.Equal(suite.T(), len(existingList), len(rogueList))
	assert.Equal(suite.T(), len(objectListNS1), 10)
	assert.Equal(suite.T(), len(objectListNS2), 10)

	objectList := []string{}
	objectList = append(objectList, objectListNS1...)
	objectList = append(objectList, objectListNS2...)
	sort.Strings(objectList)
	sort.Strings(createdList)

	assert.Equal(suite.T(), objectList, createdList)

	sort.Strings(rogueList)
	sort.Strings(existingList)
	assert.Equal(suite.T(), rogueList, existingList)
}

func (suite *RadosTestSuite) TestNewConnWithUser() {
	_, err := rados.NewConnWithUser("admin")
	assert.Equal(suite.T(), err, nil)
}

func (suite *RadosTestSuite) TestNewConnWithClusterAndUser() {
	_, err := rados.NewConnWithClusterAndUser("ceph", "client.admin")
	assert.Equal(suite.T(), err, nil)
}

func (suite *RadosTestSuite) TestReadWriteXattr() {
	suite.SetupConnection()

	oid := suite.GenObjectName()
	val := []byte("value")
	err := suite.ioctx.SetXattr(oid, "key", val)
	assert.NoError(suite.T(), err)

	out := make([]byte, len(val))
	n, err := suite.ioctx.GetXattr(oid, "key", out)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), n, len(out))

	assert.Equal(suite.T(), out, val)
}

func (suite *RadosTestSuite) TestListXattrs() {
	suite.SetupConnection()

	oid := suite.GenObjectName()
	xattrs := make(map[string][]byte)
	val := []byte("value")
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("key_%d", i)
		err := suite.ioctx.SetXattr(oid, name, val)
		assert.NoError(suite.T(), err)
		xattrs[name] = val
	}

	out, err := suite.ioctx.ListXattrs(oid)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), xattrs, out)
}

func (suite *RadosTestSuite) TestRmXattr() {
	suite.SetupConnection()

	oid := suite.GenObjectName()

	// 2 xattrs
	xattrs := make(map[string][]byte)
	xattrs["key1"] = []byte("val")
	xattrs["key2"] = []byte("val")
	assert.Equal(suite.T(), len(xattrs), 2)

	// add them to the object
	for key, value := range xattrs {
		err := suite.ioctx.SetXattr(oid, key, value)
		assert.NoError(suite.T(), err)
	}
	out, err := suite.ioctx.ListXattrs(oid)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(out), 2)
	assert.Equal(suite.T(), out, xattrs)

	// remove key1
	err = suite.ioctx.RmXattr(oid, "key1")
	assert.NoError(suite.T(), err)
	delete(xattrs, "key1")

	// verify key1 is gone
	out, err = suite.ioctx.ListXattrs(oid)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(out), 1)
	assert.Equal(suite.T(), out, xattrs)

	// remove key2
	err = suite.ioctx.RmXattr(oid, "key2")
	assert.NoError(suite.T(), err)
	delete(xattrs, "key2")

	// verify key2 is gone
	out, err = suite.ioctx.ListXattrs(oid)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(out), 0)
	assert.Equal(suite.T(), out, xattrs)
}

func (suite *RadosTestSuite) TestReadWriteOmap() {
	suite.SetupConnection()

	// set some key/value pairs on an object
	orig := map[string][]byte{
		"key1":          []byte("value1"),
		"key2":          []byte("value2"),
		"prefixed-key3": []byte("value3"),
		"empty":         []byte(""),
	}

	oid := suite.GenObjectName()
	err := suite.ioctx.SetOmap(oid, orig)
	assert.NoError(suite.T(), err)

	// verify that they can all be read back
	remaining := map[string][]byte{}
	for k, v := range orig {
		remaining[k] = v
	}

	err = suite.ioctx.ListOmapValues(oid, "", "", 4, func(key string, value []byte) {
		assert.Equal(suite.T(), remaining[key], value)
		delete(remaining, key)
	})

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(remaining), 0)

	// Get (with a fixed number of keys)
	fetched, err := suite.ioctx.GetOmapValues(oid, "", "", 4)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), orig, fetched)

	// Get All (with an iterator size bigger than the map size)
	fetched, err = suite.ioctx.GetAllOmapValues(oid, "", "", 100)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), orig, fetched)

	// Get All (with an iterator size smaller than the map size)
	fetched, err = suite.ioctx.GetAllOmapValues(oid, "", "", 1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), orig, fetched)

	// Remove
	err = suite.ioctx.RmOmapKeys(oid, []string{"key1", "prefixed-key3"})
	assert.NoError(suite.T(), err)

	fetched, err = suite.ioctx.GetOmapValues(oid, "", "", 4)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), map[string][]byte{
		"key2":  []byte("value2"),
		"empty": []byte(""),
	}, fetched)

	// Clear
	err = suite.ioctx.CleanOmap(oid)
	assert.NoError(suite.T(), err)

	fetched, err = suite.ioctx.GetOmapValues(oid, "", "", 4)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), map[string][]byte{}, fetched)
}

func (suite *RadosTestSuite) TestReadFilterOmap() {
	suite.SetupConnection()

	orig := map[string][]byte{
		"key1":          []byte("value1"),
		"prefixed-key3": []byte("value3"),
		"key2":          []byte("value2"),
	}

	oid := suite.GenObjectName()
	err := suite.ioctx.SetOmap(oid, orig)
	assert.NoError(suite.T(), err)

	// filter by prefix
	fetched, err := suite.ioctx.GetOmapValues(oid, "", "prefixed", 4)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), map[string][]byte{
		"prefixed-key3": []byte("value3"),
	}, fetched)

	// "start_after" a key
	fetched, err = suite.ioctx.GetOmapValues(oid, "key1", "", 4)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), map[string][]byte{
		"prefixed-key3": []byte("value3"),
		"key2":          []byte("value2"),
	}, fetched)

	// maxReturn
	fetched, err = suite.ioctx.GetOmapValues(oid, "", "key", 1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), map[string][]byte{
		"key1": []byte("value1"),
	}, fetched)
}

func (suite *RadosTestSuite) TestSetNamespace() {
	suite.SetupConnection()

	// create oid
	oid := suite.GenObjectName()
	bytes_in := []byte("input data")
	err := suite.ioctx.Write(oid, bytes_in, 0)
	assert.NoError(suite.T(), err)

	stat, err := suite.ioctx.Stat(oid)
	assert.Equal(suite.T(), uint64(len(bytes_in)), stat.Size)
	assert.NotNil(suite.T(), stat.ModTime)

	// oid isn't seen in space1 ns
	suite.ioctx.SetNamespace("space1")
	stat, err = suite.ioctx.Stat(oid)
	assert.Equal(suite.T(), err, rados.RadosErrorNotFound)

	// create oid2 in space1 ns
	oid2 := suite.GenObjectName()
	bytes_in = []byte("input data")
	err = suite.ioctx.Write(oid2, bytes_in, 0)
	assert.NoError(suite.T(), err)

	suite.ioctx.SetNamespace("")
	stat, err = suite.ioctx.Stat(oid2)
	assert.Equal(suite.T(), err, rados.RadosErrorNotFound)

	stat, err = suite.ioctx.Stat(oid)
	assert.Equal(suite.T(), uint64(len(bytes_in)), stat.Size)
	assert.NotNil(suite.T(), stat.ModTime)
}

func (suite *RadosTestSuite) TestListAcrossNamespaces() {
	suite.SetupConnection()

	// count objects in pool
	origObjects := 0
	err := suite.ioctx.ListObjects(func(oid string) {
		origObjects++
	})

	// create oid
	oid := suite.GenObjectName()
	bytes_in := []byte("input data")
	err = suite.ioctx.Write(oid, bytes_in, 0)
	assert.NoError(suite.T(), err)

	// create oid2 in space1 ns
	suite.ioctx.SetNamespace("space1")
	oid2 := suite.GenObjectName()
	bytes_in = []byte("input data")
	err = suite.ioctx.Write(oid2, bytes_in, 0)
	assert.NoError(suite.T(), err)

	// count objects in space1 ns
	nsFoundObjects := 0
	err = suite.ioctx.ListObjects(func(oid string) {
		nsFoundObjects++
	})
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), 1, nsFoundObjects)

	// count objects in pool
	suite.ioctx.SetNamespace(rados.RadosAllNamespaces)
	allFoundObjects := 0
	err = suite.ioctx.ListObjects(func(oid string) {
		allFoundObjects++
	})
	assert.NoError(suite.T(), err)
	assert.EqualValues(suite.T(), (origObjects + 2), allFoundObjects)
}

func (suite *RadosTestSuite) TestLocking() {
	suite.SetupConnection()

	oid := suite.GenObjectName()

	// lock ex
	res, err := suite.ioctx.LockExclusive(oid, "myLock", "myCookie", "this is a test lock", 0, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, res)

	// verify lock ex
	info, err := suite.ioctx.ListLockers(oid, "myLock")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(info.Clients))
	assert.Equal(suite.T(), true, info.Exclusive)

	// fail to lock ex again
	res, err = suite.ioctx.LockExclusive(oid, "myLock", "myCookie", "this is a description", 0, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), -17, res)

	// fail to lock sh
	res, err = suite.ioctx.LockShared(oid, "myLock", "myCookie", "", "a description", 0, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), -17, res)

	// unlock
	res, err = suite.ioctx.Unlock(oid, "myLock", "myCookie")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, res)

	// verify unlock
	info, err = suite.ioctx.ListLockers(oid, "myLock")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, len(info.Clients))

	// lock sh
	res, err = suite.ioctx.LockShared(oid, "myLock", "myCookie", "", "a description", 0, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, res)

	// verify lock sh
	info, err = suite.ioctx.ListLockers(oid, "myLock")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(info.Clients))
	assert.Equal(suite.T(), false, info.Exclusive)

	// fail to lock sh again
	res, err = suite.ioctx.LockExclusive(oid, "myLock", "myCookie", "a description", 0, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), -17, res)

	// fail to lock ex
	res, err = suite.ioctx.LockExclusive(oid, "myLock", "myCookie", "this is a test lock", 0, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res, -17)

	// break the lock
	res, err = suite.ioctx.BreakLock(oid, "myLock", info.Clients[0], "myCookie")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, res)

	// verify lock broken
	info, err = suite.ioctx.ListLockers(oid, "myLock")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, len(info.Clients))

	// lock sh with duration
	res, err = suite.ioctx.LockShared(oid, "myLock", "myCookie", "", "a description", time.Millisecond, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, res)

	// verify lock sh expired
	time.Sleep(time.Second)
	info, err = suite.ioctx.ListLockers(oid, "myLock")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, len(info.Clients))

	// lock sh with duration
	res, err = suite.ioctx.LockExclusive(oid, "myLock", "myCookie", "a description", time.Millisecond, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, res)

	// verify lock sh expired
	time.Sleep(time.Second)
	info, err = suite.ioctx.ListLockers(oid, "myLock")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, len(info.Clients))
}

func (suite *RadosTestSuite) TestOmapOnNonexistentObjectError() {
	suite.SetupConnection()
	oid := suite.GenObjectName()
	_, err := suite.ioctx.GetAllOmapValues(oid, "", "", 100)
	assert.Equal(suite.T(), err, rados.RadosErrorNotFound)
}

func TestRadosTestSuite(t *testing.T) {
	suite.Run(t, new(RadosTestSuite))
}
