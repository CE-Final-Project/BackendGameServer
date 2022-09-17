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
	Postgres    *postgres.Config
	Redis       *redis.Config
	GRPC        GRPC
	Kafka       *kafkaClient.Config
	JWT         JWT
}

type GRPC struct {
	Port        string
	Development bool
}

type JWT struct {
	Secret string `yaml:"secret"`
}

func InitConfig() (*Config, error) {

	//if err := initTimeZone(); err != nil {
	//	return nil, err
	//}
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

//func initTimeZone() error {
//	lct, err := time.LoadLocation("Asia/Bangkok")
//	if err != nil {
//		return errors.Wrap(err, "Initialize time zone error")
//	}
//	time.Local = lct
//	return nil
//}
