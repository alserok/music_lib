package main

import (
	"github.com/alserok/music_lib/internal/app"
	"github.com/alserok/music_lib/internal/config"
)

// @title Music library API
// @version 1.0
// @BasePath /v1
// @host      localhost:5000
func main() {
	app.MustStart(config.MustLoad())
}
