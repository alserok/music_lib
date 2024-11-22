package main

import (
	"github.com/alserok/music_lib/internal/app"
	"github.com/alserok/music_lib/internal/config"
)

func main() {
	app.MustStart(config.MustLoad())
}
