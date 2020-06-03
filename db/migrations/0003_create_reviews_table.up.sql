CREATE TABLE reviews (
    id uuid PRIMARY KEY,
    restaurant_id uuid REFERENCES restaurants (id) NOT NULL,
    reviewer_id uuid REFERENCES users (id) NOT NULL,
    rating SMALLINT NOT NULL CHECK (rating > 0 AND rating < 6),
    timestamp timestamp NOT NULL,
    comment VARCHAR (300) NOT NULL,
    answer VARCHAR (300),
    UNIQUE (restaurant_id, reviewer_id)
);

CREATE INDEX idx_restaurant_id ON reviews (restaurant_id, timestamp DESC);

ALTER TABLE restaurants
ADD COLUMN min_review_id uuid REFERENCES reviews (id),
ADD COLUMN max_review_id uuid REFERENCES reviews (id);