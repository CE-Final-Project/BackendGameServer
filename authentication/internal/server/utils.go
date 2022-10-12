package server

import (
	"context"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

func (s *Server) connectKafkaBrokers(ctx context.Context) error {
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

func (s *Server) initKafkaTopics(ctx context.Context) {
	controller, err := s.kafkaConn.Controller()
	if err != nil {
		s.log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	s.log.Infof("established new kafka controller connection: %s", controllerURI)

	accountRegisterTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.AccountRegister.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.AccountRegister.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.AccountRegister.ReplicationFactor,
	}

	accountUpdateTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.AccountUpdate.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.AccountUpdate.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.AccountUpdate.ReplicationFactor,
	}

	changePasswordTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.ChangePassword.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.ChangePassword.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.ChangePassword.ReplicationFactor,
	}

	banAccountTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.AccountBan.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.AccountBan.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.AccountBan.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		accountRegisterTopic,
		accountUpdateTopic,
		changePasswordTopic,
		banAccountTopic,
	); err != nil {
		s.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	s.log.Infof("kafka topics created or already exists: %+v", []kafka.TopicConfig{
		accountRegisterTopic,
		accountUpdateTopic,
		changePasswordTopic,
		banAccountTopic,
	})
}
