package lwwset

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLWWBasicAddRemove(t *testing.T) {
	s := New()
	ok := s.Lookup('a')
	require.False(t, ok)

	// add 'a'
	s.Add('a')
	ok = s.Lookup('a')
	require.True(t, ok)

	// remove 'a'
	s.Remove('a')
	ok = s.Lookup('a')
	require.False(t, ok)

	// add 'a' again
	s.Add('a')
	ok = s.Lookup('a')
	require.True(t, ok)
}

func TestLWWRemoveBias(t *testing.T) {
	now := time.Now().UnixNano()

	s := NewFromMap(Elements{'a': ElementState{
		IsRemoved: false,
		UpdatedAt: now,
	}})
	ok := s.Lookup('a')
	require.True(t, ok)

	s2 := NewFromMap(Elements{'a': ElementState{
		IsRemoved: true,
		UpdatedAt: now,
	}})
	s.Merge(s2)
	ok = s.Lookup('a')
	require.False(t, ok)
}

func TestLWWElements(t *testing.T) {
	var tests = []struct {
		name     string
		a        *LWW
		expected []interface{}
	}{
		{
			name:     "empty",
			a:        &LWW{},
			expected: []interface{}{},
		},
		{
			name: "one elements, one valid",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			expected: []interface{}{'a'},
		},
		{
			name: "two elements, one valid",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
				'b': ElementState{IsRemoved: true, UpdatedAt: 1567586021000000000},
			}},
			expected: []interface{}{'a'},
		},
		{
			name: "one elements, none valid",
			a: &LWW{m: Elements{
				'a': ElementState{IsRemoved: true, UpdatedAt: 1567586021000000000},
			}},
			expected: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.expected, tt.a.Elements())
		})
	}
}

func TestLWWEqual(t *testing.T) {
	var tests = []struct {
		name     string
		a        *LWW
		b        *LWW
		expected bool
	}{
		{
			name:     "empty",
			a:        &LWW{},
			b:        &LWW{},
			expected: true,
		},
		{
			name: "one element equal",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			b: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			expected: true,
		},
		{
			name: "unequal length",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			b: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
				'b': ElementState{UpdatedAt: 1567586022000000000},
			}},
			expected: false,
		},
		{
			name: "one element unequal values",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			b: &LWW{m: Elements{
				'b': ElementState{UpdatedAt: 1567586021000000000},
			}},
			expected: false,
		},
		{
			name: "one element unequal time",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			b: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586022000000000},
			}},
			expected: false,
		},
		{
			name: "one element unequal state",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			b: &LWW{m: Elements{
				'a': ElementState{IsRemoved: true, UpdatedAt: 1567586021000000000},
			}},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.expected, tt.a.Equal(tt.b))
		})
	}
}

func TestLWWMerge(t *testing.T) {
	var tests = []struct {
		name     string
		a        *LWW
		b        *LWW
		expected Elements
	}{
		{
			name: "merge one",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
			}},
			b: &LWW{m: Elements{
				'b': ElementState{UpdatedAt: 1567586022000000000},
			}},
			expected: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
				'b': ElementState{UpdatedAt: 1567586022000000000},
			},
		},
		{
			name: "merge one with duplicate",
			a: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
				'b': ElementState{IsRemoved: true, UpdatedAt: 1567586021000000000},
			}},
			b: &LWW{m: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
				'b': ElementState{IsRemoved: true, UpdatedAt: 1567586022000000000},
			}},
			expected: Elements{
				'a': ElementState{UpdatedAt: 1567586021000000000},
				'b': ElementState{IsRemoved: true, UpdatedAt: 1567586022000000000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.Merge(tt.b)
			require.EqualValues(t, tt.expected, tt.a.m)
			require.EqualValues(t, tt.expected, tt.a.m)
		})
	}
}
