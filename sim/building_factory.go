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

type BuildingFactory struct {
	Buildings map[types.Building]BuildingStats
}

var (
	buildingConfigInstance *BuildingFactory
	buildingOnce           sync.Once
)

type rawBuildingStat struct {
	Name  string
	HPMax uint

	ResourceCost     ResourceCost
	ConstructionTime uint

	VisionRange uint
	SizePx      uint
}

func getBuildingFactory() *BuildingFactory {
	buildingOnce.Do(func() {
		buildingConfigInstance = loadBuildingConfig()
	})
	return buildingConfigInstance
}

func GetBuildingInstance(buildingType types.Building, faction uint) BuildingInterface {
	fact := getBuildingFactory()
	stats, ok := fact.Buildings[buildingType]
	if !ok {
		stats = fact.Buildings[types.BuildingTypeInConstruction]
	}

	baseBuilding := Building{
		ID:       uuid.New(),
		Type:     buildingType,
		Stats:    &stats,
		Position: &vec2.T{},
		Rect: &image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{int(stats.SizePx), int(stats.SizePx)},
		},
		Faction: faction,
	}
	baseBuilding.SetTilePosition(0, 0)

	switch buildingType {
	case types.BuildingTypeNone:
		return nil
	case types.BuildingTypeInConstruction:
		return &InConstructionBuilding{
			Building: &baseBuilding,
		}
	case types.BuildingTypeAntHive:
		return &Hive{
			Building:   &baseBuilding,
			buildQueue: util.NewQueue[*QueuedItem](),
		}
	case types.BuildingTypeRoachHive:
		return &Hive{
			Building:   &baseBuilding,
			buildQueue: util.NewQueue[*QueuedItem](),
		}
	case types.BuildingTypeBridge:
		return &BridgeBuilding{
			Building: &baseBuilding,
		}
	default:
		return &Hive{
			Building:   &baseBuilding,
			buildQueue: util.NewQueue[*QueuedItem](),
		}
	}
}

// Returns an instance of InConstructionBuilding with the target building set to the incoming parameter
func GetInConstructionBuildingInstance(targetBuilding types.Building, faction uint) BuildingInterface {
	bld := GetBuildingInstance(types.BuildingTypeInConstruction, faction)
	bld.SetTargetBuilding(targetBuilding)
	return bld
}

func loadBuildingConfig() *BuildingFactory {
	var rawStats []rawBuildingStat
	jsonFile, err := data.Files.ReadFile("building_stats.json")
	if err != nil {
		panic("failed to read building_stats.json: " + err.Error())
	}

	if err := json.Unmarshal(jsonFile, &rawStats); err != nil {
		panic("failed to parse building_stats.json: " + err.Error())
	}

	fact := &BuildingFactory{
		Buildings: make(map[types.Building]BuildingStats),
	}

	for _, raw := range rawStats {
		buildingType := types.UtilBuildingTypeFromString(raw.Name)
		fact.Buildings[buildingType] = BuildingStats{
			Name:             raw.Name,
			HPMax:            raw.HPMax,
			HPCur:            raw.HPMax,
			ResourceCost:     raw.ResourceCost,
			ConstructionTime: raw.ConstructionTime,

			VisionRange: raw.VisionRange,
			SizePx:      raw.SizePx,
		}
	}

	return fact
}
