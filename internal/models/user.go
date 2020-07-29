package models

type User struct {
	Model
	Email          string `db:"email" json:"email"`
	EmailVerified  bool   `db:"email_verified" json:"email_verified"`
	HashedPassword string `db:"hash_password" json:"hash_password"`
}
