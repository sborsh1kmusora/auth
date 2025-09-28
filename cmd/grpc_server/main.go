package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/fvckinginsxne/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	grpcPort = ":50051"
)

func (s *server) Create(ctx context.Context, in *desc.CreateRequest) (*desc.CreateResponse, error) {
	fmt.Printf("Create request received: %+v\n", in)

	return nil, nil
}

func (s *server) Get(ctx context.Context, in *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Printf("Get request received: %+v\n", in)

	return nil, nil
}

func (s *server) Update(ctx context.Context, in *desc.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Printf("Update request received: %+v\n", in)

	return nil, nil
}

func (s *server) Delete(ctx context.Context, in *desc.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("Delete request received: %+v\n", in)

	return nil, nil
}

type server struct {
	desc.UnimplementedAuthV1Server
}

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
