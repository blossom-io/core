package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App         `yaml:"app"`
		Connections `yaml:"connections"`
	}

	App struct {
		Port     string `env-required:"true" yaml:"port" env:"BLOSSOM_CORE_PORT"`
		LogLevel string `env-required:"true" yaml:"log_level" env:"BLOSSOM_CORE_LOG_LEVEL"`
	}

	Connections struct {
		Postgres `yaml:"postgres"`
		Twitch   `yaml:"twitch"`
	}

	Postgres struct {
		URL string `env-required:"true" yaml:"url" env:"PG_URL"`
	}

	Twitch struct {
		ClientID                 string `env-required:"true" yaml:"client_id" env:"TWITCH_CLIENT_ID"`
		ClientSecret             string `env-required:"true" yaml:"client_secret" env:"TWITCH_CLIENT_SECRET"`
		AuthRedirectURL          string `env-required:"true" yaml:"auth_redirect_url" env:"TWITCH_AUTH_REDIRECT_URL"`
		SubchatInviteRedirectURL string `env-required:"true" yaml:"subchat_invite_redirect_url" env:"TWITCH_SUBCHAT_INVITE_REDIRECT_URL"`
	}
)

// New returns app config.
func New() (*Config, error) {
	c := &Config{}

	err := cleanenv.ReadEnv(c)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return c, nil
}
