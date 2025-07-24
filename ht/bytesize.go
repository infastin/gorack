package ht

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

var byteSizeUnitMap = map[string]uint64{
	"B": uint64(1),

	"k": uint64(KiB),
	"K": uint64(KiB),
	"M": uint64(MiB),
	"G": uint64(GiB),
	"T": uint64(TiB),
	"P": uint64(PiB),
	"E": uint64(EiB),

	"kB": uint64(KB),
	"KB": uint64(KB),
	"MB": uint64(MB),
	"GB": uint64(GB),
	"TB": uint64(TB),
	"PB": uint64(PB),
	"EB": uint64(EB),

	"KiB": uint64(KiB),
	"MiB": uint64(MiB),
	"GiB": uint64(GiB),
	"TiB": uint64(TiB),
	"PiB": uint64(PiB),
	"EiB": uint64(EiB),
}

func parseByteSize(s string) (uint64, error) {
	orig := s
	if s == "" {
		return 0, fmt.Errorf("ht: invalid byte size %q", orig)
	}

	// The next character must be [0-9.]
	if s[0] != '.' && (s[0] < '0' || s[0] > '9') {
		return 0, fmt.Errorf("ht: invalid byte size %q", orig)
	}

	var (
		v, f  uint64      // integers before, after decimal point
		scale float64 = 1 // byteSize = v + f/scale
	)

	// Consume [0-9]*
	pl := len(s)
	v, s, err := leadingInt(s)
	if err != nil {
		return 0, fmt.Errorf("ht: invalid byte size %q", orig)
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
		// no digits
		return 0, fmt.Errorf("ht: invalid byte size %q", orig)
	}

	// Skip one optional space
	if s != "" {
		if s[0] == ' ' {
			s = s[1:]
		}
		if s == "" {
			// trailing space is not allowed
			return 0, fmt.Errorf("ht: invalid byte size %q", orig)
		}
	}

	unit := uint64(1)
	if s != "" {
		var ok bool
		unit, ok = byteSizeUnitMap[s]
		if !ok {
			return 0, fmt.Errorf("ht: unknown unit %q in byte size %q", s, orig)
		}
	}
	if v > 1<<63/unit {
		// overflow
		return 0, fmt.Errorf("ht: invalid byte size %q", orig)
	}
	v *= unit
	if f > 0 {
		v += uint64(float64(f) * (float64(unit) / scale))
		if v > 1<<63 {
			// overflow
			return 0, fmt.Errorf("ht: invalid byte size %q", orig)
		}
	}

	return v, nil
}

func byteSizeString(byteSize, unit uint64) string {
	var result []byte
	if byteSize >= unit {
		div, exp := unit, 0
		for n := byteSize / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
		result = strconv.AppendFloat(result, float64(byteSize)/float64(div), 'f', 4, 64)
		result = bytes.TrimRight(result, ".0")
		result = append(result, "KMGTPE"[exp])
		if unit == 1000 {
			result = append(result, 'B')
		}
	} else {
		result = strconv.AppendUint(result, byteSize, 10)
		result = append(result, 'B')
	}
	return unsafe.String(unsafe.SliceData(result), len(result))
}

// ByteSizeSI is a byte size in SI units.
type ByteSizeSI uint64

const (
	KB ByteSizeSI = 1000
	MB            = KB * 1000
	GB            = MB * 1000
	TB            = GB * 1000
	PB            = TB * 1000
	EB            = PB * 1000
)

// ParseByteSizeSI parses a byte size string.
// A byte size string is a decimal number with an optional
// fraction and unit suffix, such as "1k", "1.5M" or "2MB".
//
// Valid IEC units are:
// - k, K, KiB - kibibyte
// - M, MiB - mebibyte
// - G, GiB - gibibyte
// - T, TiB - tebibyte
// - P, PiB - pebibyte
// - E, EiB - exbibyte
//
// Valid SI units are:
// - kb, KB - kilobyte
// - MB - megabyte
// - GB - gigabyte
// - TB - terabyte
// - PB - petabyte
// - EB - exabyte
//
// NOTE: Accepts both IEC and SI units.
func ParseByteSizeSI(s string) (ByteSizeSI, error) {
	byteSize, err := parseByteSize(s)
	if err != nil {
		return 0, err
	}
	return ByteSizeSI(byteSize), nil
}

func (b ByteSizeSI) String() string {
	return byteSizeString(uint64(b), 1000)
}

func (b ByteSizeSI) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *ByteSizeSI) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*b = ByteSizeSI(value)
		return nil
	case string:
		tmp, err := ParseByteSizeSI(value)
		if err != nil {
			return err
		}
		*b = tmp
		return nil
	default:
		return errors.New("ht: invalid byte size")
	}
}

func (b ByteSizeSI) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *ByteSizeSI) UnmarshalText(data []byte) error {
	tmp, err := ParseByteSizeSI(string(data))
	if err != nil {
		return err
	}
	*b = tmp
	return nil
}

// ByteSizeIEC is a byte size in IEC units.
type ByteSizeIEC uint64

const (
	KiB ByteSizeIEC = 1024
	MiB             = KiB * 1024
	GiB             = MiB * 1024
	TiB             = GiB * 1024
	PiB             = TiB * 1024
	EiB             = PiB * 1024
)

// ParseByteSizeIEC parses a byte size string.
// A byte size string is a decimal number with an optional
// fraction and unit suffix, such as "1k", "1.5M" or "2MB".
//
// Valid IEC units are:
// - k, K, KiB - kibibyte
// - M, MiB - mebibyte
// - G, GiB - gibibyte
// - T, TiB - tebibyte
// - P, PiB - pebibyte
// - E, EiB - exbibyte
//
// Valid SI units are:
// - kb, KB - kilobyte
// - MB - megabyte
// - GB - gigabyte
// - TB - terabyte
// - PB - petabyte
// - EB - exabyte
//
// NOTE: Accepts both IEC and SI units.
func ParseByteSizeIEC(s string) (ByteSizeIEC, error) {
	byteSize, err := parseByteSize(s)
	if err != nil {
		return 0, err
	}
	return ByteSizeIEC(byteSize), nil
}

func (b ByteSizeIEC) String() string {
	return byteSizeString(uint64(b), 1024)
}

func (b ByteSizeIEC) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *ByteSizeIEC) UnmarshalJSON(data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*b = ByteSizeIEC(value)
		return nil
	case string:
		tmp, err := ParseByteSizeIEC(value)
		if err != nil {
			return err
		}
		*b = tmp
		return nil
	default:
		return errors.New("ht: invalid byte size")
	}
}

func (b ByteSizeIEC) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *ByteSizeIEC) UnmarshalText(data []byte) error {
	tmp, err := ParseByteSizeIEC(string(data))
	if err != nil {
		return err
	}
	*b = tmp
	return nil
}
