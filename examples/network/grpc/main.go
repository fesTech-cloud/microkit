package main

import (
	"fmt"
	"log"
	"time"

	"github.com/festus/microkit/adapters/grpc"
)

func main() {
	client, err := grpc.NewClient("localhost:50051", 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Example with proper protobuf messages (requires .proto files)
	// req := &pb.HelloRequest{Name: "microkit"}
	// resp := &pb.HelloResponse{}
	// 
	// err = client.Call(ctx, "/hello.HelloService/SayHello", req, resp,
	//     network.WithRetry(3, 100*time.Millisecond, 2*time.Second, 2.0),
	// )

	fmt.Println("gRPC client ready - add your protobuf messages")
}