package data

import (
	"encoding/json"
	"gamejam/sim"
	"gamejam/types"
)

type UnitConfig struct {
	Units map[types.Unit]sim.UnitStats
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
	ResourceCost     sim.ResourceCost
}

func NewUnitConfig() *UnitConfig {
	var rawUnits []rawUnitStat
	jsonFile, err := Files.ReadFile("unit_config.json")
	if err != nil {
		panic("failed to read unit_config.json: " + err.Error())
	}

	if err := json.Unmarshal(jsonFile, &rawUnits); err != nil {
		panic("failed to parse unit_config.json: " + err.Error())
	}

	cfg := &UnitConfig{
		Units: make(map[types.Unit]sim.UnitStats),
	}

	for _, raw := range rawUnits {
		unitType := types.UtilUnitTypeFromString(raw.Name)
		cfg.Units[unitType] = sim.UnitStats{
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

	return cfg

}
