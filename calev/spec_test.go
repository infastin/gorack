package calev_test

import (
	"slices"
	"testing"
	"time"

	"github.com/infastin/gorack/calev"
)

func TestClock(t *testing.T) {
	date := time.Date(2025, time.January, 1, 22, 59, 0, 0, time.UTC)

	tests := []struct {
		opts     []calev.SpecOpt
		expected []time.Time
		wrapHour bool
	}{
		{
			opts: []calev.SpecOpt{calev.Minute(0)},
			expected: []time.Time{
				time.Date(2025, time.January, 1, 23, 0, 0, 0, time.UTC),
			},
			wrapHour: true,
		},
		{
			opts: []calev.SpecOpt{calev.Minute(3, 5, 7)},
			expected: []time.Time{
				time.Date(2025, time.January, 1, 23, 3, 0, 0, time.UTC),
				time.Date(2025, time.January, 1, 23, 5, 0, 0, time.UTC),
				time.Date(2025, time.January, 1, 23, 7, 0, 0, time.UTC),
			},
			wrapHour: true,
		},
		{
			opts: []calev.SpecOpt{calev.EveryMinute(3, 5, 1)},
			expected: []time.Time{
				time.Date(2025, time.January, 1, 23, 3, 0, 0, time.UTC),
				time.Date(2025, time.January, 1, 23, 4, 0, 0, time.UTC),
				time.Date(2025, time.January, 1, 23, 5, 0, 0, time.UTC),
			},
			wrapHour: true,
		},
		{
			opts: []calev.SpecOpt{calev.Hour(0), calev.Minute(0)},
			expected: []time.Time{
				time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.Hour(0), calev.Minute(0), calev.Second(22)},
			expected: []time.Time{
				time.Date(2025, time.January, 2, 0, 0, 22, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.EveryHour(10, 14, 1), calev.Minute(0)},
			expected: []time.Time{
				time.Date(2025, time.January, 2, 10, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 11, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 12, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 13, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 14, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.Hour(0), calev.EveryMinute(5, 7, 1)},
			expected: []time.Time{
				time.Date(2025, time.January, 2, 0, 5, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 0, 6, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 0, 7, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{
				calev.Hour(7, 11, 13, 17),
				calev.Minute(5, 23, 37, 41),
			},
			expected: []time.Time{
				time.Date(2025, time.January, 2, 7, 5, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 7, 23, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 7, 37, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 7, 41, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 11, 5, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 11, 23, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 11, 37, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 11, 41, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 13, 5, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 13, 23, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 13, 37, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 13, 41, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 17, 5, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 17, 23, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 17, 37, 0, 0, time.UTC),
				time.Date(2025, time.January, 2, 17, 41, 0, 0, time.UTC),
			},
		},
	}

	for i := range tests {
		tt := &tests[i]
		expected := slices.Clone(tt.expected)
		for range 2 {
			if tt.wrapHour {
				for j := range expected {
					expected[j] = expected[j].Add(time.Hour)
				}
			} else {
				for j := range expected {
					expected[j] = expected[j].AddDate(0, 0, 1)
				}
			}
			tt.expected = append(tt.expected, expected...)
		}
	}

	for i, tt := range tests {
		cur := date
		spec := calev.New(tt.opts...)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateTime), cur.Format(time.DateTime))
				return
			}
		}
	}
}

