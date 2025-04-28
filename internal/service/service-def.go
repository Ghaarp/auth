package service

import (
	"context"

	"github.com/Ghaarp/auth/internal/service/auth/model"
)

type AuthService interface {
	Create(ctx context.Context, data *model.UserDataPrivate) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserDataPublic, error)
	Update(ctx context.Context, data *model.UserDataPublic) error
	Delete(ctx context.Context, id int64) error
	StopService(ctx context.Context)
}
