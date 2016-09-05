package slf

// Receiver interface represents SLF event receiver
type Receiver interface {
	Receive(e Event)
}
