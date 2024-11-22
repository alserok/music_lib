package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port string
	Env  string

	DB Postgres

	Clients Clients
}

type Clients struct {
	SongDataAPIAddr string
}

type Postgres struct {
	Port string
	Host string
	User string
	Pass string
	Name string
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", p.Host, p.Port, p.User, p.Name)
}

func MustLoad() *Config {
	var cfg Config

	cfg.Port = os.Getenv("PORT")
	cfg.Env = os.Getenv("ENV")

	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Pass = os.Getenv("DB_PASS")
	cfg.DB.Name = os.Getenv("DB_NAME")

	cfg.Clients.SongDataAPIAddr = os.Getenv("SONG_DATA_API_ADDR")

	return &cfg
}
