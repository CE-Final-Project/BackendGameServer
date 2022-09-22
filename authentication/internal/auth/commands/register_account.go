package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	kafkaMessages "github.com/ce-final-project/backend_game_server/proto/kafka_messages"
	"github.com/ce-final-project/backend_rest_api/pkg/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type RegisterCmdHandler interface {
	Handle(ctx context.Context, command *RegisterCommand) error
}

type registerAccountHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func NewRegisterAccountHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) *registerAccountHandler {
	return &registerAccountHandler{log: log, cfg: cfg, kafkaProducer: kafkaProducer}
}

func (r *registerAccountHandler) Handle(ctx context.Context, command *RegisterCommand) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authRegisterHandler.Handle")
	defer span.Finish()

	registerDto := &kafkaMessages.RegisterAccount{
		AccountID: command.AccountID.String(),
		PlayerID:  command.PlayerID,
		Username:  command.Username,
		Email:     command.Email,
		Password:  command.Password,
	}

	dtoBytes, err := proto.Marshal(registerDto)
	if err != nil {
		return err
	}

	return r.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   r.cfg.KafkaTopics.AccountRegister.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}
