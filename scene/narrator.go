package scene

import (
	"gamejam/audio"
	"gamejam/fonts"
	"gamejam/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

type NarratorScene struct {
	BaseScene
	LevelData      LevelData
	sound          *audio.SoundManager
	songStarted    bool
	fonts          *fonts.All
	fullscreenText *ui.FullscreenText
	done           bool
}

func NewNarratorScene(fonts *fonts.All, sound *audio.SoundManager, levelData LevelData) *NarratorScene {
	return &NarratorScene{
		LevelData:      levelData,
		sound:          sound,
		fonts:          fonts,
		fullscreenText: ui.NewFullscreenText(fonts.Large, levelData.LevelIntroText, 2),
	}
}

func (n *NarratorScene) Update() error {
	if !n.songStarted {
		n.songStarted = true
		n.sound.Play("msx_narratorsong")
	}
	if n.done {
		return nil
	}
	n.fullscreenText.Update()
	if n.fullscreenText.IsDone() {
		n.done = true
		n.sound.Stop("msx_narratorsong")
		// Switch to the next scene, e.g., the play scene
		n.sm.SwitchTo(NewPlayScene(n.fonts, n.sound, n.LevelData))
	}
	return nil
}

func (n *NarratorScene) Draw(screen *ebiten.Image) {
	n.fullscreenText.Draw(screen)
}

func (n *NarratorScene) IsDone() bool {
	return n.done
}
