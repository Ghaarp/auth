package tests

import (
	"context"
	"testing"

	"github.com/Ghaarp/auth/internal/api/auth"
	"github.com/Ghaarp/auth/internal/service"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Ghaarp/auth/internal/service/auth/model"
	serviceMocks "github.com/Ghaarp/auth/internal/service/mocks"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

func TestUpdate(t *testing.T) {

	t.Parallel()

	type authServiceMockFunc func(controller *minimock.Controller) service.AuthService

	type args struct {
		ctx     context.Context
		request *generated.UpdateRequest
	}

	var (
		ctx            = context.Background()
		mockController = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = generated.Role_R_ADMIN

		request = &generated.UpdateRequest{
			Id:    id,
			Name:  wrapperspb.String(name),
			Email: wrapperspb.String(email),
		}

		result = &generated.UpdateResponse{}

		model = &model.UserDataPublic{
			Id:    id,
			Name:  name,
			Email: email,
			Role:  int64(role),
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *generated.UpdateResponse
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
				mock.UpdateMock.Expect(ctx, model).Return(nil)
				return mock
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := test.authServiceMock(mockController)
			api := auth.NewAuthImplementation(authServiceMock)

			updateResult, err := api.Update(test.args.ctx, test.args.request)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, updateResult)
		})
	}
}
