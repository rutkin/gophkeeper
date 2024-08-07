package postgress

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/rutkin/gophkeeper/internal/server/core/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(databaseDSN string) (*UserRepository, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		log.Err(err).Msg("failed connect to postgres")
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id VARCHAR (50) UNIQUE NOT NULL, name VARCHAR (50) NOT NULL, password VARCHAR(100) NOT NULL)")
	if err != nil {
		log.Err(err).Msg("failed to create users table")
		return nil, err
	}

	return &UserRepository{db: db}, nil
}

func (us *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	_, err := us.db.Exec("INSERT INTO users (id, name, password) Values ($1, $2, $3)", user.ID, user.Name, user.Password)
	if err != nil {
		log.Err(err).Msg("failed to create user")
		return err
	}
	return nil
}

func (us *UserRepository) GetUserByName(ctx context.Context, name domain.UserName) (domain.User, error) {
	row := us.db.QueryRow("SELECT id, password FROM users WHERE name=$1", name)
	var id string
	var password string
	err := row.Scan(&id, &password)
	if err != nil {
		log.Err(err).Msg("failed to get user")
		return domain.User{}, err
	}
	return domain.User{
		ID:       domain.UserID(id),
		Name:     name,
		Password: password,
	}, nil
}

func (us *UserRepository) Close() {
	us.db.Close()
}
