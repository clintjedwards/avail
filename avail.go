package avail

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// fieldType is an enum which represents different parts of a total cron expression.
// For example in the expression "0 10 15 * * *", 0 would be of type "minute".
type fieldType string

const (
	minute  fieldType = "minute"
	hour    fieldType = "hour"
	day     fieldType = "day"
	month   fieldType = "month"
	weekday fieldType = "weekday"
	year    fieldType = "year"
)

var cronExpressionRegex = regexp.MustCompile(`^((((\d+,)+\d+|(\d+(-)\d+)|\d+|\*) ?){6})$`)

// ParsedExpression represents a breakdown of a given cron time expression
type ParsedExpression struct {
	Minutes  Field
	Hours    Field
	Days     Field
	Months   Field
	Weekdays Field
	Years    Field
}

// Timeframe represents both the raw cron expression and the datastructures used to represent that
// expression for easy checking
type Timeframe struct {
	Expression       string // * * * * * * 6 fields - min, hours, day of month, month, day of week, year
	ParsedExpression ParsedExpression
}

// New will parse the given cron expression and allow user to check if the time given is within
func New(expression string) (Timeframe, error) {
	isMatch := cronExpressionRegex.MatchString(expression)
	if !isMatch {
		return Timeframe{}, fmt.Errorf("could not parse cron expression: %s; misformatted expression", expression)
	}

	terms := strings.Split(expression, " ")
	// we need this extra check to make sure there are the proper amount of fields because I am bad at regex
	if len(terms) != 6 {
		return Timeframe{}, fmt.Errorf("could not parse cron expression: %s; must have 6 terms", expression)
	}

	minutes, err := newField(minute, terms[0], 0, 59)
	if err != nil {
		return Timeframe{}, err
	}
	hours, err := newField(hour, terms[1], 0, 23)
	if err != nil {
		return Timeframe{}, err
	}
	day, err := newField(day, terms[2], 1, 31)
	if err != nil {
		return Timeframe{}, err
	}
	month, err := newField(month, terms[3], 1, 12)
	if err != nil {
		return Timeframe{}, err
	}
	weekday, err := newField(weekday, terms[4], 0, 6)
	if err != nil {
		return Timeframe{}, err
	}
	year, err := newField(year, terms[5], 1970, 2100)
	if err != nil {
		return Timeframe{}, err
	}

	return Timeframe{
		Expression: expression,
		ParsedExpression: ParsedExpression{
			Minutes:  minutes,
			Hours:    hours,
			Days:     day,
			Months:   month,
			Weekdays: weekday,
			Years:    year,
		},
	}, nil
}

// Able will evaluate if the time given is within the cron expression.
func (a *Timeframe) Able(time time.Time) bool {
	fieldTypes := []fieldType{
		minute,
		hour,
		day,
		month,
		weekday,
		year,
	}

	for _, field := range fieldTypes {
		switch field {
		case minute:
			if _, ok := a.ParsedExpression.Minutes.Values[time.Minute()]; !ok {
				return false
			}
		case hour:
			if _, ok := a.ParsedExpression.Hours.Values[time.Hour()]; !ok {
				return false
			}
		case day:
			if _, ok := a.ParsedExpression.Days.Values[time.Day()]; !ok {
				return false
			}
		case month:
			if _, ok := a.ParsedExpression.Months.Values[int(time.Month())]; !ok {
				return false
			}
		case weekday:
			if _, ok := a.ParsedExpression.Weekdays.Values[int(time.Weekday())]; !ok {
				return false
			}
		case year:
			if _, ok := a.ParsedExpression.Years.Values[time.Year()]; !ok {
				return false
			}
		}
	}

	return true
}

func generateSequentialSet(start, end int) map[int]struct{} {
	set := map[int]struct{}{}
	for i := start; i < end+1; i++ {
		set[i] = struct{}{}
	}
	return set
}
