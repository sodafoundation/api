module github.com/opensds/opensds

go 1.12

require (
	github.com/LINBIT/godrbdutils v0.0.0-20180425110027-65b98a0f103a
	github.com/astaxie/beego v1.11.1
	github.com/beorn7/perks v1.0.0 // indirect
	github.com/ceph/go-ceph v0.0.0-20170728144007-81e4191e131b
	github.com/coreos/etcd v3.3.11+incompatible
	github.com/go-ini/ini v1.41.0
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.2-0.20190517061210-b285ee9cfc6c
	github.com/gophercloud/gophercloud v0.0.0-20190528082055-3ad89c4ea008
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/prometheus/common v0.3.0 // indirect
	github.com/prometheus/procfs v0.0.0-20190425082905-87a4384529e0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.2.2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.3.0
	golang.org/x/net v0.0.0-20190125091013-d26f9f9a57f3 // indirect
	google.golang.org/genproto v0.0.0-20190128161407-8ac453e89fca // indirect
	google.golang.org/grpc v1.18.0
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	golang.org/x/crypto v0.0.0-20190131182504-b8fe1690c613 => github.com/golang/crypto v0.0.0-20190131182504-b8fe1690c613
	golang.org/x/net v0.0.0-20190125091013-d26f9f9a57f3 => github.com/golang/net v0.0.0-20190125091013-d26f9f9a57f3
	golang.org/x/sys v0.0.0-20190130150945-dca44879d564 => github.com/golang/sys v0.0.0-20190130150945-dca44879d564
	golang.org/x/test v0.3.0 => github.com/golang/test v0.3.0
	google.golang.org/genproto v0.0.0-20190128161407-8ac453e89fca => github.com/google/go-genproto v0.0.0-20190128161407-8ac453e89fca
	google.golang.org/grpc v1.18.0 => github.com/grpc/grpc-go v1.18.0
)
