package kafka

import (
	"context"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"gitlab.com/micro/user_service/config"
	"gitlab.com/micro/user_service/pkg/logger"
	"gitlab.com/micro/user_service/pkg/messagebroker"
)

type KafkaProduce struct {
	kafkaWriter *kafka.Writer
	log         logger.Logger
}

func NewKafkaProducer(conf config.Config, log logger.Logger, topic string) messagebroker.Producer {
	connStr := "kafka:9092"

	return &KafkaProduce{
		kafkaWriter: &kafka.Writer{
			Addr:         kafka.TCP(connStr),
			Topic:        topic,
			BatchTimeout: time.Millisecond * 10,
		},
		log: log,
	}
}

func (p *KafkaProduce) Start() error {
	return nil
}

func (p *KafkaProduce) Stop() error {
	err := p.kafkaWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProduce) Produce(key,body []byte, logBody string) error {
	message := kafka.Message{
		Key: key,
		Value: body,
	}

	if err :=  p.kafkaWriter.WriteMessages(context.Background(), message); err != nil {
		return err
	}

	return nil
}