package main

import (
	"gamejam/audio"
	"gamejam/config"
	"gamejam/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var Sound *audio.SoundManager

func init() {
	Sound = audio.NewSoundManager()

	Sound.LoadSound("walk", "sfx/walk/sfx_step_grass_l.wav")
	Sound.LoadSound("command1", "sfx/issue_command/bug_01.ogg")
	Sound.LoadSound("command2", "sfx/issue_command/bug_02.ogg")
	Sound.LoadSound("command3", "sfx/issue_command/bug_03.ogg")
	Sound.LoadSound("command4", "sfx/issue_command/bug_04.ogg")
	Sound.LoadSound("command5", "sfx/issue_command/bug_05.ogg")
	Sound.LoadSound("command6", "sfx/issue_command/bug_06.ogg")

}
func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	game := game.New(cfg, Sound)

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
