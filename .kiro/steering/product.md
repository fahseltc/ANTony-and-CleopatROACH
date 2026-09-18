# ANTony & CleopatROACH

A 2D real-time strategy game written in Go using the [Ebiten](https://ebitengine.org/)
engine. The player commands a swarm of ants: gathering resources (sucrose and
wood), building structures, and constructing bridges to cross water so the royal
ant (Antony) can reunite with the royal roach (Cleopatroach).

## Core gameplay loop

- Select worker units (drag-box select) and right-click to send them to harvest
  sucrose or wood, which they carry back to a hive.
- Spend resources at a hive to produce more units, or at a barracks for fighters.
- Build bridges over water crossings so units can traverse the map.
- Each level ends when the two royal units meet inside the level's completion area.

## Key domain concepts

- **Factions**: player vs. enemy. `PlayerFaction = 0`.
- **Units**: workers (DefaultAnt/DefaultRoach), fighters, and royals
  (RoyalAnt/RoyalRoach). Royals have a larger collision hitbox and cannot cross a
  single-tile-wide bridge.
- **Buildings**: hives (produce units), barracks (unlock fighters), and bridges
  (walkable structures placed only on map-authored "buildable" water tiles).
- **Bridges** open an otherwise-impassable water tile without deleting the
  underlying collision rect; the tile is force-marked walkable on the pathing grid
  and the runtime collider is made bridge-aware.

## Build & run

- Desktop: `go build ./...` then run `gamejam.exe`, or `go run gamejam`.
- Web (WASM): `GOOS=js GOARCH=wasm go build -o gamejam.wasm gamejam`, served from
  `html/`. Deployed to GitHub Pages and itch.io via `.github/workflows/gh-pages.yml`.
- After any code change, run `go build ./...` and `go vet ./...` to verify.
- There is no test suite; verify behavior by building and (where possible) running.
