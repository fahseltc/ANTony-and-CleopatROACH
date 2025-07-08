package sim

// import (
// 	"gamejam/tilemap"
// 	"image"
// 	"math"
// )

// type PathNode struct {
// 	X, Y   int
// 	G, H   int
// 	Parent *PathNode
// }

// // Convert from world to tile coordinates
// func WorldToTile(pos image.Point, tileSize int) (int, int) {
// 	return pos.X / tileSize, pos.Y / tileSize
// }

// // Convert from tile to world coordinates (top-left corner)
// func TileToWorld(x, y, tileSize int) image.Point {
// 	return image.Point{X: x * tileSize, Y: y * tileSize}
// }

// // Is tile at x, y walkable considering other units and static obstacles
// func IsWalkable(x, y int, tm *tilemap.Tilemap, sim *T, selfID string) bool {
// 	return true
// 	if x < 0 || y < 0 || x >= tm.Width || y >= tm.Height {
// 		return false
// 	}

// 	tileRect := tm.Tiles[x][y].Rect

// 	for _, obj := range tm.MapObjects {
// 		if obj.Rect.Overlaps(*tileRect) {
// 			return false
// 		}
// 	}

// 	for _, u := range sim.GetAllUnits() {
// 		if u.ID.String() == selfID {
// 			continue
// 		}
// 		if u.Rect.Overlaps(*tileRect) {
// 			return false
// 		}
// 	}

// 	return true
// }

// // Main pathfinding function
// func FindPath(start, end image.Point, unit *Unit, sim *T) []image.Point {
// 	tm := sim.world.TileMap
// 	tileSize := tm.TileSize
// 	startX, startY := WorldToTile(start, tileSize)
// 	endX, endY := WorldToTile(end, tileSize)

// 	type key struct{ x, y int }
// 	openSet := map[key]*PathNode{}
// 	closedSet := map[key]bool{}

// 	startNode := &PathNode{X: startX, Y: startY, G: 0, H: int(math.Abs(float64(startX-endX)) + math.Abs(float64(startY-endY)))}
// 	openSet[key{startX, startY}] = startNode

// 	var current *PathNode

// 	for len(openSet) > 0 {
// 		// Get node with lowest F
// 		for _, node := range openSet {
// 			if current == nil || node.G+node.H < current.G+current.H {
// 				current = node
// 			}
// 		}
// 		if current.X == endX && current.Y == endY {
// 			break
// 		}
// 		delete(openSet, key{current.X, current.Y})
// 		closedSet[key{current.X, current.Y}] = true

// 		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// 		for _, d := range dirs {
// 			nx, ny := current.X+d[0], current.Y+d[1]
// 			k := key{nx, ny}

// 			if closedSet[k] || !IsWalkable(nx, ny, tm, sim, unit.ID.String()) {
// 				continue
// 			}

// 			gCost := current.G + 1
// 			hCost := math.Abs(float64(nx-endX)) + math.Abs(float64(ny-endY))
// 			neighbor := &PathNode{X: nx, Y: ny, G: gCost, H: int(hCost), Parent: current}

// 			if existing, ok := openSet[k]; !ok || gCost < existing.G {
// 				openSet[k] = neighbor
// 			}
// 		}
// 	}

// 	// Reconstruct path
// 	if current == nil || current.X != endX || current.Y != endY {
// 		return nil
// 	}

// 	path := []image.Point{}
// 	for current != nil {
// 		worldPos := TileToWorld(current.X, current.Y, tileSize)
// 		path = append([]image.Point{worldPos}, path...)
// 		current = current.Parent
// 	}
// 	return path
// }
