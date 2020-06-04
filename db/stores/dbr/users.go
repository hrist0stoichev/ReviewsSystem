package dbr

import (
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/db/stores"
)

const usersTable = "users"

type usersStore struct {
	session *dbr.Session
}

// NewUsersStore returns a UsersStore that uses the DBR driver
func NewUsersStore(session *dbr.Session) stores.UsersStore {
	return &usersStore{
		session: session,
	}
}

// Insert adds a new user to the database
func (us *usersStore) Insert(user *models.User) error {
	if user.Id == "" {
		user.Id = uuid.NewV4().String()
	}

	_, err := us.session.
		InsertInto(usersTable).
		Columns("email", "email_confirmed", "email_confirmation_token", "hashed_password", "role").
		Record(user).
		Exec()

	return errors.Wrap(err, "could not insert new user")
}

// GetByEmail returns a user by its email
func (us *usersStore) GetByEmail(email string) (*models.User, error) {
	user := new(models.User)
	err := us.session.
		Select("*").
		From(usersTable).
		Where("email = ?", email).
		LoadOne(user)

	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, db.ErrNotFound
		}

		return nil, errors.Wrap(err, "could not load user")
	}

	return user, nil
}

// ConfirmEmail sets the user email as confirmed and removes the confirmation token from the DB
func (us *usersStore) ConfirmEmail(id string) error {
	_, err := us.session.
		Update(usersTable).
		Set("email_confirmed", true).
		Set("email_confirmation_token", nil).
		Where("id = ?", id).
		Exec()

	return errors.Wrap(err, "could not update user fields")
}
