package ui

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var MaxDuration = 120

type Notification struct {
	font            *text.Face
	textLines       []string
	maxDuration     int
	currentDuration int
	Completed       bool
}

func NewNotification(font *text.Face, text string) *Notification {
	lines := strings.Split(text, "\n")

	return &Notification{
		font:            font,
		textLines:       lines,
		maxDuration:     MaxDuration,
		currentDuration: 0,
		Completed:       false,
	}
}

func (n *Notification) Update() {
	if !n.Completed {
		n.currentDuration++
		if n.currentDuration >= n.maxDuration {
			n.Completed = true
		}
	}
}

func (n *Notification) Draw(screen *ebiten.Image) {
	if !n.Completed {
		alpha := 1.0
		halfDuration := n.maxDuration / 2
		if n.currentDuration > halfDuration {
			// Fade out in the last half
			progress := float64(n.currentDuration-halfDuration) / float64(halfDuration)
			alpha = 1.0 - progress
			if alpha < 0 {
				alpha = 0
			}
		}
		for ind, line := range n.textLines {
			tw, th := text.Measure(line, *n.font, 6)
			x := float64(400) - tw/float64(2)
			y := float64(200) - th/float64(2) + float64(ind)*25

			opts := &text.DrawOptions{}
			opts.ColorScale.ScaleAlpha(float32(alpha))
			opts.GeoM.Translate(x, y)
			text.Draw(screen, line, *n.font, opts)
		}

	}
}
