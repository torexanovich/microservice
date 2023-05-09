package kafka

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"gitlab.com/micro/post_service/config"
	"gitlab.com/micro/post_service/kafka/handler"
	"gitlab.com/micro/post_service/pkg/logger"
	"gitlab.com/micro/post_service/pkg/messagebroker"
	"gitlab.com/micro/post_service/storage"
)

type KafkaConsumer struct {
	log           logger.Logger
	KafkaConsumer *kafka.Reader
	KafkaHandler  *handler.KafkaHandler
}

func (k KafkaConsumer) Start() {
	fmt.Println("Consumer started")
	for {	
		m, err := k.KafkaConsumer.ReadMessage(context.Background())
		fmt.Println("HERE IT IS:", m, err)
		if err != nil {
			k.log.Error("Error while consuming the message", logger.Error(err))
			break
		}
		
		err = k.KafkaHandler.Handle(m.Value)
		if err != nil {
			k.log.Error("Failed to handle the consumed topic: ", logger.String("on topic", m.Topic))
		} else {
			fmt.Println("\nCheers! \n ")
			k.log.Info("Successfully consumed message",
				logger.String("on topic", m.Topic),
				logger.String("message", "success"))
			fmt.Println("\n COOL!\n ")
		}
	}

	err := k.KafkaConsumer.Close()
	if err != nil {
		k.log.Error("Failed to close consumer", logger.Error(err))
	}
}

func NewKafkaConsumer(db *sqlx.DB, conf *config.Config, log logger.Logger, topic string) messagebroker.Consumer {
	connStr := "localhost:9092"
	return &KafkaConsumer{
		KafkaConsumer: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers:  []string{connStr},
			Topic:    topic,
			MinBytes: 10e3, //10KB
			MaxBytes: 10e6, //10MB

		}),
		KafkaHandler: handler.NewKafkaHandler(*conf, log, storage.NewStoragePg(db)),
		log: log,
	}
}
