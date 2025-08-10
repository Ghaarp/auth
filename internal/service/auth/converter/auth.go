package converter

import (
	serviceDef "github.com/Ghaarp/auth/internal/service/auth/model"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

type AuthConverter struct {
}

func CreateConverter() *AuthConverter {
	return &AuthConverter{}
}

func (conv *AuthConverter) ToServiceUserDataPrivate(src *generated.PrivateUser) *serviceDef.UserDataPrivate {

	return &serviceDef.UserDataPrivate{
		Name:     src.Name,
		Email:    src.Email,
		Password: src.Password,
		Role:     int64(src.Role),
	}

}

func (conv *AuthConverter) ToServiceUserDataPublic(src *generated.PublicUser) (res *serviceDef.UserDataPublic) {

	return &serviceDef.UserDataPublic{
		Id:    src.Id,
		Name:  src.Name,
		Email: src.Email,
		Role:  int64(src.Role),
	}
}

func (conv *AuthConverter) ToProtoUserDataPublic(src *serviceDef.UserDataPublic) *generated.PublicUser {

	return &generated.PublicUser{
		Id:    src.Id,
		Name:  src.Name,
		Email: src.Email,
		Role:  generated.Role(src.Role),
	}

}

func (conv *AuthConverter) ToProtoUserDataPrivate(data *generated.CreateRequest) *generated.PrivateUser {
	return &generated.PrivateUser{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Role:     data.Role,
	}
}

func (conv *AuthConverter) ToProtoGetResponse(data *serviceDef.UserDataPublic) *generated.GetResponse {
	return &generated.GetResponse{
		Id:    data.Id,
		Name:  data.Name,
		Email: data.Email,
		Role:  generated.Role(data.Role),
	}
}

func (conv *AuthConverter) ToProtoPublicUser(data *generated.UpdateRequest) *generated.PublicUser {
	return &generated.PublicUser{
		Id:    data.Id,
		Name:  data.Name.GetValue(),
		Email: data.Email.GetValue(),
	}
}
