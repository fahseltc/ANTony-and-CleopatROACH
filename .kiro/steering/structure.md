# Project structure & architecture

Go module `gamejam` (see `go.mod`), Go 1.24. Entry point is `main.go` → `game`
package. Rendering/engine is Ebiten; pathfinding uses `github.com/quasilyte/pathing`;
maps load via `github.com/lafriks/go-tiled`.

## Package map

- `main.go` — process entry; constructs the `game.Game` and runs Ebiten.
- `game/` — `Game` implements Ebiten's `Update`/`Draw`/`Layout`. Owns the scene
  manager and font loading.
- `scene/` — scene layer and game flow (this is where most gameplay wiring lives):
  - `play.go` — `PlayScene`, the main gameplay scene. Owns the `sim.T`, `ui.Ui`,
    tilemap, sprites, cutscene actions, tutorials, and the completion condition.
    `Update` drives the sim, cutscenes, tutorials, and sprite sync each frame.
  - `level.go` — `LevelCollection` / `LevelData`: per-level setup, intro and
    completion cutscenes, and the ordered tutorial steps.
  - `cutscene.go` — `CutsceneAction` interface and concrete actions (camera pans,
    unit commands, portrait text, fog reveals, waits). Actions run one-per-frame
    and are popped when their `Update` returns true.
  - `tutorial.go` — `Tutorial` interface and `TutorialStep`. Steps advance one at
    a time via a `CompleteFunc` predicate. `NewTutorialStepNoClick` requires a
    non-nil completeFunc (steps whose action is a click must NOT dismiss on click).
  - `event_handler_manager.go` — subscribes to the event bus (build clicks,
    notifications, etc.).
  - `completion.go` — `SceneCompletion`, a pure predicate polled every frame.
  - `narrator.go`, `menu.go`, `base.go` — other scenes and the scene base type.
- `sim/` — the simulation/game logic, engine-agnostic where possible:
  - `sim.go` — `sim.T`, the world/simulation root. Unit & building maps, the
    `World` (tilemap, map objects, fog), pathing entry points (`FindClickedPath`,
    `FindNearestSurroundingWalkableTiles`), collision queries
    (`GetAllCollidersOverlapping`, `IsCoveredByBridge`), and `ConstructBuilding`.
  - `unit.go` — `Unit`, `UnitStats`, unit constructors (`NewRoyalAnt`, etc.),
    tile/position helpers. `RoyalUnitSize` is the enlarged royal hitbox.
  - `unit_state_*.go` — the unit state machine (Idle, Moving, Harvesting,
    Delivering, Attacking, AttackMove, Constructing). `unit_state_moving.go` holds
    `isColliding` and the bridge-aware collision (`mapCollides`).
  - `building*.go` — building types, factory, in-construction, hive, bridge.
  - `fog_of_war.go` — visibility grid.
- `tilemap/` — `Tilemap` wraps the loaded TMX. Builds the `pathing.Grid` in
  `GenerateTiles` (collision derived from TMX "collision" object group). Bridges
  use `SetTileWalkable` + `walkableOverrides`. `tile.go` is the per-tile record.
- `ui/` — camera, HUD, buttons, sprites/animation, drag select, construction
  ghost mouse, minimap, notifications, pause, unit groups.
- `eventing/` — a small pub/sub `EventBus` used to decouple UI clicks from sim.
- `types/` — shared enums (unit types, building types, tile types, resources).
- `vec2/` — 2D vector math.
- `audio/` — `SoundManager` (music + SFX), event-driven.
- `data/` — runtime config (`data.Config`) and `building_stats.json`.
- `assets/` — embedded assets via `go:embed` (`assets.Files`): tilemaps, sprites,
  fonts, music, SFX, tutorials, portraits.
- `util/` — image loading/scaling, queue, and `OpenURL` (build-tagged for desktop
  vs. `js/wasm`).
- `log/` — logging helper.

## Coordinate conventions

- Tiles are 128px (`TileSize = 128`, `HalfTileSize = 64`).
- A `Unit.Position` is the rect top-left (Min); `GetCenteredPosition` adds half the
  rect size. `GetTileCoordinates` rounds `Position / TileSize`.
- Tile coords vs. pixel coords are distinct — watch conversions (`*128` / `/128`).

## Cross-cutting patterns

- **Event bus**: UI publishes events; `EventHandlerManager` and `sim` subscribe.
  Prefer this over direct cross-layer calls for UI→sim actions.
- **State machine**: unit behavior is driven by `UnitStateInterface`
  (`Enter`/`Update`/`Exit`/`GetName`). Add behavior as a state, not ad-hoc flags.
- **Bridges**: never delete the water collision rect. Open the tile on the grid
  (`SetTileWalkable`) and keep runtime collision bridge-aware
  (`IsCoveredByBridge` / `mapCollides`). Only "buildable" TMX tiles accept bridges
  (`IsBridgeBuildable`).
- **Pathing grid staleness**: `Tilemap.MapObjects` is reassigned on add/remove;
  read collision from `world.TileMap.MapObjects` (live), not a cached snapshot.

## Conventions

- Match existing style; keep changes scoped. No test framework is set up — verify
  with `go build ./...` and `go vet ./...`.
- Ignore the `__debug_bin*.exe` files (Delve debug binaries) and `gamejam.exe`.
