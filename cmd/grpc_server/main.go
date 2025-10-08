// Package main implements the gRPC server
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sborsh1kmusora/auth/internal/config"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, in *desc.CreateRequest) (*desc.CreateResponse, error) {
	fmt.Printf("Create request received: %+v\n", in)

	var id int64
	err := s.pool.QueryRow(ctx, `
		insert into user_info (name, email, password, is_admin) 
		values ($1, $2, $3, $4)
		returning id;`,
		in.UserInfo.Name, in.UserInfo.Email, in.Password, in.UserInfo.IsAdmin,
	).Scan(&id)
	if err != nil {
		log.Printf("Unable to insert user info: %v\n", err)
		return nil, status.Error(codes.Internal, "Unable to insert user info")
	}

	log.Printf("Inserted user with id: %v\n", id)

	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, in *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Printf("Get request received: %+v\n", in)

	var (
		id        int64
		name      string
		email     string
		password  string
		isAdmin   bool
		createdAt time.Time
		updatedAt sql.NullTime
	)

	err := s.pool.QueryRow(ctx, `
		select id, name, email, password, is_admin, created_at, updated_at 
		from user_info 
		where id = $1`,
		in.Id,
	).Scan(&id, &name, &email, &password, &isAdmin, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Unable to get user info: %v\n", err)
		return nil, status.Error(codes.NotFound, "Unable to get user info")
	}

	var updatedAtPb *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtPb = timestamppb.New(updatedAt.Time)
	} else {
		updatedAtPb = nil
	}

	return &desc.GetResponse{
		User: &desc.User{
			Id: id,
			UserInfo: &desc.UserInfo{
				Name:    name,
				Email:   email,
				IsAdmin: isAdmin,
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: updatedAtPb,
		},
	}, nil
}

func (s *server) Update(ctx context.Context, in *desc.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Printf("Update request received: %+v\n", in)

	if in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	var exists bool
	err := s.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM user_info WHERE id = $1)
	`, in.Id).Scan(&exists)
	if err != nil {
		log.Printf("Unable to check user existence: %v\n", err)
		return nil, status.Error(codes.Internal, "database error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	updateBuilder := squirrel.Update("user_info").
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": in.Id}).
		PlaceholderFormat(squirrel.Dollar)

	if in.Name != nil && in.Name.Value != "" {
		updateBuilder = updateBuilder.Set("name", in.Name.Value)
	}

	if in.Email != nil && in.Email.Value != "" {
		updateBuilder = updateBuilder.Set("email", in.Email.Value)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		log.Printf("Unable to build SQL query: %v\n", err)
		return nil, status.Error(codes.Internal, "query building failed")
	}

	if len(args) == 1 {
		return &emptypb.Empty{}, nil
	}

	fmt.Printf("Executing query: %s with args: %v\n", query, args)

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Unable to update user info: %v\n", err)
		return nil, status.Error(codes.Internal, "update failed")
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, in *desc.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("Delete request received: %+v\n", in)

	_, err := s.pool.Exec(ctx, `
		delete from user_info
		where id = $1
	`, in.Id)
	if err != nil {
		log.Printf("Unable to delete user info: %v\n", err)
		return nil, status.Error(codes.Internal, "Unable to delete user info")
	}

	return nil, nil
}

func main() {
	flag.Parse()

	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("Unable to load config: %v\n", err)
	}

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

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
