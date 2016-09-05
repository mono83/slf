package slf

// Dispatcher is wrapper over slice of receivers
type Dispatcher struct {
	mods []Mutator
	recv []Receiver
}

// AddReceiver registers new receiver
func (d *Dispatcher) AddReceiver(r Receiver) {
	d.recv = append(d.recv, r)
}

// AddMutator registers new mutators
func (d *Dispatcher) AddMutator(m Mutator) {
	d.mods = append(d.mods, m)
}

// Receive handles incoming event
func (d *Dispatcher) Receive(e Event) {
	if len(d.recv) == 0 {
		return
	}

	// Mutating
	if len(d.mods) > 0 {
		for _, m := range d.mods {
			m.Modify(&e)
		}
	}

	// Sending
	for _, r := range d.recv {
		r.Receive(e)
	}
}
