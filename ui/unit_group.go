package ui

import (
	"gamejam/fonts"
	"gamejam/sim"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type UnitGroup struct {
	IDs  []string
	Key  ebiten.Key
	Rect image.Rectangle
}

func NewUnitGroup(IDs []string, key ebiten.Key, rect image.Rectangle) *UnitGroup {
	return &UnitGroup{
		IDs:  IDs,
		Key:  key,
		Rect: rect,
	}
}

var (
	UnitHotkeys = [...]ebiten.Key{
		ebiten.Key1, ebiten.Key2,
		ebiten.Key3, ebiten.Key4,
		ebiten.Key5, ebiten.Key6,
		ebiten.Key7, ebiten.Key8,
		ebiten.Key9, ebiten.Key0,
	}

	UnitButtonWidth   = 40
	UnitButtonHeight  = 20
	UnitButtonPadding = 2

	UnitButtonLabelTopPadding  = 8
	UnitButtonLabelLeftPadding = 10
)

type UnitGroupManager struct {
	groups            map[ebiten.Key]*UnitGroup
	rect              image.Rectangle
	lastPressedHotkey ebiten.Key

	fonts *fonts.All

	btnImg *ebiten.Image
}

func NewUnitGroupManager(fonts *fonts.All) *UnitGroupManager {
	return &UnitGroupManager{
		groups:            make(map[ebiten.Key]*UnitGroup),
		rect:              image.Rectangle{Min: image.Pt(190, 478), Max: image.Pt(610, 500)},
		lastPressedHotkey: 0,

		fonts: fonts,

		btnImg: util.LoadImage("ui/btn/unit-hotkey-btn.png"),
	}
}

func (u *UnitGroupManager) Update(selectedUnits []string, camera *Camera, sim *sim.T) []string {
	// Handle mouse clicks inside buttons
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		point := image.Pt(mx, my)
		for _, grp := range u.groups {
			if point.In(grp.Rect) {
				camera.CenterCameraOnUnitGroupByIds(grp.IDs, sim)
				return grp.IDs
			}
		}
	}

	// Handle CONTROL on keyboard pressed (for setting new hotkeys)
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		for i, k := range UnitHotkeys {
			if inpututil.IsKeyJustReleased(k) {
				u.AddGroup(i, selectedUnits, k)
			}
		}
	} else { // else CONTROL not pressed, handle regular numkey presses
		for _, k := range UnitHotkeys {
			if inpututil.IsKeyJustPressed(k) {
				unitGroup := u.groups[k]
				if unitGroup != nil {
					if u.lastPressedHotkey == k { // detect double tap?
						camera.CenterCameraOnUnitGroupByIds(unitGroup.IDs, sim)
					}
					u.lastPressedHotkey = k
					return unitGroup.IDs
				}
			}
		}
	}
	return []string{}
}

func (u *UnitGroupManager) AddGroup(index int, IDs []string, key ebiten.Key) {
	if len(IDs) <= 0 {
		return
	}

	offset := index * (UnitButtonWidth + UnitButtonPadding)
	rect := image.Rectangle{
		Min: image.Point{X: u.rect.Min.X + offset, Y: u.rect.Min.Y},
		Max: image.Point{X: u.rect.Min.X + offset + UnitButtonWidth, Y: u.rect.Min.Y + UnitButtonHeight},
	}

	copiedIDs := append([]string(nil), IDs...)
	u.groups[key] = NewUnitGroup(copiedIDs, key, rect)
}

func (u *UnitGroupManager) Draw(screen *ebiten.Image) {
	x := u.rect.Min.X
	y := u.rect.Min.Y
	for _, key := range UnitHotkeys { // there are ten of these
		if u.groups[key] != nil {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(u.btnImg, opts)
			// draw number on it
			util.DrawCenteredText(screen, u.fonts.XSmall, u.keyStringLookup(key), x+UnitButtonLabelLeftPadding, y+UnitButtonLabelTopPadding, nil)
		}
		x += UnitButtonWidth + UnitButtonPadding
	}
}

func (u *UnitGroupManager) keyStringLookup(key ebiten.Key) string {
	switch key {
	case ebiten.Key0:
		return "0"
	case ebiten.Key1:
		return "1"
	case ebiten.Key2:
		return "2"
	case ebiten.Key3:
		return "3"
	case ebiten.Key4:
		return "4"
	case ebiten.Key5:
		return "5"
	case ebiten.Key6:
		return "6"
	case ebiten.Key7:
		return "7"
	case ebiten.Key8:
		return "8"
	case ebiten.Key9:
		return "9"
	default:
		return "."
	}
}
