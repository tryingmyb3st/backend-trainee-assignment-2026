package config_test

import (
	"backend-assignment-avito/internal/core/logger"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerConfig(t *testing.T) {

	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("LOG_FOLDER", "tests/unit/conig/log")

	t.Cleanup(func() {
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_FOLDER")
	})

	want := logger.Config{
		Level:  "DEBUG",
		Folder: "tests/unit/conig/log",
	}

	cfg, err := logger.NewConfig()

	require.NoError(t, err)
	require.Equal(t, want, cfg)
}

func TestLoggerConfigDefault(t *testing.T) {

	os.Setenv("LOG_FOLDER", "tests/unit/conig/log")

	t.Cleanup(func() {
		os.Unsetenv("LOG_FOLDER")
	})

	want := logger.Config{
		Level:  "INFO",
		Folder: "tests/unit/conig/log",
	}

	cfg, err := logger.NewConfig()

	require.NoError(t, err)
	require.Equal(t, want, cfg)
}

func TestLoggerConfigRequire(t *testing.T) {

	_, err := logger.NewConfig()

	require.EqualError(t, err, "envconfig process: required key FOLDER missing value")
}

func TestLoggerConfigMust(t *testing.T) {
	require.NotPanics(t, func() {
		os.Setenv("LOG_LEVEL", "DEBUG")
		os.Setenv("LOG_FOLDER", "tests/unit/conig/log")

		t.Cleanup(func() {
			os.Unsetenv("LOG_LEVEL")
			os.Unsetenv("LOG_FOLDER")
		})

		want := logger.Config{
			Level:  "DEBUG",
			Folder: "tests/unit/conig/log",
		}

		cfg := logger.NewConfigMust()

		require.Equal(t, want, cfg)
	})
}

func TestLoggerConfigMustPanic(t *testing.T) {
	require.Panics(t, func() {
		logger.NewConfigMust()
	})
}
