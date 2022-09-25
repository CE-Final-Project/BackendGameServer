package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	kafkaMessages "github.com/ce-final-project/backend_game_server/proto/kafka_messages"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type BanAccountCmdHandler interface {
	Handle(ctx context.Context, command *BanAccountCommand) error
}

type banAccountHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewBanAccountHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *banAccountHandler {
	return &banAccountHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (u *banAccountHandler) Handle(ctx context.Context, command *BanAccountCommand) error {
	banAccountDto := &kafkaMessages.BanAccountById{
		AccountID: command.BanAccountDto.AccountID.String(),
	}

	dtoBytes, err := proto.Marshal(banAccountDto)
	if err != nil {
		return err
	}

	return u.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic: u.cfg.KafkaTopics.AccountBan.TopicName,
		Value: dtoBytes,
		Time:  time.Now().UTC(),
	})
}
