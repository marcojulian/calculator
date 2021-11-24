package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"github.com/marcojulian/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		} else {
			prime = prime + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	log.Println("ComputeAverage function was invoked")
	count := 0
	acumulator := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: float64(acumulator) / float64(count),
			})
		} else if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		acumulator = acumulator + int(req.GetNum())
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	log.Println("FindMaximum function was invoked")
	currMax := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		num := req.GetNum()
		if int(num) > currMax {
			currMax = int(num)
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Result: num,
			})
			if err != nil {
				log.Fatalf("Error while sending data to the client: %v", err)
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Println("SquareRoot function was invoked")

	num := req.GetNum()
	if num < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received negative number: %v", num),
		)
	}

	return &calculatorpb.SquareRootResponse{
		Result: math.Sqrt(float64(num)),
	}, nil
}

func (*server) Multiply(ctx context.Context, req *calculatorpb.MultiplyRequest) (*calculatorpb.MultiplyResponse, error) {
	log.Printf("Multiply function was invoked with %v", req)
	return &calculatorpb.MultiplyResponse{
		Result: req.GetMultiplicand() * req.GetMultiplier(),
	}, nil
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
