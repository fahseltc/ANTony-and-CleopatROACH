package scene

import (
	"fmt"
	"gamejam/audio"
	"gamejam/data"
	"gamejam/eventing"
	"gamejam/fonts"
	"gamejam/sim"
	"gamejam/tilemap"
	"gamejam/types"
	"gamejam/ui"
	"gamejam/util"
	"image"
	"image/color"
	"math"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/pathing"
)

var PlayerFaction = 0

type PlayScene struct {
	Config *data.Config
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

	Sprites           map[string]*ui.Sprite
	sortedSprites     []*ui.Sprite
	spritesNeedReSort bool
	selectedUnitIDs   []string

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

	ActionIssuedLocation   *image.Point
	actionIssuedFrameTimer uint

	Pause *ui.Pause

	UnitGroupManager *ui.UnitGroupManager
}

func NewPlayScene(fonts *fonts.All, sound *audio.SoundManager, levelData LevelData) *PlayScene {
	config, _ := data.NewConfig()

	tileMap := tilemap.NewTilemap(levelData.TileMapPath)
	simulation := sim.New(60, tileMap)
	scene := &PlayScene{
		Config:            config,
		sound:             sound,
		LevelData:         &levelData,
		fonts:             fonts,
		sim:               simulation,
		Ui:                ui.NewUi(fonts, tileMap, simulation),
		tileMap:           tileMap,
		drag:              ui.NewDrag(),
		constructionMouse: ui.NewConstructionMouse(),
		Sprites:           make(map[string]*ui.Sprite),
		eventBus:          simulation.EventBus,
		Pause:             ui.NewPause(sound, fonts),
		UnitGroupManager:  ui.NewUnitGroupManager(fonts),
	}
	scene.eventBus.Subscribe("MakeAntButtonClickedEvent", scene.HandleMakeAntButtonClickedEvent)
	scene.eventBus.Subscribe("BuildingButtonClickedEvent", scene.HandleBuildingButtonClickedEvent)
	scene.eventBus.Subscribe("BuildClickedEvent", scene.HandleBuildClickedEvent)
	scene.eventBus.Subscribe("NotEnoughResourcesEvent", scene.NotEnoughResourcesEvent)
	scene.eventBus.Subscribe("UnitNotUnlockedEvent", scene.UnitNotUnlockedEvent)
	scene.eventBus.Subscribe("NotificationEvent", scene.NotificationEvent)

	scene.eventBus.Subscribe("SetRallyPointEvent", scene.SetRallyPointEvent)

	scene.QueenID, scene.KingID = levelData.SetupFunc(scene)

	scene.setupSFX()
	levelData.SetupInitialCutscene(scene, scene.QueenID, scene.KingID)

	if config.SkipToGameplay {
		scene.tutorialDialogs = []Tutorial{}
		scene.cutsceneActions = []CutsceneAction{}
		scene.Ui.Camera.FadeAlpha = 0
	} else {

	}

	return scene
}

func (s *PlayScene) NotificationEvent(event eventing.Event) {
	s.CurrentNotification = ui.NewNotification(&s.fonts.Med, event.Data.(eventing.NotificationEvent).Message)
}

func (s *PlayScene) SetRallyPointEvent(event eventing.Event) {
	hiveID := s.selectedUnitIDs[0]
	unitOrHiveString := s.sim.DetermineUnitOrHiveById(hiveID)
	if unitOrHiveString == "hive" {
		hive, err := s.sim.GetBuildingByID(hiveID)
		if err == nil {
			mouseWorldX, mouseWorldY := s.Ui.Camera.MousePosToMapPos()
			pt := &image.Point{X: mouseWorldX, Y: mouseWorldY}
			s.ActionIssuedLocation = pt
			hive.SetRallyPoint(pt)
		}
	}
}
func (s *PlayScene) NotEnoughResourcesEvent(event eventing.Event) {
	resName := event.Data.(eventing.NotEnoughResourcesEvent).ResourceName
	target := event.Data.(eventing.NotEnoughResourcesEvent).UnitBeingBuilt

	var str string
	if target == "Bridge" {
		str = fmt.Sprintf("Not enough %v to build %v\nOr builder is not close enough!", resName, target)
	} else {
		str = fmt.Sprintf("Not enough %v to build %v", resName, target)
	}

	s.CurrentNotification = ui.NewNotification(&s.fonts.Med, str)
}
func (s *PlayScene) UnitNotUnlockedEvent(event eventing.Event) {
	unitName := event.Data.(eventing.UnitNotUnlockedEvent).UnitName

	s.CurrentNotification = ui.NewNotification(&s.fonts.Med, fmt.Sprintf("%v Unit not unlocked yet!", unitName))
}

