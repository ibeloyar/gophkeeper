package config

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withFlags(args ...string) func() {
	origArgs := os.Args
	origFlag := flag.CommandLine

	os.Args = append([]string{"test"}, args...)
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	return func() {
		os.Args = origArgs
		flag.CommandLine = origFlag
	}
}

func TestConfig_Read_Success(t *testing.T) {
	yamlContent := `---
grpc_server:
  addr: ":8081"
http_server:
  addr: ":8080"
database:
  dsn: "host=localhost port=5432 dbname=test"
security:
  user_password_cost: 10
  token_secret: "test-secret"
  token_lifetime: 24h
  secret_password_key: "test-key"
swagger:
  enabled: true
  json_path: "/swagger.json"
`

	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yamlContent), 0644))

	cleanup := withFlags("-c", path)
	defer cleanup()

	cfg, err := Read()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, ":8081", cfg.GrpcServer.Addr)
	assert.Equal(t, ":8080", cfg.HttpServer.Addr)
	assert.Equal(t, "host=localhost port=5432 dbname=test", cfg.Database.DSN)
	assert.Equal(t, 10, cfg.Security.UserPasswordCost)
	assert.Equal(t, 24*time.Hour, cfg.Security.TokenLifetime)
}

func TestConfig_Read_NoFlag(t *testing.T) {
	cleanup := withFlags()
	defer cleanup()

	_, err := Read()
	require.Error(t, err)

	assert.Contains(t, err.Error(), "config file read error")
}

func TestConfig_Read_FileNotFound(t *testing.T) {
	cleanup := withFlags("-c", "nonexistent.yaml")
	defer cleanup()

	_, err := Read()
	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrConfigFileReading)
}
