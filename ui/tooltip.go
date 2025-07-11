package ui

import (
	"gamejam/fonts"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	TooltipWidth  = 100
	TooltipHeight = 100
)

type TooltipInterface interface {
	OnHover(screen *ebiten.Image)
	ReAlign(sprite *Sprite)
	SetAlignment(alignment Alignment)
	GetAlignment() Alignment
	GetRect() *image.Rectangle
}

type Tooltip struct {
	rect      image.Rectangle
	bg        *ebiten.Image
	fonts     fonts.All
	alignment Alignment
}

func NewTooltip(fonts fonts.All, rect image.Rectangle, alignment Alignment) *Tooltip {
	ttRect := &image.Rectangle{
		Min: image.Pt(rect.Min.X, rect.Min.Y),
		Max: image.Pt(rect.Min.X+TooltipWidth, rect.Min.Y+TooltipHeight),
	}
	scaledBg := util.ScaleImage(util.LoadImage("ui/tooltip/tooltip-bg.png"), float32(TooltipWidth), float32(TooltipHeight))
	tt := &Tooltip{
		rect:      *ttRect,
		bg:        scaledBg,
		fonts:     fonts,
		alignment: alignment,
	}

	return tt
}

func (tt *Tooltip) OnHover(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(tt.rect.Min.X), float64(tt.rect.Min.Y))
	screen.DrawImage(tt.bg, opts)
}

func (tt *Tooltip) ReAlign(sprite *Sprite) {
	tt.alignment.Align(*sprite.Rect, &tt.rect)
}

func (tt *Tooltip) SetAlignment(alignment Alignment) {
	tt.alignment = alignment
}
func (tt *Tooltip) GetAlignment() Alignment {
	return tt.alignment
}
func (tt *Tooltip) GetRect() *image.Rectangle {
	return &tt.rect
}
