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
- [x] mouse click indicator on text dialogs
- [x] in-engine emoji stuff (ants smooching)
- [x] tutorial of any variety
- [x] level system (kinda but jank)
- [x] scripting what happens in each level, or a system for that which doesnt suck
- [x] ART - makeant+makebridge btn-pressed art
- [x] fix bug with units not finding nearest hive
- [x] STARTED - Second Level! - It should be about having two bases and using resources gathered from the other to make progress.
- [x] fullscreen text narrator
- [x] movement based on unit center instead of top-left
- [x] MAKE IT WORK ON THE WEBSITE (slashes fix in TMX)
- [x] finish level 2 scripting intro + ending scenario (love in flower area)
- [x] more tutorial
- [x] some UI text about "not enough money to build X"
- [x] SFX - need more
- [x] Animation system
- [x] walk animation
- [x] bugs carry resources visually
- [x] music
- level 2 tutorial?
- graphical feedback when clicking on things
- SFX - hive click, unit build, construction, levelsuccess
- 'selected units' UI element
- hotkeys to unit groups
- hotkeys to saved areas
- BUG - building site should be made at any distance and only progress when the builder is nearby.
- BUG - hives sometimes build units on top of other units
- [x] BUG - units are selected after initial cutscene - WHY?
- better pathfinding

- tooltips on buttons to show costs
- Patrol functionality

COMBAT

- blood stains on ground? static sprites that expire?
- attack animation
- units attacking - COULD BE CUT perhaps
- life bars?

funny chatgpt quotes
“A sting! A sting! My kingdom for a sting!”
“Buzz off, thou knavish gnat!”
“We were six legs in a pod…”
“Is this a crumb I see before me?”
“No more scuttling — this calls for full-wing flight!”
“He speaks like a drunken firefly — bright, but makes no sense.”
“She hath a thorax like no other.”
“Even the termites bow before Cleopatroach.”
“The roach is mightier than the flea.”

Antony preparing for battle

    Antony: “Let the bugle call, for war is upon us! Sharpen thy stingers!”

Cleopatroach taunting a rival

    Cleopatroach: “Thou beetle-brained scuttler! Darest thou approach the queen of Nile-larvae?”

Antony contemplating his duty

    Antony: “Between love and colony, my feelers are torn. Yet I must march!”

Cleopatroach showing affection

    Cleopatroach: “I wore perfume of the dung beetle, just for thee.”

Antony mourning a lost ant

    Antony: “O brave Foragerius! He hath fallen 'neath a human’s sandal...”

Cleopatroach, dramatic as ever

    Cleopatroach: “My thorax trembles with longing! Come, sweet Antony, before the pheromones fade!”

NARRATOR 1
In the land of Nilopolis, where the sand meets sugar and the air hums with wingèd gossip, two empires crawl toward destiny. One: the mighty Ant-tonian Legion, proud builders and brave foragers. The other: Queen Cleopatroach’s royal roachdom, ancient, secretive, and ever-scheming.
Long hath love fluttered betwixt Antony, soldier of soil, and Cleopatroach, goddess of grime. But lo! A chasm divides them, wide as a footprint and deep as a drain. Wood must be gathered. A bridge must be built. And their love… must scuttle onward.
Arise, player! Command thy swarm! The Ant Game begins!

    Narrator (Post-Level 1 Completion)

“And so, with twig and grit and tiny limbs unnumbered, the bridge was wrought! Across the void did Antony march, bold as a wasp in a wineglass, toward his queen!”

“Their antennae touched in triumph, and though their path remaineth plagued with pest and peril, this moment—this mandible-kiss of destiny—shall be sung in the tunnels of time.”

“But peace is never long in the insect world. New foes stir. New challenges arise. Onward… to Level the Second!”
