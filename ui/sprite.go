package ui

import (
	"gamejam/util"
	"image"
	"image/color"

	"github.com/google/uuid"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Id       uuid.UUID
	rect     *image.Rectangle
	img      *ebiten.Image
	Selected bool
}

func NewDefaultSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, 128, 128), "units/ant.png")
}

func NewSprite(uuid uuid.UUID, rect image.Rectangle, imgPath string) *Sprite {
	scaled := util.ScaleImage(util.LoadImage(imgPath), float32(rect.Dx()), float32(rect.Dy()))
	return &Sprite{
		Id:       uuid,
		rect:     &rect,
		img:      scaled,
		Selected: false,
	}
}

func (spr *Sprite) SetPosition(pos *image.Point) {
	spr.rect = &image.Rectangle{
		Min: *pos,
		Max: image.Point{
			X: pos.X + spr.rect.Dx(),
			Y: pos.Y + spr.rect.Dy(),
		},
	}
}

// DrawStatic draws the sprite at its world position, ignoring the camera (static on screen).
func (spr *Sprite) Draw(screen *ebiten.Image, camera *Camera) {
	opts := &ebiten.DrawImageOptions{}

	// Apply camera zoom scaling
	opts.GeoM.Scale(camera.Zoom, camera.Zoom)

	// TODO this does not seem to be translating properly between a sprite's
	// Map position and the translated screen position
	spriteX, spriteY := float64(spr.rect.Min.X), float64(spr.rect.Min.Y)
	screenX, screenY := camera.MapPosToScreenPos(spriteX, spriteY)
	// Offset sprite position by subtracting the camera's world position
	opts.GeoM.Translate(screenX, screenY)
	screen.DrawImage(spr.img, opts)

	if spr.Selected {
		x := int(float64(spr.rect.Min.X) + camera.X*camera.Zoom)
		y := int(float64(spr.rect.Min.Y) + camera.Y*camera.Zoom)
		w := int(float64(spr.rect.Dx()) * camera.Zoom)
		h := int(float64(spr.rect.Dy()) * camera.Zoom)

		green := color.RGBA{0, 255, 0, 255}

		// Top
		ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x+w), float64(y), green)
		// Bottom
		ebitenutil.DrawLine(screen, float64(x), float64(y+h), float64(x+w), float64(y+h), green)
		// Left
		ebitenutil.DrawLine(screen, float64(x), float64(y), float64(x), float64(y+h), green)
		// Right
		ebitenutil.DrawLine(screen, float64(x+w), float64(y), float64(x+w), float64(y+h), green)
	}
}
