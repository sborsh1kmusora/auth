package access

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sborsh1kmusora/auth/internal/utils"
	"google.golang.org/grpc/metadata"
)

const authPrefix = "Bearer "

var accessibleRoles map[string]string

func (s *serv) Check(ctx context.Context, endpointAddress string) (bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return false, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return false, errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := utils.VerifyToken(accessToken, s.authConfig.AccessTokenSecretKey())
	if err != nil {
		fmt.Printf("Error verifying access token: %v\n", err)
		return false, err
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		fmt.Printf("Error getting accessible roles: %v\n", err)
		return false, err
	}

	role, ok := accessibleMap[endpointAddress]
	if !ok {
		fmt.Printf("Accessible role not found for %s\n", endpointAddress)
		return true, nil
	}

	if role == claims.Role {
		return true, nil
	}

	return false, errors.New("access denied")
}

func (s *serv) accessibleRoles(ctx context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		accessInfo, err := s.accessRepo.GetList(ctx)
		if err != nil {
			return nil, err
		}

		for _, info := range accessInfo {
			accessibleRoles[info.EndpointAddress] = info.Role
		}
	}

	return accessibleRoles, nil
}
