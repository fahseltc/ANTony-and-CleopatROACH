package scene

import (
	"gamejam/environment"
	"gamejam/util"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayScene struct {
	BaseScene
}

func NewPlayScene() *PlayScene {
	scene := &PlayScene{}

	return scene
}

func (s *PlayScene) Update() error {

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{122, 122, 122, 255})
	util.DrawCenteredText(screen, environment.NewFontsCollection().Med, "play state", 300, 300, nil)

}
