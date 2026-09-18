package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var MaxDuration = 120

// NotificationMaxWidth is the maximum pixel width a notification line may occupy
// before it is wrapped onto the next line. The notification is centered on the
// 800px-wide screen, so this leaves a margin on both sides.
var NotificationMaxWidth = 700

type Notification struct {
	font            *text.Face
	textLines       []string
	maxDuration     int
	currentDuration int
	Completed       bool
}

func NewNotification(font *text.Face, msg string) *Notification {
	// wrapText also splits on explicit newlines, so a message with hard breaks
	// (e.g. the bridge "not enough resources" message) still breaks where
	// intended, while long unbroken lines are word-wrapped to fit on screen.
	lines := wrapText(msg, *font, NotificationMaxWidth)

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
