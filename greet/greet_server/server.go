package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (s *server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone func was invoked with streaming req\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Receiving the EOF\n")
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		result := "Hello from server, " + firstName + " " + lastName
		sendErr := stream.Send(&pb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending response: %v\n", err)
			return sendErr
		}
	}

}

func (s *server) GreetWithDeadline(ctx context.Context, req *pb.GreetWithDeadlineRequest) (*pb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline func was invoked with %v\n", req)

	// do some timeout
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			// Client cancles the req.
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.DeadlineExceeded, "The client canceled the request.")
		}
		time.Sleep(1 * time.Second)
	}

	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	response := pb.GreetWithDeadlineResponse{
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
	certFile := "./ssl/server.crt"
	keyFile := "./ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Failed loading certificates: %v\n", sslErr)
		return
	}
	opts := grpc.Creds(creds)

	s := grpc.NewServer(opts)
	// Register Service Server.
	pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}
