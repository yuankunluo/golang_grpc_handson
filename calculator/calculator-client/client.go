package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	doErrorUnary(c, 16)
	doErrorUnary(c, -16)
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

func doClientStreaming(c pb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a ComputeAverage Client Streaming RPC...\n")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening a stream: %v\n", err)
	}

	numbers := []int32{3, 5, 6, 78, 87}

	for _, num := range numbers {
		err := stream.Send(&pb.ComputeAverageRequest{
			Number: num,
		})
		fmt.Printf("Sending %v to server.\n", num)
		time.Sleep(1000 * time.Millisecond)
		if err != nil {
			log.Fatalf("Error while sending numbers: %v\n", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v\n", err)
	}
	fmt.Printf("The Average of %v is %v\n", numbers, res.GetAverage())
}

func doBiDiStreaming(c pb.CalculatorServiceClient) {
	fmt.Printf("Starting to do a doBiDiStreaming FindMaximum Streaming RPC...\n")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v\n", err)
	}

	waitC := make(chan struct{})

	// send routine.
	go func() {
		numbers := []int32{4, 7, 2, 18, 4, 6, 32, 5, 28, 100}

		for _, num := range numbers {
			sendErr := stream.Send(&pb.FindMaximumRequest{
				Number: num,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending number %d to server: %v.\n", num, sendErr)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive routine.
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reading server stream: %v\n", err)
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received a new maximum: %v\n", maximum)
		}
		close(waitC)
	}()

	// wait.
	<-waitC
}

func doErrorUnary(c pb.CalculatorServiceClient, number int32) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	// correct call.
	res, err := c.SquareRoot(context.Background(),
		&pb.SquareRootRequest{
			Number: number})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from grpc (user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			// framework error!
			log.Fatalf("Big Error calling SquareRoot: %v\n", err)
		}
	}

	fmt.Printf("Result of square root of %v is %v\n", number, res.GetNumberRoot())
}
