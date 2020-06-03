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
	Reviews() stores.ReviewsStore
}

type manager struct {
	users       stores.UsersStore
	restaurants stores.RestaurantsStore
	reviews     stores.ReviewsStore
}

func (m *manager) Users() stores.UsersStore {
	return m.users
}

func (m *manager) Restaurants() stores.RestaurantsStore {
	return m.restaurants
}

func (m *manager) Reviews() stores.ReviewsStore {
	return m.reviews
}

func NewManager(users stores.UsersStore, restaurants stores.RestaurantsStore, reviews stores.ReviewsStore) Manager {
	return &manager{
		users:       users,
		restaurants: restaurants,
		reviews:     reviews,
	}
}
