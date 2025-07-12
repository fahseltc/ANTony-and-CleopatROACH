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

var (
	GameResolutionW = 800
	GameResolutionH = 600
)

type Ui struct {
	log   *slog.Logger
	fonts *fonts.All

	HUD     *HUD
	MiniMap *MiniMap

	Camera   *Camera
	TileMap  *tilemap.Tilemap
	eventBus *eventing.EventBus

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
		HUD:         NewHUD(fonts, sim),
		Camera:      camera,
		TileMap:     tileMap,
		eventBus:    sim.EventBus,
		MiniMap:     mm,
		DrawEnabled: true,
	}
}

func (ui *Ui) Update(sim *sim.T, selectedUnitIDs []string) {
	ui.HUD.Update(selectedUnitIDs)
	ui.Camera.Update()

	if ui.frameCounter%20 == 0 {
		ui.frameCounter = 0
		ui.MiniMap.RenderFromTilemap(ui.TileMap, sim.GetWorld().FogOfWar)
		ui.MiniMap.DrawUnits(sim.GetAllUnits(), ui.TileMap, sim.GetWorld().FogOfWar)
		ui.MiniMap.DrawViewport(ui.Camera, ui.TileMap)
	}

	ui.frameCounter++
}

func (ui *Ui) Draw(screen *ebiten.Image, sprites map[string]*Sprite) {
	if ui.DrawEnabled {
		ui.HUD.Draw(screen, sprites)
		ui.MiniMap.Draw(screen)
	}
}
