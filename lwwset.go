package lwwset

import (
	"sync"
	"time"
)

// LWW is a thread safe LWW-Element-Set
type LWW struct {
	mu sync.RWMutex
	m  LWWElements
}

// LWWElements stores all the elements in the lww set
type LWWElements map[interface{}]LWWTime

// LWWTime stores the add and remove times of the lww element
type LWWTime struct {
	Add    time.Time
	Remove time.Time
}

// New returns a new LWW
func New() *LWW {
	return &LWW{
		m: make(LWWElements),
	}
}

// NewFromElements returns a new LWW initialised with the provided elements
func NewFromElements(m LWWElements) *LWW {
	return &LWW{
		m: copyElemnts(m),
	}
}

// Elements returs a copy of the elements currently in LWW
func (s *LWW) Elements() LWWElements {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return copyElemnts(s.m)
}

// Add inserts an element into the set
func (s *LWW) Add(e interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.m[e]
	t.Add = latestTime(t.Add, time.Now())
	s.m[e] = t
}

// Remove removes an element from the set
func (s *LWW) Remove(e interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.m[e]
	t.Remove = latestTime(t.Remove, time.Now())
	s.m[e] = t
}

// Lookup checks if an element is a member of lww set
func (s *LWW) Lookup(e interface{}) bool {
	s.mu.RLock()
	t, ok := s.m[e]
	s.mu.RUnlock()

	// first check if the element is in add set
	if ok == false {
		return false
	}

	// now compare add and remove times
	if t.Remove.Sub(t.Add) >= 0 {
		// biased towards removals for this implementation
		return false
	}
	return true
}

// Equal checks whether or not two sets are equal
func (s *LWW) Equal(new *LWW) bool {
	a := s.Elements()
	b := new.Elements()

	if len(a) != len(b) {
		return false
	}

	// compare each element
	for e, t := range a {
		t2, ok := b[e]
		if !ok ||
			!t.Add.Equal(t2.Add) ||
			!t.Remove.Equal(t2.Remove) {
			return false
		}
	}
	return true
}

// Merge merges new into s
func (s *LWW) Merge(new *LWW) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// merge each element
	for e, t := range new.Elements() {
		t2, ok := s.m[e]
		if !ok {
			s.m[e] = t
		} else {
			t.Add = latestTime(t.Add, t2.Add)
			t.Remove = latestTime(t.Remove, t2.Remove)
			s.m[e] = t
		}
	}
}

// latestTime is a helper to return the latest time
func latestTime(a time.Time, b time.Time) time.Time {
	if a.Sub(b) < 0 {
		return b
	}
	return a
}

// copyElemnts is a helper to deep copy LWWElements
func copyElemnts(m LWWElements) LWWElements {
	new := make(LWWElements, len(m))
	for e, t := range m {
		new[e] = t
	}
	return new
}
