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
	Restaurants() stores.RestaurantsStore
}

type manager struct {
	users       stores.UsersStore
	restaurants stores.RestaurantsStore
}

func (m *manager) Users() stores.UsersStore {
	return m.users
}

func (m *manager) Restaurants() stores.RestaurantsStore {
	return m.restaurants
}

func NewManager(users stores.UsersStore, restaurants stores.RestaurantsStore) Manager {
	return &manager{
		users:       users,
		restaurants: restaurants,
	}
}
