package sim

import (
	"encoding/json"
	"gamejam/data"
	"gamejam/types"
	"gamejam/util"
	"gamejam/vec2"
	"image"
	"sync"

	"github.com/google/uuid"
)

type UnitFactory struct {
	Units map[types.Unit]UnitStats
}

var (
	unitConfigInstance *UnitFactory
	once               sync.Once
)

func getUnitFactory() *UnitFactory {
	once.Do(func() {
		unitConfigInstance = loadUnitConfig()
	})
	return unitConfigInstance
}

func GetUnitInstance(unitType types.Unit, faction uint) *Unit {
	fact := getUnitFactory()
	stats, ok := fact.Units[unitType]
	if !ok {
		stats = fact.Units[types.UnitTypeDefaultAnt]
	}

	u := &Unit{
		ID:           uuid.New(),
		Type:         unitType,
		Stats:        &stats,
		Position:     &vec2.T{},
		Destinations: util.NewQueue[*vec2.T](),
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{int(stats.SizePx), int(stats.SizePx)},
		},
		Faction: uint(faction),
	}
	u.SetPosition(&vec2.T{X: 0, Y: 0}) // or wherever default position should be
	return u
}

type rawUnitStat struct {
	Name             string
	HPMax            uint
	MoveSpeed        uint
	Damage           uint
	AttackRange      uint
	AttackFrames     uint
	VisionRange      uint
	SizePx           uint
	MaxCarryCapacity uint
	ConstructionTime uint
	ResourceCost     ResourceCost
}

func loadUnitConfig() *UnitFactory {
	var rawUnits []rawUnitStat
	jsonFile, err := data.Files.ReadFile("unit_stats.json")
	if err != nil {
		panic("failed to read unit_config.json: " + err.Error())
	}

	if err := json.Unmarshal(jsonFile, &rawUnits); err != nil {
		panic("failed to parse unit_config.json: " + err.Error())
	}

	fact := &UnitFactory{
		Units: make(map[types.Unit]UnitStats),
	}

	for _, raw := range rawUnits {
		unitType := types.UtilUnitTypeFromString(raw.Name)
		fact.Units[unitType] = UnitStats{
			Name:                raw.Name,
			HPMax:               raw.HPMax,
			HPCur:               raw.HPMax,
			MoveSpeed:           raw.MoveSpeed,
			SizePx:              raw.SizePx,
			Damage:              raw.Damage,
			AttackRange:         raw.AttackRange,
			AttackFrames:        raw.AttackFrames,
			AttackFramesCur:     0,
			MaxCarryCapacity:    raw.MaxCarryCapacity,
			ResourcesCarried:    0,
			ResourceTypeCarried: types.ResourceTypeNone,
			ConstructionTime:    raw.ConstructionTime,
			ResourceCost:        raw.ResourceCost,
			VisionRange:         raw.VisionRange,
		}
	}

	return fact
}
