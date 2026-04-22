package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	err := os.WriteFile(path, []byte(yamlContent), 0644)
	require.NoError(t, err)

	cfg, err := Read(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, ":8081", cfg.GrpcServer.Addr)
	assert.Equal(t, ":8080", cfg.HttpServer.Addr)
	assert.Equal(t, "host=localhost port=5432 dbname=test", cfg.Database.DSN)

	assert.Equal(t, 10, cfg.Security.UserPasswordCost)
	assert.Equal(t, "test-secret", cfg.Security.TokenSecret)
	assert.Equal(t, 24*time.Hour, cfg.Security.TokenLifetime)
	assert.Equal(t, "test-key", cfg.Security.SecretPasswordKey)

	assert.Equal(t, true, cfg.Swagger.Enabled)
	assert.Equal(t, "/swagger.json", cfg.Swagger.JsonPath)
}

func TestConfig_Read_FileNotFound(t *testing.T) {
	_, err := Read("/this/path/does/not/exist.yaml")
	require.Error(t, err)

	assert.ErrorIs(t, err, model.ErrConfigFileReading)
}

func TestConfig_Read_InvalidYaml(t *testing.T) {
	yamlContent := `---
grpc_server:
  addr: ":8081"

http_server:
  addr: ":8080"  #  broken indent
     addr: ":8081"
`

	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.yaml")
	err := os.WriteFile(path, []byte(yamlContent), 0644)
	require.NoError(t, err)

	_, err = Read(path)
	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrConfigFileReading)
}
