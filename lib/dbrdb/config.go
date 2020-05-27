package dbrdb

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"
)

const redacted = "<redacted>"

// Config contains all properties that can be set-up for a database.
type Config struct {
	Host                  string        `env:"DBRDB_HOST" envDefault:"localhost"`
	Port                  uint          `env:"DBRDB_PORT" envDefault:"5432"`
	User                  string        `env:"DBRDB_USER" envDefault:"root"`
	Password              string        `env:"DBRDB_PASSWORD" envDefault:"root"`
	SSLMode               string        `env:"DBRDB_SSL_MODE" envDefault:"disable"`
	MaxConnectionAttempts uint          `env:"DBRDB_MAX_CONNECTION_ATTEMPTS" envDefault:"5"`
	TimeBetweenAttempts   time.Duration `env:"DBRDB_TIME_BETWEEN_ATTEMPTS" envDefault:"5s"`
	ConnectionTimeout     time.Duration `env:"DBRDB_CONNECTION_TIMEOUT" envDefault:"5s"`
	Encoding              string        `env:"DBRDB_ENCODING" envDefault:"UTF-8"`
	Driver                string        `env:"DBRDB_DRIVER" envDefault:"postgres"`
	MaxOpenConnections    uint          `env:"DBRDB_MAX_OPEN_CONNECTIONS" envDefault:"20"`
	MaxIdleConnections    uint          `env:"DBRDB_MAX_IDLE_CONNECTIONS" envDefault:"5"`
	ConnMaxLifetime       time.Duration `env:"DBRDB_CONNECTION_MAX_LIFETIME" envDefault:"1m"`
	DbName                string        `env:"DBRDB_DBNAME"`
	MigrationsDir         string        `env:"DBRDB_MIGRATIONS_DIR"`
}

// url returns the connection string for a database.
func (c *Config) url(redactPassword bool) string {
	password := c.Password
	if redactPassword {
		password = redacted
	}

	query := url.Values{}
	query.Set("sslmode", c.SSLMode)
	query.Set("connect_timeout", fmt.Sprintf("%v", int(c.ConnectionTimeout.Seconds())))
	query.Set("client_encoding", c.Encoding)

	datasource := url.URL{
		Scheme:   c.Driver,
		User:     url.UserPassword(c.User, password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(int(c.Port))),
		Path:     c.DbName,
		RawQuery: query.Encode(),
	}

	return datasource.String()
}
