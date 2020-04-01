package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (s *server) PrimeNumberDecomposition(req *pb.PrimeNumberDecompositionRequest, stream pb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received PrimeNumberDecomposition RPC: %v\n", req)
	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&pb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v", divisor)
		}
	}
	return nil
}

func (s *server) ComputeAverage(stream pb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Received ComputeAverage RPC stream\n")

	sum := int32(0)
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("stream EOF: %v", err)
			res := &pb.ComputeAverageResponse{
				Average: float64(sum) / float64(count),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("error while receiving stream: %v", err)
		}
		number := req.GetNumber()
		sum += number
		count++
	}

}

func (s *server) FindMaximum(stream pb.CalculatorService_FindMaximumServer) error {
	fmt.Println("Received FindMaximum RPC...")
	maximum := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v\n", err)
			return err
		}
		// Extract number.
		number := req.GetNumber()
		if number > maximum {
			maximum = number
			// Send Maximum back.
			sendErr := stream.Send(&pb.FindMaximumResponse{
				Maximum: maximum,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending maximum to client: %v\n", sendErr)
			}

		}

	}
}

func (s *server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &pb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
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
