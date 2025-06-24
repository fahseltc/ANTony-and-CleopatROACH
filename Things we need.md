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

- rightside UI for building bridge when a single unit is selected
- add non-standard units to sim/sprite layers (king/queen)
- static sprites added via game actions (bridges)
- removing collision rectangles when things are built
- building bridges + pathing over them
- timer for units to collect resources or build
- tooltips on buttons to show costs
- life bars / progress bars
- level system
- scripting what happens in each level, or a system for that which doesnt suck
- fullscreen text dialog
- Animation system
- walk / attack animation
- camera smooth panning for in-engine cutscenes
- units attacking - COULD BE CUT perhaps
- pathfinding
- targeting center of other objects instead of top-left corner.
- ART - makeant+makebridge btn-pressed art
- in-engine emoji stuff (ants smooching)
- text content from chez
- tutorial of any variety
- 'selected units' UI element

- blood stains on ground? static sprites that expire?
