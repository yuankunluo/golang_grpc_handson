package main

import (
	"context"
	"fmt"
	"log"

	"yuankunluo/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Count not connect: %v\n", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Start to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
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
