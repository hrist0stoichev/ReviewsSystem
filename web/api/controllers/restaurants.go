package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/middlewares"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/transfermodels"
)

type restaurantsController struct {
	restaurantsService services.RestaurantsService
	baseController
}

func NewRestaurant(restaurantsService services.RestaurantsService, logger log.Logger, validator Validator) *restaurantsController {
	return &restaurantsController{
		restaurantsService: restaurantsService,
		baseController: baseController{
			logger:    logger,
			validator: validator,
		},
	}
}

func (rs *restaurantsController) Get(res http.ResponseWriter, req *http.Request) {
	top := rs.parseIntParam(req, "top", DefaultTop, MinTop, MaxTop)
	skip := rs.parseIntParam(req, "skip", DefaultSkip, MinSkip, MaxSkip)

	userId, idErr := middlewares.UserIDFromRequest(req)
	userRole, roleErr := middlewares.UserRoleFromRequest(req)

	if idErr != nil || roleErr != nil {
		rs.logger.WithError(idErr).WithError(roleErr).Warnln("Cannot get user id or role from the request")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	restaurants, err := rs.restaurantsService.List(top, skip, *userId, *userRole)
	if err != nil {
		rs.logger.WithError(err).Warnln("could not list restaurants")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	restaurantsResponse := make([]transfermodels.RestaurantSimpleResponse, len(restaurants))
	for i, r := range restaurants {
		restaurantsResponse[i] = transfermodels.RestaurantSimpleResponse{
			Id:            r.Id,
			Name:          r.Name,
			City:          r.City,
			Address:       r.Address,
			AverageRating: r.AverageRating,
		}
	}

	rs.returnJsonResponse(res, restaurantsResponse)
}

func (rs *restaurantsController) Create(res http.ResponseWriter, req *http.Request) {
	restaurantRequest := transfermodels.CreateRestaurantRequest{}
	if err := json.NewDecoder(req.Body).Decode(&restaurantRequest); err != nil {
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	if err := rs.validator.Struct(restaurantRequest); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	userId, err := middlewares.UserIDFromRequest(req)
	if err != nil {
		rs.logger.WithError(err).Warnln("Cannot get user id from request")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	restaurant := models.Restaurant{
		OwnerId: *userId,
		Name:    restaurantRequest.Name,
		City:    restaurantRequest.City,
		Address: restaurantRequest.Address,
	}

	if err := rs.restaurantsService.Create(&restaurant); err != nil {
		rs.logger.WithError(err).Warnln("Could not create restaurant")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	restaurantResponse := transfermodels.RestaurantSimpleResponse{
		Id:            restaurant.Id,
		Name:          restaurant.Name,
		City:          restaurant.City,
		Address:       restaurant.Address,
		AverageRating: 0,
	}

	res.Header().Add("Location", fmt.Sprintf("%s%s%s/%s", req.URL.Scheme, req.Host, req.URL.Path, restaurant.Id))
	res.WriteHeader(http.StatusCreated)

	rs.returnJsonResponse(res, restaurantResponse)
}
