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
