package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/lib/log"
	"github.com/hrist0stoichev/ReviewsSystem/services"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/middlewares"
	"github.com/hrist0stoichev/ReviewsSystem/web/api/transfermodels"
)

type reviewsController struct {
	reviewsService     services.ReviewsService
	restaurantsService services.RestaurantsService
	baseController
}

func NewReviews(reviewsService services.ReviewsService, restaurantsService services.RestaurantsService, logger log.Logger, validator Validator) *reviewsController {
	return &reviewsController{
		reviewsService:     reviewsService,
		restaurantsService: restaurantsService,
		baseController: baseController{
			logger:    logger,
			validator: validator,
		},
	}
}

func (rs *reviewsController) Create(res http.ResponseWriter, req *http.Request) {
	reviewRequest := transfermodels.CreateReviewRequest{}
	if err := json.NewDecoder(req.Body).Decode(&reviewRequest); err != nil {
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	if err := rs.validator.Struct(reviewRequest); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	userId, err := middlewares.UserIDFromRequest(req)
	if err != nil {
		rs.logger.WithError(err).Warnln("Cannot get user id from request")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	restaurantExists, err := rs.restaurantsService.Exists(reviewRequest.RestaurantId)
	if err != nil {
		rs.logger.WithError(err).Warnln("Cannot determine whether restaurant exists")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	if !restaurantExists {
		http.NotFound(res, req)
		return
	}

	userHasRated, err := rs.reviewsService.HasUserReviewed(*userId, reviewRequest.RestaurantId)
	if err != nil {
		rs.logger.WithError(err).Warnln("Cannot determine whether user has rated")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	if userHasRated {
		http.Error(res, "You have already rated this restaurant", http.StatusConflict)
		return
	}

	review := models.Review{
		RestaurantId: reviewRequest.RestaurantId,
		ReviewerId:   *userId,
		Rating:       reviewRequest.Rating,
		Timestamp:    time.Now().UTC(),
		Comment:      reviewRequest.Comment,
	}

	if err := rs.reviewsService.Create(&review); err != nil {
		rs.logger.WithError(err).Warnln("Could not create review")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	reviewResponse := transfermodels.ReviewSimpleResponse{
		Id:        review.Id,
		Rating:    review.Rating,
		Timestamp: review.Timestamp,
		Comment:   review.Comment,
	}

	res.Header().Add("Location", fmt.Sprintf("%s%s%s/%s", req.URL.Scheme, req.Host, req.URL.Path, review.Id))
	res.WriteHeader(http.StatusCreated)

	rs.returnJsonResponse(res, reviewResponse)
}
