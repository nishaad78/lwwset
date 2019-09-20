package lwwset

import (
	"sync"
	"time"
)

// LWW is a thread safe LWW-Element-Set
type LWW struct {
	mu sync.RWMutex
	m  Elements
}

// Elements stores all the elements in the lww set
type Elements map[interface{}]ElementState

// ElementState stores the element state and the last modified time (in nanoseconds)
type ElementState struct {
	IsRemoved bool
	UpdatedAt int64
}

// New returns a new LWW
func New() *LWW {
	return &LWW{
		m: make(Elements),
	}
}

// NewFromMap returns a new LWW initialised with the provided elements
func NewFromMap(m Elements) *LWW {
	return &LWW{
		m: copyElements(m),
	}
}

// Map returs a copy of the elements currently in LWW
func (s *LWW) Map() Elements {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return copyElements(s.m)
}

// Add inserts an element into the set
func (s *LWW) Add(e interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.m[e]
	now := time.Now().UnixNano()
	if now > t.UpdatedAt {
		t.UpdatedAt = now
		t.IsRemoved = false
	}
	s.m[e] = t
}

// Remove removes an element from the set
func (s *LWW) Remove(e interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.m[e]
	now := time.Now().UnixNano()
	if now >= t.UpdatedAt {
		// biased towards removals for this implementation
		t.UpdatedAt = now
		t.IsRemoved = true
	}
	s.m[e] = t
}

// Lookup checks if an element is a member of the set
func (s *LWW) Lookup(e interface{}) bool {
	s.mu.RLock()
	t, ok := s.m[e]
	s.mu.RUnlock()

	return ok && t.IsRemoved == false
}

// Elements returns a list of all the valid members of the set
func (s *LWW) Elements() []interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := []interface{}{}
	for e, t := range s.m {
		if t.IsRemoved == false {
			res = append(res, e)
		}
	}

	return res
}

// Equal checks whether or not two sets are equal
func (s *LWW) Equal(new *LWW) bool {
	a := s.Map()
	b := new.Map()

	if len(a) != len(b) {
		return false
	}

	// compare each element
	for e, t := range a {
		t2, ok := b[e]
		if !ok ||
			t.IsRemoved != t2.IsRemoved ||
			t.UpdatedAt != t2.UpdatedAt {
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
	for e, t := range new.Map() {
		t2, ok := s.m[e]
		if !ok {
			s.m[e] = t
		} else if t.UpdatedAt > t2.UpdatedAt {
			s.m[e] = t
		} else if t.UpdatedAt == t2.UpdatedAt && t.IsRemoved {
			// biased towards removals
			s.m[e] = t
		}
	}
}

// copyElements is a helper to deep copy LWWElements
func copyElements(m Elements) Elements {
	new := make(Elements, len(m))
	for e, t := range m {
		new[e] = t
	}
	return new
}
