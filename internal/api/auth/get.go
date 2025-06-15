package auth

import (
	"context"

	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

func (auth *AuthImplementation) Get(ctx context.Context, req *generated.GetRequest) (*generated.GetResponse, error) {
	publicUser, err := auth.authService.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return auth.serviceConverter.ToGetResponse(publicUser), nil
}
