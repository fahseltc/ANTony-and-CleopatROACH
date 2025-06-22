package sim

import (
	"fmt"
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
	//TileData       [][]Tile
	CollisionRects []*image.Rectangle
}

const (
	NoneTile     = "none"
	ResourceTile = "resource"
	BuildingTile = "building"
)

type Tile struct {
	Type     string
	Passable bool
	Resource *Resource
}

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

func New(tps int, collisionRects []*image.Rectangle) *T {
	return &T{
		tps: tps,
		dt:  float64(1 / tps),
		world: &World{
			//TileData:       make([][]Tile, 0, 1),
			CollisionRects: collisionRects,
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
	unit.Action = action
	unit.Destination = point
	return nil
}

func (s *T) GetAllUnits() []*Unit {
	return append(s.enemyUnits, s.playerUnits...)
}

func (s *T) GetAllNearbyUnits(x, y int) []*Unit {
	var nearbyUnits []*Unit
	for _, unit := range append(s.enemyUnits, s.playerUnits...) {
		distance := unit.DistanceTo(image.Pt(x, y))
		if distance == 0 {
			continue
		}
		if distance <= NearbyDistance {
			nearbyUnits = append(nearbyUnits, unit)
		}
	}
	return nearbyUnits
}

// func (s *T) findNearestEnemy(u *Unit) *Unit {

// }