func (s *PlayScene) setupSFX() {
	s.eventBus.Subscribe("PlayIssueActionSFX", s.sound.PlayIssueActionSFX)
	s.eventBus.Subscribe("PlaySelectHiveSFX", s.sound.PlaySelectHiveSFX)
}
func (s *PlayScene) HandleMakeAntButtonClickedEvent(event eventing.Event) {
	if len(s.selectedUnitIDs) == 1 {
		hiveID := s.selectedUnitIDs[0]
		unitOrHiveString := s.sim.DetermineUnitOrHiveById(hiveID)
		innerEvent := event.Data.(*eventing.MakeAntButtonClickedEvent)
		unitType := innerEvent.UnitType
		if unitOrHiveString == "hive" {
			s.eventBus.Publish(eventing.Event{
				Type: "ConstructUnitEvent",
				Data: eventing.ConstructUnitEvent{
					HiveID:   hiveID,
					UnitType: unitType,
				},
			})
		}
	}
}

func (s *PlayScene) HandleBuildingButtonClickedEvent(event eventing.Event) {
	if len(s.selectedUnitIDs) >= 1 {
		innerEvent := event.Data.(eventing.BuildClickedEvent)
		unitID := s.selectedUnitIDs[0]
		unitOrHiveString := s.sim.DetermineUnitOrHiveById(unitID)
		if unitOrHiveString == "unit" {
			s.constructionMouse.Enabled = true
			s.constructionMouse.SetSprite(innerEvent.BuildingType)
			s.drag.Enabled = false
		}
	}
}
func (s *PlayScene) HandleBuildClickedEvent(event eventing.Event) {
	innerEvent := event.Data.(eventing.BuildClickedEvent)
	if len(s.selectedUnitIDs) >= 1 {
		success := s.sim.ConstructBuilding(innerEvent.TargetCoordinates, s.selectedUnitIDs[0], innerEvent.BuildingType)
		if !success {
			s.eventBus.Publish(eventing.Event{
				Type: "NotEnoughResourcesEvent",
				Data: eventing.NotEnoughResourcesEvent{ // todo: add reason why, for example "unit not close enough" etc
					ResourceName:   "Wood",
					UnitBeingBuilt: "Bridge",
				},
			})
		}
	}
	s.drag.Enabled = true
	s.constructionMouse.Enabled = false
}

