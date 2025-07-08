package util

import (
	"gamejam/assets"
	"gamejam/vec2"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func LoadImage(filePath string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.Files, filePath)
	if err != nil {
		img, _, _ = ebitenutil.NewImageFromFileSystem(assets.Files, "TEXTURE_MISSING.png")
	}
	return img
}

func ScaleImage(input *ebiten.Image, newWidth float32, newHeight float32) *ebiten.Image {
	imgW, imgH := input.Bounds().Dx(), input.Bounds().Dy()
	wScale := newWidth / float32(imgW)
	hScale := newHeight / float32(imgH)

	scaledImg := ebiten.NewImage(int(newWidth), int(newHeight))
	op := &ebiten.DrawImageOptions{} // Draw original onto the new image with scaling
	op.GeoM.Scale(float64(wScale), float64(hScale))
	scaledImg.DrawImage(input, op)
	return scaledImg
}

func DrawCenteredText(screen *ebiten.Image, f text.Face, s string, cx, cy int, clr color.Color) {
	tw, th := text.Measure(s, f, 6)
	x := float64(cx) - tw/float64(2)
	y := float64(cy) - th/float64(2)

	var textColor color.Color
	if clr == nil {
		textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		textColor = clr
	}

	opt := text.DrawOptions{}
	opt.ColorScale.ScaleWithColor(textColor)
	opt.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, s, f, &opt)
}

func DrawCircle(screen *ebiten.Image, x, y float64, radius float64, clr color.Color) {
	img := ebiten.NewImage(int(radius*2), int(radius*2))
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				img.Set(int(dx+radius), int(dy+radius), clr)
			}
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x-radius, y-radius)
	screen.DrawImage(img, op)
}

func DrawLine(screen *ebiten.Image, x1, y1, x2, y2 float64, clr color.Color) {
	dx := x2 - x1
	dy := y2 - y1
	steps := int(max(math.Abs(dx), math.Abs(dy)))
	if steps == 0 {
		screen.Set(int(x1), int(y1), clr)
		return
	}
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		px := x1 + dx*t
		py := y1 + dy*t
		screen.Set(int(px), int(py), clr)
	}
}

func PointToVec2(point *image.Point) *vec2.T {
	return &vec2.T{
		X: float64(point.X),
		Y: float64(point.Y),
	}
}
