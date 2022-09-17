package config

import (
	"github.com/ce-final-project/backend_game_server/pkg/kafka"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
	_ "time/tzdata"
)

type Config struct {
	ServiceName string
	Logger      *logger.Config
	HTTP        Http
	GRPC        Grpc
	KafkaTopics KafkaTopics
	JWT         JWT
}

type JWT struct {
	HeaderAuthorization string `yaml:"headerAuthorization"`
	Secret              string `yaml:"secret"`
}

type Http struct {
	Port                string
	Development         bool
	BasePath            string
	DebugHeaders        bool
	HttpClientDebug     bool
	DebugErrorsResponse bool
	IgnoreLogUrls       []string
}

type Grpc struct {
	AuthServicePort string
}

type KafkaTopics struct {
	AccountRegister kafka.TopicConfig
	AccountUpdate   kafka.TopicConfig
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
