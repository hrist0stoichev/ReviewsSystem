package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/stores/dbr"
	"github.com/hrist0stoichev/ReviewsSystem/etc"
	"github.com/hrist0stoichev/ReviewsSystem/lib/dbrdb"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/lib/server"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/controllers"
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
	restaurantsStore := dbr.NewRestaurantsStore(database.Conn().NewSession(nil))
	reviewsStore := dbr.NewReviewsStore(database.Conn().NewSession(nil))

	dbManager := db.NewManager(usersStore, restaurantsStore, reviewsStore)

	usersService := services.NewUserService(dbManager)
	tokensService := services.NewTokensService(8*time.Hour, []byte("password"))
	encryptionService := services.NewEncryptionService(services.DefaultEncryptionCost)
	emailService := services.NewEmailsService("smtp.gmail.com", "587", "", "", "", "Confirm you registration", "Click here to confirm your registration", "http://localhost:8001/api/v1/users/confirm-email", "token", "email", 30, rand.New(rand.NewSource(time.Now().UnixNano())))
	restaurantService := services.NewRestaurants(dbManager)
	reviewsService := services.NewReviews(dbManager)
	oauth2Service := services.NewOauth2(oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:9000/#",
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}, "https://graph.facebook.com/me", logger)

	usersController := controllers.NewUsers(usersService, encryptionService, tokensService, emailService, oauth2Service, logger.WithField("module", "usersController"), v)
	restaurantsController := controllers.NewRestaurant(restaurantService, logger.WithField("module", "restaurantsController"), v)
	reviewsController := controllers.NewReviews(reviewsService, restaurantService, logger.WithField("module", "reviewsController"), v)

	apiHandler := api.NewRouter(tokensService, usersController, restaurantsController, reviewsController, logger)

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
