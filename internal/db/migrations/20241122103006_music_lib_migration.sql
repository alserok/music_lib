-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs
(
    id           text PRIMARY KEY,
    song         VARCHAR(255) NOT NULL,
    release_date TIMESTAMP,
    text         TEXT,
    link         VARCHAR(255)
);

CREATE TABLE group_songs
(
    group_name VARCHAR(255) NOT NULL,
    song_id    text REFERENCES songs (id)
);

CREATE INDEX group_songs_group_index ON group_songs (group_name);
CREATE INDEX group_songs_song_index ON group_songs (song_id);
CREATE INDEX songs_song_index ON songs (song);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
