package tilemap

import (
	"fmt"
	"gamejam/ui"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

// see https://pkg.go.dev/github.com/lafriks/go-tiled
type Tilemap struct {
	renderer *render.Renderer
	tileMap  *tiled.Map
	staticBg *ebiten.Image
}

func NewTilemap() *Tilemap {
	t, err := tiled.LoadFile("assets/tilemap/untitled.tmx") // this wont work in wasm! need to embed files but it breaks
	if err != nil {
		log.Fatalf("unable to load tmx: %v", err.Error())
	}

	r, err := render.NewRenderer(t)
	if err != nil {
		log.Fatal("unable to load tmx renderer")
	}

	r.RenderLayer(0)
	if err != nil {
		log.Fatalf("layer unsupported for rendering: %v", err.Error())
	}
	staticBg := ebiten.NewImageFromImage(r.Result)
	r.Clear()
	return &Tilemap{
		renderer: r,
		tileMap:  t,
		staticBg: staticBg,
	}
}

func (tm *Tilemap) GetMap() *tiled.Map {
	return tm.tileMap
}

func (tm *Tilemap) Draw(screen *ebiten.Image, camera *ui.Camera) {
	opts := &ebiten.DrawImageOptions{}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", camera.ViewPortZoom), 1, 1)

	opts.GeoM.Scale(camera.ViewPortZoom, camera.ViewPortZoom)
	opts.GeoM.Translate(float64(camera.ViewPortX), float64(camera.ViewPortY))

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", camera.ViewPortZoom), 1, 1)
	screen.DrawImage(tm.staticBg, opts)
}
