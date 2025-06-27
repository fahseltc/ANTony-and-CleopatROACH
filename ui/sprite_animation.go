package ui

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteAnimation struct {
	SpriteSheet  *ebiten.Image
	FrameWidth   int
	FrameHeight  int
	FrameCount   int
	CurrentFrame int
	FrameTime    time.Duration
	TimeElapsed  time.Duration
	Loop         bool
	Finished     bool
}

func NewSpriteAnimation(sheet *ebiten.Image, frameWidth, frameHeight, frameCount int, frameTime time.Duration, loop bool) *SpriteAnimation {
	return &SpriteAnimation{
		SpriteSheet: sheet,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
		FrameCount:  frameCount,
		FrameTime:   frameTime,
		Loop:        loop,
	}
}

func (a *SpriteAnimation) Update(dt time.Duration) {
	if a.Finished {
		return
	}
	a.TimeElapsed += dt
	if a.TimeElapsed >= a.FrameTime {
		a.TimeElapsed -= a.FrameTime
		a.CurrentFrame++
		if a.CurrentFrame >= a.FrameCount {
			if a.Loop {
				a.CurrentFrame = 0
			} else {
				a.CurrentFrame = a.FrameCount - 1
				a.Finished = true
			}
		}
	}
}

func (a *SpriteAnimation) CurrentFrameImage() *ebiten.Image {
	x := (a.CurrentFrame * a.FrameWidth)
	rect := image.Rect(x, 0, x+a.FrameWidth, a.FrameHeight)
	return a.SpriteSheet.SubImage(rect).(*ebiten.Image)
}
