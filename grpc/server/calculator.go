package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/acoshift/backend-hardway-ep3/grpc/calculator"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(auth))
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

func (s *calculatorServer) Sum(ctx context.Context, op *pb.Sequence) (*pb.Result, error) {
	var result int64
	for _, x := range op.GetX() {
		result += x
	}
	return &pb.Result{Result: result}, nil
}

func auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	tokens := md.Get("authorization")
	if len(tokens) == 0 || tokens[0] != "super_secret" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	return handler(ctx, req)
}
