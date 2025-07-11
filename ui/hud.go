package ui

import (
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/sim"
	"gamejam/util"
	"image"
	"log/slog"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	LowerUiBgPositionX = 0.0
	LowerUiBgPositionY = 400.0
)

type RightSideHUDState int

const (
	HiddenState RightSideHUDState = iota
	HiveSelectedState
	UnitSelectedState
)

type HUD struct {
	rightSideRect  image.Rectangle
	RightSideState RightSideHUDState

	hitboxes    []image.Rectangle
	miniMapRect image.Rectangle

	//	rightUnitButtonPanel       *ButtonPanel // Buttons for unit/s selected (move/attack/hold/etc)
	rightWorkerUnitButtonPanel *ButtonPanel // Buttons for worker selected (spawn more units)
	rightHiveButtonPanel       *ButtonPanel // Buttons for building selected (spawn more units)

	resourceDisplay  *ResourceDisplay
	selectedUnitArea *SelectedUnitArea

	lowerBG *ebiten.Image

	log *slog.Logger
	sim *sim.T
}

func NewHUD(fonts *fonts.All, sim *sim.T) *HUD {
	//leftSideRect := image.Rectangle{Min: image.Pt(0, 500), Max: image.Pt(200, 600)}
	rightSideRect := image.Rectangle{Min: image.Pt(600, 500), Max: image.Pt(800, 600)}
	hud := &HUD{
		rightSideRect:  rightSideRect,
		RightSideState: HiddenState,

		rightWorkerUnitButtonPanel: NewWorkerUnitButtonPanel(fonts, sim),
		rightHiveButtonPanel:       NewHiveButtonPanel(fonts, sim),

		resourceDisplay:  NewResourceDisplay(fonts.Med),
		selectedUnitArea: NewSelectedUnitArea(),

		lowerBG: util.LoadImage("ui/bg/lower-UI-bg.png"),
		log:     log.NewLogger().With("for", "HUD"),
		sim:     sim,
	}

	// Setup UI hitbox areas
	leftSideHitbox := image.Rectangle{Min: image.Pt(0, 420), Max: image.Pt(180, 600)}
	hud.hitboxes = append(hud.hitboxes, leftSideHitbox)
	rightSideHitbox := image.Rectangle{Min: image.Pt(620, 420), Max: image.Pt(800, 600)}
	hud.hitboxes = append(hud.hitboxes, rightSideHitbox)
	middleHitbox := image.Rectangle{Min: image.Pt(180, 472), Max: image.Pt(620, 600)}
	hud.hitboxes = append(hud.hitboxes, middleHitbox)

	hud.miniMapRect = image.Rectangle{Min: image.Pt(MiniMapLeftPad, 600-MiniMapHeight-MiniMapBottomPad), Max: image.Pt(MiniMapLeftPad+MiniMapWidth, 600-MiniMapBottomPad)}
	hud.hitboxes = append(hud.hitboxes, hud.miniMapRect)

	return hud
}

func (h *HUD) Update(selectedUnitIDs []string) {
	switch h.RightSideState {
	case HiddenState:
		// do nothing
	case HiveSelectedState:
		h.rightHiveButtonPanel.Update()
	case UnitSelectedState:
		h.rightWorkerUnitButtonPanel.Update()
	}
	h.selectedUnitArea.Update(selectedUnitIDs)
}

func (c *HUD) Draw(screen *ebiten.Image) {
	// Draw UI lower BG
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(LowerUiBgPositionX, LowerUiBgPositionY)
	screen.DrawImage(c.lowerBG, op)

	switch c.RightSideState {
	case HiddenState:
		// do nothing
	case HiveSelectedState:
		c.rightHiveButtonPanel.Draw(screen)
	case UnitSelectedState:
		c.rightWorkerUnitButtonPanel.Draw(screen)
	}

	c.resourceDisplay.Draw(screen, c.sim)
	c.selectedUnitArea.Draw(screen)

	// DebugDraw UI hitboxes
	// for _, hb := range c.hitboxes {
	// 	ebitenutil.DrawRect(screen, float64(hb.Min.X), float64(hb.Min.Y), float64(hb.Dx()), float64(hb.Dy()), color.RGBA{255, 255, 255, 255})
	// }
}

func (h *HUD) IsPointInside(pt image.Point) bool {
	return slices.ContainsFunc(h.hitboxes, pt.In)
}

func (h *HUD) IsPointInsideMinimap(pt image.Point) bool {
	return pt.In(h.miniMapRect)
}
