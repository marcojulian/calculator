package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/marcojulian/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	log.Println("Hello I'm a client!")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	log.Printf("Created client: %v", c)

	doUnaryCall(c)
	doServerStreamingCall(c)
	doClientStreamingCall(c)
	doBiDiStreamingCall(c)
	doUnaryCallWithDeadline(c, time.Second)
	doErrorUnaryCall(c)
}

func doUnaryCall(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting Sum Unary RPC...")
	req := &calculatorpb.SumRequest{
		Num1: 3,
		Num2: 4,
	}
	log.Printf("Sending req: %v", req)
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}

func doServerStreamingCall(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting PrimeNumberDecomposition Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Num: 120,
	}

	log.Printf("Sending req: %v", req)
	res, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetResult())
	}
}

func doClientStreamingCall(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting ComputeAverage Client Streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		{
			Num: 1,
		},
		{
			Num: 2,
		},
		{
			Num: 3,
		},
		{
			Num: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage RPC: %v", err)
	}

	for _, req := range requests {
		log.Printf("Sending req: %v", req)
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v", err)
	}

	log.Printf("ComputeAverage Response: %v", res)
}

func doBiDiStreamingCall(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting FindMaximum BiDi Streaming RPC...")

	requests := []*calculatorpb.FindMaximumRequest{
		{
			Num: 1,
		},
		{
			Num: 5,
		},
		{
			Num: 3,
		},
		{
			Num: 6,
		},
		{
			Num: 2,
		},
		{
			Num: 20,
		},
	}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while calling BiDi RPC: %v", err)
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func(requests []*calculatorpb.FindMaximumRequest) {
		for _, req := range requests {
			log.Printf("Sending req: %v", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}(requests)

	// we receive a bunch of messages from the client (go routine)
	go func(waitc chan struct{}) {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			log.Printf("Received: %v", res.GetResult())
			time.Sleep(time.Second)
		}
		close(waitc)
	}(waitc)

	// block until everything is done
	<-waitc
}

func doErrorUnaryCall(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting SquareRoot Unary RPC...")
	req := &calculatorpb.SquareRootRequest{
		Num: -2,
	}
	log.Printf("Sending req: %v", req)
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			log.Printf("Error message: %v", resErr.Message())
			log.Fatalf("Error code: %v", resErr.Code())
		} else {
			log.Fatalf("Error while SquareRoot Sum RPC: %v", err)
		}
	}
	log.Printf("Response SquareRoot Sum: %v", res.Result)
}

func doUnaryCallWithDeadline(c calculatorpb.CalculatorServiceClient, timeout time.Duration) {
	log.Println("Starting Multiply Unary RPC...")
	req := &calculatorpb.MultiplyRequest{
		Multiplicand: 4,
		Multiplier:   3,
	}
	log.Printf("Sending req: %v", req)
	log.Printf("With timeout: %v", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	res, err := c.Multiply(ctx, req)
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok && resErr.Code() == codes.DeadlineExceeded {
			log.Println("Deadline was exceeded")
		} else {
			log.Fatalf("Error while calling Multiply RPC: %v", err)
		}
	} else {
		log.Printf("Response from Multiply: %v", res.Result)
	}
}
