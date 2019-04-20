package fileshare

import (
	"github.com/opensds/opensds/pkg/context"
	"github.com/opensds/opensds/pkg/filesharedock/client"
	pb "github.com/opensds/opensds/pkg/model/fileshareproto"
	. "github.com/opensds/opensds/testutils/collection"
	"google.golang.org/grpc"
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

// Create a volume
func (fc *fakeClient) CreateVolume(ctx context.Context, in *pb.CreateFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{
				Message: ByteVolume,
			},
		},
	}, nil
}

// Delete a volume
func (fc *fakeClient) DeleteFileShare(ctx context.Context, in *pb.DeleteFileShareOpts, opts ...grpc.CallOption) (*pb.GenericResponse, error) {
	return &pb.GenericResponse{
		Reply: &pb.GenericResponse_Result_{
			Result: &pb.GenericResponse_Result{},
		},
	}, nil
}

