package ui

import (
	"fmt"
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/tilemap"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ui struct {
	log      *slog.Logger
	fonts    *fonts.All
	hud      *HUD
	Camera   *Camera
	TileMap  *tilemap.Tilemap
	textArea *PortraitTextArea
}

func NewUi(fonts *fonts.All, tileMap *tilemap.Tilemap) *Ui {
	camera := NewCamera(tileMap.Width, tileMap.Height)
	return &Ui{
		log:     log.NewLogger().With("for", "ui"),
		fonts:   fonts,
		hud:     NewHUD(fonts.Med),
		Camera:  camera,
		TileMap: tileMap,
		textArea: NewPortraitTextArea(
			fonts,
			"We, ignorant of ourselves, beg often our own harms, which the wise powers deny us for our good; so find we profit by losing of our prayers.",
			"portraits/ant-king.png"),
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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		x, y := ui.Camera.ScreenPosToMapPos(mx, my)

		clickedTile := ui.TileMap.GetTileByPosition(x, y)
		if clickedTile != nil {
			fmt.Printf("tile clicked type:%v", clickedTile.Type)
		}
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("zoom:%v", ui.Camera.ViewPortZoom), 1, 20)

	ui.hud.Draw(screen)
	//ui.textArea.Draw(screen)
}
