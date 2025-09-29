// Package main implements the gRPC server
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	grpcAddr = ":50051"
)

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s *server) Create(ctx context.Context, in *desc.CreateRequest) (*desc.CreateResponse, error) {
	fmt.Printf("Create request received: %+v\n", in)

	_ = ctx

	return nil, nil
}

func (s *server) Get(ctx context.Context, in *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Printf("Get request received: %+v\n", in)

	_ = ctx

	return &desc.GetResponse{
		User: &desc.User{
			Id: in.GetId(),
			UserInfo: &desc.UserInfo{
				Name:  "Matvey",
				Email: "Likhanov",
				Role:  desc.Role_ADMIN,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Update(ctx context.Context, in *desc.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Printf("Update request received: %+v\n", in)

	_ = ctx

	return nil, nil
}

func (s *server) Delete(ctx context.Context, in *desc.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("Delete request received: %+v\n", in)

	_ = ctx

	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", grpcAddr) // #nosec G102
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
