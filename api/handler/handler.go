package handler

import (
	"Auth-service/storage/postgres"
	"log/slog"
)

type Handler struct {
	Logger *slog.Logger
	User   *postgres.UserRepo
}

func NewHandler(logger *slog.Logger, user *postgres.UserRepo) *Handler {
	return &Handler{Logger: logger, User: user}
}
