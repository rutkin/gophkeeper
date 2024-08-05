package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	mock_port "github.com/rutkin/gophkeeper/internal/server/core/service/mock"
	"github.com/stretchr/testify/require"
)

func TestKeeperService_SetTextData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_port.NewMockKeeperRepository(ctrl)
	ks := NewKeeperService(mockRepo)
	var storedDataCtx domain.DataContext
	var storedData []byte
	mockRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext, data []byte) error {
			storedDataCtx = dataCtx
			storedData = data
			return nil
		},
	)
	mockRepo.EXPECT().GetData(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext) ([]byte, error) {
			return storedData, nil
		},
	)

	mockRepo.EXPECT().GetMeta(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, userID domain.UserID, id domain.DataID) (domain.DataContext, error) {
			return storedDataCtx, nil
		},
	)

	ctx := context.Background()

	textData := "text"
	err := ks.SetTextData(ctx, domain.TextData{Data: textData})
	require.NoError(t, err)
	data, err := ks.GetTextData(ctx, domain.DataContext{})
	require.NoError(t, err)
	assert.Equal(t, data, textData)
}

func TestKeeperService_SetCredentialsData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_port.NewMockKeeperRepository(ctrl)
	ks := NewKeeperService(mockRepo)
	var storedDataCtx domain.DataContext
	var storedData []byte
	mockRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext, data []byte) error {
			storedDataCtx = dataCtx
			storedData = data
			return nil
		},
	)
	mockRepo.EXPECT().GetData(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext) ([]byte, error) {
			return storedData, nil
		},
	)

	mockRepo.EXPECT().GetMeta(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, userID domain.UserID, id domain.DataID) (domain.DataContext, error) {
			return storedDataCtx, nil
		},
	)

	ctx := context.Background()
	expectedData := domain.CredentialsData{
		Cred: domain.Credentials{
			Username: "user",
			Password: "password",
		},
	}
	err := ks.SetCredentialsData(ctx, expectedData)
	require.NoError(t, err)
	data, err := ks.GetCredentialsData(ctx, domain.DataContext{})
	require.NoError(t, err)
	assert.Equal(t, data, expectedData)
}

func TestKeeperService_SetBankData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_port.NewMockKeeperRepository(ctrl)
	ks := NewKeeperService(mockRepo)
	var storedDataCtx domain.DataContext
	var storedData []byte
	mockRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext, data []byte) error {
			storedDataCtx = dataCtx
			storedData = data
			return nil
		},
	)
	mockRepo.EXPECT().GetData(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext) ([]byte, error) {
			return storedData, nil
		},
	)

	mockRepo.EXPECT().GetMeta(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, userID domain.UserID, id domain.DataID) (domain.DataContext, error) {
			return storedDataCtx, nil
		},
	)

	ctx := context.Background()
	expectedData := domain.BankData{
		Card: domain.Card{
			CardNumber: "number",
			CardHolder: "holder",
			Cvv:        111,
		},
	}
	err := ks.SetBankData(ctx, expectedData)
	require.NoError(t, err)
	data, err := ks.GetBankData(ctx, domain.DataContext{})
	require.NoError(t, err)
	assert.Equal(t, data, expectedData)
}

func TestKeeperService_SetBinaryData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_port.NewMockKeeperRepository(ctrl)
	ks := NewKeeperService(mockRepo)
	var storedDataCtx domain.DataContext
	var storedData []byte
	mockRepo.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext, data []byte) error {
			storedDataCtx = dataCtx
			storedData = data
			return nil
		},
	)
	mockRepo.EXPECT().GetData(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, dataCtx domain.DataContext) ([]byte, error) {
			return storedData, nil
		},
	)

	mockRepo.EXPECT().GetMeta(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, userID domain.UserID, id domain.DataID) (domain.DataContext, error) {
			return storedDataCtx, nil
		},
	)

	ctx := context.Background()
	expectedData := domain.BinaryData{
		Data: []byte("data"),
	}
	err := ks.SetBinaryData(ctx, expectedData)
	require.NoError(t, err)
	data, err := ks.GetBinaryData(ctx, domain.DataContext{})
	require.NoError(t, err)
	assert.Equal(t, data, expectedData)
}
