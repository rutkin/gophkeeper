package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
	"github.com/rutkin/gophkeeper/internal/server/core/port"
	"github.com/rutkin/gophkeeper/internal/server/core/util"
)

type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

func NewAuthService(repo port.UserRepository, ts port.TokenService) *AuthService {
	return &AuthService{
		repo: repo,
		ts:   ts,
	}
}

func (a *AuthService) Register(ctx context.Context, user domain.User) error {
	hashPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.ID = domain.UserID(uuid.NewString())
	user.Password = hashPassword
	return a.repo.CreateUser(ctx, user)
}

func (a *AuthService) Login(ctx context.Context, user domain.User) (domain.Token, error) {
	curUser, err := a.repo.GetUserByName(ctx, user.Name)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.Token(""), domain.ErrInvalidCredentials
		}
		return domain.Token(""), err
	}

	err = util.ComparePassword(user.Password, curUser.Password)
	if err != nil {
		return domain.Token(""), domain.ErrInvalidCredentials
	}

	token, err := a.ts.CreateToken(curUser)
	if err != nil {
		return domain.Token(""), err
	}
	return token, nil
}
