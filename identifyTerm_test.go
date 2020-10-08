package avail

import "testing"

func TestIdentifyTermType(t *testing.T) {
	tests := map[string]struct {
		input string
		want  termKind
	}{
		"span":     {"1-12", span},
		"wildcard": {"*", wildcard},
		"list":     {"1,2,3,4,5,6", list},
		"value":    {"45", value},
		"unknown":  {"233)#!", unknown},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := identifyTermKind(tc.input)
			if got != tc.want {
				t.Errorf("incorrect field type identified for %s; got %s, want %s", tc.input, got, tc.want)
			}
		})
	}
}
