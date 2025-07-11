package ui

import (
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	IconDimension = 35
	IconPadding   = 8

	SelectedUnitCols = 10
	SelectedUnitRows = 2
)

type SelectedUnitArea struct {
	rect        image.Rectangle
	selectedIDs []string

	workerIcon *ebiten.Image
}

func NewSelectedUnitArea() *SelectedUnitArea {
	sua := &SelectedUnitArea{
		rect:       image.Rectangle{Min: image.Pt(190, 510), Max: image.Pt(620, 600)},
		workerIcon: util.ScaleImage(util.LoadImage("TEXTURE_MISSING.png"), float32(IconDimension), float32(IconDimension)),
	}

	return sua
}

func (s *SelectedUnitArea) Update(selectedIDs []string) {
	s.selectedIDs = selectedIDs
}

func (s *SelectedUnitArea) Draw(screen *ebiten.Image) {
	if len(s.selectedIDs) == 0 {
		return
	}
	unitCount := 0
	for rowNum := range SelectedUnitRows {
		for colNum := range SelectedUnitCols {
			x := float64(s.rect.Min.X + (colNum * (IconDimension + IconPadding)))
			y := float64(s.rect.Min.Y + (rowNum * (IconDimension + IconPadding)))
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(x, y)
			screen.DrawImage(s.workerIcon, opts)
			unitCount++
			if unitCount >= len(s.selectedIDs) {
				return
			}
		}
	}
}
