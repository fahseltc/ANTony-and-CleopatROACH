package scene

import (
	"fmt"
	"gamejam/eventing"
	"gamejam/fonts"
	"gamejam/sim"
	"gamejam/tilemap"
	"gamejam/ui"
	"image"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var PlayerFaction = 0

type PlayScene struct {
	BaseScene
	eventBus *eventing.EventBus
	sim      *sim.T
	Ui       *ui.Ui

	tileMap           *tilemap.Tilemap
	drag              *ui.Drag
	constructionMouse *ui.ConstructionMouse

	fonts *fonts.All

	sprites         map[string]*ui.Sprite
	selectedUnitIDs []string

	// Cutscene stuff
	cutsceneActions   []CutsceneAction
	inCutscene        bool
	currentDialogText string
}

func NewPlayScene(fonts *fonts.All) *PlayScene {
	tileMap := tilemap.NewTilemap()
	simulation := sim.New(60, tileMap)
	constructionMouse := &ui.ConstructionMouse{}
	scene := &PlayScene{
		fonts:             fonts,
		sim:               simulation,
		Ui:                ui.NewUi(fonts, tileMap, simulation),
		tileMap:           tileMap,
		drag:              ui.NewDrag(),
		constructionMouse: constructionMouse,
		sprites:           make(map[string]*ui.Sprite),
		eventBus:          simulation.EventBus,
	}

	scene.constructionMouse.SetSprite("tilemap/bridge.png")
	u := sim.NewDefaultAnt()
	u.SetTilePosition(10, 12)
	scene.sim.AddUnit(u)
	//scene.sim.IssueAction(u.ID.String(), sim.MovingAction, &image.Point{X: 300, Y: 400})

	u2 := sim.NewDefaultAnt()
	u2.SetTilePosition(9, 11)
	scene.sim.AddUnit(u2)
	//scene.sim.IssueAction(u2.ID.String(), sim.MovingAction, &image.Point{X: 600, Y: 200})

	king := sim.NewRoyalAnt()
	king.SetTilePosition(12, 11)
	scene.sim.AddUnit(king)

	queen := sim.NewRoyalRoach()
	queen.SetTilePosition(27, 10)
	queen.Faction = 1
	scene.sim.AddUnit(queen)

	// staticBridge := ui.NewBridgeSprite()
	// staticBridge.SetTilePosition(15, 15)
	// scene.sprites[staticBridge.Id.String()] = staticBridge

	scene.Ui.Camera.SetZoom(ui.MinZoom)
	scene.Ui.Camera.SetPosition(10, 160)

	h := sim.NewHive()
	h.SetTilePosition(8, 8)
	scene.sim.AddBuilding(h)

	scene.eventBus.Subscribe("MakeAntButtonClickedEvent", scene.HandleMakeAntButtonClickedEvent)
	scene.eventBus.Subscribe("MakeBridgeButtonClickedEvent", scene.HandleMakeBridgeButtonClickedEvent)
	scene.eventBus.Subscribe("BuildClickedEvent", scene.HandleBuildClickedEvent)

	// how to cutscene?
	// disable UI user should have access to.
	// show dialogs and force them to be accepted before moving on

	scene.BeginCutscene()
	return scene
}

func (s *PlayScene) BeginCutscene() {
	s.inCutscene = true
	s.currentDialogText = ""
	s.cutsceneActions = []CutsceneAction{
		&FadeCameraAction{Mode: "in", Speed: 3},
		//&PanCameraAction{TargetX: 500, TargetY: 300, Speed: 200},
		&ShowTextAreaAction{Text: "Welcome to Ant World!"},
		&WaitAction{Duration: 2},
		&PanCameraAction{TargetX: float64(27), TargetY: float64(10), Speed: 300},
	}
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
			s.drag.Enabled = false
			// s.eventBus.Publish(eventing.Event{
			// 	Type: "ConstructUnitEvent",
			// 	Data: eventing.ConstructUnitEvent{
			// 		HiveID: hiveID,
			// 	},
			// })
		}
	}
}
func (s *PlayScene) HandleBuildClickedEvent(event eventing.Event) {
	targetRect := event.Data.(eventing.BuildClickedEvent).TargetRect
	if len(s.selectedUnitIDs) == 1 {
		s.sim.ConstructBuilding(targetRect, s.selectedUnitIDs[0])
	}
	s.drag.Enabled = true
	s.constructionMouse.Enabled = false
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
		if s.sprites[building.GetID().String()] == nil {
			switch building.GetType() {
			case sim.BuildingTypeBridge:
				spr := ui.NewBridgeSprite(building.GetID())
				s.sprites[building.GetID().String()] = spr
				s.sprites[building.GetID().String()].SetPosition(building.GetPosition())
			case sim.BuildingTypeHive:
				spr := ui.NewHiveSprite(building.GetID())
				s.sprites[building.GetID().String()] = spr
				s.sprites[building.GetID().String()].SetPosition(building.GetPosition())
			case sim.BuildingTypeInConstruction:
				spr := ui.NewInConstructionSprite(building.GetID())
				s.sprites[building.GetID().String()] = spr
				s.sprites[building.GetID().String()].SetPosition(building.GetPosition())
			}
		} else {
			s.sprites[building.GetID().String()].SetPosition(building.GetPosition())
			s.sprites[building.GetID().String()].ProgressBar.SetProgress(building.GetProgress())
		}
	}
	// remove building & unit sprites that are no longer in the SIM
	s.UpdateRemoveInactiveSprites()

	// HANDLE CUTSCENES - we might want sim.update though

	if s.inCutscene {
		dt := 1.0 / 60.0 // or use actual delta time
		if len(s.cutsceneActions) == 0 {
			s.inCutscene = false
			s.currentDialogText = ""
		} else {
			currentCutScene := s.cutsceneActions[0]
			if currentCutScene.Update(s, dt) {
				s.cutsceneActions = s.cutsceneActions[1:]
			}
			// Early return to skip normal controls
			return nil
		}
	}

	// handle selectedIDs
	for _, spr := range s.sprites {
		if spr.Type == ui.SpriteTypeStatic {
			continue
		}
		unit, err := s.sim.GetUnitByID(spr.Id.String()) // remove unfactioned units from selection
		if err == nil {
			if unit.Faction != 0 {
				spr.Selected = false
				continue
			}
		}
		bld, err := s.sim.GetBuildingByID(spr.Id.String()) // remove unfactioned buildings from selection
		if err == nil {
			if bld.GetFaction() != 0 {
				spr.Selected = false
				continue
			}
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
					s.constructionMouse.Enabled = false
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
					s.constructionMouse.Enabled = false
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
				s.constructionMouse.Enabled = false
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
			s.constructionMouse.Enabled = false
		}
	}
	s.drag.Update(s.sprites, s.Ui.Camera, s.Ui.HUD)
	s.constructionMouse.Update(s.tileMap, s.sim)
	if !s.constructionMouse.Enabled {
		s.drag.Enabled = true
	}
	s.Ui.Update()
	s.sim.Update()

	return nil
}

