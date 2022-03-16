-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS users
(
    id           UUID                      NOT NULL,
    username     VARCHAR(255)              NOT NULL,
    password     VARCHAR(255)              NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT uix_users_username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS journals
(
    id           UUID                      NOT NULL,
    user_id      UUID                      NOT NULL,
    entry_date   VARCHAR(255)              NOT NULL,
    content      TEXT                      NULL,
    PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS journals;
