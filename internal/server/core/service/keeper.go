package service

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/gob"

	"github.com/rs/zerolog/log"

	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	"github.com/rutkin/gophkeeper/internal/server/core/port"
)

var password = []byte("secret-key")

type KeeperService struct {
	repo port.KeeperRepository
}

func NewKeeperService(repo port.KeeperRepository) *KeeperService {
	return &KeeperService{repo: repo}
}

func (ks *KeeperService) ListAll(ctx context.Context, id domain.UserID) ([]domain.DataContext, error) {
	return ks.repo.GetAllData(ctx, id)
}

func (ks *KeeperService) SetTextData(ctx context.Context, data domain.TextData) error {
	encryptedData, err := encrypt([]byte(data.Data))
	if err != nil {
		log.Err(err).Msg("failed to encrypt text data")
		return err
	}

	return ks.repo.Set(ctx, data.Ctx, encryptedData)
}

func (ks *KeeperService) GetTextData(ctx context.Context, dataCtx domain.DataContext) (string, error) {
	dataCtx, err := ks.repo.GetMeta(ctx, dataCtx.UserID, dataCtx.ID)
	if err != nil {
		log.Err(err).Msg("failed to get text data from repository")
		return "", err
	}

	data, err := ks.repo.GetData(ctx, dataCtx)
	if err != nil {
		log.Err(err).Msg("failed to get text data from repository")
		return "", err
	}

	decryptedData, err := decrypt(data)
	if err != nil {
		log.Err(err).Msg("failed to decrypt text data")
		return "", err
	}
	return string(decryptedData), nil
}

func (ks *KeeperService) SetBinaryData(ctx context.Context, data domain.BinaryData) error {
	encryptedData, err := encrypt([]byte(data.Data))
	if err != nil {
		log.Err(err).Msg("failed to encrypt text data")
		return err
	}

	return ks.repo.Set(ctx, data.Ctx, encryptedData)
}

func (ks *KeeperService) GetBinaryData(ctx context.Context, dataCtx domain.DataContext) (domain.BinaryData, error) {
	dataCtx, err := ks.repo.GetMeta(ctx, dataCtx.UserID, dataCtx.ID)
	if err != nil {
		log.Err(err).Msg("failed to get text data from repository")
		return domain.BinaryData{}, err
	}
	data, err := ks.repo.GetData(ctx, dataCtx)
	if err != nil {
		log.Err(err).Msg("failed to get binary data from repository")
		return domain.BinaryData{}, err
	}

	decryptedData, err := decrypt(data)
	if err != nil {
		log.Err(err).Msg("failed to decrypt text data")
		return domain.BinaryData{}, err
	}
	return domain.BinaryData{Ctx: dataCtx, Data: decryptedData}, nil
}

func (ks *KeeperService) SetCredentialsData(ctx context.Context, data domain.CredentialsData) error {
	var dataBuf bytes.Buffer
	encoder := gob.NewEncoder(&dataBuf)
	err := encoder.Encode(&data.Cred)
	if err != nil {
		log.Err(err).Msg("failed to encode credentials")
		return err
	}
	encrypted, err := encrypt(dataBuf.Bytes())
	if err != nil {
		log.Err(err).Msg("failed to encrypt credentials")
		return err
	}

	err = ks.repo.Set(ctx, data.Ctx, encrypted)
	if err != nil {
		log.Err(err).Msg("failed to set credentials in repository")
		return err
	}
	return nil
}

func (ks *KeeperService) GetCredentialsData(ctx context.Context, dataCtx domain.DataContext) (domain.CredentialsData, error) {
	dataCtx, err := ks.repo.GetMeta(ctx, dataCtx.UserID, dataCtx.ID)
	if err != nil {
		log.Err(err).Msg("failed to get credentials meta from repository")
		return domain.CredentialsData{}, err
	}
	data, err := ks.repo.GetData(ctx, dataCtx)
	if err != nil {
		log.Err(err).Msg("failed to credentials from repository")
		return domain.CredentialsData{}, err
	}

	decryptedData, err := decrypt(data)
	if err != nil {
		log.Err(err).Msg("failed to decrypt text data")
		return domain.CredentialsData{}, err
	}

	decoder := gob.NewDecoder(bytes.NewReader(decryptedData))
	var cred domain.Credentials
	err = decoder.Decode(&cred)
	if err != nil {
		log.Err(err).Msg("failed to decode credentials")
		return domain.CredentialsData{}, err
	}
	return domain.CredentialsData{Ctx: dataCtx, Cred: cred}, nil
}

func (ks *KeeperService) SetBankData(ctx context.Context, data domain.BankData) error {
	return SetDataImpl(ctx, ks.repo, data.Ctx, data.Card)
}

func (ks *KeeperService) GetBankData(ctx context.Context, dataCtx domain.DataContext) (domain.BankData, error) {
	dataCtx, err := ks.repo.GetMeta(ctx, dataCtx.UserID, dataCtx.ID)
	if err != nil {
		log.Err(err).Msg("failed to get credentials meta from repository")
		return domain.BankData{}, err
	}
	data, err := ks.repo.GetData(ctx, dataCtx)
	if err != nil {
		log.Err(err).Msg("failed to credentials from repository")
		return domain.BankData{}, err
	}

	decryptedData, err := decrypt(data)
	if err != nil {
		log.Err(err).Msg("failed to decrypt text data")
		return domain.BankData{}, err
	}

	decoder := gob.NewDecoder(bytes.NewReader(decryptedData))
	var card domain.Card
	err = decoder.Decode(&card)
	if err != nil {
		log.Err(err).Msg("failed to decode credentials")
		return domain.BankData{}, err
	}
	return domain.BankData{Ctx: dataCtx, Card: card}, nil
}

func (ks *KeeperService) Delete(ctx context.Context, dataCtx domain.DataContext) error {
	return ks.repo.Delete(ctx, dataCtx)
}

func SetDataImpl[TData any](ctx context.Context, repo port.KeeperRepository, dataCtx domain.DataContext, data TData) error {
	var dataBuf bytes.Buffer
	encoder := gob.NewEncoder(&dataBuf)
	err := encoder.Encode(&data)
	if err != nil {
		log.Err(err).Msg("failed to encode data")
		return err
	}
	encrypted, err := encrypt(dataBuf.Bytes())
	if err != nil {
		log.Err(err).Msg("failed to encrypt data")
		return err
	}

	err = repo.Set(ctx, dataCtx, encrypted)
	if err != nil {
		log.Err(err).Msg("failed to set data in repository")
		return err
	}
	return nil
}

func encrypt(src []byte) ([]byte, error) {
	key := sha256.Sum256([]byte(password))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	dst := aesgcm.Seal(nil, nonce, src, nil)
	return dst, err
}

func decrypt(src []byte) ([]byte, error) {
	key := sha256.Sum256([]byte(password))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	dst, err := aesgcm.Open(nil, nonce, src, nil)
	if err != nil {
		return nil, err
	}
	return dst, err
}
