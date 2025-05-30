package isstr_test

import (
	"testing"

	isstr "github.com/infastin/gorack/validation/is/str"
)

func TestCRON(t *testing.T) {
	tests := []struct {
		cron  string
		valid bool
	}{
		{"", false},
		{"*", false},
		{"* *", false},
		{"* * *", false},
		{"* * * *", true},
		{"* * * * *", true},
		{"* * * * * *", true},

		{"0 0 0 0", false},
		{"0 0 1 0", false},
		{"* * 1 1", true},
		{"* * 1 1 *", true},

		{"0/15 * * * *", true},
		{"5/15 * * * *", true},

		{"59 * * * * *", true},
		{"60 * * * * *", false},
		{"* 59 * * * *", true},
		{"* 60 * * * *", false},
		{"* * 23 * * *", true},
		{"* * 24 * * *", false},
		{"* * * 31 * *", true},
		{"* * * 32 * *", false},
		{"* * * * 12 *", true},
		{"* * * * 13 *", false},
		{"* * * * * 6", true},
		{"* * * * * 7", false},

		{"30 08 ? Jul Sun", true},
		{"30 08 * Jul Sun", true},
		{"30 08 * Jul Sun", true},
		{"30 08 15 Jul *", true},

		{"@hourly", true},
		{"@daily", true},
		{"@midnigth", true},
		{"@monthly", true},
		{"@yearly", true},
		{"@annually", true},

		{"* * 1-15/2 * *", true},
		{"* * 1-15/15 * *", false},

		{"* * * Jun-Aug *", true},
		{"* * * Aug-Jun *", false},
		{"* * * * Mon-Sun", false},
		{"* * * * Sun-Mon", true},

		{"0 22 * * 1-5", true},
		{"23 0-20/2 * * *", true},
		{"23 0-20/2,1,15,3-14 * * sun-sun", true},
		{"5 4 * * sun", true},
		{"0 0,12 1 */2 *", true},
		{"0 0 1,15 * 3", true},
	}
	for _, tt := range tests {
		err := isstr.CRON(tt.cron)
		if err != nil && tt.valid {
			t.Errorf("unexpected error: cron=%s error=%s", tt.cron, err.Error())
		} else if err == nil && !tt.valid {
			t.Errorf("expected an error: cron=%s", tt.cron)
		}
	}
}
