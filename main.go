package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/stores/dbr"
	"github.com/hrist0stoichev/ReviewsSystem/etc"
	"github.com/hrist0stoichev/ReviewsSystem/lib/dbrdb"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/lib/server"
	"github.com/hrist0stoichev/ReviewsSystem/web/api"
)

func main() {
	cfg, err := etc.GetConfig()
	if err != nil {
		fmt.Printf("could not get config: %v", err)
		os.Exit(1)
	}

	v := validator.New()
	if err = v.Struct(cfg); err != nil {
		fmt.Printf("config not valid: %v", err)
		os.Exit(1)
	}

	logger, err := log.NewLogrus(&cfg.Logging)
	if err != nil {
		fmt.Printf("could not create new logger: %v", err)
		os.Exit(1)
	}

	// Until this point there is no logger configured, so print to the console instead
	database, err := connectToDatabase(&cfg.Database, logger)
	if err != nil {
		logger.WithError(err).Fatalln("could not connect to database")
	}

	usersStore := dbr.NewUsersStore(database.Conn().NewSession(nil))

	dbManager := db.NewManager(usersStore)

	apiHandler := api.NewRouter(dbManager, logger, v)
	apiServer, err := server.New(&cfg.Server, apiHandler, logger)
	if err != nil {
		logger.WithError(err).Fatalln("could not create server")
	}

	go apiServer.ListenAndServe()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	<-sigint
	apiServer.Shutdown(context.Background())
}

func connectToDatabase(cfg *dbrdb.Config, logger log.Logger) (dbrdb.Database, error) {
	db, err := dbrdb.New(cfg, logger.WithField("module", "database"))
	if err != nil {
		return nil, errors.Wrap(err, "could not create a new db instance")
	}

	if err = db.Init(); err != nil {
		return nil, errors.Wrap(err, "could not initialize db")
	}

	if err = db.Migrate(); err != nil {
		return nil, errors.Wrap(err, "could not migrate db")
	}

	return db, nil
}
