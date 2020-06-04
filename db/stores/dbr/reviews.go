package dbr

import (
	"fmt"

	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/hrist0stoichev/ReviewsSystem/db"
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

// NewReviewsStore returns ReviewsStore that uses the DBR driver
func NewReviewsStore(session *dbr.Session) stores.ReviewsStore {
	return &reviewsStore{
		session: session,
	}
}

// GetById returns a review by its id or a ErrNotFound if it doesn't exist
func (rs *reviewsStore) GetById(revId string) (*models.Review, error) {
	review := models.Review{
		Restaurant: &models.Restaurant{},
	}

	err := rs.session.
		Select("*").
		From(reviewsTable).
		Join(restaurantsTable, fmt.Sprintf("%s.%s = %s.%s", reviewsTable, restaurantId, restaurantsTable, id)).
		Where(fmt.Sprintf("%s.%s = ?", reviewsTable, reviewId), revId).
		LoadOne(&review)

	if err != nil {
		if err == dbr.ErrNotFound {
			return nil, db.ErrNotFound
		}

		return nil, errors.Wrap(err, "could not query for review")
	}

	return &review, nil
}

// Update updates the rating, comment, and answer of a given review by its id.
func (rs *reviewsStore) Update(review *models.Review) error {
	_, err := rs.session.
		Update(reviewsTable).
		Set(rating, review.Rating).
		Set(comment, review.Comment).
		Set(answer, review.Answer).
		Where(fmt.Sprintf("%s = ?", reviewId), review.Id).
		Exec()

	return errors.Wrap(err, "could not update review")
}

// Insert starts a new transaction and makes the following changes:
// 1. Inserts the review in the reviews table
// 2. Swaps the restaurant.min_review with the current review in case it has worse score
// 3. Swaps the restaurant.max_review with the current review in case it has better score
// 4. Updates the restaurant.ratings_total and restaurant.ratings_count so that restaurant.average_rating is automatically updated by the DB
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

// ExistsForUserAndRestaurant checks whether a a particular user has already written a review for a particular restaurant
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

// ListForRestaurant returns reviews for a particular restaurantId by applying filters (pagination, orderBy, only unanswered reviews)
func (rs *reviewsStore) ListForRestaurant(restaurantId string, unanswered bool, top, skip uint64, orderBy string, isAsc bool) ([]models.Review, error) {
	query := rs.session.
		Select("reviews.id, reviews.rating, reviews.timestamp, reviews.comment, reviews.answer, users.email").
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

	rows, err := query.Rows()
	if err != nil {
		return nil, errors.Wrap(err, "could not query for reviews")
	}

	reviews := make([]models.Review, 0, top)
	for rows.Next() {
		r := models.Review{
			Reviewer: &models.User{},
		}

		err = rows.Scan(&r.Id, &r.Rating, &r.Timestamp, &r.Comment, &r.Answer, &r.Reviewer.Email)
		if err != nil {
			return nil, errors.Wrap(err, "cannot scan row")
		}

		reviews = append(reviews, r)
	}

	return reviews, nil
}
