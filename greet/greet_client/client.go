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
	// doClientStreaming(c)
	doBiDiStreaming(c)
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

func doBiDiStreaming(c pb.GreetServiceClient) {
	fmt.Println("Start to do a BiDi Streaming RPC...")

	requests := []*pb.GreetEveryoneRequest{
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Yuankun",
				LastName:  "Luo",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Wenqiang",
				LastName:  "Huang",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Xiaofei",
				LastName:  "Ren",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Dengqian",
				LastName:  "Xu",
			},
		},
	}

	// create a stream by invoking the client.
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v\n", err)
	}
	// a wait channel.
	waitC := make(chan struct{})
	// send a bunch of messages to the server (go routines).
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message %v to server.\n", req)
			sendErr := stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
			if sendErr != nil {
				log.Fatal("Error while sending to server: %v\n", sendErr)
			}
		}
		stream.CloseSend()
	}()
	// receive a bunch of messages from the server (go routine).
	go func() {
		// function to receive response
		res, err := stream.Recv()
		if err == io.EOF {
			close(waitC)
		}
		if err != nil {
			log.Fatalf("Error while receiving from server: %v\n", err)
		}
		fmt.Printf("Receiving %v", res.GetResult())
	}()
	// block until everyting is done.
	<-waitC
}
