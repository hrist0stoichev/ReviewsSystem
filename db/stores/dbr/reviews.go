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
		SET min_review_id = ?
		FROM reviews rv
		WHERE rst.id = ? AND (rst.min_review_id is NULL OR (rv.id = rst.min_review_id AND rv.rating >= ?))`,
		review.Id, review.RestaurantId, review.Rating).Exec()
	if err != nil {
		return errors.Wrap(err, "could not update min review")
	}

	_, err = tx.UpdateBySql(`
		UPDATE restaurants rst
		SET max_review_id = ?
		FROM reviews rv
		WHERE rst.id = ? AND (rst.max_review_id is NULL OR (rv.id = rst.max_review_id AND rv.rating <= ?))`,
		review.Id, review.RestaurantId, review.Rating).Exec()
	if err != nil {
		return errors.Wrap(err, "could not update max review")
	}

	_, err = tx.UpdateBySql(`
		UPDATE restaurants
		SET ratings_total = ratings_total + ?,
		ratings_count = ratings_count + 1
		WHERE id = ?`,
		review.Rating, review.RestaurantId).Exec()
	if err != nil {
		return errors.Wrap(err, "could not update rating statistics for restaurants")
	}

	return errors.Wrap(tx.Commit(), "could not commit transaction")
}

func (rs *reviewsStore) ExistsForUserAndRestaurant(userId, restaurantId string) (bool, error) {
	idFoo := ""

	err := rs.session.
		Select(reviewId).
		From(reviewsTable).
		Where("restaurant_id = ? AND reviewer_id = ?", restaurantId, userId).
		LoadOne(&idFoo)

	if err != nil {
		if err == dbr.ErrNotFound {
			return false, nil
		}

		return false, errors.Wrap(err, "cannot execute query")
	}

	return true, nil
}

func (rs *reviewsStore) ListForRestaurant(restaurantId string, unanswered bool, top, skip uint64, orderBy string, isAsc bool) ([]models.Review, error) {
	reviews := make([]models.Review, 0, top)

	query := rs.session.
		Select("*").
		From(reviewsTable).
		Join(usersTable, "reviews.reviewer_id = users.id").
		Where("restaurant_id = ?", restaurantId).
		OrderDir(orderBy, isAsc).
		Limit(top).
		Offset(skip)

	if unanswered {
		query = query.
			Where("answer is NULL")
	}

	_, err := query.Load(&reviews)
	if err != nil {
		return nil, errors.Wrap(err, "could not load reviews")
	}

	return reviews, nil
}
