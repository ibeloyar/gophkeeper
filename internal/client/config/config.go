package config

import (
	"fmt"
	"os"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GrpcServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc_server"`
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
