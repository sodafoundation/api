package metrics

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"google.golang.org/grpc"

	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
)

type fakeClient struct{}

func (fc *fakeClient) GetMetrics(ctx context.Context, in *pb.GetMetricsOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) GetUrls(ctx context.Context, in *pb.NoParams, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func NewFakeClient() client.Client {
	return &fakeClient{}
}
func (fc *fakeClient) Connect(edp string) error {
	return nil
}

func (fc *fakeClient) Close() {
	return
}
func NewFakeController(client2 *http.Client, URL string) Controller {
	return &controller{
		Client:   NewFakeClient(),
		DockInfo: &model.DockSpec{},
		API: &API{
			Client:  client2,
			baseURL: URL,
		},
	}
}
func (fc *fakeClient) CollectMetrics(ctx context.Context, in *pb.CollectMetricsOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteMetrics,
			},
		},
	}, nil

}
func Test_CollectMetrics(t *testing.T) {
	fc := NewFakeController(nil, "")
	retunMetrics, _ := fc.CollectMetrics(&pb.CollectMetricsOpts{})
	expectedMetrics := SamplemetricsSpec
	if !reflect.DeepEqual(expectedMetrics, retunMetrics) {
		t.Errorf("controller.CollectMetrics() = %v, want %v", expectedMetrics, retunMetrics)
	}
}

func Test_logMetricSpec(t *testing.T) {
	type args struct {
		spec *model.MetricSpec
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test1", args: args{spec: SamplemetricsSpec[0]}}, {name: "test1", args: args{spec: SamplemetricsSpec[1]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logMetricSpec(tt.args.spec)
		})
	}
}
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
func Test_controller_GetLatestMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/api/v1/query?query=iops")
		// Send response to be tested
		rw.Write([]byte(ByteGetMetrics))
	}))
	// Close the server when test finishes
	defer server.Close()

	type fields struct {
		Client   client.Client
		DockInfo *model.DockSpec
		API      *API
	}
	type args struct {
		opt *pb.GetMetricsOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.MetricSpec
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				opt: &pb.GetMetricsOpts{
					InstanceId:           "",
					MetricName:           "iops",
					StartTime:            "",
					EndTime:              "",
					Context:              "",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
			want:    SampleGetmetricsSpec,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFakeController(server.Client(), server.URL)
			got, err := c.GetLatestMetrics(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.GetLatestMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("controller.GetLatestMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_controller_GetInstantMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/api/v1/query?query=iops&time=1560169109")
		// Send response to be tested
		rw.Write([]byte(ByteGetMetrics))
	}))
	// Close the server when test finishes
	defer server.Close()

	type fields struct {
		Client   client.Client
		DockInfo *model.DockSpec
		API      *API
	}
	type args struct {
		opt *pb.GetMetricsOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.MetricSpec
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				opt: &pb.GetMetricsOpts{
					InstanceId:           "",
					MetricName:           "iops",
					StartTime:            "1560169109",
					EndTime:              "1560169109",
					Context:              "",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
			want:    SampleGetmetricsSpec,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFakeController(server.Client(), server.URL)
			//c. = API{server.Client(), server.URL}
			got, err := c.GetInstantMetrics(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.GetLatestMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("controller.GetLatestMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_controller_GetRangeMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "/api/v1/query?query=iops&start=1560169109&end=1560169109&step=30")
		// Send response to be tested
		rw.Write([]byte(ByteGetRangeMetrics))
	}))
	// Close the server when test finishes
	defer server.Close()

	type fields struct {
		Client   client.Client
		DockInfo *model.DockSpec
		API      *API
	}
	type args struct {
		opt *pb.GetMetricsOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.MetricSpec
		wantErr bool
	}{
		{
			name:   "test1",
			fields: fields{},
			args: args{
				opt: &pb.GetMetricsOpts{
					InstanceId:           "",
					MetricName:           "iops",
					StartTime:            "1560169109",
					EndTime:              "1560169109",
					Context:              "",
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			},
			want:    SampleGetmetricsRangeSpec,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewFakeController(server.Client(), server.URL)
			//c. = API{server.Client(), server.URL}
			got, err := c.GetRangeMetrics(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.GetLatestMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got[0], tt.want[0]) {
				t.Errorf("controller.GetLatestMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckServiceStatus(t *testing.T) {
	type args struct {
		sName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				sName: "test_telemetry_service",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckServiceStatus(tt.args.sName); (err != nil) != tt.wantErr {
				t.Errorf("CheckServiceStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_controller_GetUrls(t *testing.T) {
	type fields struct {
		Client   client.Client
		DockInfo *model.DockSpec
		API      *API
	}
	tests := []struct {
		name    string
		fields  fields
		want    *map[string]model.UrlDesc
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &controller{
				Client:   tt.fields.Client,
				DockInfo: tt.fields.DockInfo,
				API:      tt.fields.API,
			}
			var grafan_url string = "http://127.0.0.1:3000"
			var alert_mgr_url string = "http://127.0.0.1:9093"
			flag.StringVar(&grafan_url, "grafana-url", "http://127.0.0.1:3000", "Grafana listen endpoint")
			flag.StringVar(&alert_mgr_url, "alertmgr-url", "http://127.0.0.1:9093", "Alert manager listen endpoint")
			flag.Parse()
			_, err := c.GetUrls()
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.GetUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

// Fake Client dummy implementation of all required methods
func (fc *fakeClient) CreateVolume(ctx context.Context, in *pb.CreateVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolume,
			},
		},
	}, nil
}

// dummy
func (fc *fakeClient) DeleteVolume(ctx context.Context, in *pb.DeleteVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) ExtendVolume(ctx context.Context, in *pb.ExtendVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolume,
			},
		},
	}, nil
}

// dummy
func (fc *fakeClient) CreateVolumeAttachment(ctx context.Context, in *pb.CreateVolumeAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteAttachment,
			},
		},
	}, nil
}

func (fc *fakeClient) DeleteVolumeAttachment(ctx context.Context, in *pb.DeleteVolumeAttachmentOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) CreateVolumeSnapshot(ctx context.Context, in *pb.CreateVolumeSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteSnapshot,
			},
		},
	}, nil
}

