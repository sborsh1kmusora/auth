package user

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sborsh1kmusora/auth/internal/converter"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	if err := i.userService.Update(ctx, converter.ToUpdateInputFromDesc(req)); err != nil {
		log.Printf("Error updating user: %s\n", err)
		return nil, status.Error(codes.Internal, "Error updating user")
	}

	return &empty.Empty{}, nil
}
