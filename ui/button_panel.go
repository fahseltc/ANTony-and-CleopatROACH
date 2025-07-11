package ui

import (
	"gamejam/eventing"
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/sim"
	"gamejam/util"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	BtnDimension = 40
	BtnPad       = 15

	PanelWidth         = 150
	PanelHeight        = 150
	PanelHorizontalPad = 10
	PanelBottomPad     = 5
)

type ButtonPanel struct {
	log *slog.Logger

	panelRect *image.Rectangle

	btns []*Button
}

func NewUnitButtonPanel(fonts *fonts.All, s *sim.T) *ButtonPanel {
	pos := image.Point{
		X: GameResolutionW - PanelWidth - PanelHorizontalPad,
		Y: GameResolutionH - PanelHeight - PanelBottomPad,
	}
	btnPanel := &ButtonPanel{
		panelRect: &image.Rectangle{
			Min: image.Pt(pos.X, pos.Y),
			Max: image.Pt(pos.X+PanelWidth, pos.Y+PanelHeight),
		},
		log: log.NewLogger().With("for", "UnitButtonPanel"),
	}

	// Buttons laid out as follows:
	// [Attack]
	// [Move] [Stop] [Hold]

	btnX := pos.X
	btnY := pos.Y

	atkBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithClickFunc(func() {
			btnPanel.log.Info("atkbtnclicked")
			s.SetActionKeyPressed(sim.AttackKeyPressed)
		}),
		WithImage(util.LoadImage("ui/btn/atk-btn.png"), util.LoadImage("ui/btn/atk-btn-pressed.png")),
		WithKeyActivation(ebiten.KeyQ),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, atkBtn)
	btnY += BtnDimension + BtnPad

	moveBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithImage(util.LoadImage("ui/btn/move-btn.png"), util.LoadImage("ui/btn/move-btn-pressed.png")),
		WithClickFunc(func() {
			btnPanel.log.Info("movebtnclicked")
			s.SetActionKeyPressed(sim.MoveKeyPressed)
		}),
		WithKeyActivation(ebiten.KeyZ),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, moveBtn)
	btnX += BtnDimension + BtnPad

	stopBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithImage(util.LoadImage("ui/btn/stop-btn.png"), util.LoadImage("ui/btn/stop-btn-pressed.png")),
		WithClickFunc(func() {
			s.SetActionKeyPressed(sim.StopKeyPressed)
		}),
		WithKeyActivation(ebiten.KeyX),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, stopBtn)
	btnX += BtnDimension + BtnPad

	holdBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithImage(util.LoadImage("ui/btn/hold-btn.png"), util.LoadImage("ui/btn/hold-btn.png")),
		WithClickFunc(func() {
			btnPanel.log.Info("holdbtnclicked")
			s.SetActionKeyPressed(sim.HoldPositionKeyPressed)
		}),
		WithKeyActivation(ebiten.KeyC),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, holdBtn)

	return btnPanel
}

func NewHiveButtonPanel(fonts *fonts.All, s *sim.T) *ButtonPanel {
	pos := image.Point{
		X: GameResolutionW - PanelWidth - PanelHorizontalPad,
		Y: GameResolutionH - PanelHeight - PanelBottomPad,
	}
	btnPanel := &ButtonPanel{
		panelRect: &image.Rectangle{
			Min: image.Pt(pos.X, pos.Y),
			Max: image.Pt(pos.X+PanelWidth, pos.Y+PanelHeight),
		},
		log: log.NewLogger().With("for", "HiveButtonPanel"),
	}

	// Buttons laid out as follows:
	// [Worker]
	// [Fighter]
	// [Upgrade]

	btnX := pos.X
	btnY := pos.Y

	workerBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithClickFunc(func() {
			btnPanel.log.Info("workerbtnclicked")
			btnPanel.log.Info("MakeAntButtonClickedEvent")
			s.EventBus.Publish(eventing.Event{
				Type: "MakeAntButtonClickedEvent",
			})
		}),
		WithImage(util.LoadImage("ui/btn/make-ant-btn-noicon.png"), util.LoadImage("ui/btn/make-ant-btn-pressed-noicon.png")),
		WithKeyActivation(ebiten.KeyQ),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, workerBtn)
	btnY += BtnDimension + BtnPad

	fighterBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithClickFunc(func() {
			btnPanel.log.Info("fighterbtnclicked")

		}),
		WithImage(util.LoadImage("ui/btn/make-ant-btn-noicon.png"), util.LoadImage("ui/btn/make-ant-btn-pressed-noicon.png")),
		//WithKeyActivation(ebiten.KeyZ),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, fighterBtn)
	btnY += BtnDimension + BtnPad

	upgradeBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithClickFunc(func() {
			btnPanel.log.Info("upgradebtnclicked")

		}),
		WithImage(util.LoadImage("ui/btn/make-ant-btn-noicon.png"), util.LoadImage("ui/btn/make-ant-btn-pressed-noicon.png")),
		//WithKeyActivation(ebiten.KeyX),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, upgradeBtn)
	btnY += BtnDimension + BtnPad

	return btnPanel
}

func NewWorkerUnitButtonPanel(fonts *fonts.All, s *sim.T) *ButtonPanel {
	btnPanel := NewUnitButtonPanel(fonts, s)

	pos := image.Point{
		X: GameResolutionW - PanelWidth - PanelHorizontalPad,
		Y: GameResolutionH - PanelHeight - PanelBottomPad,
	}

	btnX := pos.X
	btnY := pos.Y + BtnDimension + BtnPad + BtnDimension + BtnPad

	buildBtn := NewButton(fonts,
		WithRect(image.Rectangle{Min: image.Pt(btnX, btnY), Max: image.Pt(btnX+BtnDimension, btnY+BtnDimension)}),
		WithImage(util.LoadImage("TEXTURE_MISSING.png"), util.LoadImage("TEXTURE_MISSING.png")),
		WithClickFunc(func() {
			btnPanel.log.Info("buildbtnclicked")
			//s.SetActionKeyPressed(sim.MoveKeyPressed)
		}),
		//WithKeyActivation(ebiten.KeyZ),
		WithToolTip(NewTooltip(*fonts, image.Rectangle{}, LeftAlignment)),
	)
	btnPanel.btns = append(btnPanel.btns, buildBtn)

	return btnPanel
}

func (b *ButtonPanel) Update() {
	for _, btn := range b.btns {
		btn.Update()
	}
}
func (b *ButtonPanel) Draw(screen *ebiten.Image) {
	//ebitenutil.DrawRect(screen, float64(b.panelRect.Min.X), float64(b.panelRect.Min.Y), float64(b.panelRect.Dx()), float64(b.panelRect.Dy()), color.RGBA{100, 100, 100, 255})

	for _, btn := range b.btns {
		btn.Draw(screen)
	}
}
