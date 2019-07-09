package internal

import (
	"fmt"
	"regexp"
	"strings"
)

// IsEmpty returns true if the the trimmed version of the string is empty.
func IsEmpty(value string) bool {
	return len(strings.TrimSpace(value)) == 0
}

// SanitiseFlagID sanitises the flag keys by replacing all the hyphen and whitespace characters with "_".
func SanitiseFlagID(value string) string {
	duplicates := regexp.MustCompile(`\s+|[-]+`)
	value = duplicates.ReplaceAllString(strings.TrimSpace(value), "_")
	return strings.ToUpper(value)
}

// SanitiseLongName sanitises the input by replacing all the whitespace characters with hyphen and making it lowercase.
func SanitiseLongName(name string) string {
	return SanitiseShortName(strings.ToLower(name))
}

// SanitiseShortName sanitises the input by replacing all the whitespace characters with hyphen.
func SanitiseShortName(name string) string {
	space := regexp.MustCompile(`\s+`)
	name = space.ReplaceAllString(strings.TrimSpace(name), "-")
	return strings.TrimLeft(name, "-")
}

// OutOfRangeErr creates a new out of range error.
func OutOfRangeErr(value interface{}, longName, shortName string, valid []string) error {
	if len(valid) == 0 {
		return fmt.Errorf("%v is not an acceptable value for %s.", value, GetPrintName(longName, shortName))
	}

	plural := " is"
	if len(valid) > 1 {
		plural = "s are"
	}
	return fmt.Errorf("%v is not an acceptable value for %s. The expected value%s %s.", value, GetPrintName(longName, shortName), plural, strings.Join(valid, ","))
}

// InvalidValueErr creates a new invalid flag value error.
func InvalidValueErr(value interface{}, longName, shortName, flagType string) error {
	return fmt.Errorf("'%v' is not a valid %s value for %s", value, flagType, GetPrintName(longName, shortName))
}

// GetPrintName returns the print name of the flag (i.e "-F, --flag")
func GetPrintName(long, short string) string {
	if short != "" {
		short = "-" + short + ", "
	}
	return short + "--" + long
}
