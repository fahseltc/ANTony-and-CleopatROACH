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
	// Add more sprite types as needed
)

type Sprite struct {
	Id       uuid.UUID
	rect     *image.Rectangle
	img      *ebiten.Image
	angle    float64
	Selected bool

	Type SpriteType
}

// Units
func NewRoyalAntSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, int(float64(TileDimensions)*1.5), int(float64(TileDimensions)*1.5)), "units/ant-royal.png", SpriteTypeUnit)
}
func NewRoyalRoachSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions, TileDimensions), "units/roach-royal.png", SpriteTypeUnit)
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
func NewBridgeSprite() *Sprite {
	return NewSprite(uuid.New(), image.Rect(0, 0, TileDimensions, TileDimensions), "units/bridge.png", SpriteTypeStatic)
}

func NewSprite(uuid uuid.UUID, rect image.Rectangle, imgPath string, spriteType SpriteType) *Sprite {
	scaled := util.ScaleImage(util.LoadImage(imgPath), float32(rect.Dx()), float32(rect.Dy()))
	return &Sprite{
		Id:       uuid,
		rect:     &rect,
		img:      scaled,
		Selected: false,
		Type:     spriteType,
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

func (spr *Sprite) SetTilePosition(x, y int) {
	spr.rect = &image.Rectangle{
		Min: image.Point{
			X: x * TileDimensions,
			Y: y * TileDimensions,
		},
		Max: image.Point{
			X: (x * TileDimensions) + TileDimensions,
			Y: (y * TileDimensions) + TileDimensions,
		},
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

	sprX, sprY := camera.MapPosToScreenPos(spr.rect.Min.X, spr.rect.Min.Y)
	opts.GeoM.Translate(float64(sprX), float64(sprY))

	screen.DrawImage(spr.img, opts)

	if spr.Selected && spr.Type != SpriteTypeStatic {
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
