package nfsnative

import (
	//"github.com/ceph/go-ceph/rados"
	//"github.com/ceph/go-ceph/rbd"
	log "github.com/golang/glog"
	. "github.com/opensds/opensds/contrib/drivers/utils/config"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	"github.com/opensds/opensds/pkg/utils/config"
	. "github.com/opensds/opensds/testutils/collection"
)

const (
	opensdsPrefix   = "opensds-"
	sizeShiftBit    = 30
	defaultConfPath = "/etc/opensds/driver/nfsnative.yaml"
	defaultAZ       = "default"
)

const (
	KPoolName  = "NFSNativePoolName"
	KImageName = "NFSNativeImageName"
)

type NFSNativeConfig struct {
	ConfigFile string                    `yaml:"configFile,omitempty"`
	//Pool       map[string]PoolProperties `yaml:"pool,flow"`
}

func EncodeName(id string) string {
	return opensdsPrefix + id
}

/*func NewSrcMgr(conf *NFSNativeConfig) *SrcMgr {
	return &SrcMgr{conf: conf}
}*/

/*
type SrcMgr struct {
	conn  *rados.Conn
	ioctx *rados.IOContext
	img   *rbd.Image
	conf  *CephConfig
}*/

type Driver struct {
	conf *NFSNativeConfig
}

func (d *Driver) Setup() error {
	d.conf = &NFSNativeConfig{ConfigFile: "/etc/nfsnative/nfsnative.conf"}
	p := config.CONF.OsdsDock.Backends.NFSNative.ConfigPath
	if "" == p {
		p = defaultConfPath
	}
	_, err := Parse(d.conf, p)
	return err
}

func (d *Driver) Unset() error { return nil }

func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {

	log.Infof("Create nfs native %s (%s) success.", opt.GetName(), opt.GetId())
	/*return &model.FileShareSpec{
		BaseModel: &model.BaseModel{
			Id: opt.GetId(),
		},
		Name:             opt.GetName(),
		Size:             opt.GetSize(),
		Description:      opt.GetDescription(),
		AvailabilityZone: opt.GetAvailabilityZone(),
		Metadata: map[string]string{
			KPoolName: opt.GetPoolName(),
		},
	}, nil*/
	return &SampleFileShares[0], nil
}
/*
func (d *Driver) CreateFileShare(opt *pb.CreateFileShareOpts) (*model.FileShareSpec, error) {
	// create a volume from snapshot
	//if opt.GetSnapshotId() != "" {
	//	return d.createVolumeFromSnapshot(opt)
	//}
	return d.CreateFileShare(opt)
}*/