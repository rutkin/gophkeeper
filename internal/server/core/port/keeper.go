package port

import (
	"context"

	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

type Keeper interface {
	ListAll(ctx context.Context, id domain.UserID) ([]domain.DataContext, error)
	SetTextData(ctx context.Context, data domain.TextData) error
	GetTextData(ctx context.Context, dataCtx domain.DataContext) (string, error)
	SetBinaryData(ctx context.Context, data domain.BinaryData) error
	GetBinaryData(ctx context.Context, dataCtx domain.DataContext) (domain.BinaryData, error)
	SetCredentialsData(ctx context.Context, data domain.CredentialsData) error
	GetCredentialsData(ctx context.Context, dataCtx domain.DataContext) (domain.CredentialsData, error)
	SetBankData(ctx context.Context, data domain.BankData) error
	GetBankData(ctx context.Context, dataCtx domain.DataContext) (domain.BankData, error)
	Delete(ctx context.Context, dataCtx domain.DataContext) error
}

type KeeperRepository interface {
	GetAllData(ctx context.Context, userID domain.UserID) ([]domain.DataContext, error)
	Set(ctx context.Context, dataCtx domain.DataContext, data []byte) error
	GetData(ctx context.Context, dataCtx domain.DataContext) ([]byte, error)
	GetMeta(ctx context.Context, userID domain.UserID, id domain.DataID) (domain.DataContext, error)
	Delete(ctx context.Context, dataCtx domain.DataContext) error
}
