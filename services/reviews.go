package services

import (
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type ReviewsService interface {
	Create(review *models.Review) error
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