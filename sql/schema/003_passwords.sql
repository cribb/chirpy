-- +goose Up
ALTER TABLE users
ADD COLUMN hashed_password NOT NULL

INSERT CREATE TABLE chirps
 (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    body TEXT NOT NULL,

    FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
 );


-- +goose Down
DROP TABLE IF EXISTS chirps;

