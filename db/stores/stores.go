package stores

import (
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type UsersStore interface {
	Insert(user *models.User) error
}
