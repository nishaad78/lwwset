package lwwset

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLWWBasicAddRemove(t *testing.T) {
	s := NewLWW()
	v, ok := s.Lookup('a')
	require.False(t, ok)
	require.Equal(t, time.Time{}, v)

	// add 'a'
	now := time.Now()
	s.Add('a', now)
	v, ok = s.Lookup('a')
	require.True(t, ok)
	require.Equal(t, now, v)

	// remove 'a'
	now = time.Now()
	s.Remove('a', now)
	v, ok = s.Lookup('a')
	require.False(t, ok)
	require.Equal(t, time.Time{}, v)

	// add 'a' again
	now = time.Now()
	s.Add('a', now)
	v, ok = s.Lookup('a')
	require.True(t, ok)
	require.Equal(t, now, v)
}

func TestLWWRemoveBias(t *testing.T) {
	s := NewLWW()

	// add and remove 'a'
	now := time.Now()
	s.Add('a', now)
	s.Remove('a', now)
	v, ok := s.Lookup('a')
	require.False(t, ok)
	require.Equal(t, time.Time{}, v)
}

func TestLWWEqual(t *testing.T) {
	var tests = []struct {
		name     string
		a        *LWW
		b        *LWW
		expected bool
	}{
		{
			name: "empty",
			a: &LWW{
				a: &g{m: map[interface{}]time.Time{}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			b: &LWW{
				a: &g{m: map[interface{}]time.Time{}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			expected: true,
		},
		{
			name: "one element equal",
			a: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			b: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			expected: true,
		},
		{
			name: "one element unequal",
			a: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			b: &LWW{
				a: &g{m: map[interface{}]time.Time{}},
				r: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
			},
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
		expected *LWW
	}{
		{
			name: "merge one",
			a: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			b: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'b': time.Unix(1567586022, 0),
				}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			expected: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
					'b': time.Unix(1567586022, 0),
				}},
				r: &g{m: map[interface{}]time.Time{}},
			},
		},
		{
			name: "merge one with duplicate",
			a: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
				r: &g{m: map[interface{}]time.Time{}},
			},
			b: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
				r: &g{m: map[interface{}]time.Time{
					'b': time.Unix(1567586022, 0),
				}},
			},
			expected: &LWW{
				a: &g{m: map[interface{}]time.Time{
					'a': time.Unix(1567586021, 0),
				}},
				r: &g{m: map[interface{}]time.Time{
					'b': time.Unix(1567586022, 0),
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.Merge(tt.b)
			require.EqualValues(t, tt.expected.a, tt.a.a)
			require.EqualValues(t, tt.expected.r, tt.a.r)
		})
	}
}
