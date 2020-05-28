package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/controllers"
)

func NewRouter(dbManager db.Manager, logger log.Logger, validator controllers.Validator) *mux.Router {
	usersService := services.NewUserService(dbManager)

	usersController := controllers.NewUsers(usersService, logger.WithField("module", "usersController"), validator)

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(commonMiddleware)

	apiV1Router := apiRouter.PathPrefix("/v1").Subrouter()
	apiV1Router.Methods(http.MethodPost).Path("/users").HandlerFunc(usersController.Register)
	apiV1Router.Methods(http.MethodPost).Path("/token").HandlerFunc(usersController.Login)

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
