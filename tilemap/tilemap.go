package tilemap

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

// see https://pkg.go.dev/github.com/lafriks/go-tiled
type Tilemap struct {
	Width    int
	Height   int
	TileSize int

	tileMap        *tiled.Map
	StaticBg       *ebiten.Image
	CollisionRects []*image.Rectangle
	TileSet        map[int]*tiled.TilesetTile
	Tiles          [][]*Tile
}

func NewTilemap() *Tilemap {
	tm, err := tiled.LoadFile("assets/tilemap/untitled.tmx") // this wont work in wasm! need to embed files but it breaks
	if err != nil {
		log.Fatalf("unable to load tmx: %v", err.Error())
	}

	staticBg := generateLayer0Image(tm)

	tilesIdMap := make(map[int]*tiled.TilesetTile)

	for _, tile := range tm.Tilesets[0].Tiles {
		tilesIdMap[int(tile.ID)] = tile
	}

	var mapRects []*image.Rectangle
	for _, object := range tm.ObjectGroups[0].Objects {
		rect := image.Rect(int(object.X), int(object.Y), int(object.X+object.Width), int(object.Y+object.Height))
		mapRects = append(mapRects, &rect)
	}

	tmap := &Tilemap{
		tileMap:        tm,
		StaticBg:       staticBg,
		TileSet:        tilesIdMap,
		Tiles:          make([][]*Tile, tm.Width),
		CollisionRects: mapRects,
		Width:          tm.Width,
		Height:         tm.Height,
		TileSize:       tm.TileWidth,
	}
	for i := 0; i < tm.Width; i++ {
		tmap.Tiles[i] = make([]*Tile, tm.Height)
	}
	tmap.ToWorld()
	return tmap
}

func generateLayer0Image(tm *tiled.Map) *ebiten.Image {
	r, err := render.NewRenderer(tm)
	if err != nil {
		log.Fatal("unable to load tmx renderer")
	}

	r.RenderLayer(0)
	if err != nil {
		log.Fatalf("layer unsupported for rendering: %v", err.Error())
	}
	staticBg := ebiten.NewImageFromImage(r.Result)
	r.Clear()
	return staticBg
}

func (tm *Tilemap) GetMap() *tiled.Map {
	return tm.tileMap
}

func (tm *Tilemap) ToWorld() {
	mapWidth := tm.Width
	for y := 0; y < tm.Height; y++ {
		for x := 0; x < tm.Width; x++ {
			t := tm.tileMap.Layers[0].Tiles[y*mapWidth+x]
			var tileType string
			switch t.ID {
			case 15:
				tileType = "resource"
			default:
				tileType = "none"
			}
			newTile := &Tile{
				Type:        tileType,
				TileID:      int(t.ID),
				Coordinates: &image.Point{X: x, Y: y},
				Rect: &image.Rectangle{
					Min: image.Point{X: x * tm.TileSize, Y: y * tm.TileSize},
					Max: image.Point{X: (x * tm.TileSize) + tm.TileSize, Y: (y * tm.TileSize) + tm.TileSize},
				},
			}
			tm.Tiles[x][y] = newTile
		}
	}
}

func (tm *Tilemap) GetTileByPosition(x, y int) *Tile {
	xCoord := x / tm.TileSize
	yCoord := y / tm.TileSize
	fmt.Printf("GetTileByPosition X:%v, Y:%v\n", xCoord, yCoord)
	if xCoord < tm.Width && xCoord >= 0 &&
		yCoord < tm.Height && yCoord >= 0 {
		return tm.Tiles[xCoord][yCoord]
	} else {
		return nil
	}
}

func (tm *Tilemap) Render(screen *ebiten.Image) {
	// for _, tile := range tm.tileMap.Layers[0].GetTilePosition() {
	// 	if tile != nil {
	// 		//prop := tile.Tileset.Properties.Get("anything") // nil
	// 		// tileRect := tile.Tileset.GetTileRect(tile.ID) // nil ref error ' ts.Image' is nil
	// 		// tileImage := tm.StaticBgT.SubImage(tileRect).(*ebiten.Image)
	// 		// tile.Tileset.Properties.Get("passable")
	// 		// opts := &ebiten.DrawImageOptions{}
	// 		// opts.GeoM.Translate(float64(tileRect.Min.X), float64(tileRect.Min.Y))
	// 		// screen.DrawImage(tileImage, opts)
	// 	}

	// // }
}

// // Import the library
// import (
//     etiled "github.com/bird-mtn-dev/ebitengine-tiled"
// )

// // Load the xml output from tiled during the initilization of the Scene.
// // Note that OpenTileMap will attempt to load the associated tilesets and tile images
// Tilemap = etiled.OpenTileMap("assets/tilemap/base.tmx")
// // Defines the draw parameters of the tilemap tiles
// Tilemap.Zoom = 1

// // Call Update on the Tilemap during the ebitengine Update loop
// Tilemap.Update()

// // Call Draw on the Tilemap during the ebitegine Draw loop to draw all the layers in the tilemap
// Tilemap.Draw(worldScreen)

// // This loop will draw all the Object Groups in the Tilemap.
// for idx := range Tilemap.ObjectGroups {
//     Tilemap.ObjectGroups[idx].Draw(worldScreen)
// }

// // You can draw a specific Layer by calling
// Tilemap.GetLayerByName("layer1").Draw(worldScreen)

// // You can draw a specific Object Group by calling
// Tilemap.GetObjectGroupByName("ojbect group 1").Draw(worldScreen)
