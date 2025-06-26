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

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlayScene2 struct {
	BaseScene
	eventBus *eventing.EventBus
	sim      *sim.T
	Ui       *ui.Ui

	tileMap           *tilemap.Tilemap
	drag              *ui.Drag
	constructionMouse *ui.ConstructionMouse

	fonts *fonts.All

	Sprites         map[string]*ui.Sprite
	selectedUnitIDs []string

	// Cutscene stuff
	cutsceneActions []CutsceneAction
	inCutscene      bool
	currentDialog   *ui.PortraitTextArea

	// Tutorial stuff
	tutorialDialogs []ui.Tutorial
	inTutorial      bool

	// CompletionCondition *SceneCompletion
	CompletionCondition *SceneCompletion
	SceneCompleted      bool
}

func NewPlayScene2(fonts *fonts.All) *PlayScene2 {
	tileMap := tilemap.NewTilemap("assets/tilemap/map2.tmx")
	simulation := sim.New(60, tileMap)
	constructionMouse := &ui.ConstructionMouse{}
	scene := &PlayScene2{
		fonts:             fonts,
		sim:               simulation,
		Ui:                ui.NewUi(fonts, tileMap, simulation),
		tileMap:           tileMap,
		drag:              ui.NewDrag(),
		constructionMouse: constructionMouse,
		Sprites:           make(map[string]*ui.Sprite),
		eventBus:          simulation.EventBus,
	}
	scene.constructionMouse.SetSprite("tilemap/bridge.png")

	scene.eventBus.Subscribe("MakeAntButtonClickedEvent", scene.HandleMakeAntButtonClickedEvent)
	scene.eventBus.Subscribe("MakeBridgeButtonClickedEvent", scene.HandleMakeBridgeButtonClickedEvent)
	scene.eventBus.Subscribe("BuildClickedEvent", scene.HandleBuildClickedEvent)

	scene.setupLvl1()
	return scene
}

func (scene *PlayScene) setupLvl1() {
	u := sim.NewDefaultAnt()
	u.SetTilePosition(10, 12)
	scene.sim.AddUnit(u)

	u2 := sim.NewDefaultAnt()
	u2.SetTilePosition(7, 12)
	scene.sim.AddUnit(u2)

	king := sim.NewRoyalAnt()
	king.SetTilePosition(12, 11)
	scene.sim.AddUnit(king)

	queen := sim.NewRoyalRoach()
	queen.SetTilePosition(27, 10)
	queen.Faction = 1
	scene.sim.AddUnit(queen)

	scene.Ui.Camera.SetZoom(ui.MinZoom)
	scene.Ui.Camera.SetPosition(10, 160)
	scene.Ui.Camera.FadeAlpha = 255

	h := sim.NewHive()
	h.SetTilePosition(8, 8)
	scene.sim.AddBuilding(h)

	scene.CompletionCondition = NewSceneCompletion(queen, king, scene.tileMap.MapCompletionObjects[0].Rect)

	scene.beginCutscene1(king.ID.String(), queen.ID.String())
}

