package user

import (
	"context"
	"log"

	"github.com/sborsh1kmusora/auth/internal/converter"
	desc "github.com/sborsh1kmusora/auth/pkg/auth_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToCreateInputFromDesc(req))
	if err != nil {
		log.Printf("Error creating user: %s\n", err)
		return nil, err
	}

	log.Printf("Created user with id: %d\n", id)

	return &desc.CreateResponse{Id: id}, nil
}
