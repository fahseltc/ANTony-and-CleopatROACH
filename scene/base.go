package scene

import (
	"image"

	"gamejam/fonts"

	"github.com/joelschutz/stagehand"
)

type GameState struct {
	count int
}

type BaseScene struct {
	bounds image.Rectangle
	state  *GameState
	sm     *stagehand.SceneManager[GameState]
	fonts  fonts.All
}

func (s *BaseScene) Layout(w, h int) (int, int) {
	s.bounds = image.Rect(0, 0, w, h)
	return w, h
}

func (s *BaseScene) Load(st GameState, manager stagehand.SceneController[GameState]) {
	s.state = &st
	s.sm = manager.(*stagehand.SceneManager[GameState])
}

func (s *BaseScene) Unload() GameState {
	return *s.state
}

// type FirstScene struct {
// 	BaseScene
// }

// func (s *FirstScene) Update() error {
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
// 		s.state.count++
// 	}
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
// 		s.sm.SwitchTo(&SecondScene{})
// 	}
// 	return nil
// }

// func (s *FirstScene) Draw(screen *ebiten.Image) {
// 	//screen.Fill(color.RGBA{255, 0, 0, 255}) // Fill Red
// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.state.count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
// }

// type SecondScene struct {
// 	BaseScene
// }

// func (s *SecondScene) Update() error {
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
// 		s.state.count--
// 	}
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
// 		//s.sm.SwitchWithTransition(&ThirdScene{}, stagehand.NewFadeTransition[GameState](.05))
// 	}
// 	return nil
// }

// func (s *SecondScene) Draw(screen *ebiten.Image) {
// 	//screen.Fill(color.RGBA{0, 0, 255, 255}) // Fill Blue
// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.state.count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
// }

// type ThirdScene struct {
// 	BaseScene
// }

// func (s *ThirdScene) Update() error {
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
// 		s.state.count *= 2
// 	}
// 	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
// 		//s.sm.SwitchWithTransition(&FirstScene{}, stagehand.NewSlideTransition[GameState](stagehand.RightToLeft, .05))
// 	}
// 	return nil
// }

// func (s *ThirdScene) Draw(screen *ebiten.Image) {
// 	screen.Fill(color.RGBA{0, 255, 0, 255}) // Fill Green
// 	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Count: %v, WindowSize: %s", s.state.count, s.bounds.Max), s.bounds.Dx()/2, s.bounds.Dy()/2)
// }
