package ui

import (
	"gamejam/eventing"
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/sim"
	"gamejam/tilemap"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type Ui struct {
	log      *slog.Logger
	fonts    *fonts.All
	HUD      *HUD
	Camera   *Camera
	TileMap  *tilemap.Tilemap
	eventBus *eventing.EventBus

	MiniMap *MiniMap

	DrawEnabled  bool
	frameCounter int
}

func NewUi(fonts *fonts.All, tileMap *tilemap.Tilemap, sim *sim.T) *Ui {
	camera := NewCamera(tileMap.Width, tileMap.Height)
	mm := NewMiniMap(tileMap)
	mm.RenderFromTilemap(tileMap, sim.GetWorld().FogOfWar)
	return &Ui{
		log:         log.NewLogger().With("for", "ui"),
		fonts:       fonts,
		HUD:         NewHUD(fonts.Med, sim),
		Camera:      camera,
		TileMap:     tileMap,
		eventBus:    sim.EventBus,
		MiniMap:     mm,
		DrawEnabled: true,
	}
}

func (ui *Ui) Update(sim *sim.T) {
	ui.HUD.Update()
	ui.Camera.Update()

	if ui.frameCounter%20 == 0 {
		ui.frameCounter = 0
		ui.MiniMap.RenderFromTilemap(ui.TileMap, sim.GetWorld().FogOfWar)
		ui.MiniMap.DrawUnits(sim.GetAllUnits(), ui.TileMap, sim.GetWorld().FogOfWar)
		ui.MiniMap.DrawViewport(ui.Camera, ui.TileMap)
	}

	ui.frameCounter++
}

func (ui *Ui) Draw(units []*sim.Unit, screen *ebiten.Image) {
	if ui.DrawEnabled {
		// if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// 	mx, my := ebiten.CursorPosition()
		// 	x, y := ui.Camera.ScreenPosToMapPos(mx, my)

		// 	clickedTile := ui.TileMap.GetTileByPosition(x, y)
		// 	if clickedTile != nil {
		// 		fmt.Printf("tile clicked type:%v", clickedTile.Type)
		// 	}
		// }

		ui.HUD.Draw(screen)

		ui.MiniMap.Draw(screen)
		//ui.MiniMap.DrawViewport(screen, ui.Camera, ui.TileMap)
	}
}
