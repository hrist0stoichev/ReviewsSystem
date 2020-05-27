package server

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
)

type Server interface {
	ListenAndServe()
	Shutdown(ctx context.Context)
}

type apiServer struct {
	config *Config
	server *http.Server
	logger log.Logger
}

func New(config *Config, handler http.Handler, logger log.Logger) (Server, error) {
	if config == nil || handler == nil || logger == nil {
		return nil, errors.New("config, handler, or logger is nil")
	}

	apiServer := &apiServer{
		config: config,
		server: &http.Server{
			Handler:      handler,
			Addr:         config.Addr,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		},
		logger: logger,
	}

	// TODO: Consider changing MinVersion to 1.3
	if config.IsTLSEnabled() {
		apiServer.server.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}

	return apiServer, nil
}

func (s *apiServer) ListenAndServe() {
	if err := s.listenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.WithError(err).Fatalln("Failed to listed and serve")
	}

	s.logger.Infoln("API server stopped listening for incoming connections")
}

func (s *apiServer) listenAndServe() error {
	if s.config.IsTLSEnabled() {
		s.logger.WithField("addr", s.config.Addr).Infoln("Starting API server with TLS")
		return s.server.ListenAndServeTLS(s.config.TLSCertificate, s.config.TLSKey)
	}

	s.logger.WithField("addr", s.config.Addr).Warnln("Starting API server without TLS")
	return s.server.ListenAndServe()
}

// Shutdown shuts down the server gracefully
func (s *apiServer) Shutdown(ctx context.Context) {
	s.logger.Infoln("API server shutdown started")
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.WithError(err).Errorln("Error while shutting down API server")
	}

	s.logger.Infoln("API server shutdown completed")
}
