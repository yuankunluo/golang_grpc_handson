package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "yuankunluo/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Count not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)
	// doUnary(c)
	// doServerStream(c)
	doClientStreaming(c)
}

func doUnary(c pb.GreetServiceClient) {
	fmt.Println("Start to do a Unary RPC...")
	req := &pb.GreetRequest{
		Greeting: &pb.Greeting{
			FirstName: "Yuankun",
			LastName:  "Luo",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v\n", err)
	}
	log.Printf("Response from Greet: %v\n", res.Result)
}

func doServerStream(c pb.GreetServiceClient) {
	fmt.Println("Start to do a Server Streaming RPC...")

	req := &pb.GreetManyTimesRequest{
		Greeting: &pb.Greeting{
			FirstName: "Yuankun",
			LastName:  "Luo",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// We've reached the end of the stream.
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doClientStreaming(c pb.GreetServiceClient) {
	fmt.Println("Start to do a Client Streaming RPC...")

	requests := []*pb.LongGreetRequest{
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Yuankun",
				LastName:  "Luo",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Wenqiang",
				LastName:  "Huang",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Xiaofei",
				LastName:  "Ren",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Dengqian",
				LastName:  "Xu",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error with calling LongGreat stream: %v\n", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error with CloseAndRecv LongGreat stream: %v\n", err)
	}

	fmt.Printf("LongGreet response: %v", res)

}
