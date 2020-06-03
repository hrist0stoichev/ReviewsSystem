CREATE TABLE restaurants (
    id uuid PRIMARY KEY,
    owner_id uuid REFERENCES users(id) NOT NULL,
    name VARCHAR (60) NOT NULL,
    city VARCHAR (30) NOT NULL,
    address VARCHAR (100) NOT NULL,
    img VARCHAR (200) NOT NULL,
    description VARCHAR (500) NOT NULL,
    ratings_total INTEGER NOT NULL,
    ratings_count INTEGER NOT NULL,
    average_rating REAL GENERATED ALWAYS AS (ratings_total / greatest(ratings_count, 1)) STORED
);

CREATE INDEX idx_owner_id ON restaurants (owner_id, average_rating DESC);

CREATE INDEX idx_average_rating ON restaurants (average_rating DESC);