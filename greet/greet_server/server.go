package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	pb "yuankunluo/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Printf("Greet func was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	response := pb.GreetResponse{
		Result: result,
	}
	return &response, nil
}

func (s *server) GreetManyTimes(req *pb.GreetManyTimesRequest, stream pb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes func was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " " + lastName + " number " + strconv.Itoa(i)
		res := &pb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (s *server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet func was invoked with streaming req %v\n", stream)
	result := "Hello from server, "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// We have finished the client stream
			fmt.Printf("Client sending EOF\n")
			return stream.SendAndClose(&pb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		fmt.Printf("Receive from client: %s %s\n", firstName, lastName)
		result += firstName + " " + lastName
	}

}

func main() {

	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	// Register Service Server.
	pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
