package ui

import (
	"gamejam/audio"
	"gamejam/fonts"
	"gamejam/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Pause struct {
	sound     *audio.SoundManager
	rect      *image.Rectangle
	bg        *ebiten.Image
	font      *fonts.All
	SFXSlider *Slider
	MSXSlider *Slider
	closeBtn  *Button

	Hidden bool
}

func NewPause(sound *audio.SoundManager, fonts *fonts.All) *Pause {
	rect := &image.Rectangle{Min: image.Point{X: 200, Y: 175}, Max: image.Point{X: 600, Y: 575}}
	scaled := util.ScaleImage(util.LoadImage("ui/metalPanel.png"), float32(rect.Dx()), float32(rect.Dy()))
	p := &Pause{
		sound:     sound,
		font:      fonts,
		rect:      rect,
		bg:        scaled,
		SFXSlider: NewSlider("SFX", rect.Min.X+50, rect.Min.Y+75, fonts, sound.GlobalSFXVolume),
		MSXSlider: NewSlider("Music", rect.Min.X+50, rect.Min.Y+174, fonts, sound.GlobalMSXVolume),
		Hidden:    true,
	}
	p.closeBtn = NewButton(fonts, WithText("Close"), WithRect(
		image.Rectangle{
			Min: image.Point{
				X: rect.Min.X + 100, // 400 wide
				Y: rect.Min.Y + 300, // 400 tall
			},
			Max: image.Point{
				X: rect.Min.X + 300,
				Y: rect.Min.Y + 350,
			},
		}), WithClickFunc(func() {
		p.Hidden = true
	}))

	return p
}
func (p *Pause) Update() {
	if !p.Hidden {
		p.SFXSlider.Update()
		p.MSXSlider.Update()
		p.closeBtn.Update()

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			p.sound.SetGlobalMSXVolume(p.MSXSlider.Volume)
			p.sound.GlobalSFXVolume = p.SFXSlider.Volume
			p.sound.Play("sfx_command_0")
		}
	}

}

func (p *Pause) Draw(screen *ebiten.Image) {
	if !p.Hidden {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(p.rect.Min.X), float64(p.rect.Min.Y))
		screen.DrawImage(p.bg, opts)

		p.SFXSlider.Draw(screen)
		p.MSXSlider.Draw(screen)
		p.closeBtn.Draw(screen)
	}
}
