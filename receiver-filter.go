package slf

type receiverFilter struct {
	real      Receiver
	predicate func(e Event) bool
}

func (r receiverFilter) Receive(e Event) {
	if r.predicate(e) {
		r.real.Receive(e)
	}
}

// Filter builds new filter wrapper over receiver
func Filter(r Receiver, predicate func(e Event) bool) Receiver {
	return receiverFilter{real: r, predicate: predicate}
}
