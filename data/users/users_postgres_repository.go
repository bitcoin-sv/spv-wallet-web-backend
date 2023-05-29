package users

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

const (
	postgresInsertUser = `
	INSERT INTO users(email, xpriv, created_at)
	VALUES($1, $2, $3)
	`
)

// UsersRepository is a repository for users.
type UsersRepository struct {
	db *sql.DB
}

// NewUsersRepository creates a new users repository.
func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

// InsertUser inserts a user to db.
func (r *UsersRepository) InsertUser(ctx context.Context, user *UserDto) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	stmt, err := r.db.Prepare(postgresInsertUser)
	if err != nil {
		return err
	}
	defer stmt.Close() //nolint:all
	if _, err = stmt.Exec(user.Email, user.Xpriv, user.CreatedAt); err != nil {
		return errors.Wrap(err, "failed to insert new user")
	}
	return errors.Wrap(tx.Commit(), "failed to commit tx")
}
