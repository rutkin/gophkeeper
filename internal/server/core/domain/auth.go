package domain

type UserID string
type UserName string

type User struct {
	ID       UserID
	Name     UserName
	Password string
}

type Token string

type TokenPayload struct {
	ID   UserID
	Name string
}
