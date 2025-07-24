package ht_test

import (
	"testing"

	"github.com/infastin/gorack/ht"
)

func TestParseClock(t *testing.T) {
	tests := []struct {
		clock     string
		expected  ht.Clock
		expectErr bool
	}{
		{"15:04", ht.NewClock(15, 4, 0), false},
		{"15:04:05", ht.NewClock(15, 4, 5), false},
		{"15", ht.Clock{}, true},
		{"15:", ht.Clock{}, true},
		{"15:04:", ht.Clock{}, true},
		{"33:04", ht.Clock{}, true},
		{"23:61", ht.Clock{}, true},
		{"-10:59", ht.Clock{}, true},
		{"10-59", ht.Clock{}, true},
		{"10:10:60", ht.Clock{}, true},
	}

	for i, tt := range tests {
		parsed, err := ht.ParseClock(tt.clock)
		if (err != nil) != tt.expectErr {
			t.Errorf("index %d: %s failed: error = %v, expectErr = %t", i, tt.clock, err, tt.expectErr)
		} else if tt.expected != parsed {
			t.Errorf("index %d: %s returned: %s not equal to %s", i, tt.clock, parsed, tt.expected)
		}
	}
}

func TestClock_String(t *testing.T) {
	tests := []struct {
		clock    ht.Clock
		expected string
	}{
		{ht.NewClock(15, 4, 0), "15:04"},
		{ht.NewClock(15, 4, 5), "15:04:05"},
	}

	for i, tt := range tests {
		got := tt.clock.String()
		if tt.expected != got {
			t.Errorf("index %d: %s returned: %s not equal to %s", i, tt.clock, got, tt.expected)
		}
	}
}
