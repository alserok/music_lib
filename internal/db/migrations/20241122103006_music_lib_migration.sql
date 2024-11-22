-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs
(
    song_id      text PRIMARY KEY,
    group_name   VARCHAR(255) NOT NULL,
    song         VARCHAR(255) NOT NULL,
    release_date VARCHAR(10),
    text         TEXT,
    link         VARCHAR(255)
);

CREATE INDEX song_id_index ON songs(song_id);
CREATE INDEX song_index ON songs(song);
CREATE INDEX group_index ON songs(group_name);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
