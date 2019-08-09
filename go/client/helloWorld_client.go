//The Go implementation of the gRPC client.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "../protos"
	"google.golang.org/grpc"
)

func getDouble(client pb.HelloWorldClient, value int32) {
	i := &pb.Int{Value: value}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rsp, err := client.Double(ctx, i)
	if err != nil {
		log.Fatalf("%v.GetDouble(_) = _, %v: ", client, err)
	}

	fmt.Printf("double %d => %d\n", i.Value, rsp.Value)
}

func getRange(client pb.HelloWorldClient, value int32) {
	_i := &pb.Int{Value: value}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Range(ctx, _i)
	if err != nil {
		log.Fatalf("%v.GetDouble(_) = _, %v: ", client, err)
	}

	for {
		i, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Range(_) = _, %v", client, err)
		}
		fmt.Printf("range %d => %d\n", _i.Value, i.Value)
	}
}

func getSum(client pb.HelloWorldClient, value int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Sum(ctx)
	if err != nil {
		log.Fatalf("%v.GetDouble(_) = _, %v: ", client, err)
	}

	for i := 0; i < value; i++ {
		_i := &pb.Int{Value: int32(i)}
		if err := stream.Send(_i); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, _i, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	fmt.Printf("Sum range(%d) => %d\n", value, reply.Value)
}

func getDoubleIter(client pb.HelloWorldClient, value int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.DoubleIter(ctx)
	if err != nil {
		log.Fatalf("%v.GetDouble(_) = _, %v: ", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive: %v", err)
			}
			fmt.Printf("double range(%d) => %d\n", value, in.Value*2)
		}
	}()
	for i := 0; i < value; i++ {
		_i := &pb.Int{Value: int32(i)}
		if err := stream.Send(_i); err != nil {
			log.Fatalf("Failed to send: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	serverAddr := "localhost:50051"
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewHelloWorldClient(conn)

	fmt.Println("-------------- Double-------------------")
	getDouble(client, 5)
	fmt.Println("-------------- Range -------------------")
	getRange(client, 5)
	fmt.Println("-------------- Sum----------------------")
	getSum(client, 5)
	fmt.Println("-------------- DoubleIter --------------")
	getDoubleIter(client, 5)
}

// Output
// -------------- Double-------------------
// double 5 => 10
// -------------- Range -------------------
// range 5 => 0
// range 5 => 1
// range 5 => 2
// range 5 => 3
// range 5 => 4
// -------------- Sum----------------------
// Sum range(5) => 10
// -------------- DoubleIter --------------
// double range(5) => 0
// double range(5) => 4
// double range(5) => 8
// double range(5) => 12
// double range(5) => 16
