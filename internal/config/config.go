package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"gopkg.in/yaml.v3"
)

type HTTPServer struct {
	Addr string `yaml:"addr"`
}

type Security struct {
	UserPasswordCost  int           `yaml:"user_password_cost"`
	TokenSecret       string        `yaml:"token_secret"`
	TokenLifetime     time.Duration `yaml:"token_lifetime"`
	SecretPasswordKey string        `yaml:"secret_password_key"`
}

type Config struct {
	GrpcServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc_server"`

	HttpServer HTTPServer `yaml:"http_server"`

	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`

	Security Security `yaml:"security"`

	Swagger struct {
		Enabled  bool   `yaml:"enabled"`
		JsonPath string `yaml:"json_path"`
	} `yaml:"swagger"`
}

func Read() (*Config, error) {
	var cfgPath string

	flag.StringVar(&cfgPath, "c", "", "Path to config file")

	flag.Parse()

	if cfgPath == "" {
		cfgPath = ""
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrConfigFileReading, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrConfigFileReading, err)
	}

	return &cfg, nil
}
