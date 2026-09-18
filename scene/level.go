package scene

import (
	"gamejam/sim"
	"gamejam/types"
	"gamejam/ui"
	"image"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
)

type LevelData struct {
	LevelNumber             int
	TileMapPath             string
	LevelIntroText          string
	SetupFunc               func(*PlayScene) (queenID string, kingID string)
	SetupInitialCutscene    func(*PlayScene, string, string)
	SetupCompletionCutscene func(*PlayScene, string, string)
}

type LevelCollection struct {
	Levels map[int]LevelData
}

// GetLevel returns the level data for the given level number, falling back to
// level 0 when the number is out of range. This guards against an invalid
// dev.startingLevel in config producing an empty LevelData.
func (c *LevelCollection) GetLevel(n int) LevelData {
	if ld, ok := c.Levels[n]; ok {
		return ld
	}
	return c.Levels[0]
}

func NewLevelCollection() *LevelCollection {
	coll := &LevelCollection{
		Levels: make(map[int]LevelData),
	}
	coll.Levels[0] = LevelData{
		LevelNumber: 0,
		TileMapPath: "tilemap/map1.tmx",
		LevelIntroText: `In the land of Nilopolis, where the sand meets sugar and the air hums with winged gossip, two empires crawl toward destiny.

		One: the mighty Ant-tonian Legion, proud builders and brave foragers. 

		The other: Queen Cleopatroach's royal roachdom, ancient, secretive, and ever-scheming.

		Long hath love fluttered betwixt Antony, soldier of soil, and Cleopatroach, goddess of grime. 

		But lo! A chasm divides them, wide as a footprint and deep as a drain. Wood must be gathered. A bridge must be built. And their love… must scuttle onward.
		

		Arise, player! Command thy swarm!`,
		SetupFunc: func(scene *PlayScene) (string, string) {
			u := sim.NewDefaultAntWithTilePosition(5, 9)
			scene.sim.AddUnit(u)

			u2 := sim.NewDefaultAntWithTilePosition(6, 4)
			scene.sim.AddUnit(u2)

			king := sim.NewRoyalAnt()
			king.SetTilePosition(12, 4)
			scene.sim.AddUnit(king)

			queen := sim.NewRoyalRoach()
			queen.SetTilePosition(28, 10)
			scene.sim.AddUnit(queen)

			scene.Ui.Camera.SetZoom(ui.MinZoom)
			scene.Ui.Camera.SetPosition(10, 160)
			scene.Ui.Camera.FadeAlpha = 255

			h := sim.GetBuildingInstance(types.BuildingTypeAntHive, uint(PlayerFaction))
			h.SetTilePosition(6, 6)
			scene.sim.AddBuilding(h)

			scene.CompletionCondition = NewSceneCompletion(queen, king, scene.tileMap.MapCompletionObjects[0].Rect)
			return queen.ID.String(), king.ID.String()
		},
		SetupInitialCutscene: func(s *PlayScene, cleopatroach string, antony string) {
			s.inCutscene = true
			s.Ui.DrawEnabled = false
			s.drag.Enabled = false
			s.cutsceneActions = []CutsceneAction{
				&DisableInputAction{},
				&NudgeCameraByTilesAction{OffsetTilesX: 0, OffsetTilesY: -4, Speed: 300},
				&RevealFogOfWarAction{
					TopLeft:     &image.Point{X: 14, Y: 4},
					BottomRight: &image.Point{X: 32, Y: 18},
				},
				&FadeCameraAction{Mode: "in", Speed: 2},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: O brave new bugworld! Where art thou, my chitinous queen? I must construct yon bridge, ere my love is lost!",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 18, Y: 10},
				},
				&PanCameraAction{TargetX: float64(27), TargetY: float64(10), Speed: 300},
				&IssueUnitCommandAction{
					unitID:     cleopatroach,
					targetTile: &image.Point{X: 26, Y: 10},
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach:  Love that is count'd is love too small. Rescue me, my six-legged soldier!",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&PanCameraAction{TargetX: float64(3), TargetY: float64(2), Speed: 300},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: By mandible and might, I shall summon my swarm! To toil, my brethren! Reap the crystal'd sweet!",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&IssueUnitCommandAction{
					unitID:     cleopatroach,
					targetTile: &image.Point{X: 28, Y: 10},
				},
				&EnableInputAction{},
			}
			// This should be two tutorials, one about selecting units and another about sending them to harvest the two resource types.
			s.tutorialDialogs = []Tutorial{
				NewTutorialStep( // 'left click and drag a box around units. Send the units to harvest with right click'
					"tutorials/level1/tutorial-1.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil, // trigger always
					func(ps *PlayScene) bool { // only complete once a selected WORKER is returning resources (delivering)
						for _, id := range ps.selectedUnitIDs {
							if unit, err := ps.sim.GetUnitByID(id); err == nil &&
								unit.IsWorker() &&
								unit.CurrentState != nil &&
								unit.CurrentState.GetName() == sim.UnitStateDelivering.ToString() {
								return true
							}
						}
						return false
					},
				),
				NewTutorialStep( // 'use wasd and mouse wheel to move camera'
					"tutorials/level1/tutorial-2.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil,
					func(ps *PlayScene) bool {
						return ps.Ui.Camera.PlayerMovedCameraThisFrame()
					},
				),
				NewPauseTutorialStep( // 'you can press esc to pause the game'
					"tutorials/level1/tutorial-pause.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil,
					nil,
				),
				// KEEP HARVESTING TUTORIAL??
				NewTutorialStep( // 'once youve gathered some sucrose, click the nearby hive'
					"tutorials/level1/tutorial-3.png",
					TutorialRightCenter, TutorialSizeTiny,
					func(ps *PlayScene) bool {
						if ps.sim.GetSucroseAmount() >= 50 {
							return true
						}
						return false
					},
					func(ps *PlayScene) bool {
						for _, id := range ps.selectedUnitIDs {
							if ps.sim.DetermineUnitOrHiveById(id) == "hive" {
								return true
							}
						}
						return false
					},
				),
				NewTutorialStep( // 'With the hive selected, press Q or click ICON to create a new ant for 50 sucrose'
					"tutorials/level1/tutorial-4.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil,
					func(ps *PlayScene) bool {
						for _, bld := range ps.sim.GetAllBuildings() {
							if bld.GetProgress() != 0 {
								return true
							}
						}
						return false
					},
				),
				NewTutorialStep( // wood collected + select single unit
					"tutorials/level1/tutorial-5.png",
					TutorialRightCenter, TutorialSizeTiny,
					func(ps *PlayScene) bool {
						if ps.sim.GetWoodAmount() > 30 {
							return true
						}
						return false
					},
					func(ps *PlayScene) bool {
						if len(ps.selectedUnitIDs) == 1 {
							if ps.sim.DetermineUnitOrHiveById(ps.selectedUnitIDs[0]) == "unit" {
								return true
							}
						}
						return false
					},
				),
				NewTutorialStep( // ' with the unit selected, press G then Q or ICON to begin a new bridge construction for 50 WOOD'
					"tutorials/level1/tutorial-6.png",
					TutorialTopLeft, TutorialSizeTiny,
					nil,
					func(ps *PlayScene) bool {
						return ps.constructionMouse.Enabled
					},
				),
				NewTutorialStepNoClick( // Build a bridge - this is the placement step.
					// Do NOT dismiss on the left-click that places the bridge
					// blueprint; complete once an in-construction building exists.
					"tutorials/level1/tutorial-7.png",
					TutorialTopLeft, TutorialSizeTiny,
					nil,
					func(ps *PlayScene) bool {
						for _, bld := range ps.sim.GetAllBuildings() {
							if bld.GetType() == types.BuildingTypeInConstruction {
								return true
							}
						}
						return false
					},
				),
				NewTutorialStep( // finish the bridge to re-unite them -
					"tutorials/level1/tutorial-8.png",
					TutorialCenter, TutorialSizeSmall,
					nil,
					nil,
				),
			}
		},

		SetupCompletionCutscene: func(s *PlayScene, cleopatroach string, antony string) {
			s.inCutscene = true
			s.inCutscene = true
			s.Ui.DrawEnabled = false
			s.drag.Enabled = false
			s.selectedUnitIDs = []string{} // clear selected unit IDs

			s.cutsceneActions = []CutsceneAction{
				&DisableInputAction{},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 27, Y: 10},
				},
				&WaitAction{
					Duration: 1.0,
				},
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
						"Antony: Fear not, thorax of my heart! I have crushed the peril beneath my heel",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach: Come hither, sweet thorax. Let us entwine our antennae in triumph.",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&FadeCameraAction{Mode: "out", Speed: 1},
				&WaitAction{
					Duration: 1.0, // wait for 1 second
				},
			}
			s.tutorialDialogs = []Tutorial{
				NewTutorialStep( // goal of level
					"tutorials/lvl2-tutorial-1.png",
					TutorialBottomLeft, TutorialSizeMedium,
					nil,
					nil,
				),
			}
		},
	}

	// LEVEL 1
	// LEVEL 1
	// LEVEL 1
	coll.Levels[1] = LevelData{
		LevelNumber: 1,
		TileMapPath: "tilemap/map2.tmx",
		LevelIntroText: `The Senate-mound murmurs with unrest -
	Some say Ant-tony hath bent his thorax too far,
	Given up tunnels and treaties for the shimmer of a roach's wing.

	But hark! The queen doth summon him from beyond the ravine again.
	A bridge must rise! Broods must hatch!
	And amid wood chips and whispers, history must crawl forward.`,
		// LEVEL 1 SETUP FUNC
		SetupFunc: func(s *PlayScene) (string, string) {
			// hives
			h := sim.GetBuildingInstance(types.BuildingTypeAntHive, uint(PlayerFaction))
			h.SetTilePosition(6, 8)
			s.sim.AddBuilding(h)

			rh := sim.GetBuildingInstance(types.BuildingTypeRoachHive, uint(PlayerFaction))
			rh.SetTilePosition(40, 12)
			s.sim.AddBuilding(rh)

			// royalty
			queen := sim.NewRoyalRoach()
			queen.SetTilePosition(33, 9)
			s.sim.AddUnit(queen)

			king := sim.NewRoyalAnt()
			king.SetTilePosition(10, 10)
			s.sim.AddUnit(king)

			// regular guys
			u := sim.NewDefaultAnt()
			u.SetTilePosition(4, 7)
			s.sim.AddUnit(u)
			// Start this worker harvesting a chosen resource tile right away.
			// IssueHarvestTile takes TILE coordinates (not pixels) and puts the
			// unit into its HarvestingState; the sim ticks during the intro
			// cutscene, so it walks over and begins gathering as the scene plays.
			// (3,4) is a sucrose tile in the top-left patch on map2.tmx; the
			// worker approaches it from the adjacent walkable tile.
			s.sim.IssueHarvestTile(u.ID.String(), 3, 4)

			// regular roach
			u2 := sim.NewDefaultRoach()
			u2.SetTilePosition(42, 9)
			s.sim.AddUnit(u2)
			s.sim.IssueHarvestTile(u2.ID.String(), 46, 11)

			s.Ui.Camera.SetZoom(ui.MinZoom)
			s.Ui.Camera.SetPosition(0, 105)
			s.Ui.Camera.FadeAlpha = 255

			s.inCutscene = true
			s.Ui.DrawEnabled = false
			s.drag.Enabled = false
			s.constructionMouse.Enabled = false

			s.CompletionCondition = NewSceneCompletion(queen, king, s.tileMap.MapCompletionObjects[0].Rect)
			return queen.ID.String(), king.ID.String()
		},
		// LEVEL 1 CUTSCENE
		SetupInitialCutscene: func(s *PlayScene, cleopatroach string, antony string) {
			s.cutsceneActions = []CutsceneAction{
				&DisableInputAction{},
				&PanCameraAction{TargetX: float64(4), TargetY: float64(4), Speed: 1500},
				&FadeCameraAction{Mode: "in", Speed: 3},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 14, Y: 11},
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: Yon queen doth beckon from beyond the ravine. But soft! I lack timber for my grand mandibleway...",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&PanCameraAction{TargetX: float64(12), TargetY: float64(4), Speed: 800},
				&IssueUnitCommandAction{
					unitID:     cleopatroach,
					targetTile: &image.Point{X: 31, Y: 9},
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach: The planks lie here, my love! But in return, thou must aid me in raising our mighty brood!",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&PanCameraAction{TargetX: float64(5), TargetY: float64(7), Speed: 800},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: Come, Cleopatroach, my thorax burns for thee - Let us entwine where petals crown the dirt, ",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: In yonder ring where daisies dare to bloom. ",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&PanCameraAction{TargetX: float64(5), TargetY: float64(7), Speed: 800},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 16, Y: 13},
				},
				&RevealFogOfWarAction{
					TopLeft:     &image.Point{X: 12, Y: 10},
					BottomRight: &image.Point{X: 20, Y: 25},
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: There shall we clasp antennae, love, and fate, And make a kingdom of that perfumed ground.",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 9, Y: 9},
				},
				&PanCameraAction{TargetX: float64(1), TargetY: float64(1), Speed: 300},
				&EnableInputAction{},
			}
			// LEVEL 1 TUTORIAL
			s.tutorialDialogs = []Tutorial{
				NewTutorialStep( // Get both to the flower circle
					"tutorials/level2/tutorial-1.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil, // trigger always
					nil, // only progress via left click
				),
				NewTutorialStep( // Manage two different hives
					"tutorials/level2/tutorial-2.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil, // trigger always
					nil, // only progress via left click
				),
				NewTutorialStepNoClick( // select the first hive
					"tutorials/level2/tutorial-3.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil, // trigger always
					// Complete once the ant hive is the selected building. The
					// placement/selection click must not dismiss the modal, so
					// this is a no-click step gated on the selection instead.
					func(ps *PlayScene) bool {
						return isHiveTypeSelected(ps, types.BuildingTypeAntHive)
					},
				),
				NewTutorialStepNoClick( // press control + 1 to set the hotkey
					"tutorials/level2/tutorial-4.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil, // trigger always
					// Complete once a control group has been bound to hotkey 1.
					func(ps *PlayScene) bool {
						return ps.UnitGroupManager.HasGroup(ebiten.Key1)
					},
				),
				NewTutorialStepNoClick( // scroll to roach hive and select it
					"tutorials/level2/tutorial-5.png",
					TutorialTopLeft, TutorialSizeTiny,
					nil, // trigger always
					// Complete once the roach hive is the selected building.
					func(ps *PlayScene) bool {
						return isHiveTypeSelected(ps, types.BuildingTypeRoachHive)
					},
				),
				NewTutorialStepNoClick( // press control + 2 to set the hotkey
					"tutorials/level2/tutorial-6.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil, // trigger always
					// Complete once a control group has been bound to hotkey 2.
					func(ps *PlayScene) bool {
						return ps.UnitGroupManager.HasGroup(ebiten.Key2)
					},
				),
				NewTutorialStepNoClick( // press 1 twice to move the camera back to the first hive
					"tutorials/level2/tutorial-7.png",
					TutorialRightCenter, TutorialSizeTiny,
					nil, // trigger always
					// Complete once the player double-taps hotkey 1, which jumps
					// the camera back to the first (ant hive) group.
					func(ps *PlayScene) bool {
						return ps.UnitGroupManager.DidRecenterOnGroup(ebiten.Key1)
					},
				),
			}
		},

		SetupCompletionCutscene: func(s *PlayScene, cleopatroach string, antony string) {
			s.inCutscene = true
			s.inCutscene = true
			s.Ui.DrawEnabled = false
			s.drag.Enabled = false

			s.selectedUnitIDs = []string{} // clear selected unit IDs

			s.cutsceneActions = []CutsceneAction{
				&IssueUnitCommandAction{
					unitID:     cleopatroach,
					targetTile: &image.Point{X: 18, Y: 15},
				},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 14, Y: 11},
				},
				&WaitAction{
					Duration: 1.0, // wait for 1 second
				},
				&IssueUnitCommandAction{
					unitID:     cleopatroach,
					targetTile: &image.Point{X: 17, Y: 14},
				},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 15, Y: 12},
				},
				&WaitAction{
					Duration: 0.5, // wait for 1 second
				},
				&PanCameraAction{TargetX: float64(4), TargetY: float64(4), Speed: 500},
				&ZoomCameraAction{
					TargetZoom: 0.8,
					Speed:      1,
					FocusX:     17 * 128,
					FocusY:     15 * 128,
				},
				&DrawTemporarySpriteAction{
					spr:            ui.NewHeartSprite(uuid.New()),
					TargetPosition: &image.Point{X: 2104, Y: 1724},
					MaxDuration:    180,
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: Sweet Cleopatroach, beneath these perfumed petals we meet, Yet even in this bloom,",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: the thorn of Rome doth prick my side. Octavian's shadow crawls o'er all our kingdoms vast,",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: His claws poised to snatch the crown from humble thorax and wing alike",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach: Antony, my lord, the Emperor Bugustus's gaze is cold and cruel,",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach: His legions swarm the sands, his whispers poison the air.",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach: Let us bind our broods, that none may sunder this fragile alliance.",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach: Then let the courts of Bugustus tremble and the senate-mounds whisper,",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Cleopatroach: For love, like the smallest insect, can move mountains and topple thrones.",
						ui.PortraitTypeRoyalRoach,
					),
				},
				&FadeCameraAction{Mode: "out", Speed: 1},
			}
		},
	}
	coll.Levels[2] = LevelData{
		LevelNumber: 2,
		TileMapPath: "tilemap/map3.tmx",
		LevelIntroText: `Thanks for playing the demo of ANTony & CleopatROACH! It was created for the Ebitengine Game Jam 2025, and is a work in progress.
		
		I wanted to add much more - combat, more levels, more story, more shakespeare puns (Enobarkbug!) and more features - but ran out of time in the two weeks alotted.
		
		I appreciate you playing this demo, and hope you enjoyed it!
		
		CREDITS:
		
		PROGRAMMING & EVERYTHING ELSE:
		Charles Fahselt
		
		GOLANG CONSULTANT:
		Medge

		SHAKESPEARE CONSULTANT:
		Chez Oxendine

		ART:
		ChatGPT (and I did a little bit myself)
		`,
		SetupFunc: func(s *PlayScene) (string, string) {
			s.Ui.Camera.SetZoom(ui.MinZoom)
			return "", ""
		},
		SetupInitialCutscene:    func(s *PlayScene, cleopatroach string, antony string) {},
		SetupCompletionCutscene: func(s *PlayScene, cleopatroach string, antony string) {},
	}

	// coll.Levels[3] = LevelData{
	// 	LevelNumber:    3,
	// 	TileMapPath:    "tilemap/test-map.tmx",
	// 	LevelIntroText: "",
	// 	SetupFunc: func(s *PlayScene) (string, string) {
	// 		u := sim.NewDefaultAnt()
	// 		u.SetTilePosition(9, 0)
	// 		s.sim.AddUnit(u)

	// 		// for i := 0; i < 5; i++ {
	// 		// 	u := sim.NewDefaultAnt()
	// 		// 	u.SetTilePosition(9, i)
	// 		// 	s.sim.AddUnit(u)
	// 		// }

	// 		// for i := 0; i < 5; i++ {
	// 		// 	u := sim.NewDefaultAnt()
	// 		// 	u.SetTilePosition(5, i)
	// 		// 	s.sim.AddUnit(u)
	// 		// }

	// 		// king := sim.NewRoyalAnt()
	// 		// king.SetTilePosition(12, 11)
	// 		// s.sim.AddUnit(king)

	// 		// queen := sim.NewRoyalRoach()
	// 		// queen.SetTilePosition(28, 10)
	// 		// queen.Faction = 1
	// 		// s.sim.AddUnit(queen)

	// 		s.Ui.Camera.SetZoom(ui.MinZoom)
	// 		s.Ui.Camera.SetPosition(0, 0)
	// 		s.Ui.Camera.FadeAlpha = 255

	// 		h := sim.GetBuildingInstance(types.BuildingTypeAntHive, uint(PlayerFaction))
	// 		h.SetTilePosition(6, 7)
	// 		s.sim.AddBuilding(h)

	// 		// //bad guys
	// 		// for i := 0; i < 5; i++ {
	// 		// 	u := sim.NewDefaultAnt()
	// 		// 	u.Faction = 1
	// 		// 	u.SetTilePosition(21, 14+i)
	// 		// 	s.sim.AddUnit(u)
	// 		// }
	// 		// for i := 0; i < 15; i++ {
	// 		// 	u := sim.NewDefaultAnt()
	// 		// 	u.Faction = 1
	// 		// 	u.SetTilePosition(23, 14+i)
	// 		// 	s.sim.AddUnit(u)
	// 		// }

	// 		return "", ""
	// 	},
	// 	SetupInitialCutscene:    func(s *PlayScene, cleopatroach string, antony string) {},
	// 	SetupCompletionCutscene: func(s *PlayScene, cleopatroach string, antony string) {},
	// }

	return coll
}

// isHiveTypeSelected reports whether the player's current selection contains a
// building of the given hive type. Used by the level-2 tutorial to detect when
// the ant hive or roach hive has been selected. (sim.DetermineUnitOrHiveById
// only recognises the ant hive, so we inspect the building type directly here.)
func isHiveTypeSelected(ps *PlayScene, hiveType types.Building) bool {
	for _, id := range ps.selectedUnitIDs {
		if bld, err := ps.sim.GetBuildingByID(id); err == nil && bld.GetType() == hiveType {
			return true
		}
	}
	return false
}
