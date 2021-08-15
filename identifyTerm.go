package avail

// identifyTerm is used to create a concrete type for a single cron term. Cron has several expression
// formats:
// * Span: Used to represent a range.  ex. 0-9
// * Wildcard: Used to represent all possible values within a certain term. ex. *
// * List: Used to represent an explicit list of values. ex. 1,2,3
// * Value: Used to represent a single value. ex. 2
//
// A cron term is a single field in a complete cron expression.
// Ex. in the expression: "0 15 10 * * *", "15" would be a term of type "value".

import "regexp"

// List of regexs that we use to match against a single cron term for identification.
var (
	spanRegex     = regexp.MustCompile(`^[0-9]+-[0-9]+$`)
	wildcardRegex = regexp.MustCompile(`^\*$`)
	listRegex     = regexp.MustCompile(`,+`)
	valueRegex    = regexp.MustCompile(`^([0-9]+)$`)
)

// termKind is an enum which represents different term kinds
type termKind string

const (
	span     termKind = "span"
	wildcard termKind = "wildcard"
	list     termKind = "list"
	value    termKind = "value"
	unknown  termKind = "unknown"
)

// termRegexToType stores the mapping between a term's regex representation
// and the concrete term type it is. This is used to help identify the term
// so that we can run the correct parser later.
var termRegexToType = map[*regexp.Regexp]termKind{
	spanRegex:     span,
	wildcardRegex: wildcard,
	listRegex:     list,
	valueRegex:    value,
}

func identifyTermKind(term string) termKind {
	for regex := range termRegexToType {
		isMatch := regex.MatchString(term)
		if isMatch {
			return termRegexToType[regex]
		}
	}

	return unknown
}
