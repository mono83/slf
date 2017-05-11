package context

import "sync"

// MaxPoolSize is max amount of contexts inside pool
var MaxPoolSize = 1000

var pool = make(map[string]Interface, 1000)
var mutex = sync.Mutex{}

// register places context into pool
func register(i Interface) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(pool) < MaxPoolSize {
		pool[i.ID().String()] = i
	}
}

// unregister removes context from pool
func unregister(i Interface) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(pool, i.ID().String())
}

// ListActive returns list of currently active tokens
func ListActive() []Interface {
	mutex.Lock()
	defer mutex.Unlock()

	if len(pool) == 0 {
		return []Interface{}
	}

	list := make([]Interface, len(pool))
	i := 0
	for _, c := range pool {
		list[i] = c
		i++
	}

	return list
}
