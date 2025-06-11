package calev_test

import (
	"fmt"
	"time"

	"github.com/infastin/gorack/calev"
)

func ExampleParseCron() {
	// Prefix a day with '^' to get n'th last day of a month.
	spec, err := calev.ParseCron("0 0 ^1 Feb *", calev.CronOptions{})
	if err != nil {
		panic(err)
	}

	date := time.Date(2028, time.January, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Last day of February: %s\n", spec.Next(date))

	// Prefix a day of week with '&' to get the day of week within the specified days of a month.
	spec, err = calev.ParseCron("0 0 14-26 * &Mon", calev.CronOptions{})
	if err != nil {
		panic(err)
	}

	date = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Monday within 14-26 of January: %s\n", spec.Next(date))

	// '^' and '&' can be combined.
	spec, err = calev.ParseCron("0 0 ^1-7 * &Mon", calev.CronOptions{})
	if err != nil {
		panic(err)
	}

	date = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Monday within 7 last days of January: %s\n", spec.Next(date))

	// Output:
	// Last day of February: 2028-02-29 00:00:00 +0000 UTC
	// Monday within 14-26 of January: 2025-01-20 00:00:00 +0000 UTC
	// Monday within 7 last days of January: 2025-01-27 00:00:00 +0000 UTC
}
