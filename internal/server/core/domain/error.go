package domain

import "errors"

var (
	ErrNotFound                   = errors.New("not found")
	ErrInvalidCredentials         = errors.New("invalid login/password")
	ErrInvalidToken               = errors.New("invalid token")
	ErrUserExists                 = errors.New("user exists")
	ErrEmptyAuthorizationHeader   = errors.New("empty authorizaation header")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header")
	ErrInvalidAuthorizationType   = errors.New("invalid authorization type")
	ErrBadRequest                 = errors.New("bad request")
)
