package models

import (
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"pkg.razonyang.com/gopkgs/internal/helper"
)

func init() {
	gob.Register(&User{})
}

type User struct {
	Model
	Username           string         `db:"username" json:"username"`
	Email              string         `db:"email" json:"email"`
	EmailVerified      bool           `db:"email_verified" json:"email_verified"`
	VerificationToken  sql.NullString `db:"verification_token" json:"verification_token"`
	HashedPassword     string         `db:"hashed_password" json:"hashed_password"`
	PasswordResetToken sql.NullString `db:"password_reset_token" json:"password_reset_token"`
}

func NewUser(username, email, password string) (*User, error) {
	hashedPassword, err := generatePassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:          username,
		Email:             email,
		HashedPassword:    hashedPassword,
		VerificationToken: GenerateVerificationToken(),
	}

	return user, nil
}

func FindUserByUsername(ctx context.Context, db *sqlx.DB, dest interface{}, username string) error {
	return db.GetContext(ctx, dest, "SELECT * FROM users WHERE username = ?", username)
}

func FindUserByEmail(ctx context.Context, db *sqlx.DB, dest interface{}, email string) error {
	return db.GetContext(ctx, dest, "SELECT * FROM users WHERE email = ?", email)
}

func FindUserByVerificationToken(ctx context.Context, db *sqlx.DB, dest interface{}, token string) error {
	return db.GetContext(ctx, dest, "SELECT * FROM users WHERE verification_token = ?", token)
}

func FindUserByPasswordResetToken(ctx context.Context, db *sqlx.DB, dest interface{}, token string) error {
	return db.GetContext(ctx, dest, "SELECT * FROM users WHERE password_reset_token = ?", token)
}

func (u *User) GetID() string {
	return strconv.FormatInt(u.ID, 10)
}

func (u *User) SetPassword(ctx context.Context, db *sqlx.DB, password string) (err error) {
	u.HashedPassword, err = generatePassword(password)
	if err != nil {
		return
	}
	_, err = db.ExecContext(ctx, "UPDATE users SET hashed_password = ?, password_reset_token = null WHERE id = ?", u.HashedPassword, u.ID)
	return
}

func (u *User) Insert(ctx context.Context, db *sqlx.DB) error {
	u.CreatedAt = helper.CurrentUTC()
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
	return isUniqueTokenExpired(u.VerificationToken)
}

func (u *User) IsPasswordResetTokenExpired() bool {
	return isUniqueTokenExpired(u.PasswordResetToken)
}

func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
}

func GenerateVerificationToken() sql.NullString {
	return sql.NullString{
		String: generateUniqueToken(),
		Valid:  true,
	}
}

func GeneratePasswordResetToken() sql.NullString {
	return sql.NullString{
		String: generateUniqueToken(),
		Valid:  true,
	}
}

func generateUniqueToken() string {
	return fmt.Sprintf("%s-%d", strings.ReplaceAll(uuid.New().String(), "-", ""), time.Now().Unix())
}

func isUniqueTokenExpired(token sql.NullString) bool {
	if !token.Valid {
		return true
	}
	parts := strings.Split(token.String, "-")
	if len(parts) != 2 {
		return true
	}

	expiredAt, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return true
	}

	return helper.CurrentUTC().Add(-10*time.Minute).Unix() > expiredAt
}

func generatePassword(password string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}
