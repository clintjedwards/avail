package avail

import (
	"fmt"
	"strconv"
	"strings"
)

// Field represents a single value of a cron expression sometimes called a term
// Ex. in the expression: "0 15 10 * * *", "15" would be a field.
//
// Fields have map sets so that they're easy/efficent to check.
type Field struct {
	// TODO(clintjedwards): there are some hardcoded values that exist for each field,
	// it might make sense to break this apart into its own struct so that it's easier
	// for future uses of field to be consistent
	// Example: kind, min, max values are set manually by the calling function, but almost
	// never change and as such should be hardcoded somewhere instead and the call to
	// create a field should just embed the type
	Kind fieldType
	// Term is a single field in a complete cron expression.
	// Ex. in the expression: "0 15 10 * * *", "15" would be a term.
	Term     string
	Min, Max int // The maximum and minimum values for this specific field
	// Values are sets made with structs because empty structs are 0 bytes.
	// https://dave.cheney.net/2014/03/25/the-empty-struct
	Values map[int]struct{}
}

// newField takes parameters for a given cron term and attempts to parse and returns values for it
func newField(kind fieldType, term string, min, max int) (Field, error) {
	newField := Field{
		Kind: kind,
		Term: term,
		Min:  min,
		Max:  max,
	}

	err := newField.parse()
	if err != nil {
		return Field{}, err
	}

	return newField, nil
}

// parse returns a representation of the field as a set of values
// Example: A term of "1-5" will produce "1,2,3,4,5"
func (f *Field) parse() error {
	switch identifyTermKind(f.Term) {
	case wildcard:
		f.Values = f.parseWildcardField()
		return nil
	case span:
		result, err := f.parseSpanField()
		if err != nil {
			return fmt.Errorf("could not parse %s: %w", f.Kind, err)
		}
		f.Values = result
		return nil
	case value:
		result, err := f.parseValueField()
		if err != nil {
			return fmt.Errorf("could not parse %s: %w", f.Kind, err)
		}
		f.Values = result
		return nil
	case list:
		result, err := f.parseListField()
		if err != nil {
			return fmt.Errorf("could not parse %s: %w", f.Kind, err)
		}
		f.Values = result
		return nil
	case unknown:
		return fmt.Errorf("could not parse field: %s; expression: %s", f.Kind, f.Term)
	}

	return fmt.Errorf("could not parse field: %s; expression: %s", f.Kind, f.Term)
}

func (f *Field) parseWildcardField() map[int]struct{} {
	return generateSequentialSet(f.Min, f.Max)
}

func (f *Field) parseSpanField() (map[int]struct{}, error) {
	values := strings.Split(f.Term, "-")

	min, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, fmt.Errorf("could not parse value %s: %v", values[0], err)
	}

	max, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, fmt.Errorf("could not parse value %s: %v", values[1], err)
	}

	if min >= max {
		return nil, fmt.Errorf("first value(%d) cannot be greater/equal to second(%d)", min, max)
	}

	if min < f.Min {
		return nil, fmt.Errorf("value(%d) cannot be less than min(%d)", min, f.Min)
	}

	if max > f.Max {
		return nil, fmt.Errorf("value(%d) cannot be more than max(%d)", max, f.Max)
	}

	return generateSequentialSet(min, max), nil
}

func (f *Field) parseValueField() (map[int]struct{}, error) {
	value, err := strconv.Atoi(f.Term)
	if err != nil {
		return nil, fmt.Errorf("could not parse value %s: %v", f.Term, err)
	}

	if value < f.Min {
		return nil, fmt.Errorf("value(%d) cannot be less than min(%d)", value, f.Min)
	}

	if value > f.Max {
		return nil, fmt.Errorf("value(%d) cannot be more than max(%d)", value, f.Max)
	}

	return map[int]struct{}{
		value: {},
	}, nil
}

func (f *Field) parseListField() (map[int]struct{}, error) {
	set := map[int]struct{}{}
	values := strings.Split(f.Term, ",")

	for _, rawValue := range values {
		value, err := strconv.Atoi(rawValue)
		if err != nil {
			return nil, fmt.Errorf("could not parse value %s: %v", f.Term, err)
		}

		if value < f.Min {
			return nil, fmt.Errorf("value(%d) cannot be less than min(%d)", value, f.Min)
		}

		if value > f.Max {
			return nil, fmt.Errorf("value(%d) cannot be more than max(%d)", value, f.Max)
		}

		set[value] = struct{}{}
	}

	return set, nil
}
