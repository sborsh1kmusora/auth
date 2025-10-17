package app

import (
	"context"
	"log"
	"net"

	"github.com/sborsh1kmusora/auth/internal/config"
	descAuthV1 "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	descUserV1 "github.com/sborsh1kmusora/auth/pkg/user_v1"
	"github.com/sborsh1kmusora/platform_common/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context){
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		f(ctx)
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) {
	config.Load(".env")
}

func (a *App) initServiceProvider(ctx context.Context) {
	a.serviceProvider = newServiceProvider()
}

func (a *App) initGRPCServer(ctx context.Context) {
	a.grpcServer = grpc.NewServer()
	reflection.Register(a.grpcServer)

	descUserV1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserImpl(ctx))
	descAuthV1.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))
	//descAccessV1.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessImpl(ctx))
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	if err := a.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
