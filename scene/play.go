package scene

import (
	"fmt"
	"gamejam/audio"
	"gamejam/config"
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
	Config *config.T // embedded from
	BaseScene
	LevelData   *LevelData
	sound       *audio.SoundManager
	songStarted bool

	QueenID string
	KingID  string

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
	tutorialDialogs []Tutorial
	inTutorial      bool

	// Level completion
	CompletionCondition *SceneCompletion
	SceneCompleted      bool

	// Notifications
	CurrentNotification *ui.Notification

	Pause *ui.Pause
}

func NewPlayScene(fonts *fonts.All, sound *audio.SoundManager, levelData LevelData) *PlayScene {
	config, _ := config.New()

	tileMap := tilemap.NewTilemap(levelData.TileMapPath)
	simulation := sim.New(60, tileMap)
	constructionMouse := &ui.ConstructionMouse{}
	scene := &PlayScene{
		Config:            config,
		sound:             sound,
		LevelData:         &levelData,
		fonts:             fonts,
		sim:               simulation,
		Ui:                ui.NewUi(fonts, tileMap, simulation),
		tileMap:           tileMap,
		drag:              ui.NewDrag(),
		constructionMouse: constructionMouse,
		Sprites:           make(map[string]*ui.Sprite),
		eventBus:          simulation.EventBus,
		Pause:             ui.NewPause(sound, *fonts),
	}
	scene.constructionMouse.SetSprite("tilemap/bridge.png")
	scene.eventBus.Subscribe("MakeAntButtonClickedEvent", scene.HandleMakeAntButtonClickedEvent)
	scene.eventBus.Subscribe("MakeBridgeButtonClickedEvent", scene.HandleMakeBridgeButtonClickedEvent)
	scene.eventBus.Subscribe("BuildClickedEvent", scene.HandleBuildClickedEvent)
	scene.eventBus.Subscribe("NotEnoughResourcesEvent", scene.NotEnoughResourcesEvent)

	scene.QueenID, scene.KingID = levelData.SetupFunc(scene)

	scene.setupSFX()
	levelData.SetupInitialCutscene(scene, scene.QueenID, scene.KingID)

	return scene
}

func (s *PlayScene) NotEnoughResourcesEvent(event eventing.Event) {
	resName := event.Data.(eventing.NotEnoughResourcesEvent).ResourceName
	target := event.Data.(eventing.NotEnoughResourcesEvent).TargetBeingBuilt

	var str string
	if target == "Bridge" {
		str = fmt.Sprintf("Not enough %v to build %v\nOr builder is not close enough!", resName, target)
	} else {
		str = fmt.Sprintf("Not enough %v to build %v", resName, target)
	}

	s.CurrentNotification = ui.NewNotification(&s.fonts.Med, str)
}

func (s *PlayScene) setupSFX() {
	//s.sound.GlobalVolume = 0.5
	//s.eventBus.Subscribe("PlayWalkSFX", s.sound.PlayWalkSFX)
	s.eventBus.Subscribe("PlayIssueActionSFX", s.sound.PlayIssueActionSFX)
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
		success := s.sim.ConstructBuilding(targetRect, s.selectedUnitIDs[0])
		if !success {
			s.eventBus.Publish(eventing.Event{
				Type: "NotEnoughResourcesEvent",
				Data: eventing.NotEnoughResourcesEvent{ // todo: add reason why, for example "unit not close enough" etc
					ResourceName:     "Wood",
					TargetBeingBuilt: "Bridge",
				},
			})
		}
	}
	s.drag.Enabled = true
	s.constructionMouse.Enabled = false
}