func TestDay(t *testing.T) {
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	defaultOpts := []calev.SpecOpt{
		calev.Hour(0), calev.Minute(0),
	}

	tests := []struct {
		opts     []calev.SpecOpt
		expected []time.Time
	}{
		{
			opts: []calev.SpecOpt{calev.Day(1)},
			expected: []time.Time{
				time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.Day(1, 3, 5, 7)},
			expected: []time.Time{
				time.Date(2025, time.January, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.EveryDay(1, 14, 3)},
			expected: []time.Time{
				time.Date(2025, time.January, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for i := range tests {
		tt := &tests[i]
		expected := slices.Clone(tt.expected)
		for range 2 {
			for j := range expected {
				expected[j] = expected[j].AddDate(0, 1, 0)
			}
			tt.expected = append(tt.expected, expected...)
		}
	}

	for i, tt := range tests {
		cur := date
		spec := calev.New(slices.Concat(defaultOpts, tt.opts)...)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateOnly), cur.Format(time.DateOnly))
				return
			}
		}
	}
}

func TestLastDay(t *testing.T) {
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	defaultOpts := []calev.SpecOpt{
		calev.Hour(0), calev.Minute(0),
	}

	tests := []struct {
		opts     []calev.SpecOpt
		expected []time.Time
	}{
		{
			opts: []calev.SpecOpt{calev.LastDay(1)},
			expected: []time.Time{
				time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.March, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.April, 30, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.LastDay(1, 3, 5, 7)},
			expected: []time.Time{
				time.Date(2025, time.January, 25, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 22, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 26, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.EveryLastDay(1, 5, 2)},
			expected: []time.Time{
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 26, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for i, tt := range tests {
		cur := date
		spec := calev.New(slices.Concat(defaultOpts, tt.opts)...)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateOnly), cur.Format(time.DateOnly))
				return
			}
		}
	}
}

func TestWeekday(t *testing.T) {
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	defaultOpts := []calev.SpecOpt{
		calev.Hour(0), calev.Minute(0),
	}

	tests := []struct {
		opts     []calev.SpecOpt
		expected []time.Time
	}{
		{
			opts: []calev.SpecOpt{calev.Weekday(time.Monday)},
			expected: []time.Time{
				time.Date(2025, time.January, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.Weekday(time.Monday, time.Tuesday, time.Sunday)},
			expected: []time.Time{
				time.Date(2025, time.January, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 12, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 14, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 19, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 21, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 26, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 28, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{calev.EveryWeekday(time.Monday, time.Friday, 1)},
			expected: []time.Time{
				time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 9, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 10, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{
				calev.Day(10, 20),
				calev.Weekday(time.Monday),
			},
			expected: []time.Time{
				time.Date(2025, time.January, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for i, tt := range tests {
		cur := date
		spec := calev.New(slices.Concat(defaultOpts, tt.opts)...)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateOnly), cur.Format(time.DateOnly))
				return
			}
		}
	}
}

func TestWeekdayStrict(t *testing.T) {
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	defaultOpts := []calev.SpecOpt{
		calev.Hour(0), calev.Minute(0),
	}

	tests := []struct {
		opts     []calev.SpecOpt
		expected []time.Time
	}{
		{
			opts: []calev.SpecOpt{calev.Weekday(time.Monday)},
			expected: []time.Time{
				time.Date(2025, time.January, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{
				calev.EveryDay(1, 7, 1),
				calev.WeekdayStrict(time.Monday, time.Friday),
			},
			expected: []time.Time{
				time.Date(2025, time.January, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 7, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{
				calev.EveryDay(20, 31, 1),
				calev.WeekdayStrict(time.Monday, time.Friday),
			},
			expected: []time.Time{
				time.Date(2025, time.January, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 21, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{
				calev.EveryLastDay(1, 7, 1),
				calev.WeekdayStrict(time.Monday, time.Friday),
			},
			expected: []time.Time{
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			opts: []calev.SpecOpt{
				calev.EveryLastDay(1, 20, 1),
				calev.WeekdayStrict(time.Monday, time.Friday),
			},
			expected: []time.Time{
				time.Date(2025, time.January, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.January, 31, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 14, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 21, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for i, tt := range tests {
		cur := date
		spec := calev.New(slices.Concat(defaultOpts, tt.opts)...)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateOnly), cur.Format(time.DateOnly))
				return
			}
		}
	}
}

func TestMonth(t *testing.T) {
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	defaultOpts := []calev.SpecOpt{
		calev.Day(1),
		calev.Hour(0), calev.Minute(0),
	}

	tests := []struct {
		opts     []calev.SpecOpt
		expected []time.Time
		wrapLast bool
	}{
		{
			opts:     []calev.SpecOpt{calev.Month(time.January)},
			expected: []time.Time{time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
			wrapLast: false,
		},
		{
			opts: []calev.SpecOpt{calev.Month(time.January, time.August, time.March)},
			expected: []time.Time{
				time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			wrapLast: true,
		},
		{
			opts: []calev.SpecOpt{calev.EveryMonth(time.January, time.March, 1)},
			expected: []time.Time{
				time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			wrapLast: true,
		},
		{
			opts: []calev.SpecOpt{calev.EveryMonth(time.January, time.September, 2)},
			expected: []time.Time{
				time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			wrapLast: true,
		},
		{
			opts: []calev.SpecOpt{
				calev.Month(time.February, time.November),
				calev.EveryMonth(time.January, time.September, 2),
			},
			expected: []time.Time{
				time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.September, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, time.November, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			wrapLast: true,
		},
	}

	for i := range tests {
		tt := &tests[i]
		expected := slices.Clone(tt.expected)
		for range 2 {
			if tt.wrapLast {
				year := expected[len(expected)-1].Year()
				for j := range len(expected) - 1 {
					expected[j] = expected[j].AddDate(year-expected[j].Year(), 0, 0)
				}
				expected[len(expected)-1] = expected[len(expected)-1].AddDate(1, 0, 0)
			} else {
				for j := range expected {
					expected[j] = expected[j].AddDate(1, 0, 0)
				}
			}
			tt.expected = append(tt.expected, expected...)
		}
	}

	for i, tt := range tests {
		cur := date
		spec := calev.New(slices.Concat(defaultOpts, tt.opts)...)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateOnly), cur.Format(time.DateOnly))
				return
			}
		}
	}
}

func TestDate(t *testing.T) {
	date := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

	defaultOpts := []calev.SpecOpt{
		calev.Hour(0), calev.Minute(0),
	}

	tests := []struct {
		opts     []calev.SpecOpt
		expected []time.Time
	}{
		// Days of month
		{
			opts:     []calev.SpecOpt{calev.Month(time.February), calev.Day(30)},
			expected: []time.Time{{}},
		},
		{
			opts:     []calev.SpecOpt{calev.Month(time.April), calev.Day(31)},
			expected: []time.Time{{}},
		},
		{
			opts:     []calev.SpecOpt{calev.Month(time.June), calev.Day(31)},
			expected: []time.Time{{}},
		},
		{
			opts:     []calev.SpecOpt{calev.Month(time.June), calev.Day(31)},
			expected: []time.Time{{}},
		},
		// Last days of month
		{
			opts:     []calev.SpecOpt{calev.Month(time.February), calev.LastDay(30)},
			expected: []time.Time{{}},
		},
		{
			opts:     []calev.SpecOpt{calev.Month(time.April), calev.LastDay(31)},
			expected: []time.Time{{}},
		},
		{
			opts:     []calev.SpecOpt{calev.Month(time.June), calev.LastDay(31)},
			expected: []time.Time{{}},
		},
		{
			opts:     []calev.SpecOpt{calev.Month(time.June), calev.LastDay(31)},
			expected: []time.Time{{}},
		},
	}

	for i, tt := range tests {
		spec := calev.New(slices.Concat(defaultOpts, tt.opts)...)
		cur := spec.Next(date)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateOnly), cur.Format(time.DateOnly))
				return
			}
		}
	}
}

func TestLoopYear(t *testing.T) {
	defaultOpts := []calev.SpecOpt{
		calev.Hour(0), calev.Minute(0),
	}

	tests := []struct {
		date     time.Time
		opts     []calev.SpecOpt
		expected []time.Time
	}{
		{
			date: time.Date(2028, time.January, 1, 0, 0, 0, 0, time.UTC),
			opts: []calev.SpecOpt{calev.Month(time.February), calev.Day(29)},
			expected: []time.Time{
				time.Date(2028, time.February, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2032, time.February, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2036, time.February, 29, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			date: time.Date(1996, time.January, 1, 0, 0, 0, 0, time.UTC),
			opts: []calev.SpecOpt{calev.Month(time.February), calev.Day(29)},
			expected: []time.Time{
				time.Date(1996, time.February, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2004, time.February, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2008, time.February, 29, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			date: time.Date(2096, time.January, 1, 0, 0, 0, 0, time.UTC),
			opts: []calev.SpecOpt{calev.Month(time.February), calev.Day(29)},
			expected: []time.Time{
				time.Date(2096, time.February, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2104, time.February, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2108, time.February, 29, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			date: time.Date(2028, time.January, 1, 0, 0, 0, 0, time.UTC),
			opts: []calev.SpecOpt{calev.Month(time.February), calev.LastDay(29)},
			expected: []time.Time{
				time.Date(2028, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2032, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2036, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			date: time.Date(1996, time.January, 1, 0, 0, 0, 0, time.UTC),
			opts: []calev.SpecOpt{calev.Month(time.February), calev.LastDay(29)},
			expected: []time.Time{
				time.Date(1996, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2004, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2008, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			date: time.Date(2096, time.January, 1, 0, 0, 0, 0, time.UTC),
			opts: []calev.SpecOpt{calev.Month(time.February), calev.LastDay(29)},
			expected: []time.Time{
				time.Date(2096, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2104, time.February, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2108, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			date: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			opts: []calev.SpecOpt{
				calev.Month(time.February),
				calev.Day(29), calev.WeekdayStrict(time.Monday),
			},
			expected: []time.Time{
				time.Date(2044, time.February, 29, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for i, tt := range tests {
		cur := tt.date
		spec := calev.New(slices.Concat(defaultOpts, tt.opts)...)
		for j, expected := range tt.expected {
			cur = spec.Next(cur)
			if !cur.Equal(expected) {
				t.Errorf("tests[%d][%d]: expected to get %s, but got %s",
					i, j, expected.Format(time.DateOnly), cur.Format(time.DateOnly))
				return
			}
		}
	}
}

func TestDST(t *testing.T) {
	locNY, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
		return
	}

	locLH, err := time.LoadLocation("Australia/Lord_Howe")
	if err != nil {
		t.Fatal(err)
		return
	}

	tests := []struct {
		date     time.Time
		opts     []calev.SpecOpt
		expected time.Time
	}{
		// America/New_York
		{
			date: time.Date(2025, time.March, 9, 0, 0, 0, 0, locNY),
			opts: []calev.SpecOpt{
				calev.Month(time.March),
				calev.Hour(2), calev.Minute(0),
			},
			expected: time.Date(2025, time.March, 10, 2, 0, 0, 0, locNY),
		},
		{
			date: time.Date(2025, time.March, 9, 0, 0, 0, 0, locNY),
			opts: []calev.SpecOpt{
				calev.Month(time.March),
				calev.Hour(3), calev.Minute(0),
			},
			expected: time.Date(2025, time.March, 9, 3, 0, 0, 0, locNY),
		},
		{
			date: time.Date(2025, time.January, 1, 0, 0, 0, 0, locNY),
			opts: []calev.SpecOpt{
				calev.Month(time.March),
				calev.Day(9),
				calev.Hour(2), calev.Minute(0),
			},
			expected: time.Date(2026, time.March, 9, 2, 0, 0, 0, locNY),
		},
		// Australia/Lord_Howe
		{
			date: time.Date(2025, time.October, 5, 1, 0, 0, 0, locLH),
			opts: []calev.SpecOpt{
				calev.Month(time.October),
				calev.Hour(2), calev.Minute(0),
			},
			expected: time.Date(2025, time.October, 6, 2, 0, 0, 0, locLH),
		},
		{
			date: time.Date(2025, time.October, 5, 1, 0, 0, 0, locLH),
			opts: []calev.SpecOpt{
				calev.Month(time.October),
				calev.Hour(2), calev.Minute(30),
			},
			expected: time.Date(2025, time.October, 5, 2, 30, 0, 0, locLH),
		},
		{
			date: time.Date(2025, time.October, 5, 1, 0, 0, 0, locLH),
			opts: []calev.SpecOpt{
				calev.Month(time.October),
				calev.Day(5),
				calev.Hour(2), calev.Minute(0),
			},
			expected: time.Date(2026, time.October, 5, 2, 0, 0, 0, locLH),
		},
	}

	for i, tt := range tests {
		spec := calev.New(tt.opts...)
		next := spec.Next(tt.date)
		if !next.Equal(tt.expected) {
			t.Errorf("tests[%d]: expected to get %s, but got %s",
				i, tt.expected.Format(time.DateOnly), next.Format(time.DateOnly))
			return
		}
	}
}
