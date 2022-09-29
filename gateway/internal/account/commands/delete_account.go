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

type DeleteAccountCmdHandler interface {
	Handle(ctx context.Context, command *DeleteAccountCommand) error
}

type deleteAccountHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewDeleteAccountHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *deleteAccountHandler {
	return &deleteAccountHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (u *deleteAccountHandler) Handle(ctx context.Context, command *DeleteAccountCommand) error {
	deleteAccountDto := &kafkaMessages.DeleteAccount{
		AccountID: command.DeleteAccountDto.AccountID.String(),
	}

	dtoBytes, err := proto.Marshal(deleteAccountDto)
	if err != nil {
		return err
	}

	return u.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic: u.cfg.KafkaTopics.AccountDelete.TopicName,
		Value: dtoBytes,
		Time:  time.Now().UTC(),
	})
}