func (s *PlayScene) Update() error {
	// Monitor tech unlocks
	if s.sim.GetPlayerState().TechTree.UnlockedTech[sim.TechBuildFighterUnit] {
		s.Ui.HUD.EnableFighterButton()
	}

	if !s.songStarted {
		s.songStarted = true
		s.sound.Play("msx_gamesong1")
	}
	s.sound.Update()
	s.sim.GetWorld().FogOfWar.Update(s.sim)

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

	// Every Unit in the SIM should have a sprite, if not make one.
	s.createOrUpdateUnitSprites()
	// Same for buildings
	s.createOrUpdateBuildingSprites()
	// Remove building & unit sprites that are no longer in the SIM
	s.updateRemoveInactiveSprites()

	// Handle hotkey-selected units
	hotkeyUnits := s.UnitGroupManager.Update(s.selectedUnitIDs, s.Ui.Camera, s.sim)
	if len(hotkeyUnits) != 0 {
		for _, spr := range s.Sprites {
			for _, selected := range hotkeyUnits {
				if spr.Id.String() == selected {
					spr.Selected = true
					break
				} else {
					spr.Selected = false
				}
			}
		}
	}
	// Update sim before cutscenes so things happen in the world as they play.
	s.sim.Update()
	if s.CurrentNotification != nil {
		s.CurrentNotification.Update()
	}

	// Handle cutscenes
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

	// Handle tutorials
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
			if bld.GetFaction() != uint(PlayerFaction) {
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
		if len(s.selectedUnitIDs) == 1 { // Handle 1 unit or building selected
			unitOrHiveString := s.sim.DetermineUnitOrHiveById(s.selectedUnitIDs[0])
			switch unitOrHiveString {
			case "hive":
				if s.Ui.HUD.RightSideState != ui.HiveSelectedState { // hide ui selected UI
					s.eventBus.Publish(eventing.Event{
						Type: "PlaySelectHiveSFX",
					})
					s.Ui.HUD.RightSideState = ui.HiveSelectedState
					s.constructionMouse.Enabled = false
				}
				// Handle single hive clicks
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
					mx, my := ebiten.CursorPosition()
					if !s.Ui.HUD.IsPointInside(image.Pt(mx, my)) { // if its not in the UI, set the rally point
						hive, err := s.sim.GetBuildingByID(s.selectedUnitIDs[0])
						if err != nil {
							return nil
						}
						mapX, mapY := s.Ui.Camera.ScreenPosToMapPos(mx, my)
						s.ActionIssuedLocation = &image.Point{X: mapX, Y: mapY}
						hive.SetRallyPoint(&image.Point{X: mapX, Y: mapY})
					}
				}
			case "unit":
				// hide HIVE build ui element
				if s.Ui.HUD.RightSideState != ui.UnitSelectedState { // hide hive build UI
					s.Ui.HUD.RightSideState = ui.UnitSelectedState
					s.constructionMouse.Enabled = false
				}
				// handle single unit and clicks
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) { // activate on buttonRelease to debounce
					mx, my := ebiten.CursorPosition()
					if !s.Ui.HUD.IsPointInside(image.Pt(mx, my)) {
						mapX, mapY := s.Ui.Camera.ScreenPosToMapPos(mx, my)
						s.ActionIssuedLocation = &image.Point{X: mapX, Y: mapY}
						for _, unitId := range s.selectedUnitIDs {
							s.sim.IssueAction([]string{unitId}, s.ActionIssuedLocation)
							s.eventBus.Publish(eventing.Event{
								Type: "PlayIssueActionSFX",
							})
						}
					} else if s.Ui.HUD.IsPointInsideMinimap(image.Pt(mx, my)) {
						worldX, worldY := s.Ui.MiniMap.ToWorldPixels(mx, my, s.tileMap)
						s.ActionIssuedLocation = &image.Point{X: worldX, Y: worldY}
						s.sim.IssueAction(s.selectedUnitIDs, s.ActionIssuedLocation)
						s.eventBus.Publish(eventing.Event{
							Type: "PlayIssueActionSFX",
						})
					}
				}
			default:
				s.Ui.HUD.RightSideState = ui.HiddenState
			}
		} else {
			unitOrHiveString := s.sim.DetermineUnitOrHiveById(s.selectedUnitIDs[0])
			switch unitOrHiveString {
			case "unit":
				// Handle multiple unit or building selected
				if s.Ui.HUD.RightSideState != ui.UnitSelectedState {
					s.Ui.HUD.RightSideState = ui.UnitSelectedState
					s.constructionMouse.Enabled = false
				}
				if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) { // activate on buttonRelease to debounce
					mx, my := ebiten.CursorPosition()
					if !s.Ui.HUD.IsPointInside(image.Pt(mx, my)) {
						mapX, mapY := s.Ui.Camera.ScreenPosToMapPos(mx, my)
						s.ActionIssuedLocation = &image.Point{X: mapX, Y: mapY}
						s.sim.IssueAction(s.selectedUnitIDs, s.ActionIssuedLocation)
						s.eventBus.Publish(eventing.Event{
							Type: "PlayIssueActionSFX",
						})
					} else if s.Ui.HUD.IsPointInsideMinimap(image.Pt(mx, my)) {
						worldX, worldY := s.Ui.MiniMap.ToWorldPixels(mx, my, s.tileMap)
						s.ActionIssuedLocation = &image.Point{X: worldX, Y: worldY}
						s.sim.IssueAction(s.selectedUnitIDs, s.ActionIssuedLocation)
						s.eventBus.Publish(eventing.Event{
							Type: "PlayIssueActionSFX",
						})
					}
				}
			default:
				s.Ui.HUD.RightSideState = ui.HiddenState
			}

		}
	} else {
		// zero units selected - hide the rightside HUD
		if s.Ui.HUD.RightSideState != ui.HiddenState {
			s.Ui.HUD.RightSideState = ui.HiddenState
			s.constructionMouse.Enabled = false
		}
	}

	// Set camera to minimap position pointed at
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if s.Ui.HUD.IsPointInsideMinimap(image.Pt(mx, my)) {
			worldX, worldY := s.Ui.MiniMap.ToWorldPixels(mx, my, s.tileMap)
			s.Ui.Camera.SetPosition(worldX, worldY)
		}
	}

	// Tell selected units to stop if its required
	if s.sim.ActionKeyPressed == sim.StopKeyPressed || s.sim.ActionKeyPressed == sim.HoldPositionKeyPressed {
		for _, spr := range s.Sprites {
			unit, err := s.sim.GetUnitByID(spr.Id.String())
			if err == nil && spr.Selected {
				unit.Destinations.Clear()
				if s.sim.ActionKeyPressed == sim.StopKeyPressed {
					unit.ChangeState(&sim.IdleState{})
				} else {
					//unit.Action = sim.HoldingPositionAction
				}
			}
		}
		s.sim.ActionKeyPressed = sim.NoneKeyPressed
	}
	s.drag.Update(s.Sprites, s.Ui.Camera, s.Ui.HUD)
	s.constructionMouse.Update(s.tileMap, s.sim, s.Ui.Camera)
	if !s.constructionMouse.Enabled {
		s.drag.Enabled = true
	}
	s.Ui.Update(s.sim, s.selectedUnitIDs)

	return nil
}

