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
	coll.Levels[1] = LevelData{
		LevelNumber: 1,
		TileMapPath: "assets/tilemap/untitled.tmx",
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
		SetupInitialCutscene: func(s *PlayScene, antony string, cleopatroach string) {
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
		},
	}
	return coll
}
