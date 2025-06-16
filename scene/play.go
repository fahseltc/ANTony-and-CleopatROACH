package scene

import (
	"gamejam/tilemap"
	"gamejam/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type PlayScene struct {
	BaseScene
	controls *ui.Controls
	tileMap  *tilemap.Tilemap

	viewPointX float64
	viewPointY float64
}

func NewPlayScene() *PlayScene {
	scene := &PlayScene{
		controls: ui.NewControls(),
		tileMap:  tilemap.NewTilemap(),
	}

	return scene
}

func (s *PlayScene) Update() error {
	scrollSpeed := 8.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		s.viewPointY += scrollSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		s.viewPointX += scrollSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		s.viewPointY -= scrollSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		s.viewPointX -= scrollSpeed
	}

	s.controls.Update()
	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	s.tileMap.Draw(screen, int(s.viewPointX), int(s.viewPointY))
	s.controls.Draw(screen)
}
