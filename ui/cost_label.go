package ui

import (
	"fmt"
	"gamejam/fonts"
	"gamejam/types"
	"gamejam/util"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var CostLabelResourceIconSize = float32(16)
var CostLabelWidth = 60
var CostLabelHeight = 20

type CostLabel struct {
	image *ebiten.Image
}

func NewCostLabel(resourceType types.Resource, cost int, fonts fonts.All) *CostLabel {
	img := ebiten.NewImage(CostLabelWidth, CostLabelHeight)

	// Draw background
	img.Fill(color.RGBA{230, 230, 230, 255})

	// Draw resource icon
	resIcon := utilResTypeToIcon(resourceType)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(2, 2)
	img.DrawImage(resIcon, op)

	// Draw cost number
	costText := fmt.Sprintf("%d", cost)
	x := 22
	y := 4
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	opts.ColorScale.ScaleWithColor(color.Black)
	text.Draw(img, costText, fonts.Small, opts)

	return &CostLabel{
		image: img,
	}
}

func utilResTypeToIcon(t types.Resource) *ebiten.Image {
	return util.ScaleImage(util.LoadImage(fmt.Sprintf("ui/icons/%s.png", t.ToString())), CostLabelResourceIconSize, CostLabelResourceIconSize)
}

func (c *CostLabel) Image() *ebiten.Image {
	return c.image
}
