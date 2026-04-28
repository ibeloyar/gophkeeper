package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr       string `env:"ADDRESS"`
	ConfigPath string `env:"CONFIG"`
}

type YAMLConfig struct {
	GrpcServer struct {
		Addr string `yaml:"addr"`
	} `yaml:"grpc_server"`
}

func Read() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Addr, "a", "", "The address metric grpc SERVER listen on")
	flag.StringVar(&cfg.ConfigPath, "c", "", "Path to config file")

	flag.Parse()

	if cfg.Addr == "" {
		if cfg.ConfigPath == "" {
			return nil, fmt.Errorf("without addr config file path required")
		}

		yamlConfig := YAMLConfig{}

		data, err := os.ReadFile(cfg.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", model.ErrConfigFileReading, err)
		}

		if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
			return nil, fmt.Errorf("%w: %w", model.ErrConfigFileReading, err)
		}
		fmt.Println(yamlConfig)

		if yamlConfig.GrpcServer.Addr != "" {
			cfg.Addr = yamlConfig.GrpcServer.Addr
		}
	}

	if cfg.Addr == "" {
		return nil, fmt.Errorf("no addr config file path provided")
	}

	return &cfg, nil
}
