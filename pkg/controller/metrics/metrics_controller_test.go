package metrics

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc"

	"github.com/opensds/opensds/pkg/dock/client"
	"github.com/opensds/opensds/pkg/model"
	pb "github.com/opensds/opensds/pkg/model/proto"
	. "github.com/opensds/opensds/testutils/collection"
)

type fakeClient struct{}

func NewFakeClient() client.Client {
	return &fakeClient{}
}
func (fc *fakeClient) Connect(edp string) error {
	return nil
}

func (fc *fakeClient) Close() {
	return
}
func NewFakeController() Controller {
	return &controller{
		Client:   NewFakeClient(),
		DockInfo: &model.DockSpec{},
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
	fc := NewFakeController()
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