func (s *PlayScene) createOrUpdateUnitSprites() {
	for _, unit := range s.sim.GetAllUnits() {
		unitSprite := s.Sprites[unit.ID.String()]
		if unitSprite == nil {
			switch unit.Type {
			case types.UnitTypeDefaultAnt:
				unitSprite = ui.NewDefaultAntSprite(unit.ID)
			case types.UnitTypeFighterAnt:
				unitSprite = ui.NewFighterAntSprite(unit.ID)
			case types.UnitTypeDefaultRoach:
				unitSprite = ui.NewDefaultRoachSprite(unit.ID)
			case types.UnitTypeRoyalAnt:
				unitSprite = ui.NewRoyalAntSprite(unit.ID)
			case types.UnitTypeRoyalRoach:
				unitSprite = ui.NewRoyalRoachSprite(unit.ID)
			}
			unitSprite.EventBus = s.eventBus
			unitSprite.SetPosition(unit.Position)
			unitSprite.SetAngle(unit.MovingAngle)
			unitSprite.HealthBar.Progress = float64(unit.Stats.HPCur) / float64(unit.Stats.HPMax)

		} else { // The sprite has already been created, update it
			unitSprite.SetPosition(unit.Position)
			unitSprite.SetAngle(unit.MovingAngle)
			unitSprite.HealthBar.SetProgress(float64(unit.Stats.HPCur) / float64(unit.Stats.HPMax))
			unitSprite.CarryingSucrose = (unit.Stats.ResourceTypeCarried == types.ResourceTypeSucrose && unit.Stats.ResourcesCarried > 0)
			unitSprite.CarryingWood = (unit.Stats.ResourceTypeCarried == types.ResourceTypeWood && unit.Stats.ResourcesCarried > 0)
		}
		s.Sprites[unit.ID.String()] = unitSprite
		s.spritesNeedReSort = true
	}
}

func (s *PlayScene) createOrUpdateBuildingSprites() {
	for _, building := range s.sim.GetAllBuildings() {
		spriteBuilding := s.Sprites[building.GetID().String()]
		if spriteBuilding == nil {
			switch building.GetType() {
			case types.BuildingTypeBridge:
				spriteBuilding = ui.NewBridgeSprite(building.GetID())
			case types.BuildingTypeHive:
				spriteBuilding = ui.NewHiveSprite(building.GetID())
			case types.BuildingTypeBarracks:
				spriteBuilding = ui.NewBarracksSprite(building.GetID())
			case types.BuildingTypeRoachHive:
				spriteBuilding = ui.NewRoachHiveSprite(building.GetID())
			case types.BuildingTypeInConstruction:
				spriteBuilding = ui.NewInConstructionSprite(building.GetID())
			}
			spriteBuilding.SetPosition(building.GetPosition())

		}
		spriteBuilding.ProgressBar.SetProgress(building.GetProgress())
		spriteBuilding.HealthBar.SetProgress(float64(building.GetStats().HPCur) / float64(building.GetStats().HPMax))
		s.Sprites[building.GetID().String()] = spriteBuilding
		s.spritesNeedReSort = true
	}
}

