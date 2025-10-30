package app

import (
	"context"
	"net"
	"net/http"
	"os"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sborsh1kmusora/auth/internal/config"
	"github.com/sborsh1kmusora/auth/internal/interceptor"
	"github.com/sborsh1kmusora/auth/internal/logger"
	"github.com/sborsh1kmusora/auth/internal/metrics"
	descAccessV1 "github.com/sborsh1kmusora/auth/pkg/access_v1"
	descAuthV1 "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	descUserV1 "github.com/sborsh1kmusora/auth/pkg/user_v1"
	"github.com/sborsh1kmusora/platform_common/pkg/closer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/natefinch/lumberjack.v2"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	go func() {
		err := a.runPrometheus()
		if err != nil {
			logger.Fatal("Failed to run prometheus server", zap.Error(err))
		}
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context){
		a.initLog,
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initMetrics,
	}

	for _, f := range inits {
		f(ctx)
	}

	return nil
}

func (a *App) initConfig(_ context.Context) {
	config.Load(".env")
}

func (a *App) initMetrics(_ context.Context) {
	metrics.Init()
}

func (a *App) initServiceProvider(_ context.Context) {
	a.serviceProvider = newServiceProvider()
}

func (a *App) initLog(_ context.Context) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	if env == "dev" {
		level := zap.NewAtomicLevelAt(zap.DebugLevel)
		logger.Init(getDevCore(level))
	} else {
		level := zap.NewAtomicLevelAt(zap.InfoLevel)
		logger.Init(getProdCore(level))
	}
}

func (a *App) initGRPCServer(ctx context.Context) {
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
				interceptor.MetricsInterceptor,
			),
		),
	)
	reflection.Register(a.grpcServer)

	descUserV1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAuthV1.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	descAccessV1.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))
}

func (a *App) runGRPCServer() error {
	logger.Info(
		"GRPC server is running on",
		zap.String("address", a.serviceProvider.GRPCConfig().Address()),
	)

	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	if err := a.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    ":2112",
		Handler: mux,
	}

	logger.Info("Prometheus is running on", zap.String("address", prometheusServer.Addr))

	if err := prometheusServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func getDevCore(level zap.AtomicLevel) zapcore.Core {
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)

	return zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
}

func getProdCore(level zap.AtomicLevel) zapcore.Core {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/app/logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	prodCfg := zap.NewProductionEncoderConfig()
	prodCfg.TimeKey = "timestamp"
	prodCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	jsonEncoder := zapcore.NewJSONEncoder(prodCfg)

	return zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), level),
		zapcore.NewCore(jsonEncoder, file, level),
	)
}
