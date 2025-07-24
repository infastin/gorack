package ht_test

import (
	"testing"

	"github.com/infastin/gorack/ht"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		dur      string
		expected ht.Duration
	}{
		{"1h", ht.Hour},
		{"1m", ht.Minute},
		{"1s", ht.Second},
		{"1ms", ht.Millisecond},
		{"1µs", ht.Microsecond},
		{"1us", ht.Microsecond},
		{"1ns", ht.Nanosecond},
		{"4.000000001s", 4*ht.Second + ht.Nanosecond},
		{"1h0m4.000000001s", ht.Hour + 4*ht.Second + ht.Nanosecond},
		{"1h1m0.01s", 61*ht.Minute + 10*ht.Millisecond},
		{"1h1m0.123456789s", 61*ht.Minute + 123456789*ht.Nanosecond},
		{"1.00002ms", ht.Millisecond + 20*ht.Nanosecond},
		{"1.00000002s", ht.Second + 20*ht.Nanosecond},
		{"693ns", 693 * ht.Nanosecond},
		{"10s1us693ns", 10*ht.Second + ht.Microsecond + 693*ht.Nanosecond},
		{"1ms1ns", ht.Millisecond + 1*ht.Nanosecond},
		{"1s20ns", ht.Second + 20*ht.Nanosecond},
		{"60h8ms", 60*ht.Hour + 8*ht.Millisecond},
		{"96h63s", 4*ht.Day + 63*ht.Second},
		{"2d3s96ns", 2*ht.Day + 3*ht.Second + 96*ht.Nanosecond},
		{"4d3s96ns", 4*ht.Day + 3*ht.Second + 96*ht.Nanosecond},
		{"7d3s3µs96ns", 7*ht.Day + 3*ht.Second + 3*ht.Microsecond + 96*ht.Nanosecond},
	}

	for i, tt := range tests {
		parsed, err := ht.ParseDuration(tt.dur)
		if err != nil {
			t.Errorf("index %d: %s failed: %s", i, tt.dur, err.Error())
		} else if tt.expected != parsed {
			t.Errorf("index %d: %s returned: %d not equal to %d", i, tt.dur, parsed, tt.expected)
		}

		negDur := "-" + tt.dur
		negExpected := -tt.expected

		parsed, err = ht.ParseDuration(negDur)
		if err != nil {
			t.Errorf("index %d: %s failed: %s", i, negDur, err.Error())
		} else if negExpected != parsed {
			t.Errorf("index %d: %s returned: %d not equal to %d", i, negDur, parsed, negExpected)
		}
	}
}

func TestDuration_String(t *testing.T) {
	tests := []struct {
		dur      ht.Duration
		expected string
	}{
		{ht.Hour, "1h0m0s"},
		{ht.Minute, "1m0s"},
		{ht.Second, "1s"},
		{ht.Millisecond, "1ms"},
		{ht.Microsecond, "1µs"},
		{ht.Nanosecond, "1ns"},
		{4*ht.Second + ht.Nanosecond, "4.000000001s"},
		{ht.Hour + 4*ht.Second + ht.Nanosecond, "1h0m4.000000001s"},
		{61*ht.Minute + 10*ht.Millisecond, "1h1m0.01s"},
		{61*ht.Minute + 123456789*ht.Nanosecond, "1h1m0.123456789s"},
		{ht.Millisecond + 20*ht.Nanosecond, "1.00002ms"},
		{ht.Second + 20*ht.Nanosecond, "1.00000002s"},
		{693 * ht.Nanosecond, "693ns"},
		{10*ht.Second + ht.Microsecond + 693*ht.Nanosecond, "10.000001693s"},
		{ht.Millisecond + 1*ht.Nanosecond, "1.000001ms"},
		{ht.Second + 20*ht.Nanosecond, "1.00000002s"},
		{60*ht.Hour + 8*ht.Millisecond, "2d12h0m0.008s"},
		{96*ht.Hour + 63*ht.Second, "4d1m3s"},
		{2*ht.Day + 3*ht.Second + 96*ht.Nanosecond, "2d3.000000096s"},
		{4*ht.Day + 3*ht.Second + 96*ht.Nanosecond, "4d3.000000096s"},
		{7*ht.Day + 3*ht.Second + 3*ht.Microsecond + 96*ht.Nanosecond, "7d3.000003096s"},
	}

	for i, tt := range tests {
		got := tt.dur.String()
		if tt.expected != got {
			t.Errorf("index %d: %d returned: %s not equal to %s", i, tt.dur, got, tt.expected)
		}

		neg := -tt.dur
		negExpected := "-" + tt.expected

		negGot := neg.String()
		if negExpected != negGot {
			t.Errorf("index %d: %d returned: %s not equal to %s", i, neg, negGot, negExpected)
		}
	}
}
