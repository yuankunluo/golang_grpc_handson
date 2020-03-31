package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	pb "yuankunluo.com/calculator/calculatorpb"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Count not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)
	doUnary(c)

}

func doUnary(c pb.CalculatorServiceClient) {
	fmt.Println("Start to do a Unary RPC...")
	req := &pb.SumRequest{
		FirstNumber:  1,
		SecondNumber: 2,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v\n", err)
	}
	log.Printf("Response from Greet: %v\n", res.SumResult)
}
