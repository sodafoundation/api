package ceph

import (
	"testing"
	"unsafe"
	"errors"
	"github.com/ceph/go-ceph/rbd"
	"github.com/ceph/go-ceph/rados"
	"github.com/bouk/monkey"
	"github.com/satori/go.uuid"
)

func TestCreateVolume(t *testing.T) {

	defer monkey.UnpatchAll()
	monkey.Patch((*ImageMgr).Init, func(img *ImageMgr) error {return nil })
	monkey.Patch(rbd.Create, func(ioctx *rados.IOContext, name string, size uint64, order int,
		args ...uint64)  (*rbd.Image,  error) {return nil, nil})
	monkey.Patch((*rados.Conn).Shutdown, func(c *rados.Conn) {})
	monkey.Patch((*rados.IOContext).Destroy, func(ioctx *rados.IOContext) {})

	// case 1
	plugin := CephPlugin{}
	resp, err := plugin.CreateVolume("volume001", 1)
	if err != nil {
			t.Errorf("Test Create volume error")
	}
	if resp.Size != 1{
		t.Errorf("Test Create volume size error")
	}
	if resp.Name != "volume001" {
		t.Errorf("Test Create volume name error")
	}
	if _,err = uuid.FromString(resp.Id); err != nil {
		t.Errorf("Test Create volume uuid error")
	}

	//case 2
	monkey.Unpatch((*ImageMgr).Init)
	monkey.Patch((*ImageMgr).Init, func(img *ImageMgr) error {
		return errors.New("Fake error")
	})
	plugin = CephPlugin{}
	_, err = plugin.CreateVolume("volume001", 1)
	if err == nil {
		t.Errorf("Test Create volume error")
	}

	//case 3
	monkey.Unpatch(rbd.Create)
	monkey.Patch(rbd.Create, func(ioctx *rados.IOContext, name string, size uint64, order int,
		args ...uint64)  (*rbd.Image,  error) {
		return nil,  errors.New("Fake error")
	})
	_, err = plugin.CreateVolume("volume001", 1)
	if err == nil {
		t.Errorf("Test Create volume error")
	}
}

func TestGetVolume(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch((*ImageMgr).Init, func(img *ImageMgr) error {
		return nil
	})
	monkey.Patch(rbd.GetImageNames, func(ioctx *rados.IOContext) (names []string, err error){
		nameList := []string{OPENSDS_PREFIX + ":volume001:7ee11866-1f40-4f3c-b093-7a3684523a19",}
		return nameList, nil
	})
	monkey.Patch((*rbd.Image).GetSize, func(r *rbd.Image) (size uint64, err error){
		return 1<<SIZE_SHIFT_BIT, nil
	})
	monkey.Patch((*rbd.Image).Open, func(r *rbd.Image, args ...interface{}) error{
		return nil
	})
	monkey.Patch((*rbd.Image).Close, func(r *rbd.Image) error{
		return nil
	})
	monkey.Patch((*rados.Conn).Shutdown, func(c *rados.Conn) {})
	monkey.Patch((*rados.IOContext).Destroy, func(ioctx *rados.IOContext) {})

	// case 1
	plugin := CephPlugin{}
	resp, err := plugin.GetVolume("7ee11866-1f40-4f3c-b093-7a3684523a19")
	if err != nil {
		t.Errorf("Test Get volume error")
	}
	if resp.Size != 1{
		t.Errorf("Test Get volume size error")
	}
	if resp.Name != "volume001" {
		t.Errorf("Test Get volume name error")
	}
	if resp.Id != "7ee11866-1f40-4f3c-b093-7a3684523a19" {
		t.Errorf("Test Get volume uuid error")
	}

	resp, err = plugin.GetVolume("11111111-1111-1111-1111-111111111111")
	if err != rbd.RbdErrorNotFound {
		t.Errorf("Test Get volume error")
	}
}

