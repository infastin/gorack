package ht

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type Clock struct {
	hour   uint8
	minute uint8
	second uint8
}

func NewClock(hour, minute, second int) Clock {
	minute, second = norm(minute, second, 60)
	hour, minute = norm(hour, minute, 24)
	return Clock{
		hour:   uint8(hour),
		minute: uint8(minute),
		second: uint8(second),
	}
}

func ParseClock(s string) (c Clock, err error) {
	orig := s
	if len(orig) != 5 && len(orig) != 8 {
		return Clock{}, fmt.Errorf("ht: invalid clock %q", orig)
	}

	hour, s, err := leadingInt(s)
	if err != nil || hour >= 24 {
		return Clock{}, fmt.Errorf("ht: invalid clock %q", orig)
	}

	if s[0] != ':' {
		return Clock{}, fmt.Errorf("ht: invalid clock %q", orig)
	}
	s = s[1:]

	minute, s, err := leadingInt(s)
	if err != nil || minute >= 60 {
		return Clock{}, fmt.Errorf("ht: invalid clock %q", orig)
	}

	var second uint64
	if s != "" {
		if s[0] != ':' {
			return Clock{}, fmt.Errorf("ht: invalid clock %q", orig)
		}
		s = s[1:]

		second, s, err = leadingInt(s)
		if err != nil || second >= 60 || s != "" {
			return Clock{}, fmt.Errorf("ht: invalid clock %q", orig)
		}
	}

	return Clock{
		hour:   uint8(hour),
		minute: uint8(minute),
		second: uint8(second),
	}, nil
}

func (c Clock) Hour() int {
	return int(c.hour)
}

func (c Clock) Minute() int {
	return int(c.minute)
}

func (c Clock) Second() int {
	return int(c.second)
}

func (c Clock) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%02d:%02d", c.hour, c.minute)
	if c.second != 0 {
		fmt.Fprintf(&b, ":%02d", c.second)
	}
	return b.String()
}

func (c Clock) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Clock) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		tmp, err := ParseClock(value)
		if err != nil {
			return err
		}
		*c = tmp
	default:
		return errors.New("ht: invalid clock")
	}
	return nil
}

func (c Clock) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Clock) UnmarshalText(b []byte) error {
	tmp, err := ParseClock(string(b))
	if err != nil {
		return err
	}
	*c = tmp
	return nil
}
