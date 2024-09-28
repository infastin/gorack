package xtypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type TimeOfDay struct {
	hour   uint8
	minute uint8
}

func NewTimeOfDay(hour, minute int) TimeOfDay {
	return TimeOfDay{
		hour:   uint8(hour) % 24,
		minute: uint8(minute) % 60,
	}
}

func ParseTimeOfDay(s string) (t TimeOfDay, err error) {
	if len(s) != 5 {
		return TimeOfDay{}, errors.New("time of day must be 5 characters long")
	}

	if s[2] != ':' {
		return TimeOfDay{}, errors.New("invalid time of day format, must be 15:04")
	}

	hourStr, minuteStr := s[:2], s[3:]

	hour, err := strconv.ParseUint(hourStr, 10, 8)
	if err != nil {
		return TimeOfDay{}, errors.New("invalid hour")
	}

	minute, err := strconv.ParseUint(minuteStr, 10, 8)
	if err != nil {
		return TimeOfDay{}, err
	}

	return TimeOfDay{
		hour:   uint8(hour),
		minute: uint8(minute),
	}, nil
}

func (t TimeOfDay) Hour() int {
	return int(t.hour)
}

func (t TimeOfDay) Minute() int {
	return int(t.minute)
}

func (t TimeOfDay) String() string {
	return fmt.Sprintf("%02d:%02d", t.hour, t.minute)
}

func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TimeOfDay) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		tmp, err := ParseTimeOfDay(value)
		if err != nil {
			return err
		}
		*t = tmp
	default:
		return errors.New("invalid time of day")
	}
	return nil
}

func (t TimeOfDay) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *TimeOfDay) UnmarshalText(b []byte) error {
	tmp, err := ParseTimeOfDay(string(b))
	if err != nil {
		return err
	}
	*t = tmp
	return nil
}
