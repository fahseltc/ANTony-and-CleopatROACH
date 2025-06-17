package util

import (
	"gamejam/assets"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func LoadImage(filePath string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.Files, filePath)
	if err != nil {
		return nil
	}
	if img == nil {
		img, _, _ = ebitenutil.NewImageFromFileSystem(assets.Files, "TEXTURE_MISSING.png")
	}
	return img
}

func ScaleImage(input *ebiten.Image, newWidth float32, newHeight float32) *ebiten.Image {
	imgW, imgH := input.Bounds().Dx(), input.Bounds().Dy()
	//example: input image is 300,200 and we want to draw it at 200,100 size
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
