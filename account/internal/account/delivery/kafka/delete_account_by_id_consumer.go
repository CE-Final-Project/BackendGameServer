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
)

func (s *accountMessageProcessor) processAccountDelete(ctx context.Context, r *kafka.Reader, m kafka.Message) {

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "accountMessageProcessor.processAccountDelete")
	defer span.Finish()

	msg := &kafkaMessages.DeleteAccount{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
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

	command := commands.NewDeleteAccountCommand(accUUID)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.as.Commands.DeleteAccount.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("DeleteAccount.Handle", err)
		return
	}

	s.commitMessage(ctx, r, m)
}
