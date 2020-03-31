package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"yuankunluo/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet func was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	response := greetpb.GreetResponse{
		Result: result,
	}
	return &response, nil
}

func main() {

	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	// Register Service Server.
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
