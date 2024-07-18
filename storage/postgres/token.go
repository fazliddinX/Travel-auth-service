package postgres

import (
	"time"
)

func (repo UserRepo) SaveRefreshToken(username, token string, expiresAt time.Time) error {
	_, err := repo.DB.Exec("DELETE FROM refresh_tokens WHERE username = $1", username)
	if err != nil {
		return err
	}

	_, err = repo.DB.Exec("INSERT INTO refresh_tokens (username, token, expires_at) VALUES ($1, $2, $3)",
		username, token, expiresAt)
	return err
}
