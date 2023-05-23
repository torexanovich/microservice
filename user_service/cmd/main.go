package main

import (
	"fmt"
	"net"

	"gitlab.com/micro/user_service/config"
	u "gitlab.com/micro/user_service/genproto/user"
	"gitlab.com/micro/user_service/kafka"
	"gitlab.com/micro/user_service/pkg/db"
	"gitlab.com/micro/user_service/pkg/logger"
	"gitlab.com/micro/user_service/pkg/messagebroker"
	"gitlab.com/micro/user_service/service"
	grpcclient "gitlab.com/micro/user_service/service/grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "golang")
	defer logger.Cleanup(log)

	log.Info("main:sqlxConfig",
		logger.String("'host", cfg.CommentServiceHost),
		logger.String("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDb, err := db.ConnectToDB(cfg)
	if err != nil {
		fmt.Println("failed to connect to database", err)
	}

	produceMap := make(map[string]messagebroker.Producer)
	topic := "user"
	userTopicProduce := kafka.NewKafkaProducer(cfg, log, topic)
	defer func() {
		err := userTopicProduce.Stop()
		if err != nil {
			log.Fatal("Failed to stopping Kafka", logger.Error(err))
		}
	}()
	produceMap["user"] = userTopicProduce

	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		fmt.Println("failed while grpc client", err.Error())
	}

	uService := service.NewUserService(connDb, log, produceMap, grpcClient)
	lis, err := net.Listen("tcp", cfg.UserServicePort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	s := grpc.NewServer()
	reflection.Register(s)
	u.RegisterUserServiceServer(s, uService)
	log.Info("main: server is running",
		logger.String("port", cfg.UserServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listenning: %v", logger.Error(err))
	}
}
