package ui

import (
	"fmt"
	"gamejam/eventing"
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/sim"
	"gamejam/tilemap"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Ui struct {
	log      *slog.Logger
	fonts    *fonts.All
	HUD      *HUD
	Camera   *Camera
	TileMap  *tilemap.Tilemap
	eventBus *eventing.EventBus

	DrawEnabled bool
}

func NewUi(fonts *fonts.All, tileMap *tilemap.Tilemap, sim *sim.T) *Ui {
	camera := NewCamera(tileMap.Width, tileMap.Height)
	return &Ui{
		log:         log.NewLogger().With("for", "ui"),
		fonts:       fonts,
		HUD:         NewHUD(fonts.Med, sim),
		Camera:      camera,
		TileMap:     tileMap,
		DrawEnabled: true,
		// textArea: NewPortraitTextArea(
		// 	fonts,
		// 	"We, ignorant of ourselves, beg often our own harms, which the wise powers deny us for our good; so find we profit by losing of our prayers.",
		// 	"portraits/ant-king.png"),
		eventBus: sim.EventBus,
	}
}

func (ui *Ui) Update() {
	ui.HUD.Update()
	ui.Camera.Update()
}

func (ui *Ui) Draw(screen *ebiten.Image) {
	if ui.DrawEnabled {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			x, y := ui.Camera.ScreenPosToMapPos(mx, my)

			clickedTile := ui.TileMap.GetTileByPosition(x, y)
			if clickedTile != nil {
				fmt.Printf("tile clicked type:%v", clickedTile.Type)
			}
		}

		ui.HUD.Draw(screen)
	}
}
