package game

import (
	"gamejam/config"
	"gamejam/fonts"
	"gamejam/log"
	"gamejam/scene"
	"log/slog"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

var fontPath = "fonts/PressStart2P-Regular.ttf"

type Game struct {
	LastUpdateTime time.Time
	sceneManager   *stagehand.SceneManager[scene.GameState]
	fonts          *fonts.All

	cfg *config.T
	log *slog.Logger
}

func New(cfg *config.T) *Game {
	state := scene.GameState{}
	fonts := fonts.Load(fontPath)
	levelData := scene.NewLevelCollection().Levels
	var manager *stagehand.SceneManager[scene.GameState]
	if cfg.SkipMenu {
		scene := scene.NewNarratorScene(fonts, levelData[cfg.StartingLevel])
		manager = stagehand.NewSceneManager(scene, state)
	} else {
		menu := scene.NewMenuScene(fonts)
		manager = stagehand.NewSceneManager(menu, state)
	}

	return &Game{
		sceneManager: manager,
		fonts:        fonts,
		cfg:          cfg,
		log:          log.NewLogger().With("for", "game"),
	}
}

func (g *Game) Update() error {
	g.sceneManager.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.tileMap.Draw(screen)
	g.sceneManager.Draw(screen)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.cfg.Resolutions.Internal.Width, g.cfg.Resolutions.Internal.Height
}
