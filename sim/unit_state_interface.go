package sim

type UnitStateInterface interface {
	Enter(unit *Unit)
	Update(unit *Unit, sim *T)
	Exit(unit *Unit)
	Name() string
}

// Unit States:
// AttackingState
// IdleState
// MovingState
//
// HarvestingState
// DeliveringState

// type ExampleState struct{}
// func (s *ExampleState) Enter(unit *Unit) {}
// func (s *ExampleState) Update(unit *Unit, sim *T) {}
// func (s *ExampleState) Exit(unit *Unit) {}
// func (s *ExampleState) Name() string    { return "example" }
