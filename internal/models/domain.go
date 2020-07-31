package models

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"pkg.razonyang.com/gopkgs/internal/core"
)

type Domain struct {
	Model
	UserID       string `json:"user_id" db:"user_id"`
	Name         string `json:"name" db:"name"`
	Verified     bool   `json:"verified" db:"verified"`
	ChallengeTXT string `json:"-" db:"challenge_txt"`
}

func NewDomain(name, userID string) *Domain {
	domain := &Domain{
		Name:         name,
		UserID:       userID,
		ChallengeTXT: strings.ReplaceAll(uuid.New().String(), "-", ""),
	}

	now := time.Now()
	domain.CreatedAt = now
	domain.UpdatedAt = now
	return domain
}

func (d Domain) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required, is.Domain),
	)
}

func (d *Domain) Challenge(ctx context.Context, db *sqlx.DB) error {
	if err := d.verifyChallengeTXT(ctx); err != nil {
		return err
	}

	query := "UPDATE domains SET verified = ? WHERE id = ?"
	d.Verified = true
	if _, err := db.ExecContext(ctx, query, d.Verified, d.ID); err != nil {
		return err
	}

	return nil
}

func (d *Domain) verifyChallengeTXT(ctx context.Context) error {
	multiErr := &core.MultiError{}
	for _, schema := range []string{"https", "http"} {
		u := &url.URL{
			Scheme: schema,
			Host:   d.Name,
			Path:   "/.well-known/gopkgs-challenge/" + d.ChallengeTXT,
		}
		timeOutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(timeOutCtx, http.MethodGet, u.String(), nil)
		if err != nil {
			multiErr.Add(err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			multiErr.Add(err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			multiErr.Add(fmt.Errorf("invalid response status code %d", resp.StatusCode))
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			multiErr.Add(err)
			continue
		}
		if string(body) != d.ChallengeTXT {
			multiErr.Add(fmt.Errorf("invalid challenge txt"))
			continue
		}

		return nil
	}
	return multiErr
}

func (d *Domain) Insert(ctx context.Context, db *sqlx.DB) error {
	query := "INSERT INTO domains(user_id, name, verified, challenge_txt, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)"
	res, err := db.ExecContext(ctx, query, d.UserID, d.Name, d.Verified, d.ChallengeTXT, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return err
	}
	d.ID, err = res.LastInsertId()
	return err
}

func (d *Domain) Update(ctx context.Context, db *sqlx.DB) error {
	d.UpdatedAt = time.Now()
	query := "UPDATE domains SET name = ?, verified = ?, challenge_txt = ?, updated_at = ? WHERE id = ?"
	_, err := db.ExecContext(ctx, query, d.Name, d.Verified, d.ChallengeTXT, d.UpdatedAt, d.ID)
	if err != nil {
		return err
	}
	return err
}

func (d *Domain) Delete(ctx context.Context, db *sqlx.DB) error {
	query := "DELETE FROM domains WHERE id = ?"
	_, err := db.ExecContext(ctx, query, d.ID)
	return err
}

func CountDomains(ctx context.Context, db *sqlx.DB, count *int64) error {
	query := "SELECT COUNT(id) FROM domains"
	return db.GetContext(ctx, count, query)
}

func CountDomainsByUser(ctx context.Context, db *sqlx.DB, count *int64, userID string) error {
	query := "SELECT COUNT(id) FROM domains WHERE user_id = ?"
	return db.GetContext(ctx, count, query, userID)
}

func FindDomainsByUser(ctx context.Context, db *sqlx.DB, dest interface{}, userID string) error {
	query := "SELECT * FROM domains WHERE user_id = ? ORDER BY domains.name ASC"
	return db.SelectContext(ctx, dest, query, userID)
}

func FindDomainByUser(ctx context.Context, db *sqlx.DB, dest interface{}, id int64, userID string) error {
	query := "SELECT * FROM domains WHERE id = ? AND user_id = ?"
	return db.GetContext(ctx, dest, query, id, userID)
}

func FindDomainByChallengeTXT(ctx context.Context, db *sqlx.DB, dest interface{}, name, txt string) error {
	query := "SELECT * FROM domains WHERE name = ? AND challenge_txt = ?"
	return db.GetContext(ctx, dest, query, name, txt)
}
