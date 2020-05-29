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

	restaurant := models.Restaurant{
		OwnerId: middlewares.UserIDFromRequest(req),
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
