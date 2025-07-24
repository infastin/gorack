package ht

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

// Duration is time.Duration with additional day unit (`d`) during parsing
// and with JSON and Text (un)marshaling support.
type Duration time.Duration

const (
	Nanosecond  = Duration(time.Nanosecond)
	Microsecond = Duration(time.Microsecond)
	Millisecond = Duration(time.Millisecond)
	Second      = Duration(time.Second)
	Minute      = Duration(time.Minute)
	Hour        = Duration(time.Hour)
	Day         = Duration(24 * time.Hour)
)

var durationUnitMap = map[string]uint64{
	"ns": uint64(Nanosecond),
	"us": uint64(Microsecond),
	"µs": uint64(Microsecond), // U+00B5 = micro symbol
	"μs": uint64(Microsecond), // U+03BC = Greek letter mu
	"ms": uint64(Millisecond),
	"s":  uint64(Second),
	"m":  uint64(Minute),
	"h":  uint64(Hour),
	"d":  uint64(Day),
}

// ParseDuration parses a duration string.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d".
func ParseDuration(s string) (Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d uint64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, fmt.Errorf("ht: invalid duration %q", orig)
	}
	for s != "" {
		var (
			v, f  uint64      // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if s[0] != '.' && (s[0] < '0' || s[0] > '9') {
			return 0, fmt.Errorf("ht: invalid duration %q", orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, fmt.Errorf("ht: invalid duration %q", orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, fmt.Errorf("ht: invalid duration %q", orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, fmt.Errorf("ht: missing unit in duration %q", orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := durationUnitMap[u]
		if !ok {
			return 0, fmt.Errorf("ht: unknown unit %q in duration %q", u, orig)
		}
		if v > 1<<63/unit {
			// overflow
			return 0, fmt.Errorf("ht: invalid duration %q", orig)
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += uint64(float64(f) * (float64(unit) / scale))
			if v > 1<<63 {
				// overflow
				return 0, fmt.Errorf("ht: invalid duration %q", orig)
			}
		}
		d += v
		if d > 1<<63 {
			return 0, fmt.Errorf("ht: invalid duration %q", orig)
		}
	}
	if neg {
		return -Duration(d), nil
	}
	if d > 1<<63-1 {
		return 0, fmt.Errorf("ht: invalid duration %q", orig)
	}
	return Duration(d), nil
}

func (d Duration) String() string {
	var neg bool
	if d < 0 {
		neg = true
		d = -d
	}
	var result []byte
	if neg {
		result = append(result, '-')
	}
	days, rem := d/Day, d%Day
	if days > 0 {
		result = strconv.AppendUint(result, uint64(days), 10)
		result = append(result, 'd')
	}
	if rem > 0 {
		result = append(result, time.Duration(rem).String()...)
	}
	return unsafe.String(unsafe.SliceData(result), len(result))
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(value)
		return nil
	case string:
		tmp, err := ParseDuration(value)
		if err != nil {
			return err
		}
		*d = tmp
		return nil
	default:
		return errors.New("ht: invalid duration")
	}
}

func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Duration) UnmarshalText(b []byte) error {
	dur, err := ParseDuration(string(b))
	if err != nil {
		return err
	}
	*d = dur
	return nil
}
