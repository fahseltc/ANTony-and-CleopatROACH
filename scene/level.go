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
		TileMapPath: "assets/tilemap/map1.tmx",
		LevelIntroText: `In the land of Nilopolis, where the sand meets sugar and the air hums with winged gossip, two empires crawl toward destiny.

		One: the mighty Ant-tonian Legion, proud builders and brave foragers. 

		The other: Queen Cleopatroach's royal roachdom, ancient, secretive, and ever-scheming.

		Long hath love fluttered betwixt Antony, soldier of soil, and Cleopatroach, goddess of grime. 

		But lo! A chasm divides them, wide as a footprint and deep as a drain. Wood must be gathered. A bridge must be built. And their love… must scuttle onward.
		

		Arise, player! Command thy swarm!`,
		SetupFunc: func(scene *PlayScene) (string, string) {
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
				&PanCameraAction{TargetX: float64(3), TargetY: float64(6), Speed: 300},
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: By mandible and might, I shall summon my swarm! To toil, my brethren!",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&IssueUnitCommandAction{
					unitID:     cleopatroach,
					targetTile: &image.Point{X: 28, Y: 10},
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
		LevelNumber:    1,
		TileMapPath:    "assets/tilemap/map2.tmx",
		LevelIntroText: "Test123 jeff add text later",
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
			s.Ui.Camera.SetPosition(10, 20)

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
				&ShowPortraitTextAreaAction{
					portraitTextArea: ui.NewPortraitTextArea(
						s.fonts,
						"Antony: Yon queen doth beckon from beyond the ravine. But soft! I lack timber for my grand mandibleway…",
						ui.PortraitTypeRoyalAnt,
					),
				},
				&PanCameraAction{TargetX: float64(33), TargetY: float64(9), Speed: 300},
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
			}
		},
		SetupCompletionCutscene: func(s *PlayScene, cleopatroach string, antony string) {

		},
	}
	return coll
}
