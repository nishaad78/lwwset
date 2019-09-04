package lwwset

import (
	"sync"
	"time"
)

// g is a set based on Grow-only Set which also records timestamp
type g struct {
	mu sync.RWMutex
	m  map[interface{}]time.Time
}

func newG() *g {
	s := &g{m: make(map[interface{}]time.Time)}
	return s
}

func (s *g) add(e interface{}, t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	upsert(s.m, e, t)
}

// lookup returns latest recorded time of element in set
// returned boolean is false if element does not exist in set
func (s *g) lookup(e interface{}) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.m[e]

	return t, ok
}

func (s *g) equal(new *g) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	new.mu.RLock()
	defer new.mu.RUnlock()

	if len(new.m) != len(s.m) {
		return false
	}

	// compare each element
	for e, t := range s.m {
		if t2, ok := new.m[e]; !ok || t.Sub(t2) != 0 {
			return false
		}
	}

	return true
}

func (s *g) merge(new *g) {
	s.mu.Lock()
	defer s.mu.Unlock()

	new.mu.RLock()
	defer new.mu.RUnlock()

	for e, t := range new.m {
		upsert(s.m, e, t)
	}
}

// upsert is a helper to update or add element to map
func upsert(m map[interface{}]time.Time, e interface{}, t time.Time) {
	val, ok := m[e]
	if ok && t.Sub(val) <= 0 {
		return
	}
	// assign only if t is newer or element doesn't exist in map
	m[e] = t
}
