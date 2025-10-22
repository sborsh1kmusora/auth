package main

import (
	"context"
	"flag"

	"github.com/sborsh1kmusora/auth/internal/app"
	"github.com/sborsh1kmusora/auth/internal/logger"
	"go.uber.org/zap"
)
 
func main() {
	flag.Parse()

	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("Failed to create application", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		logger.Fatal("Failed to run app", zap.Error(err))
	}
}
