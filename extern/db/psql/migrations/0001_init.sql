-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  id  uuid  PRIMARY KEY  DEFAULT gen_random_uuid(),
  email  TEXT UNIQUE  NOT NULL  DEFAULT NULL,
  password  TEXT  NOT NULL  DEFAULT NULL,
  first_name  TEXT  NOT NULL  DEFAULT '',
  last_name  TEXT NOT NULL  DEFAULT '',
  role  TEXT  NOT NULL  DEFAULT NULL,
  created_at  TIMESTAMPTZ NOT NULL  DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL  DEFAULT NOW()
);

