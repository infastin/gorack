package ht_test

import (
	"testing"

	"github.com/infastin/gorack/ht"
)

func TestParseByteSize(t *testing.T) {
	tests := []struct {
		bytesize  string
		units     []string
		expected  uint64
		expectErr bool
	}{
		{"0", []string{"", "B"}, 0, false},
		{"128", []string{"", "B"}, 128, false},
		{"32", []string{"k", "K", "KiB"}, uint64(32 * ht.KiB), false},
		{"100", []string{"KB"}, uint64(100 * ht.KB), false},
		{"2077", []string{"M", "MiB"}, uint64(2077 * ht.MiB), false},
		{"2000", []string{"MB"}, uint64(2000 * ht.MB), false},
		{"1.5", []string{"GB"}, uint64(1500 * ht.MB), false},
		{"4.5", []string{"TB"}, uint64(4500 * ht.GB), false},

		{"", []string{}, 0, true},
		{"", []string{"B", "k", "K"}, 0, true},
		{"128", []string{"s", "d"}, 0, true},
		{"abcde", []string{}, 0, true},
		{"100", []string{"  "}, 0, true},
		{"1.5", []string{" "}, 0, true},
		{"3.14159", []string{" "}, 0, true},
		{"2077", []string{"  k", "  k", "  KiB"}, 0, true},
		{"f", []string{}, 0, true},
	}

	for i, tt := range tests {
		for _, unit := range tt.units {
			for _, sep := range []string{"", " "} {
				bytesize := tt.bytesize
				if unit != "" {
					bytesize += sep + unit
				}

				si, err := ht.ParseByteSizeSI(bytesize)
				if (err != nil) != tt.expectErr {
					t.Errorf("index %d: %s failed: error = %v, expectErr = %t", i, bytesize, err, tt.expectErr)
				} else if tt.expected != uint64(si) {
					t.Errorf("index %d: %s returned: %d not equal to %d", i, bytesize, si, tt.expected)
				}

				iec, err := ht.ParseByteSizeSI(bytesize)
				if (err != nil) != tt.expectErr {
					t.Errorf("index %d: %s failed: error = %v, expectErr = %t", i, bytesize, err, tt.expectErr)
				} else if tt.expected != uint64(iec) {
					t.Errorf("index %d -> in: %s returned: %d not equal to %d", i, bytesize, iec, tt.expected)
				}

				if si != iec {
					t.Errorf("index %d: %s returned: SI(%d) not equal to IEC(%d)", i, bytesize, si, iec)
				}
			}
		}
	}
}

func TestByteSizeSI_String(t *testing.T) {
	tests := []struct {
		bytesize ht.ByteSizeSI
		expected string
	}{
		{0, "0B"},
		{128, "128B"},
		{32 * ht.KB, "32KB"},
		{2077 * ht.MB, "2.077GB"},
		{2000 * ht.MB, "2GB"},
	}

	for i, tt := range tests {
		got := tt.bytesize.String()
		if tt.expected != got {
			t.Errorf("index %d: %d returned: %s not equal to %s", i, tt.bytesize, got, tt.expected)
		}
	}
}

func TestByteSizeIEC_String(t *testing.T) {
	tests := []struct {
		bytesize ht.ByteSizeIEC
		expected string
	}{
		{0, "0B"},
		{128, "128B"},
		{32 * ht.KiB, "32K"},
		{2077 * ht.MiB, "2.0283G"},
		{2 * ht.GiB, "2G"},
	}

	for i, tt := range tests {
		got := tt.bytesize.String()
		if tt.expected != got {
			t.Errorf("index %d: %d returned: %s not equal to %s", i, tt.bytesize, got, tt.expected)
		}
	}
}
