package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sborsh1kmusora/auth/internal/app"
	"github.com/sborsh1kmusora/auth/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevel = flag.String("l", "debug", "log level")

func main() {
	flag.Parse()

	ctx := context.Background()

	logger.Init(getCore(getAtomicLevel()))

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("Failed to create application", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		logger.Fatal("Failed to run app", zap.Error(err))
	}
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
