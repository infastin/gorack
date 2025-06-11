package calev

import (
	"math/bits"
	"time"
)

// Spec represents a set of constraints for a date and time,
// including the month, day, weekday, hour, minute, and second,
// and is used to find the next date and time that satisfies these constraints.
//
// Rules determining the day of the month based on a combination of:
//   - specific days of the month
//   - last days of the month (e.g. last 5 days)
//   - days of the week (e.g. Monday, Tuesday)
//   - restricted days of the week, which refines the selection to only weekdays
//     that fall within the specified days of the month or last days of the month
//
// When multiple constraints from the set of specific days of the month,
// last days of the month, and days of the week are present,
// the nearest day that meets any of these conditions will be selected,
// effectively applying a logical OR operation between these constraints.
// The restricted days of the week constraint, on the other hand,
// acts as an additional filter that narrows down the selection to only those days
// that also meet the restricted weekday criteria.
//
// If neither days of the month nor last days of the month are specified,
// restricted days of the week behaves the same as days of the week.
//
// By default, Spec allows any month, day, hour, and minute, but only allows a second of 0,
// unless these constraints are explicitly set to other values.
type Spec struct {
	months         uint16
	days           uint32
	ldays          uint32
	weekdays       uint8
	weekdaysStrict uint8
	hours          uint32
	minutes        uint64
	seconds        uint64
}

// Makes new Spec with the given options.
func New(opts ...SpecOpt) *Spec {
	s := new(Spec)
	for _, opt := range opts {
		opt(s)
	}
	s.init()
	return s
}

func (s *Spec) init() {
	if s.months == 0 {
		s.months = 1<<13 - 1
	}
	if s.days == 0 && s.ldays == 0 {
		if s.weekdays == 0 {
			s.days = 1<<32 - 1
		} else if s.weekdaysStrict != 0 {
			s.weekdays |= s.weekdaysStrict
			s.weekdaysStrict = 0
		}
	}
	if s.hours == 0 {
		s.hours = 1<<24 - 1
	}
	if s.minutes == 0 {
		s.minutes = 1<<61 - 1
	}
}

// Returns the nearest date and time that satisfies Spec,
// or the next one if the given date and time already satisfy it.
// Returns zero if no such date and time exist.
func (s *Spec) Next(t time.Time) time.Time {
	if s.seconds != 0 {
		t = t.Truncate(time.Second)
	} else {
		t = t.Truncate(time.Minute)
	}
	cur := t

wrap:
	cur = s.nextMonth(cur)

	cur, wrapped, ok := s.nextDay(cur)
	if !ok {
		return time.Time{}
	}
	if wrapped {
		goto wrap
	}

	cur, wrapped = s.nextHour(cur)
	if wrapped {
		goto wrap
	}

	cur, wrapped = s.nextMinute(cur)
	if wrapped {
		goto wrap
	}

	if s.seconds != 0 {
		cur, wrapped = s.nextSecond(cur)
		if wrapped {
			goto wrap
		}
	}

	if cur.Equal(t) {
		if s.seconds != 0 {
			cur = cur.Add(time.Second)
		} else {
			cur = cur.Add(time.Minute)
		}
		goto wrap
	}

	return cur
}

func (s *Spec) nextMonth(t time.Time) time.Time {
	year := t.Year()
	month := int(t.Month()) - 1

	var nextMonth int

	restMonths := s.months &^ (1<<month - 1)
	if restMonths != 0 {
		nextMonth = bits.TrailingZeros16(restMonths)
		if nextMonth == month {
			return t
		}
	} else {
		restMonths = s.months & (1<<month - 1)
		nextMonth = bits.TrailingZeros16(restMonths)
		year++
	}

	return time.Date(year, time.Month(nextMonth+1), 1, 0, 0, 0, 0, t.Location())
}

func (s *Spec) nextDay(t time.Time) (nt time.Time, wrapped, ok bool) {
	if s.days != 0 || s.weekdays != 0 {
		nt, wrapped, ok = s.nextMonthDay(t)
	}
	if s.ldays != 0 {
		nt1, wrapped1, ok1 := s.nextMonthLastDay(t)
		if !ok || (wrapped && !wrapped1) || nt1.Before(nt) {
			nt, wrapped, ok = nt1, wrapped1, ok1
		}
	}
	return nt, wrapped, ok
}

