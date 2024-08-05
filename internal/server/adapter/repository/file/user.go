package repositry

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

var userStorageFile = "./user_repo"

type UserRepository struct {
	Users    map[domain.UserName]domain.User
	FilePath string
}

func NewUser(filepath string) (*UserRepository, error) {
	if len(filepath) == 0 {
		filepath = userStorageFile
	}
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Err(err).Msg("failed to open user repository file")
		return &UserRepository{}, err
	}

	userRepo := UserRepository{FilePath: filepath}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&userRepo.Users)
	if err != nil {
		if err == io.EOF {
			return &UserRepository{make(map[domain.UserName]domain.User), filepath}, nil
		}
		log.Err(err).Msg("failed to decode user repository")
		return &UserRepository{}, err
	}
	return &userRepo, nil
}

func (us *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	if _, ok := us.Users[user.Name]; ok {
		return domain.ErrUserExists
	}
	us.Users[user.Name] = user
	return nil
}

func (us *UserRepository) GetUserByName(ctx context.Context, name domain.UserName) (domain.User, error) {
	user, ok := us.Users[name]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (us *UserRepository) Close() {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(us.Users)
	if err != nil {
		log.Err(err).Msg("failed to encode user repository")
		return
	}

	err = os.WriteFile(us.FilePath, buf.Bytes(), os.ModePerm)
	if err != nil {
		log.Err(err).Msg("failed to write user repository")
		return
	}
}
