package fonts

import (
	"gamejam/assets"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type All struct {
	XXSmall text.Face
	XSmall  text.Face
	Small   text.Face
	Med     text.Face
	Large   text.Face
	XLarge  text.Face
}

func Load(path string) *All {
	fonts := &All{}
	xxs, _ := loadTTFFont(path, 4)
	fonts.XXSmall = xxs
	xs, _ := loadTTFFont(path, 8)
	fonts.XSmall = xs
	s, _ := loadTTFFont(path, 12)
	fonts.Small = s
	m, _ := loadTTFFont(path, 20)
	fonts.Med = m
	l, _ := loadTTFFont(path, 30)
	fonts.Large = l
	xl, _ := loadTTFFont(path, 42)
	fonts.XLarge = xl
	return fonts
}

func loadTTFFont(path string, size float64) (text.Face, error) {
	fontFile, err := assets.Files.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		return nil, err
	}
	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
