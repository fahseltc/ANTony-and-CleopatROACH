package scene

import (
	"fmt"
	"gamejam/fonts"
	"gamejam/sim"
	"gamejam/tilemap"
	"gamejam/ui"
	"image"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlayScene struct {
	BaseScene
	sim *sim.T
	Ui  *ui.Ui

	tileMap *tilemap.Tilemap
	drag    *ui.Drag

	fonts *fonts.All

	sprites         map[string]*ui.Sprite
	selectedUnitIDs []string
}

func NewPlayScene(fonts *fonts.All) *PlayScene {
	tileMap := tilemap.NewTilemap()
	scene := &PlayScene{
		fonts:   fonts,
		sim:     sim.New(60, tileMap),
		Ui:      ui.NewUi(fonts, tileMap),
		tileMap: tileMap,
		drag:    ui.NewDrag(),
		sprites: make(map[string]*ui.Sprite),
	}
	u := sim.NewDefaultUnit()
	u.SetTilePosition(10, 12)
	//u.Faction = 1
	scene.sim.AddUnit(u)
	//scene.sim.IssueAction(u.ID.String(), sim.MovingAction, &image.Point{X: 300, Y: 400})

	u2 := sim.NewDefaultUnit()
	u2.SetTilePosition(9, 11)
	scene.sim.AddUnit(u2)
	//scene.sim.IssueAction(u2.ID.String(), sim.MovingAction, &image.Point{X: 600, Y: 200})
	scene.Ui.Camera.SetZoom(ui.MinZoom)
	scene.Ui.Camera.SetPosition(10, 160)

	h := sim.NewHive(8, 8)
	scene.sim.AddHive(h)

	return scene
}

func (s *PlayScene) Update() error {
	// make sure all the sim units are in the list of spritess
	for _, unit := range s.sim.GetAllUnits() {
		if s.sprites[unit.ID.String()] == nil {
			spr := ui.NewDefaultSprite(unit.ID)
			s.sprites[unit.ID.String()] = spr

		} else {
			// else update sprites to match their sim positions
			s.sprites[unit.ID.String()].SetPosition(unit.Position)
			s.sprites[unit.ID.String()].SetAngle(unit.MovingAngle)
		}
	}
	// same for buildings
	for _, building := range s.sim.GetAllBuildings() {
		if s.sprites[building.ID.String()] == nil {
			spr := ui.NewHiveSprite(building.ID)
			s.sprites[building.ID.String()] = spr

		} else {
			s.sprites[building.ID.String()].SetPosition(building.Position)
		}
	}

	s.drag.Update(s.sprites, s.Ui.Camera)
	for _, spr := range s.sprites {
		if spr.Selected {
			if !slices.Contains(s.selectedUnitIDs, spr.Id.String()) {
				s.selectedUnitIDs = append(s.selectedUnitIDs, spr.Id.String())
			}
		} else if slices.ContainsFunc(s.selectedUnitIDs, func(id string) bool { return id == spr.Id.String() }) {
			s.selectedUnitIDs = slices.DeleteFunc(s.selectedUnitIDs, func(id string) bool { return id == spr.Id.String() })
		}
	}
	if len(s.selectedUnitIDs) > 0 {
		// todo: read UI state and send the right action
		// Debounce right mouse button so action only runs once per click
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			mx, my := ebiten.CursorPosition()
			for _, unitId := range s.selectedUnitIDs {
				mapX, mapY := s.Ui.Camera.ScreenPosToMapPos(mx, my)
				s.sim.IssueAction(unitId, sim.AttackMovingAction, &image.Point{X: mapX, Y: mapY})
			}
		}
	}

	s.Ui.Update()
	s.sim.Update()

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	// Draw tiles first as the BG
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(s.Ui.Camera.ViewPortZoom, s.Ui.Camera.ViewPortZoom)
	opts.GeoM.Translate(float64(s.Ui.Camera.ViewPortX), float64(s.Ui.Camera.ViewPortY))
	screen.DrawImage(s.Ui.TileMap.StaticBg, opts)

	// Then sprites on top
	for _, sprite := range s.sprites {
		sprite.Draw(screen, s.Ui.Camera)
	}
	s.Ui.Draw(screen)

	s.drag.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera:%v,%v", s.Ui.Camera.ViewPortX, s.Ui.Camera.ViewPortY), 1, 1)
}
