package server

import (
	"context"
	"github.com/ce-final-project/backend_game_server/account/config"
	accountKafka "github.com/ce-final-project/backend_game_server/account/internal/account/delivery/kafka"
	"github.com/ce-final-project/backend_game_server/account/internal/account/repository"
	"github.com/ce-final-project/backend_game_server/account/internal/account/service"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/postgres"
	redisClient "github.com/ce-final-project/backend_game_server/pkg/redis"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	log         logger.Logger
	cfg         *config.Config
	v           *validator.Validate
	kafkaConn   *kafka.Conn
	db          *sqlx.DB
	redisClient redis.UniversalClient
	as          *service.AccountService
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{
		log: log,
		v:   validator.New(),
		cfg: cfg,
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	db, err := postgres.NewPostgresDatabase(s.cfg.Postgres)
	if err != nil {
		return errors.Wrap(err, "postgres.NewPostgresDatabase")
	}
	s.log.Infof("postgres connected: %v", db.Stats().OpenConnections)
	s.db = db

	s.redisClient = redisClient.NewUniversalRedisClient(s.cfg.Redis)
	defer s.redisClient.Close()
	s.log.Infof("Redis connected: %+v", s.redisClient.PoolStats())

	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close()

	pgRepo, err := repository.NewAccountRepository(s.log, s.db)
	if err != nil {
		return errors.Wrap(err, "repository.NewAccountRepository")
	}
	redisRepo := repository.NewAccountCacheRepository(s.log, s.redisClient)

	s.as = service.NewAccountService(s.log, s.cfg, pgRepo, redisRepo)

	accountMessageProcessor := accountKafka.NewAccountMessageProcessor(s.log, s.cfg, s.v, s.as)

	s.log.Info("Starting Account Kafka consumers")
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), accountKafka.PoolSize, accountMessageProcessor.ProcessMessages)

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close()

	closeGrpcServer, grpcServer, err := s.NewAuthGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
