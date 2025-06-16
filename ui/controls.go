package ui

import (
	"gamejam/log"
	"gamejam/util"
	"image"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
)

type Controls struct {
	bg        *ebiten.Image
	rect      image.Rectangle
	attackBtn *Button
	//attackLabel *ebiten.Image
	moveBtn *Button
	stopBtn *Button
	log     *slog.Logger
}

func NewControls() *Controls {
	c := &Controls{
		rect: image.Rectangle{Min: image.Pt(0, 450), Max: image.Pt(300, 600)},
		//attackLabel: util.LoadImage("ui/keys/z.png"),
		log: log.NewLogger().With("for", "ui"),
	}
	c.bg = util.ScaleImage(util.LoadImage("ui/btn/controls-bg.png"), float32(c.rect.Dx()), float32(c.rect.Dy()))

	c.attackBtn = NewButton(
		WithRect(image.Rectangle{Min: image.Pt(c.rect.Min.X+20, c.rect.Min.Y+20), Max: image.Pt(c.rect.Min.X+70, c.rect.Min.Y+70)}),
		WithClickFunc(func() {
			c.log.Info("atkbtnclicked")
		}),
		WithImage(util.LoadImage("ui/btn/atk-btn.png"), util.LoadImage("ui/btn/atk-btn.png")),
	)
	c.moveBtn = NewButton(
		WithRect(image.Rectangle{Min: image.Pt(c.rect.Min.X+80, c.rect.Min.Y+20), Max: image.Pt(c.rect.Min.X+130, c.rect.Min.Y+70)}),
		WithImage(util.LoadImage("ui/btn/move-btn.png"), util.LoadImage("ui/btn/move-btn.png")),
		WithClickFunc(func() {
			c.log.Info("movebtnclicked")
		}),
	)
	c.stopBtn = NewButton(
		WithRect(image.Rectangle{Min: image.Pt(c.rect.Min.X+140, c.rect.Min.Y+20), Max: image.Pt(c.rect.Min.X+190, c.rect.Min.Y+70)}),
		WithImage(util.LoadImage("ui/btn/stop-btn.png"), util.LoadImage("ui/btn/stop-btn.png")),
		WithClickFunc(func() {
			c.log.Info("stopbtnclicked")
		}),
	)
	return c
}

func (c *Controls) Update() {
	c.attackBtn.Update()
	c.stopBtn.Update()
	c.moveBtn.Update()
}

func (c *Controls) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(c.rect.Min.X), float64(c.rect.Min.Y))
	screen.DrawImage(c.bg, opts)
	c.attackBtn.Draw(screen)
	c.stopBtn.Draw(screen)
	c.moveBtn.Draw(screen)
}
