package services

import (
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type ReviewsService interface {
	Create(review *models.Review) error
	HasUserReviewed(userId, restaurantId string) (bool, error)
	ListForRestaurant(restaurantId string, unanswered bool, top, skip uint64, orderBy string, isAsc bool) ([]models.Review, error)
}

type reviewsService struct {
	db db.Manager
}

func NewReviews(db db.Manager) ReviewsService {
	return &reviewsService{
		db: db,
	}
}

func (rs *reviewsService) Create(review *models.Review) error {
	err := rs.db.Reviews().Insert(review)
	return errors.Wrap(err, "could not insert review")
}

func (rs *reviewsService) HasUserReviewed(userId, restaurantId string) (bool, error) {
	exists, err := rs.db.Reviews().ExistsForUserAndRestaurant(userId, restaurantId)
	return exists, errors.Wrap(err, "could not determine whether review exists")
}

func (rs *reviewsService) ListForRestaurant(restaurantId string, unanswered bool, top, skip uint64, orderBy string, isAsc bool) ([]models.Review, error) {
	reviews, err := rs.db.Reviews().ListForRestaurant(restaurantId, unanswered, top, skip, orderBy, isAsc)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get reviews for restaurant")
	}

	return reviews, nil
}
