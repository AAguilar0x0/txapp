-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id  uuid  PRIMARY KEY  DEFAULT gen_random_uuid(),
  email  TEXT UNIQUE  NOT NULL  DEFAULT NULL,
  password  TEXT  NOT NULL,
  first_name  TEXT  NOT NULL,
  last_name  TEXT NOT NULL,
  role  TEXT  NOT NULL
);

