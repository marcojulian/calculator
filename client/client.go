package main

import (
	"context"
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
