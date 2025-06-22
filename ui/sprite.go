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
	opts.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)
	// Offset sprite position by subtracting the camera's world position
	sprX, sprY := camera.MapPosToScreenPos(spr.rect.Min.X, spr.rect.Min.Y)
	opts.GeoM.Translate(float64(sprX), float64(sprY))

	screen.DrawImage(spr.img, opts)

	if spr.Selected {
		spr.drawSelectedBox(screen, camera)
	}
}

func (spr *Sprite) drawSelectedBox(screen *ebiten.Image, camera *Camera) {
	boxX, boxY := camera.MapPosToScreenPos(spr.rect.Min.X, spr.rect.Min.Y)
	w := int(float64(spr.rect.Dx()) * camera.ViewPortZoom)
	h := int(float64(spr.rect.Dy()) * camera.ViewPortZoom)

	green := color.RGBA{0, 255, 0, 255}

	// Top
	ebitenutil.DrawLine(screen, float64(boxX), float64(boxY), float64(boxX+w), float64(boxY), green)
	// Bottom
	ebitenutil.DrawLine(screen, float64(boxX), float64(boxY+h), float64(boxX+w), float64(boxY+h), green)
	// Left
	ebitenutil.DrawLine(screen, float64(boxX), float64(boxY), float64(boxX), float64(boxY+h), green)
	// Right
	ebitenutil.DrawLine(screen, float64(boxX+w), float64(boxY), float64(boxX+w), float64(boxY+h), green)
}
