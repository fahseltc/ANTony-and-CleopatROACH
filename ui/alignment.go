package ui

import "image"

type Alignment int

var (
	AlignedPadding = 5
)

const (
	CenterAlignment Alignment = iota
	TopAlignment
	BottomAlignment
	LeftAlignment
	RightAlignment
)

func (alignment Alignment) Align(source image.Rectangle, toBeAligned *image.Rectangle) {
	w := toBeAligned.Dx()
	h := toBeAligned.Dy()

	switch alignment {
	case CenterAlignment:
		centerX := source.Min.X + (source.Dx()-w)/2
		centerY := source.Min.Y + (source.Dy()-h)/2
		toBeAligned.Min = image.Pt(centerX, centerY)
		toBeAligned.Max = image.Pt(centerX+w, centerY+h)

	case TopAlignment:
		centerX := source.Min.X + (source.Dx()-w)/2
		toBeAligned.Min = image.Pt(centerX, source.Min.Y-AlignedPadding-h)
		toBeAligned.Max = image.Pt(centerX+w, source.Min.Y-AlignedPadding)

	case BottomAlignment:
		centerX := source.Min.X + (source.Dx()-w)/2
		toBeAligned.Min = image.Pt(centerX, source.Max.Y+AlignedPadding)
		toBeAligned.Max = image.Pt(centerX+w, source.Max.Y+AlignedPadding+h)

	case LeftAlignment:
		centerY := source.Min.Y + (source.Dy()-h)/2
		toBeAligned.Min = image.Pt(source.Min.X-AlignedPadding-w, centerY)
		toBeAligned.Max = image.Pt(source.Min.X-AlignedPadding, centerY+h)

	case RightAlignment:
		centerY := source.Min.Y + (source.Dy()-h)/2
		toBeAligned.Min = image.Pt(source.Max.X+AlignedPadding, centerY)
		toBeAligned.Max = image.Pt(source.Max.X+AlignedPadding+w, centerY+h)
	}
}
