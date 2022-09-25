package kafka

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/ce-final-project/backend_game_server/account/internal/account/commands"
	"github.com/ce-final-project/backend_game_server/pkg/tracing"
	kafkaMessages "github.com/ce-final-project/backend_game_server/proto/kafka_messages"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

func (s *accountMessageProcessor) processAccountRegister(ctx context.Context, r *kafka.Reader, m kafka.Message) {

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "accountMessageProcessor.processAccountRegister")
	defer span.Finish()

	var msg kafkaMessages.RegisterAccount
	if err := proto.Unmarshal(m.Value, &msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	accUUID, err := uuid.FromString(msg.GetAccountID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	command := commands.NewCreateAccountCommand(accUUID, msg.GetPlayerID(), msg.GetUsername(), msg.GetEmail(), msg.GetPassword(), false, time.Now(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.as.Commands.CreateAccount.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("CreateAccount.Handle", err)
		return
	}

	s.commitMessage(ctx, r, m)
}
