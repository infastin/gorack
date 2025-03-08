package xtypes_test

import (
	"testing"
	"time"

	"github.com/infastin/gorack/xtypes"
)

func TestParseDuration(t *testing.T) {
	for i, tt := range []struct {
		dur      string
		expected time.Duration
	}{
		{"1h", time.Duration(time.Hour)},
		{"1m", time.Duration(time.Minute)},
		{"1s", time.Duration(time.Second)},
		{"1ms", time.Duration(time.Millisecond)},
		{"1µs", time.Duration(time.Microsecond)},
		{"1us", time.Duration(time.Microsecond)},
		{"1ns", time.Duration(time.Nanosecond)},
		{"4.000000001s", time.Duration(4*time.Second + time.Nanosecond)},
		{"1h0m4.000000001s", time.Duration(time.Hour + 4*time.Second + time.Nanosecond)},
		{"1h1m0.01s", time.Duration(61*time.Minute + 10*time.Millisecond)},
		{"1h1m0.123456789s", time.Duration(61*time.Minute + 123456789*time.Nanosecond)},
		{"1.00002ms", time.Duration(time.Millisecond + 20*time.Nanosecond)},
		{"1.00000002s", time.Duration(time.Second + 20*time.Nanosecond)},
		{"693ns", time.Duration(693 * time.Nanosecond)},
		{"10s1us693ns", time.Duration(10*time.Second + time.Microsecond + 693*time.Nanosecond)},
		{"1ms1ns", time.Duration(time.Millisecond + 1*time.Nanosecond)},
		{"1s20ns", time.Duration(time.Second + 20*time.Nanosecond)},
		{"60h8ms", time.Duration(60*time.Hour + 8*time.Millisecond)},
		{"96h63s", time.Duration(96*time.Hour + 63*time.Second)},
		{"2d3s96ns", time.Duration(48*time.Hour + 3*time.Second + 96*time.Nanosecond)},
		{"4d3s96ns", time.Duration(96*time.Hour + 3*time.Second + 96*time.Nanosecond)},
		{"7d3s3µs96ns", time.Duration(168*time.Hour + 3*time.Second + 3*time.Microsecond + 96*time.Nanosecond)},
	} {
		parsed, err := xtypes.ParseDuration(tt.dur)
		if err != nil {
			t.Logf("index %d -> in: %s returned: %s\tnot equal to %s", i, tt.dur, err.Error(), tt.expected.String())
		} else if tt.expected != time.Duration(parsed) {
			t.Errorf("index %d -> in: %s returned: %d\tnot equal to %d", i, tt.dur, parsed, tt.expected)
		}

		negDur := "-" + tt.dur
		negExpected := -tt.expected

		parsed, err = xtypes.ParseDuration(negDur)
		if err != nil {
			t.Logf("index %d -> in: %s returned: %s\tnot equal to %s", i, negDur, err.Error(), negExpected.String())
		} else if negExpected != time.Duration(parsed) {
			t.Errorf("index %d -> in: %s returned: %d\tnot equal to %d", i, negDur, parsed, negExpected)
		}
	}
}
