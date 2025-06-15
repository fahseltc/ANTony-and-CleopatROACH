package game

import (
	"gamejam/environment"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Env            *environment.Env
	LastUpdateTime time.Time
}

func NewGame(env *environment.Env) *Game {
	// logger uses JSON structure as follows
	env.Logger.Info("Game Constructor", "exampleKey", "exampleValue")
	return &Game{
		Env: env,
	}
}

func (g *Game) Update() error {
	// Pt1: Calculate DT
	if g.LastUpdateTime.IsZero() {
		g.LastUpdateTime = time.Now()
	}
	//dt := time.Since(g.LastUpdateTime).Seconds()

	//
	// call game object updates here
	//

	// Pt2: Calculate DT for next loop
	g.LastUpdateTime = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Env.Config.Get("resolution.internal.w").(int), g.Env.Config.Get("resolution.internal.h").(int)
}
