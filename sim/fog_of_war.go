package sim

type FogState uint8

const (
	FogUnexplored FogState = iota
	FogMemory
	FogVisible
)

type FogOfWar struct {
	Enabled       bool
	Width, Height int
	Tiles         [][]FogState
}

func NewFogOfWar(width, height int) *FogOfWar {
	tiles := make([][]FogState, height)
	for y := range tiles {
		tiles[y] = make([]FogState, width)
	}
	return &FogOfWar{
		Enabled: true,
		Width:   width,
		Height:  height,
		Tiles:   tiles,
	}
}

func (f *FogOfWar) Update(sim *T) {
	// Step 1: Clear previous visibility (convert visible -> memory)
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			if f.Tiles[y][x] == FogVisible {
				f.Tiles[y][x] = FogMemory
			}
		}
	}

	for _, unit := range sim.GetAllPlayerUnits() {
		radius := int(unit.Stats.VisionRange) // in tiles
		unitPt := unit.GetTileCoordinates().ToPoint()
		cx, cy := unitPt.X, unitPt.Y

		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				x := cx + dx
				y := cy + dy
				if x < 0 || x >= f.Width || y < 0 || y >= f.Height {
					continue
				}
				if dx*dx+dy*dy <= radius*radius {
					f.Tiles[y][x] = FogVisible
				}
			}
		}
	}
	for _, bldg := range sim.GetAllBuildings() {
		radius := int(bldg.GetVisionRange())
		bldgPt := bldg.GetTilePosition()
		cx, cy := int(bldgPt.X), int(bldgPt.Y)

		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				x := cx + dx
				y := cy + dy
				if x < 0 || x >= f.Width || y < 0 || y >= f.Height {
					continue
				}
				if dx*dx+dy*dy <= radius*radius {
					f.Tiles[y][x] = FogVisible
				}
			}
		}

	}
}
