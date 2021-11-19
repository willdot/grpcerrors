package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/willdot/grpcerror/proto"
)

func main() {
	fmt.Println("server started")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterTestServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "failed to serve grpc server"))
	}
}

const (
	addr = "localhost:50051"
)

type Server struct {
}

func (s *Server) DoSomething(ctx context.Context, req *pb.Input) (*pb.Output, error) {
	st := status.New(codes.InvalidArgument, "request not valid. see details for more info")

	// validate request and if there are any invalid fields, add details to the status
	if req.Id == 0 {
		var err error
		st, err = st.WithDetails(&pb.ValidationError{
			Field:  "id",
			Reason: "should be > 0",
		})
		if err != nil {
			log.Fatal("failed to set validation error")
		}
	}

	if req.Name == "" {
		var err error
		st, err = st.WithDetails(&pb.ValidationError{
			Field:  "name",
			Reason: "can't be empty",
		})
		if err != nil {
			log.Fatal("failed to set validation error")
		}
	}

	// if there were any validation errors, return
	if len(st.Details()) > 0 {
		return nil, st.Err()
	}

	return &pb.Output{
		Result: fmt.Sprintf("Name: '%s'\nId: %v", req.Name, req.Id),
	}, nil
}
