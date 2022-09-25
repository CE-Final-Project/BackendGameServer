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

type ChangePasswordCmdHandler interface {
	Handle(ctx context.Context, command *ChangePasswordCommand) error
}

type changePasswordHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewChangePasswordHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *changePasswordHandler {
	return &changePasswordHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (u *changePasswordHandler) Handle(ctx context.Context, command *ChangePasswordCommand) error {
	changePasswordDto := &kafkaMessages.ChangePassword{
		AccountID:   command.ChangePasswordDto.AccountID.String(),
		OldPassword: command.ChangePasswordDto.OldPassword,
		NewPassword: command.ChangePasswordDto.NewPassword,
	}

	dtoBytes, err := proto.Marshal(changePasswordDto)
	if err != nil {
		return err
	}

	return u.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic: u.cfg.KafkaTopics.ChangePassword.TopicName,
		Value: dtoBytes,
		Time:  time.Now().UTC(),
	})
}
