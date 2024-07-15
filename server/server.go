package server

import (
	"Auth-service/genproto/auth_service"
	"Auth-service/storage/postgres"
	"log/slog"
)

type Server struct {
	auth_service.UnimplementedAuthServiceServer
	Logger *slog.Logger
	User   *postgres.UserRepo
}

func NewServer(user *postgres.UserRepo, logger *slog.Logger) *Server {
	return &Server{User: user, Logger: logger}
}
