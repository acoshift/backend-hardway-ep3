package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/acoshift/backend-hardway-ep3/grpc/calculator"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, &calculatorServer{})
	grpcServer.Serve(lis)
}

type calculatorServer struct{}

func (s *calculatorServer) Add(ctx context.Context, op *pb.Operand) (*pb.Result, error) {
	result := op.GetX() + op.GetY()
	return &pb.Result{Result: result}, nil
}

func (s *calculatorServer) Mul(ctx context.Context, op *pb.Operand) (*pb.Result, error) {
	result := op.GetX() * op.GetY()
	return &pb.Result{Result: result}, nil
}

func (s *calculatorServer) Div(ctx context.Context, op *pb.Operand) (*pb.Result, error) {
	if op.GetY() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "divided by zero")
	}
	result := op.GetX() / op.GetY()
	return &pb.Result{Result: result}, nil
}
