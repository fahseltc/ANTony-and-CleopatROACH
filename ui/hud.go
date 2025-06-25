package ui

import (
	"gamejam/eventing"
	"gamejam/log"
	"gamejam/sim"
	"gamejam/util"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type RightSideHUDState int

const (
	HiddenState RightSideHUDState = iota
	HiveSelectedState
	UnitSelectedState
)

type HUD struct {
	leftSideBg   *ebiten.Image
	leftSideRect image.Rectangle

	rightSideBg            *ebiten.Image
	rightSideRect          image.Rectangle
	RightSideState         RightSideHUDState
	rightSideMakeAntBtn    *Button
	rightSideMakeBridgeBtn *Button
	rightSideZImg          *ebiten.Image

	resourceDisplay *ResourceDisplay
	//attackBtn       *Button
	//attackLabel *ebiten.Image
	// moveBtn *Button
	// stopBtn *Button
	log *slog.Logger
	sim *sim.T
}

func NewHUD(font text.Face, sim *sim.T) *HUD {
	leftSideRect := image.Rectangle{Min: image.Pt(0, 500), Max: image.Pt(200, 600)}
	rightSideRect := image.Rectangle{Min: image.Pt(600, 500), Max: image.Pt(800, 600)}
	c := &HUD{
		leftSideRect: leftSideRect,
		leftSideBg:   util.ScaleImage(util.LoadImage("ui/btn/controls-bg-left.png"), float32(leftSideRect.Dx()), float32(leftSideRect.Dy())),

		rightSideRect:  rightSideRect,
		rightSideBg:    util.ScaleImage(util.LoadImage("ui/btn/controls-bg-right.png"), float32(leftSideRect.Dx()), float32(leftSideRect.Dy())),
		RightSideState: HiddenState,
		rightSideZImg:  util.ScaleImage(util.LoadImage("ui/keys/z.png"), float32(40), float32(40)),

		resourceDisplay: NewResourceDisplay(font),
		log:             log.NewLogger().With("for", "HUD"),
		sim:             sim,
	}

	c.rightSideMakeAntBtn = NewButton(font,
		WithRect(image.Rectangle{
			Min: image.Pt(c.rightSideRect.Min.X+20, c.rightSideRect.Min.Y+15),
			Max: image.Pt(c.rightSideRect.Min.X+70, c.rightSideRect.Min.Y+65)}),
		WithClickFunc(func() {
			c.log.Info("MakeAntButtonClickedEvent")
			sim.EventBus.Publish(eventing.Event{
				Type: "MakeAntButtonClickedEvent",
			})
		}),
		WithImage(util.LoadImage("ui/btn/make-ant-btn.png"), util.LoadImage("ui/btn/make-ant-btn-pressed.png")),
		WithKeyActivation(ebiten.KeyZ),
	)

	c.rightSideMakeBridgeBtn = NewButton(font,
		WithRect(image.Rectangle{
			Min: image.Pt(c.rightSideRect.Min.X+20, c.rightSideRect.Min.Y+15),
			Max: image.Pt(c.rightSideRect.Min.X+70, c.rightSideRect.Min.Y+65)}),
		WithClickFunc(func() {
			c.log.Info("MakeBridgeButtonClickedEvent")
			sim.EventBus.Publish(eventing.Event{
				Type: "MakeBridgeButtonClickedEvent",
			})
		}),
		WithImage(util.LoadImage("ui/btn/make-bridge-btn.png"), util.LoadImage("ui/btn/make-bridge-btn-pressed.png")),
		WithKeyActivation(ebiten.KeyZ),
	)

	// c.attackBtn = NewButton(font,
	// 	WithRect(image.Rectangle{Min: image.Pt(c.rect.Min.X+20, c.rect.Min.Y+20), Max: image.Pt(c.rect.Min.X+70, c.rect.Min.Y+70)}),
	// 	WithClickFunc(func() {
	// 		c.log.Info("atkbtnclicked")
	// 	}),
	// 	WithImage(util.LoadImage("ui/btn/atk-btn.png"), util.LoadImage("ui/btn/atk-btn-pressed.png")),
	// 	WithKeyActivation(ebiten.KeyZ),
	// )
	// c.moveBtn = NewButton(font,
	// 	WithRect(image.Rectangle{Min: image.Pt(c.rect.Min.X+80, c.rect.Min.Y+20), Max: image.Pt(c.rect.Min.X+130, c.rect.Min.Y+70)}),
	// 	WithImage(util.LoadImage("ui/btn/move-btn.png"), util.LoadImage("ui/btn/move-btn-pressed.png")),
	// 	WithClickFunc(func() {
	// 		c.log.Info("movebtnclicked")
	// 	}),
	// 	WithKeyActivation(ebiten.KeyX),
	// )
	// c.stopBtn = NewButton(font,
	// 	WithRect(image.Rectangle{Min: image.Pt(c.rect.Min.X+140, c.rect.Min.Y+20), Max: image.Pt(c.rect.Min.X+190, c.rect.Min.Y+70)}),
	// 	WithImage(util.LoadImage("ui/btn/stop-btn.png"), util.LoadImage("ui/btn/stop-btn-pressed.png")),
	// 	WithClickFunc(func() {
	// 		c.log.Info("stopbtnclicked")
	// 	}),
	// 	WithKeyActivation(ebiten.KeyC),
	// )
	return c
}

func (c *HUD) Update() {
	switch c.RightSideState {
	case HiddenState:
		// do nothing
	case HiveSelectedState:
		c.rightSideMakeAntBtn.Update()
	case UnitSelectedState:
		c.rightSideMakeBridgeBtn.Update()
	}

	//c.attackBtn.Update()
	//c.stopBtn.Update()
	//c.moveBtn.Update()

}

func (c *HUD) Draw(screen *ebiten.Image) {
	// draw left side BG
	// opts := &ebiten.DrawImageOptions{}
	// opts.GeoM.Translate(float64(c.leftSideRect.Min.X), float64(c.leftSideRect.Min.Y))
	// screen.DrawImage(c.leftSideBg, opts)
	//c.attackBtn.Draw(screen)
	//c.stopBtn.Draw(screen)
	//c.moveBtn.Draw(screen)

	c.DrawRightSide(screen)
	// draw resource display
	c.resourceDisplay.Draw(screen, c.sim)
}

func (c *HUD) DrawRightSide(screen *ebiten.Image) {
	// setup right side BG options
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(c.rightSideRect.Min.X), float64(c.rightSideRect.Min.Y))

	switch c.RightSideState {
	case HiddenState:
		// draw nothing
	case HiveSelectedState:
		screen.DrawImage(c.rightSideBg, opts)
		c.rightSideMakeAntBtn.Draw(screen)
		c.DrawRightSideZImg(screen)
	case UnitSelectedState:
		screen.DrawImage(c.rightSideBg, opts)
		c.rightSideMakeBridgeBtn.Draw(screen)
		c.DrawRightSideZImg(screen)

	}

}
func (c *HUD) DrawRightSideZImg(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(c.rightSideRect.Min.X+25), float64(c.rightSideRect.Min.Y+64))
	screen.DrawImage(c.rightSideZImg, opts)
}