func (s *PlayScene) UpdateRemoveInactiveSprites() {
	activeIDs := make(map[string]struct{})
	for _, building := range s.sim.GetAllBuildings() {
		activeIDs[building.GetID().String()] = struct{}{}
	}
	for _, unit := range s.sim.GetAllUnits() {
		activeIDs[unit.ID.String()] = struct{}{}
	}
	for id := range s.sprites {
		if _, exists := activeIDs[id]; !exists {
			delete(s.sprites, id)
		}
	}
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	if s.currentDialogText != "" {
		ebitenutil.DebugPrintAt(screen, s.currentDialogText, 100, 100)
	}

	// Draw tiles first as the BG
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(s.Ui.Camera.ViewPortZoom, s.Ui.Camera.ViewPortZoom)
	opts.GeoM.Translate(float64(s.Ui.Camera.ViewPortX), float64(s.Ui.Camera.ViewPortY))
	screen.DrawImage(s.Ui.TileMap.StaticBg, opts)
	//s.DebugDraw(screen)

	// Then Static Sprites
	for _, sprite := range s.sprites {
		if sprite.Type == ui.SpriteTypeStatic {
			sprite.Draw(screen, s.Ui.Camera)
		}
	}

	// Then non-static
	for _, sprite := range s.sprites {
		if sprite.Type != ui.SpriteTypeStatic {
			sprite.Draw(screen, s.Ui.Camera)
		}
	}
	s.Ui.Draw(screen)

	s.drag.Draw(screen)

	s.constructionMouse.Draw(screen, s.Ui.Camera)
	s.Ui.Camera.DrawFade(screen)
}

func (s *PlayScene) DebugDraw(screen *ebiten.Image) {
	for _, mo := range s.tileMap.MapObjects {
		rect := mo.Rect
		x0, y0 := s.Ui.Camera.MapPosToScreenPos(rect.Min.X, rect.Min.Y)
		x1, y1 := s.Ui.Camera.MapPosToScreenPos(rect.Max.X, rect.Max.Y)
		// Draw rectangle outline in red, scaled to viewport
		for x := x0; x < x1; x++ {
			screen.Set(x, y0, color.RGBA{255, 0, 0, 255})
			screen.Set(x, y1-1, color.RGBA{255, 0, 0, 255})
		}
		for y := y0; y < y1; y++ {
			screen.Set(x0, y, color.RGBA{255, 0, 0, 255})
			screen.Set(x1-1, y, color.RGBA{255, 0, 0, 255})
		}
	}
	for _, spr := range s.sprites {
		rect := spr.Rect
		x0, y0 := s.Ui.Camera.MapPosToScreenPos(rect.Min.X, rect.Min.Y)
		x1, y1 := s.Ui.Camera.MapPosToScreenPos(rect.Max.X, rect.Max.Y)
		for x := x0; x < x1; x++ {
			screen.Set(x, y0, color.RGBA{255, 255, 0, 255})
			screen.Set(x, y1-1, color.RGBA{255, 255, 0, 255})
		}
		for y := y0; y < y1; y++ {
			screen.Set(x0, y, color.RGBA{255, 255, 0, 255})
			screen.Set(x1-1, y, color.RGBA{255, 255, 0, 255})
		}
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("camera:%v,%v", s.Ui.Camera.ViewPortX, s.Ui.Camera.ViewPortY), 1, 1)
	mx, my := ebiten.CursorPosition()
	x, y := s.Ui.Camera.ScreenPosToMapPos(mx, my)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mouseCoords:%v,%v", x, y), 1, 40)
}