func (s *PlayScene) updateRemoveInactiveSprites() {
	activeIDs := make(map[string]struct{})
	for _, building := range s.sim.GetAllBuildings() {
		activeIDs[building.GetID().String()] = struct{}{}
	}
	for _, unit := range s.sim.GetAllUnits() {
		activeIDs[unit.ID.String()] = struct{}{}
	}
	for id, spr := range s.Sprites {
		if spr.Type == ui.SpriteTypeStatic {
			continue // static sprites are never removed automatically - they dont exist in the SIM, just in UI
		}
		if _, exists := activeIDs[id]; !exists {
			if spr.Type == ui.SpriteTypeWorker || spr.Type == ui.SpriteTypeFighter { // if it was a unit, replace it with a blood splat
				bloodSprite := ui.NewBloodSprite(uuid.New())
				bloodSprite.SetCenteredPosition(spr.GetCenteredPosition())
				s.Sprites[bloodSprite.Id.String()] = bloodSprite
			}
			delete(s.Sprites, id) // then delete the old sprite
			s.spritesNeedReSort = true
			s.selectedUnitIDs = slices.DeleteFunc(s.selectedUnitIDs, func(id string) bool { return id == spr.Id.String() })
		}
	}
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	// Draw tilemap static BG FIRST
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(s.Ui.Camera.ViewPortZoom, s.Ui.Camera.ViewPortZoom)
	opts.GeoM.Translate(float64(s.Ui.Camera.ViewPortX), float64(s.Ui.Camera.ViewPortY))
	screen.DrawImage(s.Ui.TileMap.StaticBg, opts)

	// only sort sprites if needed
	if s.spritesNeedReSort {
		s.sortedSprites = make([]*ui.Sprite, 0, len(s.Sprites))
		for _, spr := range s.Sprites {
			s.sortedSprites = append(s.sortedSprites, spr)
		}
		slices.SortFunc(s.sortedSprites, sortSprites) // sorted by Id so we don't get flickering.
		s.spritesNeedReSort = false
	}

	// Draw all sorted sprites
	for _, sprite := range s.sortedSprites {
		sprite.Draw(screen, s.Ui.Camera)
	}

	// Then fog of war
	s.drawFogOfWar(screen)

	if s.Config.DebugDraw {
		s.DebugDraw(screen)
	}
	s.drawExpandingActionIssuedCircle(screen)
	s.Ui.Draw(screen, s.Sprites)
	s.UnitGroupManager.Draw(screen)
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

	if !s.Pause.Hidden { // this should always be drawn last
		s.Pause.Draw(screen)
		return
	}
}

func sortSprites(a, b *ui.Sprite) int {
	// First, compare by static vs non-static (static = 0, non-static = 1)
	getTypeRank := func(s *ui.Sprite) int {
		if s.Type == ui.SpriteTypeStatic {
			return 0
		}
		return 1
	}

	aRank := getTypeRank(a)
	bRank := getTypeRank(b)
	if aRank != bRank {
		if aRank < bRank {
			return -1
		}
		return 1
	}

	// If same rank, sort by Id string
	return strings.Compare(a.Id.String(), b.Id.String())
}

func (s *PlayScene) drawFogOfWar(screen *ebiten.Image) {
	fow := s.sim.GetWorld().FogOfWar
	if fow.Enabled {
		tileSize := 128.0
		for y := 0; y < fow.Height; y++ {
			for x := 0; x < fow.Width; x++ {
				worldX := float64(x) * tileSize
				worldY := float64(y) * tileSize

				// Apply camera transformation using viewPortX, viewPortY, and viewPortZoom
				screenX, screenY := s.Ui.Camera.MapPosToScreenPos(int(worldX), int(worldY))

				size := tileSize * s.Ui.Camera.ViewPortZoom

				switch fow.Tiles[y][x] {
				case sim.FogUnexplored:
					vector.DrawFilledRect(screen, float32(screenX-1), float32(screenY-1), float32(size+2), float32(size+2), color.Black, false)
				case sim.FogMemory:
					vector.DrawFilledRect(screen, float32(screenX-1), float32(screenY-1), float32(size+2), float32(size+2), color.RGBA{0, 0, 0, 128}, false)
				}
			}
		}
	}
}

