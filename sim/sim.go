package sim

import (
	"fmt"
	"gamejam/eventing"
	"gamejam/tilemap"
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"log/slog"
	"math"
	"slices"
	"sync"
)

var NearbyDistance = uint(300)
var UnitSucroseCost = uint16(50)
var BuildingWoodCost = uint16(50)
var BuilderMaxDistance = uint(340)

type T struct {
	EventBus *eventing.EventBus
	tps      int
	dt       float64
	world    *World

	playerState PlayerState
	stateMu     sync.RWMutex

	unitMap     map[int][]*Unit
	buildingMap map[int][]BuildingInterface

	ActionKeyPressed ActionKeyPressed
}

type Collider struct {
	Rect    *image.Rectangle
	Center  image.Point
	Radius  uint
	OwnerID string
}

type PlayerState struct {
	Sucrose uint16
	Wood    uint16
}

type ResourceType uint

const (
	ResourceTypeNone ResourceType = iota
	ResourceTypeSucrose
	ResourceTypeWood
)

type ActionKeyPressed uint

const (
	NoneKeyPressed ActionKeyPressed = iota
	AttackKeyPressed
	MoveKeyPressed
	StopKeyPressed
	HoldPositionKeyPressed
)

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
			FogOfWar:   NewFogOfWar(tileMap.Width, tileMap.Height),
		},

		// TODO Spawn Points
		playerState: PlayerState{
			Sucrose: 9999,
		},
		// playerUnits: make([]*Unit, 0, 10),
		// enemyUnits:  make([]*Unit, 0, 10),
		unitMap:     make(map[int][]*Unit),
		buildingMap: make(map[int][]BuildingInterface),
	}
	bus.Subscribe("ConstructUnitEvent", sim.HandleConstructUnitEvent)
	return sim
}
func (s *T) HandleConstructUnitEvent(event eventing.Event) {
	hiveID := event.Data.(eventing.ConstructUnitEvent).HiveID
	success := s.ConstructUnit(hiveID)
	if !success {
		s.EventBus.Publish(eventing.Event{
			Type: "NotEnoughResourcesEvent",
			Data: eventing.NotEnoughResourcesEvent{
				ResourceName:     "Sucrose",
				TargetBeingBuilt: "Ant",
			},
		})
	}
}

func (s *T) Update() {
	for _, um := range s.unitMap {
		for _, unit := range um {
			unit.Update(s)
		}
	}
	for _, bm := range s.buildingMap {
		for _, building := range bm {
			building.Update(s)
		}
	}
}

func (s *T) RemoveUnit(u *Unit) {
	faction := u.Faction
	if s.unitMap[int(faction)] == nil {
		return
	}
	s.unitMap[int(faction)] = slices.DeleteFunc(s.unitMap[int(faction)], func(other *Unit) bool {
		return other.ID.String() == u.ID.String()
	})
}

func (s *T) AddUnit(u *Unit) {
	faction := u.Faction
	if s.unitMap[int(faction)] == nil {
		s.unitMap[int(faction)] = make([]*Unit, 0)
	}
	s.unitMap[int(faction)] = append(s.unitMap[int(faction)], u)
}

func (s *T) AddBuilding(b BuildingInterface) {
	faction := b.GetFaction()
	if s.buildingMap[int(faction)] == nil {
		s.buildingMap[int(faction)] = make([]BuildingInterface, 0)
	}
	s.buildingMap[int(faction)] = append(s.buildingMap[int(faction)], b)
}
func (s *T) RemoveBuilding(b BuildingInterface) {
	faction := b.GetFaction()
	if s.buildingMap[int(faction)] == nil {
		return
	}
	s.buildingMap[int(faction)] = slices.DeleteFunc(s.buildingMap[int(faction)], func(other BuildingInterface) bool {
		return other.GetID() == b.GetID()
	})
	s.world.TileMap.RemoveCollisionRect(b.GetRect())
}

