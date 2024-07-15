package postgres

import (
	"database/sql"
	"log/slog"
)

type UserRepo struct {
	Logger *slog.Logger
	DB     *sql.DB
}

func NewUserRepo(logger *slog.Logger, db *sql.DB) *UserRepo {
	return &UserRepo{Logger: logger, DB: db}
}
