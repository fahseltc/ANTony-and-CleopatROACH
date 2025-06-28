package scene

import (
	"gamejam/sim"
	"gamejam/ui"
	"image"

	"github.com/google/uuid"
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
			u := sim.NewDefaultAnt()
			u.SetTilePosition(6, 12)
			scene.sim.AddUnit(u)

			u2 := sim.NewDefaultAnt()
			u2.SetTilePosition(6, 5)
			scene.sim.AddUnit(u2)

			king := sim.NewRoyalAnt()
			king.SetTilePosition(12, 11)
			scene.sim.AddUnit(king)

			queen := sim.NewRoyalRoach()
			queen.SetTilePosition(28, 10)
			queen.Faction = 1
			scene.sim.AddUnit(queen)

			scene.Ui.Camera.SetZoom(ui.MinZoom)
			scene.Ui.Camera.SetPosition(10, 160)
			scene.Ui.Camera.FadeAlpha = 255

			h := sim.NewHive()
			h.SetTilePosition(8, 8)
			scene.sim.AddBuilding(h)

			scene.CompletionCondition = NewSceneCompletion(queen, king, scene.tileMap.MapCompletionObjects[0].Rect)
			return queen.ID.String(), king.ID.String()
		},
		SetupInitialCutscene: func(s *PlayScene, cleopatroach string, antony string) {
			s.inCutscene = true
			s.Ui.DrawEnabled = false
			s.drag.Enabled = false
			s.cutsceneActions = []CutsceneAction{
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
			}

			s.tutorialDialogs = []Tutorial{
				NewTutorialStep( // click and drag units
					"tutorials/tutorial-1.png",
					&image.Rectangle{Min: image.Point{X: 412, Y: 341}, Max: image.Point{X: 800, Y: 600}},
					nil, // trigger always
					func(ps *PlayScene) bool { // only complete once a unit is selected
						if len(ps.selectedUnitIDs) > 0 {
							return true
						}
						return false
					},
				),
				NewTutorialStep( // move camera
					"tutorials/tutorial-2.png",
					&image.Rectangle{Min: image.Point{X: 412, Y: 341}, Max: image.Point{X: 800, Y: 600}},
					nil,
					func(ps *PlayScene) bool { // only complete once a unit is selected
						if ps.Ui.Camera.ViewPortX != 0 && ps.Ui.Camera.ViewPortY != 0 { // TODO fragile!!
							return true
						}
						return false
					},
				),
				NewTutorialStep( // collected some sucrose + select hive
					"tutorials/tutorial-3.png",
					&image.Rectangle{Min: image.Point{X: 0, Y: 341}, Max: image.Point{X: 388, Y: 600}},
					func(ps *PlayScene) bool {
						if ps.sim.GetSucroseAmount() > 30 {
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
				NewTutorialStep( // hive selected + build unit
					"tutorials/tutorial-4.png",
					&image.Rectangle{Min: image.Point{X: 0, Y: 341}, Max: image.Point{X: 388, Y: 600}},
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
					"tutorials/tutorial-5.png",
					&image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 388, Y: 259}},
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
				NewTutorialStep( // unit selected + start building bridge
					"tutorials/tutorial-6.png",
					&image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 388, Y: 259}},
					nil,
					func(ps *PlayScene) bool {
						return ps.constructionMouse.Enabled
					},
				),
				NewTutorialStep( // info about building bridges
					"tutorials/tutorial-7.png",
					&image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 388, Y: 259}},
					nil,
					nil,
				),
				NewTutorialStep( // Build a bridge
					"tutorials/tutorial-8.png",
					&image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 388, Y: 259}},
					nil,
					func(ps *PlayScene) bool {
						for _, bld := range ps.sim.GetAllBuildings() {
							if bld.GetType() == sim.BuildingTypeInConstruction {
								return true
							}
						}
						return false
					},
				),
				NewTutorialStep( // finish the bridge
					"tutorials/tutorial-9.png",
					&image.Rectangle{Min: image.Point{X: 0, Y: 341}, Max: image.Point{X: 388, Y: 600}},
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
			}
		},
	}

	coll.Levels[1] = LevelData{
		LevelNumber: 1,
		TileMapPath: "tilemap/map2.tmx",
		LevelIntroText: `The Senate-mound murmurs with unrest -
	Some say Ant-tony hath bent his thorax too far,
	Given up tunnels and treaties for the shimmer of a roach's wing.

	But hark! The queen doth summon him from beyond the ravine again.
	A bridge must rise! Broods must hatch!
	And amid wood chips and whispers, history must crawl forward.`,
		SetupFunc: func(s *PlayScene) (string, string) {
			// hives
			h := sim.NewHive()
			h.SetTilePosition(6, 8)
			s.sim.AddBuilding(h)

			rh := sim.NewRoachHive()
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
			// starting them mining here doesnt work!
			//s.sim.IssueAction(u.ID.String(), &image.Point{X: 70, Y: 235}) // start him mining

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
		SetupInitialCutscene: func(s *PlayScene, cleopatroach string, antony string) {
			s.cutsceneActions = []CutsceneAction{
				&FadeCameraAction{Mode: "in", Speed: 2},
				// &PanCameraAction{TargetX: float64(2), TargetY: float64(4), Speed: 300},

				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: Yon queen doth beckon from beyond the ravine. But soft! I lack timber for my grand mandibleway...",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&PanCameraAction{TargetX: float64(12), TargetY: float64(4), Speed: 400},
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
				&PanCameraAction{TargetX: float64(4), TargetY: float64(4), Speed: 500},
				&IssueUnitCommandAction{
					unitID:     antony,
					targetTile: &image.Point{X: 15, Y: 12},
				},
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
	return coll
}
