package services

import (
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type RestaurantsService interface {
	Create(restaurant *models.Restaurant) error
}

type restaurantsService struct {
	db db.Manager
}

func NewRestaurants(db db.Manager) RestaurantsService {
	return &restaurantsService{
		db: db,
	}
}

func (rs *restaurantsService) Create(restaurant *models.Restaurant) error {
	err := rs.db.Restaurants().Insert(restaurant)
	return errors.Wrap(err, "could not insert restaurant")
}
