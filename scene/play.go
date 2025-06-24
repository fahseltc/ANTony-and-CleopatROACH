package scene

import (
	"fmt"
	"gamejam/eventing"
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

	tileMap           *tilemap.Tilemap
	drag              *ui.Drag
	constructionMouse *ui.ConstructionMouse

	fonts *fonts.All

	sprites         map[string]*ui.Sprite
	selectedUnitIDs []string

	eventBus *eventing.EventBus
}

func NewPlayScene(fonts *fonts.All) *PlayScene {
	tileMap := tilemap.NewTilemap()
	simulation := sim.New(60, tileMap)
	scene := &PlayScene{
		fonts:             fonts,
		sim:               simulation,
		Ui:                ui.NewUi(fonts, tileMap, simulation),
		tileMap:           tileMap,
		drag:              ui.NewDrag(),
		constructionMouse: &ui.ConstructionMouse{},
		sprites:           make(map[string]*ui.Sprite),
		eventBus:          simulation.EventBus,
	}

	scene.constructionMouse.SetSprite("tilemap/map_tiles/bridge.png")
	u := sim.NewDefaultAnt()
	u.SetTilePosition(10, 12)
	scene.sim.AddUnit(u)
	//scene.sim.IssueAction(u.ID.String(), sim.MovingAction, &image.Point{X: 300, Y: 400})

	u2 := sim.NewDefaultAnt()
	u2.SetTilePosition(9, 11)
	scene.sim.AddUnit(u2)
	//scene.sim.IssueAction(u2.ID.String(), sim.MovingAction, &image.Point{X: 600, Y: 200})

	u3 := sim.NewRoyalAnt()
	u3.SetTilePosition(12, 11)
	scene.sim.AddUnit(u3)

	staticBridge := ui.NewBridgeSprite()
	staticBridge.SetTilePosition(15, 15)
	scene.sprites[staticBridge.Id.String()] = staticBridge

	scene.Ui.Camera.SetZoom(ui.MinZoom)
	scene.Ui.Camera.SetPosition(10, 160)

	h := sim.NewHive(8, 8)
	scene.sim.AddHive(h)

	scene.eventBus.Subscribe("MakeAntButtonClickedEvent", scene.HandleMakeAntButtonClickedEvent)
	scene.eventBus.Subscribe("MakeBridgeButtonClickedEvent", scene.HandleMakeBridgeButtonClickedEvent)

	return scene
}
func (s *PlayScene) HandleMakeAntButtonClickedEvent(event eventing.Event) {
	if len(s.selectedUnitIDs) == 1 {
		hiveID := s.selectedUnitIDs[0]
		unitOrHiveString := s.sim.DetermineUnitOrHiveById(hiveID)
		if unitOrHiveString == "hive" {
			s.eventBus.Publish(eventing.Event{
				Type: "ConstructUnitEvent",
				Data: eventing.ConstructUnitEvent{
					HiveID: hiveID,
				},
			})
		}
	}
}

func (s *PlayScene) HandleMakeBridgeButtonClickedEvent(event eventing.Event) {
	if len(s.selectedUnitIDs) == 1 {
		unitID := s.selectedUnitIDs[0]
		unitOrHiveString := s.sim.DetermineUnitOrHiveById(unitID)
		if unitOrHiveString == "unit" {
			s.constructionMouse.Enabled = true
			// s.eventBus.Publish(eventing.Event{
			// 	Type: "ConstructUnitEvent",
			// 	Data: eventing.ConstructUnitEvent{
			// 		HiveID: hiveID,
			// 	},
			// })
		}
	}
}

