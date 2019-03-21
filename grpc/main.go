package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/acoshift/backend-hardway-ep3/grpc/calculator"
)

// go get -u google.golang.org/grpc
// brew install protobuf
// go get -u github.com/golang/protobuf/protoc-gen-go
// mkdir -p calculator
// protoc calculator.proto --go_out=plugins=grpc:calculator
func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewCalculatorClient(conn)
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "super_secret"))

	fmt.Println("4 + 3 =")
	result, err := client.Add(ctx, &pb.Operand{X: 4, Y: 3})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result.Result)
	}

	fmt.Println("4 * 3 =")
	result, err = client.Mul(ctx, &pb.Operand{X: 4, Y: 3})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result.Result)
	}

	fmt.Println("1024 / 8 =")
	result, err = client.Div(ctx, &pb.Operand{X: 1024, Y: 8})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result.Result)
	}

	fmt.Println("1 / 0 =")
	result, err = client.Div(ctx, &pb.Operand{X: 1, Y: 0})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result.Result)
	}
}
