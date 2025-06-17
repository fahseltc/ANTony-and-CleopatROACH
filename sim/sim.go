package sim

import (
	"slices"
	"sync/atomic"
)

// T is the simulated world for the game. This Sim makes
// a couple assumptions:
//
//   - The Game ticks at a fixed rate (tps)
//   - The time delta between ticks is (mostly) fixed
type T struct {
	tps   int
	dt    float64
	world World

	// ID generator
	unitIdx atomic.Uint32

	playerSpawnX, playerSpawnY float64
	playerWorkers              []Worker
	playerUnits                []Unit

	enemyUnits               []Unit
	enemySpawnX, enemySpawnY float64

	selectedUnits []*Unit
}

type World struct {
	TileData [][]Tile
}

type Tile struct {
	Type     string
	Resource *Resource
}

// Resource represents a collectable resource on the map
type Resource struct {
	Type      string
	X, Y      float64
	Available uint
	// TODO do we want to limit number of assigned units?
	// FreeSlots uint
}

func New(tps int) *T {
	return &T{
		tps: tps,
		dt:  float64(1 / tps),
		world: World{
			TileData: make([][]Tile, 0, 1),
		},
		playerWorkers: make([]Worker, 1),
		playerUnits:   make([]Unit, 0, 10),
		enemyUnits:    make([]Unit, 0, 10),
	}
}

func (s *T) Tick() {
	// update resource counts
	// Update unit movement
	// calculate damage done
	// process unit removals
	// process unit additions
}

func (s *T) NewUnitID() uint32 {
	return s.unitIdx.Add(1)
}

func (s *T) RemoveUnit(u *Unit) {
	s.playerUnits = slices.DeleteFunc(s.playerUnits, func(other Unit) bool {
		return other.ID == u.ID
	})
}
