package testutil

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	bufSize = 1024 * 1024
)

// GrpcTester provide helper method to setup grpc test environment
type GrpcTester struct {
	GRPConn *grpc.ClientConn
}

// SetupGrpcServer setup grpc server and grpc client connection
func (tester *GrpcTester) SetupGrpcServer(serv *grpc.Server) error {
	lis := newBufferListner(bufSize)

	var err error

	go func(lis *bufferListner) {
		if err := serv.Serve(lis); err != nil {
			fmt.Printf("testutil: grpc server setup failed, err: %+v", err)
		}
	}(lis)

	if err != nil {
		return err
	}

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(lis.Dialer), grpc.WithInsecure())
	if err != nil {
		return err
	}
	tester.GRPConn = conn
	return nil
}

type bufferListner struct {
	*bufconn.Listener
}

func newBufferListner(size int64) *bufferListner {
	return &bufferListner{
		Listener: bufconn.Listen(int(size)),
	}
}

func (lis bufferListner) Dialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
