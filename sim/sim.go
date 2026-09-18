package sim

import (
	"fmt"
	"gamejam/eventing"
	"gamejam/log"
	"gamejam/tilemap"
	"gamejam/types"
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"log/slog"
	"math"
	"slices"
)

var NearbyDistance = uint(300)

// BuilderMaxDistance is the maximum edge-to-target distance (in pixels) a
// selected unit may be from the tile it is trying to build on. Roughly 1.5
// tiles (tiles are 128px), so the builder must be close to the build site.
var BuilderMaxDistance = uint(325)

type T struct {
	EventBus *eventing.EventBus
	tps      int
	dt       float64
	world    *World

	playerState *PlayerState

	unitMap     map[int][]*Unit
	buildingMap map[int][]BuildingInterface

	ActionKeyPressed ActionKeyPressed

	log *slog.Logger
}

type Collider struct {
	Rect    *image.Rectangle
	Center  image.Point
	Radius  uint
	OwnerID string
}

type PlayerState struct {
	Sucrose  uint
	Wood     uint
	TechTree *TechTree
}

func (ps *PlayerState) GetTechTree() *TechTree {
	return ps.TechTree
}

type ActionKeyPressed uint

const (
	NoneKeyPressed ActionKeyPressed = iota
	AttackKeyPressed
	MoveKeyPressed
	StopKeyPressed
	HoldPositionKeyPressed
)

