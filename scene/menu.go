package scene

import (
	"gamejam/audio"
	"gamejam/data"
	"gamejam/fonts"
	"gamejam/ui"
	"gamejam/util"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// menuBtnSwayAmplitude is the maximum tilt (in radians) of the menu buttons.
const menuBtnSwayAmplitude = 0.06 // ~3.4 degrees

// menuAntLineCount is the number of single-file ant lines marching around the
// menu background. Tweak this to add or remove lines.
const menuAntLineCount = 6

// External links surfaced from the menu.
const (
	githubURL = "https://github.com/fahseltc/ANTony-and-CleopatROACH"
	koFiURL   = "https://ko-fi.com/fahseltc"
)

type MenuScene struct {
	BaseScene
	startBtn  *ui.Button
	optsBtn   *ui.Button
	githubBtn *ui.Button
	koFiBtn   *ui.Button
	bg        *ebiten.Image
	txt       string
	fonts     *fonts.All
	sound     *audio.SoundManager
	cfg       *data.Config
	started   bool

	antColumns []*ui.MenuAntColumn
	tick       float64

	pause *ui.Pause
}

func NewMenuScene(fonts *fonts.All, sound *audio.SoundManager, cfg *data.Config) *MenuScene {
	scene := &MenuScene{
		bg:    util.LoadImage("ui/bg/menu-bg.png"),
		txt:   "ANTony & CleopatROACH",
		fonts: fonts,
		sound: sound,
		cfg:   cfg,
		pause: ui.NewPause(sound, fonts),
	}
	scene.startBtn = ui.NewButton(fonts, ui.WithText("START"), ui.WithRect(image.Rectangle{
		Min: image.Point{X: 200, Y: 520},
		Max: image.Point{X: 390, Y: 570},
	}), ui.WithClickFunc(func() {
		levelData := NewLevelCollection().GetLevel(scene.cfg.Dev.StartingLevel)
		scene.sound.Stop("msx_menusong")
		scene.sm.SwitchTo(NewNarratorScene(scene.fonts, scene.sound, levelData))
	}))

	scene.optsBtn = ui.NewButton(fonts, ui.WithText("OPTIONS"), ui.WithRect(image.Rectangle{
		Min: image.Point{X: 410, Y: 520},
		Max: image.Point{X: 600, Y: 570},
	}), ui.WithClickFunc(func() {
		scene.pause.Hidden = false
	}))

	// Small square icon buttons flanking the main START/OPTIONS row. They open
	// the browser on desktop and in a new tab in the web build via util.OpenURL.
	// The main row spans X 200-600 at Y 520-570 (center Y 545); these sit just
	// outside it, vertically centered, and stay within the 800px screen width.
	// Rects match each image's native aspect ratio so they aren't squished.
	// github.png is 1248x285 (~4.38:1), kofi.png is 672x356 (~1.89:1). Both are
	// drawn at ~34px tall, flanking the START/OPTIONS row (center Y 545).
	// Rects match each image's native aspect ratio so they aren't squished.
	// github.png is 1248x285 (~4.38:1), kofi.png is 672x356 (~1.89:1). Both sit
	// low in the frame, below the START/OPTIONS row, with the Ko-fi icon pushed
	// well clear of OPTIONS (which ends at X 600).
	githubImg := util.LoadImage("github.png")
	scene.githubBtn = ui.NewButton(fonts,
		ui.WithRect(image.Rectangle{
			Min: image.Point{X: 40, Y: 562},
			Max: image.Point{X: 154, Y: 588}, // 114x26, ~4.38:1
		}),
		ui.WithImage(githubImg, githubImg),
		ui.WithClickFunc(func() {
			util.OpenURL(githubURL)
		}),
	)

	koFiImg := util.LoadImage("kofi.png")
	scene.koFiBtn = ui.NewButton(fonts,
		ui.WithRect(image.Rectangle{
			Min: image.Point{X: 655, Y: 543},
			Max: image.Point{X: 745, Y: 591}, // 90x48, ~1.89:1, centered in the right gap
		}),
		ui.WithImage(koFiImg, koFiImg),
		ui.WithClickFunc(func() {
			util.OpenURL(koFiURL)
		}),
	)

	// Decorative ants that march in single-file lines around the background.
	for i := 0; i < menuAntLineCount; i++ {
		length := 5 + i // vary the line lengths a little
		scene.antColumns = append(scene.antColumns, ui.NewMenuAntColumn(length, ui.ScreenWidth, ui.ScreenHeight))
	}

	return scene
}

func (s *MenuScene) Update() error {
	if !s.started {
		s.started = true
		s.sound.Play("msx_menusong")
	}
	s.tick++

	// Sway the buttons gently left and right, slightly out of phase so they
	// don't move in lockstep.
	s.startBtn.Rotation = math.Sin(s.tick*0.04) * menuBtnSwayAmplitude
	s.optsBtn.Rotation = math.Sin(s.tick*0.04+math.Pi/3) * menuBtnSwayAmplitude
	s.githubBtn.Rotation = math.Sin(s.tick*0.04+2*math.Pi/3) * 0.02
	s.koFiBtn.Rotation = math.Sin(s.tick*0.04+math.Pi) * 0.02

	for _, col := range s.antColumns {
		col.Update(ui.ScreenWidth, ui.ScreenHeight)
	}

	s.startBtn.Update(true)
	s.optsBtn.Update(true)
	s.githubBtn.Update(true)
	s.koFiBtn.Update(true)
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
	s.githubBtn.Draw(screen)
	s.koFiBtn.Draw(screen)
	s.pause.Draw(screen)

	// Ant lines march on top of everything else.
	for _, col := range s.antColumns {
		col.Draw(screen)
	}
}
