package dbr

import (
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/db/stores"
)

const (
	restaurantsTable = "restaurants"
	id               = "id"
	ownerId          = "owner_id"
	name             = "name"
	city             = "city"
	address          = "address"
	img              = "img"
	description      = "description"
	ratingsTotal     = "ratings_total"
	ratingsCount     = "ratings_count"
	averageRating    = "average_rating"
	minReviewId      = "min_review_id"
	maxReviewId      = "max_review_id"
)

type restaurantsStore struct {
	session *dbr.Session
}

func NewRestaurantsStore(session *dbr.Session) stores.RestaurantsStore {
	return &restaurantsStore{
		session: session,
	}
}

func (rs *restaurantsStore) Insert(restaurant *models.Restaurant) error {
	if restaurant.Id == "" {
		restaurant.Id = uuid.NewV4().String()
	}

	_, err := rs.session.
		InsertInto(restaurantsTable).
		Columns(id, ownerId, name, city, address, img, description, ratingsTotal, ratingsCount, minReviewId, maxReviewId).
		Record(restaurant).
		Exec()

	return errors.Wrap(err, "could not insert into restaurants table")
}

func (rs *restaurantsStore) GetByRating(top, skip int, forOwnerId *string, minRating, maxRating float32) ([]models.Restaurant, error) {
	query := rs.session.
		Select(id, name, city, address, img, description, averageRating).
		From(restaurantsTable).
		Where(fmt.Sprintf("%s >= ? AND %s <= ?", averageRating, averageRating), minRating, maxRating).
		OrderDesc(averageRating).
		Offset(uint64(skip)).
		Limit(uint64(top))

	if forOwnerId != nil {
		query = query.Where(fmt.Sprintf("%s = ?", ownerId), forOwnerId)
	}

	restaurants := make([]models.Restaurant, 0, top)

	_, err := query.Load(&restaurants)
	if err != nil {
		return nil, errors.Wrap(err, "could not get restaurants from db")
	}

	return restaurants, nil
}
