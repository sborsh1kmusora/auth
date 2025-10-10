package main

import (
	"context"
	"log"

	"github.com/sborsh1kmusora/auth/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer pool.Close()

	userRepo := userRepository.NewRepository(pool)
	userSrv := userService.NewService(userRepo)

	if err := a.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
