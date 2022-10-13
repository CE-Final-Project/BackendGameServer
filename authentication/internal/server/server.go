package server

import (
	"context"
	"github.com/ce-final-project/backend_game_server/authentication/config"
	accountRepository "github.com/ce-final-project/backend_game_server/authentication/internal/account/repository"
	accountService "github.com/ce-final-project/backend_game_server/authentication/internal/account/service"
	roleRepository "github.com/ce-final-project/backend_game_server/authentication/internal/role/repository"
	roleService "github.com/ce-final-project/backend_game_server/authentication/internal/role/service"
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

type Server struct {
	log         logger.Logger
	cfg         *config.Config
	v           *validator.Validate
	acc         *accountService.AccountService
	rs          *roleService.RoleService
	db          *sqlx.DB
	redisClient redis.UniversalClient
	kafkaConn   *kafka.Conn
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{
		log: log,
		v:   validator.New(),
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// TODO: init db and repository

	db, err := postgres.NewPostgresDatabase(s.cfg.Postgres)
	if err != nil {
		return errors.Wrap(err, "postgres.NewPostgresDatabase")
	}
	s.log.Infof("postgres connected: %v", db.Stats().OpenConnections)
	s.db = db
	defer s.db.Close()

	s.redisClient = redisClient.NewUniversalRedisClient(s.cfg.Redis)
	defer s.redisClient.Close()
	s.log.Infof("Redis connected: %+v", s.redisClient.PoolStats())

	accountRepo := accountRepository.NewAccountRepository(s.db, s.log)
	cacheRepo := accountRepository.NewCacheRepository(s.redisClient, s.log)

	roleRepo := roleRepository.NewRoleRepository(s.db, s.log)

	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close()

	s.acc = accountService.NewAccountService(s.log, s.cfg, accountRepo, cacheRepo)
	s.rs = roleService.NewRoleService(s.log, s.cfg, roleRepo)

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close()

	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}

	closeGrpcServer, grpcServer, err := s.NewAuthGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer()

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
