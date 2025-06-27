package sim

import (
	"fmt"
	"gamejam/eventing"
	"gamejam/tilemap"
	"image"
	"slices"
	"sync"
)

var NearbyDistance = uint(300)
var UnitSucroseCost = uint16(50)
var BuildingWoodCost = uint16(10)
var BuilderMaxDistance = uint(340)

type T struct {
	EventBus *eventing.EventBus
	tps      int
	dt       float64
	world    *World

	playerState PlayerState
	stateMu     sync.RWMutex

	playerSpawnX, playerSpawnY float64
	playerUnits                []*Unit
	playerBuildings            []BuildingInterface
	enemyUnits                 []*Unit
	enemySpawnX, enemySpawnY   float64

	selectedUnits []*Unit
}

type World struct {
	TileMap    *tilemap.Tilemap
	TileData   [][]*tilemap.Tile
	MapObjects []*tilemap.MapObject
}

type Collider struct {
	Rect    *image.Rectangle
	OwnerID string
}

type PlayerState struct {
	Sucrose uint16
	Wood    uint16
}

func (s *T) GetPlayerState() PlayerState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.playerState
}

func New(tps int, tileMap *tilemap.Tilemap) *T {
	bus := eventing.NewEventBus()

	sim := &T{
		EventBus: bus,
		tps:      tps,
		dt:       float64(1 / tps),
		world: &World{
			TileMap:    tileMap,
			TileData:   tileMap.Tiles,
			MapObjects: tileMap.MapObjects,
		},

		// TODO Spawn Points
		playerState: PlayerState{},
		//playerWorkers: make([]Worker, 1),
		playerUnits: make([]*Unit, 0, 10),
		enemyUnits:  make([]*Unit, 0, 10),
	}
	bus.Subscribe("ConstructUnitEvent", sim.HandleConstructUnitEvent)
	return sim
}
func (s *T) HandleConstructUnitEvent(event eventing.Event) {
	hiveID := event.Data.(eventing.ConstructUnitEvent).HiveID
	s.ConstructUnit(hiveID)
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
func (s *T) AddBuilding(b BuildingInterface) {
	s.playerBuildings = append(s.playerBuildings, b)
}
func (s *T) RemoveBuilding(b BuildingInterface) {
	s.playerBuildings = slices.DeleteFunc(s.playerBuildings, func(other BuildingInterface) bool {
		return other.GetID() == b.GetID()
	})
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
func (s *T) GetBuildingByID(id string) (BuildingInterface, error) {
	for _, hive := range s.playerBuildings {
		if hive.GetID().String() == id {
			return hive, nil
		}
	}
	return nil, fmt.Errorf("unable to find unit with ID:%v", id)
}

func (s *T) IssueAction(id string, point *image.Point) error {
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
	if tile != nil && (tile.Type == "wood" || tile.Type == "sucrose") {
		return ResourceDestination
	}

	return LocationDestination
}

func (s *T) GetAllUnits() []*Unit {
	return append(s.enemyUnits, s.playerUnits...)
}

func (s *T) GetAllNearbyCollidersHarvesting(x, y int) []*image.Rectangle {
	var nearbyColliders []*image.Rectangle
	for _, unit := range s.enemyUnits {
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
			nearbyColliders = append(nearbyColliders, building.GetRect())
		}
	}
	return nearbyColliders
}
func (s *T) GetAllNearbyColliders(x, y int) []*Collider {
	var nearbyColliders []*Collider
	for _, unit := range append(s.enemyUnits, s.playerUnits...) {
		if unit == nil {
			continue
		}
		distance := unit.DistanceTo(image.Pt(x, y))
		if distance <= NearbyDistance {
			nearbyColliders = append(nearbyColliders, &Collider{
				Rect:    unit.Rect,
				OwnerID: unit.ID.String(),
			})
		}
	}
	for _, building := range s.playerBuildings {
		distance := building.DistanceTo(image.Pt(x, y))
		if distance <= NearbyDistance {
			nearbyColliders = append(nearbyColliders, &Collider{
				Rect:    building.GetRect(),
				OwnerID: building.GetID().String(),
			})
		}
	}
	return nearbyColliders
}

func (s *T) GetAllCollidersOverlapping(rect *image.Rectangle) []*Collider {
	var colliders []*Collider
	for _, unit := range append(s.enemyUnits, s.playerUnits...) {
		if unit == nil {
			continue
		}
		if unit.Rect.Overlaps(*rect) {
			colliders = append(colliders, &Collider{
				Rect:    unit.Rect,
				OwnerID: unit.ID.String(),
			})
		}
	}
	for _, building := range s.playerBuildings {
		if building.GetType() == BuildingTypeBridge { // bridges dont have collision!
			continue
		}
		if building.GetRect().Overlaps(*rect) {
			colliders = append(colliders, &Collider{
				Rect:    building.GetRect(),
				OwnerID: building.GetID().String(),
			})
		}
	}
	return colliders
}

func (s *T) GetAllBuildings() []BuildingInterface {
	return s.playerBuildings
}

func (s *T) DetermineUnitOrHiveById(id string) string { // TODO use building.GetType()
	_, err := s.GetBuildingByID(id)
	if err == nil {
		return "hive"
	}
	_, err2 := s.GetUnitByID(id)
	if err2 == nil {
		return "unit"
	}
	return "neither"
}

func (s *T) AddWood(amount uint) {
	s.playerState.Wood += uint16(amount)
}
func (s *T) GetWoodAmount() uint16 {
	return s.playerState.Wood
}
func (s *T) AddSucrose(amount uint) {
	s.playerState.Sucrose += uint16(amount)
}
func (s *T) GetSucroseAmount() uint16 {
	return s.playerState.Sucrose
}

func (s *T) ConstructUnit(hiveId string) bool {
	hive, err := s.GetBuildingByID(hiveId)
	if err != nil {
		return false
	}
	if s.playerState.Sucrose >= UnitSucroseCost {
		s.playerState.Sucrose -= UnitSucroseCost
		hive.AddUnitToBuildQueue()
		return true
	} else {
		return false
	}
}

func (s *T) ConstructBuilding(target *image.Rectangle, builderID string) bool {
	if s.playerState.Wood < BuildingWoodCost { // cant afford it
		return false
	}
	unit, err := s.GetUnitByID(builderID)
	if err != nil {
		return false // todo print builder doesnt exist
	}

	targetCenter := image.Pt(
		target.Min.X+(target.Dx()/2),
		target.Min.Y+(target.Dy()/2),
	)

	if unit.DistanceTo(targetCenter) > BuilderMaxDistance { // todo min should be center!
		return false
	} else {
		// actually build the thing
		s.playerState.Wood -= BuildingWoodCost
		inConstructionBuilding := NewInConstructionBuilding(target.Min.X, target.Min.Y, BuildingTypeBridge) // always bridge for now, but easy to change
		s.playerBuildings = append(s.playerBuildings, inConstructionBuilding)
		return true
	}
}
