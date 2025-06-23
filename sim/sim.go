package sim

import (
	"fmt"
	"gamejam/tilemap"
	"image"
	"slices"
	"sync"
	"sync/atomic"
)

var NearbyDistance = uint(300)

// T is the simulated world for the game. This Sim makes
// a couple assumptions:
//
//   - The Game ticks at a fixed rate (tps)
//   - The time delta between ticks is (mostly) fixed
type T struct {
	tps   int
	dt    float64
	world *World

	// ID generator
	unitIdx atomic.Uint32

	playerState PlayerState
	stateMu     sync.RWMutex

	playerSpawnX, playerSpawnY float64
	playerUnits                []*Unit
	playerBuildings            []*Hive
	enemyUnits                 []*Unit
	enemySpawnX, enemySpawnY   float64

	selectedUnits []*Unit
}

type World struct {
	TileMap          *tilemap.Tilemap
	TileData         [][]*tilemap.Tile
	CollisionObjects []*image.Rectangle
}

const (
	NoneTile     = "none"
	ResourceTile = "resource"
	BuildingTile = "building"
)

type PlayerState struct {
	Wood  uint16
	Food  uint16
	Water uint16
}

func (s *T) GetPlayerState() PlayerState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.playerState
}

const (
	Protein = "protein"
	Sucrose = "sucrose"
)

// Resource represents a collectable resource on the map
type Resource struct {
	Type      string
	X, Y      float64
	Available uint
}

func New(tps int, tileMap *tilemap.Tilemap) *T {
	return &T{
		tps: tps,
		dt:  float64(1 / tps),
		world: &World{
			TileMap:          tileMap,
			TileData:         tileMap.Tiles,
			CollisionObjects: tileMap.CollisionRects,
		},

		// TODO Spawn Points
		playerState: PlayerState{},
		//playerWorkers: make([]Worker, 1),
		playerUnits: make([]*Unit, 0, 10),
		enemyUnits:  make([]*Unit, 0, 10),
	}
}

func (s *T) Update() {

	for _, unit := range s.playerUnits {
		//nearestEnemy := findNearestEnemy()
		//unit.SetNearestEnemy()
		unit.Update(s)
	}
	for _, building := range s.playerBuildings {
		building.Update(s)
	}
	for _, unit := range s.enemyUnits {
		unit.Update(s)
	}
	// update resource counts
	// Update unit movement
	// calculate damage done
	// process unit removals
	// process unit additions
}

func (s *T) RemoveUnit(u *Unit) {
	s.playerUnits = slices.DeleteFunc(s.playerUnits, func(other *Unit) bool {
		return other.ID == u.ID
	})
}

func (s *T) AddUnit(u *Unit) {
	s.playerUnits = append(s.playerUnits, u)
}
func (s *T) AddHive(h *Hive) {
	s.playerBuildings = append(s.playerBuildings, h)
}

func (s *T) GetUnitByID(id string) (*Unit, error) {
	for _, unit := range s.playerUnits {
		if unit.ID.String() == id {
			return unit, nil
		}
	}
	for _, unit := range s.enemyUnits {
		if unit.ID.String() == id {
			return unit, nil
		}
	}
	return nil, fmt.Errorf("unable to find unit with ID:%v", id)
}

func (s *T) IssueAction(id string, action Action, point *image.Point) error {
	unit, err := s.GetUnitByID(id)
	if err != nil {
		return err
	}
	unit.Destination = point
	unit.DestinationType = s.DetermineDestinationType(point)
	// TODO: take passed in ACTION into account as it might matter for some UI buttons
	switch unit.DestinationType {
	case EnemyDestination:
		unit.Action = AttackMovingAction
	case ResourceDestination:
		unit.Action = CollectingAction
	case LocationDestination:
		unit.Action = AttackMovingAction
	}

	return nil
}

func (s *T) DetermineDestinationType(point *image.Point) DestinationType {
	// This should use factions instead
	for _, enemy := range s.enemyUnits {
		if point.In(*enemy.Rect) {
			return EnemyDestination
		}
	}
	tile := s.world.TileMap.GetTileByPosition(point.X, point.Y)
	if tile != nil && tile.Type == "resource" {
		return ResourceDestination
	}

	return LocationDestination
}

func (s *T) GetAllUnits() []*Unit {
	return append(s.enemyUnits, s.playerUnits...)
}

func (s *T) GetAllNearbyColliders(x, y int) []*image.Rectangle {
	var nearbyColliders []*image.Rectangle
	for _, unit := range append(s.enemyUnits, s.playerUnits...) {
		distance := unit.DistanceTo(image.Pt(x, y))
		if distance == 0 {
			continue
		}
		if distance <= NearbyDistance {
			nearbyColliders = append(nearbyColliders, unit.Rect)
		}
	}
	for _, building := range s.playerBuildings {
		distance := building.DistanceTo(image.Pt(x, y))
		if distance == 0 {
			continue
		}
		if distance <= NearbyDistance {
			nearbyColliders = append(nearbyColliders, building.Rect)
		}
	}
	return nearbyColliders
}

func (s *T) GetAllBuildings() []*Hive {
	return s.playerBuildings
}

// func (s *T) findNearestEnemy(u *Unit) *Unit {

// }
