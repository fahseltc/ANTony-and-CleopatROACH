package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type ProgressBar struct {
	X, Y          int
	Width, Height int
	Progress      float64 // 0.0 to 1.0
	BgColor       color.Color
	FgColor       color.Color
	BorderColor   color.Color
}

func NewProgressBar(x, y, width, height int) *ProgressBar {
	return &ProgressBar{
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Progress:    0.0,
		BgColor:     color.RGBA{64, 64, 64, 255},
		FgColor:     color.RGBA{155, 155, 155, 255},
		BorderColor: color.RGBA{255, 255, 255, 255},
	}
}
func NewHealthBar(x, y, width, height int) *ProgressBar {
	return &ProgressBar{
		X:           x,
		Y:           y + 25,
		Width:       width,
		Height:      height,
		Progress:    0.0,
		BgColor:     color.RGBA{64, 64, 64, 255},
		FgColor:     color.RGBA{0, 200, 0, 255},
		BorderColor: color.RGBA{255, 255, 255, 255},
	}
}

func (pb *ProgressBar) SetProgress(p float64) {
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	pb.Progress = p
}

func (pb *ProgressBar) Draw(screen *ebiten.Image, camera *Camera) {
	// if pb.Progress == 0 {
	// 	return
	// }
	// Map world position to screen position using camera
	screenX, screenY := camera.MapPosToScreenPos(pb.X, pb.Y)

	// Draw background
	bg := ebiten.NewImage(pb.Width, pb.Height)
	//bg.Fill(pb.BgColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)
	op.GeoM.Translate(float64(screenX), float64(screenY))
	screen.DrawImage(bg, op)

	// Draw progress
	fgWidth := int(float64(pb.Width) * pb.Progress)
	if fgWidth > 0 {
		fg := ebiten.NewImage(fgWidth, pb.Height)
		fg.Fill(pb.FgColor)
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)
		op2.GeoM.Translate(float64(screenX), float64(screenY))
		screen.DrawImage(fg, op2)
	}
}
