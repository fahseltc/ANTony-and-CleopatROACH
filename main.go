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

	//	Sound.LoadSound("walk", "sfx/walk/sfx_step_grass_l.wav")
	Sound.LoadSound("sfx_command_0", "sfx/issue_command/bug_01.wav")
	Sound.LoadSound("sfx_command_1", "sfx/issue_command/bug_02.wav")
	Sound.LoadSound("sfx_command_2", "sfx/issue_command/bug_03.wav")
	Sound.LoadSound("sfx_command_3", "sfx/issue_command/bug_04.wav")
	Sound.LoadSound("sfx_command_4", "sfx/issue_command/bug_05.wav")

	Sound.LoadSound("msx_gamesong1", "music/Sketchbook 2024-11-07.wav")
	Sound.LoadSound("msx_menusong", "music/Sketchbook 2024-01-24_02.wav")
	Sound.LoadSound("msx_narratorsong", "music/JDSherbert Desert Sirocco.wav")

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
