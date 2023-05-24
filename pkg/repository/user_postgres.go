package repository

import (
	"bux-wallet/pkg/domain"
	"bux-wallet/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type UserPGX struct {
	pool   *pgxpool.Pool
	logger zerolog.Logger
}

func (b UserPGX) Create() (domain.User, error) {
	return domain.User{}, nil
}

func (b UserPGX) GetByCredentials(email, password string) (domain.User, error) {
	return domain.User{}, nil
}

func NewUserPGXRepo(pool *pgxpool.Pool) *UserPGX {
	return &UserPGX{
		pool:   pool,
		logger: logger.NewRepositoryLogger(),
	}
}
