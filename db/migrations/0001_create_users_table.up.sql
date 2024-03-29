CREATE TYPE role AS ENUM ('regular', 'owner', 'admin');

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR (64) UNIQUE NOT NULL,
    email_confirmed boolean NOT NULL,
    email_confirmation_token VARCHAR (30),
    hashed_password VARCHAR (60) NOT NULL,
    role role NOT NULL
);

CREATE INDEX idx_email ON users USING hash (email);