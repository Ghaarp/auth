package repository

import (
	"context"

	"github.com/Ghaarp/auth/internal/repository/auth/model"
)

type AuthRepository interface {
	Create(ctx context.Context, data *model.UserDataPrivate) (int64, error)
	Get(ctx context.Context, id int64) (*model.UserDataPublic, error)
	Update(ctx context.Context, data *model.UserDataPublic) error
	Delete(ctx context.Context, id int64) error
	ClosePool(ctx context.Context)
}
