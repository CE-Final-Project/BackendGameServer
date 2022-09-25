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

func (s *accountMessageProcessor) processAccountUpdate(ctx context.Context, r *kafka.Reader, m kafka.Message) {

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "accountMessageProcessor.processUpdateAccount")
	defer span.Finish()

	msg := &kafkaMessages.AccountUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}
	account := msg.GetAccount()
	accUUID, err := uuid.FromString(account.GetAccountID())
	if err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	command := commands.NewUpdateAccountCommand(accUUID, account.GetUsername(), account.GetEmail(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.as.Commands.UpdateAccount.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("UpdateAccount.Handle", err)
		return
	}

	s.commitMessage(ctx, r, m)
}
