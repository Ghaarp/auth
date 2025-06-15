package service

import (
	"github.com/Ghaarp/auth/internal/service/auth/model"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

type ServiceConverter interface {
	ToServiceUserDataPrivate(src *generated.PrivateUser) *model.UserDataPrivate
	ToServiceUserDataPublic(src *generated.PublicUser) *model.UserDataPublic
	ToProtoUserDataPublic(src *model.UserDataPublic) *generated.PublicUser
	ToApiUserDataPrivate(data *generated.CreateRequest) *generated.PrivateUser
	ToGetResponse(data *model.UserDataPublic) *generated.GetResponse
	ToPublicUser(data *generated.UpdateRequest) *generated.PublicUser
}
