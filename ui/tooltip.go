package ui

import (
	"gamejam/fonts"
	"gamejam/sim"
	"gamejam/types"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	TooltipWidth  = 150
	TooltipHeight = 150
)

type TooltipInterface interface {
	OnHover(screen *ebiten.Image)
	ReAlign(sprite *Sprite)
	ReAlignToRect(rect *image.Rectangle)
	SetAlignment(alignment Alignment)
	GetAlignment() Alignment
	GetRect() *image.Rectangle
}

type Tooltip struct {
	rect      image.Rectangle
	bg        *ebiten.Image
	fonts     fonts.All
	alignment Alignment

	ta *TextArea
}

func NewTooltip(fonts fonts.All, rect image.Rectangle, alignment Alignment, text string) *Tooltip {
	ttRect := image.Rectangle{
		Min: image.Pt(rect.Min.X, rect.Min.Y),
		Max: image.Pt(rect.Min.X+TooltipWidth, rect.Min.Y+TooltipHeight),
	}
	textArea := NewPlainTextArea(fonts.XSmall, text, &ttRect)
	scaledBg := util.ScaleImage(util.LoadImage("ui/tooltip/tooltip-bg.png"), float32(TooltipWidth), float32(TooltipHeight))

	// cl := NewCostLabel(types.ResourceTypeSucrose, 100, fonts)
	// scaledBg.DrawImage(cl.Image(), nil)

	tt := &Tooltip{
		rect:      ttRect,
		bg:        scaledBg,
		fonts:     fonts,
		alignment: alignment,

		ta: textArea,
	}

	return tt
}

func NewUnitCostToolTip(fonts fonts.All, unitType types.Unit, rect image.Rectangle, alignment Alignment) *Tooltip {
	ttRect := image.Rectangle{
		Min: image.Pt(rect.Min.X, rect.Min.Y),
		Max: image.Pt(rect.Min.X+TooltipWidth, rect.Min.Y+TooltipHeight),
	}
	scaledBg := util.ScaleImage(util.LoadImage("ui/tooltip/tooltip-bg.png"), float32(TooltipWidth), float32(TooltipHeight))

	unit := sim.GetUnitInstance(unitType, 0)

	// Sucrose Cost Label
	scl := NewCostLabel(types.ResourceTypeSucrose, int(unit.Stats.ResourceCost.Sucrose), fonts)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(10, 5)
	scaledBg.DrawImage(scl.Image(), opts)

	// Wood Cost Label
	wcl := NewCostLabel(types.ResourceTypeWood, int(unit.Stats.ResourceCost.Wood), fonts)
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(72, 5)
	scaledBg.DrawImage(wcl.Image(), opts)

	textArea := NewPlainTextArea(fonts.XSmall, unit.Stats.ToolTipString, &ttRect)

	tt := &Tooltip{
		rect:      ttRect,
		bg:        scaledBg,
		fonts:     fonts,
		alignment: alignment,
		ta:        textArea,
	}
	return tt

}

func (tt *Tooltip) OnHover(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(tt.rect.Min.X), float64(tt.rect.Min.Y))
	screen.DrawImage(tt.bg, opts)
	tt.ta.Draw(screen)
}

func (tt *Tooltip) ReAlign(sprite *Sprite) {
	tt.alignment.Align(*sprite.Rect, &tt.rect)
	tt.alignment.Align(*sprite.Rect, tt.ta.textRect)
	if tt.ta.bgRect != nil {
		tt.alignment.Align(*sprite.Rect, tt.ta.bgRect)
	}
	tt.clampOnScreen()
}

func (tt *Tooltip) ReAlignToRect(rect *image.Rectangle) {
	tt.alignment.Align(*rect, &tt.rect)
	tt.alignment.Align(*rect, tt.ta.textRect)
	if tt.ta.bgRect != nil {
		tt.alignment.Align(*rect, tt.ta.bgRect)
	}
	tt.clampOnScreen()
}

// clampOnScreen shifts the whole tooltip back onto the screen if alignment
// placed it partially or fully offscreen (e.g. a LeftAlignment tooltip on a
// button near the left edge). The offset is computed from the background rect
// (tt.rect) and applied uniformly to every sub-rect so the text and background
// stay locked together.
func (tt *Tooltip) clampOnScreen() {
	var dx, dy int

	if tt.rect.Min.X < 0 {
		dx = -tt.rect.Min.X
	} else if tt.rect.Max.X > GameResolutionW {
		dx = GameResolutionW - tt.rect.Max.X
	}

	if tt.rect.Min.Y < 0 {
		dy = -tt.rect.Min.Y
	} else if tt.rect.Max.Y > GameResolutionH {
		dy = GameResolutionH - tt.rect.Max.Y
	}

	if dx == 0 && dy == 0 {
		return
	}

	offset := image.Pt(dx, dy)
	tt.rect = tt.rect.Add(offset)
	*tt.ta.textRect = tt.ta.textRect.Add(offset)
	if tt.ta.bgRect != nil {
		*tt.ta.bgRect = tt.ta.bgRect.Add(offset)
	}
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
