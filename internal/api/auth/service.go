package auth

import (
	"github.com/Ghaarp/auth/internal/service"
	converter "github.com/Ghaarp/auth/internal/service/auth/converter"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

type AuthImplementation struct {
	generated.UnimplementedAuthV1Server
	authService      service.AuthService
	serviceConverter service.ServiceConverter
}

func NewAuthImplementation(authService service.AuthService) *AuthImplementation {
	return &AuthImplementation{
		authService:      authService,
		serviceConverter: converter.CreateConverter(),
	}
}
