package main

import (
	"gamejam/config"
	"gamejam/fonts"
	"gamejam/game"

	"github.com/hajimehoshi/ebiten/v2"
)

var fontPath = "fonts/PressStart2P-Regular.ttf"

func main() {
	cfg := config.New()
	fonts := fonts.Load(fontPath)
	game := game.New(cfg, fonts)

	ebiten.SetWindowTitle(cfg.WindowTitle)
	// set external window resolution
	ebiten.SetWindowSize(cfg.Resolutions.External.Width, cfg.Resolutions.External.Height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(int(cfg.TargetFPS))

	err := ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
