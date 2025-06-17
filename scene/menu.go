package scene

import (
	"gamejam/ui"
	"gamejam/util"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScene struct {
	BaseScene
	startBtn *ui.Button
	txt      string
	font     text.Face
}

func NewMenuScene(font text.Face) *MenuScene {
	scene := &MenuScene{
		txt:  "Anteo and Antiet",
		font: font,
	}
	scene.startBtn = ui.NewButton(font, ui.WithText("START"), ui.WithRect(image.Rectangle{
		Min: image.Point{X: 300, Y: 300},
		Max: image.Point{X: 500, Y: 350},
	}), ui.WithClickFunc(func() {
		scene.BaseScene.sm.SwitchTo(NewPlayScene(scene.font))
	}))
	return scene
}

func (s *MenuScene) Update() error {
	s.startBtn.Update()
	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{22, 0, 0, 255}) // Fill Red
	util.DrawCenteredText(screen, s.font, s.txt, 400, 50, nil)

	s.startBtn.Draw(screen)
}
