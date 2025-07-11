package ui

import (
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/sim"
	"gamejam/util"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	LowerUiBgPositionX = 0.0
	LowerUiBgPositionY = 400.0

	RightSideButtonX       = 440.0
	RightSideButtonY       = 515.0
	RightSideButtonPadding = 20.0
)

// image.Rectangle{Min: image.Pt(600, 500), Max: image.Pt(800, 600)}
//Min: image.Pt(c.rightSideRect.Min.X+20, c.rightSideRect.Min.Y+15),
//Max: image.Pt(c.rightSideRect.Min.X+70, c.rightSideRect.Min.Y+65)}),

type RightSideHUDState int

const (
	HiddenState RightSideHUDState = iota
	HiveSelectedState
	UnitSelectedState
)

type HUD struct {
	rightSideRect  image.Rectangle
	RightSideState RightSideHUDState

	//	rightUnitButtonPanel       *ButtonPanel // Buttons for unit/s selected (move/attack/hold/etc)
	rightWorkerUnitButtonPanel *ButtonPanel // Buttons for worker selected (spawn more units)
	rightHiveButtonPanel       *ButtonPanel // Buttons for building selected (spawn more units)

	resourceDisplay *ResourceDisplay

	lowerBG *ebiten.Image

	log *slog.Logger
	sim *sim.T
}

func NewHUD(fonts *fonts.All, sim *sim.T) *HUD {
	//leftSideRect := image.Rectangle{Min: image.Pt(0, 500), Max: image.Pt(200, 600)}
	rightSideRect := image.Rectangle{Min: image.Pt(600, 500), Max: image.Pt(800, 600)}
	c := &HUD{
		rightSideRect:  rightSideRect,
		RightSideState: HiddenState,

		//rightUnitButtonPanel: NewUnitButtonPanel(fonts, sim),
		rightWorkerUnitButtonPanel: NewWorkerUnitButtonPanel(fonts, sim),
		rightHiveButtonPanel:       NewHiveButtonPanel(fonts, sim),

		resourceDisplay: NewResourceDisplay(fonts.Med),

		lowerBG: util.LoadImage("ui/bg/lower-UI-bg.png"),
		log:     log.NewLogger().With("for", "HUD"),
		sim:     sim,
	}

	return c
}

func (c *HUD) Update() {
	switch c.RightSideState {
	case HiddenState:
		// do nothing
	case HiveSelectedState:
		c.rightHiveButtonPanel.Update()
	case UnitSelectedState:
		c.rightWorkerUnitButtonPanel.Update()
	}

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
}
