package environment

import (
	"bytes"
	"gamejam/assets"
	"io"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Fonts struct {
	XSmall text.Face
	Small  text.Face
	Med    text.Face
	Large  text.Face
}

var fontPath = "fonts/PressStart2P-Regular.ttf"

func NewFontsCollection() *Fonts {
	fonts := &Fonts{}
	xs, _ := loadTTFFont(fontPath, 8)
	fonts.XSmall = xs
	s, _ := loadTTFFont(fontPath, 12)
	fonts.Small = s
	m, _ := loadTTFFont(fontPath, 20)
	fonts.Med = m
	l, _ := loadTTFFont(fontPath, 30)
	fonts.Large = l
	return fonts
}

func loadTTFFont(path string, size float64) (text.Face, error) {
	fontFile, err := assets.Files.Open(path)
	if err != nil {
		return nil, err
	}
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {
		return nil, err
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		return nil, err
	}
	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