func (s *PlayScene) beginCutscene1(antony string, cleopatroach string) {
	s.inCutscene = true
	s.Ui.DrawEnabled = false
	s.drag.Enabled = false
	s.cutsceneActions = []CutsceneAction{
		&FadeCameraAction{Mode: "in", Speed: 2},
		&ShowPortraitTextAreaAction{
			portraitTextArea: ui.NewPortraitTextArea(
				s.fonts,
				"Antony: Welcome to da Ant Game! Where my gurl at? I need to build a bridge to get to her!",
				"portraits/ant-royalty.png",
			),
		},
		&IssueUnitCommandAction{
			unitID:     antony,
			targetTile: &image.Point{X: 14, Y: 11},
		},
		&PanCameraAction{TargetX: float64(27), TargetY: float64(10), Speed: 300},
		&IssueUnitCommandAction{
			unitID:     cleopatroach,
			targetTile: &image.Point{X: 26, Y: 10},
		},
		&ShowPortraitTextAreaAction{
			portraitTextArea: ui.NewPortraitTextArea(
				s.fonts,
				"Cleopatroach: There's beggary in the love that can be reckon'd. Save me Antony!",
				"portraits/cockroach-royalty.png",
			),
		},
		&PanCameraAction{TargetX: float64(3), TargetY: float64(6), Speed: 300},
		&ShowPortraitTextAreaAction{
			portraitTextArea: ui.NewPortraitTextArea(
				s.fonts,
				"Antony: Gotchu babe! Put these ants to work!",
				"portraits/ant-royalty.png",
			),
		},
	}

	s.tutorialDialogs = []ui.Tutorial{
		ui.NewTutorialStep(
			"tutorials/tutorial-1.png",
			&image.Rectangle{Min: image.Point{X: 412, Y: 341}, Max: image.Point{X: 800, Y: 600}},
			nil, nil,
		),
		ui.NewTutorialStep(
			"tutorials/tutorial-2.png",
			&image.Rectangle{Min: image.Point{X: 412, Y: 341}, Max: image.Point{X: 800, Y: 600}},
			nil, nil,
		),
	}

}
func (s *PlayScene) SceneCompleteCutscene() {
	s.inCutscene = true
	s.inCutscene = true
	s.Ui.DrawEnabled = false
	s.drag.Enabled = false

	s.selectedUnitIDs = []string{} // clear selected unit IDs

	// make them move towards each other
	cleopatroach := s.CompletionCondition.UnitOne.ID.String()
	antony := s.CompletionCondition.UnitTwo.ID.String()

	s.cutsceneActions = []CutsceneAction{
		&IssueUnitCommandAction{
			unitID:     cleopatroach,
			targetTile: &image.Point{X: 27, Y: 5},
		},
		&IssueUnitCommandAction{
			unitID:     antony,
			targetTile: &image.Point{X: 27, Y: 13},
		},
		&WaitAction{
			Duration: 1.0, // wait for 1 second
		},
		&IssueUnitCommandAction{
			unitID:     cleopatroach,
			targetTile: &image.Point{X: 27, Y: 8},
		},

		&IssueUnitCommandAction{
			unitID:     antony,
			targetTile: &image.Point{X: 27, Y: 10},
		},
		&WaitAction{
			Duration: 1.0, // wait for 1 second
		},
		&PanCameraAction{TargetX: float64(30), TargetY: float64(10), Speed: 300},
		&ZoomCameraAction{
			TargetZoom: 0.8,
			Speed:      1,
			FocusX:     29 * 128,
			FocusY:     9 * 128,
		},
		&DrawTemporarySpriteAction{
			spr:            ui.NewHeartSprite(uuid.New()),
			TargetPosition: &image.Point{X: 3525, Y: 1200},
			MaxDuration:    180,
		},
		&ShowPortraitTextAreaAction{
			portraitTextArea: ui.NewPortraitTextArea(
				s.fonts,
				"Antony: Saved u gurl",
				"portraits/ant-royalty.png",
			),
		},
		&ShowPortraitTextAreaAction{
			portraitTextArea: ui.NewPortraitTextArea(
				s.fonts,
				"Cleopatroach: **smooch noises**",
				"portraits/cockroach-royalty.png",
			),
		},
		&FadeCameraAction{Mode: "out", Speed: 1},
	}
}

