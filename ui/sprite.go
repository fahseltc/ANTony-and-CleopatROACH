package ui

import (
	"gamejam/eventing"
	"gamejam/util"
	"image"
	"image/color"
	"time"

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
	Id   uuid.UUID
	Type SpriteType
	Rect *image.Rectangle
	img  *ebiten.Image

	EventBus *eventing.EventBus

	Animation         *SpriteAnimation
	defaultSS         *ebiten.Image
	carryingSucroseSS *ebiten.Image
	carryingWoodSS    *ebiten.Image

	angle    float64
	Selected bool
	lastPos  image.Point

	CarryingSucrose bool
	CarryingWood    bool

	ProgressBar *ProgressBar
}

// Units
func NewRoyalAntSprite(uuid uuid.UUID) *Sprite {
	size := 192
	spr := NewSprite(uuid, image.Rect(0, 0, size, size), "units/ants/ant-royal.png", SpriteTypeUnit)
	spr.carryingSucroseSS = util.LoadImage("units/ants/ant-royal-carrying-sucrose.png")
	spr.carryingWoodSS = util.LoadImage("units/ants/ant-royal-carrying-wood.png")
	spr.defaultSS = util.LoadImage("units/ants/ant-royal-walk.png")
	spr.Animation = NewSpriteAnimation(
		spr.defaultSS,
		192, 192, 4, 4, true,
	)
	return spr
}

func NewRoyalRoachSprite(uuid uuid.UUID) *Sprite {
	size := 192
	spr := NewSprite(uuid, image.Rect(0, 0, size, size), "units/roaches/roach-royal.png", SpriteTypeUnit)
	spr.carryingSucroseSS = util.LoadImage("units/roaches/roach-royal-carrying-sucrose.png")
	spr.carryingWoodSS = util.LoadImage("units/roaches/roach-royal-carrying-wood.png")
	spr.defaultSS = util.LoadImage("units/roaches/roach-royal-walk.png")
	spr.Animation = NewSpriteAnimation(
		spr.defaultSS,
		192, 192, 4, 4, true,
	)
	return spr
}
func NewDefaultAntSprite(uuid uuid.UUID) *Sprite {
	spr := NewSprite(uuid, image.Rect(0, 0, TileDimensions, TileDimensions), "units/ants/ant.png", SpriteTypeUnit)
	spr.carryingSucroseSS = util.LoadImage("units/ants/ant-carrying-sucrose.png")
	spr.carryingWoodSS = util.LoadImage("units/ants/ant-carrying-wood.png")
	spr.defaultSS = util.LoadImage("units/ants/ant-walk.png")
	spr.Animation = NewSpriteAnimation(
		spr.defaultSS,
		128, 128, 4, 4, true,
	)

	return spr
}
func NewDefaultRoachSprite(uuid uuid.UUID) *Sprite {
	spr := NewSprite(uuid, image.Rect(0, 0, TileDimensions, TileDimensions), "units/roaches/roach.png", SpriteTypeUnit)
	spr.carryingSucroseSS = util.LoadImage("units/roaches/roach-carrying-sucrose.png")
	spr.carryingWoodSS = util.LoadImage("units/roaches/roach-carrying-wood.png")
	spr.defaultSS = util.LoadImage("units/roaches/roach-walk.png")
	spr.Animation = NewSpriteAnimation(
		spr.defaultSS,
		128, 128, 4, 4, true,
	)
	return spr
}

// Buildings
func NewHiveSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions*2, TileDimensions*2), "units/ant-hill.png", SpriteTypeHive)
}
func NewRoachHiveSprite(uuid uuid.UUID) *Sprite {
	return NewSprite(uuid, image.Rect(0, 0, TileDimensions*2, TileDimensions*2), "units/roach-hill.png", SpriteTypeHive)
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
	if spr.Rect != nil {
		spr.lastPos = spr.Rect.Min
	}
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
	if spr.Rect != nil {
		spr.lastPos = spr.Rect.Min
	}
	newX := x * TileDimensions
	newY := y * TileDimensions
	spr.Rect = &image.Rectangle{
		Min: image.Point{X: newX, Y: newY},
		Max: image.Point{X: newX + TileDimensions, Y: newY + TileDimensions},
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
	if spr.Animation != nil {
		spr.UpdateAnimation(1)
		if spr.CarryingSucrose {
			spr.Animation.SpriteSheet = spr.carryingSucroseSS
		} else if spr.CarryingWood {
			spr.Animation.SpriteSheet = spr.carryingWoodSS
		} else {
			spr.Animation.SpriteSheet = spr.defaultSS
		}
	}

	opts := &ebiten.DrawImageOptions{}

	sz := spr.img.Bounds().Size()

	opts.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)

	w, h := float64(sz.X), float64(sz.X)
	opts.GeoM.Translate(-w/2*camera.ViewPortZoom, -h/2*camera.ViewPortZoom)
	opts.GeoM.Rotate(spr.angle)

	opts.GeoM.Translate(w/2*camera.ViewPortZoom, h/2*camera.ViewPortZoom)

	sprX, sprY := camera.MapPosToScreenPos(spr.Rect.Min.X, spr.Rect.Min.Y)
	opts.GeoM.Translate(float64(sprX), float64(sprY))

	if spr.Animation != nil {
		screen.DrawImage(spr.Animation.CurrentFrameImage(), opts)
	} else {
		screen.DrawImage(spr.img, opts)
	}

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

func (spr *Sprite) UpdateAnimation(dt time.Duration) {
	if spr.Animation == nil || spr.Rect == nil {
		return
	}
	if spr.Rect.Min != spr.lastPos {
		spr.Animation.Update(dt)
		spr.lastPos = spr.Rect.Min // update for next frame
		//spr.SendWalkSFXEvent()
	}
}

// func (spr *Sprite) SendWalkSFXEvent() {
// 	if spr.EventBus != nil {
// 		spr.EventBus.Publish(eventing.Event{
// 			Type: "PlayWalkSFX",
// 		})
// 	}
// }
