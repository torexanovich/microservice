package handler

import (
	"gitlab.com/micro/post_service/config"
	"gitlab.com/micro/post_service/pkg/logger"
	"gitlab.com/micro/post_service/storage"
	p "gitlab.com/micro/post_service/genproto/post"
)

type KafkaHandler struct {
	log logger.Logger
	cfg config.Config
	storage storage.IStorage
}

func NewKafkaHandler(config config.Config, log logger.Logger, storage storage.IStorage) *KafkaHandler {
	return &KafkaHandler{
		log: log,
		cfg: config,
		storage: storage,
	}
}

func (h *KafkaHandler) Handle(value []byte) error {
	post := p.PostRequest{}
	err := post.Unmarshal(value)
	if err != nil {
		return err
	}

	_, err = h.storage.Post().CreatePost(&post)
	if err != nil {
		return err
	}

	return nil
}
