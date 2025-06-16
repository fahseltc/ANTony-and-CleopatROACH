package main

import (
	"fmt"
	"gamejam/config"
	"gamejam/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cfg := config.New()
	game := game.NewGame(cfg)

	fmt.Printf("%#v\n", cfg)

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
