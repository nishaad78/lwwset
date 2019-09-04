package lwwset

import (
	"time"
)

// LWW is a thread safe LWW-Element-Set
type LWW struct {
	a *g
	r *g
}

// NewLWW returns a new LWW
func NewLWW() *LWW {
	s := &LWW{
		a: newG(),
		r: newG(),
	}
	return s
}

// Add inserts an element into the set
func (s *LWW) Add(e interface{}, t time.Time) {
	s.a.add(e, t)
}

// Remove removes an element from the set
func (s *LWW) Remove(e interface{}, t time.Time) {
	s.r.add(e, t)
}

// Lookup checks if an element is a member of lww set and returns latest recorded time
// returned bool is false if element is not a member of the set
func (s *LWW) Lookup(e interface{}) (time.Time, bool) {
	t, ok := s.a.lookup(e)
	removeT, removeOk := s.r.lookup(e)

	// first check if the element is in add set
	if ok == false {
		return t, false
	}

	// now check if element exists in remove set and compare the timestamps
	if removeOk == true {
		timeDiff := removeT.Sub(t)
		if timeDiff >= 0 {
			// using remove bias for this implementation
			return time.Time{}, false
		}
	}
	return t, true
}

// Equal checks whether or not two sets are equal
func (s *LWW) Equal(new *LWW) bool {
	return s.a.equal(new.a) && s.r.equal(new.r)
}

// Merge merges new into s
func (s *LWW) Merge(new *LWW) {
	s.a.merge(new.a)
	s.r.merge(new.r)
}
