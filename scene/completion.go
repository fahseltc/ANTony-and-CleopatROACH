package scene

import (
	"gamejam/log"
	"gamejam/sim"
	"image"
	"log/slog"
)

type SceneCompletion struct {
	CompletionArea *image.Rectangle // Not used in this example, but could be used to define a specific area for completion
	UnitOne        *sim.Unit
	UnitTwo        *sim.Unit
	log            *slog.Logger
}

func NewSceneCompletion(first *sim.Unit, second *sim.Unit, area *image.Rectangle) *SceneCompletion {
	return &SceneCompletion{
		CompletionArea: area,
		UnitOne:        first,
		UnitTwo:        second,
		log:            log.NewLogger().With("for", "SceneCompletion"),
	}
}

func (sc *SceneCompletion) IsComplete(sim *sim.T) bool {
	if sc.CompletionArea != nil {
		if sc.CompletionArea.Overlaps(*sc.UnitOne.Rect) && sc.CompletionArea.Overlaps(*sc.UnitTwo.Rect) {
			sc.log.Warn("Scene completion condition met")
			return true
		}
	}
	return false
}
