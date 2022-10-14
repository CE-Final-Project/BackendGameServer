package config

import (
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/postgres"
	"github.com/ce-final-project/backend_game_server/pkg/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
	_ "time/tzdata"
)

type Config struct {
	ServiceName string
	Logger      *logger.Config
	GRPC        GRPC
	Kafka       *kafkaClient.Config
	KafkaTopics KafkaTopics
	Postgres    *postgres.Config
	Redis       *redis.Config
	JWT         JWT
}

type GRPC struct {
	Port               string
	Development        bool
	AccountServicePort string
}

type KafkaTopics struct {
	AccountRegister kafkaClient.TopicConfig
	AccountUpdate   kafkaClient.TopicConfig
	ChangePassword  kafkaClient.TopicConfig
	AccountBan      kafkaClient.TopicConfig
	AccountDelete   kafkaClient.TopicConfig
}

type JWT struct {
	Secret     string `yaml:"secret"`
	ExpireTime string `yaml:"expireTime"`
}

func InitConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	config := &Config{}

	if err := viper.Unmarshal(config); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return config, nil
}
