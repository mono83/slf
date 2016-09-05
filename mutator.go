package slf

// Mutator interface describes components, that modifies events
type Mutator interface {
	Modify(*Event)
}
