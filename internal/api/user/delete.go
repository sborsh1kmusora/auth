package user

import (
	"context"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	if err := i.userService.Delete(ctx, req.Id); err != nil {
		log.Printf("Error deleting user: %s\n", err)
		return nil, status.Error(codes.Internal, "Error deleting user")
	}

	return &empty.Empty{}, nil
}
