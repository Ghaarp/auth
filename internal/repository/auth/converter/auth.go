package converter

import (
	repo "github.com/Ghaarp/auth/internal/repository/auth/model"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

type AuthConverter struct {
}

func (conv *AuthConverter) ToRepoUserDataPrivate(src *generated.PrivateUser) *repo.UserDataPrivate {

	return &repo.UserDataPrivate{
		Name:     src.Name,
		Email:    src.Email,
		Password: src.Password,
		Role:     int64(src.Role),
	}

}

func (conv *AuthConverter) ToRepoUserDataPublic(src *generated.PublicUser) (res *repo.UserDataPublic) {

	res = &repo.UserDataPublic{}
	res.Id = src.Id
	res.Name.Valid = len(src.Name) == 0
	res.Name.String = src.Name
	res.Email.Valid = len(src.Email) == 0
	res.Email.String = src.Email
	res.Role = int64(src.Role)

	return

}

func (conv *AuthConverter) ToProtoUserDataPublic(src *repo.UserDataPublic) *generated.PublicUser {

	return &generated.PublicUser{
		Id:    src.Id,
		Name:  src.Name.String,
		Email: src.Email.String,
		Role:  generated.Role(src.Role),
	}

}
