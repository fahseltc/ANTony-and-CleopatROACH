package ui

import (
	"gamejam/util"
	"image"

	"github.com/google/uuid"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Id   uuid.UUID
	rect image.Rectangle
	img  *ebiten.Image
}

func NewDefaultSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, 128, 128), "units/ant.png")
}

func NewSprite(uuid uuid.UUID, rect image.Rectangle, imgPath string) *Sprite {
	scaled := util.ScaleImage(util.LoadImage(imgPath), float32(rect.Dx()), float32(rect.Dy()))
	return &Sprite{
		Id:   uuid,
		rect: rect,
		img:  scaled,
	}
}

func (spr *Sprite) SetPosition(pos *image.Point) {
	spr.rect.Min.X = pos.X
	spr.rect.Min.Y = pos.Y
}

// DrawStatic draws the sprite at its world position, ignoring the camera (static on screen).
func (spr *Sprite) Draw(screen *ebiten.Image, camera *Camera) {
	opts := &ebiten.DrawImageOptions{}
	// Offset sprite position by adding the camera's world position
	opts.GeoM.Translate(
		float64(spr.rect.Min.X+camera.ViewPortX),
		float64(spr.rect.Min.Y+camera.ViewPortY),
	)
	screen.DrawImage(spr.img, opts)
}
