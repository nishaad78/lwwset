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
	now := time.Now()
	// add and remove 'a' at the same time
	m := LWWElements{'a': LWWTime{now, now}}
	s := NewFromElements(m)
	ok := s.Lookup('a')
	require.False(t, ok)
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
			a: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
			}},
			b: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
			}},
			expected: true,
		},
		{
			name: "unequal length",
			a: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
			}},
			b: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
				'b': LWWTime{Add: time.Unix(1567586022, 0)},
			}},
			expected: false,
		},
		{
			name: "one element unequal",
			a: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
			}},
			b: &LWW{m: LWWElements{
				'b': LWWTime{Add: time.Unix(1567586021, 0)},
			}},
			expected: false,
		},
		{
			name: "one element unequal time",
			a: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
			}},
			b: &LWW{m: LWWElements{
				'a': LWWTime{Remove: time.Unix(1567586021, 0)},
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
		expected LWWElements
	}{
		{
			name: "merge one",
			a: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
			}},
			b: &LWW{m: LWWElements{
				'b': LWWTime{Add: time.Unix(1567586022, 0)},
			}},
			expected: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
				'b': LWWTime{Add: time.Unix(1567586022, 0)},
			},
		},
		{
			name: "merge one with duplicate",
			a: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
				'b': LWWTime{Remove: time.Unix(1567586021, 0)},
			}},
			b: &LWW{m: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
				'b': LWWTime{Remove: time.Unix(1567586022, 0)},
			}},
			expected: LWWElements{
				'a': LWWTime{Add: time.Unix(1567586021, 0)},
				'b': LWWTime{Remove: time.Unix(1567586022, 0)},
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
