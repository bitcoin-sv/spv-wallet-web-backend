package users

import (
	"context"
	"database/sql"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
)

const (
	postgresInsertUser = `
	INSERT INTO users(email, xpriv, paymail, created_at)
	VALUES($1, $2, $3, $4)
	`

	postgresGetUserByEmail = `
	SELECT id, email, xpriv, paymail, created_at
	FROM users
	WHERE email = $1
	`

	postgresGetUserById = `
	SELECT id, email, xpriv, paymail, created_at
	FROM users
	WHERE id = $1
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
func (r *UsersRepository) InsertUser(ctx context.Context, user *users.User) error {
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
	if _, err = stmt.Exec(user.Email, user.Xpriv, user.Paymail, user.CreatedAt); err != nil {
		return err
	}
	return tx.Commit()
}

// GetUserByEmail returns user by email.
func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	var user UserDto
	row := r.db.QueryRowContext(ctx, postgresGetUserByEmail, email)
	if err := row.Scan(&user.Id, &user.Email, &user.Xpriv, &user.Paymail, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user.toUser(), nil
}

// GetUserById returns user by id.
func (r *UsersRepository) GetUserById(ctx context.Context, id int) (*users.User, error) {
	var user UserDto
	row := r.db.QueryRowContext(ctx, postgresGetUserById, id)
	if err := row.Scan(&user.Id, &user.Email, &user.Xpriv, &user.Paymail, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user.toUser(), nil
}
