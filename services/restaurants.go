package services

import (
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type RestaurantsService interface {
	Create(restaurant *models.Restaurant) error
	ListByRating(top, skip int, userId string, userRole models.Role, minrRating, maxRating float32) ([]models.Restaurant, error)
	GetSingle(id string) (*models.Restaurant, error)
	Exists(id string) (bool, error)
}

var (
	ErrRestaurantNotFound = errors.New("restaurant not found")
)

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

func (rs *restaurantsService) ListByRating(top, skip int, userId string, userRole models.Role, minrRating, maxRating float32) ([]models.Restaurant, error) {
	var ownerId *string = nil

	if userRole == models.Owner {
		ownerId = &userId
	}

	restaurants, err := rs.db.Restaurants().GetByRating(top, skip, ownerId, minrRating, maxRating)
	if err != nil {
		return nil, errors.Wrap(err, "could not get restaurants")
	}

	return restaurants, nil
}

func (rs *restaurantsService) GetSingle(id string) (*models.Restaurant, error) {
	restaurant, err := rs.db.Restaurants().GetSingle(id)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrRestaurantNotFound
		}

		return nil, errors.Wrap(err, "could not get single restaurant")
	}

	return restaurant, nil
}

func (rs *restaurantsService) Exists(id string) (bool, error) {
	exists, err := rs.db.Restaurants().Exists(id)
	if err != nil {
		return false, errors.Wrap(err, "error checking if restaurant exists")
	}

	return exists, nil
}
