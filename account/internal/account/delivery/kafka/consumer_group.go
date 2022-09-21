package kafka

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	"github.com/ce-final-project/backend_game_server/account/internal/service"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-playground/validator"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	PoolSize = 30
)

type accountMessageProcessor struct {
	log logger.Logger
	cfg *config.Config
	v   *validator.Validate
	as  *service.AccountService
}

func NewAccountMessageProcessor(log logger.Logger, cfg *config.Config, v *validator.Validate, as *service.AccountService) *accountMessageProcessor {
	return &accountMessageProcessor{log: log, cfg: cfg, v: v, as: as}
}

func (s *accountMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			s.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		s.logProcessMessage(m, workerID)

		switch m.Topic {
		case s.cfg.KafkaTopics.AccountRegister.TopicName:
			s.processAccountRegisterd(ctx, r, m)
		case s.cfg.KafkaTopics.AccountUpdate.TopicName:
			s.processAccountUpdated(ctx, r, m)
		case s.cfg.KafkaTopics.AccountDelete.TopicName:
			s.processAccountDeleted(ctx, r, m)
		}
	}
}
