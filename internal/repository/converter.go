package repository

import (
	"github.com/Ghaarp/auth/internal/repository/auth/model"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

type RepoConverter interface {
	ToRepoUserDataPrivate(src *generated.PrivateUser) *model.UserDataPrivate
	ToRepoUserDataPublic(src *generated.PublicUser) *model.UserDataPublic
	ToProtoUserDataPublic(src *model.UserDataPublic) *generated.PublicUser
}
