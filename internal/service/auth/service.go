package auth

import (
	"context"

	repositoryDef "github.com/Ghaarp/auth/internal/repository"
	serviceDef "github.com/Ghaarp/auth/internal/service"
	"github.com/Ghaarp/auth/internal/service/auth/model"
)

type AuthService struct {
	repository repositoryDef.AuthRepository
	converter  repositoryDef.RepoConverter
}

func CreateService(r repositoryDef.AuthRepository, c repositoryDef.RepoConverter) serviceDef.AuthService {
	return &AuthService{
		repository: r,
		converter:  c,
	}
}

func (service *AuthService) Create(ctx context.Context, data *model.UserDataPrivate) (int64, error) {
	userData := service.converter.ToRepoUserDataPrivate(data)
	return service.repository.Create(ctx, userData)
}

func (service *AuthService) Get(ctx context.Context, id int64) (*model.UserDataPublic, error) {
	result, err := service.repository.Get(ctx, id)
	return service.converter.ToServiceUserDataPublic(result), err
}

func (service *AuthService) Update(ctx context.Context, data *model.UserDataPublic) error {
	return service.repository.Update(ctx, service.converter.ToRepoUserDataPublic(data))
}

func (service *AuthService) Delete(ctx context.Context, id int64) error {
	return service.repository.Delete(ctx, id)
}

func (service *AuthService) StopService(ctx context.Context) {
	service.repository.ClosePool(ctx)
}
