package main

import (
	"context"
	"log"
	"time"

	desc "github.com/fvckinginsxne/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := desc.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	r, err := c.Create(ctx, &desc.CreateRequest{
		UserInfo: &desc.UserInfo{
			Name:  "Matvey",
			Email: "matvey@gmail.com",
			Role:  desc.Role_ADMIN,
		},
	})
	if err != nil {
		log.Fatalf("could not create: %v", err)
	}

	log.Printf("response: %+v", r)
}
