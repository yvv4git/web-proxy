package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for Server overrides only
	os.Setenv("SERVER_HOST", "env-host")
	os.Setenv("SERVER_PORT", "8080")
	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
	}()

	tmpFile, err := os.CreateTemp("", "config-*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
log_level = "info"

[server]
host = "file-host"
port = 1234

[auth]
predifined_auth = { accounts = [
  { username = "file-user", password = "file-pass" }
]}
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	var cfg Config
	err = Load(tmpFile.Name(), &cfg)
	require.NoError(t, err)

	// Expected config matches file values since Server fields don't support env overrides
	expected := Config{
		LogLevel: "info",
		Server: Server{
			Host: "file-host",
			Port: 1234,
		},
		Auth: Auth{
			PredifinedAuth: PredifinedAuth{
				Accounts: []Account{
					{
						Username: "file-user",
						Password: "file-pass",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, cfg) {
		t.Logf("Expected: %+v", expected)
		t.Logf("Actual: %+v", cfg)
		t.FailNow()
	}
}
