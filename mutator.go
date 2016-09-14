package slf

// Mutator interface describes components, that modifies events
type Mutator interface {
	Modify(*Event)
}

// FuncMutator returns Mutator built over func
func FuncMutator(f func(*Event)) Mutator {
	return callableMutator{invoke: f}
}

type callableMutator struct {
	invoke func(*Event)
}

func (cm callableMutator) Modify(e *Event) {
	cm.invoke(e)
}
