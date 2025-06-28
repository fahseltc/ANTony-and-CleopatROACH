package audio

import (
	"bytes"
	"embed"
	"fmt"
	"gamejam/assets"
	"gamejam/eventing"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var soundFiles embed.FS

var (
	audioContext = audio.NewContext(44100)
)

type SoundManager struct {
	GlobalSFXVolume  float64
	GlobalMSXVolume  float64
	sounds           map[string][]byte
	activePlayers    map[string][]*audio.Player
	maxOverlaps      map[string]int
	lastPlayedFrame  map[string]int
	minFrameInterval map[string]int
	currentFrame     int
}

func NewSoundManager() *SoundManager {
	return &SoundManager{
		GlobalSFXVolume: 0.3,
		GlobalMSXVolume: 0.4,
		sounds:          make(map[string][]byte),
		activePlayers:   make(map[string][]*audio.Player),
		lastPlayedFrame: make(map[string]int),
		minFrameInterval: map[string]int{
			"walk": 30, // walk sound must wait X frames between plays
		},
		maxOverlaps: map[string]int{ // songs need to be in this list if we want to be able to stop them entirely.
			"sfx_command_0":    1,
			"sfx_command_1":    1,
			"sfx_command_2":    1,
			"sfx_command_3":    1,
			"sfx_command_4":    1,
			"sfx_command_5":    1,
			"msx_menusong":     1,
			"msx_narratorsong": 1,
			"msx_gamesong1":    1,
		},
	}
}

func (sm *SoundManager) LoadSound(name string, path string) {
	data, err := assets.Files.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to load sound %s: %v", name, err)
	}
	sm.sounds[name] = data
}
func (sm *SoundManager) Play(name string) {
	// Enforce frame interval limit
	if minInterval, ok := sm.minFrameInterval[name]; ok {
		if last, ok := sm.lastPlayedFrame[name]; ok && sm.currentFrame-last < minInterval {
			return // Too soon
		}
		sm.lastPlayedFrame[name] = sm.currentFrame
	}

	// Enforce overlap limit
	if limit, hasLimit := sm.maxOverlaps[name]; hasLimit {
		active := sm.activePlayers[name]
		filtered := active[:0]
		for _, p := range active {
			if p.IsPlaying() {
				filtered = append(filtered, p)
			}
		}
		sm.activePlayers[name] = filtered

		if len(filtered) >= limit {
			return
		}
	}

	data, ok := sm.sounds[name]
	if !ok {
		log.Printf("sound not found: %s", name)
		return
	}

	stream, err := wav.DecodeWithSampleRate(44100, bytes.NewReader(data))
	if err != nil {
		log.Printf("failed to decode sound %s: %v", name, err)
		return
	}
	var player *audio.Player
	if len(name) >= 4 && name[:4] == "msx_" {
		loop := audio.NewInfiniteLoop(stream, stream.Length())
		player, err = audioContext.NewPlayer(loop)
	} else {
		player, err = audioContext.NewPlayer(stream)
	}

	if err != nil {
		log.Printf("failed to create player for %s: %v", name, err)
		return
	}
	// Set volume based on prefix
	if len(name) >= 4 && name[:4] == "msx_" {
		player.SetVolume(sm.GlobalMSXVolume)
	} else {
		player.SetVolume(sm.GlobalSFXVolume)
	}
	player.Play()

	if _, hasLimit := sm.maxOverlaps[name]; hasLimit {
		sm.activePlayers[name] = append(sm.activePlayers[name], player)
	}
}

func (sm *SoundManager) Update() {
	sm.currentFrame++
}

func (sm *SoundManager) PlayRandom(prefix string, count int) {
	n := rand.Intn(count)
	sm.Play(fmt.Sprintf("%s_%d", prefix, n))
}

func (sm *SoundManager) PlayWalkSFX(event eventing.Event) {
	sm.Play("walk")
}

func (sm *SoundManager) PlayIssueActionSFX(event eventing.Event) {
	sm.PlayRandom("sfx_command", 5) //format is 'sfx_command1'
}

func (sm *SoundManager) Stop(name string) {
	players, ok := sm.activePlayers[name]
	if !ok {
		return
	}

	for _, player := range players {
		if player.IsPlaying() {
			player.Pause() // or .Close() if you want to release resources
		}
	}
	sm.activePlayers[name] = nil
}

func (sm *SoundManager) SetGlobalMSXVolume(volume float64) {
	sm.GlobalMSXVolume = volume
	for name, players := range sm.activePlayers {
		if len(name) >= 4 && name[:4] == "msx_" {
			for _, player := range players {
				if player.IsPlaying() {
					player.SetVolume(volume)
				}
			}
		}
	}
}
