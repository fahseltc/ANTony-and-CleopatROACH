Things we need

- tilemap where tiles have traits like collision, resources, etc
- Basic commands needed: Move(only), AttackMove, Stop, HoldPosition
- unit

  - attacking
  - moving (w/ pathfinding)
  - gathering
  - states (moving, attackMoving, hold position, none, gathering)

- building
- multi-selection of units
- resources (on map / on UI / in memory)
- camera scrolling / zooming
- unit groups
- hotkeys to zoom to map locations

Simulator interface

- getUnitById(id)
- issueUnitCommand(id)
- getUnitsInArea(rect)

~3 days left
WHAT REMAINS

- [x] rightside UI for building bridge when a single unit is selected
- [x] add non-standard units to sim/sprite layers (king/queen)
- [x] tweaks to collision system
- [x] fix resource collection jank due to not using center of objects
- [x] static sprites added via game actions (bridges)
- [x] removing collision rectangles when things are built
- [x] building bridges + (CHECK) pathing over them
- [x] Indicators for keyboard hotkeys on UI
- [x] progress bars
- [x] targeting center of other objects instead of top-left corner. (still not 100% but better)
- [x] fix janky large ant collision
- [x] camera smooth panning for in-engine cutscenes
- [x] camera fading
- [x] timer for units to collect resources and build
- [x] builder must be nearby new building for it to start construction - KINDA WRONG
- BUG - building site should be made at any distance and only progress when the builder is nearby.

- level system
- scripting what happens in each level, or a system for that which doesnt suck
- fullscreen text dialog
- Animation system
- walk / attack animation
- units attacking - COULD BE CUT perhaps
- life bars?
- pathfinding
- ART - makeant+makebridge btn-pressed art
- in-engine emoji stuff (ants smooching)
- text content from chez
- tutorial of any variety
- 'selected units' UI element

other

- tooltips on buttons to show costs
- blood stains on ground? static sprites that expire?
- bugs carry resources visually
- Patrol functionality