func TestDeleteVolme(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch((*ImageMgr).Init, func(img *ImageMgr) error {
		return nil
	})
	monkey.Patch(rbd.GetImageNames, func(ioctx *rados.IOContext) (names []string, err error){
		nameList := []string{OPENSDS_PREFIX + ":volume001:7ee11866-1f40-4f3c-b093-7a3684523a19",}
		return nameList, nil
	})

	monkey.Patch((*rbd.Image).GetSize, func(r *rbd.Image) (size uint64, err error){
		return 1<<SIZE_SHIFT_BIT, nil
	})
	monkey.Patch((*rbd.Image).Remove, func(r *rbd.Image) error{
		return nil
	})
	monkey.Patch((*rados.Conn).Shutdown, func(c *rados.Conn) {})
	monkey.Patch((*rados.IOContext).Destroy, func(ioctx *rados.IOContext) {})

	// case 1
	plugin := CephPlugin{}
	err := plugin.DeleteVolume("7ee11866-1f40-4f3c-b093-7a3684523a19")
	if err != nil {
		t.Errorf("Test Delete volume error")
	}
}

func TestAttachVolume(t *testing.T) {
	plugin := CephPlugin{}
	err := plugin.AttachVolume("7ee11866-1f40-4f3c-b093-7a3684523a19",
		"opensds-server","/mnt")
	if err != nil {
		t.Errorf("Test attach volume error")
	}
}

func TestDetachVolume(t *testing.T) {
	plugin := CephPlugin{}
	err := plugin.DetachVolume("7ee11866-1f40-4f3c-b093-7a3684523a19")
	if err != nil {
		t.Errorf("Test detach volume error")
	}
}

func TestCreateSnapshot(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch((*ImageMgr).Init, func(img *ImageMgr) error {
		return nil
	})
	monkey.Patch(rbd.GetImageNames, func(ioctx *rados.IOContext) (names []string, err error){
		nameList := []string{OPENSDS_PREFIX + ":volume001:7ee11866-1f40-4f3c-b093-7a3684523a19",}
		return nameList, nil
	})

	monkey.Patch((*rbd.Image).GetSize, func(r *rbd.Image) (size uint64, err error){
		return 1<<SIZE_SHIFT_BIT, nil
	})
	//
	type Snapshot struct {
		image *rbd.Image
		name  string
	}
	monkey.Patch((*rbd.Image).CreateSnapshot, func(image *rbd.Image, snapname string) (*rbd.Snapshot, error) {
		snapshot := &rbd.Snapshot{}
		p := (*Snapshot)(unsafe.Pointer(snapshot))
		p.name = "snapshot001"
		p.image = nil
		return snapshot, nil
	})

	monkey.Patch((*rbd.Image).Open, func(r *rbd.Image, args ...interface{}) error{return nil })
	monkey.Patch((*rbd.Image).Close, func(r *rbd.Image) error{return nil })
	monkey.Patch((*rados.Conn).Shutdown, func(c *rados.Conn) {})
	monkey.Patch((*rados.IOContext).Destroy, func(ioctx *rados.IOContext) {})

	// case 1
	plugin := CephPlugin{}
	resp, err := plugin.CreateSnapshot("snapshot001", "7ee11866-1f40-4f3c-b093-7a3684523a19",
		"unite test")
	if err != nil {
		t.Errorf("Test Create snapshot error")
	}
	if resp.Name != "snapshot001" {
		t.Errorf("Test Create snapshot name error")
	}
	if resp.VolumeId != "7ee11866-1f40-4f3c-b093-7a3684523a19" {
		t.Errorf("Test Create snapshot name error")
	}
	if _,err = uuid.FromString(resp.Id); err != nil {
		t.Errorf("Test Create snapshot error")
	}
}

