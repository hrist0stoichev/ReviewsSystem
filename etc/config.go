package etc

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/lib/dbrdb"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/lib/server"
)

type Config struct {
	Server   server.Config
	Database dbrdb.Config
	Logging  log.Config
}

func GetConfig() (*Config, error) {
	cfg := new(Config)

	if err := env.Parse(cfg); err != nil {
		return nil, errors.Wrap(err, "could not parse env variables")
	}

	return cfg, nil
}
