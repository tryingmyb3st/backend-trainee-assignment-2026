package config_test

import (
	postgres_pool "backend-assignment-avito/internal/core/repository/postgres"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRepoConfig(t *testing.T) {

	os.Setenv("POSTGRES_USER", "user")
	os.Setenv("POSTGRES_PASSWORD", "pass")
	os.Setenv("POSTGRES_DB", "db")
	os.Setenv("POSTGRES_HOST", "host")

	t.Cleanup(func() {
		os.Unsetenv("POSTGRES_USER")
		os.Unsetenv("POSTGRES_PASSWORD")
		os.Unsetenv("POSTGRES_DB")
		os.Unsetenv("POSTGRES_HOST")
	})

	want := postgres_pool.Config{
		User:     "user",
		Password: "pass",
		Database: "db",
		Host:     "host",
		Timeout:  time.Duration(10 * time.Second),
	}

	cfg, err := postgres_pool.NewConfig()

	require.NoError(t, err)
	require.Equal(t, want, cfg)
}

func TestRepoConfigRequire(t *testing.T) {

	_, err := postgres_pool.NewConfig()

	require.EqualError(t, err, "envconfig process: required key USER missing value")
}

func TestRepoConfigMust(t *testing.T) {
	require.NotPanics(t, func() {
		os.Setenv("POSTGRES_USER", "user")
		os.Setenv("POSTGRES_PASSWORD", "pass")
		os.Setenv("POSTGRES_DB", "db")
		os.Setenv("POSTGRES_HOST", "host")

		t.Cleanup(func() {
			os.Unsetenv("POSTGRES_USER")
			os.Unsetenv("POSTGRES_PASSWORD")
			os.Unsetenv("POSTGRES_DB")
			os.Unsetenv("POSTGRES_HOST")
		})

		want := postgres_pool.Config{
			User:     "user",
			Password: "pass",
			Database: "db",
			Host:     "host",
			Timeout:  time.Duration(10 * time.Second),
		}

		cfg := postgres_pool.NewConfigMust()

		require.Equal(t, want, cfg)
	})
}

func TestRepoConfigMustPanic(t *testing.T) {
	require.Panics(t, func() {
		postgres_pool.NewConfigMust()
	})
}