func (s *PlayScene) Update() error {
	if !s.songStarted {
		s.songStarted = true
		s.sound.Play("msx_gamesong1")
	}
	s.sound.Update()

	// Determine Pause State
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.Pause.Hidden = !s.Pause.Hidden
	}
	if !s.Pause.Hidden { // stop the game processing when paused!
		s.Pause.Update()
		return nil
	}

	if s.CompletionCondition.IsComplete(s.sim) && !s.SceneCompleted {
		s.SceneCompleted = true
		s.LevelData.SetupCompletionCutscene(s, s.QueenID, s.KingID)
	}
	// make sure all the sim units are in the list of spritess
	for _, unit := range s.sim.GetAllUnits() {
		if s.Sprites[unit.ID.String()] == nil {
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
			s.Sprites[unit.ID.String()].EventBus = s.eventBus
			s.Sprites[unit.ID.String()].SetPosition(unit.Position)
			s.Sprites[unit.ID.String()].SetAngle(unit.MovingAngle)
			s.Sprites[unit.ID.String()].CarryingSucrose = (unit.Stats.ResourceTypeCarried == "sucrose" && unit.Stats.ResourceCarried > 0)
			s.Sprites[unit.ID.String()].CarryingWood = (unit.Stats.ResourceTypeCarried == "wood" && unit.Stats.ResourceCarried > 0)
		}
	}
	// same for buildings
	for _, building := range s.sim.GetAllBuildings() {
		if s.Sprites[building.GetID().String()] == nil {
			var spr *ui.Sprite
			switch building.GetType() {
			case sim.BuildingTypeBridge:
				spr = ui.NewBridgeSprite(building.GetID())
			case sim.BuildingTypeHive:
				spr = ui.NewHiveSprite(building.GetID())
			case sim.BuildingTypeRoachHive:
				spr = ui.NewRoachHiveSprite(building.GetID())
			case sim.BuildingTypeInConstruction:
				spr = ui.NewInConstructionSprite(building.GetID())
			}
			s.Sprites[building.GetID().String()] = spr
			s.Sprites[building.GetID().String()].SetPosition(building.GetPosition())
			s.Sprites[building.GetID().String()].ProgressBar.SetProgress(building.GetProgress())
		} else {
			s.Sprites[building.GetID().String()].SetPosition(building.GetPosition())
			s.Sprites[building.GetID().String()].ProgressBar.SetProgress(building.GetProgress())
		}
	}
	// remove building & unit sprites that are no longer in the SIM
	s.UpdateRemoveInactiveSprites()

	// Update sim before cutscenes so things happen in the world as they play.
	s.sim.Update()
	if s.CurrentNotification != nil {
		s.CurrentNotification.Update()
	}
	// HANDLE CUTSCENES - we might want sim.update though

	if s.inCutscene {
		dt := 1.0 / 60.0 // or use actual delta time
		if len(s.cutsceneActions) == 0 {
			if s.SceneCompleted {
				LevelData := NewLevelCollection().Levels[s.LevelData.LevelNumber+1]
				s.sound.Stop("msx_gamesong1")
				s.BaseScene.sm.SwitchTo(NewNarratorScene(s.fonts, s.sound, LevelData)) // switch to next level
			}
			s.inCutscene = false
			s.Ui.DrawEnabled = true
			s.drag.Enabled = true
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
		s.tutorialDialogs[0].CheckTrigger(s) // Check the first tutorial dialog trigger
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
					s.eventBus.Publish(eventing.Event{
						Type: "PlayIssueActionSFX",
					})
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
						s.eventBus.Publish(eventing.Event{
							Type: "PlayIssueActionSFX",
						})
					}
				}
			}
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
					s.eventBus.Publish(eventing.Event{
						Type: "PlayIssueActionSFX",
					})
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

func (s *PlayScene) UpdateRemoveInactiveSprites() {
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

func (s *PlayScene) Draw(screen *ebiten.Image) {
	// Draw tiles first as the BG
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(s.Ui.Camera.ViewPortZoom, s.Ui.Camera.ViewPortZoom)
	opts.GeoM.Translate(float64(s.Ui.Camera.ViewPortX), float64(s.Ui.Camera.ViewPortY))
	screen.DrawImage(s.Ui.TileMap.StaticBg, opts)

	if s.Config.DebugDraw {
		s.DebugDraw(screen)
	}

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
	if s.CurrentNotification != nil {
		s.CurrentNotification.Draw(screen)
	}

	s.Ui.Camera.DrawFade(screen) // this should always be drawn second to last

	if !s.Pause.Hidden { // pause should always come last
		s.Pause.Draw(screen)
		return
	}
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("zoom:%v", s.Ui.Camera.ViewPortZoom), 1, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mouseScreenCoords:%v,%v", mx, my), 1, 40)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mouseMapCoords:%v,%v", x, y), 1, 60)

	// print hovered tile coordinates
	if s.tileMap != nil {
		tile := s.tileMap.GetTileByPosition(x, y)
		if tile != nil {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HoveredTileCoords: %v,%v", tile.Coordinates.X, tile.Coordinates.Y), 1, 80)
		}
	}
}
