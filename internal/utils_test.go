package internal_test

import (
	"testing"

	"go.xitonix.io/flags/internal"
)

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
		valid       string
		rangeLength int
		expectedErr string
	}{
		{
			title:       "list with single item",
			value:       "abc",
			longName:    "long",
			valid:       "A, ",
			rangeLength: 1,
			expectedErr: "abc is not an acceptable value for --long. The expected value is A.",
		},
		{
			title:       "list with more than one items",
			value:       "abc",
			longName:    "long",
			valid:       "A, B and C, ",
			rangeLength: 2,
			expectedErr: "abc is not an acceptable value for --long. The expected values are A, B and C.",
		},
		{
			title:       "empty valid range string with more than two items in the list",
			value:       "abc",
			longName:    "long",
			valid:       "",
			rangeLength: 2,
			expectedErr: "abc is not an acceptable value for --long.",
		},
		{
			title:       "empty valid range string with one item in the list",
			value:       "abc",
			longName:    "long",
			valid:       "",
			rangeLength: 1,
			expectedErr: "abc is not an acceptable value for --long.",
		},
		{
			title:       "white space valid range string with more than two items in the list",
			value:       "abc",
			longName:    "long",
			valid:       "   ",
			rangeLength: 2,
			expectedErr: "abc is not an acceptable value for --long.",
		},
		{
			title:       "white space valid range string with one item in the list",
			value:       "abc",
			longName:    "long",
			valid:       "   ",
			rangeLength: 1,
			expectedErr: "abc is not an acceptable value for --long.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			err := internal.OutOfRangeErr(tc.value, tc.longName, tc.valid, tc.rangeLength)
			if err.Error() != tc.expectedErr {
				t.Errorf("Expected error: '%v', Actual: '%v'", tc.expectedErr, err)
			}
		})
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
