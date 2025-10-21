package app

import (
	"context"

	"github.com/sborsh1kmusora/auth/internal/api/access"
	"github.com/sborsh1kmusora/auth/internal/api/auth"
	"github.com/sborsh1kmusora/auth/internal/api/user"
	"github.com/sborsh1kmusora/auth/internal/config"
	"github.com/sborsh1kmusora/auth/internal/logger"
	accessRepo "github.com/sborsh1kmusora/auth/internal/repository/access"
	userRepo "github.com/sborsh1kmusora/auth/internal/repository/user"
	accessService "github.com/sborsh1kmusora/auth/internal/service/access"
	authService "github.com/sborsh1kmusora/auth/internal/service/auth"
	userService "github.com/sborsh1kmusora/auth/internal/service/user"
	"github.com/sborsh1kmusora/platform_common/pkg/closer"
	"github.com/sborsh1kmusora/platform_common/pkg/db"
	"github.com/sborsh1kmusora/platform_common/pkg/db/pg"
	"github.com/sborsh1kmusora/platform_common/pkg/db/transaction"
	"go.uber.org/zap"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig
	authConfig config.AuthConfig

	dbClient   db.Client
	txManager  db.TxManager
	userRepo   userRepo.Repository
	accessRepo accessRepo.Repository

	userService   userService.Service
	authService   authService.Service
	accessService accessService.Service

	userImpl   *user.Implementation
	authImpl   *auth.Implementation
	accessImpl *access.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			logger.Fatal("Failed to initialize grpc config", zap.Error(err))
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			logger.Fatal("Failed to initialize pg config", zap.Error(err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			logger.Fatal("Failed to initialize auth config", zap.Error(err))
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			logger.Fatal("Failed to initialize db client", zap.Error(err))
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			logger.Fatal("Failed to ping db", zap.Error(err))
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) userRepo.Repository {
	if s.userRepo == nil {
		s.userRepo = userRepo.NewRepository(s.DBClient(ctx))
	}

	return s.userRepo
}

func (s *serviceProvider) AccessRepository(ctx context.Context) accessRepo.Repository {
	if s.accessRepo == nil {
		s.accessRepo = accessRepo.NewRepository(s.DBClient(ctx))
	}

	return s.accessRepo
}

func (s *serviceProvider) UserService(ctx context.Context) userService.Service {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) AuthService(ctx context.Context) authService.Service {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthConfig(),
			s.UserRepository(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) AccessService(ctx context.Context) accessService.Service {
	if s.accessService == nil {
		s.accessService = accessService.NewService(
			s.AuthConfig(),
			s.AccessRepository(ctx),
		)
	}

	return s.accessService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) AccessImpl(ctx context.Context) *access.Implementation {
	if s.accessImpl == nil {
		s.accessImpl = access.NewImplementation(s.AccessService(ctx))
	}

	return s.accessImpl
}
