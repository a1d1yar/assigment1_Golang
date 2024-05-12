package validator

import (
	"regexp"
	"testing"
)

func TestMatches(t *testing.T) {
	testCases := []struct {
		name   string
		value  string
		rx     *regexp.Regexp
		expect bool
	}{
		{"ValidEmail", "test@example.com", EmailRX, true},
		{"InvalidEmail", "invalid", EmailRX, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Matches(tc.value, tc.rx)
			if result != tc.expect {
				t.Errorf("Expected %v but got %v", tc.expect, result)
			}
		})
	}
}

func TestPermittedValue(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		permitted []string
		expect    bool
	}{
		{"ValidValue", "value1", []string{"value1", "value2"}, true},
		{"InvalidValue", "value3", []string{"value1", "value2"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := PermittedValue(tc.value, tc.permitted...)
			if result != tc.expect {
				t.Errorf("Expected %v but got %v", tc.expect, result)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	testCases := []struct {
		name   string
		values []string
		expect bool
	}{
		{"UniqueValues", []string{"value1", "value2", "value3"}, true},
		{"NonUniqueValues", []string{"value1", "value1", "value2"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Unique(tc.values)
			if result != tc.expect {
				t.Errorf("Expected %v but got %v", tc.expect, result)
			}
		})
	}
}

func TestUnique_EmptySlice(t *testing.T) {
	result := Unique([]string{})
	if !result {
		t.Error("Expected true for empty slice, but got false")
	}
}
