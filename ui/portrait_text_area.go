package ui

import (
	"gamejam/fonts"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type PortraitTextArea struct {
	ta          *TextArea
	portrait    *ebiten.Image
	portraitPos *image.Point
}

func NewPortraitTextArea(fonts *fonts.All, text string, portraitPath string) *PortraitTextArea {
	pta := &PortraitTextArea{
		ta: NewTextArea(
			fonts, text,
		),
		portrait:    util.LoadImage(portraitPath),
		portraitPos: &image.Point{X: 6, Y: 406},
	}
	pta.ta.bgRect = &image.Rectangle{
		Min: image.Point{X: 0, Y: 400},
		Max: image.Point{X: 800, Y: 600},
	}
	pta.ta.textRect = &image.Rectangle{
		Min: image.Point{X: 200, Y: 400},
		Max: image.Point{X: 800, Y: 600},
	}
	pta.ta.splitTextOntoLines()
	pta.ta.bg = util.LoadImage("ui/textbox-bg-portrait.png")
	return pta
}

func (pta *PortraitTextArea) Draw(screen *ebiten.Image) {
	pta.ta.Draw(screen)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(pta.portraitPos.X), float64(pta.portraitPos.Y))
	screen.DrawImage(pta.portrait, opts)
}
