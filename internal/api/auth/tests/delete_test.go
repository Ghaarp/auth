package tests

import (
	"context"
	"testing"

	"github.com/Ghaarp/auth/internal/api/auth"
	"github.com/Ghaarp/auth/internal/service"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	serviceMocks "github.com/Ghaarp/auth/internal/service/mocks"
	generated "github.com/Ghaarp/auth/pkg/auth_v1"
)

func TestDelete(t *testing.T) {

	t.Parallel()

	type authServiceMockFunc func(controller *minimock.Controller) service.AuthService

	type args struct {
		ctx     context.Context
		request *generated.DeleteRequest
	}

	var (
		ctx            = context.Background()
		mockController = minimock.NewController(t)

		id = gofakeit.Int64()

		request = &generated.DeleteRequest{
			Id: id,
		}

		result = &generated.DeleteResponse{}
	)

	tests := []struct {
		name            string
		args            args
		want            *generated.DeleteResponse
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
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			authServiceMock := test.authServiceMock(mockController)
			api := auth.NewAuthImplementation(authServiceMock)

			deleteResult, err := api.Delete(test.args.ctx, test.args.request)
			require.Equal(t, test.err, err)
			require.Equal(t, test.want, deleteResult)
		})
	}
}
