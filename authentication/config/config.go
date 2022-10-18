package config

import (
	"crypto/rsa"
	kafkaClient "github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/ce-final-project/backend_game_server/pkg/postgres"
	"github.com/ce-final-project/backend_game_server/pkg/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
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
	FriendInvite    kafkaClient.TopicConfig
}

type JWT struct {
	ExpireTime string `yaml:"expireTime"`
	VerifyKey  *rsa.PublicKey
	SignKey    *rsa.PrivateKey
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

	var signBytes []byte
	signBytes, err = os.ReadFile(viper.GetString("jwt.privateKeyPath"))
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile.signBytes")
	}

	config.JWT.SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseRSAPrivateKeyFromPEM")
	}

	var verifyBytes []byte
	verifyBytes, err = os.ReadFile(viper.GetString("jwt.pubKeyPath"))
	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile.verifyBytes")
	}

	config.JWT.VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseRSAPublicKeyFromPEM")
	}

	return config, nil
}
