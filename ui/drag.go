package ui

import (
	"gamejam/log"
	"image"
	"image/color"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Drag struct {
	Enabled          bool
	dragRect         image.Rectangle
	firstClickPoint  image.Point
	secondClickPoint image.Point
	log              *slog.Logger
}

func NewDrag() *Drag {
	return &Drag{
		Enabled:         true,
		dragRect:        image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(0, 0)},
		firstClickPoint: image.Pt(0, 0),
		log:             log.NewLogger().With("for", "Drag"),
	}
}

func (d *Drag) Update(sprites map[string]*Sprite, camera *Camera, HUD *HUD) {
	if !d.Enabled {
		return
	}
	mx, my := ebiten.CursorPosition()
	pt := image.Point{X: mx, Y: my}

	// pt.In(HUD.leftSideRect) || REMOVED not inUSE
	if HUD.RightSideState != HiddenState && pt.In(HUD.rightSideRect) { // abort updating selected units if the click is inside the UI elements
		d.dragRect = image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(0, 0)}
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		d.firstClickPoint = image.Point{X: mx, Y: my}
	}
	// Detect if the mouse is being held down
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if d.firstClickPoint.Eq(image.Pt(0, 0)) {
			return
		}
		d.secondClickPoint = image.Point{X: mx, Y: my}
		minX := d.firstClickPoint.X
		maxX := mx
		if mx < d.firstClickPoint.X {
			minX = mx
			maxX = d.firstClickPoint.X
		}
		minY := d.firstClickPoint.Y
		maxY := my
		if my < d.firstClickPoint.Y {
			minY = my
			maxY = d.firstClickPoint.Y
		}

		d.dragRect = image.Rectangle{Min: image.Pt(minX, minY), Max: image.Pt(maxX, maxY)}
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		var selectedIDs []string
		mapRect := image.Rectangle{
			Min: image.Pt(camera.ScreenPosToMapPos(d.dragRect.Min.X, d.dragRect.Min.Y)),
			Max: image.Pt(camera.ScreenPosToMapPos(d.dragRect.Max.X+4, d.dragRect.Max.Y+4))}
		d.dragRect = image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(0, 0)}
		for _, sprite := range sprites {
			if sprite.Type == SpriteTypeStatic {
				sprite.Selected = false
				continue
			}
			if sprite.Rect.Overlaps(mapRect) {
				selectedIDs = append(selectedIDs, sprite.Id.String())
				sprite.Selected = true
			} else {
				sprite.Selected = false
			}
		}

		d.log.Info("Units Selected", "array", selectedIDs)
	}
}

func (d *Drag) Draw(screen *ebiten.Image) {
	if !d.firstClickPoint.Eq(image.Pt(0, 0)) {
		x, y := float64(d.dragRect.Min.X), float64(d.dragRect.Min.Y)
		w, h := float64(d.dragRect.Dx()), float64(d.dragRect.Dy())
		c := color.RGBA{127, 255, 0, 255}
		ebitenutil.DrawLine(screen, x, y, x+w, y, c)     // top
		ebitenutil.DrawLine(screen, x, y, x, y+h, c)     // left
		ebitenutil.DrawLine(screen, x+w, y, x+w, y+h, c) // right
		ebitenutil.DrawLine(screen, x, y+h, x+w, y+h, c) // bottom
	}

	// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("drag:(%v,%v) to (%v,%v)",
	// 	d.dragRect.Min.X,
	// 	d.dragRect.Min.Y,
	// 	d.dragRect.Max.X,
	// 	d.dragRect.Max.Y), 1, 40)
}
