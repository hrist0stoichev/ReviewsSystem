package dbr

import (
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hrist0stoichev/ReviewsSystem/db/models"
	"github.com/hrist0stoichev/ReviewsSystem/db/stores"
)

const (
	reviewsTable = "reviews"
	reviewId     = "id"
	restaurantId = "restaurant_id"
	reviewerId   = "reviewer_id"
	rating       = "rating"
	timestamp    = "timestamp"
	comment      = "comment"
	answer       = "answer"
)

type reviewsStore struct {
	session *dbr.Session
}

func NewReviewsStore(session *dbr.Session) stores.ReviewsStore {
	return &reviewsStore{
		session: session,
	}
}

func (rs *reviewsStore) Insert(review *models.Review) error {
	if review.Id == "" {
		review.Id = uuid.NewV4().String()
	}

	tx, err := rs.session.Begin()
	if err != nil {
		return errors.Wrap(err, "could not begin transaction")
	}

	defer tx.RollbackUnlessCommitted()

	_, err = tx.
		InsertInto(reviewsTable).
		Columns(reviewId, restaurantId, reviewerId, rating, timestamp, comment, answer).
		Record(review).
		Exec()
	if err != nil {
		return errors.Wrap(err, "could not insert review")
	}

	_, err = tx.UpdateBySql(`
		UPDATE restaurants rst 
		LEFT JOIN reviews rv 
		ON rst.min_review_id = rv.id 
		SET rst.min_review_id = ? 
		WHERE rst.id = ? AND rv.rating >= ?`,
		review.Id, review.RestaurantId, review.Rating).Exec()
	if err != nil {
		return errors.Wrap(err, "could not update min review")
	}

	_, err = tx.UpdateBySql(`
		UPDATE restaurants rst 
		LEFT JOIN reviews rv 
		ON rst.max_review_id = rv.id 
		SET rst.max_review_id = ? 
		WHERE rst.id = ? AND rv.rating <= ?`,
		review.Id, review.RestaurantId, review.Rating).Exec()
	if err != nil {
		return errors.Wrap(err, "could not update max review")
	}

	return errors.Wrap(tx.Commit(), "could not commit transaction")
}
