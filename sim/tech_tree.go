package sim

type TechID string

const (
	TechFasterGathering  TechID = "faster_gathering"
	TechBuildFighterUnit TechID = "build_fighter"
)

type Tech struct {
	ID              TechID
	Name            string
	Description     string
	Cost            ResourceCost
	ResearchSeconds uint
	Dependencies    []TechID
	Unlocked        bool
}

type TechTree struct {
	AvailableTech map[TechID]*Tech // all known techs
	UnlockedTech  map[TechID]bool  // which techs have been unlocked
}

func NewTechTree() *TechTree {
	return &TechTree{
		AvailableTech: map[TechID]*Tech{
			TechFasterGathering: {
				ID:              TechFasterGathering,
				Name:            "Faster Gathering",
				Description:     "Workers collect resources 20% faster",
				Cost:            ResourceCost{Sucrose: 250},
				ResearchSeconds: 60,
				Dependencies:    nil,
			},
			TechBuildFighterUnit: {
				ID:              TechBuildFighterUnit,
				Name:            "Build Fighter Unit",
				Description:     "Unlocks the ability to build Fighter units",
				Cost:            ResourceCost{Sucrose: 0, Wood: 0}, // Unlocked via building
				ResearchSeconds: 0,
				Dependencies:    nil,
			},
		},
		UnlockedTech: map[TechID]bool{},
	}
}

func (tt *TechTree) CanResearch(id TechID) bool {
	tech, ok := tt.AvailableTech[id]
	if !ok || tt.UnlockedTech[id] {
		return false
	}
	for _, dep := range tech.Dependencies {
		if !tt.UnlockedTech[dep] {
			return false
		}
	}
	return true
}

func (tt *TechTree) Unlock(id TechID, playerState *PlayerState) bool {
	tech, ok := tt.AvailableTech[id]
	if !ok || tt.UnlockedTech[id] || !tt.CanResearch(id) {
		return false
	}
	if tech.Cost.CanAfford(*playerState) {
		success := tech.Cost.Purchase(playerState)
		if success {
			tt.UnlockedTech[id] = true
			return true
		}
	}
	return false
}
func (tt *TechTree) GetDescription(id TechID) string {
	return tt.AvailableTech[id].Description
}
