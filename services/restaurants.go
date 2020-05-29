package services

import (
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type RestaurantsService interface {
	Create(restaurant *models.Restaurant) error
	List(top, skip int, userId string, userRole models.Role) ([]models.Restaurant, error)
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

func (rs *restaurantsService) List(top, skip int, userId string, userRole models.Role) ([]models.Restaurant, error) {
	var ownerId *string = nil

	if userRole == models.Owner {
		ownerId = &userId
	}

	restaurants, err := rs.db.Restaurants().Get(top, skip, ownerId)
	if err != nil {
		return nil, errors.Wrap(err, "could not get restaurants")
	}

	return restaurants, nil
}
