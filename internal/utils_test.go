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
			title:    "whitespace input",
			input:    "   ",
			expected: true,
		},
		{
			title:    "tab input",
			input:    "\t",
			expected: true,
		},
		{
			title:    "tab and whitespace input",
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

func TestSanitiseEnvVarName(t *testing.T) {
	testCases := []struct {
		title    string
		input    string
		expected string
	}{
		{
			title: "empty input",
		},
		{
			title: "whitespace input",
			input: "   ",
		},
		{
			title:    "lower case input",
			input:    "environment_variable",
			expected: "ENVIRONMENT_VARIABLE",
		},
		{
			title:    "hyphened input",
			input:    "environment-variable",
			expected: "ENVIRONMENT_VARIABLE",
		},
		{
			title:    "prefixed with hyphened input",
			input:    "--environment-variable",
			expected: "_ENVIRONMENT_VARIABLE",
		},
		{
			title:    "spaced input",
			input:    "environment variable",
			expected: "ENVIRONMENT_VARIABLE",
		},
		{
			title:    "multi spaced input",
			input:    "environment  variable",
			expected: "ENVIRONMENT_VARIABLE",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			actual := internal.SanitiseEnvVarName(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %s, Actual: %s", tc.expected, actual)
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
			title: "whitespace input",
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
			title: "whitespace input",
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
			input:    "Flag Name",
			expected: "Flag-Name",
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