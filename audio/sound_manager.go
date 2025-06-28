package audio

import (
	"bytes"
	"embed"
	"gamejam/assets"
	"gamejam/eventing"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var soundFiles embed.FS

var (
	audioContext = audio.NewContext(44100)
)

type SoundManager struct {
	sounds           map[string][]byte
	activePlayers    map[string][]*audio.Player
	maxOverlaps      map[string]int
	lastPlayedFrame  map[string]int
	minFrameInterval map[string]int
	currentFrame     int
}

func NewSoundManager() *SoundManager {
	return &SoundManager{
		sounds:          make(map[string][]byte),
		activePlayers:   make(map[string][]*audio.Player),
		lastPlayedFrame: make(map[string]int),
		minFrameInterval: map[string]int{
			"walk": 30, // walk sound must wait X frames between plays
		},
		maxOverlaps: map[string]int{
			"walk": 3, // still keep overlap limit
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
func (sm *SoundManager) Play(name string, volume float64) {
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

	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		log.Printf("failed to create player for %s: %v", name, err)
		return
	}
	player.SetVolume(volume)
	player.Play()

	if _, hasLimit := sm.maxOverlaps[name]; hasLimit {
		sm.activePlayers[name] = append(sm.activePlayers[name], player)
	}
}

func (sm *SoundManager) Update() {
	sm.currentFrame++
}

// func (sm *SoundManager) PlayRandom(prefix string, count int) {
// 	n := rand.Intn(count)
// 	sm.Play(fmt.Sprintf("%s_%d", prefix, n), 1.0)
// }

func (sm *SoundManager) PlayWalkSFX(event eventing.Event) {
	sm.Play("walk", 0.1)
}
