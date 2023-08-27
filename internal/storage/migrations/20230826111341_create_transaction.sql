-- +goose Up
-- +goose StatementBegin
CREATE TABLE transaction (
  id serial primary key,
  num_transaction uuid NOT NULL,
  user_id int NOT NULL,
  CREATED_AT timestamp NOT NULL DEFAULT NOW()

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction;
--SELECT 'down SQL query';
-- +goose StatementEnd
