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
	reviewsService services.ReviewsService
	baseController
}

func NewReviews(reviewsService services.ReviewsService, logger log.Logger, validator Validator) *reviewsController {
	return &reviewsController{
		reviewsService: reviewsService,
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

	// TODO: Check if the restaurant exists

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
