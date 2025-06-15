package converter

import (
	"log"

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

	log.Printf("Converting PublicUser to UserDataPublic: %+v", src)
	res = &serviceDef.UserDataPublic{}
	res.Id = src.Id
	res.Name.Valid = len(src.Name) != 0
	res.Name.String = src.Name
	res.Email.Valid = len(src.Email) != 0
	res.Email.String = src.Email
	res.Role = int64(src.Role)
	log.Printf("Converted UserDataPublic: %+v", res)
	return

}

func (conv *AuthConverter) ToProtoUserDataPublic(src *serviceDef.UserDataPublic) *generated.PublicUser {

	return &generated.PublicUser{
		Id:    src.Id,
		Name:  src.Name.String,
		Email: src.Email.String,
		Role:  generated.Role(src.Role),
	}

}

func (conv *AuthConverter) ToApiUserDataPrivate(data *generated.CreateRequest) *generated.PrivateUser {
	return &generated.PrivateUser{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Role:     data.Role,
	}
}

func (conv *AuthConverter) ToGetResponse(data *serviceDef.UserDataPublic) *generated.GetResponse {
	return &generated.GetResponse{
		Id:    data.Id,
		Name:  data.Name.String,
		Email: data.Email.String,
		Role:  generated.Role(data.Role),
	}
}

func (conv *AuthConverter) ToPublicUser(data *generated.UpdateRequest) *generated.PublicUser {
	return &generated.PublicUser{
		Id:    data.Id,
		Name:  data.Name.GetValue(),
		Email: data.Email.GetValue(),
	}
}
