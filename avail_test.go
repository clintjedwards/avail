package avail

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseable(t *testing.T) {
	tests := map[string]struct {
		expression string
	}{
		"wildcard": {
			expression: "* * * * * *",
		},
		"ranges": {
			expression: "* * 25-30 12 * 2020",
		},
		"single values": {
			expression: "* * 1 1 * *",
		},
		"range + single value": {
			expression: "* * * 6,7,8 * 2020",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := New(tc.expression)
			if err != nil {
				t.Errorf("expression %s should be parsed successfully", tc.expression)
			}
		})
	}
}

func TestUnparseable(t *testing.T) {
	tests := map[string]struct {
		expression string
	}{
		"too many arguments": {
			expression: "* * * * * * *",
		},
		"too few arguments": {
			expression: "* * * *",
		},
		"out of bounds single value": {
			expression: "* * * * 22222 *",
		},
		"out of bounds range": {
			expression: "10-500 * * * * *",
		},
		"out of bounds list": {
			expression: "* 1,40,100 * * * *",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := New(tc.expression)
			if err == nil {
				t.Errorf("expression %s should not be parsed successfully", tc.expression)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		expression string
		want       Avail
	}{
		"wildcard": {"* * * * * *", Avail{
			Expression: "* * * * * *",
			ParsedExpression: ParsedExpression{
				Minutes: Field{
					Kind:   minute,
					Term:   "*",
					Min:    0,
					Max:    59,
					Values: generateSequentialSet(0, 59),
				},
				Hours: Field{
					Kind:   hour,
					Term:   "*",
					Min:    0,
					Max:    23,
					Values: generateSequentialSet(0, 23),
				},
				Days: Field{
					Kind:   day,
					Term:   "*",
					Min:    1,
					Max:    31,
					Values: generateSequentialSet(1, 31),
				},
				Months: Field{
					Kind:   month,
					Term:   "*",
					Min:    1,
					Max:    12,
					Values: generateSequentialSet(1, 12),
				},
				Weekdays: Field{
					Kind:   weekday,
					Term:   "*",
					Min:    0,
					Max:    6,
					Values: generateSequentialSet(0, 6),
				},
				Years: Field{
					Kind:   year,
					Term:   "*",
					Min:    1970,
					Max:    2100,
					Values: generateSequentialSet(1970, 2100),
				},
			},
		}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := New(tc.expression)
			if err != nil {
				t.Error(err)
			}

			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Errorf("result is different than expected(-want +got):\n%s", diff)
			}
		})
	}
}

func TestParseWildcard(t *testing.T) {
	want := Field{
		Kind:   minute,
		Term:   "*",
		Min:    0,
		Max:    59,
		Values: generateSequentialSet(0, 59),
	}
	got, err := newField(minute, "*", 0, 59)
	if err != nil {
		t.Error(err)
	}

	diff := cmp.Diff(want, got)
	if diff != "" {
		t.Errorf("result is different than expected(-want +got):\n%s", diff)
	}
}

func TestParseSpan(t *testing.T) {
	want := Field{
		Kind:   hour,
		Term:   "4-14",
		Min:    0,
		Max:    23,
		Values: generateSequentialSet(4, 14),
	}
	got, err := newField(hour, "4-14", 0, 23)
	if err != nil {
		t.Error(err)
	}

	diff := cmp.Diff(want, got)
	if diff != "" {
		t.Errorf("result is different than expected(-want +got):\n%s", diff)
	}
}

func TestAble(t *testing.T) {

	tests := map[string]struct {
		expression string
		time       time.Time
		want       bool
	}{
		"wildcard":           {"* * * * * *", time.Now(), true},
		"year out of range":  {"* * * * * 2019", time.Now(), false},
		"year list in range": {"* * * * * 2019,2020,2021", time.Date(2020, 1, 1, 1, 1, 1, 1, time.UTC), true},
		"specific day of week": {"* * * 6 2 2020",
			time.Date(2020, 6, 2, 1, 1, 1, 1, time.UTC), true},
		"specific day of week; out of range": {"* * * 6 2 2020",
			time.Date(2020, 8, 4, 1, 1, 1, 1, time.UTC), false},
		"every day at noon in January only": {"0 12 * 1 * *",
			time.Date(2020, 1, 24, 12, 0, 0, 0, time.UTC), true},
		"every day from 6am to 2pm": {"* 6-14 * * * *",
			time.Date(2020, 1, 24, 12, 0, 0, 0, time.UTC), true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			avail, err := New(tc.expression)
			if err != nil {
				t.Error(err)
			}

			if avail.Able(tc.time) != tc.want {
				t.Errorf("want %t, got %t", tc.want, !tc.want)
			}
		})
	}

}

func ExampleAvail_Able() {
	avail, _ := New("* * * * * *")

	now := time.Now()

	fmt.Println(avail.Able(now))
	// Output: true
}
