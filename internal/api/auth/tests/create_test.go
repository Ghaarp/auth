package tests

import (
	"context"
	"testing"

	"github.com/Ghaarp/auth/internal/api/auth"
	"github.com/Ghaarp/auth/internal/service"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Ghaarp/auth/internal/service/auth/model"
	serviceMocks "github.com/Ghaarp/auth/internal/service/mocks"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

func TestCreate(t *testing.T) {

	t.Parallel()

	type authServiceMockFunc func(controller *minimock.Controller) service.AuthService

	type args struct {
		ctx     context.Context
		request *generated.CreateRequest
	}

	var (
		ctx            = context.Background()
		mockController = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		pass  = "123"
		role  = generated.Role_R_ADMIN

		request = &generated.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        pass,
			PasswordConfirm: pass,
			Role:            role,
		}

		result = &generated.CreateResponse{
			Id: id,
		}

		model = &model.UserDataPrivate{
			Name:     name,
			Email:    email,
			Password: pass,
			Role:     int64(role),
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *generated.CreateResponse
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "Case 1",
			args: args{
				ctx:     ctx,
				request: request,
			},
			want: result,
			err:  nil,
			authServiceMock: func(controller *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(t)
				mock.CreateMock.Expect(ctx, model).Return(id, nil)
				return mock
			},
		},
	}

	for _, test := range tests {

		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := test.authServiceMock(mockController)
			api := auth.NewAuthImplementation(authServiceMock)

			createResult, err := api.Create(test.args.ctx, test.args.request)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, createResult)
		})
	}
}
