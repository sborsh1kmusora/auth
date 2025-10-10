package user

import (
	"context"
	"log"

	"github.com/sborsh1kmusora/auth/internal/converter"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := i.userService.Get(ctx, req.Id)
	if err != nil {
		log.Printf("Error getting user: %s\n", err)
		return nil, status.Error(codes.Internal, "Error getting user")
	}

	log.Printf("getting user %+v\n", userObj)

	return &desc.GetResponse{
		User: converter.ToUserDescFromService(userObj),
	}, nil
}
