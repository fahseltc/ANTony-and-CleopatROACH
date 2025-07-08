package tilemap

import (
	"gamejam/assets"
	"gamejam/vec2"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"github.com/quasilyte/pathing"
)

// see https://pkg.go.dev/github.com/lafriks/go-tiled
type Tilemap struct {
	Width    int
	Height   int
	TileSize int

	tileMap  *tiled.Map
	StaticBg *ebiten.Image

	MapObjects           []*MapObject
	MapCompletionObjects []*MapCompletionObject

	Pathing   *pathing.AStar
	pathLayer pathing.GridLayer
	PathGrid  *pathing.Grid

	TileSet map[int]*tiled.TilesetTile
	Tiles   [][]*Tile
}

type MapObject struct {
	Rect        *image.Rectangle
	IsBuildable bool
}
type MapCompletionObject struct {
	Rect *image.Rectangle
}

const (
	UnwalkableTile = iota
	WalkableTile
)

func NewTilemap(mapPath string) *Tilemap {
	tm, err := tiled.LoadFile(mapPath, tiled.WithFileSystem(assets.Files))

	if err != nil {
		log.Fatalf("unable to load tmx: %v", err.Error())
	}

	staticBg := generateLayer0Image(tm)

	tilesIdMap := make(map[int]*tiled.TilesetTile)

	for _, tile := range tm.Tilesets[0].Tiles {
		tilesIdMap[int(tile.ID)] = tile
	}

	// setup collision objects
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

	tmap := &Tilemap{
		tileMap:  tm,
		StaticBg: staticBg,

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

	// setup Pathfinding
	tmap.Pathing = pathing.NewAStar(pathing.AStarConfig{
		NumCols: uint(tm.Width),
		NumRows: uint(tm.Height),
	})
	tmap.pathLayer = pathing.MakeGridLayer([4]uint8{
		WalkableTile:   1, // passable
		UnwalkableTile: 0, // not passable
	})

	tmap.GenerateTiles()

	return tmap
}

func generateLayer0Image(tm *tiled.Map) *ebiten.Image {
	r, err := render.NewRendererWithFileSystem(tm, assets.Files)
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

func (tm *Tilemap) GenerateTiles() {
	for i := 0; i < tm.Width; i++ {
		tm.Tiles[i] = make([]*Tile, tm.Height)
	}
	newGrid := pathing.NewGrid(pathing.GridConfig{
		WorldWidth:  uint(tm.Width * tm.TileSize),
		WorldHeight: uint(tm.Height * tm.TileSize),
		CellWidth:   uint(tm.TileSize),
		CellHeight:  uint(tm.TileSize),
	})
	mapWidth := tm.Width
	for y := 0; y < tm.Height; y++ {
		for x := 0; x < tm.Width; x++ {
			t := tm.tileMap.Layers[0].Tiles[y*mapWidth+x]
			tileRect := &image.Rectangle{
				Min: image.Point{X: x * tm.TileSize, Y: y * tm.TileSize},
				Max: image.Point{X: (x * tm.TileSize) + tm.TileSize, Y: (y * tm.TileSize) + tm.TileSize},
			}
			var hasCollision bool
			for _, mo := range tm.MapObjects {
				if mo.Rect.Overlaps(*tileRect) {
					hasCollision = true
					break
				}
			}
			var tileTag uint8
			if hasCollision {
				tileTag = UnwalkableTile
			} else {
				tileTag = WalkableTile
			}
			newGrid.SetCellTile(pathing.GridCoord{X: x, Y: y}, tileTag)

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
				Type:         tileType,
				TileID:       int(t.ID),
				Coordinates:  &image.Point{X: x, Y: y},
				Rect:         tileRect,
				HasCollision: hasCollision,
			}
			tm.Tiles[x][y] = newTile
		}
	}
	tm.PathGrid = newGrid
}

func (tm *Tilemap) GetTileByPosition(x, y int) *Tile {
	xCoord := x / tm.TileSize
	yCoord := y / tm.TileSize
	return tm.GetTileByCoordinates(xCoord, yCoord)
}

func (tm *Tilemap) GetTileByCoordinates(xCoord, yCoord int) *Tile {
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
	tm.GenerateTiles()
	return removed
}

func (tm *Tilemap) AddCollisionRect(rectToAdd *image.Rectangle) bool {
	// Check if the rectangle already exists
	for _, mo := range tm.MapObjects {
		if mo.Rect.Min == rectToAdd.Min && mo.Rect.Max == rectToAdd.Max {
			return false // Already exists
		}
	}
	mo := &MapObject{
		Rect:        rectToAdd,
		IsBuildable: false,
	}
	tm.MapObjects = append(tm.MapObjects, mo)
	tm.GenerateTiles()
	return true
}

func (tm *Tilemap) FindPath(start *vec2.T, end *vec2.T) []*vec2.T {
	bpr := tm.Pathing.BuildPath(tm.PathGrid, pathing.GridCoord{X: int(start.X), Y: int(start.Y)}, pathing.GridCoord{X: int(end.X), Y: int(end.Y)}, tm.pathLayer)
	if !bpr.Partial {
		var currentPos *vec2.T
		currentPos = start
		var nav []*vec2.T
		for bpr.Steps.HasNext() {
			switch bpr.Steps.Next() {
			case pathing.DirRight:
				vec := &vec2.T{X: currentPos.X + 1, Y: currentPos.Y}
				nav = append(nav, vec)
				currentPos = vec
			case pathing.DirDown:
				vec := &vec2.T{X: currentPos.X, Y: currentPos.Y + 1}
				nav = append(nav, vec)
				currentPos = vec
			case pathing.DirLeft:
				vec := &vec2.T{X: currentPos.X - 1, Y: currentPos.Y}
				nav = append(nav, vec)
				currentPos = vec
			case pathing.DirUp:
				vec := &vec2.T{X: currentPos.X, Y: currentPos.Y - 1}
				nav = append(nav, vec)
				currentPos = vec
			default:
				break
			}
		}
		return nav
	}
	return nil
}
