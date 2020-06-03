DROP INDEX idx_restaurant_id;

ALTER TABLE restaurants
    DROP COLUMN min_review_id,
    DROP COLUMN max_review_id;

DROP TABLE reviews;
