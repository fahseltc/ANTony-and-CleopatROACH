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
	log      *slog.Logger
	fonts    *fonts.All
	hud      *HUD
	Camera   *Camera
	TileMap  *tilemap.Tilemap
	textArea *TextArea
}

func NewUi(cam *Camera, fonts *fonts.All, tileMap *tilemap.Tilemap) *Ui {
	return &Ui{
		log:      log.NewLogger().With("for", "ui"),
		fonts:    fonts,
		hud:      NewHUD(fonts.Med),
		Camera:   cam,
		TileMap:  tileMap,
		textArea: NewTextArea(fonts, "test, this is a big chunk of text in the area!"),
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

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("zoom: %0.2f", ui.Camera.Zoom), 1, 20)

	ui.hud.Draw(screen)
	ui.textArea.Draw(screen)
}
