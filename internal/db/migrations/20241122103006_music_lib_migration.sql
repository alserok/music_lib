-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE songs
(
    song_id      text PRIMARY KEY,
    group_name   VARCHAR(255) NOT NULL,
    song         VARCHAR(255) NOT NULL,
    release_date VARCHAR(10),
    text         TEXT,
    link         VARCHAR(255)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
