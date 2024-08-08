package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	mock_port "github.com/rutkin/gophkeeper/internal/server/core/service/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_RegisterAndLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_port.NewMockUserRepository(ctrl)
	var storedUser domain.User
	mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, user domain.User) error {
			storedUser = user
			return nil
		},
	)
	mockRepo.EXPECT().GetUserByName(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, name domain.UserName) (domain.User, error) {
			return storedUser, nil
		},
	)
	mockToken := mock_port.NewMockTokenService(ctrl)
	expectedToken := domain.Token("token")
	mockToken.EXPECT().CreateToken(gomock.Any()).DoAndReturn(
		func(user domain.User) (domain.Token, error) {
			return expectedToken, nil
		},
	)
	as := NewAuthService(mockRepo, mockToken)
	user := domain.User{
		ID:       "id",
		Name:     "name",
		Password: "password",
	}
	ctx := context.Background()
	err := as.Register(ctx, user)
	require.NoError(t, err)
	token, err := as.Login(ctx, user)
	require.NoError(t, err)
	require.Equal(t, token, expectedToken)
}