func (s *PlayScene2) HandleMakeAntButtonClickedEvent(event eventing.Event) {
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

func (s *PlayScene2) HandleMakeBridgeButtonClickedEvent(event eventing.Event) {
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
func (s *PlayScene2) HandleBuildClickedEvent(event eventing.Event) {
	targetRect := event.Data.(eventing.BuildClickedEvent).TargetRect
	if len(s.selectedUnitIDs) == 1 {
		s.sim.ConstructBuilding(targetRect, s.selectedUnitIDs[0])
	}
	s.drag.Enabled = true
	s.constructionMouse.Enabled = false
}

func (s *PlayScene2) Update() error {
	if s.CompletionCondition.IsComplete(s.sim) && !s.SceneCompleted {
		s.SceneCompleted = true
		s.SceneCompleteCutscene()
	}
	// make sure all the sim units are in the list of spritess
	for _, unit := range s.sim.GetAllUnits() {
		if s.Sprites[unit.ID.String()] == nil {
			// TODO switch based on sim.unitType and use the right sprite for each unit
			switch unit.Type {
			case sim.UnitTypeDefaultAnt:
				s.Sprites[unit.ID.String()] = ui.NewDefaultAntSprite(unit.ID)
			case sim.UnitTypeDefaultRoach:
				s.Sprites[unit.ID.String()] = ui.NewDefaultRoachSprite(unit.ID)
			case sim.UnitTypeRoyalAnt:
				s.Sprites[unit.ID.String()] = ui.NewRoyalAntSprite(unit.ID)
			case sim.UnitTypeRoyalRoach:
				s.Sprites[unit.ID.String()] = ui.NewRoyalRoachSprite(unit.ID)
			}

		} else {
			// else update sprites to match their sim positions
			s.Sprites[unit.ID.String()].SetPosition(unit.Position)
			s.Sprites[unit.ID.String()].SetAngle(unit.MovingAngle)
		}
	}
	// same for buildings
	for _, building := range s.sim.GetAllBuildings() {
		if s.Sprites[building.GetID().String()] == nil {
			switch building.GetType() {
			case sim.BuildingTypeBridge:
				spr := ui.NewBridgeSprite(building.GetID())
				s.Sprites[building.GetID().String()] = spr
				s.Sprites[building.GetID().String()].SetPosition(building.GetPosition())
			case sim.BuildingTypeHive:
				spr := ui.NewHiveSprite(building.GetID())
				s.Sprites[building.GetID().String()] = spr
				s.Sprites[building.GetID().String()].SetPosition(building.GetPosition())
			case sim.BuildingTypeInConstruction:
				spr := ui.NewInConstructionSprite(building.GetID())
				s.Sprites[building.GetID().String()] = spr
				s.Sprites[building.GetID().String()].SetPosition(building.GetPosition())
			}
		} else {
			s.Sprites[building.GetID().String()].SetPosition(building.GetPosition())
			s.Sprites[building.GetID().String()].ProgressBar.SetProgress(building.GetProgress())
		}
	}
	// remove building & unit sprites that are no longer in the SIM
	s.UpdateRemoveInactiveSprites()

	// Update sim before cutscenes so things happen in the world as they play.
	s.sim.Update()
	// HANDLE CUTSCENES - we might want sim.update though

	if s.inCutscene {
		dt := 1.0 / 60.0 // or use actual delta time
		if len(s.cutsceneActions) == 0 {
			if s.SceneCompleted {
				//s.BaseScene.sm.SwitchTo(NewPlayScene(s.fonts))
			}
			s.inCutscene = false
			s.Ui.DrawEnabled = true
			s.drag.Enabled = true
			s.constructionMouse.Enabled = true
		} else {
			currentCutScene := s.cutsceneActions[0]
			if s.currentDialog != nil {
				s.currentDialog.Update()
			}
			if currentCutScene.Update(s, dt) {
				s.cutsceneActions = s.cutsceneActions[1:]
			}
			// Early return to skip normal controls
			return nil
		}
	}

	if len(s.tutorialDialogs) > 0 && !s.inCutscene {
		// Check if any tutorial dialog is active
		s.inTutorial = true
		s.tutorialDialogs[0].CheckTrigger(s.sim) // Check the first tutorial dialog trigger
		if s.tutorialDialogs[0].IsComplete() {
			s.tutorialDialogs = s.tutorialDialogs[1:] // Remove the completed dialog
			if len(s.tutorialDialogs) == 0 {
				s.inTutorial = false // No more tutorial dialogs
			}
		}
		if !s.inTutorial {
			s.currentDialog = nil // no active tutorial dialog
		}
	} else {
		s.inTutorial = false
		s.currentDialog = nil // no active tutorial dialog
	}

	// handle selectedIDs
	for _, spr := range s.Sprites {
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
						s.sim.IssueAction(unitId, &image.Point{X: mapX, Y: mapY})
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
					s.sim.IssueAction(unitId, &image.Point{X: mapX, Y: mapY})
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
	s.drag.Update(s.Sprites, s.Ui.Camera, s.Ui.HUD)
	s.constructionMouse.Update(s.tileMap, s.sim)
	if !s.constructionMouse.Enabled {
		s.drag.Enabled = true
	}
	s.Ui.Update()

	return nil
}

func (s *PlayScene2) UpdateRemoveInactiveSprites() {
	activeIDs := make(map[string]struct{})
	for _, building := range s.sim.GetAllBuildings() {
		activeIDs[building.GetID().String()] = struct{}{}
	}
	for _, unit := range s.sim.GetAllUnits() {
		activeIDs[unit.ID.String()] = struct{}{}
	}
	for id, spr := range s.Sprites {
		if spr.Type == ui.SpriteTypeStatic {
			continue // static sprites are not removed
		}
		if _, exists := activeIDs[id]; !exists {
			delete(s.Sprites, id)
		}
	}
}

func (s *PlayScene2) Draw(screen *ebiten.Image) {
	// Draw tiles first as the BG
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(s.Ui.Camera.ViewPortZoom, s.Ui.Camera.ViewPortZoom)
	opts.GeoM.Translate(float64(s.Ui.Camera.ViewPortX), float64(s.Ui.Camera.ViewPortY))
	screen.DrawImage(s.Ui.TileMap.StaticBg, opts)
	//s.DebugDraw(screen)

	// Then Static Sprites
	for _, sprite := range s.Sprites {
		if sprite.Type == ui.SpriteTypeStatic {
			sprite.Draw(screen, s.Ui.Camera)
		}
	}

	// Then non-static
	for _, sprite := range s.Sprites {
		if sprite.Type != ui.SpriteTypeStatic {
			sprite.Draw(screen, s.Ui.Camera)
		}
	}
	s.Ui.Draw(screen)

	s.drag.Draw(screen)

	s.constructionMouse.Draw(screen, s.Ui.Camera)

	if s.currentDialog != nil {
		s.currentDialog.Draw(screen)
	}

	if len(s.tutorialDialogs) > 0 && !s.inCutscene {
		// Check if any tutorial dialog is active
		s.inTutorial = true
		s.tutorialDialogs[0].Draw(screen)
	}
	s.Ui.Camera.DrawFade(screen) // this should always be drawn last
}

func (s *PlayScene2) DebugDraw(screen *ebiten.Image) {
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
	for _, spr := range s.Sprites {
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mouseScreenCoords:%v,%v", mx, my), 1, 40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mouseMapCoords:%v,%v", x, y), 1, 60)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("zoom:%v", s.Ui.Camera.ViewPortZoom), 1, 20)

}
