package internal_test

import (
	"testing"

	"github.com/xitonix/flags/internal"
)

func TestGetPrintName(t *testing.T) {
	testCases := []struct {
		title       string
		long, short string
		expected    string
	}{
		{
			title:    "empty input",
			expected: "--",
		},
		{
			title:    "long name only",
			long:     "long",
			expected: "--long",
		},
		{
			title:    "short name only",
			short:    "s",
			expected: "-s, --",
		},
		{
			title:    "long and short names",
			short:    "s",
			long:     "long",
			expected: "-s, --long",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := internal.GetPrintName(tc.long, tc.short)
			if actual != tc.expected {
				t.Errorf("Expected %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected bool
	}{
		{
			title:    "empty input",
			expected: true,
		},
		{
			title:    "white space input",
			input:    "   ",
			expected: true,
		},
		{
			title:    "tab input",
			input:    "\t",
			expected: true,
		},
		{
			title:    "tab and white space input",
			input:    "  \t   \t    ",
			expected: true,
		},
		{
			title:    "none empty",
			input:    "flags",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := internal.IsEmpty(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}

func TestSanitiseFlagID(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected string
	}{
		{
			title: "empty input",
		},
		{
			title: "white space input",
			input: "   ",
		},
		{
			title:    "lower case input",
			input:    "flag_id",
			expected: "FLAG_ID",
		},
		{
			title:    "hyphened input",
			input:    "flag-id",
			expected: "FLAG_ID",
		},
		{
			title:    "prefixed with hyphened input",
			input:    "--flag-id",
			expected: "_FLAG_ID",
		},
		{
			title:    "spaced input",
			input:    "flag id",
			expected: "FLAG_ID",
		},
		{
			title:    "multi spaced input",
			input:    "flag      id",
			expected: "FLAG_ID",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := internal.SanitiseFlagID(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %s, Actual: %s", tc.expected, actual)
			}
		})
	}
}

func TestOutOfRangeErr(t *testing.T) {
	testCases := []struct {
		title       string
		value       interface{}
		longName    string
		shortName   string
		valid       []string
		expectedErr string
	}{
		{
			title:       "list with single item",
			value:       "abc",
			longName:    "long",
			valid:       []string{"A"},
			expectedErr: "abc is not an acceptable value for --long. The expected value is A.",
		},
		{
			title:       "list with two items",
			value:       "abc",
			longName:    "long",
			valid:       []string{"A", "B"},
			expectedErr: "abc is not an acceptable value for --long. The expected values are A,B.",
		},
		{
			title:       "list with three items",
			value:       "abc",
			longName:    "long",
			valid:       []string{"A", "B", "C"},
			expectedErr: "abc is not an acceptable value for --long. The expected values are A,B,C.",
		},
		{
			title:       "empty valid range string",
			value:       "abc",
			longName:    "long",
			valid:       []string{},
			expectedErr: "abc is not an acceptable value for --long.",
		},
		{
			title:       "short named flag and a list with single item",
			value:       "abc",
			longName:    "long",
			shortName:   "S",
			valid:       []string{"A"},
			expectedErr: "abc is not an acceptable value for -S, --long. The expected value is A.",
		},
		{
			title:       "short named flag and a list with two items",
			value:       "abc",
			longName:    "long",
			shortName:   "S",
			valid:       []string{"A", "B"},
			expectedErr: "abc is not an acceptable value for -S, --long. The expected values are A,B.",
		},
		{
			title:       "short named flag and a list with three items",
			value:       "abc",
			longName:    "long",
			shortName:   "S",
			valid:       []string{"A", "B", "C"},
			expectedErr: "abc is not an acceptable value for -S, --long. The expected values are A,B,C.",
		},
		{
			title:       "short named flag and a empty valid range string",
			value:       "abc",
			longName:    "long",
			shortName:   "S",
			valid:       []string{},
			expectedErr: "abc is not an acceptable value for -S, --long.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := internal.OutOfRangeErr(tc.value, tc.longName, tc.shortName, tc.valid)
			if err.Error() != tc.expectedErr {
				t.Errorf("Expected error: '%v', Actual: '%v'", tc.expectedErr, err)
			}
		})
	}
}

func TestInvalidValueErr(t *testing.T) {
	expected := `'abc' is not a valid type value for -s, --flag`
	actual := internal.InvalidValueErr("abc", "flag", "s", "type")
	if actual.Error() != expected {
		t.Errorf("Expected: '%s', Actual: '%s'", expected, actual)
	}
}

func TestSanitiseLongName(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected string
	}{
		{
			title: "empty input",
		},
		{
			title: "white space input",
			input: "   ",
		},
		{
			title:    "lower case input",
			input:    "flag",
			expected: "flag",
		},
		{
			title:    "upper case input",
			input:    "FLAG",
			expected: "flag",
		},
		{
			title:    "spaced case input",
			input:    "   flag   ",
			expected: "flag",
		},
		{
			title:    "multi words input",
			input:    "flag name",
			expected: "flag-name",
		},
		{
			title:    "hyphened input",
			input:    "flag--name",
			expected: "flag--name",
		},
		{
			title:    "hyphened multi words input",
			input:    "flag--name   suffix",
			expected: "flag--name-suffix",
		},
		{
			title:    "prefixed with hyphen",
			input:    "--flag",
			expected: "flag",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := internal.SanitiseLongName(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %s, Actual: %s", tc.expected, actual)
			}
		})
	}
}

func TestSanitiseShortName(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected string
	}{
		{
			title: "empty input",
		},
		{
			title: "white space input",
			input: "   ",
		},
		{
			title:    "lower case input",
			input:    "flag",
			expected: "flag",
		},
		{
			title:    "upper case input",
			input:    "FLAG",
			expected: "FLAG",
		},
		{
			title:    "spaced case input",
			input:    "   flag   ",
			expected: "flag",
		},
		{
			title:    "multi words input",
			input:    "flag name",
			expected: "flag-name",
		},
		{
			title:    "multi words mixed case input",
			input:    "Flag Value",
			expected: "Flag-Value",
		},
		{
			title:    "hyphened input",
			input:    "flag--name",
			expected: "flag--name",
		},
		{
			title:    "hyphened multi words input",
			input:    "flag--name   suffix",
			expected: "flag--name-suffix",
		},
		{
			title:    "prefixed with hyphen",
			input:    "--flag",
			expected: "flag",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := internal.SanitiseShortName(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %s, Actual: %s", tc.expected, actual)
			}
		})
	}
}
