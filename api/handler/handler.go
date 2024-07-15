package handler

import (
	s"Auth-service/server"
	"Auth-service/storage/postgres"
	"log/slog"
)
b
type Handler struct {
	Logger *slog.Logger
	Server   *s.Server
}

func NewHandler(logger *slog.Logger, server *s.Server) *Handler {
	return &Handler{Logger: logger, Server: server}
}
