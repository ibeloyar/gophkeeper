package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GrpcServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc_server"`

	HttpServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"http_server"`

	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`

	Security struct {
		UserPasswordCost  int           `yaml:"user_password_cost"`
		TokenSecret       string        `yaml:"token_secret"`
		TokenLifetime     time.Duration `yaml:"token_lifetime"`
		SecretPasswordKey string        `yaml:"secret_password_key"`
	} `yaml:"security"`

	Swagger struct {
		Enabled  bool   `yaml:"enabled"`
		JsonPath string `yaml:"json_path"`
	} `yaml:"swagger"`
}

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrConfigFileReading, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrConfigFileReading, err)
	}

	return &cfg, nil
}
