package stores

import (
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type UsersStore interface {
	Insert(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	ConfirmEmail(id string) error
}

type RestaurantsStore interface {
	Insert(restaurant *models.Restaurant) error
	GetByRating(top, skip int, forOwnerId *string, minRating, maxRating float32) ([]models.Restaurant, error)
}

type ReviewsStore interface {
	Insert(review *models.Review) error
}
