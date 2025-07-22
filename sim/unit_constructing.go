package sim

type ConstructingState struct{}

func (s *ConstructingState) Enter(unit *Unit)          {}
func (s *ConstructingState) Update(unit *Unit, sim *T) {}
func (s *ConstructingState) Exit(unit *Unit)           {}
func (s *ConstructingState) GetName() string {
	return UnitStateConstructing.ToString()
}
