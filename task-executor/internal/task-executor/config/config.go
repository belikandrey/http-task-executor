package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents full configuration for app.
type Config struct {
	// Env - environment name (local/prod).
	Env string `env-required:"true" yaml:"env"`
	// Postgres - postgres configuration presented in PostgresConfig.
	Postgres PostgresConfig `                    yaml:"postgres"`
	// LoggerConfig - logger configuration presented in LoggerConfig.
	LoggerConfig LoggerConfig `                    yaml:"logger"`
	// KafkaCfg - kafka configuration presented in KafkaConfig.
	KafkaCfg KafkaConfig `                    yaml:"kafka"`
	// ExternalServiceTimeout - external service timeout time.Duration.
	ExternalServiceTimeout time.Duration `env-required:"true" yaml:"external_service_timeout"`
	// ExecutorNumGoroutines - max number of goroutines in executor.
	ExecutorNumGoroutines int `env-required:"true" yaml:"executor_num_goroutines"`
}

// PostgresConfig - represents postgres configuration.
type PostgresConfig struct {
	// Host - represents postgres server host.
	Host string `env-required:"true" yaml:"host"`
	// Port - represents postgres server port.
	Port int64 `env-required:"true" yaml:"port"`
	// Name - represents postgres database name.
	Name string `env-required:"true" yaml:"name"`
	// User - represents postgres user.
	User string `env-required:"true" yaml:"user"`
	// Password - represents postgres password.
	Password string `env-required:"true" yaml:"password"`
	// SslMode - represents postgres sslmode.
	SslMode string `env-required:"true" yaml:"sslmode"`
	// Driver - represents postgres driver.
	Driver string `env-required:"true" yaml:"driver"`
	// MaxOpenConnections - number of max db connections.
	MaxOpenConnections int64 `env-required:"true" yaml:"max_open_connections"`
	// ConnectionMaxLifetime - max lifetime of connection.
	ConnectionMaxLifetime time.Duration `env-required:"true" yaml:"connection_max_lifetime"`
	// MaxIdleConnections - number of max db idle connections.
	MaxIdleConnections int64 `env-required:"true" yaml:"max_idle_connections"`
	// ConnectionMaxIdleTime - max lifetime of idle connection.
	ConnectionMaxIdleTime time.Duration `env-required:"true" yaml:"connection_max_idle_time"`
}

// LoggerConfig - represents logger configuration.
type LoggerConfig struct {
	// Filename - output logger filename.
	Filename string `env-required:"true" yaml:"filename"`
	// Level - logger level.
	Level string `env-required:"true" yaml:"level"`
	// Format - logger format (json/text).
	Format string `env-required:"true" yaml:"format"`
}

// KafkaConfig - represents kafka configuration.
type KafkaConfig struct {
	// Addresses - brokers addresses.
	Addresses []string `env-required:"true" yaml:"addresses"`
	// Topic - kafka topic.
	Topic string `env-required:"true" yaml:"topic"`
	// ConsumerGroup - kafka consumer group ID.
	ConsumerGroup string `env-required:"true" yaml:"consumer_group"`
}

// MustLoad - loads full app configuration and returns *Config.
// Panic on any error.
func MustLoad() *Config {
	path := getConfigPath()

	if path == "" {
		panic("config file path not found")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var config Config

	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		panic("can't read config file " + path + " : " + err.Error())
	}

	return &config
}

func getConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "config file path")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
