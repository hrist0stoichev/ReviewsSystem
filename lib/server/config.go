package server

import (
	"time"
)

// Config contains all properties that can be set-up for an API server
type Config struct {
	Addr           string        `env:"SERVER_ADDR" envDefault:":8001"`
	ReadTimeout    time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"1s"`
	WriteTimeout   time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"1s"`
	IdleTimeout    time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"1s"`
	TLSCertificate string        `env:"SERVER_TLS_CERTIFICATE"`
	TLSKey         string        `env:"SERVER_TLS_KEY"`
}

// IsTLSEnabled indicates whether cert and key a provided in the configuration
func (c *Config) IsTLSEnabled() bool {
	return c.TLSCertificate != "" && c.TLSKey != ""
}
