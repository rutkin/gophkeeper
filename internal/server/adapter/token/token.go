package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

var (
	secretKey   = []byte("secret-key")
	userIDKey   = "userid"
	userNameKey = "username"
)

type TokenService struct {
	Exp time.Duration
}

func New(exp time.Duration) *TokenService {
	return &TokenService{Exp: exp}
}

func (ts *TokenService) CreateToken(user domain.User) (domain.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			userIDKey:   user.ID,
			userNameKey: user.Name,
			"exp":       time.Now().Add(ts.Exp).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return domain.Token(tokenString), nil
}

func (ts *TokenService) VerifyToken(tokenStr string) (domain.TokenPayload, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return domain.TokenPayload{}, err
	}

	if !token.Valid {
		log.Error().Msg("failed to validate token")
		return domain.TokenPayload{}, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.TokenPayload{}, domain.ErrInvalidToken
	}

	userID, ok := claims[userIDKey].(string)
	if !ok {
		return domain.TokenPayload{}, domain.ErrInvalidToken
	}

	userName, ok := claims[userNameKey].(string)
	if !ok {
		return domain.TokenPayload{}, domain.ErrInvalidToken
	}

	return domain.TokenPayload{
		ID:   domain.UserID(userID),
		Name: userName,
	}, nil
}
