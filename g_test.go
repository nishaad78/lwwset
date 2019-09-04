package lwwset

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGBasicAddLookup(t *testing.T) {
	s := newG()
	v, ok := s.lookup("empty")
	require.False(t, ok)
	require.Equal(t, time.Time{}, v)

	now := time.Now()
	s.add('a', now)

	v, ok = s.lookup('a')
	require.True(t, ok)
	require.Equal(t, now, v)
}

func TestGEqual(t *testing.T) {
	var tests = []struct {
		name     string
		a        *g
		b        *g
		expected bool
	}{
		{
			name:     "empty",
			a:        &g{m: map[interface{}]time.Time{}},
			b:        &g{m: map[interface{}]time.Time{}},
			expected: true,
		},
		{
			name: "one element equal",
			a: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			}},
			b: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			}},
			expected: true,
		},
		{
			name: "one element unequal",
			a: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			}},
			b: &g{m: map[interface{}]time.Time{
				'b': time.Unix(1567586021, 0),
			}},
			expected: false,
		},
		{
			name: "one element different time",
			a: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			}},
			b: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586022, 0),
			}},
			expected: false,
		},
		{
			name: "one empty",
			a: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			}},
			b:        &g{m: map[interface{}]time.Time{}},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.expected, tt.a.equal(tt.b))
		})
	}
}

func TestGMerge(t *testing.T) {
	var tests = []struct {
		name     string
		a        *g
		b        *g
		expected map[interface{}]time.Time
	}{
		{
			name: "merge one",
			a: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			}},
			b: &g{m: map[interface{}]time.Time{
				'b': time.Unix(1567586022, 0),
			}},
			expected: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
				'b': time.Unix(1567586022, 0),
			},
		},
		{
			name: "merge one with duplicate",
			a: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
				'b': time.Unix(1567586021, 0),
			}},
			b: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
				'b': time.Unix(1567586022, 0),
			}},
			expected: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
				'b': time.Unix(1567586022, 0),
			},
		},
		{
			name: "merge one with a empty",
			a:    &g{m: map[interface{}]time.Time{}},
			b: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
				'b': time.Unix(1567586022, 0),
			}},
			expected: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
				'b': time.Unix(1567586022, 0),
			},
		},
		{
			name: "merge one with b empty",
			a: &g{m: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			}},
			b: &g{m: map[interface{}]time.Time{}},
			expected: map[interface{}]time.Time{
				'a': time.Unix(1567586021, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.merge(tt.b)
			require.EqualValues(t, tt.expected, tt.a.m)
		})
	}
}
