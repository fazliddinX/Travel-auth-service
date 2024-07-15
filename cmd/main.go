package main

import (
	"Auth-service/api"
	"Auth-service/api/handler"
	"Auth-service/genproto/auth_service"
	logger2 "Auth-service/logger"
	"Auth-service/server"
	"Auth-service/storage/postgres"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger := logger2.InitLogger()

	db, err := postgres.Connection()
	if err != nil {
		logger.Error("Failed to connect to database", "error", err.Error())
		log.Fatalln(err)
	}

	user := postgres.NewUserRepo(logger, db)

	server1 := server.NewServer(user, logger)

	handler1 := handler.NewHandler(logger, server1)

	router := api.Router(handler1)

	go log.Fatalln(router.Run(":7070"))

	//--------------------------------------------------------------------------------------

	listner, err := net.Listen("tcp", ":5051")
	if err != nil {
		logger.Error("Failed to listen", "error", err.Error())
		log.Fatalln(err)
	}
	defer listner.Close()

	grpcServer := grpc.NewServer()
	auth_service.RegisterAuthServiceServer(grpcServer, server1)
	log.Fatalln(grpcServer.Serve(listner))
}