func (s *T) GetUnitByID(id string) (*Unit, error) {
	for _, um := range s.unitMap {
		for _, unit := range um {
			if unit.ID.String() == id {
				return unit, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to find unit with ID:%v", id)
}
func (s *T) GetBuildingByID(id string) (BuildingInterface, error) {
	for _, bm := range s.buildingMap {
		for _, bld := range bm {
			if bld.GetID().String() == id {
				return bld, nil
			}
		}
	}
	return nil, fmt.Errorf("unable to find building with ID:%v", id)
}

func (s *T) IssueAction(ids []string, point *image.Point) error {
	fmt.Printf("currentActionKey: %v\n", s.ActionKeyPressed)
	if len(ids) == 0 {
		return fmt.Errorf("no unit IDs passed")
	} else if len(ids) == 1 {
		s.issueSingleAction(ids[0], point)
	} else {
		s.issueGroupAction(ids, point)
	}
	s.ActionKeyPressed = NoneKeyPressed
	return nil
}

func (s *T) issueSingleAction(id string, point *image.Point) error {
	unit, err := s.GetUnitByID(id)
	if err != nil {
		return err
	}

	clickedTile := s.world.TileMap.GetTileByPosition(point.X, point.Y)
	if clickedTile == nil {
		slog.Warn("tile clicked was not found")
		return fmt.Errorf("tile clicked was not found")
	}

	unit.DestinationType = s.DetermineDestinationType(point, unit.Faction)
	switch s.ActionKeyPressed {
	case NoneKeyPressed:
	case AttackKeyPressed:
	case MoveKeyPressed:
	case StopKeyPressed:
	case HoldPositionKeyPressed:
	}
	switch unit.DestinationType {
	case EnemyDestination:
		unit.Action = AttackMovingAction
	case ResourceDestination:
		unit.Action = CollectingAction
		unit.LastResourcePos = &vec2.T{
			X: float64(clickedTile.Coordinates.X*int(TileSize) + int(HalfTileSize)),
			Y: float64(clickedTile.Coordinates.Y*int(TileSize) + int(HalfTileSize)),
		}
	case LocationDestination:
		unit.Action = AttackMovingAction
	}

	start := unit.GetTileCoordinates()
	end := util.PointToVec2(clickedTile.Coordinates)

	// parse A* path into a series of destinations.
	unit.Destinations.Clear()
	steps := s.FindClickedPath(start, end)
	for _, step := range steps {
		unit.Destinations.Enqueue(&vec2.T{X: step.X*TileSize + HalfTileSize, Y: step.Y*TileSize + HalfTileSize}) // add 64 for half tile?
	}

	return nil
}
func (s *T) issueGroupAction(ids []string, point *image.Point) error {
	if len(ids) == 0 {
		return nil
	}

	// Gather units by ID
	var units []*Unit
	var positions []*vec2.T
	for _, id := range ids {
		unit, err := s.GetUnitByID(id)
		if err == nil && unit != nil {
			units = append(units, unit)
			positions = append(positions, unit.GetCenteredPosition())
		}
	}
	if len(units) == 0 {
		return nil
	}

	centroid := s.FindUnitGroupCenter(positions)
	// If units are very far apart, migrate them closer to the centroid first
	const maxSpread = 120.0 // adjust as needed
	var maxDist float64
	for _, unit := range units {
		dist := unit.GetCenteredPosition().Distance(centroid)
		if dist > maxDist {
			maxDist = dist
		}
	}
	if maxDist > maxSpread {
		// Scale each unit's offset from the centroid so that the group contracts toward the centroid,
		// but maintains relative formation, and then move them to the scaled positions.
		scale := maxSpread / maxDist
		destPoint := vec2.T{X: float64(point.X), Y: float64(point.Y)}
		for _, unit := range units {
			unitPos := unit.GetCenteredPosition()
			offset := vec2.T{X: unitPos.X - centroid.X, Y: unitPos.Y - centroid.Y}
			scaledOffset := vec2.T{X: offset.X * scale, Y: offset.Y * scale}
			target := image.Point{
				X: int(destPoint.X + scaledOffset.X),
				Y: int(destPoint.Y + scaledOffset.Y),
			}
			_ = s.IssueAction([]string{unit.ID.String()}, &target)
		}
		return nil
	}

	// Calculate offset for each unit from the centroid and issue final move
	destPoint := vec2.T{X: float64(point.X), Y: float64(point.Y)}
	for _, unit := range units {
		unitPos := unit.GetCenteredPosition()
		offset := vec2.T{X: unitPos.X - centroid.X, Y: unitPos.Y - centroid.Y}
		target := image.Point{
			X: int(destPoint.X + offset.X),
			Y: int(destPoint.Y + offset.Y),
		}
		_ = s.IssueAction([]string{unit.ID.String()}, &target)
	}

	return nil
}

// accepts integar map coordinates (not pixels)
func (s *T) FindClickedPath(start *vec2.T, end *vec2.T) []*vec2.T {
	path := s.world.TileMap.FindPath(start, end)
	if len(path) != 0 {
		return s.optimizePath(path)
	}
	firstEndingPos := end

	for len(path) == 0 {
		end = s.FindNearestSurroundingWalkableTiles(start, end)
		if end == nil || (firstEndingPos.X == end.X && firstEndingPos.Y == end.Y) { // water or somewhere completely unwalkable was clicked
			return nil
		}
		path = s.world.TileMap.FindPath(start, end)
		if path == nil {
			return nil
		}
	}
	return s.optimizePath(path)
}

// Accepts unwalkableTile with integer Tile coordinates (not pixel coordinates)
func (s *T) FindNearestSurroundingWalkableTiles(currentPos *vec2.T, unwalkableCoords *vec2.T) *vec2.T {
	var walkableTiles []*vec2.T
	for _, bldg := range s.GetAllBuildings() {
		rect := bldg.GetRect()
		ux, uy := int(unwalkableCoords.X*TileSize), int(unwalkableCoords.Y*TileSize)
		if ux >= rect.Min.X && ux < rect.Max.X && uy >= rect.Min.Y && uy < rect.Max.Y {
			walkableTiles = bldg.GetAdjacentCoordinates()
			// TODO: check if these are all walkable?
		}
	}
	if walkableTiles == nil {
		directions := []struct{ dx, dy int }{
			{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // cardinal directions
			//{-1, -1}, {-1, 1}, {1, -1}, {1, 1}, // diagonals
		}

		x, y := int(unwalkableCoords.X), int(unwalkableCoords.Y)
		for _, dir := range directions {
			nx, ny := x+dir.dx, y+dir.dy
			tile := s.world.TileMap.GetTileByCoordinates(nx, ny)
			if tile != nil && !tile.HasCollision {
				walkableTiles = append(walkableTiles, &vec2.T{X: float64(nx), Y: float64(ny)})
			}
		}
	}

	if len(walkableTiles) == 0 {
		return nil
	}
	closest := walkableTiles[0]
	minDist := currentPos.Distance(*closest)
	for _, tile := range walkableTiles[1:] {
		dist := currentPos.Distance(*tile)
		if dist < minDist {
			minDist = dist
			closest = tile
		}
	}
	return closest
}

func (s *T) DetermineDestinationType(point *image.Point, selfFaction uint) DestinationType {
	for factionIndex, um := range s.unitMap {
		if factionIndex == int(selfFaction) {
			continue
		}
		for _, u := range um {
			if point.In(*u.Rect) {
				return EnemyDestination
			}
		}
	}

	tile := s.world.TileMap.GetTileByPosition(point.X, point.Y)
	if tile != nil && (tile.Type == "wood" || tile.Type == "sucrose") {
		return ResourceDestination
	}

	return LocationDestination
}

func (s *T) GetAllUnits() []*Unit {
	var allUnits []*Unit
	for _, um := range s.unitMap {
		allUnits = append(allUnits, um...)
	}
	return allUnits
}

func (s *T) GetAllUnitsByFaction(faction uint) []*Unit {
	if s.unitMap[int(faction)] == nil {
		return nil
	}
	return s.unitMap[int(faction)]
}

func (s *T) GetAllPlayerUnits() []*Unit {
	if s.unitMap[PlayerFaction] == nil {
		return nil
	}
	return s.unitMap[PlayerFaction]
}

func (s *T) GetAllEnemyUnitsByFaction(faction uint) []*Unit {
	var allEnemyUnits []*Unit
	for factionIndex, um := range s.unitMap {
		if factionIndex == int(faction) { // skip self-faction
			continue
		}
		allEnemyUnits = append(allEnemyUnits, um...)
	}
	return allEnemyUnits
}

func (s *T) GetAllNearbyFriendlyUnits(sourceUnit *Unit) []*Unit {
	var nearbyUnits []*Unit
	if s.unitMap[int(sourceUnit.Faction)] == nil {
		return nil
	}
	for _, unit := range s.unitMap[int(sourceUnit.Faction)] {
		if sourceUnit.ID.String() == unit.ID.String() {
			continue
		}
		unitDist := unit.GetCenteredPosition().Distance(*sourceUnit.GetCenteredPosition())
		if unitDist <= 150 {
			nearbyUnits = append(nearbyUnits, unit)
		}
	}
	return nearbyUnits
}

func (s *T) GetAllCollidersOverlapping(rect *image.Rectangle) []*Collider {
	var colliders []*Collider
	for _, unit := range s.GetAllUnits() {
		if unit == nil {
			continue
		}
		if unit.Rect.Overlaps(*rect) {
			colliders = append(colliders, &Collider{
				Rect:    unit.Rect,
				Center:  unit.GetCenteredPosition().ToPoint(),
				Radius:  uint(unit.Rect.Dx() / 2),
				OwnerID: unit.ID.String(),
			})
		}
	}
	for _, building := range s.GetAllBuildings() {
		if building.GetType() == BuildingTypeBridge { // bridges dont have collision!
			continue
		}
		if building.GetRect().Overlaps(*rect) {
			colliders = append(colliders, &Collider{
				Rect:    building.GetRect(),
				Center:  building.GetCenteredPosition().ToPoint(),
				Radius:  uint(building.GetRect().Dx() / 2),
				OwnerID: building.GetID().String(),
			})
		}
	}
	for _, mapObj := range s.world.MapObjects {
		if mapObj.Rect.Overlaps(*rect) {
			colliders = append(colliders, &Collider{
				Rect:    mapObj.Rect,
				Center:  image.Point{},
				Radius:  0,
				OwnerID: "map",
			})
		}
	}

	return colliders
}

func (s *T) GetAllBuildings() []BuildingInterface {
	var allBuildings []BuildingInterface
	for _, bm := range s.buildingMap {
		allBuildings = append(allBuildings, bm...)
	}
	return allBuildings
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

func (s *T) AddResource(amount uint, resType ResourceType) {
	switch resType {
	case ResourceTypeSucrose:
		s.playerState.Sucrose += uint16(amount)
	case ResourceTypeWood:
		s.playerState.Wood += uint16(amount)
	}
}

func (s *T) GetWoodAmount() uint16 {
	return s.playerState.Wood
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

	targetCenter := vec2.T{X: float64(target.Min.X + (target.Dx() / 2)), Y: float64(target.Min.Y + (target.Dy() / 2))}

	if unit.DistanceTo(&targetCenter) > BuilderMaxDistance { // todo min should be center!
		return false
	} else {
		// actually build the thing
		s.playerState.Wood -= BuildingWoodCost
		inConstructionBuilding := NewInConstructionBuilding(target.Min.X, target.Min.Y, BuildingTypeBridge) // always bridge for now, but easy to change
		s.buildingMap[int(unit.Faction)] = append(s.buildingMap[int(unit.Faction)], inConstructionBuilding)
		return true
	}
}

func (s *T) GetWorld() *World {
	return s.world
}
func (s *T) optimizePath(nav []*vec2.T) []*vec2.T {
	if len(nav) <= 2 {
		return nav // nothing to optimize
	}

	first := nav[0]
	last := nav[len(nav)-1]

	// Check if straight path is walkable
	if s.isLineWalkable(first, last) {
		return []*vec2.T{first, last} // optimized path: straight line
	}

	// fallback to original path
	return nav
}

func (s *T) isLineWalkable(start, end *vec2.T) bool {
	x0 := int(start.X)
	y0 := int(start.Y)
	x1 := int(end.X)
	y1 := int(end.Y)

	dx := math.Abs(float64(x1 - x0))
	dy := math.Abs(float64(y1 - y0))
	sx := -1
	sy := -1

	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}

	err := dx - dy

	for {
		// Check tile at current position
		tile := s.world.TileMap.GetTileByCoordinates(x0, y0)
		if tile == nil || tile.HasCollision {
			return false
		}

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
	return true
}

func (s *T) SetActionKeyPressed(key ActionKeyPressed) {
	s.ActionKeyPressed = key
}

func (s *T) FindUnitGroupCenter(positions []*vec2.T) vec2.T {
	// Calculate the centroid of the selected units
	var sumX, sumY float64
	for _, pos := range positions {
		sumX += pos.X
		sumY += pos.Y
	}
	return vec2.T{X: sumX / float64(len(positions)), Y: sumY / float64(len(positions))}
}
