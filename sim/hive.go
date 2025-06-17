package sim

import "sync"

type Hive struct {
	X, Y       float64
	Workers    []*Worker
	MaxWorkers uint16
	Resources  map[string]uint16
	resourceMu sync.RWMutex
}
