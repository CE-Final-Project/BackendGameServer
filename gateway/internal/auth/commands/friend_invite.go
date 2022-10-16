package commands

import (
	"context"
	"github.com/ce-final-project/backend_game_server/gateway/config"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	kafkaMessages "github.com/ce-final-project/backend_game_server/proto/kafka_messages"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type FriendInviteCmdHandler interface {
	Handle(ctx context.Context, command *FriendInviteCommand) error
}

type friendInviteHandler struct {
	log           logger.Logger
	cfg           *config.Config
	kafkaProducer kafkaClient.Producer
}

func (f *friendInviteHandler) Handle(ctx context.Context, command *FriendInviteCommand) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "friendInviteHandler.Handle")
	defer span.Finish()

	friendInviteDto := &kafkaMessages.FriendInvite{
		PlayerID:       command.FriendInviteDto.PlayerID,
		FriendPlayerID: command.FriendInviteDto.FriendPlayerID,
	}

	dtoBytes, err := proto.Marshal(friendInviteDto)
	if err != nil {
		return err
	}

	return f.kafkaProducer.PublishMessage(ctx, kafka.Message{
		Topic:   f.cfg.KafkaTopics.FriendInvite.TopicName,
		Value:   dtoBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	})
}

func NewFriendInviteHandler(log logger.Logger, cfg *config.Config, kafkaProducer kafkaClient.Producer) FriendInviteCmdHandler {
	return &friendInviteHandler{
		log:           log,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
	}
}
