// Package main implements the gRPC server
package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	userApi "github.com/sborsh1kmusora/auth/internal/api/user"
	"github.com/sborsh1kmusora/auth/internal/config"
	userRepository "github.com/sborsh1kmusora/auth/internal/repository/user"
	userService "github.com/sborsh1kmusora/auth/internal/service/user"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	config.Load(configPath)

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("Unable to load config: %v\n", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("Unable to load config: %v\n", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address()) // #nosec G102
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}
	defer pool.Close()

	userRepo := userRepository.NewRepository(pool)
	userSrv := userService.NewService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterAuthV1Server(s, userApi.NewImplementation(userSrv))

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
