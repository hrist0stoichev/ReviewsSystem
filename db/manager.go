package db

import (
	"errors"

	"github.com/hrist0stoichev/ReviewsSystem/db/stores"
)

var (
	ErrNotFound = errors.New("not found")
)

type Manager interface {
	Users() stores.UsersStore
}

type manager struct {
	users stores.UsersStore
}

func (m *manager) Users() stores.UsersStore {
	return m.users
}

func NewManager(users stores.UsersStore) Manager {
	return &manager{
		users: users,
	}
}
