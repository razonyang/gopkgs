package models

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Username          string `db:"username" json:"username"`
	Email             string `db:"email" json:"email"`
	EmailVerified     bool   `db:"email_verified" json:"email_verified"`
	VerificationToken string `db:"verification_token" json:"verification_token"`
	HashedPassword    string `db:"hash_password" json:"hash_password"`
}

func NewUser(username, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:          username,
		Email:             email,
		HashedPassword:    string(hashedPassword),
		VerificationToken: strings.ReplaceAll(uuid.New().String(), "-", ""),
	}

	return user, nil
}

func (u *User) Insert(ctx context.Context, db *sqlx.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	query := `
INSERT INTO users(username, email, email_verified, verification_token, hashed_password, created_at, updated_at) VALUES
(?, ?, ?, ?, ?, ?, ?)
`
	res, err := db.ExecContext(ctx, query, u.Username, u.Email, u.EmailVerified, u.VerificationToken, u.HashedPassword, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}
	if u.ID, err = res.LastInsertId(); err != nil {
		return err
	}

	return nil
}
