package main

import (
	"fmt"
	"net"

	"gitlab.com/micro/user_service/config"
	u "gitlab.com/micro/user_service/genproto/user"
	"gitlab.com/micro/user_service/pkg/db"
	"gitlab.com/micro/user_service/pkg/logger"
	"gitlab.com/micro/user_service/service"
	grpcclient "gitlab.com/micro/user_service/service/grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "golang")
	defer logger.Cleanup(log)

	connDb, err := db.ConnectToDB(cfg)
	if err != nil {
		fmt.Println("failed connect database", err)
	}

	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		fmt.Println("failed while grpc client", err.Error())
	}

	userService := service.NewUserService(connDb, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.UserServicePort)
	if err != nil {
		log.Fatal("failed while listening port: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	u.RegisterUserServiceServer(s, userService)

	log.Info("main: server running",
		logger.String("port", cfg.UserServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed while listening: %v", logger.Error(err))
	}
}
