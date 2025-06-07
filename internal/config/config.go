package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
}

type ServerConfig struct {
	Port        int
	MaxFileSize int64
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	Producer ProducerConfig
}

type ProducerConfig struct {
	MaxRetries   int
	RetryBackoff string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
