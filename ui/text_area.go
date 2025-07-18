package ui

import (
	"gamejam/util"
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var LineSpacingPx = 15.0
var LineLeftPadding = 25.0

type TextArea struct {
	bg          *ebiten.Image
	currentFont text.Face
	bgRect      *image.Rectangle
	textRect    *image.Rectangle
	text        string
	lines       []string

	TextOverflows bool

	Dismissed bool
}

func NewPlainTextArea(font text.Face, text string, rect *image.Rectangle) *TextArea {
	ta := &TextArea{
		bg:            nil,
		currentFont:   font,
		bgRect:        nil,
		textRect:      rect,
		text:          text,
		TextOverflows: false,
	}
	ta.splitTextOntoLines()
	return ta
}

func NewTextArea(font text.Face, text string) *TextArea {
	bgRect := &image.Rectangle{Min: image.Pt(0, 400), Max: image.Pt(800, 600)}
	ta := &TextArea{
		bg:            util.LoadImage("ui/bg/textbox-bg.png"),
		currentFont:   font,
		bgRect:        bgRect,
		textRect:      bgRect,
		text:          text,
		TextOverflows: false,
	}
	ta.splitTextOntoLines()
	return ta
}

func (ta *TextArea) splitTextOntoLines() {
	ta.lines = nil
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
		tw, th := text.Measure(testLine, ta.currentFont, LineSpacingPx)
		totalHeight += th

		if tw+LineSpacingPx+LineLeftPadding >= float64(ta.textRect.Dx()) && currentLine != "" {
			ta.lines = append(ta.lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}

	}
	ta.TextOverflows = totalHeight >= float64(ta.textRect.Dy())
	if currentLine != "" {
		ta.lines = append(ta.lines, currentLine)
	}
}

func (ta *TextArea) Draw(screen *ebiten.Image) {
	// draw textbox BG
	if ta.bgRect != nil && ta.bg != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(ta.bgRect.Min.X), float64(ta.bgRect.Min.Y))
		screen.DrawImage(ta.bg, opts)
	}

	// draw text lines

	_, th := text.Measure(ta.text, ta.currentFont, LineSpacingPx)
	th += LineSpacingPx
	y := float64(ta.textRect.Min.Y) + 0.5*th

	//ebitenutil.DrawRect(screen, float64(ta.textRect.Min.X), float64(ta.textRect.Min.Y), float64(ta.textRect.Dx()), float64(ta.textRect.Dy()), color.RGBA{22, 22, 22, 225})
	for _, line := range ta.lines {
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(ta.textRect.Min.X)+LineLeftPadding, y)
		opts.ColorScale.ScaleWithColor(color.RGBA{R: 0, G: 0, B: 0, A: 255})
		text.Draw(screen, line, ta.currentFont, opts)
		y += th
	}
}

func (ta *TextArea) ChangeText(newText string) {
	ta.text = newText
	ta.splitTextOntoLines()
}
