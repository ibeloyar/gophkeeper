package config

import (
	"errors"
	"os"
	"testing"

	"github.com/ibeloyar/gophkeeper/internal/model"
)

func TestRead_Success(t *testing.T) {
	yamlContent := []byte(`
grpc_server:
  addr: "localhost:8080"
`)

	tmpFile := createTempFile(t, yamlContent)
	defer os.Remove(tmpFile)

	cfg, err := Read(tmpFile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.GrpcServer.Addr != "localhost:8080" {
		t.Errorf("expected addr 'localhost:8080', got '%s'", cfg.GrpcServer.Addr)
	}
}

func TestRead_FileNotFound(t *testing.T) {
	_, err := Read("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}

	if !errors.Is(err, model.ErrConfigFileReading) {
		t.Errorf("expected model.ErrConfigFileReading wrapped error, got %v", err)
	}
}

func TestRead_InvalidYAML(t *testing.T) {
	invalidYAML := []byte("invalid: yaml: content")
	tmpFile := createTempFile(t, invalidYAML)
	defer os.Remove(tmpFile)

	_, err := Read(tmpFile)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}

	if !errors.Is(err, model.ErrConfigFileReading) {
		t.Errorf("expected model.ErrConfigFileReading wrapped error, got %v", err)
	}
}

func TestRead_EmptyFile(t *testing.T) {
	emptyContent := []byte("")
	tmpFile := createTempFile(t, emptyContent)
	defer os.Remove(tmpFile)

	cfg, err := Read(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.GrpcServer.Addr != "" {
		t.Errorf("expected empty addr, got '%s'", cfg.GrpcServer.Addr)
	}
}

func createTempFile(t *testing.T, content []byte) string {
	tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}
