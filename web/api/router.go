package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/controllers"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/middlewares"
)

func NewRouter(
	tokensService services.TokensService,
	usersController *controllers.Users,
	restaurantsController *controllers.Restaurants,
	reviewsController *controllers.Reviews,
	logger log.Logger,
) *mux.Router {
	authMiddleware := middlewares.NewAuth(tokensService, logger)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()
	apiV1Router.Use(middlewares.SetCORS, middlewares.SetJsonContentType)

	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/users").HandlerFunc(usersController.Register)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/users/confirm-email").HandlerFunc(usersController.ConfirmEmail)
	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/token").HandlerFunc(usersController.Login)

	apiV1Router.Methods(http.MethodGet).Path("/facebookauth").HandlerFunc(usersController.RedirectToFacebookAuth)
	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/token/facebook").HandlerFunc(usersController.FacebookLogin)

	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/restaurants").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Owner.String())(http.HandlerFunc(restaurantsController.Create)).ServeHTTP)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/restaurants").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String(), models.Owner.String(), models.Admin.String())(http.HandlerFunc(restaurantsController.ListByRating)).ServeHTTP)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/restaurants/{id}").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String(), models.Owner.String(), models.Admin.String())(http.HandlerFunc(restaurantsController.GetSingle)).ServeHTTP)

	apiV1Router.Methods(http.MethodPost, http.MethodOptions).Path("/reviews").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String())(http.HandlerFunc(reviewsController.Create)).ServeHTTP)
	apiV1Router.Methods(http.MethodGet, http.MethodOptions).Path("/reviews").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Regular.String(), models.Owner.String(), models.Admin.String())(http.HandlerFunc(reviewsController.ListForRestaurant)).ServeHTTP)
	apiV1Router.Methods(http.MethodPut, http.MethodOptions).Path("/reviews/{id}/answer").HandlerFunc(authMiddleware.AuthorizeForRoles(models.Owner.String())(http.HandlerFunc(reviewsController.Answer)).ServeHTTP)

	return router
}
