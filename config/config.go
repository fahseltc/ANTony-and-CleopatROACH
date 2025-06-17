package config

import (
	"encoding/json"
	"fmt"

	"gamejam/data"
)

// T represents a Config file for the Game
type T struct {
	WindowTitle string `json:"windowTitle"`
	TargetFPS   uint   `json:"targetFPS"`
	DebugDraw   bool   `json:"debugDraw"`
	Resolutions struct {
		Internal Resolution `json:"internal"`
		External Resolution `json:"external"`
	} `json:"resolution"`
}

type Resolution struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

// New loads config keys from a root-level config.json
func New() (*T, error) {
	var cfg T
	jsonFile, err := data.Files.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}
	return &cfg, nil
}
