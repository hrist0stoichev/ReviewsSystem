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
	GetSingle(id string) (*models.Restaurant, error)
	Exists(id string) (bool, error)
}

type ReviewsStore interface {
	GetById(revId string) (*models.Review, error)
	Update(review *models.Review) error
	Insert(review *models.Review) error
	ExistsForUserAndRestaurant(userId, restaurantId string) (bool, error)
	ListForRestaurant(restaurantId string, unanswered bool, top, skip uint64, orderBy string, isAsc bool) ([]models.Review, error)
}
