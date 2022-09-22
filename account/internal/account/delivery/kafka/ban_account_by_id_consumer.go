package kafka

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/ce-final-project/backend_game_server/account/internal/account/commands"
	kafkaMessages "github.com/ce-final-project/backend_game_server/proto/kafka_messages"
	"github.com/ce-final-project/backend_rest_api/pkg/tracing"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

func (s *accountMessageProcessor) processBanAccountById(ctx context.Context, r *kafka.Reader, m kafka.Message) {

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "accountMessageProcessor.processBanAccountById")
	defer span.Finish()

	msg := &kafkaMessages.BanAccountById{}
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

	command := commands.NewBanAccountByIdCommand(accUUID, msg.GetIsBan(), time.Now())
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.as.Commands.BanAccountById.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("BanAccountById.Handle", err)
		return
	}

	s.commitMessage(ctx, r, m)
}
