package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	Database      DatabaseConfig      `mapstructure:"database"`
	Redis         RedisConfig         `mapstructure:"redis"`
	RabbitMQ      RabbitMQConfig      `mapstructure:"rabbitmq"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	Meilisearch   MeilisearchConfig   `mapstructure:"meilisearch"`
	Asynq         AsynqConfig         `mapstructure:"asynq"`
	Storage       StorageConfig       `mapstructure:"storage"`
	Logger        LoggerConfig        `mapstructure:"logger"`
}

type ServerConfig struct {
	Port   string `mapstructure:"port"`
	Mode   string `mapstructure:"mode"`
	NodeID int64  `mapstructure:"node_id"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMQConfig struct {
	URL string `mapstructure:"url"`
}

type ElasticsearchConfig struct {
	Addresses []string `mapstructure:"addresses"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`
}

type MeilisearchConfig struct {
	Host   string `mapstructure:"host"`
	APIKey string `mapstructure:"api_key"`
}

type AsynqConfig struct {
	Addr        string `mapstructure:"addr"`
	Password    string `mapstructure:"password"`
	DB          int    `mapstructure:"db"`
	Concurrency int    `mapstructure:"concurrency"`
}

type StorageConfig struct {
	Type  string      `mapstructure:"type"`
	Local LocalConfig `mapstructure:"local"`
}

type LocalConfig struct {
	Path string `mapstructure:"path"`
	URL  string `mapstructure:"url"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
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