// dummy
func (fc *fakeClient) DeleteVolumeSnapshot(ctx context.Context, in *pb.DeleteVolumeSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) CreateVolumeGroup(ctx context.Context, in *pb.CreateVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolumeGroup,
			},
		},
	}, nil
}

// dummy
func (fc *fakeClient) UpdateVolumeGroup(ctx context.Context, in *pb.UpdateVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolumeGroup,
			},
		},
	}, nil
}

// dummy
func (fc *fakeClient) DeleteVolumeGroup(ctx context.Context, in *pb.DeleteVolumeGroupOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) AttachVolume(ctx context.Context, in *pb.AttachVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: "",
			},
		},
	}, nil
}

// dummy
func (fc *fakeClient) DetachVolume(ctx context.Context, in *pb.DetachVolumeOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) CreateReplication(ctx context.Context, in *pb.CreateReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteReplication,
			},
		},
	}, nil
}

// dummy
func (fc *fakeClient) DeleteReplication(ctx context.Context, in *pb.DeleteReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) EnableReplication(ctx context.Context, in *pb.EnableReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) DisableReplication(ctx context.Context, in *pb.DisableReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

// dummy
func (fc *fakeClient) FailoverReplication(ctx context.Context, in *pb.FailoverReplicationOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

//dummy fileshare functions
func (c *fakeClient) GetLatestMetrics(opt *pb.GetMetricsOpts) ([]*model.MetricSpec, error) {
	return nil, nil
}
func (fc *fakeClient) CreateFileShare(ctx context.Context, in *pb.CreateFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) CreateFileShareAcl(ctx context.Context, in *pb.CreateFileShareAclOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) DeleteFileShareAcl(ctx context.Context, in *pb.DeleteFileShareAclOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}
func (fc *fakeClient) DeleteFileShare(ctx context.Context, in *pb.DeleteFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) CreateFileShareSnapshot(ctx context.Context, in *pb.CreateFileShareSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}

func (fc *fakeClient) DeleteFileShareSnapshot(ctx context.Context, in *pb.DeleteFileShareSnapshotOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return nil, nil
}
