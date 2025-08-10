package auth

import (
	"context"

	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

func (auth *AuthImplementation) Update(ctx context.Context, req *generated.UpdateRequest) (*generated.UpdateResponse, error) {
	publicUser := auth.serviceConverter.ToProtoPublicUser(req)
	publicUserService := auth.serviceConverter.ToServiceUserDataPublic(publicUser)
	err := auth.authService.Update(ctx, publicUserService)
	if err != nil {
		return nil, err
	}

	return &generated.UpdateResponse{}, nil
}