func (s *T) GetPlayerState() *PlayerState {
	return s.playerState
}
func (s *T) SetPlayerState(ps *PlayerState) {
	s.playerState = ps
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
		playerState: &PlayerState{
			Sucrose:  0, //9000,
			Wood:     0, //100,
			TechTree: NewTechTree(),
		},
		// playerUnits: make([]*Unit, 0, 10),
		// enemyUnits:  make([]*Unit, 0, 10),
		unitMap:     make(map[int][]*Unit),
		buildingMap: make(map[int][]BuildingInterface),

		log: log.NewLogger().With("for", "sim"),
	}
	bus.Subscribe("ConstructUnitEvent", sim.HandleConstructUnitEvent)
	return sim
}
func (s *T) HandleConstructUnitEvent(event eventing.Event) {
	hiveID := event.Data.(eventing.ConstructUnitEvent).HiveID
	unitType := types.UtilUnitTypeFromString(event.Data.(eventing.ConstructUnitEvent).UnitType)
	success := s.ConstructUnit(hiveID, unitType)
	if !success {
		// Report the resource the player is actually short on, rather than
		// always blaming sucrose.
		cost := UtilUnitTypeToUnit(unitType).Stats.ResourceCost
		resName := "Sucrose"
		if s.playerState.Wood < cost.Wood && s.playerState.Sucrose >= cost.Sucrose {
			resName = "Wood"
		}
		s.EventBus.Publish(eventing.Event{
			Type: "NotEnoughResourcesEvent",
			Data: eventing.NotEnoughResourcesEvent{
				ResourceName:   resName,
				UnitBeingBuilt: unitType.DisplayName(),
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
	if u == nil {
		return
	}
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
	// Bridges are walkable: they exist precisely to open up an otherwise
	// impassable tile, so they must not add a collision rect back onto it.
	// (The runtime collider check in GetAllCollidersOverlapping already skips
	// bridges for the same reason.)
	if b.GetType() != types.BuildingTypeBridge {
		s.world.TileMap.AddCollisionRect(b.GetRect())
	} else {
		// Explicitly mark the bridge's tile walkable on the pathing grid. The
		// underlying chasm collision object can't be removed by exact-rect match
		// (it may span multiple tiles), so instead we clear the grid cell for
		// this tile directly. The runtime collision (isColliding) is made
		// bridge-aware via IsCoveredByBridge.
		rect := b.GetRect()
		tileX := (rect.Min.X + rect.Dx()/2) / int(TileSize)
		tileY := (rect.Min.Y + rect.Dy()/2) / int(TileSize)
		s.world.TileMap.SetTileWalkable(tileX, tileY)
	}
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
	s.log.Debug("issuing action", "currentActionKey", s.ActionKeyPressed)
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

// IssueHarvestTile orders a single unit (by ID) to harvest the resource on the
// given TILE coordinate. It converts the tile to its center pixel and routes
// through IssueAction, which starts a HarvestingState when the tile is a
// resource. Intended for level setup / cutscenes where you want a worker to
// begin gathering a specific node without a player click.
//
// The tile must be a wood or sucrose tile; if it isn't (or is off-map), this
// logs a warning and the unit will simply move there (or do nothing) instead of
// harvesting, so authoring mistakes are visible in the logs.
func (s *T) IssueHarvestTile(unitID string, tileX, tileY int) error {
	tile := s.world.TileMap.GetTileByCoordinates(tileX, tileY)
	if tile == nil {
		s.log.Warn("IssueHarvestTile: tile not found", "unitID", unitID, "tileX", tileX, "tileY", tileY)
		return fmt.Errorf("harvest tile (%d,%d) not found", tileX, tileY)
	}
	if tile.Type != types.TileTypeWood && tile.Type != types.TileTypeSucrose {
		s.log.Warn("IssueHarvestTile: tile is not a resource; unit will not harvest",
			"unitID", unitID, "tileX", tileX, "tileY", tileY, "tileType", tile.Type)
	}
	// Target the tile CENTER so DetermineDestinationType reliably samples this
	// tile (an edge pixel could round into a neighbour).
	center := &image.Point{
		X: tileX*int(TileSize) + int(HalfTileSize),
		Y: tileY*int(TileSize) + int(HalfTileSize),
	}
	return s.IssueAction([]string{unitID}, center)
}

func (s *T) issueSingleAction(id string, point *image.Point) error {
	unit, err := s.GetUnitByID(id)
	if err != nil {
		return err
	}

	clickedTile := s.world.TileMap.GetTileByPosition(point.X, point.Y)
	if clickedTile == nil {
		s.log.Warn("tile clicked was not found", "x", point.X, "y", point.Y)
		return fmt.Errorf("tile clicked was not found")
	}

	unit.DestinationType = s.DetermineDestinationType(point, unit.Faction)
	// switch s.ActionKeyPressed {
	// case NoneKeyPressed:
	// case AttackKeyPressed:
	// case MoveKeyPressed:
	// case StopKeyPressed:
	// case HoldPositionKeyPressed:
	// }
	switch unit.DestinationType {
	case types.DestinationTypeEnemy:
		unit.ChangeState(&AttackMoveState{})
	case types.DestinationTypeResource:
		unit.LastResourcePos = &vec2.T{
			X: float64(clickedTile.Coordinates.X*int(TileSize) + int(HalfTileSize)),
			Y: float64(clickedTile.Coordinates.Y*int(TileSize) + int(HalfTileSize)),
		}
		unit.ChangeState(&HarvestingState{})
	case types.DestinationTypeTile:
		if s.ActionKeyPressed == AttackKeyPressed {
			unit.ChangeState(&AttackMoveState{})
		} else {
			unit.ChangeState(&MovingState{})
		}
	}

	start := unit.GetTileCoordinates()

	unit.Destinations.Clear()

	if unit.DestinationType == types.DestinationTypeResource {
		// The clicked resource tile is impassable, so route to a walkable tile
		// adjacent to it (this worker's own approach slot) and path all the way
		// there. This avoids A* aiming at the resource tile itself and avoids
		// the final waypoint jumping to an unvalidated offset.
		approach := unit.HarvestApproachPos(s, unit.LastResourcePos)
		approachTile := &vec2.T{
			X: math.Floor(approach.X / TileSize),
			Y: math.Floor(approach.Y / TileSize),
		}
		steps := s.FindClickedPath(start, approachTile)
		for _, step := range steps {
			unit.Destinations.Enqueue(&vec2.T{X: step.X*TileSize + HalfTileSize, Y: step.Y*TileSize + HalfTileSize})
		}
		unit.Destinations.Enqueue(approach)
		return nil
	}

	// Non-resource orders: path to the clicked tile and finish at the exact
	// clicked pixel.
	end := util.PointToVec2(clickedTile.Coordinates)
	steps := s.FindClickedPath(start, end)
	for _, step := range steps {
		unit.Destinations.Enqueue(&vec2.T{X: step.X*TileSize + HalfTileSize, Y: step.Y*TileSize + HalfTileSize})
	}
	unit.Destinations.Enqueue(&vec2.T{X: float64(point.X), Y: float64(point.Y)})

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
	// If the START tile is unwalkable, A* can't leave it and returns nothing.
	// This happens when a unit is standing on a tile the grid treats as blocked
	// - e.g. a bridge EDGE tile, or a tile its rounded position snapped onto
	// that has no bridge over the water. Snap the start to the nearest walkable
	// neighbour so pathing can proceed, mirroring the END recovery below.
	if startTile := s.world.TileMap.GetTileByCoordinates(int(start.X), int(start.Y)); startTile != nil && startTile.HasCollision {
		if fixedStart := s.FindNearestSurroundingWalkableTiles(start, start); fixedStart != nil {
			start = fixedStart
		}
	}

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

func (s *T) DetermineDestinationType(point *image.Point, selfFaction uint) types.Destination {
	for factionIndex, um := range s.unitMap {
		if factionIndex == int(selfFaction) {
			continue
		}
		for _, u := range um {
			if point.In(*u.Rect) {
				return types.DestinationTypeEnemy
			}
		}
	}

	tile := s.world.TileMap.GetTileByPosition(point.X, point.Y)
	if tile != nil && (tile.Type == types.TileTypeWood || tile.Type == types.TileTypeSucrose) {
		return types.DestinationTypeResource
	}

	return types.DestinationTypeTile
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
		if building.GetType() == types.BuildingTypeBridge { // bridges dont have collision!
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
	// Use the tilemap's live collision objects. s.world.MapObjects is a snapshot
	// taken at sim construction and goes stale after collision rects are added or
	// removed (the tilemap reassigns its own slice), so it can diverge from the
	// pathing grid.
	for _, mapObj := range s.world.TileMap.MapObjects {
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

// IsCoveredByBridge reports whether the given pixel point lies on a completed
// bridge. Bridges are walkable and exist to open up an otherwise-impassable
// tile, so a point over a bridge should ignore the underlying chasm collision.
func (s *T) IsCoveredByBridge(x, y int) bool {
	pt := image.Point{X: x, Y: y}
	for _, building := range s.GetAllBuildings() {
		if building.GetType() != types.BuildingTypeBridge {
			continue
		}
		if pt.In(*building.GetRect()) {
			return true
		}
	}
	return false
}

func (s *T) DetermineUnitOrHiveById(id string) string { // TODO use building.GetType()
	b, err := s.GetBuildingByID(id)
	if err == nil && (b.GetType() == types.BuildingTypeAntHive || b.GetType() == types.BuildingTypeRoachHive) {
		return "hive"
	}
	_, err2 := s.GetUnitByID(id)
	if err2 == nil {
		return "unit"
	}
	return "neither"
}

func (s *T) AddResource(amount uint, resType types.Resource) {
	switch resType {
	case types.ResourceTypeSucrose:
		s.playerState.Sucrose += amount
	case types.ResourceTypeWood:
		s.playerState.Wood += amount
	}
}

func (s *T) GetWoodAmount() uint {
	return s.playerState.Wood
}
func (s *T) GetSucroseAmount() uint {
	return s.playerState.Sucrose
}
func (s *T) ConstructUnit(hiveId string, unitType types.Unit) bool {
	hive, err := s.GetBuildingByID(hiveId)
	if err != nil {
		return false
	}
	// A roach hive produces roaches from the same "worker" button an ant hive
	// uses ants for. The button is generic (it always requests the default
	// worker), so pick the faction-appropriate worker based on the hive type.
	if unitType == types.UnitTypeDefaultAnt && hive.GetType() == types.BuildingTypeRoachHive {
		unitType = types.UnitTypeDefaultRoach
	}
	newUnit := UtilUnitTypeToUnit(unitType)
	if s.playerState.Sucrose >= newUnit.Stats.ResourceCost.Sucrose &&
		s.playerState.Wood >= newUnit.Stats.ResourceCost.Wood {
		s.playerState.Sucrose -= newUnit.Stats.ResourceCost.Sucrose
		s.playerState.Wood -= newUnit.Stats.ResourceCost.Wood

		hive.AddItemToBuildQueue(&QueuedItem{
			Type:             types.QueuedItemTypeUnit,
			Unit:             newUnit,
			ConstructionTime: newUnit.Stats.ConstructionTime,
		})
		return true
	} else {
		return false
	}
}

func (s *T) ConstructBuilding(tileCoords image.Point, builderID string, buildingType types.Building) bool {
	unit, err := s.GetUnitByID(builderID)
	if err != nil {
		return false // todo print builder doesnt exist
	}

	building := UtilBuildingTypeToBuilding(buildingType)

	// Bridges may only be placed on water the map author marked buildable.
	// Reject placement anywhere else BEFORE spending resources so the click is
	// a no-op on invalid ground, and tell the player why.
	if buildingType == types.BuildingTypeBridge &&
		!s.world.TileMap.IsBridgeBuildable(tileCoords.X, tileCoords.Y) {
		s.EventBus.Publish(eventing.Event{
			Type: "NotificationEvent",
			Data: eventing.NotificationEvent{
				Message: "Bridges can only be built on water.",
			},
		})
		return false
	}

	// The builder must be close to the target tile. Measure from the unit's
	// nearest edge to the target tile center so the requirement is forgiving of
	// which side the unit stands on.
	targetCenter := &vec2.T{
		X: float64(tileCoords.X)*TileSize + HalfTileSize,
		Y: float64(tileCoords.Y)*TileSize + HalfTileSize,
	}
	if unit.EdgeDistanceTo(targetCenter) > BuilderMaxDistance {
		s.EventBus.Publish(eventing.Event{
			Type: "NotificationEvent",
			Data: eventing.NotificationEvent{
				Message: "Move closer to build here.",
			},
		})
		return false
	}

	if !building.GetStats().ResourceCost.CanAfford(*s.playerState) {
		s.EventBus.Publish(eventing.Event{
			Type: "NotEnoughResourcesEvent",
			Data: eventing.NotEnoughResourcesEvent{
				ResourceName:   "Wood",
				UnitBeingBuilt: buildingType.ToString(),
			},
		})
		return false
	} else {
		// actually build the thing
		bought := building.GetStats().ResourceCost.Purchase(s.playerState)
		if bought {
			icb := GetInConstructionBuildingInstance(buildingType, uint(PlayerFaction))
			icb.SetTilePosition(tileCoords.X, tileCoords.Y)
			s.buildingMap[int(unit.Faction)] = append(s.buildingMap[int(unit.Faction)], icb)
			return true
		} else {
			return false
		}
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

	// Only collapse to a straight line if that line is walkable on the grid AND
	// does not clip any collision object. isLineWalkable samples tiles via
	// Bresenham and can skip a thin obstacle at a tile corner; the mover in
	// MoveToDestination collides against world.MapObjects, so we must verify the
	// segment against those too. Otherwise the optimizer can route a worker
	// straight through an impassable resource/obstacle and the mover then wedges
	// there forever (never "arrives", never dequeues).
	if s.isLineWalkable(first, last) && s.isSegmentClearOfMapObjects(first, last) {
		return []*vec2.T{first, last} // optimized path: straight line
	}

	// fallback to original path
	return nav
}

// isSegmentClearOfMapObjects reports whether the straight segment between two
// TILE coordinates (converted to tile-center pixels) avoids every collision
// MapObject rect. This matches the collision model the mover actually uses, so
// an "optimized" straight path can't cut a corner through an obstacle the mover
// will refuse to cross.
func (s *T) isSegmentClearOfMapObjects(start, end *vec2.T) bool {
	// Convert tile coords to pixel-space endpoints at tile centers.
	x0 := start.X*TileSize + HalfTileSize
	y0 := start.Y*TileSize + HalfTileSize
	x1 := end.X*TileSize + HalfTileSize
	y1 := end.Y*TileSize + HalfTileSize

	// Sample the segment finely enough that no tile-sized obstacle can slip
	// between samples (step well under one tile).
	dist := math.Hypot(x1-x0, y1-y0)
	if dist == 0 {
		return true
	}
	step := HalfTileSize // 64px steps: never skips a 128px collision rect
	steps := int(math.Ceil(dist/step)) + 1

	for _, mo := range s.world.TileMap.MapObjects {
		for i := 0; i <= steps; i++ {
			t := float64(i) / float64(steps)
			px := x0 + (x1-x0)*t
			py := y0 + (y1-y0)*t
			if int(px) >= mo.Rect.Min.X && int(px) < mo.Rect.Max.X &&
				int(py) >= mo.Rect.Min.Y && int(py) < mo.Rect.Max.Y {
				// A bridge opens the tile beneath a collision rect, so a sample
				// point covered by a bridge does not block the straight path.
				if s.IsCoveredByBridge(int(px), int(py)) {
					continue
				}
				return false
			}
		}
	}
	return true
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

func (s *T) RevealFogOfWar(topLeft *vec2.T, bottomRight *vec2.T) {
	s.world.RevealFogOfWar(topLeft, bottomRight)
}
