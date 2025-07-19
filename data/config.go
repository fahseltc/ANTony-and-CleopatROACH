package data

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	WindowTitle    string `json:"windowTitle"`
	TargetFPS      uint   `json:"targetFPS"`
	SkipMenu       bool   `json:"skipMenu"`
	SkipToGameplay bool   `json:"skipToGameplay"`
	StartingLevel  int    `json:"startingLevel"`
	DebugDraw      bool   `json:"debugDraw"`
	MuteAudio      bool   `json:"muteAudio"`
	Resolutions    struct {
		Internal Resolution `json:"internal"`
		External Resolution `json:"external"`
	} `json:"resolution"`
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
