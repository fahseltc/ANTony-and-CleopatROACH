package scene

import (
	"gamejam/sim"
	"image"
)

type SceneCompletion struct {
	CompletionArea *image.Rectangle // Not used in this example, but could be used to define a specific area for completion
	UnitOne        *sim.Unit
	UnitTwo        *sim.Unit
}

func NewSceneCompletion(first *sim.Unit, second *sim.Unit, area *image.Rectangle) *SceneCompletion {
	return &SceneCompletion{
		CompletionArea: area,
		UnitOne:        first,
		UnitTwo:        second,
	}
}

// IsComplete is a pure predicate polled every frame. It must not log or have
// side-effects: once the condition holds it stays true for the rest of the
// scene (including the whole completion cutscene), so logging here spams the
// output. The caller logs the transition exactly once.
func (sc *SceneCompletion) IsComplete(sim *sim.T) bool {
	if sc != nil && sc.CompletionArea != nil {
		if sc.CompletionArea.Overlaps(*sc.UnitOne.Rect) && sc.CompletionArea.Overlaps(*sc.UnitTwo.Rect) {
			return true
		}
	}
	return false
}
