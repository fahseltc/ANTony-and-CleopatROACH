package game

import (
	"gamejam/environment"
	"gamejam/scene"
	"gamejam/tilemap"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type Game struct {
	Env            *environment.Env
	LastUpdateTime time.Time

	tileMap      *tilemap.Tilemap
	sceneManager *stagehand.SceneManager[scene.GameState]
}

func NewGame(env *environment.Env) *Game {
	state := scene.GameState{}
	sceneInstance := scene.NewMenuScene()
	manager := stagehand.NewSceneManager(sceneInstance, state)

	// logger uses JSON structure as follows
	env.Logger.Info("Game Constructor", "exampleKey", "exampleValue")
	return &Game{
		Env:          env,
		tileMap:      tilemap.NewTilemapLoader(),
		sceneManager: manager,
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
	g.sceneManager.Update()

	// Pt2: Calculate DT for next loop
	g.LastUpdateTime = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.tileMap.Draw(screen)
	g.sceneManager.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.Env.Config.Get("resolution.internal.w").(int), g.Env.Config.Get("resolution.internal.h").(int)
}