func (s *Spec) nextMonthDay(t time.Time) (nt time.Time, wrapped, ok bool) {
	year, month, day := t.Date()
	day--

	var nextDay int

	days := s.days
	if s.weekdaysStrict != 0 || s.weekdays != 0 {
		firstWd := int(time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).Weekday())
		if s.weekdaysStrict != 0 {
			days &= wdToMd(s.weekdaysStrict, firstWd)
		}
		if s.weekdays != 0 {
			days |= wdToMd(s.weekdays, firstWd)
		}
		if days == 0 {
			nt = time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
			return nt, true, true
		}
	}

	restDays := days &^ (1<<day - 1)
	if restDays != 0 {
		nextDay = bits.TrailingZeros32(restDays)
		if nextDay == day {
			return t, false, true
		}
		if !datePossible(year, month, nextDay+1, t.Location()) {
			if !s.dayPossible(nextDay + 1) {
				return time.Time{}, false, false
			}
			nt = time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
			return nt, true, true
		}
	} else {
		month++
		wrapped = true

		days := s.days
		if s.weekdaysStrict != 0 || s.weekdays != 0 {
			firstWd := int(time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).Weekday())
			if s.weekdaysStrict != 0 {
				days &= wdToMd(s.weekdaysStrict, firstWd)
			}
			if s.weekdays != 0 {
				days |= wdToMd(s.weekdays, firstWd)
			}
			if days == 0 {
				nt = time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
				return nt, true, true
			}
		}

		restDays = days & (1<<day - 1)
		nextDay = bits.TrailingZeros32(restDays)
		if !datePossible(year, month, nextDay+1, t.Location()) {
			if !s.dayPossible(nextDay + 1) {
				return time.Time{}, false, false
			}
			nt = time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
			return nt, true, true
		}
	}

	nt = time.Date(year, month, nextDay+1, 0, 0, 0, 0, t.Location())
	return nt, wrapped, true
}

func (s *Spec) nextMonthLastDay(t time.Time) (nt time.Time, wrapped, ok bool) {
	year, month, day := t.Date()
	day--

	var nextDay int

	maxDays := monthDays(year, month, t.Location())
	sh := 31 - maxDays

	// Last days of the month are represented as a set of bits
	// where the 30th bit represents the last day of the month,
	// the 29th bit represents the second last day of the month, etc.
	//
	// This means that if there are no 31 days in a month, we have to shift
	// masks to the left by the difference between 31 and the number of days in the month,
	// so that 30th bit still represents the last day of the month.
	//
	// This also means that we have to subtract the shift amount
	// when calculating the next day of a month.

	ldays := s.ldays
	if s.weekdaysStrict != 0 {
		firstWd := int(time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).Weekday())
		ldays &= wdToMd(s.weekdaysStrict, firstWd) << sh
		if ldays == 0 {
			nt = time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
			return nt, true, true
		}
	}

	restDays := ldays &^ (1<<(day+sh) - 1)

	if restDays != 0 {
		nextDay = bits.TrailingZeros32(restDays) - sh
		if nextDay == day {
			return t, false, true
		}
	} else {
		closestDay := 30 - (bits.LeadingZeros32(ldays) - 1) - sh
		if closestDay < 0 && !s.dayPossible(maxDays-closestDay) {
			return time.Time{}, false, false
		}

		month++
		wrapped = true

		maxDays := monthDays(year, month, t.Location())
		nextSh := 31 - maxDays

		ldays := s.ldays
		if s.weekdaysStrict != 0 {
			firstWd := int(time.Date(year, month, 1, 0, 0, 0, 0, t.Location()).Weekday())
			ldays &= wdToMd(s.weekdaysStrict, firstWd) << nextSh
			if ldays == 0 {
				nt = time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
				return nt, true, true
			}
		}

		restDays = ldays & (1<<(day+sh) - 1)
		nextDay = bits.TrailingZeros32(restDays) - nextSh
	}

	nt = time.Date(year, month, nextDay+1, 0, 0, 0, 0, t.Location())
	return nt, wrapped, true
}

