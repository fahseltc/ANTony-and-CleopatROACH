package scene

import (
	"fmt"
	"gamejam/fonts"
	"gamejam/sim"
	"gamejam/tilemap"
	"gamejam/ui"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PlayScene struct {
	BaseScene
	sim *sim.T
	ui  *ui.Ui

	tileMap *tilemap.Tilemap

	fonts *fonts.All

	sprites []*ui.Sprite
}

func NewPlayScene(fonts *fonts.All) *PlayScene {
	tileMap := tilemap.NewTilemap()
	scene := &PlayScene{
		fonts:   fonts,
		sim:     sim.New(60),
		ui:      ui.NewUi(fonts, tileMap),
		tileMap: tileMap,
	}

	sim.AddUnit()
	ant := ui.NewSprite(image.Rect(50, 50, 128, 128), "units/ant.png")
	scene.sprites = append(scene.sprites, ant)
	return scene
}

func (s *PlayScene) Update() error {
	s.ui.Update()
	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	s.ui.Draw(screen)
	for _, sprite := range s.sprites {
		sprite.Draw(screen, s.ui.Camera)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera:%v,%v", s.ui.Camera.ViewPortX, s.ui.Camera.ViewPortY), 1, 1)
}
