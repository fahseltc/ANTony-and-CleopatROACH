package scene

import (
	"fmt"
	"gamejam/config"
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

func NewPlayScene(cfg *config.T, fonts *fonts.All) *PlayScene {
	screenWidth := cfg.Resolutions.Internal.Width
	screenHeight := cfg.Resolutions.Internal.Height

	// TODO this should come from the Sim?
	mapWidth := cfg.Resolutions.External.Width
	mapHeight := cfg.Resolutions.External.Height

	cam := ui.NewCamera(screenWidth, screenHeight, mapWidth, mapHeight)

	tileMap := tilemap.NewTilemap()
	s := sim.New(60, tileMap.CollisionRects)
	scene := &PlayScene{
		fonts:   fonts,
		sim:     s,
		ui:      ui.NewUi(cam, fonts, tileMap),
		tileMap: tileMap,
		drag:    ui.NewDrag(),
		sprites: make(map[string]*ui.Sprite),
	}
	u := sim.NewDefaultUnit()
	u.SetPosition(&image.Point{400, 500})
	scene.sim.AddUnit(u)
	//scene.sim.IssueAction(u.ID.String(), sim.MovingAction, &image.Point{X: 300, Y: 400})

	u2 := sim.NewDefaultUnit()
	u2.SetPosition(&image.Point{300, 300})
	scene.sim.AddUnit(u2)
	//scene.sim.IssueAction(u2.ID.String(), sim.MovingAction, &image.Point{X: 600, Y: 200})

	ant := ui.NewDefaultSprite(u.ID)
	scene.sprites[u.ID.String()] = ant
	return scene
}

func (s *PlayScene) Update() error {
	// make sure all the sim units are in the list of sprites
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
	// Draw tiles first as the BG
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(s.ui.Camera.Zoom, s.ui.Camera.Zoom)
	opts.GeoM.Translate(float64(s.ui.Camera.X), float64(s.ui.Camera.Y))
	screen.DrawImage(s.ui.TileMap.StaticBg, opts)
	// Then sprites on top
	for _, sprite := range s.sprites {
		// if s.ui.Camera.IsVisible(sprite.X, sprite.Y, sprite.Width, sprite.Height) {
		sprite.Draw(screen, s.ui.Camera)
		// }
	}
	s.ui.Draw(screen)

	s.drag.Draw(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera:%0.2f, %0.2f", s.ui.Camera.X, s.ui.Camera.Y), 1, 1)
}
