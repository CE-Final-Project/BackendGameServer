package server

import (
	"context"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/pkg/errors"
)

func (s *server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	s.kafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.log.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (s *server) getConsumerGroupTopics() []string {
	return []string{
		s.cfg.KafkaTopics.AccountRegistered.TopicName,
		s.cfg.KafkaTopics.AccountUpdated.TopicName,
		s.cfg.KafkaTopics.AccountDeleted.TopicName,
	}
}
