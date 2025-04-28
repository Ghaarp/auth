package repository

import (
	repoModel "github.com/Ghaarp/auth/internal/repository/auth/model"
	serviceModel "github.com/Ghaarp/auth/internal/service/auth/model"
)

type RepoConverter interface {
	ToRepoUserDataPrivate(src *serviceModel.UserDataPrivate) *repoModel.UserDataPrivate
	ToRepoUserDataPublic(src *serviceModel.UserDataPublic) *repoModel.UserDataPublic
	ToServiceUserDataPublic(src *repoModel.UserDataPublic) *serviceModel.UserDataPublic
}
