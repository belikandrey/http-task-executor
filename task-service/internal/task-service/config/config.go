package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

// Config represents full configuration for app.
type Config struct {
	// Env - environment name (local/prod).
	Env string `yaml:"env" env-required:"true"`
	// ServerConfig - http server configuration presented in HttpServerConfig.
	ServerConfig HttpServerConfig `yaml:"http_server"`
	// Postgres - postgres configuration presented in PostgresConfig.
	Postgres PostgresConfig `yaml:"postgres"`
	// LoggerConfig - logger configuration presented in LoggerConfig.
	LoggerConfig LoggerConfig `yaml:"logger"`
	// ExternalServiceTimeout - external service timeout time.Duration.
	ExternalServiceTimeout time.Duration `yaml:"external_service_timeout"`
	// KafkaCfg - kafka configuration presented in KafkaConfig.
	KafkaCfg KafkaConfig `yaml:"kafka"`
}

// HttpServerConfig represents http server configuration.
type HttpServerConfig struct {
	// Host - http-server host.
	Host string `yaml:"host" env-required:"true"`
	// Port - http-server port.
	Port int64 `yaml:"port" env-required:"true"`
	// ReadTimeout - http-server read timeout.
	ReadTimeout time.Duration `yaml:"read_timeout" env-default:"5s"`
	// WriteTimeout - http-server write timeout.
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
}

// PostgresConfig - represents postgres configuration.
type PostgresConfig struct {
	// Host - represents postgres server host.
	Host string `yaml:"host" env-required:"true"`
	// Port - represents postgres server port.
	Port int64 `yaml:"port" env-required:"true"`
	// Name - represents postgres database name.
	Name string `yaml:"name" env-required:"true"`
	// User - represents postgres user.
	User string `yaml:"user" env-required:"true"`
	// Password - represents postgres password.
	Password string `yaml:"password" env-required:"true"`
	// SslMode - represents postgres sslmode.
	SslMode string `yaml:"sslmode" env-required:"true"`
	// Driver - represents postgres driver.
	Driver string `yaml:"driver" env-required:"true"`
	// MaxOpenConnections - number of max db connections.
	MaxOpenConnections int64 `yaml:"max_open_connections" env-required:"true"`
	// ConnectionMaxLifetime - max lifetime of connection.
	ConnectionMaxLifetime time.Duration `yaml:"connection_max_lifetime" env-required:"true"`
	// MaxIdleConnections - number of max db idle connections.
	MaxIdleConnections int64 `yaml:"max_idle_connections" env-required:"true"`
	// ConnectionMaxIdleTime - max lifetime of idle connection.
	ConnectionMaxIdleTime time.Duration `yaml:"connection_max_idle_time" env-required:"true"`
}

// LoggerConfig - represents logger configuration.
type LoggerConfig struct {
	// Filename - output logger filename.
	Filename string `yaml:"filename" env-required:"true"`
	// Level - logger level.
	Level string `yaml:"level" env-required:"true"`
	// Format - logger format (json/text).
	Format string `yaml:"format" env-required:"true"`
}

// KafkaConfig - represents kafka configuration.
type KafkaConfig struct {
	// Addresses - brokers addresses.
	Addresses []string `yaml:"addresses" env-required:"true"`
	// Topic - kafka topic.
	Topic string `yaml:"topic" env-required:"true"`
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
