package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

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

func (rs *reviewsController) ListForRestaurant(res http.ResponseWriter, req *http.Request) {
	restaurantId := req.URL.Query().Get("restaurantId")
	if restaurantId == "" {
		http.Error(res, "You need to specify restaurant id", http.StatusBadRequest)
		return
	}

	orderBy := req.URL.Query().Get("orderBy")
	if orderBy == "" {
		orderBy = "timestamp"
	}

	top := rs.parseFloatParam(req, "top", DefaultTop, MinTop, MaxTop)
	skip := rs.parseFloatParam(req, "skip", DefaultSkip, MinSkip, MaxSkip)
	unanswered := req.URL.Query().Get("unanswered") == "true"
	asc := req.URL.Query().Get("orderByAsc") == "true"

	reviews, err := rs.reviewsService.ListForRestaurant(restaurantId, unanswered, uint64(top), uint64(skip), orderBy, asc)
	if err != nil {
		rs.logger.WithError(err).Warnln("Cannot get reviews")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	responseReviews := make([]transfermodels.ReviewSimpleResponse, len(reviews))
	for i, r := range reviews {
		responseReviews[i] = transfermodels.ReviewSimpleResponse{
			Id:        r.Id,
			Reviewer:  r.Reviewer.Email,
			Rating:    r.Rating,
			Timestamp: r.Timestamp,
			Comment:   r.Comment,
			Answer:    r.Answer,
		}
	}

	rs.returnJsonResponse(res, responseReviews)
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

func (rs *reviewsController) Answer(res http.ResponseWriter, req *http.Request) {
	answerRequest := transfermodels.AnswerReviewRequest{}
	if err := json.NewDecoder(req.Body).Decode(&answerRequest); err != nil {
		http.Error(res, ModelDecodeError, http.StatusBadRequest)
		return
	}

	if err := rs.validator.Struct(answerRequest); err != nil {
		http.Error(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id := mux.Vars(req)["id"]

	review, err := rs.reviewsService.GetById(id)
	if err != nil {
		if err == services.ErrReviewNotFound {
			http.NotFound(res, req)
			return
		}

		rs.logger.WithError(err).Warnln("could not get review by id")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	userId, err := middlewares.UserIDFromRequest(req)
	if err != nil {
		rs.logger.WithError(err).Warnln("Cannot get user id from request")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	if review.Restaurant.OwnerId != *userId {
		http.NotFound(res, req)
		return
	}

	if review.Answer != nil {
		http.Error(res, "You have already answered this review!", http.StatusConflict)
		return
	}

	review.Answer = &answerRequest.Answer

	err = rs.reviewsService.Update(review)
	if err != nil {
		rs.logger.WithError(err).Warnln("could not update review")
		http.Error(res, InternalServerError, http.StatusInternalServerError)
		return
	}

	reviewResponse := transfermodels.ReviewSimpleResponse{
		Id:        review.Id,
		Rating:    review.Rating,
		Timestamp: review.Timestamp,
		Comment:   review.Comment,
		Answer:    review.Answer,
	}

	rs.returnJsonResponse(res, reviewResponse)
}
