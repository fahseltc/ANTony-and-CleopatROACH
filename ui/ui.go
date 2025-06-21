package ui

import (
	"fmt"
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/tilemap"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Ui struct {
	log     *slog.Logger
	fonts   *fonts.All
	hud     *HUD
	Camera  *Camera
	TileMap *tilemap.Tilemap
}

func NewUi(fonts *fonts.All, tileMap *tilemap.Tilemap) *Ui {
	return &Ui{
		log:     log.NewLogger().With("for", "ui"),
		fonts:   fonts,
		hud:     NewHUD(fonts.Med),
		Camera:  NewCamera(),
		TileMap: tileMap,
	}
}

func (ui *Ui) Update() {
	ui.hud.Update()
	ui.Camera.Update()
}

func (ui *Ui) Draw(screen *ebiten.Image) {
	// opts := &ebiten.DrawImageOptions{}

	// opts.GeoM.Scale(ui.Camera.ViewPortZoom, ui.Camera.ViewPortZoom)
	// opts.GeoM.Translate(float64(ui.Camera.ViewPortX), float64(ui.Camera.ViewPortY))
	//	ui.tileMap.Render(screen)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("zoom:%v", ui.Camera.ViewPortZoom), 1, 20)

	ui.hud.Draw(screen)
}
