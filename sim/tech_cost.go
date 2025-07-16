package sim

type TechCost struct {
	Sucrose uint
	Wood    uint
}

func (tc *TechCost) CanAfford(playerState PlayerState) bool {
	if playerState.Sucrose >= tc.Sucrose && playerState.Wood >= tc.Wood {
		return true
	}
	return false
}
func (tc *TechCost) Purchase(playerState *PlayerState) bool {
	if tc.CanAfford(*playerState) {
		playerState.Sucrose = playerState.Sucrose - tc.Sucrose
		playerState.Wood = playerState.Wood - tc.Wood
		return true
	}
	return false
}
