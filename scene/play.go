package scene

import (
	"fmt"
	"gamejam/audio"
	"gamejam/data"
	"gamejam/eventing"
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/sim"
	"gamejam/tilemap"
	"gamejam/types"
	"gamejam/ui"
	"gamejam/util"
	"image"
	"image/color"
	"log/slog"
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

	eventBus            *eventing.EventBus
	eventHandlerManager *EventHandlerManager
	sim                 *sim.T
	Ui                  *ui.Ui

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
	// lastLoggedCutsceneAction is the cutscene action most recently logged as
	// "started", so each action is logged once as it becomes active rather than
	// every frame it runs.
	lastLoggedCutsceneAction CutsceneAction
	// cutsceneTotalActions is the number of actions the current cutscene started
	// with. cutsceneActions is consumed (sliced) as actions finish, so this is
	// captured once at cutscene start to report "current / max" progress in logs.
	cutsceneTotalActions int

	// When true, all keyboard input (movement, hotkeys, button keys) and drag
	// selecting are ignored. Mouse clicks still work. Toggled by DisableInputAction.
	inputDisabled bool

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

	log *slog.Logger
}

func NewPlayScene(fonts *fonts.All, sound *audio.SoundManager, levelData LevelData) *PlayScene {
	config, _ := data.NewConfig()

	tileMap := tilemap.NewTilemap(levelData.TileMapPath)
	simulation := sim.New(60, tileMap)

	// Grant debug starting resources if configured. Only positive values take
	// effect; the player starts with this amount of both sucrose and wood.
	if config.Dev.DebugStartingResources > 0 {
		amount := uint(config.Dev.DebugStartingResources)
		simulation.AddResource(amount, types.ResourceTypeSucrose)
		simulation.AddResource(amount, types.ResourceTypeWood)
	}

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
		log:               log.NewLogger().With("for", "PlayScene"),
	}
	scene.eventHandlerManager = NewEventHandlerManager(simulation.EventBus, scene)

	scene.QueenID, scene.KingID = levelData.SetupFunc(scene)

	scene.setupSFX()
	levelData.SetupInitialCutscene(scene, scene.QueenID, scene.KingID)

	// Skip the intro cutscene if configured, or if jumping straight to gameplay.
	if config.Dev.SkipCutscenes || config.Dev.SkipToGameplay {
		scene.cutsceneActions = []CutsceneAction{}
		scene.inCutscene = false
		scene.Ui.DrawEnabled = true
		scene.drag.Enabled = true
		scene.inputDisabled = false
		scene.Ui.Camera.FadeAlpha = 0 // don't leave the screen faded to black
	}

	// Skip the tutorial dialogs if configured, or if jumping straight to gameplay.
	if config.Dev.SkipTutorial || config.Dev.SkipToGameplay {
		scene.tutorialDialogs = []Tutorial{}
	}

	return scene
}

