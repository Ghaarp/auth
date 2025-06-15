package auth

import (
	"context"
	"fmt"

	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

func (auth *AuthImplementation) Create(ctx context.Context, req *generated.CreateRequest) (*generated.CreateResponse, error) {
	if auth == nil {
		return &generated.CreateResponse{}, fmt.Errorf("auth implementation is not initialized")
	}
	if auth.serviceConverter == nil {
		return &generated.CreateResponse{}, fmt.Errorf("service converter is not initialized")
	}

	privateUser := auth.serviceConverter.ToApiUserDataPrivate(req)
	id, err := auth.authService.Create(ctx, auth.serviceConverter.ToServiceUserDataPrivate(privateUser))
	if err != nil {
		return nil, err
	}

	return &generated.CreateResponse{
		Id: id,
	}, nil
}