func (s *PlayScene) drawExpandingActionIssuedCircle(screen *ebiten.Image) {
	if s.ActionIssuedLocation != nil && s.actionIssuedFrameTimer < 20 {
		mx, my := s.Ui.Camera.MapPosToScreenPos(s.ActionIssuedLocation.X, s.ActionIssuedLocation.Y)
		radius := 2 + int(float64(s.actionIssuedFrameTimer)*1.5)
		// Draw a simple circle using Set (not efficient, but fine for debug/notification)
		for angle := 0; angle < 360; angle++ {
			rad := float64(angle) * (3.14159265 / 180)
			x := mx + int(float64(radius)*math.Cos(rad))
			y := my + int(float64(radius)*math.Sin(rad))
			if x >= 0 && y >= 0 && x < screen.Bounds().Dx() && y < screen.Bounds().Dy() {
				screen.Set(x, y, color.RGBA{127, 255, 0, 255})
			}
		}
		s.actionIssuedFrameTimer++
		if s.actionIssuedFrameTimer >= 20 {
			s.ActionIssuedLocation = nil
			s.actionIssuedFrameTimer = 0
		}
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
	// for _, spr := range s.Sprites {
	// 	if spr.Type == ui.SpriteTypeStatic {
	// 		continue
	// 	}
	// 	rect := spr.Rect
	// 	x0, y0 := s.Ui.Camera.MapPosToScreenPos(rect.Min.X, rect.Min.Y)
	// 	x1, y1 := s.Ui.Camera.MapPosToScreenPos(rect.Max.X, rect.Max.Y)
	// 	for x := x0; x < x1; x++ {
	// 		screen.Set(x, y0, color.RGBA{255, 255, 0, 255})
	// 		screen.Set(x, y1-1, color.RGBA{255, 255, 0, 255})
	// 	}
	// 	for y := y0; y < y1; y++ {
	// 		screen.Set(x0, y, color.RGBA{255, 255, 0, 255})
	// 		screen.Set(x1-1, y, color.RGBA{255, 255, 0, 255})
	// 	}

	// 	// Draw debug circles for unit circular hitboxes
	// 	// center := spr.GetCenter()
	// 	// x0, y0 = s.Ui.Camera.MapPosToScreenPos(center.X, center.Y)
	// 	// zoomedRadius := 64 * s.Ui.Camera.ViewPortZoom
	// 	// util.DrawCircle(screen, float64(x0), float64(y0), zoomedRadius, color.RGBA{255, 25, 255, 255})
	// }
	for y := 0; y < s.tileMap.Height; y++ {
		for x := 0; x < s.tileMap.Width; x++ {
			tileType := s.tileMap.PathGrid.GetCellTile(pathing.GridCoord{X: x, Y: y})
			rect := image.Rect(x*s.tileMap.TileSize, y*s.tileMap.TileSize, (x+1)*s.tileMap.TileSize, (y+1)*s.tileMap.TileSize)
			x0, y0 := s.Ui.Camera.MapPosToScreenPos(rect.Min.X, rect.Min.Y)
			// Draw tileType uint in top-left corner
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", tileType), x0+2, y0+2)
		}
	}

	// Draw Unit destination paths
	for _, unit := range s.sim.GetAllUnits() {
		if !unit.Destinations.IsEmpty() {
			lastPos := unit.GetCenteredPosition()
			for _, dest := range unit.Destinations.Items {
				if dest == nil || lastPos == nil {
					return
				}
				x0, y0 := s.Ui.Camera.MapPosToScreenPos(int(lastPos.X), int(lastPos.Y))
				x1, y1 := s.Ui.Camera.MapPosToScreenPos(int(dest.X), int(dest.Y))
				util.DrawLine(screen, float64(x0), float64(y0), float64(x1), float64(y1), color.RGBA{0, 0, 255, 255})
				lastPos = dest
			}
		}
		// Draw unit state
		cx, cy := s.Ui.Camera.MapPosToScreenPos(unit.Position.ToPoint().X, unit.Position.ToPoint().Y)
		if unit.CurrentState != nil {
			ebitenutil.DebugPrintAt(screen, unit.CurrentState.Name(), cx, cy)
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

func (s *PlayScene) SetSelectedSprites(IDs []string) {
	for _, spr := range s.Sprites {
		spr.Selected = false
		for _, id := range IDs {
			if id == spr.Id.String() {
				spr.Selected = true
				continue
			}
		}
	}
}
