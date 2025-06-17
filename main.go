package main

import (
	"gamejam/config"
	"gamejam/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	game := game.New(cfg)

	ebiten.SetWindowTitle(cfg.WindowTitle)
	// set external window resolution
	ebiten.SetWindowSize(cfg.Resolutions.External.Width, cfg.Resolutions.External.Height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(int(cfg.TargetFPS))

	err = ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
