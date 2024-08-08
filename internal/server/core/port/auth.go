package port

import (
	"context"

	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

type TokenService interface {
	CreateToken(user domain.User) (domain.Token, error)
	VerifyToken(token string) (domain.TokenPayload, error)
}

type AuthService interface {
	Register(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (domain.Token, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByName(ctx context.Context, name domain.UserName) (domain.User, error)
	Close()
}
