package main

import (
	"context"
	"io"
	"log"

	"github.com/marcojulian/calculator/calculatorpb"
	"google.golang.org/grpc"
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
}

func doUnaryCall(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting Sum Unary RPC...")
	req := &calculatorpb.SumRequest{
		Num1: 3,
		Num2: 4,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v + %v = %v", req.Num1, req.Num2, res.Result)
}

func doServerStreamingCall(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting PrimeNumberDecomposition Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Num: 120,
	}

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
