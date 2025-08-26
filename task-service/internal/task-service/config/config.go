package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env                    string           `yaml:"env" env-required:"true"`
	ServerConfig           HttpServerConfig `yaml:"http_server"`
	Postgres               PostgresConfig   `yaml:"postgres"`
	LoggerConfig           LoggerConfig     `yaml:"logger"`
	ExternalServiceTimeout time.Duration    `yaml:"external_service_timeout"`
	KafkaCfg               KafkaConfig      `yaml:"kafka"`
}

type HttpServerConfig struct {
	Host         string        `yaml:"host" env-required:"true"`
	Port         int64         `yaml:"port" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
}

type PostgresConfig struct {
	Host                  string        `yaml:"host" env-required:"true"`
	Port                  int64         `yaml:"port" env-required:"true"`
	Name                  string        `yaml:"name" env-required:"true"`
	User                  string        `yaml:"user" env-required:"true"`
	Password              string        `yaml:"password" env-required:"true"`
	SslMode               string        `yaml:"sslmode" env-required:"true"`
	Driver                string        `yaml:"driver" env-required:"true"`
	MaxOpenConnections    int64         `yaml:"max_open_connections" env-required:"true"`
	ConnectionMaxLifetime time.Duration `yaml:"connection_max_lifetime" env-required:"true"`
	MaxIdleConnections    int64         `yaml:"max_idle_connections" env-required:"true"`
	ConnectionMaxIdleTime time.Duration `yaml:"connection_max_idle_time" env-required:"true"`
}

type LoggerConfig struct {
	Filename string `yaml:"filename" env-required:"true"`
	Level    string `yaml:"level" env-required:"true"`
	Format   string `yaml:"format" env-required:"true"`
}

type KafkaConfig struct {
	Addresses []string `yaml:"addresses" env-required:"true"`
	Topic     string   `yaml:"topic" env-required:"true"`
}

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
