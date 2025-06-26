package tilemap

import (
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

	tileMap              *tiled.Map
	StaticBg             *ebiten.Image
	MapObjects           []*MapObject
	MapCompletionObjects []*MapCompletionObject
	TileSet              map[int]*tiled.TilesetTile
	Tiles                [][]*Tile
}

type MapObject struct {
	Rect        *image.Rectangle
	IsBuildable bool
}
type MapCompletionObject struct {
	Rect *image.Rectangle
}

func NewTilemap(mapPath string) *Tilemap {
	tm, err := tiled.LoadFile(mapPath) // this wont work in wasm! need to embed files but it breaks
	if err != nil {
		log.Fatalf("unable to load tmx: %v", err.Error())
	}

	staticBg := generateLayer0Image(tm)

	tilesIdMap := make(map[int]*tiled.TilesetTile)

	for _, tile := range tm.Tilesets[0].Tiles {
		tilesIdMap[int(tile.ID)] = tile
	}

	var mapCollisionObjects []*MapObject
	var mapCompletionObjects []*MapCompletionObject
	for _, objectGroup := range tm.ObjectGroups {
		if objectGroup.Name == "collision" {
			for _, object := range objectGroup.Objects {
				mo := &MapObject{
					Rect: &image.Rectangle{
						Min: image.Point{X: int(object.X), Y: int(object.Y)},
						Max: image.Point{X: int(object.X + object.Width), Y: int(object.Y + object.Height)},
					},
				}
				if len(object.Properties) > 0 {
					mo.IsBuildable = object.Properties.GetBool("buildable")
				} else {
					mo.IsBuildable = false
				} // default to true if no properties
				mapCollisionObjects = append(mapCollisionObjects, mo)
			}
		}
		if objectGroup.Name == "completion-area" {
			for _, object := range objectGroup.Objects {
				mo := &MapCompletionObject{
					Rect: &image.Rectangle{
						Min: image.Point{X: int(object.X), Y: int(object.Y)},
						Max: image.Point{X: int(object.X + object.Width), Y: int(object.Y + object.Height)},
					},
				}
				mapCompletionObjects = append(mapCompletionObjects, mo)
			}
		}
	}

	// for _, object := range tm.ObjectGroups[0].Objects {
	// 	mo := &MapObject{
	// 		Rect: &image.Rectangle{
	// 			Min: image.Point{X: int(object.X), Y: int(object.Y)},
	// 			Max: image.Point{X: int(object.X + object.Width), Y: int(object.Y + object.Height)},
	// 		},
	// 	}

	// 	if len(object.Properties) > 0 {
	// 		mo.IsBuildable = object.Properties.GetBool("buildable")
	// 	}
	// 	mapObjects = append(mapObjects, mo)
	// }

	tmap := &Tilemap{
		tileMap:              tm,
		StaticBg:             staticBg,
		TileSet:              tilesIdMap,
		Tiles:                make([][]*Tile, tm.Width),
		MapObjects:           mapCollisionObjects,
		MapCompletionObjects: mapCompletionObjects,
		Width:                tm.Width,
		Height:               tm.Height,
		TileSize:             tm.TileWidth,
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
				tileType = "sucrose"
			case 6:
				tileType = "wood"
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
	//fmt.Printf("GetTileByPosition X:%v, Y:%v\n", xCoord, yCoord)
	if xCoord < tm.Width && xCoord >= 0 &&
		yCoord < tm.Height && yCoord >= 0 {
		return tm.Tiles[xCoord][yCoord]
	} else {
		return nil
	}
}

func (tm *Tilemap) RemoveCollisionRect(rectToRemove *image.Rectangle) bool {
	newObjs := tm.MapObjects[:0]
	removed := false
	for _, mo := range tm.MapObjects {
		if mo.Rect.Min == rectToRemove.Min && mo.Rect.Max == rectToRemove.Max {
			// This is the one we remove
			removed = true
			continue
		}
		newObjs = append(newObjs, mo)
	}
	tm.MapObjects = newObjs
	return removed
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
