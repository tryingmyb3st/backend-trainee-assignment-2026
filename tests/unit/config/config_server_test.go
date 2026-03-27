package config_test

import (
	"backend-assignment-avito/internal/core/transport/server"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServerConfig(t *testing.T) {

	os.Setenv("HTTP_ADDR", ":5050")
	os.Setenv("HTTP_TIMEOUT", "10s")

	t.Cleanup(func() {
		os.Unsetenv("HTTP_ADDR")
		os.Unsetenv("HTTP_TIMEOUT")
	})

	want := server.Config{
		Addr:            ":5050",
		ShutdownTimeout: time.Duration(10 * time.Second),
	}

	cfg, err := server.NewConfig()

	require.NoError(t, err)
	require.Equal(t, want, cfg)
}

func TestServerConfigDefault(t *testing.T) {

	os.Setenv("HTTP_ADDR", ":5050")

	t.Cleanup(func() {
		os.Unsetenv("HTTP_ADDR")
	})

	want := server.Config{
		Addr:            ":5050",
		ShutdownTimeout: time.Duration(30 * time.Second),
	}

	cfg, err := server.NewConfig()

	require.NoError(t, err)
	require.Equal(t, want, cfg)
}

func TestServerConfigRequire(t *testing.T) {

	_, err := server.NewConfig()

	require.EqualError(t, err, "envconfig process: required key ADDR missing value")
}

func TestServerConfigMust(t *testing.T) {
	require.NotPanics(t, func() {
		os.Setenv("HTTP_ADDR", ":5050")
		os.Setenv("HTTP_TIMEOUT", "10s")

		t.Cleanup(func() {
			os.Unsetenv("HTTP_ADDR")
			os.Unsetenv("HTTP_TIMEOUT")
		})

		want := server.Config{
			Addr:            ":5050",
			ShutdownTimeout: time.Duration(10 * time.Second),
		}

		cfg := server.NewConfigMust()

		require.Equal(t, want, cfg)
	})
}

func TestServerConfigMustPanic(t *testing.T) {
	require.Panics(t, func() {
		server.NewConfigMust()
	})
}
