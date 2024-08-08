package repositry

import (
	"context"
	"testing"

	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	"github.com/stretchr/testify/require"
)

func TestUserRepository(t *testing.T) {
	filePath := "./test_user_repo"
	//defer os.Remove(filePath)
	userRepo, err := NewUser(filePath)
	require.NoError(t, err)
	ctx := context.Background()
	user := domain.User{
		ID:       "id",
		Name:     "name",
		Password: "password",
	}
	err = userRepo.CreateUser(ctx, user)
	require.NoError(t, err)
	actualUser, err := userRepo.GetUserByName(ctx, "name")
	require.NoError(t, err)
	require.Equal(t, user, actualUser)
}
