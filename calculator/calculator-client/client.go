package main

import (
	"context"
	"fmt"
	"io"
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
	doServerStreaming(c)
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

func doServerStreaming(c pb.CalculatorServiceClient) {
	fmt.Println("Start to do a PrimeDecomposition streaming RPC...")
	req := &pb.PrimeNumberDecompositionRequest{
		Number: 1239039284012,
	}
	stream, err := c.PrimeNumberDecomposition(
		context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Printf("PrimeDecompo Prime Factor: %v\n", res.GetPrimeFactor())
	}
}
