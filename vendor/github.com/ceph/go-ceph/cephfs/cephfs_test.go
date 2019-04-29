package cephfs_test

import (
	"fmt"
	"github.com/ceph/go-ceph/cephfs"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
)

var (
	CephMountTest string = "/tmp/ceph/mds/mnt/"
)

func TestCreateMount(t *testing.T) {
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)
}

func TestMountRoot(t *testing.T) {
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)

	err = mount.ReadDefaultConfigFile()
	assert.NoError(t, err)

	err = mount.Mount()
	assert.NoError(t, err)
}

func TestSyncFs(t *testing.T) {
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)

	err = mount.ReadDefaultConfigFile()
	assert.NoError(t, err)

	err = mount.Mount()
	assert.NoError(t, err)

	err = mount.SyncFs()
	assert.NoError(t, err)
}

func TestChangeDir(t *testing.T) {
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)

	err = mount.ReadDefaultConfigFile()
	assert.NoError(t, err)

	err = mount.Mount()
	assert.NoError(t, err)

	dir1 := mount.CurrentDir()
	assert.NotNil(t, dir1)

	err = mount.MakeDir("/asdf", 0755)
	assert.NoError(t, err)

	err = mount.ChangeDir("/asdf")
	assert.NoError(t, err)

	dir2 := mount.CurrentDir()
	assert.NotNil(t, dir2)

	assert.NotEqual(t, dir1, dir2)
	assert.Equal(t, dir1, "/")
	assert.Equal(t, dir2, "/asdf")
}

func TestRemoveDir(t *testing.T) {
	dirname := "one"
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)

	err = mount.ReadDefaultConfigFile()
	assert.NoError(t, err)

	err = mount.Mount()
	assert.NoError(t, err)

	err = mount.MakeDir(dirname, 0755)
	assert.NoError(t, err)

	err = mount.SyncFs()
	assert.NoError(t, err)

	// os.Stat the actual mounted location to verify Makedir/RemoveDir
	_, err = os.Stat(CephMountTest + dirname)
	assert.NoError(t, err)

	err = mount.RemoveDir(dirname)
	assert.NoError(t, err)

	_, err = os.Stat(CephMountTest + dirname)
	assert.EqualError(t, err,
		fmt.Sprintf("stat %s: no such file or directory", CephMountTest+dirname))
}

func TestUnmountMount(t *testing.T) {
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)
	fmt.Printf("%#v\n", mount.IsMounted())

	err = mount.ReadDefaultConfigFile()
	assert.NoError(t, err)

	err = mount.Mount()
	assert.NoError(t, err)
	assert.True(t, mount.IsMounted())

	err = mount.Unmount()
	assert.NoError(t, err)
	assert.False(t, mount.IsMounted())
}

func TestReleaseMount(t *testing.T) {
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)

	err = mount.Release()
	assert.NoError(t, err)
}

func TestChmodDir(t *testing.T) {
	dirname := "two"
	var stats_before uint32 = 0755
	var stats_after uint32 = 0700
	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)

	err = mount.ReadDefaultConfigFile()
	assert.NoError(t, err)

	err = mount.Mount()
	assert.NoError(t, err)

	err = mount.MakeDir(dirname, stats_before)
	assert.NoError(t, err)

	err = mount.SyncFs()
	assert.NoError(t, err)

	// os.Stat the actual mounted location to verify Makedir/RemoveDir
	stats, err := os.Stat(CephMountTest + dirname)
	assert.NoError(t, err)

	assert.Equal(t, uint32(stats.Mode().Perm()), stats_before)

	err = mount.Chmod(dirname, stats_after)
	assert.NoError(t, err)

	stats, err = os.Stat(CephMountTest + dirname)
	assert.Equal(t, uint32(stats.Mode().Perm()), stats_after)
}

// Not cross-platform, go's os does not specifiy Sys return type
func TestChown(t *testing.T) {
	dirname := "three"
	// dockerfile creates bob user account
	var bob uint32 = 1010
	var root uint32 = 0

	mount, err := cephfs.CreateMount()
	assert.NoError(t, err)
	assert.NotNil(t, mount)

	err = mount.ReadDefaultConfigFile()
	assert.NoError(t, err)

	err = mount.Mount()
	assert.NoError(t, err)

	err = mount.MakeDir(dirname, 0755)
	assert.NoError(t, err)

	err = mount.SyncFs()
	assert.NoError(t, err)

	// os.Stat the actual mounted location to verify Makedir/RemoveDir
	stats, err := os.Stat(CephMountTest + dirname)
	assert.NoError(t, err)

	assert.Equal(t, uint32(stats.Sys().(*syscall.Stat_t).Uid), root)
	assert.Equal(t, uint32(stats.Sys().(*syscall.Stat_t).Gid), root)

	err = mount.Chown(dirname, bob, bob)
	assert.NoError(t, err)

	stats, err = os.Stat(CephMountTest + dirname)
	assert.NoError(t, err)
	assert.Equal(t, uint32(stats.Sys().(*syscall.Stat_t).Uid), bob)
	assert.Equal(t, uint32(stats.Sys().(*syscall.Stat_t).Gid), bob)

}
