package dbrdb

import (
	"fmt"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
)

const fileSchemaFormat = "file://%s"

// Database interface allows clients to connect to, migrate, and query a dbr compliant database (mysql, postgres, sqlite).
type Database interface {
	Init() error
	Migrate() error
	Close() error
	Conn() *dbr.Connection
	Logger() log.Logger
}

type dbrDb struct {
	config *Config
	conn   *dbr.Connection
	logger log.Logger
}

// New returns an instance of dbrDb that implements the Database interface.
// It is mandatory that a config and logger are passed.
// Once the clients has a Database instance, they need to call Init() and, optionally, Migrate()
// in this order before using the connection (provided by Conn()) to query the database.
func New(config *Config, logger log.Logger) (Database, error) {
	if config == nil || logger == nil {
		return nil, errors.New("config or logger is nil")
	}

	return &dbrDb{
		config: config,
		logger: logger,
	}, nil
}

// Init opens a connection to a database and Pings it N times (specified in the config) until it becomes available or until N is reached.
func (db *dbrDb) Init() error {
	conn, err := dbr.Open(db.config.Driver, db.config.url(false), nil)
	if err != nil {
		return errors.Wrapf(err, "could not open dbrDb connection with connection string: %s", db.config.url(true))
	}

	conn.SetMaxOpenConns(int(db.config.MaxOpenConnections))
	conn.SetMaxIdleConns(int(db.config.MaxIdleConnections))
	conn.SetConnMaxLifetime(db.config.ConnMaxLifetime)

	for attempt := uint(0); attempt < db.config.MaxConnectionAttempts; attempt++ {
		if err = conn.Ping(); err == nil {
			db.logger.Infoln("Successfully connected to database")
			db.conn = conn

			return nil
		}

		db.logger.Warnf("[%v/%v] Could not ping dbrDb. Retrying after %v", attempt+1, db.config.MaxConnectionAttempts, db.config.TimeBetweenAttempts)
		time.Sleep(db.config.TimeBetweenAttempts)
	}

	return errors.Wrap(err, "could not ping dbrDb")
}

// Migrate applies all up migrations from a folder (specified in the config) to the database.
func (db *dbrDb) Migrate() error {
	m, err := migrate.New(fmt.Sprintf(fileSchemaFormat, db.config.MigrationsDir), db.config.url(false))
	if err != nil {
		return errors.Wrap(err, "could not create a migrate instance")
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "could not migrate database to latest version")
	}

	db.logger.Infoln("Successfully migrated dbrDb to latest version")

	return nil
}

// Close closes the connection to the database.
func (db *dbrDb) Close() error {
	return errors.Wrap(db.conn.Close(), "could not close database connection")
}

// Conn returns the connection that can be used to query the database.
func (db *dbrDb) Conn() *dbr.Connection {
	return db.conn
}

// Logger returns the logger that is passed when initializing a new instance of dbrDb.
func (db *dbrDb) Logger() log.Logger {
	return db.logger
}
