package ui

import (
	"gamejam/fonts"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PortraitTextArea struct {
	Ta          *TextArea // this should be embedded
	portrait    *ebiten.Image
	portraitPos *image.Point
}

func NewPortraitTextArea(fonts *fonts.All, text string, portraitType PortraitType) *PortraitTextArea {
	pta := &PortraitTextArea{
		Ta: NewTextArea(
			fonts.Med, text,
		),
		portrait:    util.LoadImage(portraitType.String()),
		portraitPos: &image.Point{X: 6, Y: 406},
	}
	pta.Ta.bgRect = &image.Rectangle{
		Min: image.Point{X: 0, Y: 400},
		Max: image.Point{X: 800, Y: 600},
	}
	pta.Ta.textRect = &image.Rectangle{
		Min: image.Point{X: 200, Y: 400},
		Max: image.Point{X: 800, Y: 600},
	}
	pta.Ta.splitTextOntoLines()
	pta.Ta.bg = util.LoadImage("ui/bg/textbox-bg-portrait.png")
	return pta
}

func (pta *PortraitTextArea) Draw(screen *ebiten.Image) {
	pta.Ta.Draw(screen)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(pta.portraitPos.X), float64(pta.portraitPos.Y))
	screen.DrawImage(pta.portrait, opts)
}

func (pta *PortraitTextArea) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		pta.Ta.Dismissed = true
	}
}
