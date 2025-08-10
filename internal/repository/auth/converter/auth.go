package converter

import (
	repoModel "github.com/Ghaarp/auth/internal/repository/auth/model"
	serviceModel "github.com/Ghaarp/auth/internal/service/auth/model"
)

type AuthConverter struct {
}

func CreateConverter() *AuthConverter {
	return &AuthConverter{}
}

func (conv *AuthConverter) ToRepoUserDataPrivate(src *serviceModel.UserDataPrivate) *repoModel.UserDataPrivate {

	return &repoModel.UserDataPrivate{
		Name:     src.Name,
		Email:    src.Email,
		Password: src.Password,
		Role:     int64(src.Role),
	}

}

func (conv *AuthConverter) ToRepoUserDataPublic(src *serviceModel.UserDataPublic) *repoModel.UserDataPublic {

	res := &repoModel.UserDataPublic{
		Id:   src.Id,
		Role: src.Role,
	}

	res.Name.Valid = len(src.Name) != 0
	res.Name.String = src.Name
	res.Email.Valid = len(src.Email) != 0
	res.Email.String = src.Email

	return res
}

func (conv *AuthConverter) ToServiceUserDataPublic(src *repoModel.UserDataPublic) *serviceModel.UserDataPublic {

	return &serviceModel.UserDataPublic{
		Id:    src.Id,
		Name:  src.Name.String,
		Email: src.Email.String,
		Role:  src.Role,
	}

}
