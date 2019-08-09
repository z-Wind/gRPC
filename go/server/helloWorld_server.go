//The Go implementation of the gRPC server.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "../protos"
)

type helloWorldServer struct {
}

func (s *helloWorldServer) Double(ctx context.Context, Int *pb.Int) (*pb.Int, error) {
	rsp := &pb.Int{Value: Int.Value * 2}

	return rsp, nil
}

func (s *helloWorldServer) Range(Int *pb.Int, stream pb.HelloWorld_RangeServer) error {
	var i int32 = 0
	for ; i < Int.Value; i++ {
		rsp := &pb.Int{Value: i}
		if err := stream.Send(rsp); err != nil {
			return err
		}
	}
	return nil
}

func (s *helloWorldServer) Sum(stream pb.HelloWorld_SumServer) error {
	var result int32 = 0
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			rsp := &pb.Int{Value: result}
			return stream.SendAndClose(rsp)
		}
		if err != nil {
			return err
		}
		result += in.Value
	}
}

func (s *helloWorldServer) DoubleIter(stream pb.HelloWorld_DoubleIterServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		rsp := &pb.Int{Value: in.Value * 2}
		if err := stream.Send(rsp); err != nil {
			return err
		}
	}
}

func newServer() *helloWorldServer {
	s := &helloWorldServer{}
	return s
}

func main() {
	port := 50051
	addr := fmt.Sprintf("localhost:%d", port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listen to %s\n", addr)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterHelloWorldServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
