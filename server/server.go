package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/marcojulian/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Printf("Sum function was invoked with %v", req)
	return &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("PrimeNumberDecomposition function was invoked with %v", req)
	var prime int32 = 2
	n := req.GetNum()
	for n > 1 {
		if n%prime == 0 {
			n = n / prime
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: prime,
			}
			stream.Send(res)
			time.Sleep(time.Second)
		} else {
			prime = prime + 1
		}
	}
	return nil
}

func main() {
	log.Println("Hello I'm a server!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
