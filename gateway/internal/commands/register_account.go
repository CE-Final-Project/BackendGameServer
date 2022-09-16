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

type RegisterAccountCmdHandler interface {
	Handle(ctx context.Context, command *RegisterAccountCommand) error
}

type registerAccountHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewRegisterAccountHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *registerAccountHandler {
	return &registerAccountHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}

func (r *registerAccountHandler) Handle(ctx context.Context, command *RegisterAccountCommand) error {
	registerDto := &kafkaMessages.RegisterAccount{
		Username: command.RegisterDto.Username,
		Email:    command.RegisterDto.Email,
		Password: command.RegisterDto.Password,
	}

	dtoBytes, err := proto.Marshal(registerDto)
	if err != nil {
		return err
	}

	return r.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic: r.cfg.KafkaTopics.AccountRegister.TopicName,
		Value: dtoBytes,
		Time:  time.Now().UTC(),
	})
}
