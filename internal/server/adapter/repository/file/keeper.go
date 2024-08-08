package repositry

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

type KeeperRepository struct {
	storagePath string
}

func NewKeeper() (*KeeperRepository, error) {
	storagePath := "./keeper_storage"
	err := os.Mkdir(storagePath, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Err(err).Msg("Failed to create repository")
		return nil, err
	}
	return &KeeperRepository{storagePath: storagePath}, nil
}

func (ks *KeeperRepository) GetAllData(ctx context.Context, userID domain.UserID) ([]domain.DataContext, error) {
	userPath := ks.storagePath + "/" + string(userID)
	entries, err := os.ReadDir(userPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		log.Err(err).Msg("Failed to get read dir")
		return nil, err
	}
	var result []domain.DataContext
	for _, entry := range entries {
		if entry.IsDir() {
			meta, err := ks.GetMeta(ctx, userID, domain.DataID(entry.Name()))
			if err != nil {
				log.Err(err).Msgf("failed to get meta '%s'", entry.Name())
			}
			result = append(result, meta)
		}
	}

	return result, nil
}

func (ks *KeeperRepository) Set(ctx context.Context, dataCtx domain.DataContext, data []byte) error {
	dataPath := ks.storagePath + "/" + string(dataCtx.UserID) + "/" + string(dataCtx.ID)
	err := os.MkdirAll(dataPath, os.ModePerm)
	if err != nil && err != os.ErrExist {
		log.Err(err).Msgf("Failed to create directory '%s'", dataPath)
		return err
	}

	var metaBuf bytes.Buffer
	encoder := gob.NewEncoder(&metaBuf)
	err = encoder.Encode(dataCtx)
	if err != nil {
		log.Err(err).Msg("Failed to encode data context")
		return err
	}

	err = os.WriteFile(dataPath+"/meta", metaBuf.Bytes(), os.ModePerm)
	if err != nil {
		log.Err(err).Msgf("Failed to write meta '%s'", dataPath)
		return err
	}

	err = os.WriteFile(dataPath+"/data", data, os.ModePerm)
	if err != nil {
		log.Err(err).Msgf("Failed to write file '%s'", dataPath)
		return err
	}
	return nil
}

func (ks *KeeperRepository) GetData(ctx context.Context, dataCtx domain.DataContext) ([]byte, error) {
	dataPath := ks.storagePath + "/" + string(dataCtx.UserID) + "/" + string(dataCtx.ID)

	data, err := os.ReadFile(dataPath + "/data")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		log.Err(err).Msgf("Failed to get data '%s'", dataPath)
		return nil, err
	}
	return data, err
}

func (ks *KeeperRepository) GetMeta(ctx context.Context, userID domain.UserID, id domain.DataID) (domain.DataContext, error) {
	dataPath := ks.storagePath + "/" + string(userID) + "/" + string(id)

	meta, err := os.Open(dataPath + "/meta")
	if err != nil {
		if os.IsNotExist(err) {
			return domain.DataContext{}, domain.ErrNotFound
		}
		log.Err(err).Msgf("Failed to get data '%s'", dataPath)
		return domain.DataContext{}, err
	}

	var dataCtx domain.DataContext
	encoder := gob.NewDecoder(meta)
	err = encoder.Decode(&dataCtx)
	if err != nil {
		log.Err(err).Msg("Failed to decode data context")
		return domain.DataContext{}, err
	}
	return dataCtx, nil
}

func (ks *KeeperRepository) Delete(ctx context.Context, dataCtx domain.DataContext) error {
	dataPath := ks.storagePath + "/" + string(dataCtx.UserID) + "/" + string(dataCtx.ID)
	err := os.RemoveAll(dataPath)
	if err != nil {
		log.Err(err).Msgf("failed to remove data: %s", dataPath)
		return err
	}
	return nil
}
