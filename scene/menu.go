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
	optsBtn  *ui.Button
	bg       *ebiten.Image
	txt      string
	fonts    *fonts.All
	sound    *audio.SoundManager
	started  bool

	pause *ui.Pause
}

func NewMenuScene(fonts *fonts.All, sound *audio.SoundManager) *MenuScene {
	scene := &MenuScene{
		bg:    util.LoadImage("ui/bg/menu-bg.png"),
		txt:   "ANTony & CleopatROACH",
		fonts: fonts,
		sound: sound,
		pause: ui.NewPause(sound, *fonts),
	}
	scene.startBtn = ui.NewButton(fonts, ui.WithText("START"), ui.WithRect(image.Rectangle{
		Min: image.Point{X: 200, Y: 520},
		Max: image.Point{X: 390, Y: 570},
	}), ui.WithClickFunc(func() {
		levelData := NewLevelCollection().Levels[0]
		scene.sound.Stop("msx_menusong")
		scene.sm.SwitchTo(NewNarratorScene(scene.fonts, scene.sound, levelData))
	}))

	scene.optsBtn = ui.NewButton(fonts, ui.WithText("OPTIONS"), ui.WithRect(image.Rectangle{
		Min: image.Point{X: 410, Y: 520},
		Max: image.Point{X: 600, Y: 570},
	}), ui.WithClickFunc(func() {
		scene.pause.Hidden = false
	}))

	return scene
}

func (s *MenuScene) Update() error {
	if !s.started {
		s.started = true
		s.sound.Play("msx_menusong")
	}
	s.startBtn.Update()
	s.optsBtn.Update()
	s.pause.Update()
	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.bg, nil)
	util.DrawCenteredText(screen, s.fonts.XLarge, "ANTony", 400, 50, nil)
	util.DrawCenteredText(screen, s.fonts.XLarge, "&", 400, 120, nil)
	util.DrawCenteredText(screen, s.fonts.XLarge, "CleopatROACH", 400, 190, nil)

	s.startBtn.Draw(screen)
	s.optsBtn.Draw(screen)
	s.pause.Draw(screen)
}
