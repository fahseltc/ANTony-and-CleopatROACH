package ui

import (
	"image/color"
	"strings"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type FullscreenText struct {
	TextLines    []string
	ScrollY      float64
	ScrollSpeed  float64
	FontFace     text.Face
	Done         bool
	screenHeight int
	lineHeight   int
	LineSpacing  float64

	PaddingLeft int
}

var (
	ScreenWidth  = 800
	ScreenHeight = 600
	HPadding     = 60
	ScrollSpeed  = 1.5
)

func NewFullscreenText(font text.Face, rawText string, lineSpacing float64) *FullscreenText {
	maxWidth := ScreenWidth - 2*HPadding
	lines := wrapText(rawText, font, maxWidth)

	_, th := text.Measure("A", font, 1.0)
	lineHeight := int(th)

	return &FullscreenText{
		TextLines:    lines,
		ScrollY:      float64(ScreenHeight), // start offscreen bottom
		ScrollSpeed:  ScrollSpeed,
		FontFace:     font,
		screenHeight: ScreenHeight,
		lineHeight:   lineHeight,
		LineSpacing:  lineSpacing,
		PaddingLeft:  HPadding,
	}
}

func (f *FullscreenText) Update() {
	if f.Done {
		return
	}

	minScroll := float64(f.screenHeight) - float64(f.TotalTextHeight())

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if f.ScrollY > minScroll {
			f.ScrollY = minScroll // jump immediately to bottom
		} else {
			f.Done = true // finish on next click
		}
		return
	}

	// scroll up
	f.ScrollY -= f.ScrollSpeed

	// clamp scroll between top and bottom positions
	if f.ScrollY < minScroll {
		f.ScrollY = minScroll
	}
	if f.ScrollY > float64(f.screenHeight) {
		f.ScrollY = float64(f.screenHeight)
	}
}

func (f *FullscreenText) Draw(screen *ebiten.Image) {
	if f.Done {
		return
	}

	y := int(f.ScrollY)
	for _, line := range f.TextLines {
		opts := &text.DrawOptions{}
		opts.GeoM.Translate(float64(f.PaddingLeft), float64(y))
		opts.ColorScale.ScaleWithColor(color.Black)
		text.Draw(screen, line, f.FontFace, opts)
		y += int(float64(f.lineHeight) * f.LineSpacing)
	}
}

func (f *FullscreenText) TotalTextHeight() int {
	return int(float64(f.lineHeight) * f.LineSpacing * float64(len(f.TextLines)))
}

func (f *FullscreenText) IsDone() bool {
	return f.Done
}

func wrapText(textStr string, face text.Face, maxWidth int) []string {
	var lines []string
	paragraphs := strings.Split(textStr, "\n")

	for _, para := range paragraphs {
		words := strings.FieldsFunc(para, unicode.IsSpace)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		line := words[0]
		for _, word := range words[1:] {
			testLine := line + " " + word
			tw, _ := text.Measure(testLine, face, 1.0)
			if tw > float64(maxWidth) {
				lines = append(lines, line)
				line = word
			} else {
				line = testLine
			}
		}
		lines = append(lines, line)
	}

	return lines
}
