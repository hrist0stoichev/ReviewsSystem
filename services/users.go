package services

import (
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UsersService interface {
	CreateUser(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	ConfirmEmail(id string) error
}

type usersService struct {
	db db.Manager
}

func NewUserService(dbManager db.Manager) UsersService {
	return &usersService{
		db: dbManager,
	}
}

func (us *usersService) CreateUser(user *models.User) error {
	err := us.db.Users().Insert(user)
	return errors.Wrap(err, "could not insert user in the database")
}

func (us *usersService) GetByEmail(email string) (*models.User, error) {
	user, err := us.db.Users().GetByEmail(email)
	if err != nil {
		if err == db.ErrNotFound {
			return nil, ErrUserNotFound
		}

		return nil, errors.Wrap(err, "could not get user by email")
	}

	return user, nil
}

func (us *usersService) ConfirmEmail(id string) error {
	err := us.db.Users().ConfirmEmail(id)
	return errors.Wrap(err, "could not confirm email")
}
