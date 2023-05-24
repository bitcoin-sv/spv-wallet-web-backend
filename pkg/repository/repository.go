package repository

import (
	"bux-wallet/pkg/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	Create() (domain.User, error)
	GetByCredentials(email, password string) (domain.User, error)
}

type Repositories struct {
	User User
}

func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User: NewUserPGXRepo(pool),
	}
}
