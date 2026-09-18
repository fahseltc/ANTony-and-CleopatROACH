package ui

import (
	"fmt"
	"gamejam/fonts"
	"gamejam/util"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Slider struct {
	Type                string
	X, Y, Width, Height int
	HandleX             int
	Dragging            bool
	Volume              float64 // 0.0 - 1.0
	font                text.Face
}

func NewSlider(sliderType string, x, y int, font *fonts.All, startingVolume float64) *Slider {
	// Clamp starting volume between 0 and 1
	if startingVolume < 0 {
		startingVolume = 0
	}
	if startingVolume > 1 {
		startingVolume = 1
	}

	width := 300
	handleX := x + int(startingVolume*float64(width))

	return &Slider{
		Type:    sliderType,
		X:       x,
		Y:       y,
		Width:   width,
		Height:  40,
		HandleX: handleX,
		Volume:  startingVolume,
		font:    font.Small,
	}
}

func (s *Slider) Update() {
	mouseX, mouseY := ebiten.CursorPosition()

	// Check if mouse pressed on handle or inside slider bar
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !s.Dragging {
			// Start drag if cursor is over handle
			handleRect := s.handleRect()
			if pointInRect(mouseX, mouseY, handleRect) {
				s.Dragging = true
			}
		}
	} else {
		// Release drag on mouse up
		s.Dragging = false
	}

	if s.Dragging {
		// Clamp handleX within slider bounds
		minX := s.X
		maxX := s.X + s.Width
		if mouseX < minX {
			mouseX = minX
		}
		if mouseX > maxX {
			mouseX = maxX
		}
		s.HandleX = mouseX

		// Calculate volume (0.0 to 1.0)
		s.Volume = float64(s.HandleX-s.X) / float64(s.Width)
	}
}

func (s *Slider) Draw(screen *ebiten.Image) {
	// Draw slider bar background (gray)
	barRect := s.barRect()
	ebitenutil.DrawRect(screen, float64(barRect.Min.X), float64(barRect.Min.Y), float64(s.Width), float64(s.Height), color.RGBA{100, 100, 100, 255})

	// Draw slider handle (white)
	handleRect := s.handleRect()
	ebitenutil.DrawRect(screen, float64(handleRect.Min.X), float64(handleRect.Min.Y), float64(handleRect.Dx()), float64(handleRect.Dy()), color.RGBA{255, 255, 255, 255})

	// Draw volume text
	util.DrawCenteredText(screen, s.font, fmt.Sprintf("%v Volume: %.2f", s.Type, s.Volume), int(float64(s.X)+float64(s.barRect().Dx())*0.5), s.Y-20, color.RGBA{0, 0, 0, 255})
}

func (s *Slider) barRect() image.Rectangle {
	return image.Rect(s.X, s.Y, s.X+s.Width, s.Y+s.Height)
}

func (s *Slider) handleRect() image.Rectangle {
	handleWidth := 30
	handleHeight := s.Height + 10
	return image.Rect(s.HandleX-handleWidth/2, s.Y-(handleHeight-s.Height)/2, s.HandleX+handleWidth/2, s.Y+(handleHeight+s.Height)/2)
}

func pointInRect(x, y int, rect image.Rectangle) bool {
	return x >= rect.Min.X && x <= rect.Max.X && y >= rect.Min.Y && y <= rect.Max.Y
}
