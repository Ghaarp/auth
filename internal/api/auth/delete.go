package auth

import (
	"context"

	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

func (auth *AuthImplementation) Delete(ctx context.Context, req *generated.DeleteRequest) (*generated.DeleteResponse, error) {
	err := auth.authService.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &generated.DeleteResponse{}, nil
}
