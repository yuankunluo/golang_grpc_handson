package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "yuankunluo.com/calculator/calculatorpb"
)

type server struct {
}

func (s *server) Sum(context context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	fmt.Printf("Received Sum RPC: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber

	return &pb.SumResponse{
		SumResult: firstNumber + secondNumber,
	}, nil
}

func main() {

	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	// Register Service Server.
	pb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
