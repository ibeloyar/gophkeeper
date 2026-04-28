package config

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestRead_FlagAddrProvided(t *testing.T) {
	resetFlags()

	// -a, without -c
	os.Args = []string{"./test", "-a", "localhost:8080"}
	cfg, err := Read()
	require.NoError(t, err)

	assert.Equal(t, "localhost:8080", cfg.Addr)
	assert.Equal(t, "", cfg.ConfigPath)
}

func TestRead_FlagAddrAndConfigPath(t *testing.T) {
	resetFlags()

	os.Args = []string{"./test", "-a", "localhost:8080", "-c", "/tmp/test.yaml"}
	cfg, err := Read()
	require.NoError(t, err)

	assert.Equal(t, "localhost:8080", cfg.Addr)
	assert.Equal(t, "/tmp/test.yaml", cfg.ConfigPath)
}

func TestRead_NoAddr_NoConfigPath(t *testing.T) {
	resetFlags()

	// not -a and -c
	os.Args = []string{"./test"}
	cfg, err := Read()
	assert.EqualError(t, err, "without addr config file path required")
	assert.Nil(t, cfg)
}

func TestRead_NoAddr_ButConfigPath(t *testing.T) {
	resetFlags()

	// temp YAML‑файл
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(`
grpc_server:
  addr: "localhost:1234"
`)
	require.NoError(t, err)
	tmpfile.Close()

	os.Args = []string{"./test", "-c", tmpfile.Name()}

	cfg, err := Read()
	require.NoError(t, err)
	assert.Equal(t, "localhost:1234", cfg.Addr)
	assert.Equal(t, tmpfile.Name(), cfg.ConfigPath)
}

func TestRead_NoAddr_ConfigPath_ButNoAddrInYAML(t *testing.T) {
	resetFlags()

	// YAML без grpc_server.addr
	tmpfile, err := ioutil.TempFile("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(`
grpc_server:
  # addr: "..."
`)
	require.NoError(t, err)
	tmpfile.Close()

	os.Args = []string{"./test", "-c", tmpfile.Name()}

	cfg, err := Read()
	assert.EqualError(t, err, "no addr config file path provided")
	assert.Nil(t, cfg)
}

func TestRead_ConfigFileReading_Error(t *testing.T) {
	resetFlags()

	os.Args = []string{"./test", "-c", "/no/such/config.yaml"}

	cfg, err := Read()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), model.ErrConfigFileReading.Error())
}

func TestRead_YAML_Unmarshal_Error(t *testing.T) {
	resetFlags()

	// invalid YAML
	tmpfile, err := ioutil.TempFile("", "bad-config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString("::: invalid yaml :")
	require.NoError(t, err)
	tmpfile.Close()

	os.Args = []string{"./test", "-c", tmpfile.Name()}

	cfg, err := Read()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), model.ErrConfigFileReading.Error())
}

func TestRead_AddrFromFlagOverridesYAML(t *testing.T) {
	resetFlags()

	// addr
	tmpfile, err := ioutil.TempFile("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(`
grpc_server:
  addr: "localhost:1234"
`)
	require.NoError(t, err)
	tmpfile.Close()

	os.Args = []string{"./test", "-a", "127.0.0.1:5555", "-c", tmpfile.Name()}

	cfg, err := Read()
	require.NoError(t, err)

	// from flag, not from YAML
	assert.Equal(t, "127.0.0.1:5555", cfg.Addr)
}

func TestRead_NoAddr_FromYAML(t *testing.T) {
	resetFlags()

	// addr only in YAML
	tmpfile, err := ioutil.TempFile("", "config-*.yaml")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(`
grpc_server:
  addr: "localhost:1234"
`)
	require.NoError(t, err)
	tmpfile.Close()

	os.Args = []string{"./test", "-c", tmpfile.Name()}

	cfg, err := Read()
	require.NoError(t, err)

	assert.Equal(t, "localhost:1234", cfg.Addr)
}
