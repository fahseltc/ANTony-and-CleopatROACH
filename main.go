package main

import (
	"gamejam/environment"
	"gamejam/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	env := environment.NewEnv()
	game := game.NewGame(env)

	ebiten.SetWindowTitle(env.Config.Get("windowTitle").(string))
	// set external window resolution
	ebiten.SetWindowSize(env.Config.Get("resolution.external.w").(int), env.Config.Get("resolution.external.h").(int))
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(env.Config.Get("targetFPS").(int))

	err := ebiten.RunGame(game)
	if err != nil {
		panic(err)
	}
}
