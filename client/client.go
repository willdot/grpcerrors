package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	pb "github.com/willdot/grpcerror/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTestClient(conn)

	emptyRequest(c)

	oneInvalidFieldRequest(c)

	validRequest(c)
}

func emptyRequest(client pb.TestClient) {
	fmt.Println("making request without any valid fields")
	_, err := client.DoSomething(context.Background(), &pb.Input{})
	if err != nil {
		extractErrorDetails(err)
	}
}

func oneInvalidFieldRequest(client pb.TestClient) {
	fmt.Println("making request with only 1 valid field")
	_, err := client.DoSomething(context.Background(), &pb.Input{Id: 1})
	if err != nil {
		extractErrorDetails(err)
	}
}

func validRequest(client pb.TestClient) {
	fmt.Println("making valid request")
	resp, err := client.DoSomething(context.Background(), &pb.Input{Id: 1, Name: "Will"})
	if err != nil {
		extractErrorDetails(err)
	}

	fmt.Println(resp.Result)
}

func extractErrorDetails(err error) {
	fmt.Println(err.Error())

	st := status.Convert(err)
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *pb.ValidationError:
			fmt.Printf("error with field '%s': %s\n", t.Field, t.Reason)
		}
	}
}
