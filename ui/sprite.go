package ui

import (
	"gamejam/util"
	"image"
	"image/color"

	"github.com/google/uuid"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var TileDimensions = 128

type SpriteType int

const (
	SpriteTypeDefault SpriteType = iota
	SpriteTypeHive
	SpriteTypeUnit
	SpriteTypeStatic
	SpriteTypeInConstruction
	// Add more sprite types as needed
)

type Sprite struct {
	Id       uuid.UUID
	Type     SpriteType
	Rect     *image.Rectangle
	img      *ebiten.Image
	angle    float64
	Selected bool

	ProgressBar *ProgressBar
}

// Units
func NewRoyalAntSprite(uuid uuid.UUID) *Sprite {
	size := 192
	return NewSprite(uuid, image.Rect(0, 0, size, size), "units/ant-royal.png", SpriteTypeUnit)
}

func NewRoyalRoachSprite(uuid uuid.UUID) *Sprite {
	size := 192
	return NewSprite(uuid, image.Rect(0, 0, size, size), "units/roach-royal.png", SpriteTypeUnit)
}
func NewDefaultAntSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions, TileDimensions), "units/ant.png", SpriteTypeUnit)
}
func NewDefaultRoachSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions, TileDimensions), "units/roach.png", SpriteTypeUnit)
}

// Buildings
func NewHiveSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions*2, TileDimensions*2), "units/anthill.png", SpriteTypeHive)
}

// Static Sprites
func NewBridgeSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions, TileDimensions), "tilemap/bridge.png", SpriteTypeStatic)
}
func NewInConstructionSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions, TileDimensions), "tilemap/in-construction.png", SpriteTypeInConstruction)
}
func NewHeartSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions/2, TileDimensions/2), "ui/heart.png", SpriteTypeStatic)
}

func NewSprite(uuid uuid.UUID, Rect image.Rectangle, imgPath string, spriteType SpriteType) *Sprite {
	scaled := util.ScaleImage(util.LoadImage(imgPath), float32(Rect.Dx()), float32(Rect.Dy()))
	return &Sprite{
		Id:          uuid,
		Type:        spriteType,
		Rect:        &Rect,
		img:         scaled,
		Selected:    false,
		ProgressBar: NewProgressBar(Rect.Min.X, Rect.Min.Y, Rect.Dx(), 6),
	}
}

func (spr *Sprite) SetPosition(pos *image.Point) {
	spr.Rect = &image.Rectangle{
		Min: *pos,
		Max: image.Point{
			X: pos.X + spr.Rect.Dx(),
			Y: pos.Y + spr.Rect.Dy(),
		},
	}
	spr.SetProgressBarPosition(pos.X, pos.Y)
}

func (spr *Sprite) SetTilePosition(x, y int) {
	spr.Rect = &image.Rectangle{
		Min: image.Point{
			X: x * TileDimensions,
			Y: y * TileDimensions,
		},
		Max: image.Point{
			X: (x * TileDimensions) + TileDimensions,
			Y: (y * TileDimensions) + TileDimensions,
		},
	}
	spr.SetProgressBarPosition(x, y)
}

func (spr *Sprite) SetProgressBarPosition(x, y int) {
	if spr.ProgressBar != nil && spr.Rect != nil {
		barX := spr.Rect.Min.X
		barY := spr.Rect.Max.Y - spr.ProgressBar.Height
		spr.ProgressBar.X = barX
		spr.ProgressBar.Y = barY
	}
}

func (spr *Sprite) SetAngle(angle float64) {
	spr.angle = angle
}
func (spr *Sprite) Draw(screen *ebiten.Image, camera *Camera) {
	opts := &ebiten.DrawImageOptions{}

	sz := spr.img.Bounds().Size()

	opts.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)

	w, h := float64(sz.X), float64(sz.X)
	opts.GeoM.Translate(-w/2*camera.ViewPortZoom, -h/2*camera.ViewPortZoom)
	opts.GeoM.Rotate(spr.angle)

	opts.GeoM.Translate(w/2*camera.ViewPortZoom, h/2*camera.ViewPortZoom)

	sprX, sprY := camera.MapPosToScreenPos(spr.Rect.Min.X, spr.Rect.Min.Y)
	opts.GeoM.Translate(float64(sprX), float64(sprY))

	screen.DrawImage(spr.img, opts)

	if spr.Selected && spr.Type != SpriteTypeStatic {
		spr.drawSelectedBox(screen, camera)
	}
	if spr.ProgressBar != nil {
		spr.ProgressBar.Draw(screen, camera)
	}

}

func (spr *Sprite) drawSelectedBox(screen *ebiten.Image, camera *Camera) {
	boxX, boxY := camera.MapPosToScreenPos(spr.Rect.Min.X, spr.Rect.Min.Y)
	w := int(float64(spr.Rect.Dx()) * camera.ViewPortZoom)
	h := int(float64(spr.Rect.Dy()) * camera.ViewPortZoom)

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
