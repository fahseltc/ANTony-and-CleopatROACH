package data

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	WindowTitle string `json:"windowTitle"`
	TargetFPS   uint   `json:"targetFPS"`
	MuteAudio   bool   `json:"muteAudio"`
	Resolutions struct {
		Internal Resolution `json:"internal"`
		External Resolution `json:"external"`
	} `json:"resolution"`
	// Dev holds developer/debug-only toggles. These are intended to be off in
	// shipping builds; keeping them under a single key makes it obvious which
	// settings are real app configuration and which are dev conveniences.
	Dev DevConfig `json:"dev"`
}

// DevConfig groups developer/debug toggles that should not be enabled in a
// shipping build.
type DevConfig struct {
	// DebugDraw enables the debug overlay (collision rects, etc.).
	DebugDraw bool `json:"debugDraw"`
	// DebugStartingResources, when positive, grants the player this amount of
	// each resource type (sucrose and wood) at the start of a play scene.
	// Zero or negative values are ignored.
	DebugStartingResources int `json:"debugStartingResources"`
	// SkipMenu jumps straight to the starting level's intro narration.
	SkipMenu bool `json:"skipMenu"`
	// SkipToGameplay jumps straight into the starting level's play scene.
	SkipToGameplay bool `json:"skipToGameplay"`
	// SkipCutscenes removes the intro cutscene from play scenes.
	SkipCutscenes bool `json:"skipCutscenes"`
	// SkipTutorial removes the tutorial dialogs from play scenes.
	SkipTutorial bool `json:"skipTutorial"`
	// StartingLevel selects which level to start on when using the menu START
	// button or the skip toggles above.
	StartingLevel int `json:"startingLevel"`
}

type Resolution struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

// New loads config keys from a root-level config.json
func NewConfig() (*Config, error) {
	var cfg Config
	jsonFile, err := Files.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}
	return &cfg, nil
}
