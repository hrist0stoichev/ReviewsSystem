package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/controllers"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/middlewares"
)

func NewRouter(dbManager db.Manager, logger log.Logger, validator controllers.Validator) *mux.Router {
	usersService := services.NewUserService(dbManager)
	// TODO: Move this password to docker-compose secrets
	tokensService := services.NewTokensService(8*time.Hour, []byte("password"))
	encryptionService := services.NewEncryptionService(services.DefaultEncryptionCost)

	usersController := controllers.NewUsers(
		usersService,
		encryptionService,
		tokensService,
		logger.WithField("module", "usersController"),
		validator)

	_ = middlewares.NewAuth(tokensService)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middlewares.SetJsonContentType)

	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()
	apiV1Router.Methods(http.MethodPost).Path("/users").HandlerFunc(usersController.Register)
	apiV1Router.Methods(http.MethodPost).Path("/token").HandlerFunc(usersController.Login)

	return router
}
