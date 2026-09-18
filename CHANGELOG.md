# Changelog

All notable changes to ANTony & CleopatROACH are documented here.

## V2

Post-jam improvements: a playable second faction, UI/UX polish, more robust
cutscenes and level flow, and several gameplay/logging fixes.

### Added
- **Roach hive unit production.** The roach hive can now build units the same
  way ant hives do. Selecting a roach hive shows the build panel, and its worker
  button produces roach workers (`UnitTypeDefaultRoach`) instead of ants.
- **Dedicated roach hive UI.** The roach hive panel uses the new
  `make-roach-btn` icon (and its pressed variant) and shows the roach worker's
  own cost tooltip.
- **`sim.IssueHarvestTile(unitID, tileX, tileY)`.** A tile-based helper to send
  a worker to start harvesting a chosen resource tile from level setup or a
  cutscene (used to start a worker mining in the Level 1 setup).
- **Level 2 tutorial completion conditions.** The multi-hive / control-group
  tutorial steps now complete on the right in-game action (select the ant hive,
  bind Ctrl+1, select the roach hive, bind Ctrl+2, double-tap 1 to recenter),
  backed by new `UnitGroupManager.HasGroup` and `DidRecenterOnGroup` helpers.
- **Cutscene logging with progress.** Each cutscene action is logged through the
  JSON logging system as it starts and finishes, including `current` / `max`
  progress through the cutscene.

### Changed
- **Notifications wrap on screen.** Long notification messages are now
  word-wrapped to fit the screen (still honoring explicit newlines) instead of
  running off both edges.
- **Warnings routed through the logging system.** The "tile clicked was not
  found" warning (and related sim messages) now go through the structured JSON
  logger and include the offending coordinates, instead of the default
  plain-text logger.
- **End-of-campaign flow.** Finishing the final level returns to the main menu
  instead of loading an empty level.

### Fixed
- **Harvest distance too long.** Workers no longer harvest from a tile or two
  away. Previously a worker committed to harvesting after a single approach
  regardless of distance; it now keeps closing to within harvest range (with a
  bounded retry so a genuinely boxed-in worker still gathers).
- **Button tooltips misaligned.** Fixed hive-panel (and all) button tooltips
  whose background and text desynced — bottom-row buttons drew their tooltip in
  the wrong place. Tooltips now stay aligned and clamped on screen as one unit.
- **Hotkey buttons ignored hidden UI.** The unit control-group hotkey buttons
  now hide along with the rest of the HUD (e.g. during cutscenes).
- **Roach hive advanced unit.** The roach hive panel no longer shows the
  fighter/advanced unit button.
- **Out-of-bounds cutscene target.** Fixed a `missing type in composite literal`
  compile error and a cutscene unit command that targeted an off-map tile (so
  the unit silently never moved); out-of-bounds harvest/command targets are now
  logged.

## V1

Original game jam release.

- Real-time strategy game built in Go with the Ebiten engine.
- Command a swarm of ants: drag-box select workers and right-click to harvest
  sucrose or wood, carried back to a hive.
- Spend resources at a hive to produce more units.
- Build bridges over water so units can cross, opening otherwise-impassable
  tiles without removing the underlying collision.
- Royal units (Antony the ant and Cleopatroach the roach) with a larger hitbox;
  each level ends when the two royals meet in the level's completion area.
- Fog of war, minimap, camera panning/zoom, unit control groups, and a
  pub/sub event bus wiring UI actions to the simulation.
- Tiled (TMX) maps, embedded assets, music and SFX, portrait dialog cutscenes,
  and per-level intro narration and tutorials.
