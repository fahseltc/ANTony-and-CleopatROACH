package sim

// // Worker is a resource collector that has one job: collect
// // resources and bring them to a Hive. It is very similar to
// // a Unit but with far less capability, so it represented
// // separately for that reason
// type Worker struct {
// 	X, Y             float64
// 	Action           Action
// 	AssignedResource *Resource
// 	AssignedHive     Hive

// 	// destX / destY denote where the worker is currently instructed to move to
// 	destX, destY float64

// 	// CollectionSpeed denotes the amount of resource per Tick()
// 	// a Worker collects. So if you want 1 resource per second,
// 	// CollectionSpeed: 1/tps
// 	collectionSpeed float64
// 	carrying        float64
// 	maxCarry        float64

// 	moveSpeed float64
// }

// // Tick processes Worker state per simulation tick. Before you
// // call Tick(), make sure the Worker is assigned to a Hive
// func (w *Worker) Tick() {
// 	if w.Action == Idle {
// 		return
// 	}
// 	if w.Action == Collecting {
// 		w.carrying += w.collectionSpeed
// 		w.carrying = max(w.carrying, w.maxCarry)
// 		if w.carrying == w.maxCarry {
// 			w.ReturnToHive()
// 			return
// 		}
// 	}
// 	if w.Action == Delivering {
// 		// TODO - does it take time or is it instantaneous?
// 		// For now: instantaneous
// 		// TODO Hive.AddResource(type, w.carrying)
// 		w.carrying = 0
// 		w.GoToResource()
// 		return
// 	}

// 	if w.Action == Moving {
// 		// If we're at the resource: start collection
// 		if w.X == w.AssignedResource.X && w.Y == w.AssignedResource.Y {
// 			w.Action = Collecting
// 			return
// 		}
// 		// If we're at the Hive: start delivery
// 		if w.X == w.AssignedHive.X && w.Y == w.AssignedHive.Y {
// 			w.Action = Delivering
// 			return
// 		}
// 		// If we're at the destination, which is neither the
// 		// intended resource OR the hive, idle like a boss
// 		if w.X == w.destX && w.Y == w.destY {
// 			w.Action = Idle
// 			return
// 		}

// 		// TODO pathfinding
// 		if w.X < w.destX {
// 			w.X += w.moveSpeed
// 		} else if w.X > w.destX {
// 			w.X -= w.moveSpeed
// 		}
// 		if w.Y < w.destY {
// 			w.Y += w.moveSpeed
// 		} else if w.Y > w.destY {
// 			w.Y -= w.moveSpeed
// 		}
// 	}
// }

// // ReturnToHive instructs the worker to stop what it's doing
// // and start moving back to the assigned Hive
// func (w *Worker) ReturnToHive() {
// 	w.Action = Moving
// 	w.destX = w.AssignedHive.X
// 	w.destY = w.AssignedHive.Y
// }

// // GoToResource instructs the worker to stop what it's doing
// // and start moving back to the assigned resource
// func (w *Worker) GoToResource() {
// 	w.Action = Moving
// 	w.destX = w.AssignedResource.X
// 	w.destY = w.AssignedResource.Y
// }

// // Collect instructs a unit to constantly move between
// // its spawn point and the given resource
// func (w *Worker) Collect(resource *Resource) {
// 	w.AssignedResource = resource
// 	w.AssignedResource.X = resource.X
// 	w.AssignedResource.Y = resource.Y
// }
