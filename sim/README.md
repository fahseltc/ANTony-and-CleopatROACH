# Simulation Container

## Game Start

Player starts with a single worker, spawned at the starting hive location. This worker
must be assigned by the player to go collect a particular resource.

## Workers

Workers are a special type of unit in that their only function is to move between
a Hive and a resource location, collect at the resource location, then return to
the Hive to increment that resource. It cannot attack, but it _can_ be moved by
the player to a manual location, if necessary. Doing so interrupts the collection
loop until a unit is told to resume collection.

## Units
