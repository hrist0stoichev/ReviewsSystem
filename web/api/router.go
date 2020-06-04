package api

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
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
	emailService := services.NewEmailsService("smtp.gmail.com", "587", "", "", "", "Confirm you registration", "Click here to confirm your registration", "http://localhost:8001/api/v1/users/confirm-email", "token", "email", 30, rand.New(rand.NewSource(time.Now().UnixNano())))
	restaurantService := services.NewRestaurants(dbManager)
	reviewsService := services.NewReviews(dbManager)
	oauth2Service := services.NewOauth2(oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8001/api/v1/token/facebook",
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}, "https://graph.facebook.com/me", logger)

	usersController := controllers.NewUsers(
		usersService,
		encryptionService,
		tokensService,
		emailService,
		oauth2Service,
		logger.WithField("module", "usersController"),
		validator)

	restaurantsController := controllers.NewRestaurant(restaurantService, logger.WithField("module", "restaurantsController"), validator)
	reviewsController := controllers.NewReviews(reviewsService, restaurantService, logger.WithField("module", "reviewsController"), validator)

	authMiddleware := middlewares.NewAuth(tokensService, logger)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middlewares.SetCORS, middlewares.SetJsonContentType)

	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()
	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/users").HandlerFunc(usersController.Register)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/users/confirm-email").HandlerFunc(usersController.ConfirmEmail)
	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/token").HandlerFunc(usersController.Login)

	apiV1Router.Methods(http.MethodGet).Path("/login/facebook").HandlerFunc(usersController.RedirectToFacebookAuth)
	apiV1Router.Methods(http.MethodGet).Path("/token/facebook").HandlerFunc(usersController.HandleFacebookLoginCallback)

	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/restaurants").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Owner.String())(http.HandlerFunc(restaurantsController.Create)).ServeHTTP)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/restaurants").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String(), models.Owner.String(), models.Admin.String())(http.HandlerFunc(restaurantsController.ListByRating)).ServeHTTP)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/restaurants/{id}").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String(), models.Owner.String(), models.Admin.String())(http.HandlerFunc(restaurantsController.GetSingle)).ServeHTTP)

	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/reviews").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String())(http.HandlerFunc(reviewsController.Create)).ServeHTTP)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/reviews").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String(), models.Owner.String(), models.Admin.String())(http.HandlerFunc(reviewsController.ListForRestaurant)).ServeHTTP)
	apiV1Router.Methods(http.MethodPut, http.MethodOptions).Path("/reviews/{id}/answer").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Owner.String())(http.HandlerFunc(reviewsController.Answer)).ServeHTTP)

	return router
}
