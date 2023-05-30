package users

import (
	"context"
	"database/sql"
)

const (
	postgresInsertUser = `
	INSERT INTO users(email, xpriv, paymail, created_at)
	VALUES($1, $2, $3, $4)
	`

	postgresGetUserByEmail = `
	SELECT email, xpriv, paymail, created_at
	FROM users
	WHERE email = $1
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
	if _, err = stmt.Exec(user.Email, user.Xpriv, user.Paymail, user.CreatedAt); err != nil {
		return err
	}
	return tx.Commit()
}

// GetUserByEmail returns user by email.
func (r *UsersRepository) GetUserByEmail(ctx context.Context, email string) (*UserDto, error) {
	var user UserDto
	row := r.db.QueryRowContext(ctx, postgresGetUserByEmail, email)
	if err := row.Scan(&user.Email, &user.Xpriv, &user.Paymail, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
