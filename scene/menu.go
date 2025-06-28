package scene

import (
	"gamejam/audio"
	"gamejam/fonts"
	"gamejam/ui"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScene struct {
	BaseScene
	startBtn *ui.Button
	bg       *ebiten.Image
	txt      string
	fonts    *fonts.All
	sound    *audio.SoundManager
	started  bool
}

func NewMenuScene(fonts *fonts.All, sound *audio.SoundManager) *MenuScene {
	scene := &MenuScene{
		bg:    util.LoadImage("ui/menu-bg.png"),
		txt:   "ANTony & CleopatROACH",
		fonts: fonts,
		sound: sound,
	}
	scene.startBtn = ui.NewButton(fonts.Med, ui.WithText("START"), ui.WithRect(image.Rectangle{
		Min: image.Point{X: 250, Y: 520},
		Max: image.Point{X: 550, Y: 570},
	}), ui.WithClickFunc(func() {
		levelData := NewLevelCollection().Levels[0]
		scene.sound.Stop("msx_menusong")
		scene.sm.SwitchTo(NewNarratorScene(scene.fonts, scene.sound, levelData))
	}))

	return scene
}

func (s *MenuScene) Update() error {
	if !s.started {
		s.started = true
		s.sound.Play("msx_menusong")
	}
	s.startBtn.Update()
	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.bg, nil)
	util.DrawCenteredText(screen, s.fonts.XLarge, "ANTony", 400, 50, nil)
	util.DrawCenteredText(screen, s.fonts.XLarge, "&", 400, 120, nil)
	util.DrawCenteredText(screen, s.fonts.XLarge, "CleopatROACH", 400, 190, nil)

	s.startBtn.Draw(screen)
}