func (s *PlayScene) setupSFX() {
	s.eventBus.Subscribe("PlayIssueActionSFX", s.sound.PlayIssueActionSFX)
	s.eventBus.Subscribe("PlaySelectHiveSFX", s.sound.PlaySelectHiveSFX)
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
		// If Escape just OPENED the pause menu and the active tutorial opts in
		// (only the pause tutorial does), mark it dismissed now so it's already
		// gone when the player unpauses. Other tutorials are never affected.
		if !s.Pause.Hidden && len(s.tutorialDialogs) > 0 && s.tutorialDialogs[0].WantsDismissOnPauseOpen() {
			s.tutorialDialogs[0].Dismiss()
		}
	}
	if !s.Pause.Hidden { // stop the game processing when paused!
		s.Pause.Update()
		return nil
	}

	if s.CompletionCondition.IsComplete(s.sim) && !s.SceneCompleted {
		s.SceneCompleted = true
		// Log the transition exactly once (guarded by SceneCompleted). The
		// predicate itself stays side-effect free so it can be polled every
		// frame without spamming the log during the completion cutscene.
		s.log.Debug("scene completion condition met")
		s.LevelData.SetupCompletionCutscene(s, s.QueenID, s.KingID)
	}

	// Every Unit in the SIM should have a sprite, if not make one.
	s.createOrUpdateUnitSprites()
	// Same for buildings
	s.createOrUpdateBuildingSprites()
	// Remove building & unit sprites that are no longer in the SIM
	s.updateRemoveInactiveSprites()

	// Handle hotkey-selected units (keyboard driven, so skip when input disabled)
	if !s.inputDisabled {
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
				s.sound.Stop("msx_gamesong1")
				nextLevelNum := s.LevelData.LevelNumber + 1
				if nextLevel, ok := NewLevelCollection().Levels[nextLevelNum]; ok {
					// There's another level: play its intro narration.
					s.BaseScene.sm.SwitchTo(NewNarratorScene(s.fonts, s.sound, nextLevel))
				} else {
					// Last level finished (no level nextLevelNum in the map).
					// Return to the main menu instead of loading an empty level.
					s.log.Info("no next level; returning to main menu", "finishedLevel", s.LevelData.LevelNumber)
					s.BaseScene.sm.SwitchTo(NewMenuScene(s.fonts, s.sound, s.Config))
				}
			}
			s.inCutscene = false
			s.cutsceneTotalActions = 0 // reset so the next cutscene recaptures its total
			s.Ui.DrawEnabled = true
			s.drag.Enabled = true
			s.inputDisabled = false // restore input when the cutscene finishes
			// The click that dismissed the final cutscene dialog was pressed
			// during the cutscene; its mouse-release lands on this first
			// post-cutscene frame. Return early so drag.Update doesn't treat
			// that dangling release as a drag-select (which would otherwise
			// select a phantom unit and trip tutorial completion checks).
			return nil
		} else {
			// Capture the cutscene's total action count once, at its start.
			// cutsceneActions is sliced down as actions finish, so len() alone
			// can't report the full length after the first action completes.
			if s.cutsceneTotalActions == 0 {
				s.cutsceneTotalActions = len(s.cutsceneActions)
			}
			// current is the 1-based index of the active action within the whole
			// cutscene: total minus the not-yet-started remainder.
			current := s.cutsceneTotalActions - len(s.cutsceneActions) + 1

			currentCutScene := s.cutsceneActions[0]
			// Log each cutscene action once, when it becomes the active action,
			// routing it through the JSON logging system like other components.
			if currentCutScene != s.lastLoggedCutsceneAction {
				s.log.Info("cutscene action started",
					"action", fmt.Sprintf("%T", currentCutScene),
					"current", current,
					"max", s.cutsceneTotalActions)
				s.lastLoggedCutsceneAction = currentCutScene
			}
			if s.currentDialog != nil {
				s.currentDialog.Update()
			}
			if currentCutScene.Update(s, dt) {
				s.log.Info("cutscene action finished",
					"action", fmt.Sprintf("%T", currentCutScene),
					"current", current,
					"max", s.cutsceneTotalActions)
				s.cutsceneActions = s.cutsceneActions[1:]
				s.lastLoggedCutsceneAction = nil
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

	if !s.inCutscene {
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
	}

	if !s.inCutscene {
		if len(s.selectedUnitIDs) > 0 {
			if len(s.selectedUnitIDs) == 1 { // Handle 1 unit or building selected
				unitOrHiveString := s.sim.DetermineUnitOrHiveById(s.selectedUnitIDs[0])
				switch unitOrHiveString {
				case "hive":
					// Ant and roach hives share selection behaviour but show
					// different build panels (roach panel produces roaches).
					hiveState := ui.HiveSelectedState
					if hive, err := s.sim.GetBuildingByID(s.selectedUnitIDs[0]); err == nil &&
						hive.GetType() == types.BuildingTypeRoachHive {
						hiveState = ui.RoachHiveSelectedState
					}
					if s.Ui.HUD.RightSideState != hiveState { // hide ui selected UI
						s.eventBus.Publish(eventing.Event{
							Type: "PlaySelectHiveSFX",
						})
						s.Ui.HUD.RightSideState = hiveState
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
	// Disable drag selecting while input is locked.
	if s.inputDisabled {
		s.drag.Enabled = false
	}
	if !s.inCutscene {
		s.drag.Update(s.Sprites, s.Ui.Camera, s.Ui.HUD)
		s.constructionMouse.Update(s.tileMap, s.sim, s.Ui.Camera)
	}
	if !s.constructionMouse.Enabled && !s.inputDisabled {
		s.drag.Enabled = true
	}
	// Gate keyboard-driven camera panning/zoom and button hotkeys. Mouse clicks
	// on buttons still work because those aren't gated by this flag.
	s.Ui.Camera.InputEnabled = !s.inputDisabled
	s.Ui.Update(s.sim, s.selectedUnitIDs, !s.inputDisabled)

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
			case types.BuildingTypeAntHive:
				spriteBuilding = ui.NewAntHiveSprite(building.GetID())
			case types.BuildingTypeBarracks:
				spriteBuilding = ui.NewBarracksSprite(building.GetID())
			case types.BuildingTypeRoachHive:
				spriteBuilding = ui.NewRoachHiveSprite(building.GetID())
			case types.BuildingTypeInConstruction:
				spriteBuilding = ui.NewInConstructionSprite(building.GetID())
			}
			spriteBuilding.SetPosition(building.GetPosition())
			// originalBuildingPos := building.GetPosition()
			// spriteBuilding.SetPosition(&vec2.T{X: originalBuildingPos.X / 128.0, Y: originalBuildingPos.Y / 128.0})

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

	if s.Config.Dev.DebugDraw {
		s.DebugDraw(screen)
	}
	s.drawExpandingActionIssuedCircle(screen)
	s.Ui.Draw(screen, s.Sprites)
	// The unit-hotkey group buttons are part of the HUD, so hide them whenever
	// the rest of the UI is hidden (e.g. during cutscenes). s.Ui.Draw already
	// gates itself on DrawEnabled internally, but this panel is drawn separately.
	if s.Ui.DrawEnabled {
		s.UnitGroupManager.Draw(screen)
	}
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
			ebitenutil.DebugPrintAt(screen, unit.CurrentState.GetName(), cx, cy)
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
