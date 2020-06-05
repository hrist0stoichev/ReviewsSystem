package etc

import (
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/lib/dbrdb"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/lib/server"
)

type Config struct {
	Server       server.Config
	Database     dbrdb.Config
	Logging      log.Config
	Tokens       TokensConfig
	FacebookAuth FacebookAuthConfig
	Email        EmailConfig
}

type TokensConfig struct {
	ValidFor   time.Duration `env:"TOKENS_VALID_FOR"`
	SigningKey string        `env:"TOKENS_SIGNING_KEY"`
}

type FacebookAuthConfig struct {
	ClientId     string   `env:"FACEBOOK_CLIENT_ID"`
	ClientSecret string   `env:"FACEBOOK_CLIENT_SECRET"`
	RedirectURL  string   `env:"FACEBOOK_REDIRECT_URL"`
	Scopes       []string `env:"FACEBOOK_SCOPES"`
}

type EmailConfig struct {
	SMTPHost              string `env:"EMAIL_SMTP_HOST"`
	SMTPPort              string `env:"EMAIL_SMTP_PORT"`
	Username              string `env:"EMAIL_SMTP_USERNAME"`
	Password              string `env:"EMAIL_SMTP_PASSWORD"`
	ConfirmationEndpoint  string `env:"EMAIL_CONFIRMATION_ENDPOINT"`
	RedirectionEndpoint   string `env:"EMAIL_REDIRECTION_ENDPOINT"`
	SkipEmailVerification bool   `env:"SKIP_EMAIL_VERIFICATION"`
}

func GetConfig() (*Config, error) {
	cfg := new(Config)

	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "could not parse env variables")
	}

	return cfg, nil
}
