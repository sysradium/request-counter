package ephemeral_test

import (
	"testing"
	"time"

	"github.com/sysradium/request-counter/internal/counter/ephemeral"
)

func TestIsEmptyOnStart(t *testing.T) {
	s := ephemeral.New(time.Second)

	if l := s.Len(); l != 0 {
		t.Errorf("expected storage to empty, got %d items instead", l)
	}
}

func TestLen(t *testing.T) {
	now := time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	window := 5 * time.Second

	var tests = []struct {
		name     string
		given    []time.Time
		expected int
	}{
		{
			name:     "initially empty",
			expected: 0,
		},
		{
			name:     "1 fresh request",
			given:    []time.Time{now.Add(-time.Second)},
			expected: 1,
		}, {
			name: "a few requests within window",
			given: []time.Time{
				now.Add(-2 * time.Second),
				now.Add(-3 * time.Second),
				now.Add(-time.Second),
			},
			expected: 3,
		}, {
			name: "all stale",
			given: []time.Time{
				now.Add(-window).Add(-2 * time.Second),
				now.Add(-window).Add(-time.Second),
			},
			expected: 0,
		}, {
			name: "mix of out of and in a window requests",
			given: []time.Time{
				now.Add(-window).Add(-2 * time.Second),
				now.Add(-window).Add(-time.Second),
				now.Add(-2 * time.Second),
				now.Add(-time.Second),
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := ephemeral.New(window,
				ephemeral.WithClock(func() time.Time {
					return now
				}),
			)

			for _, ts := range tt.given {
				if err := s.Add(ts); err != nil {
					t.Fatalf("add returned unexpected err: %v", err)
				}
			}

			if actual := s.Len(); actual != tt.expected {
				t.Errorf("(%s): expected %v, actual %v", tt.given, tt.expected, actual)
			}

		})
	}
}

func TestPruneKeepsRelevantData(t *testing.T) {
	now := time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	window := 5 * time.Second

	s := ephemeral.New(window,
		ephemeral.WithClock(func() time.Time {
			return now
		}),
	)

	windowBoundary := now.Add(-window)
	for _, ts := range []time.Time{
		windowBoundary.Add(-time.Second),
		now.Add(-time.Second),
		now.Add(-time.Second),
		now.Add(-time.Second),
	} {
		if err := s.Add(ts); err != nil {
			t.Fatalf("add returned unexpected err: %v", err)
		}
	}

	s.Prune()

	expected := 3
	if l := s.Len(); l != expected {
		t.Errorf("expected len to be %d got %d instead", expected, l)
	}
}

func TestGet(t *testing.T) {
	now := time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	window := 5 * time.Second

	var tests = []struct {
		name        string
		expectedLen int
		given       []time.Time
	}{
		{
			name:        "nothing",
			expectedLen: 0,
			given:       []time.Time{},
		}, {
			name:        "3 fresh items",
			expectedLen: 3,
			given: []time.Time{
				now.Add(-time.Second),
				now.Add(-time.Second),
				now.Add(-time.Second),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := ephemeral.New(window,
				ephemeral.WithClock(func() time.Time {
					return now
				}),
			)

			for _, ts := range tt.given {
				if err := s.Add(ts); err != nil {
					t.Fatalf("add returned unexpected err: %v", err)
				}
			}

			items := s.Get()
			actual := len(items)

			if actual != tt.expectedLen {
				t.Errorf("(%s): expected %d items, got %d", tt.name, tt.expectedLen, actual)
			}

		})
	}
}
