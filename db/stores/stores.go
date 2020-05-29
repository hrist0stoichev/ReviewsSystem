package stores

import (
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type UsersStore interface {
	Insert(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	ConfirmEmail(id string) error
}
