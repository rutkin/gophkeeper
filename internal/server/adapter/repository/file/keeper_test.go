package repositry

import (
	"context"
	"os"
	"testing"

	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	"github.com/stretchr/testify/require"
)

func TestKeeperRepository(t *testing.T) {
	err := os.Mkdir("./test_repo", os.ModePerm)
	defer os.RemoveAll("./test_repo")
	require.NoError(t, err)
	repo := KeeperRepository{"./test_repo"}
	ctx := context.Background()
	_, err = repo.GetAllData(ctx, domain.UserID("id"))
	require.Equal(t, domain.ErrNotFound, err)
	dataCtx := domain.DataContext{
		ID:     "id",
		UserID: "user_id",
		Meta:   "meta",
		Title:  "title",
		Type:   domain.BinaryType,
	}
	data := []byte("data")
	err = repo.Set(ctx, dataCtx, data)
	require.NoError(t, err)
	actualData, err := repo.GetData(ctx, dataCtx)
	require.NoError(t, err)
	require.Equal(t, data, actualData)
	actualMeta, err := repo.GetMeta(ctx, domain.UserID("user_id"), domain.DataID("id"))
	require.NoError(t, err)
	require.Equal(t, dataCtx, actualMeta)
	expectedData, err := repo.GetAllData(ctx, domain.UserID("user_id"))
	require.NoError(t, err)
	require.Equal(t, 1, len(expectedData))
	require.Equal(t, dataCtx, expectedData[0])
	err = repo.Delete(ctx, dataCtx)
	require.NoError(t, err)
	_, err = repo.GetData(ctx, dataCtx)
	require.Equal(t, domain.ErrNotFound, err)
}
