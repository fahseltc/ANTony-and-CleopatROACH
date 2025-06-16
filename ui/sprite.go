package ui

import (
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	rect image.Rectangle
	img  *ebiten.Image
}

func NewSprite(rect image.Rectangle, imgPath string) *Sprite {
	scaled := util.ScaleImage(util.LoadImage(imgPath), float32(rect.Dx()), float32(rect.Dy()))

	return &Sprite{
		rect: rect,
		img:  scaled,
	}
}

func (spr *Sprite) Draw(screen *ebiten.Image, camera *Camera) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(spr.rect.Min.X), float64(spr.rect.Min.Y))
	screen.DrawImage(spr.img, opts)
}
