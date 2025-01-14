package models

import "time"

type Song struct {
	SongID string   `json:"songID" db:"id"`
	Group  string   `json:"group" db:"group_name"`
	Song   string   `json:"song"`
	Data   SongData `json:"data"`
}

type NewSong struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type SongData struct {
	ReleaseDate time.Time `json:"releaseDate" db:"release_date"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type SongFilter struct {
	SongID      string `db:"id"`
	Group       string
	Song        string
	ReleaseDate *time.Time
	Text        string
	Link        string
	Lim         int
	Off         int
}