func (s *PlayScene) Update() error {
	// make sure all the sim units are in the list of spritess
	for _, unit := range s.sim.GetAllUnits() {
		if s.sprites[unit.ID.String()] == nil {
			// TODO switch based on sim.unitType and use the right sprite for each unit
			switch unit.Type {
			case sim.UnitTypeDefaultAnt:
				s.sprites[unit.ID.String()] = ui.NewDefaultAntSprite(unit.ID)
			case sim.UnitTypeDefaultRoach:
				s.sprites[unit.ID.String()] = ui.NewDefaultRoachSprite(unit.ID)
			case sim.UnitTypeRoyalAnt:
				s.sprites[unit.ID.String()] = ui.NewRoyalAntSprite(unit.ID)
			case sim.UnitTypeRoyalRoach:
				s.sprites[unit.ID.String()] = ui.NewRoyalRoachSprite(unit.ID)
			}

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

	for _, spr := range s.sprites {
		if spr.Type == ui.SpriteTypeStatic {
			continue
		}
		if spr.Selected {
			if !slices.Contains(s.selectedUnitIDs, spr.Id.String()) {
				s.selectedUnitIDs = append(s.selectedUnitIDs, spr.Id.String())
			}
		} else if slices.ContainsFunc(s.selectedUnitIDs, func(id string) bool { return id == spr.Id.String() }) {
			s.selectedUnitIDs = slices.DeleteFunc(s.selectedUnitIDs, func(id string) bool { return id == spr.Id.String() })
		}
	}

	if len(s.selectedUnitIDs) > 0 {
		// Handle 1 unit or building selected
		if len(s.selectedUnitIDs) == 1 {
			unitOrHiveString := s.sim.DetermineUnitOrHiveById(s.selectedUnitIDs[0])
			switch unitOrHiveString {
			case "hive":
				// handle hive
				// show hive UI elements
				if s.Ui.HUD.RightSideState != ui.HiveSelectedState {
					s.Ui.HUD.RightSideState = ui.HiveSelectedState
				}
				// s.eventBus.Publish(eventing.Event{
				// 	Type: "ToggleRightSideHUDEvent",
				// 	Data: eventing.ToggleRightSideHUDEvent{
				// 		Show: true,
				// 	},
				// })
			case "unit":
				// hide HIVE build ui element
				if s.Ui.HUD.RightSideState != ui.UnitSelectedState {
					s.Ui.HUD.RightSideState = ui.UnitSelectedState
				}
				// handle unit and clicks
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) { // activate on buttonRelease to debounce
					mx, my := ebiten.CursorPosition()
					for _, unitId := range s.selectedUnitIDs {
						mapX, mapY := s.Ui.Camera.ScreenPosToMapPos(mx, my)
						s.sim.IssueAction(unitId, sim.AttackMovingAction, &image.Point{X: mapX, Y: mapY})
					}
				}

			default:

			}

			// if inpututil.IsKeyJustReleased(ebiten.KeyQ) {
			// 	for _, unitId := range s.selectedUnitIDs {
			// 		hive, err := s.sim.GetHiveByID(unitId)
			// 		if err != nil {
			// 			continue
			// 		}
			// 		s.eventBus.Publish(eventing.Event{
			// 			Type: "ConstructUnitEvent",
			// 			Data: eventing.ConstructUnitEvent{
			// 				HiveID: hive.ID.String(),
			// 			},
			// 		})
			// 	}
			// }
		} else {
			if s.Ui.HUD.RightSideState != ui.HiddenState {
				s.Ui.HUD.RightSideState = ui.HiddenState
			}

			// handle multiple units/buildings selected
			if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) { // activate on buttonRelease to debounce
				mx, my := ebiten.CursorPosition()
				for _, unitId := range s.selectedUnitIDs {
					mapX, mapY := s.Ui.Camera.ScreenPosToMapPos(mx, my)
					s.sim.IssueAction(unitId, sim.AttackMovingAction, &image.Point{X: mapX, Y: mapY})
				}
			}
		}
	} else {
		// zero units selected - hide the rightside HUD
		if s.Ui.HUD.RightSideState != ui.HiddenState {
			s.Ui.HUD.RightSideState = ui.HiddenState
		}
	}
	s.drag.Update(s.sprites, s.Ui.Camera, s.Ui.HUD)
	s.constructionMouse.Update()
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
	s.constructionMouse.Draw(screen, s.Ui.Camera)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera:%v,%v", s.Ui.Camera.ViewPortX, s.Ui.Camera.ViewPortY), 1, 1)
}
