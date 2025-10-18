package access

import (
	"github.com/sborsh1kmusora/auth/internal/service/access"
	desc "github.com/sborsh1kmusora/auth/pkg/access_v1"
)

type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService access.Service
}

func NewImplementation(accessService access.Service) *Implementation {
	return &Implementation{
		accessService: accessService,
	}
}
