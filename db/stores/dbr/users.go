package dbr

import (
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"

	"github.com/hrist0stoichev/ReviewsSystem/db"
	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/db/stores"
)

const usersTable = "users"

type usersStore struct {
	session *dbr.Session
}

func NewUsersStore(session *dbr.Session) stores.UsersStore {
	return &usersStore{
		session: session,
	}
}

func (us *usersStore) Insert(user *models.User) error {
	_, err := us.session.
		InsertInto(usersTable).
		Columns("email", "email_confirmed", "hashed_password", "role").
		Record(user).
		Exec()

	return errors.Wrap(err, "could not insert new user")
}

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
