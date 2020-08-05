package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Username          string         `db:"username" json:"username"`
	Email             string         `db:"email" json:"email"`
	EmailVerified     bool           `db:"email_verified" json:"email_verified"`
	VerificationToken sql.NullString `db:"verification_token" json:"verification_token"`
	HashedPassword    string         `db:"hashed_password" json:"hashed_password"`
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
		VerificationToken: GenerateVerificationToken(),
	}

	return user, nil
}

func FindUserByEmail(ctx context.Context, db *sqlx.DB, dest interface{}, email string) error {
	return db.GetContext(ctx, dest, "SELECT * FROM users WHERE email = ?", email)
}

func FindUserByVerificationToken(ctx context.Context, db *sqlx.DB, dest interface{}, token string) error {
	return db.GetContext(ctx, dest, "SELECT * FROM users WHERE verification_token = ?", token)
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

func (u *User) IsVerificationExpired() bool {
	if !u.VerificationToken.Valid {
		return false
	}
	parts := strings.Split(u.VerificationToken.String, "-")
	if len(parts) != 2 {
		return true
	}

	expiredAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return true
	}

	return time.Now().Add(-10*time.Minute).Unix() > expiredAt
}

func GenerateVerificationToken() sql.NullString {
	return sql.NullString{
		String: fmt.Sprintf("%s-%d", strings.ReplaceAll(uuid.New().String(), "-", ""), time.Now().Unix()),
		Valid:  true,
	}
}