func (s *Spec) nextHour(t time.Time) (nt time.Time, wrapped bool) {
	day := t.Day()
	hour := t.Hour()

	var nextHour int

	restHours := s.hours &^ (1<<hour - 1)
	if restHours != 0 {
		nextHour = bits.TrailingZeros32(restHours)
		if nextHour == hour {
			return t, false
		}
	} else {
		restHours = s.hours & (1<<hour - 1)
		nextHour = bits.TrailingZeros32(restHours)
		day++
		wrapped = true
	}

	nt = time.Date(t.Year(), t.Month(), day, nextHour, 0, 0, 0, t.Location())

	// Handle DST.
	if hour := nt.Hour(); hour != nextHour {
		day++
		wrapped = true
		nt = time.Date(t.Year(), t.Month(), day, nextHour, 0, 0, 0, t.Location())
	}

	return nt, wrapped
}

func (s *Spec) nextMinute(t time.Time) (nt time.Time, wrapped bool) {
	hour := t.Hour()
	minute := t.Minute()

	var nextMinute int

	restMinutes := s.minutes &^ (1<<minute - 1)
	if restMinutes != 0 {
		nextMinute = bits.TrailingZeros64(restMinutes)
		if nextMinute == minute {
			return t, false
		}
	} else {
		restMinutes = s.minutes & (1<<minute - 1)
		nextMinute = bits.TrailingZeros64(restMinutes)
		hour++
		wrapped = true
	}

	year, month, day := t.Date()
	nt = time.Date(year, month, day, hour, nextMinute, 0, 0, t.Location())

	return nt, wrapped
}

func (s *Spec) nextSecond(t time.Time) (nt time.Time, wrapped bool) {
	minute := t.Minute()
	second := t.Second()

	var nextSecond int

	restSeconds := s.seconds &^ (1<<second - 1)
	if restSeconds != 0 {
		nextSecond = bits.TrailingZeros64(restSeconds)
		if nextSecond == second {
			return t, false
		}
	} else {
		restSeconds = s.seconds & (1<<second - 1)
		nextSecond = bits.TrailingZeros64(restSeconds)
		minute++
		wrapped = true
	}

	year, month, day := t.Date()
	nt = time.Date(year, month, day, t.Hour(), minute, nextSecond, 0, t.Location())

	return nt, wrapped
}

func (s *Spec) dayPossible(day int) bool {
	for i := 0; i < 12; i++ {
		if s.months&(1<<i) == 0 {
			continue
		}
		if monthContainsDay(time.Month(i+1), day) {
			return true
		}
	}
	return false
}

func monthContainsDay(month time.Month, day int) bool {
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return day >= 1 && day <= 31
	case time.February:
		return day >= 1 && day <= 29
	case time.April, time.June, time.September, time.November:
		return day >= 1 && day <= 30
	}
	return false
}

func wdToMd(x uint8, s int) uint32 {
	const n = 7

	s %= n
	x = x & (1<<n - 1)
	x = x>>s | x<<(n-s)

	x0 := uint32(x)
	x1 := uint32(x) << 7
	x2 := uint32(x) << 14
	x3 := uint32(x) << 21
	x4 := uint32(x) << 28

	return x0 | x1 | x2 | x3 | x4
}

func datePossible(year int, month time.Month, day int, loc *time.Location) bool {
	if month > 12 {
		year += int(month-1) / 12
		month = (month-1)%12 + 1
	}
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return day >= 1 && day <= 31
	case time.February:
		if day <= 28 {
			return day >= 1 && day <= 28
		}
		maxDays := 32 - time.Date(year, month, 32, 0, 0, 0, 0, loc).Day()
		return day >= 1 && day <= maxDays
	case time.April, time.June, time.September, time.November:
		return day >= 1 && day <= 30
	}
	return false
}

func monthDays(year int, month time.Month, loc *time.Location) int {
	if month > 12 {
		year += int(month-1) / 12
		month = (month-1)%12 + 1
	}
	switch month {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.February:
		return 32 - time.Date(year, month, 32, 0, 0, 0, 0, loc).Day()
	case time.April, time.June, time.September, time.November:
		return 30
	}
	return 0
}
