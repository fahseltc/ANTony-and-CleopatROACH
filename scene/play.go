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
	ui  *ui.Ui

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
		sim:     sim.New(60),
		ui:      ui.NewUi(fonts, tileMap),
		tileMap: tileMap,
		drag:    ui.NewDrag(),
		sprites: make(map[string]*ui.Sprite),
	}
	u := sim.NewDefaultUnit()
	scene.sim.AddUnit(u)
	scene.sim.IssueAction(u.ID.String(), sim.MovingAction, &image.Point{X: 200, Y: 200})

	u2 := sim.NewDefaultUnit()
	scene.sim.AddUnit(u2)
	scene.sim.IssueAction(u2.ID.String(), sim.MovingAction, &image.Point{X: 400, Y: 200})

	ant := ui.NewDefaultSprite(u.ID)
	scene.sprites[u.ID.String()] = ant
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
		}
	}

	s.drag.Update(s.sprites)
	for _, spr := range s.sprites {
		if spr.Selected {
			s.selectedUnitIDs = append(s.selectedUnitIDs, spr.Id.String())
		} else if slices.ContainsFunc(s.selectedUnitIDs, func(id string) bool { return id == spr.Id.String() }) {
			s.selectedUnitIDs = slices.DeleteFunc(s.selectedUnitIDs, func(id string) bool { return id == spr.Id.String() })
		}
	}
	if len(s.selectedUnitIDs) > 0 {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			mx, my := ebiten.CursorPosition()
			for _, unitId := range s.selectedUnitIDs {
				mapX, mapY := s.ui.Camera.ScreenPosToMapPos(mx, my)
				s.sim.IssueAction(unitId, sim.AttackMovingAction, &image.Point{X: mapX, Y: mapY})
			}
		}
	}

	s.ui.Update()
	s.sim.Update()

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	s.ui.Draw(screen)
	for _, sprite := range s.sprites {
		sprite.Draw(screen, s.ui.Camera)
	}
	s.drag.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera:%v,%v", s.ui.Camera.ViewPortX, s.ui.Camera.ViewPortY), 1, 1)
}
