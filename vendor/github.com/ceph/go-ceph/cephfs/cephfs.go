package cephfs

/*
#cgo LDFLAGS: -lcephfs
#cgo CPPFLAGS: -D_FILE_OFFSET_BITS=64
#include <stdlib.h>
#include <cephfs/libcephfs.h>
*/
import "C"

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"syscall"
	"unsafe"
)

type cephError int

func (e cephError) Error() string {
	if e == 0 {
		return fmt.Sprintf("cephfs: no error given")
	}
	err := syscall.Errno(uint(math.Abs(float64(e))))
	return fmt.Sprintf("cephfs: ret=(%d) %v", e, err)
}

// MountInfo exports ceph's ceph_mount_info from libcephfs.cc
type MountInfo struct {
	mount *C.struct_ceph_mount_info
}

// CreateMount creates a mount handle for interacting with Ceph.
func CreateMount() (*MountInfo, error) {
	mount := &MountInfo{}
	ret := C.ceph_create(&mount.mount, nil)
	if ret != 0 {
		log.Errorf("CreateMount: Failed to create mount")
		return nil, cephError(ret)
	}
	return mount, nil
}

// ReadDefaultConfigFile loads the ceph configuration from the specified config file.
func (mount *MountInfo) ReadDefaultConfigFile() error {
	ret := C.ceph_conf_read_file(mount.mount, nil)
	if ret != 0 {
		log.Errorf("ReadDefaultConfigFile: Failed to read ceph config")
		return cephError(ret)
	}
	return nil
}

// Mount mounts the mount handle.
func (mount *MountInfo) Mount() error {
	ret := C.ceph_mount(mount.mount, nil)
	if ret != 0 {
		log.Errorf("Mount: Failed to mount")
		return cephError(ret)
	}
	return nil
}

// Unmount unmounts the mount handle.
func (mount *MountInfo) Unmount() error {
	ret := C.ceph_unmount(mount.mount)
	if ret != 0 {
		log.Errorf("Unmount: Failed to unmount")
		return cephError(ret)
	}
	return nil
}

// Release destroys the mount handle.
func (mount *MountInfo) Release() error {
	ret := C.ceph_release(mount.mount)
	if ret != 0 {
		log.Errorf("Release: Failed to release mount")
		return cephError(ret)
	}
	return nil
}

// SyncFs synchronizes all filesystem data to persistent media.
func (mount *MountInfo) SyncFs() error {
	ret := C.ceph_sync_fs(mount.mount)
	if ret != 0 {
		log.Errorf("Mount: Failed to sync filesystem")
		return cephError(ret)
	}
	return nil
}

// CurrentDir gets the current working directory.
func (mount *MountInfo) CurrentDir() string {
	cDir := C.ceph_getcwd(mount.mount)
	return C.GoString(cDir)
}

// ChangeDir changes the current working directory.
func (mount *MountInfo) ChangeDir(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_chdir(mount.mount, cPath)
	if ret != 0 {
		log.Errorf("ChangeDir: Failed to change directory")
		return cephError(ret)
	}
	return nil
}

// MakeDir creates a directory.
func (mount *MountInfo) MakeDir(path string, mode uint32) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_mkdir(mount.mount, cPath, C.mode_t(mode))
	if ret != 0 {
		log.Errorf("MakeDir: Failed to make directory %s", path)
		return cephError(ret)
	}
	return nil
}

// RemoveDir removes a directory.
func (mount *MountInfo) RemoveDir(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_rmdir(mount.mount, cPath)
	if ret != 0 {
		log.Errorf("RemoveDir: Failed to remove directory")
		return cephError(ret)
	}
	return nil
}

// Chmod changes the mode bits (permissions) of a file/directory.
func (mount *MountInfo) Chmod(path string, mode uint32) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_chmod(mount.mount, cPath, C.mode_t(mode))
	if ret != 0 {
		log.Errorf("Chmod: Failed to chmod :%s", path)
		return cephError(ret)
	}
	return nil
}

// Chown changes the ownership of a file/directory.
func (mount *MountInfo) Chown(path string, user uint32, group uint32) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	ret := C.ceph_chown(mount.mount, cPath, C.int(user), C.int(group))
	if ret != 0 {
		log.Errorf("Chown: Failed to chown :%s", path)
		return cephError(ret)
	}
	return nil
}

// IsMounted checks mount status.
func (mount *MountInfo) IsMounted() bool {
	ret := C.ceph_is_mounted(mount.mount)
	return ret == 1
}
