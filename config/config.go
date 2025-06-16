package config

import (
	"encoding/json"

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
func New() *T {
	var cfg T
	jsonFile, _ := data.Files.ReadFile("config.json")
	_ = json.Unmarshal(jsonFile, &cfg)
	return &cfg
}
