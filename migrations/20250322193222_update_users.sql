-- +goose Up
ALTER TABLE users
ADD user_role int;

-- +goose Down
ALTER TABLE users
DROP COLUMN user_role;