func TestGetSnapshot(t *testing.T) {
	defer monkey.UnpatchAll()
	monkey.Patch((*ImageMgr).Init, func(img *ImageMgr) error {
		return nil
	})
	monkey.Patch(rbd.GetImageNames, func(ioctx *rados.IOContext) (names []string, err error){
		nameList := []string{OPENSDS_PREFIX + ":volume001:7ee11866-1f40-4f3c-b093-7a3684523a19",}
		return nameList, nil
	})
	monkey.Patch((*rbd.Image).GetSnapshotNames, func(*rbd.Image) (snaps []rbd.SnapInfo, err error){
		snaps = make([]rbd.SnapInfo, 1)
		snaps[0] = rbd.SnapInfo{Id: uint64(1),
			Size: uint64(1<<SIZE_SHIFT_BIT),
			Name: OPENSDS_PREFIX + ":snapshot001:25f5d7a2-553d-4d6c-904d-179a9e698cf8",
		}
		return snaps, nil
	})
	monkey.Patch((*rbd.Image).GetSize, func(r *rbd.Image) (size uint64, err error){
		return 1<<SIZE_SHIFT_BIT, nil
	})

	type Snapshot struct {
		image *rbd.Image
		name  string
	}
	monkey.Patch((*rbd.Image).CreateSnapshot, func(image *rbd.Image, snapname string) (*rbd.Snapshot, error){
		snapshot := &rbd.Snapshot{}
		p := (*Snapshot)(unsafe.Pointer(snapshot))
		p.name = snapname
		p.image = image
		return snapshot, nil
	})

	monkey.Patch((*rbd.Image).Open, func(r *rbd.Image, args ...interface{}) error{return nil })
	monkey.Patch((*rbd.Image).Close, func(r *rbd.Image) error{return nil })
	monkey.Patch((*rados.Conn).Shutdown, func(c *rados.Conn) {})
	monkey.Patch((*rados.IOContext).Destroy, func(ioctx *rados.IOContext) {})

	// case 1
	plugin := CephPlugin{}
	resp, err := plugin.GetSnapshot("25f5d7a2-553d-4d6c-904d-179a9e698cf8")
	if err != nil {
		t.Errorf("Test Get snapshot error")
	}
	if resp.Name != "snapshot001"{
		t.Errorf("Test Get snapshot name error")
	}
	if resp.Size != 1{
		t.Errorf("Test Get snapshot size error")
	}

	// case 2
	_, err = plugin.GetSnapshot("11111111-1111-1111-1111-111111111111")
	if err != rbd.RbdErrorNotFound {
		t.Errorf("Test Get snapshot error")
	}
}

func TestDeleteSnapshot(t *testing.T) {


	defer monkey.UnpatchAll()
	monkey.Patch((*ImageMgr).Init, func(img *ImageMgr) error {
		return nil
	})
	monkey.Patch(rbd.GetImageNames, func(ioctx *rados.IOContext) (names []string, err error){
		nameList := []string{OPENSDS_PREFIX + ":volume001:7ee11866-1f40-4f3c-b093-7a3684523a19",}
		return nameList, nil
	})
	monkey.Patch((*rbd.Image).GetSnapshotNames, func(*rbd.Image) (snaps []rbd.SnapInfo, err error){
		snaps = make([]rbd.SnapInfo, 1)
		snaps[0] = rbd.SnapInfo{Id: uint64(1),
			Size: uint64(1<<SIZE_SHIFT_BIT),
			Name: OPENSDS_PREFIX + ":snapshot001:25f5d7a2-553d-4d6c-904d-179a9e698cf8",
		}
		return snaps, nil
	})
	monkey.Patch((*rbd.Image).GetSize, func(r *rbd.Image) (size uint64, err error){
		return 1<<SIZE_SHIFT_BIT, nil
	})

	type Snapshot struct {
		image *rbd.Image
		name  string
	}
	monkey.Patch((*rbd.Snapshot).Remove, func(*rbd.Snapshot) error{
		return  nil
	})

	monkey.Patch((*rbd.Image).Open, func(r *rbd.Image, args ...interface{}) error{return nil })
	monkey.Patch((*rbd.Image).Close, func(r *rbd.Image) error{return nil })
	monkey.Patch((*rados.Conn).Shutdown, func(c *rados.Conn) {})
	monkey.Patch((*rados.IOContext).Destroy, func(ioctx *rados.IOContext) {})

	// case 1
	plugin := CephPlugin{}
	err := plugin.DeleteSnapshot("25f5d7a2-553d-4d6c-904d-179a9e698cf8")
	if err != nil {
		t.Errorf("Test Delete snapshot error")
	}
}
