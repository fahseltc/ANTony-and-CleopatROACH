package ui

import (
	"fmt"
	"gamejam/sim"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ResourceDisplay struct {
	bg   *ebiten.Image
	font text.Face
	rect *image.Rectangle
}

func NewResourceDisplay(font text.Face) *ResourceDisplay {
	img := util.LoadImage("ui/resource-hud.png")
	rect := image.Rectangle{Min: image.Point{X: 680, Y: 0}, Max: image.Point{X: 800, Y: 80}}
	scaled := util.ScaleImage(img, float32(rect.Dx()), float32(rect.Dy()))
	return &ResourceDisplay{
		bg:   scaled,
		font: font,
		rect: &rect,
	}
}

func (rd *ResourceDisplay) Draw(screen *ebiten.Image, sim *sim.T) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(rd.rect.Min.X), float64(rd.rect.Min.Y))

	screen.DrawImage(rd.bg, opts)

	sucrose := sim.GetSucroseAmount()
	util.DrawCenteredText(screen, rd.font, fmt.Sprintf("%v", sucrose), rd.rect.Min.X+82, rd.rect.Min.Y+20, nil)

	wood := sim.GetWoodAmount()
	util.DrawCenteredText(screen, rd.font, fmt.Sprintf("%v", wood), rd.rect.Min.X+82, rd.rect.Min.Y+55, nil)
}
