package isstr

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var cronMonthAlt = map[string]int{
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

var cronWeekdayAlt = map[string]int{
	"sun": 0,
	"mon": 1,
	"tue": 2,
	"wed": 3,
	"thu": 4,
	"fri": 5,
	"sat": 6,
}

func cronValid(cron string) error {
	scanner := bufio.NewScanner(strings.NewReader(cron))
	scanner.Split(bufio.ScanWords)

	const minFields = 4
	const maxFields = 6

	fields := make([]string, 0, maxFields)

	for scanner.Scan() {
		if len(fields) == maxFields {
			return fmt.Errorf("at most %d fields are allowed", maxFields)
		}
		fields = append(fields, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(fields) == 1 && fields[0][0] == '@' {
		field := fields[0]
		if field[0] == '@' {
			return cronValidPreset(field)
		}
	}

	if len(fields) < minFields {
		return fmt.Errorf("at least %d fields are required", minFields)
	}

	return cronValidFields(fields)
}

func cronValidPreset(field string) error {
	switch field {
	case "@hourly", "@daily", "@midnigth", "@weekly", "@monthly", "@annually", "@yearly":
		return nil
	default:
		return fmt.Errorf("unknown preset: %s", field)
	}
}

func cronValidFields(fields []string) error {
	if len(fields) == 6 {
		if err := cronValidField(fields[0], 0, 59, nil); err != nil {
			return fmt.Errorf("invalid seconds field: %w", err)
		}
		fields = fields[1:]
	}

	if err := cronValidField(fields[0], 0, 59, nil); err != nil {
		return fmt.Errorf("invalid minutes field: %w", err)
	}

	if err := cronValidField(fields[1], 0, 23, nil); err != nil {
		return fmt.Errorf("invalid hours field: %w", err)
	}

	if err := cronValidField(fields[2], 1, 31, nil); err != nil {
		return fmt.Errorf("invalid day of month field: %w", err)
	}

	if err := cronValidField(fields[3], 1, 12, cronMonthAlt); err != nil {
		return fmt.Errorf("invalid month field: %w", err)
	}

	if len(fields) == 5 {
		if err := cronValidField(fields[4], 0, 7, cronWeekdayAlt); err != nil {
			return fmt.Errorf("invalid day of week field: %w", err)
		}
	}

	return nil
}

func cronValidField(field string, min, max int, alt map[string]int) error {
	value := field
	for n := 1; field != ""; n++ {
		pos := strings.IndexByte(field, ',')
		if pos != -1 {
			value = field[:pos]
			field = field[pos+1:]
		} else {
			field = ""
		}
		if err := cronValidValue(value, min, max, alt); err != nil {
			if pos == -1 && n == 1 {
				return fmt.Errorf("invalid value: %w", err)
			}
			return fmt.Errorf("invalid %d'th value in list: %w", n, err)
		}
	}
	return nil
}

func cronValidValue(value string, min, max int, alt map[string]int) (err error) {
	var (
		ch      byte
		rest    string
		low     = value
		lowNum  = min
		highNum = max
	)

	pos := strings.IndexAny(value, "-/")
	if pos != -1 {
		ch = value[pos]
		low = value[:pos]
		rest = value[pos+1:]
	}

	if low == "*" || low == "?" {
		if ch == '-' {
			return fmt.Errorf("range is not allowed with '%s'", low)
		}
	} else if lowNum, err = cronValidBasicValue(low, min, max, alt); err != nil {
		if ch == '-' {
			return fmt.Errorf("invalid low value of range: %w", err)
		}
		if ch == '/' {
			return fmt.Errorf("invalid base of step value: %w", err)
		}
		return err
	}

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

		highNum, err = cronValidBasicValue(high, min, max, alt)
		if err != nil {
			return fmt.Errorf("invalid high value of range: %w", err)
		}

		if lowNum > highNum {
			return fmt.Errorf("invalid range: %d > %d", lowNum, highNum)
		}
	}

	if ch == '/' {
		step, err := cronValidNumber(rest, min, max)
		if err != nil {
			return fmt.Errorf("invalid step: %w", err)
		}
		if step > (highNum - lowNum) {
			return fmt.Errorf("step %d is greater than range %d-%d", step, lowNum, highNum)
		}
	}

	return nil
}

func cronValidBasicValue(value string, min, max int, alt map[string]int) (int, error) {
	if alt != nil {
		val, ok := alt[strings.ToLower(value)]
		if ok {
			return val, nil
		}
	}
	val, err := cronValidNumber(value, min, max)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func cronValidNumber(value string, min, max int) (int, error) {
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
