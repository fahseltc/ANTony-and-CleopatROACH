package ui

import (
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	IconDimension = 35
	IconPadding   = 8

	UiLifeBarWidth  = 29
	UiLifeBarHeight = 3

	SelectedUnitCols = 10
	SelectedUnitRows = 2
)

type SelectedUnitArea struct {
	rect        image.Rectangle
	selectedIDs []string

	workerIcon *ebiten.Image
	hiveIcon   *ebiten.Image
}

func NewSelectedUnitArea() *SelectedUnitArea {
	sua := &SelectedUnitArea{
		rect:       image.Rectangle{Min: image.Pt(190, 510), Max: image.Pt(620, 600)},
		workerIcon: util.ScaleImage(util.LoadImage("ui/icon/ant.png"), float32(IconDimension), float32(IconDimension)),
		hiveIcon:   util.ScaleImage(util.LoadImage("ui/icon/hive.png"), float32(IconDimension), float32(IconDimension)),
	}

	return sua
}

func (s *SelectedUnitArea) Update(selectedIDs []string) {
	s.selectedIDs = selectedIDs
}

func (s *SelectedUnitArea) Draw(screen *ebiten.Image, sprites map[string]*Sprite) {
	if len(s.selectedIDs) == 0 {
		return
	}
	unitCount := 0
	for rowNum := range SelectedUnitRows {
		for colNum := range SelectedUnitCols {
			if unitCount > len(s.selectedIDs) {
				return
			}
			id := s.selectedIDs[unitCount]
			spr := sprites[id]
			if spr == nil {
				return
			}

			x := float64(s.rect.Min.X + (colNum * (IconDimension + IconPadding)))
			y := float64(s.rect.Min.Y + (rowNum * (IconDimension + IconPadding)))
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(x, y)
			switch spr.Type {
			case SpriteTypeUnit:
				screen.DrawImage(s.workerIcon, opts)
			case SpriteTypeHive:
				screen.DrawImage(s.hiveIcon, opts)
			}
			s.drawLifeBar(x, y, spr, screen)

			unitCount++
			if unitCount >= len(s.selectedIDs) {
				return
			}
		}
	}
}
func (s *SelectedUnitArea) drawLifeBar(x, y float64, spr *Sprite, screen *ebiten.Image) {
	fgWidth := int(float64(UiLifeBarWidth) * spr.HealthBar.Progress)
	if fgWidth > 0 {
		fg := ebiten.NewImage(fgWidth, UiLifeBarHeight)
		fg.Fill(spr.HealthBar.FgColor)
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(x+float64((IconDimension-UiLifeBarWidth)/2), y+float64(IconDimension)-4)
		screen.DrawImage(fg, op2)
	}
}
