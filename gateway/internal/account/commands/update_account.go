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

type UpdateAccountCmdHandler interface {
	Handle(ctx context.Context, command *UpdateAccountCommand) error
}

type updateAccountHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewUpdateAccountHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *updateAccountHandler {
	return &updateAccountHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (u *updateAccountHandler) Handle(ctx context.Context, command *UpdateAccountCommand) error {
	updateAccountDto := &kafkaMessages.UpdateAccount{
		AccountID: command.UpdateDto.AccountID.String(),
		Username:  command.UpdateDto.Username,
		Email:     command.UpdateDto.Email,
	}

	dtoBytes, err := proto.Marshal(updateAccountDto)
	if err != nil {
		return err
	}

	return u.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic: u.cfg.KafkaTopics.AccountUpdate.TopicName,
		Value: dtoBytes,
		Time:  time.Now().UTC(),
	})
}
