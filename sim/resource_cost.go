package sim

type ResourceCost struct {
	Sucrose uint
	Wood    uint
}

func (tc *ResourceCost) CanAfford(playerState PlayerState) bool {
	if playerState.Sucrose >= tc.Sucrose && playerState.Wood >= tc.Wood {
		return true
	}
	return false
}
func (tc *ResourceCost) Purchase(playerState *PlayerState) bool {
	if tc.CanAfford(*playerState) {
		playerState.Sucrose = playerState.Sucrose - tc.Sucrose
		playerState.Wood = playerState.Wood - tc.Wood
		return true
	}
	return false
}
