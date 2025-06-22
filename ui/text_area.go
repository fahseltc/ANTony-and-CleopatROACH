package ui

import (
	"gamejam/fonts"
	"gamejam/util"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var LineSpacingPx = 10.0
var LineLeftPadding = 20.0

type TextArea struct {
	bg    *ebiten.Image
	fonts *fonts.All
	rect  *image.Rectangle
	text  string
	lines []string
}

func NewTextArea(fonts *fonts.All, text string) *TextArea {
	ta := &TextArea{
		bg:    util.LoadImage("ui/textbox-bg.png"),
		fonts: fonts,
		rect:  &image.Rectangle{Min: image.Pt(0, 400), Max: image.Pt(600, 600)},
		text:  text,
	}
	ta.splitTextOntoLines()
	return ta
}

func (ta *TextArea) splitTextOntoLines() {
	font := ta.fonts.Med
	// Split the text into words
	words := strings.Split(ta.text, " ") // if we use commas or something else, this will have bugs
	var currentLine string
	var totalHeight float64

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		// Measure the width of the testLine
		tw, th := text.Measure(testLine, font, LineSpacingPx)
		totalHeight += th

		if tw > float64(ta.rect.Dx()) && currentLine != "" {
			ta.lines = append(ta.lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}

	}
	if currentLine != "" {
		ta.lines = append(ta.lines, currentLine)
	}
}

func (ta *TextArea) Draw(screen *ebiten.Image) {
	// draw textbox BG
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(ta.rect.Min.X), float64(ta.rect.Min.Y))
	screen.DrawImage(ta.bg, opts)

	// draw text lines
	font := ta.fonts.Med
	_, th := text.Measure(ta.text, font, LineSpacingPx)
	th += LineSpacingPx
	y := float64(ta.rect.Min.Y) + th

	for _, line := range ta.lines {
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(ta.rect.Min.X)+LineLeftPadding, y)
		opts.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		text.Draw(screen, line, font, opts)
		y += th
	}

}

func (ta *TextArea) ChangeText(newText string) {
	ta.text = newText
	ta.splitTextOntoLines()
}
