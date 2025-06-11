package calev

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	cronMonthAlt = map[string]int{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}

	cronWeekdayAlt = map[string]int{
		"sun": 0,
		"mon": 1,
		"tue": 2,
		"wed": 3,
		"thu": 4,
		"fri": 5,
		"sat": 6,
	}
)

type CronOptions struct {
	Seconds         bool // Must contain seconds field.
	SecondsOptional bool // May contain seconds field.
	WeekdayOptional bool // Day of week field may be omitted.
}

// Parses cron expression and returns Spec that represents it.
// Cron expression can contain last days of month and restricted days of week.
// Last days of week must be preceded by '^' character.
// Restricted days of week must be preceded by '&' character.
func ParseCron(cron string, opts CronOptions) (*Spec, error) {
	scanner := bufio.NewScanner(strings.NewReader(cron))
	scanner.Split(bufio.ScanWords)

	minFields, maxFields := 5, 5
	if opts.Seconds || opts.SecondsOptional {
		maxFields++
	}
	if opts.Seconds {
		minFields++
	}
	if opts.WeekdayOptional {
		minFields--
	}

	fields := make([]string, 0, maxFields)

	for scanner.Scan() {
		if len(fields) == maxFields {
			return nil, fmt.Errorf("at most %d fields are allowed", maxFields)
		}
		fields = append(fields, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(fields) == 1 && fields[0][0] == '@' {
		field := fields[0]
		if field[0] == '@' {
			return parsePreset(field)
		}
	}

	if len(fields) < minFields {
		return nil, fmt.Errorf("at least %d fields are required", minFields)
	}

	return parseFields(fields)
}

func parsePreset(field string) (*Spec, error) {
	s := new(Spec)

	switch field {
	case "@hourly": // 0 * * * *
		s.minutes = 1 << 0
	case "@daily", "@midnigth": // 0 0 * * *
		s.minutes = 1 << 0
		s.hours = 1 << 0
	case "@weekly": // 0 0 * * 0
		s.minutes = 1 << 0
		s.hours = 1 << 0
		s.weekdays = 1 << 0
	case "@monthly": // 0 0 1 * *
		s.minutes = 1 << 0
		s.hours = 1 << 0
		s.days = 1 << 0
	case "@annually", "@yearly": // 0 0 1 1 *
		s.minutes = 1 << 0
		s.hours = 1 << 0
		s.days = 1 << 0
		s.months = 1 << 0
	default:
		return nil, fmt.Errorf("unknown preset: %s", field)
	}

	s.init()
	return s, nil
}

func parseFields(fields []string) (s *Spec, err error) {
	s = new(Spec)

	if len(fields) == 6 {
		s.seconds, err = cronParseField[uint64](fields[0], 0, 59, nil)
		if err != nil {
			return nil, fmt.Errorf("invalid seconds field: %w", err)
		}
		fields = fields[1:]
	}

	s.minutes, err = cronParseField[uint64](fields[0], 0, 59, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid minutes field: %w", err)
	}

	s.hours, err = cronParseField[uint32](fields[1], 0, 23, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid hours field: %w", err)
	}

	s.days, s.ldays, err = cronParseDay(fields[2])
	if err != nil {
		return nil, fmt.Errorf("invalid day of month field: %w", err)
	}

	s.months, err = cronParseField[uint16](fields[3], 1, 12, cronMonthAlt)
	if err != nil {
		return nil, fmt.Errorf("invalid month field: %w", err)
	}

	if len(fields) == 5 {
		s.weekdays, s.weekdaysStrict, err = cronParseWeekday(fields[4])
		if err != nil {
			return nil, fmt.Errorf("invalid day of week field: %w", err)
		}
	}

	s.init()
	return s, nil
}

func cronParseDay(field string) (days, lastDays uint32, err error) {
	for n := 1; field != ""; n++ {
		var value string
		commaPos := strings.IndexByte(field, ',')
		if commaPos != -1 {
			value = field[:commaPos]
			field = field[commaPos+1:]
		} else {
			value = field
			field = ""
		}
		var reverse bool
		if value[0] == '^' {
			value = value[1:]
			reverse = true
		}
		vi, err := cronParseValue[uint32](value, 1, 31, nil, reverse)
		if err != nil {
			if commaPos == -1 && n == 1 {
				return 0, 0, fmt.Errorf("invalid value: %w", err)
			}
			return 0, 0, fmt.Errorf("invalid %d'th value in list: %w", n, err)
		}
		if reverse {
			lastDays |= vi
		} else {
			days |= vi
		}
	}
	return days, lastDays, nil
}

func cronParseWeekday(field string) (weekdays, weekdaysStrict uint8, err error) {
	for n := 1; field != ""; n++ {
		var value string
		commaPos := strings.IndexByte(field, ',')
		if commaPos != -1 {
			value = field[:commaPos]
			field = field[commaPos+1:]
		} else {
			value = field
			field = ""
		}
		v := &weekdays
		if value[0] == '&' {
			value = value[1:]
			v = &weekdaysStrict
		}
		vi, err := cronParseValue[uint8](value, 0, 7, cronWeekdayAlt, false)
		if err != nil {
			if commaPos == -1 && n == 1 {
				return 0, 0, fmt.Errorf("invalid value: %w", err)
			}
			return 0, 0, fmt.Errorf("invalid %d'th value in list: %w", n, err)
		}
		*v |= vi
	}
	// 0 and 7 are both for Sunday, but time.Sunday is 0,
	// so we clear 7'th bit and set 0 bit in both sets.
	if weekdays&(1<<7) != 0 {
		weekdays = weekdays&^(1<<7) | 1<<0
	}
	if weekdaysStrict&(1<<7) != 0 {
		weekdaysStrict = weekdaysStrict&^(1<<7) | 1<<0
	}
	return weekdays, weekdaysStrict, nil
}

func cronParseField[T uintSet](field string, min, max int, alt map[string]int) (v T, err error) {
	for n := 1; field != ""; n++ {
		var value string
		commaPos := strings.IndexByte(field, ',')
		if commaPos != -1 {
			value = field[:commaPos]
			field = field[commaPos+1:]
		} else {
			value = field
			field = ""
		}
		vi, err := cronParseValue[T](value, min, max, alt, false)
		if err != nil {
			if commaPos == -1 && n == 1 {
				return 0, fmt.Errorf("invalid value: %w", err)
			}
			return 0, fmt.Errorf("invalid %d'th value in list: %w", n, err)
		}
		v |= vi
	}
	return v, nil
}

func cronParseValue[T uintSet](value string, min, max int, alt map[string]int, reverse bool) (v T, err error) {
	var (
		ch      byte
		rest    string
		low     = value
		lowNum  = min
		highNum = max
		step    = 1
	)

	pos := strings.IndexAny(value, "-/")
	if pos != -1 {
		ch = value[pos]
		low = value[:pos]
		rest = value[pos+1:]
	}

	if low == "*" || low == "?" {
		if ch == '-' {
			return 0, fmt.Errorf("range is not allowed with '%s'", low)
		}
		if ch == 0 {
			return 0, nil
		}
	} else if lowNum, err = cronParseBasicValue(low, min, max, alt); err != nil {
		if ch == '-' {
			return 0, fmt.Errorf("invalid low value of range: %w", err)
		}
		if ch == '/' {
			return 0, fmt.Errorf("invalid base of step value: %w", err)
		}
		return 0, err
	}

	if ch != 0 {
		if ch == '-' {
			pos := strings.IndexByte(rest, '/')

			high := rest
			if pos != -1 {
				high = rest[:pos]
				ch = '/'
				rest = rest[pos+1:]
			} else {
				ch = 0
				rest = ""
			}

			highNum, err = cronParseBasicValue(high, min, max, alt)
			if err != nil {
				return 0, fmt.Errorf("invalid high value of range: %w", err)
			}

			if lowNum > highNum {
				return 0, fmt.Errorf("invalid range: %d > %d", lowNum, highNum)
			}
		}

		if ch == '/' {
			step, err = cronParseNumber(rest, min, max)
			if err != nil {
				return 0, fmt.Errorf("invalid step: %w", err)
			}
			if step > (highNum - lowNum) {
				return 0, fmt.Errorf("step %d is greater than range %d-%d", step, lowNum, highNum)
			}
		}
	} else {
		highNum = lowNum
		step = 0
	}

	lowNum -= min
	highNum -= min

	if reverse {
		v = rangeToSetReverse[T](lowNum, highNum, step, max-min)
	} else {
		v = rangeToSet[T](lowNum, highNum, step)
	}

	return v, nil
}

func cronParseBasicValue(value string, min, max int, alt map[string]int) (int, error) {
	if alt != nil {
		val, ok := alt[strings.ToLower(value)]
		if ok {
			return val, nil
		}
	}
	val, err := cronParseNumber(value, min, max)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func cronParseNumber(value string, min, max int) (int, error) {
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return 0, errors.New("not a number")
		}
	}
	num, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return 0, err
	}
	val := int(num)
	if val < min {
		return 0, fmt.Errorf("must not be lower than %d", min)
	}
	if val > max {
		return 0, fmt.Errorf("must not be greater than %d", max)
	}
	return val, nil
}
