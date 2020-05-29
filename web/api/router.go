package api

import (
	"math/rand"
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
	// TODO: Move the secrets to docker-compose secrets
	usersService := services.NewUserService(dbManager)
	tokensService := services.NewTokensService(8*time.Hour, []byte("password"))
	encryptionService := services.NewEncryptionService(services.DefaultEncryptionCost)
	emailService := services.NewEmailsService("smtp.gmail.com", "587", "hnstoychev@gmail.com", "", "", "Confirm you registration", "Click here to confirm your registration", "https://website.com/api/v1/confirm-email", "token", 30, rand.New(rand.NewSource(time.Now().UnixNano())))

	usersController := controllers.NewUsers(
		usersService,
		encryptionService,
		tokensService,
		emailService,
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
