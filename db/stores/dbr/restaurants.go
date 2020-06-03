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
			INNER JOIN users min_usr ON min_rv.reviewer_id = min_usr.id
			LEFT JOIN reviews max_rv ON res.max_review_id = max_rv.id
			INNER JOIN users max_usr ON max_rv.reviewer_id = max_usr.id 
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
