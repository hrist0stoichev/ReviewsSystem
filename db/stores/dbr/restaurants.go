package dbr

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hrist0stoichev/ReviewsSystem/db"
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

// NewRestaurantsStore returns a RestaurantsStore that uses the DBR driver
func NewRestaurantsStore(session *dbr.Session) stores.RestaurantsStore {
	return &restaurantsStore{
		session: session,
	}
}

// Insert generates a new ID for the restaurant and inserts it in the database. The ID can then be used by the callers of this method
// in case an error is not returned.
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

// GetByRating returns a list of restaurants ordered by average rating, applying a number of filters (pagination, rating range, specific owner)
// There is an index on the averageRating column so that this query executes faster.
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

// Exists check if a restaurant with a given ID exists.
func (rs *restaurantsStore) Exists(restId string) (bool, error) {
	idFoo := ""

	err := rs.session.
		Select(id).
		From(restaurantsTable).
		Where("id = ?", restId).
		LoadOne(&idFoo)

	if err != nil {
		if err == dbr.ErrNotFound {
			return false, nil
		}

		return false, errors.Wrap(err, "cannot execute query")
	}

	return true, nil
}

// GetSingle returns a restaurant by id, populating its min_review and max_review fields. This operation is extremely optimized
// as the min_review_id and max_review_id are stored within the restaurant record and updated only when new reviews are added to the
// restaurant. This allows for getting the restaurant and its worst and best reviews using a single query without searching in the
// reviews table every time.
func (rs *restaurantsStore) GetSingle(resId string) (*models.Restaurant, error) {
	r := struct {
		Id                 string
		OwnerId            string
		Name               string
		City               string
		Address            string
		Img                string
		Description        string
		AverageRating      float32
		MinReviewId        *string
		MinReviewRating    *uint8
		MinReviewTimestamp *time.Time
		MinReviewComment   *string
		MinReviewAnswer    *string
		MinReviewReviewer  *string
		MaxReviewId        *string
		MaxReviewRating    *uint8
		MaxReviewTimestamp *time.Time
		MaxReviewComment   *string
		MaxReviewAnswer    *string
		MaxReviewReviewer  *string
	}{}

	// Get the restaurant with its min and max reviews
	err := rs.session.QueryRow(`
			SELECT res.id, res.owner_id, res.name, res.city, res.address, res.img, res.description, res.average_rating, min_rv.id, min_rv.rating, min_rv.timestamp, min_rv.comment, min_rv.answer, min_usr.email, max_rv.id, max_rv.rating, max_rv.timestamp, max_rv.comment, max_rv.answer, max_usr.email
			FROM restaurants res
			LEFT JOIN reviews min_rv ON res.min_review_id = min_rv.id
			LEFT JOIN users min_usr ON min_rv.reviewer_id = min_usr.id
			LEFT JOIN reviews max_rv ON res.max_review_id = max_rv.id
			LEFT JOIN users max_usr ON max_rv.reviewer_id = max_usr.id 
			WHERE res.id = $1`, resId).
		Scan(&r.Id, &r.OwnerId, &r.Name, &r.City, &r.Address, &r.Img, &r.Description, &r.AverageRating, &r.MinReviewId, &r.MinReviewRating, &r.MinReviewTimestamp, &r.MinReviewComment, &r.MinReviewAnswer, &r.MinReviewReviewer, &r.MaxReviewId, &r.MaxReviewRating, &r.MaxReviewTimestamp, &r.MaxReviewComment, &r.MaxReviewAnswer, &r.MaxReviewReviewer)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNotFound
		}

		return nil, errors.Wrap(err, "could not scan restaurant row")
	}

	restaurant := models.Restaurant{
		Id:            r.Id,
		OwnerId:       r.OwnerId,
		MinReviewId:   r.MinReviewId,
		MaxReviewId:   r.MaxReviewId,
		Name:          r.Name,
		City:          r.City,
		Address:       r.Address,
		Img:           r.Img,
		Description:   r.Description,
		AverageRating: r.AverageRating,
	}

	if r.MinReviewId != nil {
		restaurant.MinReview = &models.Review{
			Id:           *r.MinReviewId,
			RestaurantId: r.Id,
			Rating:       *r.MinReviewRating,
			Timestamp:    *r.MinReviewTimestamp,
			Comment:      *r.MinReviewComment,
			Answer:       r.MinReviewAnswer,
			Reviewer: &models.User{
				Email: *r.MinReviewReviewer,
			},
		}
	}

	if r.MaxReviewId != nil {
		restaurant.MaxReview = &models.Review{
			Id:           *r.MaxReviewId,
			RestaurantId: r.Id,
			Rating:       *r.MaxReviewRating,
			Timestamp:    *r.MaxReviewTimestamp,
			Comment:      *r.MaxReviewComment,
			Answer:       r.MaxReviewAnswer,
			Reviewer: &models.User{
				Email: *r.MaxReviewReviewer,
			},
		}
	}

	return &restaurant, nil
}
