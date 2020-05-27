package services

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

type UsersService interface {
	CreateUser(user *models.User, rawPassword *string) error
}

type usersService struct {
	db db.Manager
}

func NewUserService(dbManager db.Manager) UsersService {
	return &usersService{
		db: dbManager,
	}
}

func (us *usersService) CreateUser(user *models.User, rawPassword *string) error {
	saltedHash, err := bcrypt.GenerateFromPassword([]byte(*rawPassword), bcrypt.DefaultCost)

	// Remove the raw password from memory as fast as possible
	*rawPassword = ""

	if err != nil {
		return errors.Wrap(err, "could not generate salted hash from password")
	}

	user.HashedPassword = string(saltedHash)

	// TODO: Return custom error if email already exists
	if err = us.db.Users().Insert(user); err != nil {
		return errors.Wrap(err, "could not insert user in database")
	}

	return nil
}
