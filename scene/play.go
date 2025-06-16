package scene

import (
	"fmt"
	"gamejam/tilemap"
	"gamejam/ui"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PlayScene struct {
	BaseScene
	unitControls *ui.Controls
	tileMap      *tilemap.Tilemap
	camera       *ui.Camera

	sprites []*ui.Sprite
}

func NewPlayScene() *PlayScene {
	tileMap := tilemap.NewTilemap()
	scene := &PlayScene{
		unitControls: ui.NewControls(),
		tileMap:      tileMap,
		camera:       ui.NewCamera(),
	}

	ant := ui.NewSprite(image.Rect(50, 50, 128, 128), "units/ant.png")

	scene.sprites = append(scene.sprites, ant)

	return scene
}

func (s *PlayScene) Update() error {
	s.camera.Update()
	s.unitControls.Update()
	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	s.tileMap.Draw(screen, s.camera)
	s.unitControls.Draw(screen)
	for _, sprite := range s.sprites {
		sprite.Draw(screen, s.camera)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera:%v,%v", s.camera.ViewPortX, s.camera.ViewPortY), 1, 1)
}
