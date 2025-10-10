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
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
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
