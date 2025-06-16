package scene

import (
	"gamejam/environment"
	"gamejam/ui"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScene struct {
	BaseScene
	startBtn *ui.Button
	bg       *ebiten.Image
	txt      string
}

func NewMenuScene() *MenuScene {
	scene := &MenuScene{
		bg:  util.LoadImage("ui/menu-bg.png"),
		txt: "ANTony & CleopatROACH",
	}
	scene.startBtn = ui.NewButton(ui.WithText("START"), ui.WithRect(image.Rectangle{
		Min: image.Point{X: 250, Y: 520},
		Max: image.Point{X: 550, Y: 570},
	}), ui.WithClickFunc(func() {
		scene.BaseScene.sm.SwitchTo(NewPlayScene())
	}))
	return scene
}

func (s *MenuScene) Update() error {
	s.startBtn.Update()
	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.bg, nil)
	util.DrawCenteredText(screen, environment.NewFontsCollection().Large, s.txt, 400, 50, nil)

	s.startBtn.Draw(screen)
}
